package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

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

var (
	requestCounts = make(map[string]int)
	mutex         = &sync.Mutex{}
	dataDir       = func() string {
		if dir := os.Getenv("XDG_DATA_HOME"); dir != "" {
			return filepath.Join(dir, "homework-collection")
		}
		home, err := os.UserHomeDir()
		if err != nil {
			// As a last resort, just use the current directory
			return "."
		}
		return filepath.Join(home, ".local", "share", "homework-collection")
	}()
)

// 每分钟清除所有 IP 的访问计数
func startRateLimitReset() {
	go func() {
		for {
			time.Sleep(1 * time.Minute)
			mutex.Lock()
			requestCounts = make(map[string]int)
			mutex.Unlock()
		}
	}()
}

// Rate limiter middleware
func RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := GetClientIP(r)

		mutex.Lock()
		requestCounts[ip]++
		count := requestCounts[ip]
		mutex.Unlock()

		if false && count > 120 { // TODO: remove `false`
			http.Error(w, "Too many requests, please try again later.", http.StatusTooManyRequests)
			log.Println("IP throttled:", ip)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func RedirectToHome(w http.ResponseWriter, r *http.Request) {
	// Only redict when in /
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, "/home/", http.StatusFound)
}

func ShowLogin(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, dataDir+"/www/http/auth/login.html")
}

func ShowHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/home/" {
		http.NotFound(w, r)
		return
	}

	http.ServeFile(w, r, filepath.Join(dataDir, "www", "html", "index.html"))
}

func ShowHomeworks(w http.ResponseWriter, r *http.Request) {
	// TODO: implement this
}

func GetProgress(w http.ResponseWriter, r *http.Request) {
	taskId := r.FormValue("taskId")

	progress, err := fileUploader.GetProgress(taskId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(progress)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

type Students map[Account]*Student

var testMode bool
var students = Students{}
var assignments = make(map[string]Assignment)
var fileUploader = MakeFileUploader()

func main() {
	log.Println("Homework-collection system running...")
	testMode = len(os.Args) == 2 && os.Args[1] == "--test"

	if testMode {
		log.Println("RUN in TEST MODE")
	}

	log.Println("Loading students.")
	err := LoadStudents()
	if err != nil {
		log.Println(err.Error())
	}
	log.Println("Loading assignments.")
	err = LoadAssignments(&assignments)
	if err != nil {
		log.Println(err.Error())
	}

	os.Mkdir(filepath.Join(dataDir, "homeworks"), 0755)

	startRateLimitReset()
	http.Handle("/", http.HandlerFunc(RedirectToHome))
	http.Handle("POST /api/process-homework/", RateLimit(http.HandlerFunc(ProcessHomework)))
	http.Handle("POST /api/progress/", RateLimit(http.HandlerFunc(GetProgress)))
	http.Handle("/api/export-to-zip/", RateLimit(http.HandlerFunc(ExportToZip)))
	http.Handle("POST /auth/login/", RateLimit(http.HandlerFunc(ShowLogin)))
	http.Handle("/home/", http.HandlerFunc(ShowHome))
	http.Handle("/home/homeworks/", http.HandlerFunc(ListFiles))
	http.Handle("/home/list-files/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { http.Redirect(w, r, "/home/homeworks", http.StatusFound) }))
	// http.Handle("/home/homeworks/", http.StripPrefix("/home/homeworks", http.FileServer(http.Dir("./homeworks"))))

	log.Println("Server is listening on port 8080")
	err = http.ListenAndServe(":8080", nil)
	log.Println(err.Error())
}
