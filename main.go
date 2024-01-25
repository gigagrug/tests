package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", Home)
	err := http.ListenAndServe(":8000", addCORS(mux))
	if err != nil {
		log.Fatal(err)
	}
}
func addCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		h.ServeHTTP(w, r)
	})
}

type Blog struct {
	ID      string
	Title   string
	Content string
}

func Home(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	blogs := []Blog{
		{ID: "1", Title: "so", Content: "ok"},
		{ID: "2", Title: "ao", Content: "ak"},
	}

	jsonData, err := json.Marshal(blogs)
	if err != nil {
		http.Error(w, "Unable to marshal JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(jsonData)
}
