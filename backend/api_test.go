package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
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

func TestReadAllTasksRoute(t *testing.T) {
	// Test that will enter multiple tasks into the DB and then will call the /read route and will check if the tasks get returned by the server

	initDB()
	defer db.Close()

	_, err := db.Exec("INSERT INTO tasks VALUES (0, 'Task0', 'TaskDesc0', 'Status0', 'Date0');")
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec("INSERT INTO tasks VALUES (1, 'Task1','TaskDesc1', 'Status1', 'Date1');")
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec("INSERT INTO tasks VALUES (2, 'Task2', 'TaskDesc2', 'Status2', 'Date2');")
	if err != nil {
		t.Fatal(err)
	}

	req, _ := http.NewRequest("GET", "/read", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(readAllTasksHandler)

	want := "Task ID: 0, Title: Task0\nTask ID: 1, Title: Task1\nTask ID: 2, Title: Task2\n"

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
	initDB()
	defer db.Close()

	want := "Task created successfully!"
	// In testing don't need to catch the error from NewRequest as we know our route /create will correctly resolve in the test server
	req, _ := http.NewRequest("POST", "/create", nil)

	testTast := Task{
		Title:       "Test Task",
		Description: "This is a test task",
		Status:      "Pending",
		DueDateTime: "2023-10-01T10:00:00Z",
	}
	// req.Form takes a url.Values type so we'll convert our test data
	formData := url.Values{
		"title":         {testTast.Title},
		"description":   {testTast.Description},
		"status":        {testTast.Status},
		"due_date_time": {testTast.DueDateTime},
	}
	req.Form = formData

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(createTaskHandler)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("got %v, want %v", rr.Code, http.StatusOK)
	}
	got := rr.Body.String()
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
