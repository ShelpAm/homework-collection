package main

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

func MakeZip(out string, dir string) error {
	zipFile, err := os.Create(out)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			filename := filepath.Base(path)

			writer, err := zipWriter.Create(filename)
			if err != nil {
				return err
			}

			_, err = io.Copy(writer, file)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func IsBadFilename(filename string) bool {
	if strings.ContainsAny(filename, "/\\") {
		return true
	}

	ext := filepath.Ext(filename)
	allowed := []string{".zip", ".7z", ".gz", ".rar", ".xz"}
	contains := slices.IndexFunc(allowed, func(e string) bool {
		return e == ext
	}) != -1

	if !contains {
		return true
	}

	return false
}

func main() {
	os.Mkdir("homeworks", 0755)

	http.HandleFunc("/api/process-homework", func(w http.ResponseWriter, r *http.Request) {
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
			http.Error(w, "Don't attack my server plz", http.StatusInternalServerError)
			return
		}

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

		fmt.Fprintf(w, "Homework submitted successfully")
		log.Println("Received file", filename)
	})

	http.HandleFunc("/api/export-to-zip", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Invalid request method, please use GET", http.StatusBadRequest)
			return
		}

		os.Mkdir("zip", 0755)
		zipPath := "zip/exported.zip"
		os.Remove(zipPath)

		err := MakeZip(zipPath, "homeworks")
		if err != nil {
			http.Error(w, "Cannot make zip file, 请联系杨扬骁。", http.StatusInternalServerError)
			log.Println("Cannot make zip file")
		}

		w.Header().Set("Content-Disposition", "attachment; filename=exported.zip")
		w.Header().Set("Content-Type", "application/zip")
		http.ServeFile(w, r, zipPath)

		log.Println("Zip was exported successfully")
	})

	log.Println("Server is listening on port 8080")
	http.ListenAndServe(":8080", nil)
}
