package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/sajidzamanme/emi-tracker/models"
	"github.com/sajidzamanme/emi-tracker/repo"
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
	err = repo.FindRecordByRecordID(recordID, &er)
	if errors.Is(err, repo.ErrorRecordNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = utils.EncodeJson(w, er)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Add EMIRecord to Database
func InsertRecordByUserID(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.PathValue("userID"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Save EMIRecord from Request Body to er
	var er models.EMIRecord
	err = json.NewDecoder(r.Body).Decode(&er)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("JSON Decoding failed")
		return
	}
	defer r.Body.Close()

	// Add New EMIRecord to Database
	err = repo.InsertEMIRecord(userID, &er)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Update User TotalLoaned, TotalPaid, activeEMI as new EMIRecord is added
	err = repo.UpdateUserForEMIChange(userID, er.TotalAmount, er.PaidAmount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "EMI Record Added to ID:", userID)
}

// Update EMIRecord in Database
func UpdateRecordByRecordID(w http.ResponseWriter, r *http.Request) {
	recordID, err := strconv.Atoi(r.PathValue("recordID"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var er models.EMIRecord
	err = repo.FindRecordByRecordID(recordID, &er)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Save previous amount to calculate change later
	prevTotalAmount := er.TotalAmount
	prevPaidAmount := er.PaidAmount

	// Save the new information from Request Body to er
	err = json.NewDecoder(r.Body).Decode(&er)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("JSON Decoding failed")
		return
	}
	defer r.Body.Close()

	// Save the updated er in Database
	err = repo.UpdateEMIRecord(recordID, &er)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Update user with changed EMIRecord
	changeOfTotalAmount := er.TotalAmount - prevTotalAmount
	changeofPaidAmount := er.PaidAmount - prevPaidAmount
	err = repo.UpdateUserForEMIChange(er.OwnerID, changeOfTotalAmount, changeofPaidAmount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
	err = repo.FindRecordByRecordID(recordID, &er)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Update User TotalLoaned, TotalPaid, activeEMI as EMIRecord is being deleted
	err = repo.UpdateUserForEMIChange(er.OwnerID, (-1 * er.TotalAmount), (-1 * er.PaidAmount))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Delete EMIRecord from database
	err = repo.DeleteEMIRecord(er.OwnerID, recordID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "EMI Record deleted with ID: %d\n", recordID)
}

// Increase TotalPaidAmount & CurrentlyPaidAmount by InstallmentAmount
func PayInstallment(w http.ResponseWriter, r *http.Request) {
	recordID, err := strconv.Atoi(r.PathValue("recordID"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var er models.EMIRecord
	err = repo.FindRecordByRecordID(recordID, &er)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return if EMI is already completed
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
	err = repo.UpdateEMIRecord(recordID, &er)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Update User TotalPaid and CurrentlyPaid using previously saved extra
	changeAmount := er.InstallmentAmount - extra
	err = repo.UpdateUserForInstallment(er.OwnerID, -changeAmount, changeAmount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// If paying installment completes EMI:
	if er.TotalAmount == er.PaidAmount {
		err = repo.CompleteEMI(er.OwnerID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	fmt.Fprintln(w, "Installment paid of ID:", recordID)
}
