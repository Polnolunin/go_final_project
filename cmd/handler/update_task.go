package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"go_final_project/cmd/task"
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

	status, err := UpdateTask(db, req)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), status)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))
}

func UpdateTask(db *sql.DB, req *http.Request) (int, error) {
	task, status, err := CheckTask(req)
	if err != nil {
		return status, err
	}

	_, err = db.Exec(`UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`,
		task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func TaskID(db *sql.DB, id string) ([]byte, int, error) {
	var task task.Task

	row := db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = :id",
		sql.Named("id", id))

	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return []byte{}, http.StatusInternalServerError, fmt.Errorf(`{"error":"ошибка записи %v"}`, err)
	}

	if err := row.Err(); err != nil {
		return []byte{}, http.StatusInternalServerError, fmt.Errorf(`{"error":"ошибка записи %v"}`, err)
	}

	result, err := json.Marshal(task)
	if err != nil {
		return []byte{}, http.StatusInternalServerError, err
	}

	return result, http.StatusOK, nil
}

func UptadeTaskID(db *sql.DB, req *http.Request) ([]byte, int, error) {
	taskID, responseStatus, err := CheckTask(req)
	if err != nil {
		return nil, responseStatus, err
	}

	res, err := db.Exec(`UPDATE scheduler SET
	date = :date, title = :title, comment = :comment, repeat = :repeat
	WHERE id = :id`,
		sql.Named("date", taskID.Date),
		sql.Named("title", taskID.Title),
		sql.Named("comment", taskID.Comment),
		sql.Named("repeat", taskID.Repeat),
		sql.Named("id", taskID.ID))
	if err != nil {
		return nil, http.StatusInternalServerError, formatError("ошибка при обновлении задачи", err)
	}

	result, err := res.RowsAffected()
	if err != nil {
		return nil, http.StatusInternalServerError, formatError("ошибка при получении количества затронутых строк", err)
	}
	if result == 0 {
		return nil, http.StatusBadRequest, errors.New(`{"error":"задача не найдена"}`)
	}

	return []byte("{}"), http.StatusOK, nil
}

func formatError(message string, err error) error {
	return fmt.Errorf(`{"error":"%s: %v"}`, message, err)
}
