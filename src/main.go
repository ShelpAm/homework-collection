package main

import (
	"log"
	"net/http"
	"os"
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

func RedirectToHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.Redirect(w, r, "/home", http.StatusFound)
	} else {
		http.NotFound(w, r) // or let other handlers take care of it
	}
}

func ShowLogin(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "www/http/auth/login.html")
}

func ShowHome(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "www/html/index.html")
}

func ShowHomeworks(w http.ResponseWriter, r *http.Request) {
	// TODO: implement this
}

var testMode bool
var accounts = make(map[Student]struct{})
var assignments = make(map[string]Assignment)

func main() {
	testMode = len(os.Args) == 2 && os.Args[1] == "--test"

	if testMode {
		log.Println("RUN in TEST MODE")
	}

	LoadStudents(&accounts)
	LoadAssignments(&assignments)

	os.Mkdir("homeworks", 0755)

	http.Handle("/api/process-homework", RateLimit(http.HandlerFunc(ProcessHomework)))
	http.Handle("/api/export-to-zip", RateLimit(http.HandlerFunc(ExportToZip)))
	http.Handle("/auth/login", RateLimit(http.HandlerFunc(ShowLogin)))
	http.Handle("/home", http.HandlerFunc(ShowHome))
	http.Handle("/home/list-files", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { http.Redirect(w, r, "/home/homeworks", http.StatusFound) }))
	http.Handle("/home/homeworks", http.HandlerFunc(ListFiles))
	// http.Handle("/home/homeworks", http.StripPrefix("/home/homeworks", http.FileServer(http.Dir("./homeworks"))))
	http.Handle("/", http.HandlerFunc(RedirectToHome)) // Move to last line to lastly match this.

	log.Println("Server is listening on port 8080")
	http.ListenAndServe(":8080", nil)
}
