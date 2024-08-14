package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go_final_project/cmd/date"
	"go_final_project/cmd/task"
	"net/http"
	"regexp"
	"time"
)

func CheckTask(req *http.Request) (task.Task, int, error) {
	var task task.Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		return task, http.StatusInternalServerError, fmt.Errorf(`{"error":"ошибка чтения тела запроса"}`)
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		return task, http.StatusBadRequest, fmt.Errorf(`{"error":"неправильный формат запроса"}`)
	}

	if task.Title == "" {
		return task, http.StatusBadRequest, fmt.Errorf(`{"error":"отсутствует заголовок задачи"}`)
	}

	now := time.Now()
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)

	if task.Date == "" || task.Date == "today" {
		task.Date = now.Format(date.DateFormat)
	} else {
		if matched, _ := regexp.MatchString(`^\d{8}$`, task.Date); !matched {
			return task, http.StatusBadRequest, fmt.Errorf(`{"error":"неправильный формат даты"}`)
		}
		dateParse, err := time.Parse(date.DateFormat, task.Date)
		if err != nil {
			return task, http.StatusBadRequest, fmt.Errorf(`{"error":"неверная дата"}`)
		}

		if dateParse.Before(now) {
			if task.Repeat == "" {
				task.Date = now.Format(date.DateFormat)
			} else {
				nextDate, err := date.NextDate(now, task.Date, task.Repeat)
				if err != nil {
					return task, http.StatusBadRequest, fmt.Errorf(`{"error":"неверное правило повторения: %v"}`, err)
				}
				task.Date = nextDate
			}
		}
	}

	if task.Repeat != "" {
		_, err := date.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return task, http.StatusBadRequest, fmt.Errorf(`{"error":"неправильное правило повтора: %v"}`, err)
		}
	}

	return task, http.StatusOK, nil
}
