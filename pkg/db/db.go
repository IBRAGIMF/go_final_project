package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

const (
	Schema = `
		CREATE TABLE scheduler (id INTEGER PRIMARY KEY AUTOINCREMENT,
								date CHAR(8) NOT NULL DEFAULT "",
								title VARCHAR(256),
								comment TEXT,
								repeat VARCHAR(128)
							   );
		CREATE INDEX idx_date ON scheduler (date);					   
							`
)

var DB *sql.DB

func Init(dbFile string) error {
	_, err := os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}
	// Открываем (или создаём) базу данных
	DB, err = sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatal(err)
		return fmt.Errorf("db open error: %w", err)
	}
	// Если install == true, создаём таблицу и индекс
	if install {
		_, err = DB.Exec(Schema)
		if err != nil {
			return fmt.Errorf("Error create table: %w", err)
		}
		fmt.Println("база создана")
	} else {

		fmt.Println("база существуют")
	}

	return nil
}
