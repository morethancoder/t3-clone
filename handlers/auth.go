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

	datastar "github.com/starfederation/datastar/sdk/go"
)

func GETAuthRedirect() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := r.Cookie("jwt")
		if err == nil {
			utils.Log.Info("Already authenticated")
			http.Redirect(w, r, "/", http.StatusMovedPermanently)
			return
		}

		err = layouts.MainLayout(pages.AuthRedirect("Authenticating...")).Render(r.Context(), w)
		if err != nil {
			utils.Log.Error(err.Error())
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
	}
}

func POSTAuthRedirect() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//get user session
		cookie, err := r.Cookie(KeySessionID)
		if err != nil {
			utils.Log.Error(err.Error())
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		//get search params
		state := r.URL.Query().Get("state")
		code := r.URL.Query().Get("code")

		session := services.Session{
			ID: cookie.Value,
		}

		ok := session.Load()
		if !ok {
			utils.Log.Error("Session not found")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		//get auth providers
		authProviders, ok := session.Data[KeyOAuth2Providers].([]db.OAuth2Provider)
		if !ok {
			utils.Log.Error("Auth providers not found")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		for _, provider := range authProviders {
			if provider.State == state {
				//authenticate
				utils.Log.Debug("Authenticating with %s code: %s", provider.Name, code)

				res, err := db.Db.AuthWithOAuth2(db.OAuth2AuthRequest{
					Provider:     provider.Name,
					Code:         code,
					CodeVerifier: provider.CodeVerifier,
					RedirectURL:  os.Getenv("APP_URL") + "/auth-redirect",
				}, map[db.QueryParam]string{})

				if err != nil {
					utils.Log.Error(err.Error())
					http.Error(w, "Something went wrong", http.StatusInternalServerError)
					return
				}

				// save jwt as cookie to auto auth
				jwtCookie := &http.Cookie{
					Name:     "jwt",
					Value:    res.Token,
					Path:     "/",
					HttpOnly: true,
					SameSite: http.SameSiteStrictMode,
				}

				if os.Getenv("ENV") != "dev" {
					jwtCookie.Secure = true
				}

				http.SetCookie(w, jwtCookie)


				ctx := context.WithValue(r.Context(), "AuthRecord", res.Record)
				r = r.WithContext(ctx)

				//redirect to home
				sse := datastar.NewSSE(w, r)

				err = sse.MergeFragmentTempl(pages.HomePage())
				if err != nil {
					utils.Log.Error(err.Error())
					http.Error(w, "Something went wrong", http.StatusInternalServerError)
				}
				return

			}
		}

		utils.Log.Error("State not found")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		//maybe render an update in the frontend
		return
	}
}

func GETSignOut() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:     "jwt",
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
		})
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
}
