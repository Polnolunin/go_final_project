package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"go_final_project/cmd/task"
	"net/http"
	"time"

	_ "modernc.org/sqlite"
)

func TasksGet(w http.ResponseWriter, req *http.Request) {
	tasks := make(map[string][]task.Task)

	par := req.URL.Query().Get("search")

	db, err := sql.Open("sqlite", task.FileDB)
	if err != nil {
		fmt.Println("open db")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	if par != "" {
		var tasksParam []task.Task
		tasksParam, ResponseStatus, err := ConditionalTask(db, par)
		if err != nil {
			http.Error(w, err.Error(), ResponseStatus)
			return
		}
		tasks["tasks"] = tasksParam

		if tasks["tasks"] == nil {
			tasks["tasks"] = []task.Task{}
		}

		response, err := json.Marshal(tasks)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
		return
	}

	rows, err := db.Query(`SELECT id, date, title, comment, repeat FROM scheduler 
	ORDER BY date LIMIT 20`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		task := task.Task{}

		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := rows.Err(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tasks["tasks"] = append(tasks["tasks"], task)
	}

	if tasks["tasks"] == nil {
		tasks["tasks"] = []task.Task{}
	}

	response, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func ConditionalTask(db *sql.DB, search string) ([]task.Task, int, error) {
	var tasks []task.Task

	var date bool
	timeSearch, err := time.Parse("02.01.2006", search)
	if err == nil {
		date = true
	}

	var rows *sql.Rows

	if date {
		dateFormat := timeSearch.Format("20060102")
		rows, err = db.Query(`SELECT id, date, title, comment, repeat FROM scheduler
        WHERE date = :date LIMIT 20`,
			sql.Named("date", dateFormat))
	} else {
		rows, err = db.Query(`SELECT id, date, title, comment, repeat FROM scheduler
    WHERE title LIKE :search OR comment LIKE :search ORDER BY date LIMIT 20`,
			sql.Named("search", "%"+search+"%"))
	}

	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	defer rows.Close()

	for rows.Next() {
		task := task.Task{}
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return tasks, http.StatusOK, nil
}
