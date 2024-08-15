package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"go_final_project/cmd/repository"
	"go_final_project/cmd/task"
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
	var err error

	switch req.Method {
	case http.MethodGet:
		if par == "" {
			http.Error(w, `{"error":"неверный id"}`, http.StatusBadRequest)
			return
		}
		task, err := h.repo.TaskID(par)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w, `{"error":"задача не найдена"}`, http.StatusNotFound)
			} else {
				http.Error(w, fmt.Sprintf(`{"error":"ошибка при получении задачи: %v"}`, err), http.StatusInternalServerError)
			}
			return
		}
		response, err = json.Marshal(task)
	case http.MethodPost:
		task, err := CheckTask(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		addedTask, err := h.repo.AddTask(task)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"ошибка при добавлении задачи: %v"}`, err), http.StatusInternalServerError)
			return
		}
		response, err = json.Marshal(addedTask)
	case http.MethodPut:
		task, err := CheckTask(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = h.repo.UpdateTask(task)
		if err != nil {
			if errors.Is(err, repository.ErrTaskNotFound) {
				http.Error(w, `{"error":"задача не найдена"}`, http.StatusNotFound)
			} else {
				http.Error(w, fmt.Sprintf(`{"error":"ошибка при обновлении задачи: %v"}`, err), http.StatusInternalServerError)
			}
			return
		}
		response = []byte("{}")
	case http.MethodDelete:
		if par == "" {
			http.Error(w, `{"error":"неверный id"}`, http.StatusBadRequest)
			return
		}
		err = h.repo.DeleteTask(par)
		if err != nil {
			if errors.Is(err, repository.ErrTaskNotFound) {
				http.Error(w, `{"error":"задача не найдена"}`, http.StatusNotFound)
			} else {
				http.Error(w, fmt.Sprintf(`{"error":"ошибка при удалении задачи: %v"}`, err), http.StatusInternalServerError)
			}
			return
		}
		response = []byte("{}")
	default:
		http.Error(w, `{"error":"неподдерживаемый метод"}`, http.StatusMethodNotAllowed)
		return
	}

	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"внутренняя ошибка сервера: %v"}`, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	if err != nil {
		log.Printf("Ошибка при записи ответа: %v", err)
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}
}
func (h *TaskHandler) TasksGet(w http.ResponseWriter, req *http.Request) {
	par := req.URL.Query().Get("search")

	tasks, err := h.repo.GetTasks(par)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"ошибка при получении задач: %v"}`, err), http.StatusInternalServerError)
		return
	}

	response := map[string][]task.Task{"tasks": tasks}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"ошибка при формировании ответа: %v"}`, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(jsonResponse)
	if err != nil {
		log.Printf("Ошибка при записи ответа: %v", err)
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

}

func (h *TaskHandler) TaskDone(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")

	err := h.repo.TaskDone(id)
	if err != nil {
		if errors.Is(err, repository.ErrTaskNotFound) {
			http.Error(w, `{"error":"задача не найдена"}`, http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf(`{"error":"ошибка при выполнении задачи: %v"}`, err), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write([]byte("{}"))
	if err != nil {
		log.Printf("Ошибка при записи ответа: %v", err)
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

}

func (h *TaskHandler) NextDateHandl(w http.ResponseWriter, req *http.Request) {
	param := req.URL.Query()
	now := param.Get("now")
	day := param.Get("date")
	repeat := param.Get("repeat")

	nextDay, err := h.repo.NextDate(now, day, repeat)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"ошибка при вычислении следующей даты: %v"}`, err), http.StatusBadRequest)
		return
	}

	_, err = w.Write([]byte(nextDay))
	if err != nil {
		log.Printf("Ошибка при записи ответа: %v", err)
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

}
