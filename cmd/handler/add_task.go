package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

func AddTask(db *sql.DB, req *http.Request) ([]byte, int, error) {
	var resp ResponseForPostTask

	task, ResponseStatus, err := CheckTask(req)
	if err != nil {
		return []byte{}, ResponseStatus, err
	}

	result, err := db.Exec(`INSERT INTO scheduler (date, title, comment, repeat)
		VALUES (:date, :title, :comment, :repeat)`,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
	)
	if err != nil {
		return []byte{}, http.StatusInternalServerError, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return []byte{}, http.StatusInternalServerError, err
	}

	resp.Id = id

	idResult, err := json.Marshal(resp)
	if err != nil {
		return []byte{}, http.StatusInternalServerError, err
	}
	return idResult, http.StatusOK, nil
}
