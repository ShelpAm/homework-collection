package main

import (
	"archive/zip"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
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
			// Compute the relative path for the zip entry
			relPath, err := filepath.Rel(dir, path)
			if err != nil {
				return err
			}

			// Open the file
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			// Create zip entry with relative path
			writer, err := zipWriter.Create(relPath)
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

func ExportToZip(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Invalid request method, please use GET", http.StatusBadRequest)
		return
	}

	tmpDir := filepath.Join(dataDir, "zip")
	os.Mkdir(tmpDir, 0755)
	zipFile := filepath.Join(tmpDir, "exported.zip")
	os.Remove(zipFile)

	if err := MakeZip(zipFile, homeworksDir); err != nil {
		http.Error(w, "Failed to make zip file, 请联系服务器管理员。"+err.Error(), http.StatusInternalServerError)
		log.Println("Failed to make zip file: " + err.Error())
	}

	w.Header().Set("Content-Disposition", "attachment; filename=exported.zip")
	w.Header().Set("Content-Type", "application/zip")
	http.ServeFile(w, r, zipFile)

	log.Println("Zip was exported successfully")
}
