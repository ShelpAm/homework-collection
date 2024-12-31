package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"slices"
	"strings"
)

type SizeableReader struct {
	Reader io.Reader
	Size   int64
}

func (r *SizeableReader) Read(p []byte) (int, error) {
	return r.Reader.Read(p)
}

func (r *SizeableReader) Close() {
	r.Close()
}

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

type ProcessHomeworkResult struct {
	TaskId TaskId
}

func ProcessHomework(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*") // Allow all origins
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	// Since router only pass POST to here, we don't need to test POST method.
	// if r.Method != "POST" {
	// 	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	// 	return
	// }

	// Authenticates student's info.
	username := r.FormValue("username")
	schoolId := r.FormValue("schoolId")
	assignmentName := r.FormValue("assignmentName")

	s := Student{username, schoolId}
	if _, exists := accounts[s]; !exists {
		http.Error(w, "Student doesn't exist", http.StatusBadRequest)
		return
	}

	// Verifies assignment
	assignment, assignmentExists := assignments[assignmentName]
	if !assignmentExists {
		http.Error(w, "Specified assignment doesn't exist", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("homework")
	if err != nil {
		http.Error(w, "Failed to get file", http.StatusBadRequest)
		return
	}
	// defer file.Close() // We pass the ownwership to `s.Submit`

	filename := header.Filename
	fileSize := header.Size

	// We don't need to check if is bad file now, since we used permission system.
	// if IsBadFilename(filename) {
	// 	http.Error(w, "You received this message due to that you have uploaded suspicious file. "+
	// 		"If you have further questions, please contact the admininstrator of this server (yyx). "+
	// 		"Sorry for the inconvenience caused.",
	// 		http.StatusBadRequest)
	// 	log.Println("Bad file received:", filename)
	// 	return
	// }

	if testMode {
		http.Error(w, "In test mode, nothing actually uploaded.", http.StatusOK)
		log.Println("In TEST mode, received file")
		return
	}

	// Since this is basically schedules the execution of copying file, the file
	// will be closed after the function exits. So we need transfer the ownwership
	// of the file to `s.Submit`.
	taskId, err := s.Submit(&assignment, &SizeableReader{Reader: file, Size: fileSize}, filename)
	if err != nil {
		http.Error(w, "Failed to submit file: `"+err.Error()+"`, please contact server admin (yyx).", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(taskId)
}
