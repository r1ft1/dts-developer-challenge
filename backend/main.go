package main

import (
	"fmt"
	"log"
	"net/http"
)

type Server struct {
	mux *http.ServeMux
}

// Override the default Go Mux to Attach a header to all route handlers to check if the request is coming from our frontend
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// because we can only set a single origin in the header, we will check what the request's origin is and compare it against a list, because we want the production url and development localhost url to be whitelisted
	allowedOrigins := []string{
		"http://localhost:5173",
		"http://127.0.0.1:5173",
		"https://frontend-production-f438.up.railway.app/",
		"https://frontend-production-f438.up.railway.app",
	}
	origin := r.Header.Get("Origin")
	if origin != "" {
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
				break
			}
		}
	}
	s.mux.ServeHTTP(w, r)
}

func main() {
	initDB()
	defer db.Close()

	s := &Server{mux: http.NewServeMux()}

	s.mux.HandleFunc("/create", createTaskHandler)
	s.mux.HandleFunc("/delete", deleteTaskHandler)
	s.mux.HandleFunc("/update", updateTaskHandler)
	s.mux.HandleFunc("/get", getTaskHandler)
	s.mux.HandleFunc("/", indexHandler)

	fmt.Println("Server is running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", s))
}
