package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	_ "modernc.org/sqlite"
	"net/http"
)

type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	DueDateTime string `json:"due_date_time"`
}

type Response struct {
	Message string `json:"message,omitempty"`
	Tasks   []Task `json:"tasks,omitempty"`
}

var db *sql.DB

func initDB() {
	var err error
	db, err = sql.Open("sqlite", "./app.db")
	if err != nil {
		log.Fatal(err)
	}

	// SQL statement which deletes the table if it exists and creates one if it doesn't
	sqlStmt := `
	DROP TABLE IF EXISTS tasks;

	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		title TEXT,
		description TEXT,
		status TEXT,
		due_date_time TEXT
	);`

	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatalf("Error creating table: %q: %s\n", err, sqlStmt)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, title, description, status, due_date_time FROM tasks")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var tasks []Task
	defer rows.Close()
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.DueDateTime)
		if err != nil {
			log.Println("Error scanning row:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tasks = append(tasks, task)
	}

	writeTasksAndMessage(w, "Welcome to the Task Manager!", tasks)
}

func createTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		title := r.FormValue("title")
		description := r.FormValue("description")
		status := r.FormValue("status")
		dueDateTime := r.FormValue("due_date_time")

		//Validate the input to check if title, status and duedate are present, reject and return an error back to the client if not
		if title == "" || status == "" || dueDateTime == "" {
			// var err error
			// err = fmt.Errorf("Title, status and due date time are required")
			writeTasksAndMessage(w, "Title, status and due date time are required", nil)
			w.WriteHeader(http.StatusBadRequest)
			return
			// http.Error(w, err.Error(), http.StatusBadRequest)
		}

		_, err := db.Exec("INSERT INTO tasks (title, description, status, due_date_time) VALUES (?, ?, ?, ?)", title, description, status, dueDateTime)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tasks := []Task{}

		task := Task{
			Title:       title,
			Description: description,
			Status:      status,
			DueDateTime: dueDateTime,
		}

		tasks = append(tasks, task)
		message := "Task Created Successfully"
		writeTasksAndMessage(w, message, tasks)
	}
}
func writeTasksAndMessage(w http.ResponseWriter, message string, tasks []Task) {
	response := Response{
		Message: message,
		Tasks:   tasks,
	}
	responseJson, err := json.Marshal(response)
	if err != nil {
		// Write error to server console is converting to json fails
		fmt.Println("Error marshalling tasks to JSON:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJson)
}

func readAllTasksHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, title FROM tasks")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Title)
		if err != nil {
			log.Println("Error scanning row:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "Task ID: %d, Title: %s\n", task.ID, task.Title)
	}
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	row := db.QueryRow("SELECT id, title, description, status, due_date_time FROM tasks WHERE id = ?", id)
	var task Task
	err := row.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.DueDateTime)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Task not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	fmt.Fprintf(w, "Task ID: %d, Title: %s, Description: %s, Status: %s, DueDateTime: %s", task.ID, task.Title, task.Description, task.Status, task.DueDateTime)
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	_, err := db.Exec("DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeTasksAndMessage(w, fmt.Sprintf("Task ID %s deleted successfully!", id), nil)
}

// The specification only asks for the update function to update the status of the task
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	status := r.URL.Query().Get("status")

	_, err := db.Exec("UPDATE tasks SET status = ? WHERE id = ?", status, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeTasksAndMessage(w, fmt.Sprintf("Task ID %s updated status to %s successfully!", id, status), nil)
}
