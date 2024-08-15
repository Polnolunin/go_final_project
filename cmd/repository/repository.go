package repository

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"go_final_project/cmd/date"
	"go_final_project/cmd/task"
	"net/http"
	"time"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) AddTask(req *http.Request) ([]byte, int, error) {
	var resp struct {
		Id int64 `json:"id"`
	}

	task, ResponseStatus, err := CheckTask(req)
	if err != nil {
		return []byte{}, ResponseStatus, err
	}

	result, err := r.db.Exec(`INSERT INTO scheduler (date, title, comment, repeat)
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

func (r *Repository) DeleteTask(id string) (int, error) {
	task, err := r.db.Exec("DELETE FROM scheduler WHERE id = :id", sql.Named("id", id))
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

func (r *Repository) TaskDone(id string) (int, error) {
	var taskID task.Task

	row := r.db.QueryRow("SELECT * FROM scheduler WHERE id = :id", sql.Named("id", id))

	err := row.Scan(&taskID.ID, &taskID.Date, &taskID.Title, &taskID.Comment, &taskID.Repeat)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf(`{"error":"writing date %v"}`, err)
	}

	if taskID.Repeat == "" {
		return r.DeleteTask(id)
	}

	now := time.Now()
	dataNew, err := date.NextDate(now, taskID.Date, taskID.Repeat)
	if err != nil {
		return http.StatusBadRequest, err
	}

	res, err := r.db.Exec(`UPDATE scheduler SET date = :date WHERE id = :id`,
		sql.Named("date", dataNew),
		sql.Named("id", taskID.ID))
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf(`{"error":"task is not found %v"}`, err)
	}

	result, err := res.RowsAffected()
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf(`{"error":"task is not found %v"}`, err)
	}
	if result == 0 {
		return http.StatusBadRequest, errors.New(`{"error":"task is not found"}`)
	}

	return http.StatusOK, nil
}

func (r *Repository) GetTasks(search string) ([]task.Task, int, error) {
	var tasks []task.Task

	var rows *sql.Rows
	var err error

	if search != "" {
		rows, err = r.db.Query(`SELECT id, date, title, comment, repeat FROM scheduler
		WHERE title LIKE :search OR comment LIKE :search OR date = :date ORDER BY date LIMIT 20`,
			sql.Named("search", "%"+search+"%"),
			sql.Named("date", search))
	} else {
		rows, err = r.db.Query(`SELECT id, date, title, comment, repeat FROM scheduler 
		ORDER BY date LIMIT 20`)
	}

	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	defer rows.Close()

	for rows.Next() {
		var task task.Task
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

func (r *Repository) TaskID(id string) ([]byte, int, error) {
	var task task.Task

	row := r.db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = :id",
		sql.Named("id", id))

	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return []byte{}, http.StatusInternalServerError, fmt.Errorf(`{"error":"ошибка записи %v"}`, err)
	}

	result, err := json.Marshal(task)
	if err != nil {
		return []byte{}, http.StatusInternalServerError, err
	}

	return result, http.StatusOK, nil
}

func (r *Repository) UpdateTask(req *http.Request) ([]byte, int, error) {
	taskID, responseStatus, err := CheckTask(req)
	if err != nil {
		return nil, responseStatus, err
	}

	res, err := r.db.Exec(`UPDATE scheduler SET
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

func (r *Repository) NextDate(now, day, repeat string) (string, error) {
	timeNow, err := time.Parse(date.DateFormat, now)
	if err != nil {
		return "", err
	}

	nextDay, err := date.NextDate(timeNow, day, repeat)
	if err != nil {
		return "", err
	}

	return nextDay, nil
}

func formatError(message string, err error) error {
	return fmt.Errorf(`{"error":"%s: %v"}`, message, err)
}
func (r *Repository) ConditionalTask(db *sql.DB, search string) ([]task.Task, int, error) {
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
