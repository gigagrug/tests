package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/go-playground/validator/v10"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", Home)
	mux.HandleFunc("POST /createBlog/{$}", createBlog)
	mux.HandleFunc("POST /upload/{$}", uploadHandler)

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

const MAX_UPLOAD_SIZE = 1024 * 1024 // 1MB
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)
	if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
		http.Error(w, "The uploaded file is too big. Please choose an file that's less than 1MB in size", http.StatusBadRequest)
		return
	}
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer file.Close()

	// Create the uploads folder if it doesn't
	// already exist
	err = os.MkdirAll("./uploads", os.ModePerm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create a new file in the uploads directory
	dst, err := os.Create(fmt.Sprintf("./uploads/%d%s", time.Now().UnixNano(), filepath.Ext(fileHeader.Filename)))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer dst.Close()

	// Copy the uploaded file to the filesystem
	// at the specified destination
	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Upload successful")
}

type Blog struct {
	Id      int    `json:"id"`
	Title   string `json:"title" validate:"required,min=3,max=5"`
	Article string `json:"article" validate:"required"`
}

func createBlog(w http.ResponseWriter, r *http.Request) {
	var blog Blog
	json.NewDecoder(r.Body).Decode(&blog)
	validate := validator.New()
	err := validate.Struct(&blog)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		http.Error(w, fmt.Sprintf("Validation error: %s", errors), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(blog)
}
func Home(w http.ResponseWriter, r *http.Request) {
	blogs := []Blog{
		{Id: 1, Title: "so", Article: "ok"},
		{Id: 2, Title: "ao", Article: "ak"},
	}

	jsonData, err := json.Marshal(blogs)
	if err != nil {
		http.Error(w, "Unable to marshal JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(jsonData)
}
