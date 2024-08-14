package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"go_final_project/cmd/repository"
	"go_final_project/cmd/task"
	"net/http"
)

type TaskHandler struct {
	repo *repository.Repository
}

func NewTaskHandler(repo *repository.Repository) *TaskHandler {
	return &TaskHandler{repo: repo}
}

func (h *TaskHandler) TaskHandler(w http.ResponseWriter, req *http.Request) {
	par := req.URL.Query().Get("id")

	var response []byte
	var ResponseStatus int
	var err error

	switch req.Method {
	case http.MethodGet:
		if par == "" {
			http.Error(w, `{"error":"неверный id"}`, http.StatusBadRequest)
			return
		}
		response, ResponseStatus, err = h.repo.TaskID(par)
	case http.MethodPost:
		response, ResponseStatus, err = h.repo.AddTask(req)
	case http.MethodPut:
		response, ResponseStatus, err = h.repo.UpdateTask(req)
	case http.MethodDelete:
		ResponseStatus, err = h.repo.DeleteTask(par)
		if err == nil {
			response, err = json.Marshal(map[string]interface{}{})
		}
	}

	if err != nil {
		http.Error(w, err.Error(), ResponseStatus)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (h *TaskHandler) TasksGet(w http.ResponseWriter, req *http.Request) {
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
		tasksParam, ResponseStatus, err := h.repo.ConditionalTask(db, par)
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

	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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

func (h *TaskHandler) TaskDone(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")

	ResponseStatus, err := h.repo.TaskDone(id)
	if err != nil {
		http.Error(w, err.Error(), ResponseStatus)
		return
	}

	response, err := json.Marshal(map[string]interface{}{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (h *TaskHandler) NextDateHandl(w http.ResponseWriter, req *http.Request) {
	param := req.URL.Query()
	now := param.Get("now")
	day := param.Get("date")
	repeat := param.Get("repeat")

	nextDay, err := h.repo.NextDate(now, day, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte(nextDay))
}
