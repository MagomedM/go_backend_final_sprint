package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"golf/pkg/date"
	"unicode"
)
import "golf/pkg/api/db"

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	body, err := io.ReadAll(r.Body)
	if err != nil {
		WriteError(w, "Ошибка чтения тела запроса")
		return
	}
	defer r.Body.Close()
	if len(body) == 0 {
		WriteError(w, "Пустое тело запроса")
		return
	}
	err3 := json.Unmarshal(body, &task)
	if err3 != nil {
		WriteError(w, "Ошибка десериализации")
		return
	}
	if task.Title == "" {
		WriteError(w, "Пустой заголовок")
		return
	}
	err1 := checkDate(&task)
	if err1 != nil {
		WriteError(w, "Ошибка проверки")
		return
	}
	id, err2 := db.AddTask(&task)
	if err2 != nil {
		WriteError(w, "Ошибка добавления задачи")
		return
	}
	response := map[string]interface{}{
		"id": id,
	}
	writeJson(w, response)
}

func checkDate(task *db.Task) error {
	now := time.Now()
	if task.Date == "" {
		task.Date = now.Format("20060102")
		return nil
	}
	if task.Date == now.Format("20060102") {
		return nil
	}
	t, err := time.Parse("20060102", task.Date)
	if err != nil {
		return err
	}
	if task.Repeat == "" {
		if t.Before(now) {
			task.Date = now.Format("20060102")
		}
		return nil
	}
	if task.Repeat == "d 1" {
		task.Date = now.Format("20060102")
		return nil
	}
	next, err := date.NextDate(now, task.Date, task.Repeat)
	if err != nil {
		return err
	}
	// если сегодня (now) больше task.Date (t)
	if date.AfterNow(now, t) {
		if task.Repeat == "" {
			// если правила повторения нет, то берём сегодняшнее число
			task.Date = now.Format("20060102")
		} else {
			// в противном случае, берём вычисленную ранее следующую дату
			task.Date = next
		}
	}
	return nil
}
func writeJson(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	jsonData, err := json.Marshal(data)
	if err != nil {
		WriteError(w, "Ошибка сериализации")
		return
	}
	_, err = w.Write(jsonData)
	if err != nil {
		WriteError(w, "Ошибка отправления")
		return
	}
}

func WriteError(w http.ResponseWriter, errorMessage string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	response := map[string]interface{}{
		"error": errorMessage,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Ошибка при отправке ошибки: %v", err)
	}
}

func GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		WriteError(w, "Id не указан")
		return
	}
	task, err := db.GetTask(id)
	if err != nil {
		WriteError(w, "Задача "+id+" не найдена")
		return
	}
	writeJson(w, task)
}

func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	body, err := io.ReadAll(r.Body)
	if err != nil {
		WriteError(w, "Ошибка чтения тела запроса")
		return
	}
	defer r.Body.Close()
	if len(body) == 0 {
		WriteError(w, "Пустое тело запроса")
		return
	}
	err3 := json.Unmarshal(body, &task)
	if err3 != nil {
		WriteError(w, "Ошибка десериализации")
		return
	}
	if task.Title == "" {
		WriteError(w, "Пустой заголовок")
		return
	}
	err1 := checkDate(&task)
	if err1 != nil {
		WriteError(w, "Ошибка проверки")
		return
	}
	err2 := db.UpdateTask(&task)
	if err2 != nil {
		WriteError(w, "Ошибка добавления задачи")
		return
	}
	var nullTask db.Task
	nullTask = db.Task{}
	writeJson(w, nullTask)
}

func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		WriteError(w, "Пустой id")
		return
	}
	if !unicode.IsDigit(rune(id[0])) {
		WriteError(w, "Неправильный id")
		return
	}
	db.DeleteTask(id)
	writeJson(w, nil)
}
