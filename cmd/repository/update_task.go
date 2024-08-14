package repository

import (
	"database/sql"
	"fmt"
	"net/http"
)

func UpdateTaskHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPut {
		http.Error(w, `{"error":"Метод не поддерживается"}`, http.StatusMethodNotAllowed)
		return
	}

	db, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		http.Error(w, `{"error":"Ошибка открытия базы данных"}`, http.StatusInternalServerError)
		return
	}
	defer db.Close()

	repo := NewRepository(db)
	result, status, err := repo.UpdateTask(req)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), status)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}
