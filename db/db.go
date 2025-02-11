package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

func ConnectDB() {
	// Получение пути к файлу базы данных из переменной окружения
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}

	var install bool
	if _, err := os.Stat(dbFile); err != nil {
		if os.IsNotExist(err) {
			install = true
			fmt.Println("database does not exist")
		} else {
			panic(err)
		}
	} else {
		fmt.Println("database already exists")
	}

	// Открытие базы данных
	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		fmt.Println(err)
		return
		//log.Fatal(err)
	}
	defer db.Close()

	if install {
		// Создание базы данных и индекса
		CreateScheduler := `CREATE TABLE scheduler (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            date CHAR(8) NOT NULL DEFAULT "",
            title TEXT NOT NULL DEFAULT "",
            comment TEXT,
            repeat VARCHAR(128) NOT NULL DEFAULT ""
        );`
		_, err := db.Exec(CreateScheduler)
		if err != nil {
			log.Fatal(err)
		}

		// Создание индекса по полю date
		createIndex := `CREATE INDEX date_scheduler ON scheduler(date);`
		_, err = db.Exec(createIndex)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("database created successfully")
	}
}
