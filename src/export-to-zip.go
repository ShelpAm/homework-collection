package main

import (
	"archive/zip"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// Define the FileData struct
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

func ExportToZip(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Invalid request method, please use GET", http.StatusBadRequest)
		return
	}

	os.Mkdir("zip", 0755)
	zipPath := "zip/exported.zip"
	os.Remove(zipPath)

	err := MakeZip(zipPath, "homeworks")
	if err != nil {
		http.Error(w, "Cannot make zip file, 请联系服务器管理员。", http.StatusInternalServerError)
		log.Println("Cannot make zip file")
	}

	w.Header().Set("Content-Disposition", "attachment; filename=exported.zip")
	w.Header().Set("Content-Type", "application/zip")
	http.ServeFile(w, r, zipPath)

	log.Println("Zip was exported successfully")
}
