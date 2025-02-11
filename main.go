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
	// Загружаем переменные окружения из файла .env
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	// Плдключение БД
	db.ConnectDB()

	// Определение порта
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	//r := chi.NewRouter()

	//r.Handle("/*", http.FileServer(http.Dir("./web")))

	// Запуск сервера
	//log.Printf("Starting server on :%s\n", port)
	//if err := http.ListenAndServe(":"+port, r); err != nil {
	//	fmt.Printf("Start server error: %s", err.Error())
	//	return
	//}

	mux := http.NewServeMux()

	webDir := "./web"
	mux.Handle("/", http.FileServer(http.Dir(webDir)))
	mux.HandleFunc("/api/nextdate", handlers.NextDateHandler)

	// Запуск сервера
	log.Printf("Starting server on :%s\n", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		fmt.Printf("Start server error: %s", err.Error())
		return
	}
}
