package handlers

import (
	"context"
	"morethancoder/t3-clone/db"
	"morethancoder/t3-clone/services"
	"morethancoder/t3-clone/utils"
	"morethancoder/t3-clone/views/components"
	"net/http"

	datastar "github.com/starfederation/datastar/sdk/go"
)

func SSEHub(w http.ResponseWriter, r *http.Request) {
	jwtCookie, err := r.Cookie("jwt")
	if err != nil {
		utils.Log.Error("No jwt cookie found")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	res, err := db.Db.AuthRefresh(jwtCookie.Value)
	if err != nil {
		utils.Log.Error(err.Error())
		http.Error(w, "Failed to auth refresh", http.StatusUnauthorized)
		return
	}

	ctx := context.WithValue(r.Context(), "AuthRecord", res.Record)
	r = r.WithContext(ctx)

	uid := res.Record.ID
	sse := datastar.NewSSE(w, r)
	services.UserSSEHub.Add(uid, sse)
	defer services.UserSSEHub.Remove(uid, sse)

	utils.Log.Debug("ssehub %+v", services.UserSSEHub)

	//send connected alert
	services.UserSSEHub.BroadcastFragments(uid,
		components.Alert(components.AlertData{
			Level:   "success",
			Message: "Connected",
		}))

	<-r.Context().Done()
}
