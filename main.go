package main

import (
	"final/pkg/db"
	"final/pkg/server"
	"log"
)

func main() {
	if err := db.Init("scheduler.db"); err != nil {
		log.Fatalf("Ошибка при создании БД: %v", err)
	}
	defer db.DB.Close()
	server.Run()

}
