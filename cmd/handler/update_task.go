package handler

import (
	"encoding/json"
	"fmt"
	"go_final_project/cmd/repository"
	"net/http"
)

func (h *TaskHandler) UpdateTaskHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPut {
		http.Error(w, `{"error":"Метод не поддерживается"}`, http.StatusMethodNotAllowed)
		return
	}

	task, err := CheckTask(req)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	err = h.repo.UpdateTask(task)
	if err != nil {
		if err == repository.ErrTaskNotFound {
			http.Error(w, `{"error":"задача не найдена"}`, http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf(`{"error":"ошибка при обновлении задачи: %v"}`, err), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{})
}
