package api

import (
	"net/http"

	"golf/pkg/api/db"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := db.Tasks(50) // в параметре максимальное количество записей
	if err != nil {
		// здесь вызываете функцию, которая возвращает ошибку в JSON
		// её желательно было реализовать на предыдущем шаге
		WriteError(w, "Ошибка получения задач")
		return
	}
	if tasks == nil {
		tasks = []*db.Task{}
	}
	writeJson(w, TasksResp{
		Tasks: tasks,
	})
}
