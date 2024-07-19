package handler

import (
	"database/sql"
	"encoding/json"
	"go_final_project/cmd/task"
	"net/http"
)

type ResponseForPostTask struct {
	Id int64 `json:"id"`
}

var ResponseStatus int

func TaskHandler(w http.ResponseWriter, req *http.Request) {
	par := req.URL.Query().Get("id")

	db, err := sql.Open("sqlite", task.FileDB)
	if err != nil {
		http.Error(w, "ошибка открытия базы данных"+err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var response []byte

	switch req.Method {
	case http.MethodGet:
		if par == "" {
			http.Error(w, `{"error":"неверный id"}`, http.StatusBadRequest)
			return
		}
		response, ResponseStatus, err = TaskID(db, par)
	case http.MethodPost:
		response, ResponseStatus, err = AddTask(db, req)
	case http.MethodPut:
		response, ResponseStatus, err = UptadeTaskID(db, req)
	case http.MethodDelete:
		ResponseStatus, err = DeleteTask(db, par)
		if err == nil {
			str := map[string]interface{}{}
			response, err = json.Marshal(str)
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
