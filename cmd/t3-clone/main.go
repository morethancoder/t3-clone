package main

import (
	"morethancoder/t3-clone/handlers"
	"morethancoder/t3-clone/services"
	"morethancoder/t3-clone/utils"
	"net/http"
	"os"
)

func main() {
	if os.Getenv("ENV") == "dev" {
		utils.Log.DebugMode = true
	}

	go services.LoopAndCleanSessionStore()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /sse", handlers.SSEHub)
	mux.HandleFunc("GET /chat/new", handlers.GETNewChat)
	mux.HandleFunc("POST /chat", handlers.POSTChat)

	mux.HandleFunc("GET /login", handlers.GETLogin())
	mux.HandleFunc("GET /sign-out", handlers.GETSignOut())

	mux.HandleFunc("POST /auth-redirect", handlers.POSTAuthRedirect())
	mux.HandleFunc("GET /auth-redirect", handlers.GETAuthRedirect())

	if os.Getenv("ENV") == "dev" {
		mux.HandleFunc("GET /hotreload", handlers.SSEHotreload)
	}

	mux.Handle("GET /static/{filepath...}", handlers.GETStatic())
	mux.HandleFunc("GET /", handlers.GETHome())

	utils.Log.Debug("Listening on port 8080")

	http.ListenAndServe(":8080", mux)
}
