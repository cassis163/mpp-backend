package database

import (
	"database/sql"
	"fmt"
	"mpp/util"
	"os"
	"path/filepath"
	
	_ "github.com/mattn/go-sqlite3"
)

func Init() *sql.DB {
	wd, wdErr := os.Getwd()
	util.CheckErr(wdErr)
	databasePath := filepath.Join(wd, "/movies.db")
	fmt.Println(databasePath)
	db, err := sql.Open("sqlite3", databasePath)
	util.CheckErr(err)
	db.SetMaxOpenConns(1)
	fmt.Println("Started database connection (sqlite3)")

	return db
}
