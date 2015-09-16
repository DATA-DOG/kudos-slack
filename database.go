package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var db sql.DB

func loadDatabase() {
	if _, err := os.Stat(config.Database); os.IsNotExist(err) {
		createDatabase()
		return
	}

	fmt.Println("Loading database...")

	db, err := sql.Open("sqlite3", config.Database)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}

func createDatabase() {
	db, err := sql.Open("sqlite3", config.Database)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Println("Creating database...")

	sqlStmt := `CREATE TABLE kudos (kudo text);`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
}

func saveKudo(kudo Kudo) {
	//db.Exec(query string, args ...interface{})
}
