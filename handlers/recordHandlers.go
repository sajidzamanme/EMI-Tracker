package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/sajidzamanme/emi-tracker/database"
	"github.com/sajidzamanme/emi-tracker/models"
	"github.com/sajidzamanme/emi-tracker/utils"
)

// JSON Response with Record Details
func GetRecordByRecordID(w http.ResponseWriter, r *http.Request) {
	recordID, err := strconv.Atoi(r.PathValue("recordID"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var er models.EMIRecord
	err = utils.FindRecordByRecordID(recordID, &er)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error scanning emiRecords row: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(er)
}

// Add EMIRecord to Database
func PostRecordByUserID(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.PathValue("userID"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Save EMIRecord from request body to er
	var er models.EMIRecord
	err = json.NewDecoder(r.Body).Decode(&er)
	if err != nil {
		http.Error(w, "Invalid Record Entry", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Add New EMIRecord to Database
	query := `INSERT INTO
						emiRecords (ownerID, title, totalAmount, paidAmount, installmentAmount, startDate, endDate, deductDay)
						VALUES(?, ?, ?, ?, ?, ?, ?, ?);`

	_, err = database.DB.Exec(query, userID, er.Title, er.TotalAmount, er.PaidAmount, er.InstallmentAmount, er.StartDate, er.EndDate, er.DeductDay)
	if err != nil {
		http.Error(w, "Internal Database Error", http.StatusInternalServerError)
		log.Printf("Database error: %v\n", err)
		return
	}

	// Update User TotalLoaned, TotalPaid, CurrentlyLoaned & CurrentlyPaid as new EMIRecord is added
	query = `UPDATE users
					SET totalLoaned = totalLoaned + ?,
							totalPaid = totalPaid + ?,
							currentlyLoaned = currentlyLoaned + ?,
							currentlyPaid = currentlyPaid + ?
					WHERE userID = ?`

	_, err = database.DB.Exec(query, er.TotalAmount, er.PaidAmount, er.TotalAmount, er.PaidAmount, userID)
	if err != nil {
		http.Error(w, "Internal Database Error", http.StatusInternalServerError)
		log.Printf("Database error: %v\n", err)
		return
	}

	fmt.Fprintln(w, "EMI Record Added to ID:", userID)
}

// Update EMIRecord in Database
func PutRecordByRecordID(w http.ResponseWriter, r *http.Request) {
	recordID, err := strconv.Atoi(r.PathValue("recordID"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var er models.EMIRecord
	err = utils.FindRecordByRecordID(recordID, &er)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error scanning emiRecords row: %v", err)
		return
	}

	// Save previous amount to calculate change later
	prevTotalAmount := er.TotalAmount
	prevPaidAmount := er.PaidAmount

	// Save the new information in er
	err = json.NewDecoder(r.Body).Decode(&er)
	if err != nil {
		http.Error(w, "Invalid Record Entry", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Save the updated er in Database
	query := `UPDATE emiRecords
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

	// Update User TotalLoaned, TotalPaid, CurrentlyLoaned & CurrentlyPaid as EMIRecord is updated
	// Use previous amount to calculate the change, and add it
	query = `UPDATE users
					SET totalLoaned = totalLoaned + ?,
							totalPaid = totalPaid + ?,
							currentlyLoaned = currentlyLoaned + ?,
							currentlyPaid = currentlyPaid + ?
					WHERE userID = ?`

	_, err = database.DB.Exec(query, er.TotalAmount-prevTotalAmount, er.PaidAmount-prevPaidAmount,
		er.TotalAmount-prevTotalAmount, er.PaidAmount-prevPaidAmount, er.OwnerID)
	if err != nil {
		http.Error(w, "Internal Database Error", http.StatusInternalServerError)
		log.Printf("Database error: %v\n", err)
		return
	}

	fmt.Fprintln(w, "EMI Record Updated of ID:", recordID)
}

// Delete EMIRecord from Database
func DeleteRecordByRecordID(w http.ResponseWriter, r *http.Request) {
	recordID, err := strconv.Atoi(r.PathValue("recordID"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var er models.EMIRecord
	err = utils.FindRecordByRecordID(recordID, &er)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error scanning emiRecords row: %v", err)
		return
	}

	// Update User TotalLoaned, TotalPaid, CurrentlyLoaned & CurrentlyPaid as EMIRecord is being deleted
	query := `UPDATE users
						SET totalLoaned = totalLoaned - ?,
								totalPaid = totalPaid - ?,
								currentlyLoaned = currentlyLoaned - ?,
								currentlyPaid = currentlyPaid - ?
						WHERE userID = ?`

	_, err = database.DB.Exec(query, er.TotalAmount, er.PaidAmount, er.TotalAmount, er.PaidAmount, er.OwnerID)
	if err != nil {
		http.Error(w, "Internal Database Error", http.StatusInternalServerError)
		log.Printf("Database error: %v\n", err)
		return
	}

	// Delete EMIRecord from database
	query = `DELETE FROM emiRecords WHERE recordID = ?`

	_, err = database.DB.Exec(query, recordID)
	if err != nil {
		http.Error(w, "Internal Database Error", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "EMI Record deleted with ID: %d\n", recordID)
}

// Increase TotalPaidAmount & CurrentlyPaidAmount by InstallmentAmount
func GetPayInstallment(w http.ResponseWriter, r *http.Request) {
	recordID, err := strconv.Atoi(r.PathValue("recordID"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var er models.EMIRecord
	err = utils.FindRecordByRecordID(recordID, &er)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error scanning emiRecords row: %v", err)
		return
	}

	if er.PaidAmount == er.TotalAmount {
		http.Error(w, "EMI Already Paid", http.StatusNotAcceptable)
		return
	}

	// Increase PaidAmount by InstallAmount
	// If it overflows then save the extra amount to use later
	er.PaidAmount += er.InstallmentAmount
	extra := 0
	if er.PaidAmount > er.TotalAmount {
		extra = er.PaidAmount - er.TotalAmount
		er.PaidAmount = er.TotalAmount
	}

	// Update EMIRecord in the Database
	query := `UPDATE emiRecords
						SET paidAmount = ?
						WHERE recordID = ?`

	_, err = database.DB.Exec(query, er.PaidAmount, recordID)
	if err != nil {
		http.Error(w, "Internal Database Error", http.StatusInternalServerError)
		log.Printf("Database error: %v\n", err)
		return
	}

	// Update User TotalPaid and CurrentlyPaid using previously save extra
	query = `UPDATE users
					SET totalPaid = totalPaid + ?,
							currentlyPaid = currentlyPaid + ?
					WHERE userID = ?`

	_, err = database.DB.Exec(query, er.InstallmentAmount-extra, er.InstallmentAmount-extra, er.OwnerID)
	if err != nil {
		http.Error(w, "Internal Database Error", http.StatusInternalServerError)
		log.Printf("Database error: %v\n", err)
		return
	}

	// If paying installment completes EMI:
	if er.TotalAmount == er.PaidAmount {
		query = `UPDATE users
						SET currentlyLoaned = currentlyLoaned - ?,
								currentlyPaid = currentlyPaid - ?,
								completedEMI = completedEMI + 1
						WHERE userID = ?`

		_, err = database.DB.Exec(query, er.TotalAmount, er.TotalAmount, er.OwnerID)
		if err != nil {
			http.Error(w, "Internal Database Error", http.StatusInternalServerError)
			log.Printf("Database error: %v\n", err)
			return
		}
	}

	fmt.Fprintln(w, "Installment paid of ID:", recordID)
}
