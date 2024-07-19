package main

import (
	"go_final_project/cmd/db"
	"go_final_project/cmd/handler"

	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "modernc.org/sqlite"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Printf("ошибка загрузки .env файла: %v", err)
		log.Println("используются значения по умолчанию")
	}

	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	db.CheckDB()
	defer db.DB.Close()

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("./web")))
	mux.HandleFunc("/api/nextdate", handler.NextDateHandl)
	mux.HandleFunc("/api/task", handler.TaskHandler)
	mux.HandleFunc("/api/tasks", handler.TasksGet)
	mux.HandleFunc("/api/task/done", handler.TaskDone)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}


	log.Println("запуск сервера: http://localhost:" + port + "/")
	err = server.ListenAndServe()
	if err != nil {
		log.Printf("ошибка при запуске сервера: %s", err)
	}
}
