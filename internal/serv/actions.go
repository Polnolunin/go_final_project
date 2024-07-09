package actions

import (
	"os"

	_ "modernc.org/sqlite"
)

// CheckPort извлекает номер порта из переменной окружения "TODO_PORT".
// Если "TODO_PORT" не установлен, по умолчанию используется "7540".
//
// Возвращает:
// Номер порта в виде строки.
func CheckPort() string {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}
	return port
}
