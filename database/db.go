package database

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() {
	// move to .env
	dsn := "root:sajid123@tcp(127.0.0.1:3306)/emiTracker"

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalln("sql.Open error:", err)
	}

	if err := DB.Ping(); err != nil {
		log.Fatalln("Database not reachable:", err)
	}

	log.Println("Connected to database")
}