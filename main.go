package main

import (
	"database/sql"
	"go_final_project/cmd/handler"
	"go_final_project/cmd/repository"
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

	dbFile := "scheduler.db"
	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		log.Fatalf("ошибка открытия базы данных: %v", err)
	}
	defer db.Close()

	repo := repository.NewRepository(db)
	taskHandler := handler.NewTaskHandler(repo)

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("./web")))
	mux.HandleFunc("/api/nextdate", taskHandler.NextDateHandl)
	mux.HandleFunc("/api/task", taskHandler.TaskHandler)
	mux.HandleFunc("/api/tasks", taskHandler.TasksGet)
	mux.HandleFunc("/api/task/done", taskHandler.TaskDone)

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
