package handlers

import (
	"context"
	"morethancoder/t3-clone/db"
	"morethancoder/t3-clone/utils"
	"morethancoder/t3-clone/views/layouts"
	"morethancoder/t3-clone/views/pages"
	"net/http"
)

func GETHome() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jwtCookie, err := r.Cookie("jwt")
		if err != nil {
			utils.Log.Error("No jwt cookie found")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		res, err := db.Db.AuthRefresh(jwtCookie.Value)
		if err != nil {
			utils.Log.Error(err.Error())
			http.SetCookie(w, &http.Cookie{
				Name:   "jwt",
				MaxAge: -1,
				HttpOnly: true,
			})
			http.Error(w, "Failed to auth refresh", http.StatusUnauthorized)
			return
		}

		getModelsRes, err := db.Db.GetModelRecords(map[db.QueryParam]string{})
		if err != nil {
			utils.Log.Error(err.Error())
			http.Error(w, "Failed to get models", http.StatusInternalServerError)
			return
		}

		mappedModels := db.GroupModelRecordsByCompany(getModelsRes.Items)
		ctx := context.WithValue(r.Context(), "Models", mappedModels)
		ctx = context.WithValue(ctx, "AuthRecord", res.Record)
		r = r.WithContext(ctx)

		err = layouts.MainLayout(pages.HomePage()).Render(r.Context(), w)
		if err != nil {
			utils.Log.Error(err.Error())
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}

	}
}
