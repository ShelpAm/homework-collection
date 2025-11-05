package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type FileData struct {
	Name         string
	Path         string
	LastModified time.Time
}

// params:
//
//	AssignmentName: string
func ListSubmissions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	assignmentName := r.URL.Query().Get("AssignmentName")
	dir := filepath.Join(dataDir, "homeworks", assignmentName)

	// Read directory contents
	files, err := os.ReadDir(dir)
	if err != nil {
		// No directory is regarded as no submissions.
		json.NewEncoder(w).Encode([]FileData{})
		return
		// http.Error(w, "Unable to read directory "+dir, http.StatusInternalServerError)
		// return
	}

	// Prepare file data to pass to the template
	var fileData []FileData
	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			log.Println("An error occurred: ", err)
		}
		fileData = append(fileData, FileData{
			Name:         file.Name(),
			Path:         filepath.Join(dir, file.Name()),
			LastModified: info.ModTime().Truncate(time.Second),
		})
	}

	json.NewEncoder(w).Encode(fileData)

	// // Parse and execute the template
	// tmpl := `
	// <!DOCTYPE html>
	// <html lang="en">
	// <head>
	// 	<meta charset="UTF-8">
	// 	<title>File List</title>
	// </head>
	// <body>
	// 	<a href="/home">Go Back to Home</a>
	// 	<h1>Files in Directory</h1>
	// 	<ul>
	// 		{{range .}}
	// 			<li><a href="/home/homeworks/{{.Name}}">{{.Name}}</a> Last Modified {{.LastModified | formatTime}}</li>
	// 		{{end}}
	// 	</ul>
	// 	<p>Totally {{len .}} files.</p>
	// </body>
	// </html>
	// `
	// formatTimeF := func(t time.Time) string {
	// 	return t.Format("on 2006-01-02 at 15:04:05")
	// }
	// formatTime := template.FuncMap{"formatTime": formatTimeF}
	// t, err := template.New("filelist").Funcs(formatTime).Parse(tmpl)
	// if err != nil {
	// 	http.Error(w, "Error rendering template", http.StatusInternalServerError)
	// 	return
	// }
	//
	// // Serve the rendered template with file data
	// err = t.Execute(w, fileData)
	// if err != nil {
	// 	http.Error(w, "Error rendering template", http.StatusInternalServerError)
	// }
}
