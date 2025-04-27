package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
)

type Task struct {
	ID          int
	Title       string
	Description string
	Status      string
	DueDateTime string
}

var db *sql.DB

func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "./app.db")
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
	fmt.Fprint(w, "Welcome to the Task Manager!")
}

func createTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		title := r.FormValue("title")
		description := r.FormValue("description")
		status := r.FormValue("status")
		dueDateTime := r.FormValue("due_date_time")

		_, err := db.Exec("INSERT INTO tasks (title, description, status, due_date_time) VALUES (?, ?, ?, ?)", title, description, status, dueDateTime)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.ResponseWriter.Write(w, []byte("Task created successfully!"))
		http.ResponseWriter.WriteHeader(w, http.StatusOK)
	}
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
	fmt.Fprintf(w, "Task with ID %s deleted successfully!", id)
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	status := r.URL.Query().Get("status")

	_, err := db.Exec("UPDATE tasks SET status = ? WHERE id = ?", status, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Task with ID %s updated status to %s successfully!", id, status)
}
