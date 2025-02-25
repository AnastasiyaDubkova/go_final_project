package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"go_final_project/db"
	"go_final_project/models"
)

func GetTasksHandler(w http.ResponseWriter, r *http.Request) {

	searchParam := r.URL.Query().Get("search")

	var tasks []models.Task
	var err error

	if searchParam != "" {

		searchDate, err := time.Parse("02.01.2006", searchParam)

		if err != nil {
			tasks, err = db.GetTasksSearch(searchParam)
		} else {
			searchDateStr := searchDate.Format("20060102")
			tasks, err = db.GetTasksSearchDate(searchDateStr)
		}
	} else {
		tasks, err = db.GetTasks()
	}

	if err != nil {
		respErr, _ := json.Marshal(map[string]string{"error": "error on get tasks"})
		http.Error(w, string(respErr), http.StatusInternalServerError)
		return
	}

	res := map[string]interface{}{
		"tasks": tasks,
	}

	if tasks == nil {
		res["tasks"] = []models.Task{}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	resp, err := json.Marshal(res)
	if err != nil {
		respErr, _ := json.Marshal(map[string]string{"error": err.Error()})
		http.Error(w, string(respErr), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(resp)

}
