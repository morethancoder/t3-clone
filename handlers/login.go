package handlers

import (
	"context"
	"morethancoder/t3-clone/db"
	"morethancoder/t3-clone/services"
	"morethancoder/t3-clone/utils"
	"morethancoder/t3-clone/views/layouts"
	"morethancoder/t3-clone/views/pages"
	"net/http"
	"os"
)

const KeyOAuth2Providers = "OAuth2Providers"
const KeySessionID = "SessionID" 

func GETLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := r.Cookie("jwt")
		if err == nil {
			utils.Log.Info("Already authenticated" )
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		utils.Log.Debug(err.Error())

		// no jwt cookie found
		// fetching session cookie
		cookie, err := r.Cookie(KeySessionID)
		if err == nil {
			// cookie found user has session stored
			session := services.Session{
				ID: cookie.Value,
			}

			loaded := session.Load()
			if loaded {
				authProviders, ok := session.Data[KeyOAuth2Providers]
				if ok {
					ctx := context.WithValue(r.Context(), KeyOAuth2Providers, authProviders)
					r = r.WithContext(ctx)

					err = layouts.MainLayout(pages.LoginPage()).Render(r.Context(), w)
					if err != nil {
						utils.Log.Error(err.Error())
						http.Error(w, "Something went wrong", http.StatusInternalServerError)
						return
					}
					return
				}
			}
		}

		// fetch available auth methods
		authMethods, err := db.Db.GetAuthMethods()
		if err != nil {
			utils.Log.Error(err.Error())
			http.Error(w, "Failed to get auth methods", http.StatusInternalServerError)
			return
		}

		// create session
		session := services.NewSession(map[string]interface{}{
			KeyOAuth2Providers: authMethods.OAuth2.Providers,
		})

		// send session id to client as cookie
		cookie = &http.Cookie{
			Name:    KeySessionID,
			Value:   session.ID,
			Expires: session.Expiry,
			Path:    "/",
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		}

		if os.Getenv("ENV") != "dev" {
			cookie.Secure = true
		}

		http.SetCookie(w,cookie)

		// update context to update front-end
		ctx := context.WithValue(r.Context(), KeyOAuth2Providers, authMethods.OAuth2.Providers)
		r = r.WithContext(ctx)

		err = layouts.MainLayout(pages.LoginPage()).Render(r.Context(), w)
		if err != nil {
			utils.Log.Error(err.Error())
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
	}
}
