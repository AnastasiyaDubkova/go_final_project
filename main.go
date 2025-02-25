package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"go_final_project/db"
	"go_final_project/handlers"

	"github.com/joho/godotenv"
)

func main() {
	// Загрузка переменных окружения из файла .env
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	// Подключение БД
	db.ConnectDB()
	defer db.DB.Close()

	// Определение порта
	// *1
	// Если существует переменная окружения TODO_PORT и/или она прописана в .env
	// присвоим данное значение в переменную port
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	mux := http.NewServeMux()

	webDir := "./web"
	mux.Handle("/", http.FileServer(http.Dir(webDir)))
	mux.HandleFunc("/api/nextdate", handlers.NextDateHandler)
	mux.HandleFunc("/api/task", handlers.TaskHandler)
	mux.HandleFunc("/api/tasks", handlers.GetTasksHandler)
	mux.HandleFunc("/api/task/done", handlers.PostTaskDoneHandler)

	// Запуск сервера
	fmt.Fprintln(os.Stdout, "Starting server", port)
	log.Printf("Starting server on :%s\n", port)

	if err := http.ListenAndServe(":"+port, mux); err != nil {
		fmt.Printf("Start server error: %s", err.Error())
		return
	}

}
