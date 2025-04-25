package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIndexRoute(t *testing.T) {
	want := "Welcome to the Task Manager!"
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(indexHandler)

	initDB()
	defer db.Close()

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("got %v, want %v", rr.Code, http.StatusOK)
	}
	got := rr.Body.String()
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestCreateTaskRoute(t *testing.T) {
	want := "Task created successfully!"
	req, err := http.NewRequest("POST", "/create", nil)
	testTast := Task{
		Title:       "Test Task",
		Description: "This is a test task",
		Status:      "Pending",
		DueDateTime: "2023-10-01T10:00:00Z",
	}
	// Simulate the form data
	req.Form = make(map[string][]string)
	req.Form["title"] = []string{testTast.Title}
	req.Form["description"] = []string{testTast.Description}
	req.Form["status"] = []string{testTast.Status}
	req.Form["due_date_time"] = []string{testTast.DueDateTime}

	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(createTaskHandler)

	initDB()
	defer db.Close()

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("got %v, want %v", rr.Code, http.StatusOK)
	}
	got := rr.Body.String()
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
