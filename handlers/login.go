package handlers

import (
	"log"
	"morethancoder/t3-clone/views/layouts"
	"morethancoder/t3-clone/views/pages"
	"net/http"
	datastar "github.com/starfederation/datastar/sdk/go"
)

func LoginGET() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dest := r.Header.Get("Sec-Fetch-Dest")

		switch dest {
		case "document":
			err := layouts.MainLayout(pages.LoginPage()).Render(r.Context(), w)
			if err != nil {
				http.Error(w, "Something went wrong", http.StatusInternalServerError)
				log.Println(err)
				return
			}
		case "empty":

			sse := datastar.NewSSE(w, r)

			err := sse.MergeFragmentTempl(pages.LoginPage())
			if err != nil {
				http.Error(w, "Something went wrong", http.StatusInternalServerError)
				log.Println(err)
				return
			}

		default:
			err := pages.LoginPage().Render(r.Context(), w)
			if err != nil {
				http.Error(w, "Something went wrong", http.StatusInternalServerError)
				log.Println(err)
				return
			}
		}

	}
}




