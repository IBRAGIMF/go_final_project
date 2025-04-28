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
	// проверяем что структура не пустая
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}
	// проверим что Title передан
	if task.Title == "" {
		writeJSON(w, map[string]string{"error": "Не указан заголовок задачи"})
		return
	}
	// проверим что Repeat передан
	if task.Repeat == "" {
		writeJSON(w, map[string]string{"error": "Не указано правило повторения"})
		return
	}

	// Проверяем на корректность полученное значение Date в струрктуре task
	if err := checkDate(&task); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}
	id, err := db.AddTask(&task)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, map[string]string{"id": fmt.Sprintf("%d", id)})
}

func checkDate(task *db.Task) error {
	now := time.Now()
	if task.Date == "" {
		task.Date = now.Format("20060102")
	}

	t, err := time.Parse("20060102", task.Date)
	if err != nil {
		return fmt.Errorf("Не верный формат даты")
	}
	// Если дата в прошлом и есть правило повторения
	if afterNow(now, t) && task.Repeat != "" {
		next, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return fmt.Errorf("Ошибка в правиле повторения: %v", err)
		}

		task.Date = next
	}

	// Если дата в прошлом и НЕТ правила повторения
	if afterNow(now, t) && task.Repeat == "" {
		task.Date = now.Format("20060102")
	}

	return nil
}

func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(data)
}

func afterNow(now, t time.Time) bool {
	return now.After(t)
}
