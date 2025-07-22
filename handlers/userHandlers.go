package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/sajidzamanme/emi-tracker/database"
	"github.com/sajidzamanme/emi-tracker/models"
	"github.com/sajidzamanme/emi-tracker/utils"
)

// JSON Response with all Users
func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	// Select all users from database
	rows, err := database.DB.Query("SELECT * FROM users")
	if err != nil {
		http.Error(w, "Internal Database Error", http.StatusInternalServerError)
		log.Printf("Database error: %v\n", err)
		return
	}
	defer rows.Close()

	// save queried users in users slice
	var users []models.User
	for rows.Next() {
		var u models.User
		err = rows.Scan(&u.UserID, &u.Name, &u.Email, &u.Password, &u.TotalLoaned, &u.TotalPaid, &u.CurrentlyLoaned, &u.CurrentlyPaid, &u.CompletedEMI)
		if err != nil {
			log.Printf("Error scanning user: %v", err)
			return
		}
		users = append(users, u)
	}

	if len(users) == 0 {
		http.Error(w, "No Users Found", http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// JSON Response with User (through userID)
func GetUserByID(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.PathValue("userID"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var u models.User
	err = utils.FindUserByUserID(userID, &u)
	if err != nil {
		log.Printf("Error scanning users row: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(u)
}

// Add User to Database
func PostUser(w http.ResponseWriter, r *http.Request) {
	// Save User from request body to u
	var u models.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Invalid User Details", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Hash the password
	var err error
	u.Password, err = utils.HashPassword(u.Password)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Password hashing failed. Error: %v", err)
		return
	}

	// Set default values
	u.TotalLoaned = 0
	u.TotalPaid = 0
	u.CurrentlyLoaned = 0
	u.CurrentlyPaid = 0
	u.CompletedEMI = 0

	// Insert user to database
	query := `INSERT INTO
						users(name, email, password, totalLoaned, totalPaid, currentlyLoaned, currentlyPaid, completedEMI)
						VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	res, err := database.DB.Exec(query, u.Name, u.Email, u.Password, u.TotalLoaned, u.TotalPaid, u.CurrentlyLoaned, u.CurrentlyPaid, u.CompletedEMI)
	if err != nil {
		http.Error(w, "Internal Database Error", http.StatusInternalServerError)
		log.Printf("Database error: %v\n", err)
		return
	}

	// Get id of new user
	id, err := res.LastInsertId()
	if err != nil {
		http.Error(w, "Internal Server Error:", http.StatusInternalServerError)
		log.Printf("Server Error: %v\n", err)
		return
	}

	fmt.Fprintln(w, "User Added. ID:", int(id))
}

// Update User in Database
func PutUser(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.PathValue("userID"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var u models.User
	err = utils.FindUserByUserID(userID, &u)
	if err != nil {
		log.Printf("Error scanning users row: %v", err)
	}

	// Overwrite the new info
	if err = json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Invalid Record Entry", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Update entry in database
	query := `UPDATE users
						SET name = ?,
								email = ?,
								password = ?,
								totalLoaned = ?,
								totalPaid = ?,
								currentlyLoaned = ?,
								currentlyPaid = ?,
								completedEMI = ?
						WHERE userID = ?;`

	_, err = database.DB.Exec(query, u.Name, u.Email, u.Password, u.TotalLoaned,
		u.TotalPaid, u.CurrentlyLoaned, u.CurrentlyPaid, u.CompletedEMI, userID)
	if err != nil {
		http.Error(w, "Internal Database Error", http.StatusInternalServerError)
		log.Printf("Database error: %v", err)
		return
	}

	fmt.Fprintln(w, "User Updated with ID:", userID)
}

// Delete User from Database
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.PathValue("userID"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Delete user from database
	query := `DELETE FROM users WHERE userID = ?;`

	if _, err := database.DB.Exec(query, userID); err != nil {
		http.Error(w, "Internal Database Error", http.StatusInternalServerError)
		log.Printf("Database error: %v", err)
		return
	}

	fmt.Fprintln(w, "User Deleted with ID:", userID)
}

// JSON Response with all EMIRecords added to an individual User
func GetAllRecordsByUserID(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.PathValue("userID"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Get all EMIRecords of the requested user
	query := `SELECT * FROM emiRecords WHERE ownerID = ?;`
	rows, err := database.DB.Query(query, userID)
	if err != nil {
		http.Error(w, "Internal Database Error", http.StatusInternalServerError)
		log.Printf("Database error: %v\n", err)
		return
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
			return
		}
		records = append(records, er)
	}

	if len(records) == 0 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "No records found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(records)
}

// Takes email and password, verifies with hashed password for login
func PostLogin(w http.ResponseWriter, r *http.Request) {
	// Parsing input from request body
	type loginInput struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var inputU loginInput
	err := json.NewDecoder(r.Body).Decode(&inputU)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("JSON Decoding failed")
		return
	}

	// Finding hashedpassword from database by email
	var hashedPassword string
	query := `SELECT password FROM users WHERE email = ?`
	err = database.DB.QueryRow(query, inputU.Email).Scan(&hashedPassword)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "Invalid email", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Internal Database Error", http.StatusInternalServerError)
		log.Println("QueryRow Error:", err)
		return
	}

	// If password is wrong, show error and return
	if !utils.CheckPassword(hashedPassword, inputU.Password) {
		http.Error(w, "Incorrect Password", http.StatusBadRequest)
		return
	}

	// Create Token and set in cookies

	fmt.Fprintln(w, "Login Successful")
}
