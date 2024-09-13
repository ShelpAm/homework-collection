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
	"sync"
	"time"
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

func GetClientIP(r *http.Request) string {
    // Check if the request has the X-Forwarded-For header
    forwarded := r.Header.Get("X-Forwarded-For")
    if forwarded != "" {
        // X-Forwarded-For may contain multiple IPs, the first one is the original client IP
        return strings.Split(forwarded, ",")[0]
    }

    // If X-Forwarded-For is not set, fall back to RemoteAddr
    ip := r.RemoteAddr
    ip = strings.Split(ip, ":")[0] // Extract the IP without the port
    return ip
}

var requestCounts = make(map[string]int)
var mutex = &sync.Mutex{}

// Rate limiter middleware
func RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := GetClientIP(r)

		mutex.Lock()

		// Increment request count for the IP
		requestCounts[ip]++

		// Set a time window to reset the counter every minute
		go func() {
			time.Sleep(1 * time.Minute)
			mutex.Lock()
			delete(requestCounts, ip) // Reset request count after 1 minute
			mutex.Unlock()
		}()

		// Check if the IP exceeds the limit (e.g., 60 requests per minute)
		if requestCounts[ip] > 10 {
			mutex.Unlock()
			http.Error(w, "Too many requests, please try again later.", http.StatusTooManyRequests)
			log.Println("IP banned due to excessive requests:", ip)
			return
		}

		mutex.Unlock()

		// Proceed to the next handler
		next.ServeHTTP(w, r)
	})
}

func main() {
	testMode := len(os.Args) == 2 && os.Args[1] == "--test"

	if testMode {
		log.Println("RUN in TEST MODE")
	}

	os.Mkdir("homeworks", 0755)

	http.Handle("/api/process-homework", RateLimit(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	})))

	http.Handle("/api/export-to-zip", RateLimit(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	})))

	log.Println("Server is listening on port 8080")
	http.ListenAndServe(":8080", nil)
}
