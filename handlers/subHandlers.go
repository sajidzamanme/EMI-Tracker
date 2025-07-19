package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/sajidzamanme/emi-tracker/database"
	"github.com/sajidzamanme/emi-tracker/models"
)

// DONE (DATE ERROR)
func GetSubByID(w http.ResponseWriter, r *http.Request) {
	subID, err := strconv.Atoi(r.PathValue("subID"))
	if err != nil {
		fmt.Fprintln(w, "Invalid ID", http.StatusBadRequest)
	}

	query := `SELECT * FROM subscriptions WHERE subID = ?;`

	rows, err := database.DB.Query(query, subID)
	if err != nil {
		log.Fatalln(err)
	}
	defer rows.Close()

	if !rows.Next() {
		http.Error(w, "Subscription not found", http.StatusNotFound)
		return
	}
	var s models.Subscription
	rows.Scan(&s.SubID, &s.OwnerID, &s.SubName, &s.TotalAmount,
		&s.PaidAmount, &s.PaymentAmount, &s.StartDate, &s.EndDate, &s.DeductDay)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s)
}

func PostSubByUserID(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.PathValue("userID"))
	if err != nil {
		fmt.Fprintln(w, "Invalid ID", http.StatusBadRequest)
	}

	fmt.Fprintln(w, "Subscription Added to ID:", userID)
}

func PutSubBySubID(w http.ResponseWriter, r *http.Request) {
	subID, err := strconv.Atoi(r.PathValue("subID"))
	if err != nil {
		fmt.Fprintln(w, "Invalid ID", http.StatusBadRequest)
	}

	fmt.Fprintln(w, "Subscription Updated of ID:", subID)
}

func DeleteSubBySubID(w http.ResponseWriter, r *http.Request) {
	subID, err := strconv.Atoi(r.PathValue("subID"))
	if err != nil {
		fmt.Fprintln(w, "Invalid ID", http.StatusBadRequest)
	}

	fmt.Fprintln(w, "Subscription Deleted with ID:", subID)
}
