package repository

import (
	"database/sql"
	"errors"
	"strconv"
	"time"

	"go_final_project/cmd/date"
	"go_final_project/cmd/task"

	_ "modernc.org/sqlite"
)

var ErrTaskNotFound = errors.New("задача не найдена")

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) AddTask(t task.Task) (task.Task, error) {

	result, err := r.db.Exec(`INSERT INTO scheduler (date, title, comment, repeat)
        VALUES (:date, :title, :comment, :repeat)`,
		sql.Named("date", t.Date),
		sql.Named("title", t.Title),
		sql.Named("comment", t.Comment),
		sql.Named("repeat", t.Repeat),
	)
	if err != nil {
		return task.Task{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return task.Task{}, err
	}

	t.ID = strconv.FormatInt(id, 10)
	return t, nil
}

func (r *Repository) DeleteTask(id string) error {
	result, err := r.db.Exec("DELETE FROM scheduler WHERE id = ?", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrTaskNotFound
	}

	return nil
}

func (r *Repository) TaskDone(id string) error {
	var taskID task.Task

	err := r.db.QueryRow("SELECT * FROM scheduler WHERE id = ?", id).Scan(&taskID.ID, &taskID.Date, &taskID.Title, &taskID.Comment, &taskID.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrTaskNotFound
		}
		return err
	}

	if taskID.Repeat == "" {
		return r.DeleteTask(id)
	}

	now := time.Now()
	dataNew, err := date.NextDate(now, taskID.Date, taskID.Repeat)
	if err != nil {
		return err
	}

	result, err := r.db.Exec(`UPDATE scheduler SET date = ? WHERE id = ?`, dataNew, taskID.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrTaskNotFound
	}

	return nil
}

func (r *Repository) GetTasks(search string) ([]task.Task, error) {
	var tasks []task.Task
	var rows *sql.Rows
	var err error

	if search != "" {
		rows, err = r.db.Query(`SELECT id, date, title, comment, repeat FROM scheduler
        WHERE title LIKE ? OR comment LIKE ? OR date = ? ORDER BY date LIMIT 20`,
			"%"+search+"%", "%"+search+"%", search)
	} else {
		rows, err = r.db.Query(`SELECT id, date, title, comment, repeat FROM scheduler 
        ORDER BY date LIMIT 20`)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var task task.Task
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if tasks == nil {
		tasks = []task.Task{}
	}

	return tasks, rows.Err()
}

func (r *Repository) TaskID(id string) (task.Task, error) {
	var t task.Task

	err := r.db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", id).Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			return task.Task{}, ErrTaskNotFound
		}
		return task.Task{}, err
	}

	return t, nil
}

func (r *Repository) UpdateTask(t task.Task) error {
	result, err := r.db.Exec(`UPDATE scheduler SET
    date = ?, title = ?, comment = ?, repeat = ?
    WHERE id = ?`,
		t.Date, t.Title, t.Comment, t.Repeat, t.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrTaskNotFound
	}

	return nil
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

func (r *Repository) ConditionalTask(search string) ([]task.Task, error) {
	var tasks []task.Task
	var rows *sql.Rows
	var err error

	timeSearch, err := time.Parse("02.01.2006", search)
	if err == nil {
		dateFormat := timeSearch.Format("20060102")
		rows, err = r.db.Query(`SELECT id, date, title, comment, repeat FROM scheduler
        WHERE date = ? LIMIT 20`, dateFormat)
	} else {
		rows, err = r.db.Query(`SELECT id, date, title, comment, repeat FROM scheduler
        WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT 20`,
			"%"+search+"%", "%"+search+"%")
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var task task.Task
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}
