package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
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

	var body Response
	err = json.Unmarshal([]byte(got), &body)
	if err != nil {
		t.Error(err)
	}

	if body.Message != want {
		t.Errorf("got %q, want %q", body.Message, want)
	}

}

func TestReadAllTasksRoute(t *testing.T) {
	// Test that will enter multiple tasks into the DB and then will call the /read route and will check if the tasks get returned by the server

	initDB()
	defer db.Close()

	_, err := db.Exec("INSERT INTO tasks (title, description, status, due_date_time) VALUES ('Task0', 'TaskDesc0', 'Status0', 'Date0');")
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec("INSERT INTO tasks (title, description, status, due_date_time) VALUES ('Task1','TaskDesc1', 'Status1', 'Date1');")
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec("INSERT INTO tasks (title, description, status, due_date_time) VALUES ('Task2', 'TaskDesc2', 'Status2', 'Date2');")
	if err != nil {
		t.Fatal(err)
	}

	req, _ := http.NewRequest("GET", "/read", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(readAllTasksHandler)

	want := "Task ID: 1, Title: Task0\nTask ID: 2, Title: Task1\nTask ID: 3, Title: Task2\n"

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("got %v, want %v", rr.Code, http.StatusOK)
	}
	got := rr.Body.String()
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestDeleteTaskRoute(t *testing.T) {
	//Test will take a task's ID through a /delete route. We will check if the route handler has deleted the task successfully by checking the database to see if that record is still present
	initDB()
	defer db.Close()

	_, err := db.Exec("INSERT INTO tasks (title, description, status, due_date_time) VALUES ('TaskToDelete', 'TaskDeleteDesc', 'Status3', 'Date3');")
	if err != nil {
		t.Fatal(err)
	}

	// show that in the database there is an inserted task, making the count 1
	var statementOutput string
	response := db.QueryRow("SELECT COUNT(*) FROM tasks;")
	err = response.Scan(&statementOutput)
	if err != nil {
		t.Fatal(err)
	}

	//convert the string to an int
	i, err := strconv.Atoi(statementOutput)
	if err != nil {
		t.Fatal(err)
	}

	if i != 1 {
		t.Errorf("got %q, want %q", i, 1)
	}

	req, _ := http.NewRequest("GET", "/delete?id=1", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(deleteTaskHandler)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("got %v, want %v", rr.Code, http.StatusOK)
	}
	// Got will be the output of a SQL Select statement
	want := "Task ID 1 deleted successfully!"
	got := rr.Body.String()

	var body Response
	err = json.Unmarshal([]byte(got), &body)
	if err != nil {
		t.Error(err)
	}

	if body.Message != want {
		t.Errorf("got %q, want %q", body.Message, want)
	}

	response = db.QueryRow("SELECT COUNT(*) FROM tasks;")
	err = response.Scan(&statementOutput)
	//check if the sql statement returns 0 for count matching that ID
	if err != nil {
		t.Fatal(err)
	}

	//convert the string to an int
	i, err = strconv.Atoi(statementOutput)
	if err != nil {
		t.Fatal(err)
	}

	if i != 0 {
		t.Errorf("got %q, want %q", i, 0)
	}
}

func TestGetTaskRoute(t *testing.T) {
	//Test will insert a task into the Database directly and then make a request to the API's /get route, supplying the ID as a query param. We will check if the route handler returns the task associated with the supplied ID successfully by checking if the returned task is the same as the one we inserted
	initDB()
	defer db.Close()

	_, err := db.Exec("INSERT INTO tasks (title, description, status, due_date_time) VALUES ('TaskToGet', 'TaskGetDesc', 'Status3', 'Date3');")
	if err != nil {
		t.Fatal(err)
	}

	req, _ := http.NewRequest("GET", "/get?id=1", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getTaskHandler)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("got %v, want %v", rr.Code, http.StatusOK)
	}
	// Got will be the output of a SQL Select statement
	got := rr.Body.String()
	want := "Task ID: 1, Title: TaskToGet, Description: TaskGetDesc, Status: Status3, DueDateTime: Date3"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestCreateTaskRoute(t *testing.T) {
	// Test that will input a task into the Database through the /create route, will prepare a form to append to a request. To check we will read all the rows from the database and confirm that the new task is present.
	initDB()
	defer db.Close()

	want := "Test Task"

	// In testing don't need to catch the error from NewRequest as we know our route /create will correctly resolve in the test server
	req, _ := http.NewRequest("POST", "/create", nil)

	testTask := Task{
		Title:       "Test Task",
		Description: "This is a test task",
		Status:      "Pending",
		DueDateTime: "2023-10-01T10:00:00Z",
	}
	// req.Form takes a url.Values type so we'll convert our test data
	formData := url.Values{
		"title":         {testTask.Title},
		"description":   {testTask.Description},
		"status":        {testTask.Status},
		"due_date_time": {testTask.DueDateTime},
	}
	req.Form = formData

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(createTaskHandler)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("got %v, want %v", rr.Code, http.StatusOK)
	}

	var statementOutput string

	response := db.QueryRow("SELECT title FROM tasks WHERE title = ?;", want)
	err := response.Scan(&statementOutput)
	if err != nil {
		t.Fatal(err)
	}

	// Want the sql statement output to be the title of the task we just inserted
	if statementOutput != "Test Task" {
		t.Errorf("got %q, want %q", statementOutput, want)
	}
}

func TestUpdateTask(t *testing.T) {
	// Test that will first insert a task into the database, and then through request to the /update route handler will update a task. And then check if the update was successful by comparing the output of an sql select statement with the expected output
	initDB()
	defer db.Close()

	_, err := db.Exec("INSERT INTO tasks (title, description, status, due_date_time) VALUES ('TaskToUpdate', 'TaskUpdateDesc', 'Not Started', 'Date');")
	if err != nil {
		t.Fatal(err)
	}

	req, _ := http.NewRequest("PATCH", "/update?id=1&status=Completed", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(updateTaskHandler)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("got %v, want %v", rr.Code, http.StatusOK)
	}
	// Got will be the output of a SQL Select statement
	want := "Completed"

	var statementOutput string

	response := db.QueryRow("SELECT status FROM tasks WHERE id = ?;", 1)
	err = response.Scan(&statementOutput)
	if err != nil {
		t.Fatal(err)
	}

	if statementOutput != want {
		t.Errorf("got %q, want %q", statementOutput, want)
	}
}
