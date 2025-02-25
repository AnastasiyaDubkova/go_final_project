package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"go_final_project/db"
	"go_final_project/nextdate"
)

// Отметка о выполнении задачи
func PostTaskDoneHandler(w http.ResponseWriter, r *http.Request) {
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

	// Проверка существования задачи (получаем задачу)
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

	// Проверка поля repeat (периодичности задачи)
	if task.Repeat == "" {
		err := db.DeleteTask(idInt)
		if err != nil {
			respErr, _ := json.Marshal(map[string]string{"error": err.Error()})
			http.Error(w, string(respErr), http.StatusInternalServerError)
			return
		}
	} else {
		nextDate, err := nextdate.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			respErr, _ := json.Marshal(map[string]string{"error": err.Error()})
			http.Error(w, string(respErr), http.StatusBadRequest)
			return
		}

		errUpdate := db.UpdateTask(task.ID, nextDate, task.Title, task.Comment, task.Repeat)
		if errUpdate != nil {
			respErr, _ := json.Marshal(map[string]string{"error": errUpdate.Error()})
			http.Error(w, string(respErr), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	resp, _ := json.Marshal(map[string]string{})
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(resp)
}
