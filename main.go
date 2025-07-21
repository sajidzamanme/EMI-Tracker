package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/sajidzamanme/emi-tracker/database"
	"github.com/sajidzamanme/emi-tracker/router"
)

func main() {
	fmt.Println("EMI TRACKER")

	database.InitDB()

	err := http.ListenAndServe(":8080", router.NewMux())
	if err != nil {
		log.Fatalln("Server Crashed")
	}
}
