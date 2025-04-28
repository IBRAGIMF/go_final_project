package api

import (
	"net/http"

	"github.com/go-chi/chi"
)

func Init() http.Handler {
	r := chi.NewRouter()

	r.Post("/task", addTaskHandler)
	r.Get("/nextdate", NextDayHandler)
	r.Get("/tasks", tasksHandler)
	r.Put("/task", updTaskHandler)
	r.Get("/task", getTaskHandler)
	r.Post("/task/done", doneTaskHandler)
	r.Delete("/task", deleteTaskHandler)

	return r
}
