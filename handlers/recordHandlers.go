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
func GetRecordByRecordID(w http.ResponseWriter, r *http.Request) {
	recordID, err := strconv.Atoi(r.PathValue("recordID"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	query := `SELECT * FROM emiRecords WHERE recordID = ?;`

	rows := database.DB.QueryRow(query, recordID)
	var er models.EMIRecord
	err = rows.Scan(&er.RecordID, &er.OwnerID, &er.Title, &er.TotalAmount,
		&er.PaidAmount, &er.InstallmentAmount, &er.StartDate, &er.EndDate, &er.DeductDay)
	if err != nil {
		log.Printf("Error scanning emiRecords row: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(er)
}

// DONE
func PostRecordByUserID(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.PathValue("userID"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var er models.EMIRecord
	err = json.NewDecoder(r.Body).Decode(&er)
	if err != nil {
		http.Error(w, "Invalid Record Entry", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	query := `INSERT INTO
	emiRecords (ownerID, title, totalAmount, paidAmount, installmentAmount, startDate, endDate, deductDay)
	VALUES(?, ?, ?, ?, ?, ?, ?, ?);`

	_, err = database.DB.Exec(query, userID, er.Title, er.TotalAmount, er.PaidAmount, er.InstallmentAmount, er.StartDate, er.EndDate, er.DeductDay)
	if err != nil {
		http.Error(w, "Internal Database Error", http.StatusInternalServerError)
		log.Printf("Database error: %v\n", err)
		return
	}

	query = `UPDATE users
					SET totalLoaned = totalLoaned + ?, totalPaid = totalPaid + ?
					WHERE userID = ?`

	_, err = database.DB.Exec(query, er.TotalAmount, er.PaidAmount, userID)
	if err != nil {
		http.Error(w, "Internal Database Error", http.StatusInternalServerError)
		log.Printf("Database error: %v\n", err)
		return
	}

	fmt.Fprintln(w, "EMI Record Added to ID:", userID)
}

// DONE
func PutRecordByRecordID(w http.ResponseWriter, r *http.Request) {
	recordID, err := strconv.Atoi(r.PathValue("recordID"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	query := `SELECT * FROM emiRecords WHERE recordID = ?`

	rows := database.DB.QueryRow(query, recordID)
	var er models.EMIRecord
	err = rows.Scan(&er.RecordID, &er.OwnerID, &er.Title, &er.TotalAmount,
		&er.PaidAmount, &er.InstallmentAmount, &er.StartDate, &er.EndDate, &er.DeductDay)
	if err != nil {
		log.Printf("Error scanning emiRecords row: %v", err)
		return
	}

	prevTotalAmount := er.TotalAmount
	prevPaidAmount := er.PaidAmount

	// new data from request
	err = json.NewDecoder(r.Body).Decode(&er)
	if err != nil {
		http.Error(w, "Invalid Record Entry", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	query = `UPDATE emiRecords
	SET ownerID = ?,
			title = ?,
			totalAmount = ?,
			paidAmount = ?,
			installmentAmount = ?,
			startDate = ?,
			endDate = ?,
			deductDay = ?
	WHERE recordID = ?`

	_, err = database.DB.Exec(query, er.OwnerID, er.Title, er.TotalAmount, er.PaidAmount,
		er.InstallmentAmount, er.StartDate, er.EndDate, er.DeductDay, recordID)
	if err != nil {
		http.Error(w, "Internal Database Error", http.StatusInternalServerError)
		log.Printf("Database error: %v\n", err)
		return
	}

	query = `UPDATE users
	SET totalLoaned = totalLoaned + ?, totalPaid = totalPaid + ?
	WHERE userID = ?`

	_, err = database.DB.Exec(query, er.TotalAmount-prevTotalAmount, er.PaidAmount-prevPaidAmount, er.OwnerID)
	if err != nil {
		http.Error(w, "Internal Database Error", http.StatusInternalServerError)
		log.Printf("Database error: %v\n", err)
		return
	}

	fmt.Fprintln(w, "EMI Record Updated of ID:", recordID)
}

// DONE
func DeleteRecordByRecordID(w http.ResponseWriter, r *http.Request) {
	recordID, err := strconv.Atoi(r.PathValue("recordID"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	query := `SELECT * FROM emiRecords WHERE recordID = ?;`

	rows := database.DB.QueryRow(query, recordID)
	var er models.EMIRecord
	err = rows.Scan(&er.RecordID, &er.OwnerID, &er.Title, &er.TotalAmount,
		&er.PaidAmount, &er.InstallmentAmount, &er.StartDate, &er.EndDate, &er.DeductDay)
	if err != nil {
		log.Printf("Error scanning emiRecords row: %v", err)
		return
	}

	query = `UPDATE users
					SET totalLoaned = totalLoaned - ?, totalPaid = totalPaid - ?
					WHERE userID = ?`

	_, err = database.DB.Exec(query, er.TotalAmount, er.PaidAmount, er.OwnerID)
	if err != nil {
		http.Error(w, "Internal Database Error", http.StatusInternalServerError)
		log.Printf("Database error: %v\n", err)
		return
	}

	query = `DELETE FROM emiRecords WHERE recordID = ?`

	_, err = database.DB.Exec(query, recordID)
	if err != nil {
		http.Error(w, "Internal Database Error", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "EMI Record deleted with ID: %d\n", recordID)
}
