package main

import (
	"html/template"
	"net/http"
	"os"
	"path/filepath"
)

type FileData struct {
	Name string
	Path string
}

func ListFiles(w http.ResponseWriter, _ *http.Request) {
	dir := "./homeworks/第四周" // Change this for listed directory

	// Read directory contents
	files, err := os.ReadDir(dir)
	if err != nil {
		http.Error(w, "Unable to read directory", http.StatusInternalServerError)
		return
	}

	// Prepare file data to pass to the template
	var fileData []FileData
	for _, file := range files {
		fileData = append(fileData, FileData{
			Name: file.Name(),
			Path: filepath.Join(dir, file.Name()),
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
		<h1>Files in Directory</h1>
		<ul>
			{{range .}}
				<li><a href="/home/homeworks/{{.Name}}">{{.Name}}</a></li>
			{{end}}
		</ul>
    <p>Totally {{len .}} files.</p>
	</body>
	</html>
	`
	t, err := template.New("filelist").Parse(tmpl)
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
