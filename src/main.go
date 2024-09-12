package main

import (
    "fmt"
    "io"
    "log"
    "net/http"
    "os"
    "path/filepath"
)

func main() {
    http.HandleFunc("/process-homework", func(w http.ResponseWriter, r *http.Request) {
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
    })

    log.Println("Server is listening on port 8080")
    http.ListenAndServe(":8080", nil)
}

