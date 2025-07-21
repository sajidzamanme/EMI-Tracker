package utils

import (
	"github.com/sajidzamanme/emi-tracker/database"
	"github.com/sajidzamanme/emi-tracker/models"
	"golang.org/x/crypto/bcrypt"
)

func FindUserByUserID(userID int, u *models.User) error {
	query := `SELECT * FROM users WHERE userID = ?`

	rows := database.DB.QueryRow(query, userID)
	err := rows.Scan(&u.UserID, &u.Name, &u.Email, &u.Password, &u.TotalLoaned, &u.TotalPaid, &u.CurrentlyLoaned, &u.CurrentlyPaid, &u.CompletedEMI)
	if err != nil {
		return err
	}

	return nil
}

func FindRecordByRecordID(recordID int, er *models.EMIRecord) error {
	query := `SELECT * FROM emiRecords WHERE recordID = ?;`

	rows := database.DB.QueryRow(query, recordID)
	err := rows.Scan(&er.RecordID, &er.OwnerID, &er.Title, &er.TotalAmount,
		&er.PaidAmount, &er.InstallmentAmount, &er.StartDate, &er.EndDate, &er.DeductDay)
	if err != nil {
		return err
	}

	return nil
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
}
