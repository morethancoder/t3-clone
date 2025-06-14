package main

import (
	"log"
	"morethancoder/t3-clone/handlers"
	"net/http"
)

func main() {
	mux := http.NewServeMux()	

	mux.Handle("GET /static/{filepath...}", handlers.StaticGET())
	
	mux.HandleFunc("GET /hotreload", handlers.HotReloadSSE)

	mux.Handle("GET /login", handlers.LoginGET())

	mux.Handle("GET /", handlers.HomeGET())

	log.Println("Server is running on port 8080")

	http.ListenAndServe(":8080", mux)
}
