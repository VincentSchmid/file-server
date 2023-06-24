package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

var ApiKey string

func loadAPIKey() error {
	godotenv.Load(".env")

	ApiKey = os.Getenv("API_KEY")
	return nil
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("X-API-Key") != ApiKey {
		http.Error(w, "Invalid API Key", http.StatusForbidden)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	filePath := filepath.Join("uploads", handler.Filename)
	out, err := os.Create(filePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Successfully uploaded file\n")
}

func main() {
	err := loadAPIKey()
	if err != nil {
		fmt.Println("Failed to load API key:", err)
		return
	}

	http.HandleFunc("/upload", uploadFile)

	// Serve files in the uploads folder
	http.Handle("/files/", http.StripPrefix("/files/", http.FileServer(http.Dir("uploads"))))

	http.ListenAndServe(":8080", nil)
}
