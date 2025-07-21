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
func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT * FROM users")
	if err != nil {
		log.Fatalln(err)
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		rows.Scan(&u.UserID, &u.Name, &u.Email, &u.Password, &u.TotalLoaned, &u.TotalPaid)
		users = append(users, u)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// DONE
func GetUserByID(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.PathValue("userID"))
	if err != nil {
		fmt.Fprintln(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	query := `SELECT * FROM users WHERE userID = ?`

	rows, err := database.DB.Query(query, userID)
	if err != nil {
		log.Fatalln(err)
	}
	defer rows.Close()

	if !rows.Next() {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	var u models.User
	rows.Scan(&u.UserID, &u.Name, &u.Email, &u.Password, &u.TotalLoaned, &u.TotalPaid)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(u)
}

// DONE
func PostUser(w http.ResponseWriter, r *http.Request) {
	var u models.User

	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Invalid User Details", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO users(
		name, email, pass, totalLoaned, totalPaid
		)
		VALUES (?, ?, ?, ?, ?)`

	// error here
	res, err := database.DB.Exec(query, u.Name, u.Email, u.Password, u.TotalLoaned, u.TotalPaid)
	if err != nil {
		http.Error(w, "Internal Database Error", http.StatusInternalServerError)
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		http.Error(w, "Internal Server Error:", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "User Added. ID:", int(id))
}

// DONE
func PutUser(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.PathValue("userID"))
	if err != nil {
		fmt.Fprintln(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	query := `SELECT * FROM users WHERE userID = ?`

	rows, err := database.DB.Query(query, userID)
	if err != nil {
		log.Fatalln(err)
	}
	defer rows.Close()

	// check if found
	rows.Next()
	var u models.User
	rows.Scan(&u.UserID, &u.Name, &u.Email, &u.Password, &u.TotalLoaned, &u.TotalPaid)

	if err = json.NewDecoder(r.Body).Decode(&u); err != nil {
		log.Fatalln(err)
	}

	query = `UPDATE users
						SET name = ?,
  					email = ?,
  					pass = ?,
						totalLoaned = ?,
						totalPaid = ?
						WHERE userID = ?;`

	if _, err := database.DB.Exec(query, u.Name, u.Email, u.Password, u.TotalLoaned, u.TotalPaid, userID); err != nil {
		http.Error(w, "Database Error", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "User Updated with ID:", userID)
}

// DONE
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.PathValue("userID"))
	if err != nil {
		fmt.Fprintln(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	query := `DELETE FROM users WHERE userID = ?;`

	if _, err := database.DB.Exec(query, userID); err != nil {
		http.Error(w, "Internal Database Error", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "User Deleted with ID:", userID)
}

// DONE
func GetAllRecordsByUserID(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.PathValue("userID"))
	if err != nil {
		fmt.Fprintln(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	query := `SELECT * FROM emiRecords WHERE ownerID = ?;`

	rows, err := database.DB.Query(query, userID)
	if err != nil {
		log.Fatalln(err)
	}
	defer rows.Close()

	var records []models.EMIRecord
	for rows.Next() {
		var er models.EMIRecord
		rows.Scan(&er.RecordID, &er.OwnerID, &er.Title, &er.TotalAmount,
			&er.PaidAmount, &er.InstallmentAmount, &er.StartDate, &er.EndDate, &er.DeductDay)
		records = append(records, er)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(records)
}
