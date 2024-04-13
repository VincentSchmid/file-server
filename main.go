package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

var (
	ApiKey string
	logger *zap.Logger
)

func initLogger() {
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		panic(err)
	}
	logger.Info("Logger initialized")
}

func loadAPIKey() error {
	godotenv.Load(".env")

	ApiKey = os.Getenv("API_KEY")
	if ApiKey == "" {
		logger.Error("API key is not set")
		return fmt.Errorf("API key is not set")
	}
	logger.Info("API key loaded")
	return nil
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		logger.Error("Invalid request method", zap.String("method", r.Method))
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("X-API-Key") != ApiKey {
		logger.Warn("Invalid API Key provided", zap.String("provided_api_key", r.Header.Get("X-API-Key")))
		http.Error(w, "Invalid API Key", http.StatusForbidden)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		logger.Error("Error retrieving file from form", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	filePath := filepath.Join("uploads", handler.Filename)
	out, err := os.Create(filePath)
	if err != nil {
		logger.Error("Error creating file", zap.String("file_path", filePath), zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		logger.Error("Error copying file", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.Info("Successfully uploaded file", zap.String("file_path", filePath))
	fmt.Fprintf(w, "Successfully uploaded file\n")
}

func main() {
	initLogger()

	err := loadAPIKey()
	if err != nil {
		logger.Fatal("Failed to load API key", zap.Error(err))
		return
	}

	http.HandleFunc("/upload", uploadFile)
	http.Handle("/files/", http.StripPrefix("/files/", http.FileServer(http.Dir("uploads"))))

	logger.Info("Server starting on port 8080")
	http.ListenAndServe(":8080", nil)
}
