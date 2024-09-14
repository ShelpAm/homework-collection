package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func ProcessHomework(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, "You received this message due to that you have uploaded suspicious file. Don't attack my server plz", http.StatusInternalServerError)
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
