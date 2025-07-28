package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/sajidzamanme/emi-tracker/database"
	"github.com/sajidzamanme/emi-tracker/router"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Failed to load Environment Variables:", err)
	}
}

func main() {
	fmt.Println("EMI TRACKER")

	database.InitDB()

	PORT := os.Getenv("PORT")
	err := http.ListenAndServe(PORT, router.NewMux())
	if err != nil {
		log.Fatalln("Server Crashed")
	}
}
