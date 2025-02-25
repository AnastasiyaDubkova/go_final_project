package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"go_final_project/db"
	"go_final_project/models"
	"go_final_project/nextdate"
)

// Переключение методов
func TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		postTaskHandler(w, r)
	case http.MethodPut:
		putTaskHandler(w, r)
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	default:
		http.Error(w, "method is not allowed", http.StatusMethodNotAllowed)
	}
}

// Создание задачи
func postTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var task models.Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		respErr, _ := json.Marshal(map[string]string{"error": err.Error()})
		http.Error(w, string(respErr), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		respErr, _ := json.Marshal(map[string]string{"error": err.Error()})
		http.Error(w, string(respErr), http.StatusBadRequest)
		return
	}

	// Проверка заголовка задачи
	if task.Title == "" {
		respErr, _ := json.Marshal(map[string]string{"error": "title is empty"})
		http.Error(w, string(respErr), http.StatusBadRequest)
		return
	}

	// Проверка наличия поля date и формата даты
	if task.Date == "" {
		task.Date = time.Now().Format(nextdate.FormatTime)
	} else {
		_, err := time.Parse(nextdate.FormatTime, task.Date)
		if err != nil {
			respErr, _ := json.Marshal(map[string]string{"error": "date format is wrong"})
			http.Error(w, string(respErr), http.StatusBadRequest)
			return
		}
	}

	nowTimeStr := time.Now().Format(nextdate.FormatTime)

	if task.Date < nowTimeStr {
		if task.Repeat == "" {
			task.Date = nowTimeStr
		} else {
			nextDate, err := nextdate.NextDate(time.Now(), task.Date, task.Repeat)
			if err != nil {
				respErr, _ := json.Marshal(map[string]string{"error": err.Error()})
				http.Error(w, string(respErr), http.StatusBadRequest)
				return
			}
			task.Date = nextDate
		}
	}

	// Получение id добавленной задачи
	id, err := db.AddTask(task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		respErr, _ := json.Marshal(map[string]string{"error": "error on get task id"})
		http.Error(w, string(respErr), http.StatusInternalServerError)
		return
	}

	resp, _ := json.Marshal(map[string]string{"id": strconv.Itoa(int(id))})
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write(resp)
}

// Получение задачи по ID
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		respErr, _ := json.Marshal(map[string]string{"error": "id is empty"})
		http.Error(w, string(respErr), http.StatusBadRequest)
		return
	}

	idInt, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respErr, _ := json.Marshal(map[string]string{"error": "id format is wrong"})
		http.Error(w, string(respErr), http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(idInt)

	if err != nil {
		if err == sql.ErrNoRows {
			respErr, _ := json.Marshal(map[string]string{"error": "empty string, task not found"})
			http.Error(w, string(respErr), http.StatusNotFound)
		} else {
			respErr, _ := json.Marshal(map[string]string{"error": err.Error()})
			http.Error(w, string(respErr), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	response, err := json.Marshal(task)
	if err != nil {
		respErr, _ := json.Marshal(map[string]string{"error": err.Error()})
		http.Error(w, string(respErr), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)
}

// Обновление задачи
func putTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var task models.Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		respErr, _ := json.Marshal(map[string]string{"error": err.Error()})
		http.Error(w, string(respErr), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		respErr, _ := json.Marshal(map[string]string{"error": err.Error()})
		http.Error(w, string(respErr), http.StatusBadRequest)
		return
	}

	// Проверка ID задачи
	if task.ID == "" {
		respErr, _ := json.Marshal(map[string]string{"error": "id is empty"})
		http.Error(w, string(respErr), http.StatusBadRequest)
		return
	}

	// Проверка заголовка задачи
	if task.Title == "" {
		respErr, _ := json.Marshal(map[string]string{"error": "title is empty"})
		http.Error(w, string(respErr), http.StatusBadRequest)
		return
	}

	// Проверка наличия поля date и формата даты
	if task.Date == "" {
		task.Date = time.Now().Format(nextdate.FormatTime)
	} else {
		_, err := time.Parse(nextdate.FormatTime, task.Date)
		if err != nil {
			respErr, _ := json.Marshal(map[string]string{"error": "date format is wrong"})
			http.Error(w, string(respErr), http.StatusBadRequest)
			return
		}
	}
	nowTimeStr := time.Now().Format(nextdate.FormatTime)
	fmt.Fprintln(os.Stdout, "nowTimeStr", nowTimeStr)

	if task.Date < nowTimeStr {
		if task.Repeat == "" {
			task.Date = nowTimeStr
		} else {
			nextDate, err := nextdate.NextDate(time.Now(), task.Date, task.Repeat)
			if err != nil {
				respErr, _ := json.Marshal(map[string]string{"error": err.Error()})
				http.Error(w, string(respErr), http.StatusBadRequest)
				return
			}
			task.Date = nextDate
		}
	}

	errUpdate := db.UpdateTask(task.ID, task.Date, task.Title, task.Comment, task.Repeat)
	if errUpdate != nil {
		respErr, _ := json.Marshal(map[string]string{"error": errUpdate.Error()})
		http.Error(w, string(respErr), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(os.Stdout, "Задача", task)
	resp, _ := json.Marshal(map[string]string{})
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(resp)
}

// Удаление задачи
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")

	if idStr == "" {
		respErr, _ := json.Marshal(map[string]string{"error": "id is empty"})
		http.Error(w, string(respErr), http.StatusBadRequest)
		return
	}

	idInt, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respErr, _ := json.Marshal(map[string]string{"error": "id format is wrong"})
		http.Error(w, string(respErr), http.StatusBadRequest)
		return
	}

	// Удаление задачи
	errDelete := db.DeleteTask(idInt)
	if errDelete != nil {
		respErr, _ := json.Marshal(map[string]string{"error": errDelete.Error()})
		http.Error(w, string(respErr), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	resp, _ := json.Marshal(map[string]string{})
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(resp)
}
