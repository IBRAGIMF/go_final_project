package api

import (
	"encoding/json"
	"final/pkg/db"
	"fmt"
	"net/http"
	"time"
)

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}
	// проверим, что задан заголовок
	if task.Title == "" {
		writeJSON(w, map[string]string{"error": "Не указан заголовок задачи"})
		return
	}
	if err := checkDate(&task); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}
	// сохраняем задачу в БД
	id, err := db.AddTask(&task)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, map[string]string{"id": fmt.Sprintf("%d", id)})
}

func checkDate(task *db.Task) error {
	now := time.Now()
	today := now.Format("20060102")
	if task.Date == "" {
		task.Date = today
	}
	// проверяем формат даты
	if _, err := time.Parse("20060102", task.Date); err != nil {
		return fmt.Errorf("Не верный формат даты")
	}
	if task.Date < today {
		if task.Repeat != "" {
			next, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				return fmt.Errorf("Ошибка в правиле повторения: %v", err)
			}
			task.Date = next
		} else {
			task.Date = today
		}
	}
	return nil
}

func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(data)
}
