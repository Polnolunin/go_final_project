package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
)

func DeleteTask(db *sql.DB, id string) (int, error) {
	task, err := db.Exec("DELETE FROM scheduler WHERE id = :id", sql.Named("id", id))
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf(`{"error":"%s"}`, err)
	}

	rowsAffected, err := task.RowsAffected()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if rowsAffected == 0 {
		return http.StatusBadRequest, errors.New(`{"error":"задача не найдена"}`)
	}
	return http.StatusOK, nil
}
