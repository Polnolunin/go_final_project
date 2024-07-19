package db

import (
	"database/sql"
	"fmt"
	"go_final_project/cmd/task"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func CheckDB() error {

	appPath, err := os.Executable()
	if err != nil {
		return err
	}

	dbFile := filepath.Join(filepath.Dir(appPath), task.FileDB)
	_, err = os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}

	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}
	defer db.Close()

	if install {
		statement, err := db.Prepare(`CREATE TABLE IF NOT EXISTS scheduler 
		(id INTEGER PRIMARY KEY AUTOINCREMENT,
		date CHAR(8) NOT NULL DEFAULT "",
   		title VARCHAR(128) NOT NULL DEFAULT "",
   		comment TEXT NOT NULL DEFAULT "",
  		repeat VARCHAR(128) NOT NULL DEFAULT "");

		CREATE INDEX IF NOT EXISTS date_indx ON scheduler (date);
		`)
		if err != nil {
			return fmt.Errorf("ошибка создания базы данных. %v", err)
		}
		statement.Exec()
	}
	return nil
}
