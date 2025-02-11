package handlers

import (
	"net/http"
	"time"

	"go_final_project/nextdate"
)

func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	now := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	nowFormat, err := time.Parse(nextdate.FormatTime, now)
	if err != nil {
		http.Error(w, "incorrect date", http.StatusBadRequest)
		return
	}

	nextDate, err := nextdate.NextDate(nowFormat, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, _ = w.Write([]byte(nextDate))
}
