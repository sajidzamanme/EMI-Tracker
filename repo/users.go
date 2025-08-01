package repo

import (
	"database/sql"
	"errors"
	"log"

	"github.com/sajidzamanme/emi-tracker/database"
	"github.com/sajidzamanme/emi-tracker/models"
)

var ErrorUserNotFound = errors.New("User Doesn't Exist")
var ErrorServerError = errors.New("Internal Server Error")

func FindUserByUserID(userID int, u *models.User) error {
	query := `SELECT * FROM users WHERE userID = ?`

	rows := database.DB.QueryRow(query, userID)
	err := rows.Scan(&u.UserID, &u.Name, &u.Email, &u.Password, &u.TotalLoaned, &u.TotalPaid, &u.ActiveEMI, &u.CompletedEMI)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrorUserNotFound
	} else if err != nil {
		log.Printf("Error scanning users row: %v", err)
		return ErrorServerError
	}
	return nil
}

func GetAllUsers() ([]models.User, error) {
	// Select all users from database
	rows, err := database.DB.Query("SELECT * FROM users")
	if err != nil {
		log.Printf("Database error: %v\n", err)
		return nil, ErrorServerError
	}
	defer rows.Close()

	// save queried users in users slice
	var users []models.User
	for rows.Next() {
		var u models.User
		err = rows.Scan(&u.UserID, &u.Name, &u.Email, &u.Password, &u.TotalLoaned, &u.TotalPaid, &u.ActiveEMI, &u.CompletedEMI)
		if err != nil {
			log.Printf("Error scanning user: %v", err)
			return nil, ErrorServerError
		}
		users = append(users, u)
	}

	if len(users) == 0 {
		return nil, errors.New("No Users Found")
	}

	return users, nil
}

func InsertUser(u models.User) (int, error) {
	tx, err := database.DB.Begin()
	if err != nil {
		log.Printf("Transaction Begin Error: %v", err)
		return -1, ErrorServerError
	}

	query := `INSERT INTO
						users(name, email, password, totalLoaned, totalPaid, activeEMI, completedEMI)
						VALUES (?, ?, ?, ?, ?, ?, ?)`

	res, err := tx.Exec(query, u.Name, u.Email, u.Password, u.TotalLoaned, u.TotalPaid, u.ActiveEMI, u.CompletedEMI)
	if err != nil {
		tx.Rollback()
		log.Printf("Database error: %v\n", err)
		return -1, ErrorServerError
	}

	// Get id of new user
	id, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		log.Printf("Server Error: %v\n", err)
		return -1, ErrorServerError
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("Commit Error: %v", err)
		return -1, ErrorServerError
	}

	return int(id), nil
}

func UpdateUser(u models.User) error {
	query := `UPDATE users
						SET name = ?,
								email = ?,
								password = ?,
								totalLoaned = ?,
								totalPaid = ?,
								activeEMI = ?,
								completedEMI = ?
						WHERE userID = ?;`

	_, err := database.DB.Exec(query, u.Name, u.Email, u.Password, u.TotalLoaned,
		u.TotalPaid, u.ActiveEMI, u.CompletedEMI, u.UserID)
	if err != nil {
		log.Printf("Database error: %v", err)
		return ErrorServerError
	}
	return nil
}

func DeleteUser(userID int) error {
	query := `DELETE FROM users WHERE userID = ?;`

	if _, err := database.DB.Exec(query, userID); err != nil {
		log.Printf("Database error: %v", err)
		return ErrorServerError
	}
	return nil
}

func GetAllEMIRecordByUserID(userID int) ([]models.EMIRecord, error) {
	// Select all emiRecords of a user from database
	query := `SELECT * FROM emiRecords WHERE ownerID = ?;`
	rows, err := database.DB.Query(query, userID)
	if err != nil {
		log.Printf("Database error: %v\n", err)
		return nil, ErrorServerError
	}
	defer rows.Close()

	// Save the EMIRecords in records slice
	var records []models.EMIRecord
	for rows.Next() {
		var er models.EMIRecord
		err = rows.Scan(&er.RecordID, &er.OwnerID, &er.Title, &er.TotalAmount,
			&er.PaidAmount, &er.InstallmentAmount, &er.StartDate, &er.EndDate, &er.DeductDay)
		if err != nil {
			log.Printf("Error scanning record: %v", err)
			return nil, ErrorServerError
		}
		records = append(records, er)
	}

	if len(records) == 0 {
		return nil, errors.New("No EMI Records found")
	}

	return records, nil
}

func GetHashedPasswordByEmail(email string) (string, error) {
	var hashedPassword string
	query := `SELECT password FROM users WHERE email = ?`
	err := database.DB.QueryRow(query, email).Scan(&hashedPassword)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrorUserNotFound
	} else if err != nil {
		log.Println("QueryRow Error:", err)
		return "", ErrorServerError
	}
	return hashedPassword, nil
}
