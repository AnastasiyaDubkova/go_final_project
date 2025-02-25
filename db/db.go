package db

import (
	"database/sql"
	"errors"
	"fmt"
	"go_final_project/models"
	"log"
	"os"
	"time"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func ConnectDB() {
	// *2
	// Определяем путь к файлу базы данных через переменную окружения
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
	}

	DB = db

	if install {
		// Создание базы данных и индекса
		CreateScheduler := `CREATE TABLE scheduler (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            date CHAR(8) NOT NULL DEFAULT "",
            title TEXT NOT NULL DEFAULT "",
            comment TEXT,
            repeat VARCHAR(128) NOT NULL DEFAULT ""
        );`
		_, err := DB.Exec(CreateScheduler)
		if err != nil {
			log.Fatal(err)
		}

		// Создание индекса по полю date
		createIndex := `CREATE INDEX date_scheduler ON scheduler(date);`
		_, err = DB.Exec(createIndex)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("database created successfully")
	}
}

// Добавляем задачу в БД
func AddTask(date, title, comment, repeat string) (int64, error) {
	fmt.Fprintln(os.Stdout, "request modify data", date, title, comment, repeat)
	res, err := DB.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
		sql.Named("date", date),
		sql.Named("title", title),
		sql.Named("comment", comment),
		sql.Named("repeat", repeat))

	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

// Получаем список ближайших задач
func GetTasks() ([]models.Task, error) {

	now := time.Now().Format("20060102")

	var tasks []models.Task

	rows, err := DB.Query("SELECT id, date, title, comment, repeat FROM scheduler WHERE date >= :now ORDER BY date LIMIT 50",
		sql.Named("now", now))

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var task models.Task

		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)

		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

// Получаем список задач на конкретную дату
func GetTasksSearchDate(date string) ([]models.Task, error) {

	var tasks []models.Task

	rows, err := DB.Query("SELECT id, date, title, comment, repeat FROM scheduler WHERE date = :date LIMIT 50",
		sql.Named("date", date))

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var task models.Task

		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)

		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

// Получаем список задач на основе текста, используя поля title и comment
func GetTasksSearch(search string) ([]models.Task, error) {

	var tasks []models.Task

	searchParam := "%" + search + "%"

	rows, err := DB.Query("SELECT * FROM scheduler WHERE title LIKE :search OR comment LIKE :search ORDER BY date LIMIT 50",
		sql.Named("search", searchParam))

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var task models.Task

		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)

		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

// Получаем задачу по ID
func GetTask(id int64) (models.Task, error) {

	task := models.Task{}

	row := DB.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = :id",
		sql.Named("id", id))

	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)

	if err != nil {
		return models.Task{}, err
	}
	return task, nil
}

// Обновляем задачу
func UpdateTask(id, date, title, comment, repeat string) error {
	response, err := DB.Exec("UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id",
		sql.Named("date", date),
		sql.Named("title", title),
		sql.Named("comment", comment),
		sql.Named("repeat", repeat),
		sql.Named("id", id))

	if err != nil {
		return err
	}

	rowsAffected, err := response.RowsAffected()

	if err != nil || rowsAffected == 0 {
		return errors.New("task not found")
	}
	return nil
}

// Удаляем задачу из базы данных
func DeleteTask(id int64) error {
	response, err := DB.Exec("DELETE FROM scheduler WHERE id = :id", sql.Named("id", id))
	if err != nil {
		return err
	}

	rowsAffected, err := response.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return errors.New("task not found")
	}
	return nil
}
