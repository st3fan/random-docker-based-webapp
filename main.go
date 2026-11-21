package main

import (
	"embed"
	"log"
	"net/http"
	"os"
)

//go:embed static/index.html static/style.css
var assets embed.FS

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFileFS(w, r, assets, "static/index.html")
	})
	mux.HandleFunc("GET /style.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFileFS(w, r, assets, "static/style.css")
	})
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
