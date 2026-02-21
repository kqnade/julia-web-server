package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"

	"github.com/kqnade/julia-web-server/internal/handler"
)

//go:embed web
var webFS embed.FS

func main() {
	webContent, err := fs.Sub(webFS, "web")
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	// Serve index.html at /satori/julia
	mux.HandleFunc("GET /satori/julia", func(w http.ResponseWriter, r *http.Request) {
		data, err := fs.ReadFile(webContent, "index.html")
		if err != nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(data)
	})

	// Serve app.js
	mux.HandleFunc("GET /satori/julia/app.js", func(w http.ResponseWriter, r *http.Request) {
		data, err := fs.ReadFile(webContent, "app.js")
		if err != nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
		w.Write(data)
	})

	// Julia set computation API
	mux.HandleFunc("GET /satori/julia/api", handler.JuliaAPI)

	addr := ":8080"
	fmt.Printf("Julia Set server listening on http://localhost%s/satori/julia\n", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
