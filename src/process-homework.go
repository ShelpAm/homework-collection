package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

func IsBadFilename(filename string) bool {
	if strings.ContainsAny(filename, "/\\") {
		return true
	}

	ext := filepath.Ext(filename)
	allowed := []string{".zip", ".7z", ".gz", ".rar", ".xz"}
	if !slices.Contains(allowed, ext) {
		return true
	}

	if !strings.Contains(filename, "-") {
		return true
	}

	return false
}

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func Respond(w http.ResponseWriter, r Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // HTTP 200 OK
	json.NewEncoder(w).Encode(r)
}

func ProcessHomework(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*") // Allow all origins
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("homework")
	if err != nil {
		http.Error(w, "Failed to get file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	filename := header.Filename
	filepath := filepath.Join("homeworks", filename)

	if IsBadFilename(filename) {
		http.Error(w, "You received this message due to that you have uploaded suspicious file. "+
			"If you have further questions, please contact the admininstrator of this server (yyx). "+
			"Sorry for the inconvenience caused.",
			http.StatusBadRequest)
		log.Println("Bad file received:", filename)
		return
	}

	if !testMode {
		f, err := os.Create(filepath)
		if err != nil {
			http.Error(w, "Failed to create file", http.StatusInternalServerError)
			return
		}
		defer f.Close()

		_, err = io.Copy(f, file)
		if err != nil {
			http.Error(w, "Failed to copy file", http.StatusInternalServerError)
			return
		}
	}

	if testMode {
		log.Print("TEST MODE: ")
	}

	fmt.Fprintf(w, "Homework submitted successfully")
	log.Println("Received file", filename)
}
