package server

import (
	"final/pkg/api"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

const (
	WebDir      = "web"
	DefaultPort = "5050"
)

// Запускаем Web сервер
func Run() {

	r := chi.NewRouter()
	r.Mount("/api", api.Init())
	r.Handle("/*", http.FileServer(http.Dir(WebDir)))

	if err := http.ListenAndServe(":5050", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}

}
