package main

import "fmt"
import "golf/pkg/server"
import "golf/pkg/db"
import "golf/pkg/api"

func main() {
	fmt.Println("Запуск сервера")
	db.Init("scheduler.db")
	api.Init()
	server.Run()
}
