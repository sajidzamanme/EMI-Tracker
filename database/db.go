package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() {
	ROOT := os.Getenv("ROOT")
	PASSWORD := os.Getenv("PASSWORD")
	dsn := fmt.Sprintf("%v:%v@tcp(127.0.0.1:3306)/emiTracker?parseTime=true", ROOT, PASSWORD)

	// Configure database connection
	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalln("sql.Open error:", err)
	}

	// Establish database connection
	if err := DB.Ping(); err != nil {
		log.Fatalln("Database not reachable:", err)
	}

	log.Println("Connected to database")
}
