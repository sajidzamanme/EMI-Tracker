package repo

import (
	"log"

	"github.com/sajidzamanme/emi-tracker/database"
	"github.com/sajidzamanme/emi-tracker/models"
)

// Find record from databse and set to sent pointer
func FindRecordByRecordID(recordID int, er *models.EMIRecord) error {
	query := `SELECT * FROM emiRecords WHERE recordID = ?;`

	rows := database.DB.QueryRow(query, recordID)
	err := rows.Scan(&er.RecordID, &er.OwnerID, &er.Title, &er.TotalAmount,
		&er.PaidAmount, &er.InstallmentAmount, &er.StartDate, &er.EndDate, &er.DeductDay)
	if err != nil {
		log.Printf("Database error: %v\n", err)
		return err
	}

	return nil
}

func InsertEMIRecord(userID int, er *models.EMIRecord) error {
	query := `INSERT INTO
						emiRecords (ownerID, title, totalAmount, paidAmount, installmentAmount, startDate, endDate, deductDay)
						VALUES(?, ?, ?, ?, ?, ?, ?, ?);`

	_, err := database.DB.Exec(query, userID, er.Title, er.TotalAmount, er.PaidAmount, er.InstallmentAmount, er.StartDate, er.EndDate, er.DeductDay)
	if err != nil {
		log.Printf("Database error: %v\n", err)
		return err
	}
	return nil
}

func UpdateUserForEMIChange(userID int, totalAmount, paidAmount int) error {
	query := `UPDATE users
						SET totalLoaned = totalLoaned + ?,
							totalPaid = totalPaid + ?,
							currentlyLoaned = currentlyLoaned + ?,
							currentlyPaid = currentlyPaid + ?
						WHERE userID = ?`

	_, err := database.DB.Exec(query, totalAmount, paidAmount, totalAmount, paidAmount, userID)
	if err != nil {
		log.Printf("Database error: %v\n", err)
		return err
	}
	return nil
}

func UpdateUserForInstallment(userID int, totalAmount, paidAmount int) error {
	query := `UPDATE users
						SET totalPaid = totalPaid + ?,
								currentlyLoaned = currentlyLoaned + ?,
								currentlyPaid = currentlyPaid + ?
						WHERE userID = ?`

	_, err := database.DB.Exec(query, paidAmount, totalAmount, paidAmount, userID)
	if err != nil {
		log.Printf("Database error: %v\n", err)
		return err
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
		return err
	}
	return nil
}

func DeleteEMIRecord(recordID int) error {
	query := `DELETE FROM emiRecords WHERE recordID = ?`

	_, err := database.DB.Exec(query, recordID)
	if err != nil {
		log.Printf("Database error: %v\n", err)
		return err
	}
	return nil
}

func CompleteEMI(ownerID int) error {
	query := `UPDATE users
						SET completedEMI = completedEMI + 1
						WHERE userID = ?`

	_, err := database.DB.Exec(query, ownerID)
	if err != nil {
		log.Printf("Database error: %v\n", err)
		return err
	}
	return nil
}
