package main

import (
	"html/template"
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

func ServeHomework(w http.ResponseWriter, r *http.Request) {
	// if r.URL.String()
	http.ServeFile(w, r, filepath.Join(homeworksDir))

	dir := filepath.Join(dataDir, "./homeworks/五个一") // Change this for listed directory

	// Read directory contents
	files, err := os.ReadDir(dir)
	if err != nil {
		http.Error(w, "Unable to read directory", http.StatusInternalServerError)
		return
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
			LastModified: info.ModTime(),
		})
	}

	// Parse and execute the template
	tmpl := `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<title>File List</title>
	</head>
	<body>
		<a href="/home">Go Back to Home</a>
		<h1>Files in Directory</h1>
		<ul>
			{{range .}}
				<li><a href="/home/homeworks/{{.Name}}">{{.Name}}</a> Last Modified {{.LastModified | formatTime}}</li>
			{{end}}
		</ul>
		<p>Totally {{len .}} files.</p>
	</body>
	</html>
	`
	formatTimeF := func(t time.Time) string {
		return t.Format("on 2006-01-02 at 15:04:05")
	}
	formatTime := template.FuncMap{"formatTime": formatTimeF}
	t, err := template.New("filelist").Funcs(formatTime).Parse(tmpl)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}

	// Serve the rendered template with file data
	err = t.Execute(w, fileData)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}
