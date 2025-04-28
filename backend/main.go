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
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
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
