package api

import (
	"golf/pkg/constant"
	"golf/pkg/date"
	"net/http"
	"time"

	"github.com/MagomedM/go_backend_final_sprint/pkg/api/db"
)

func Init() {
	http.HandleFunc("/api/nextdate", nextDayHandler)
	http.HandleFunc("/api/task", taskHandler)
	http.HandleFunc("/api/tasks", tasksHandler)
	http.HandleFunc("/api/task/done", doneHandler)
}

func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	if nowStr == "" {
		nowStr = time.Now().Format("20060102")
	}
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")
	now, err := time.Parse(constant.DateFormat, nowStr)
	if err != nil {
		WriteError(w, "Invalid 'now' parameter")
		return
	}
	nextDate, err := date.NextDate(now, dateStr, repeat)
	if err != nil {
		WriteError(w, "Error calculating next date")
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(nextDate))
	return
}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	// обработка других методов будет добавлена на следующих шагах
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodGet:
		GetTaskHandler(w, r)
	case http.MethodPut:
		UpdateTaskHandler(w, r)
	case http.MethodDelete:
		DeleteTaskHandler(w, r)
	}
}

func doneHandler(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	id := r.URL.Query().Get("id")
	task, err := db.GetTask(id)
	if err != nil {
		WriteError(w, "Ошибка получения задачи")
		return
	}
	if task.Repeat == "" {
		err = db.DeleteTask(id)
		if err != nil {
			WriteError(w, "Ошибка удаления задачи")
			return
		}
	} else {
		next, err := date.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			WriteError(w, "Ошибка получения следующей даты")
			return
		}
		err = db.UpdateDate(next, id)
		if err != nil {
			WriteError(w, "Ошибка обновления задачи")
			return
		}
	}
	writeJson(w, nil)
}
