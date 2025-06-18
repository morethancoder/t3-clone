package handlers

import (
	"morethancoder/t3-clone/services"
	"morethancoder/t3-clone/utils"
	"morethancoder/t3-clone/views/components"
	"net/http"
)

func POSTPingHub(w http.ResponseWriter, r *http.Request) {
	uid := r.URL.Query().Get("uid")
	if uid == "" {
		utils.Log.Error("No uid found")
		http.Error(w, "No uid found", http.StatusBadRequest)
		return
	}

	// alert user
	services.UserSSEHub.BroadcastFragments(uid,
		components.Alert(components.AlertData{
			Level:   "info",
			Message: "Pong",
		}))

}
