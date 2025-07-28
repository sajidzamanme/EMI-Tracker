package repo

import (
	"database/sql"
	"errors"
	"log"

	"github.com/sajidzamanme/emi-tracker/database"
	"github.com/sajidzamanme/emi-tracker/models"
)

var ErrorRecordNotFound = errors.New("Record Not Found")

func FindRecordByRecordID(recordID int, er *models.EMIRecord) error {
	query := `SELECT * FROM emiRecords WHERE recordID = ?;`

	rows := database.DB.QueryRow(query, recordID)
	err := rows.Scan(&er.RecordID, &er.OwnerID, &er.Title, &er.TotalAmount,
		&er.PaidAmount, &er.InstallmentAmount, &er.StartDate, &er.EndDate, &er.DeductDay)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrorRecordNotFound
	} else if err != nil {
		log.Printf("Error scanning emiRecords row: %v", err)
		return ErrorServerError
	}

	return nil
}

func InsertEMIRecord(userID int, er *models.EMIRecord) error {
	// Insert new EMI Record in emiRecords
	query := `INSERT INTO
						emiRecords (ownerID, title, totalAmount, paidAmount, installmentAmount, startDate, endDate, deductDay)
						VALUES(?, ?, ?, ?, ?, ?, ?, ?);`

	_, err := database.DB.Exec(query, userID, er.Title, er.TotalAmount, er.PaidAmount, er.InstallmentAmount, er.StartDate, er.EndDate, er.DeductDay)
	if err != nil {
		log.Printf("Database error: %v\n", err)
		return ErrorServerError
	}

	// update Number of Active EMI in user
	query = `UPDATE users SET activeEMI = activeEMI + 1 WHERE userID = ?`
	_, err = database.DB.Exec(query, userID)
	if err != nil {
		log.Printf("Database error: %v\n", err)
		return ErrorServerError
	}

	return nil
}

func UpdateUserForEMIChange(userID int, totalAmountChange, paidAmountChange int) error {
	query := `UPDATE users
						SET totalLoaned = totalLoaned + ?,
								totalPaid = totalPaid + ?
						WHERE userID = ?`

	_, err := database.DB.Exec(query, totalAmountChange, paidAmountChange, userID)
	if err != nil {
		log.Printf("Database error: %v\n", err)
		return ErrorServerError
	}
	return nil
}

func UpdateUserForInstallment(userID int, paidAmountChange int) error {
	query := `UPDATE users
						SET totalPaid = totalPaid + ?
						WHERE userID = ?`

	_, err := database.DB.Exec(query, paidAmountChange, userID)
	if err != nil {
		log.Printf("Database error: %v\n", err)
		return ErrorServerError
	}
	return nil
}

func UpdateEMIRecord(recordID int, er *models.EMIRecord) error {
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

	_, err := database.DB.Exec(query, er.OwnerID, er.Title, er.TotalAmount, er.PaidAmount,
		er.InstallmentAmount, er.StartDate, er.EndDate, er.DeductDay, recordID)
	if err != nil {
		log.Printf("Database error: %v\n", err)
		return ErrorServerError
	}
	return nil
}

func DeleteEMIRecord(userID, recordID int, reduceComplete bool) error {
	if reduceComplete {
		// update Number of Completed EMI in user
		query := `UPDATE users SET completedEMI = completedEMI - 1 WHERE userID = ?`
		_, err := database.DB.Exec(query, userID)
		if err != nil {
			log.Printf("Database error: %v\n", err)
			return ErrorServerError
		}
	} else {
		// update Number of Active EMI in user
		query := `UPDATE users SET activeEMI = activeEMI - 1 WHERE userID = ?`
		_, err := database.DB.Exec(query, userID)
		if err != nil {
			log.Printf("Database error: %v\n", err)
			return ErrorServerError
		}
	}

	// delete emiRecord from database
	query := `DELETE FROM emiRecords WHERE recordID = ?`
	_, err := database.DB.Exec(query, recordID)
	if err != nil {
		log.Printf("Database error: %v\n", err)
		return ErrorServerError
	}

	return nil
}

func CompleteEMI(userID int) error {
	query := `UPDATE users
						SET completedEMI = completedEMI + 1,
								activeEMI = activeEMI - 1
						WHERE userID = ?`

	_, err := database.DB.Exec(query, userID)
	if err != nil {
		log.Printf("Database error: %v\n", err)
		return ErrorServerError
	}
	return nil
}
