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

// DONE
func GetSubByID(w http.ResponseWriter, r *http.Request) {
	subID, err := strconv.Atoi(r.PathValue("subID"))
	if err != nil {
		fmt.Fprintln(w, "Invalid ID", http.StatusBadRequest)
		return
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
	err = rows.Scan(&s.SubID, &s.OwnerID, &s.SubName, &s.TotalAmount,
		&s.PaidAmount, &s.PaymentAmount, &s.StartDate, &s.EndDate, &s.DeductDay)

	if err != nil {
		log.Println("Scan error:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s)
}

// DONE
func PostSubByUserID(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.PathValue("userID"))
	if err != nil {
		fmt.Fprintln(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var s models.Subscription
	err = json.NewDecoder(r.Body).Decode(&s)
	if err != nil {
		log.Fatalln(err)
	}

	query := `INSERT INTO
	subscriptions (ownerID, subName, totalAmount, paidAmount, paymentAmount, startDate, endDate, deductDay)
	VALUES(?, ?, ?, ?, ?, ?, ?, ?);`

	_, err = database.DB.Exec(query, userID, s.SubName, s.TotalAmount, s.PaidAmount, s.PaymentAmount, s.StartDate, s.EndDate, s.DeductDay)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "Subscription Added to ID:", userID)
}

func PutSubBySubID(w http.ResponseWriter, r *http.Request) {
	subID, err := strconv.Atoi(r.PathValue("subID"))
	if err != nil {
		fmt.Fprintln(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	query := `SELECT * FROM subscriptions WHERE subID = ?`

	rows, err := database.DB.Query(query, subID)
	if err != nil {
		log.Fatalln(err)
		return
	}
	defer rows.Close()

	if !rows.Next() {
		log.Fatalln(err)
	}
	var s models.Subscription

	rows.Scan(&s.SubID, &s.OwnerID, &s.SubName, &s.TotalAmount, &s.PaidAmount, &s.PaymentAmount, &s.StartDate, &s.EndDate, &s.DeductDay)

	err = json.NewDecoder(r.Body).Decode(&s)
	if err != nil {
		log.Fatalln(err)
	}

	query = `UPDATE subscriptions
	SET ownerID = ?,
			subName = ?,
			totalAmount = ?,
			paidAmount = ?,
			paymentAmount = ?,
			startDate = ?,
			endDate = ?,
			deductDay = ?
	WHERE subID = ?`

	_, err = database.DB.Exec(query, s.OwnerID, s.SubName, s.TotalAmount, s.PaidAmount,
		s.PaymentAmount, s.StartDate, s.EndDate, s.DeductDay, subID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Fatalln(err)
	}

	fmt.Fprintln(w, "Subscription Updated of ID:", subID)
}

func DeleteSubBySubID(w http.ResponseWriter, r *http.Request) {
	subID, err := strconv.Atoi(r.PathValue("subID"))
	if err != nil {
		fmt.Fprintln(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	query := `DELETE FROM subscriptions WHERE subID = ?`

	_, err = database.DB.Exec(query, subID)
	if err != nil {
		log.Println("Database error:", err)
		http.Error(w, "Internal Database Error", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Subscription deleted with ID: %d\n", subID)
}
