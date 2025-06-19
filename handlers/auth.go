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
		err := layouts.MainLayout(pages.AuthRedirect("Authenticating...")).Render(r.Context(), w)
		if err != nil {
			utils.Log.Error(err.Error())
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
	}
}

func POSTAuthRedirect() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := r.Cookie("jwt")
		if err == nil {
			utils.Log.Info("Already authenticated")
		sse := datastar.NewSSE(w, r)
			sse.ExecuteScript("window.location.href = '/'")
			return
		}
		//get user session
		cookie, err := r.Cookie(KeySessionID)
		if err != nil {
			utils.Log.Error(err.Error())
		sse := datastar.NewSSE(w, r)
			sse.ExecuteScript("window.location.href = '/login'")
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
		sse := datastar.NewSSE(w, r)
			sse.ExecuteScript("window.location.href = '/login'")
			return
		}

		//get auth providers
		authProviders, ok := session.Data[KeyOAuth2Providers].([]db.OAuth2Provider)
		if !ok {
			utils.Log.Error("Auth providers not found")
		sse := datastar.NewSSE(w, r)
			sse.ExecuteScript("window.location.href = '/login'")
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
					utils.Log.Error("Failed to auth with %s: %s", provider.Name, err.Error())
		sse := datastar.NewSSE(w, r)
					sse.ExecuteScript("window.location.href = '/login'")
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

				modelrecords, err := db.Db.GetModelRecords(map[db.QueryParam]string{})
				if err != nil {
					utils.Log.Error(err.Error())
					http.Error(w, "Failed to get models", http.StatusInternalServerError)
					return
				}

				ctx := context.WithValue(r.Context(), "AuthRecord", res.Record)
				ctx = context.WithValue(ctx, "Models", db.GroupModelRecordsByCompany(modelrecords.Items))
				r = r.WithContext(ctx)

				http.SetCookie(w, jwtCookie)

				sse := datastar.NewSSE(w, r)
				//redirect to home
				err = sse.MergeFragmentTempl(pages.HomePage())
				if err != nil {
					utils.Log.Debug("Failed to merge fragment: %s", err.Error())
					sse.ExecuteScript("window.location.href = '/'")
				}
				return

			}
		}

		utils.Log.Error("State not found")
		sse := datastar.NewSSE(w, r)
		sse.ExecuteScript("window.location.href = '/login'")
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
