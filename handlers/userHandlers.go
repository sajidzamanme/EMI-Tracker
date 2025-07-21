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
		http.Error(w, "Internal Database Error", http.StatusInternalServerError)
		log.Printf("Database error: %v\n", err)
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		err = rows.Scan(&u.UserID, &u.Name, &u.Email, &u.Password, &u.TotalLoaned, &u.TotalPaid)
		if err != nil {
			log.Printf("Error scanning user: %v", err)
			return
		}
		users = append(users, u)
	}

	if len(users) == 0 {
		w.WriteHeader(http.StatusNoContent)
		fmt.Println("No users found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// DONE
func GetUserByID(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.PathValue("userID"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	query := `SELECT * FROM users WHERE userID = ?`

	rows := database.DB.QueryRow(query, userID)
	var u models.User
	err = rows.Scan(&u.UserID, &u.Name, &u.Email, &u.Password, &u.TotalLoaned, &u.TotalPaid)
	if err != nil {
		log.Printf("Error scanning user: %v", err)
		return
	}

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
	defer r.Body.Close()

	u.TotalLoaned = 0
	u.TotalPaid = 0

	query := `INSERT INTO users(
		name, email, pass, totalLoaned, totalPaid
		)
		VALUES (?, ?, ?, ?, ?)`

	res, err := database.DB.Exec(query, u.Name, u.Email, u.Password, u.TotalLoaned, u.TotalPaid)
	if err != nil {
		http.Error(w, "Internal Database Error", http.StatusInternalServerError)
		log.Printf("Database error: %v\n", err)
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		http.Error(w, "Internal Server Error:", http.StatusInternalServerError)
		log.Printf("Server Error: %v", err)
		return
	}

	fmt.Fprintln(w, "User Added. ID:", int(id))
}

// DONE
func PutUser(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.PathValue("userID"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	query := `SELECT * FROM users WHERE userID = ?`

	rows := database.DB.QueryRow(query, userID)
	var u models.User
	err = rows.Scan(&u.UserID, &u.Name, &u.Email, &u.Password, &u.TotalLoaned, &u.TotalPaid)
	if err != nil {
		log.Printf("Error scanning user: %v", err)
		return
	}

	if err = json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Invalid Record Entry", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	query = `UPDATE users
						SET name = ?,
  					email = ?,
  					pass = ?,
						totalLoaned = ?,
						totalPaid = ?
						WHERE userID = ?;`

	if _, err := database.DB.Exec(query, u.Name, u.Email, u.Password, u.TotalLoaned, u.TotalPaid, userID); err != nil {
		http.Error(w, "Internal Database Error", http.StatusInternalServerError)
		log.Printf("Database error: %v", err)
		return
	}

	fmt.Fprintln(w, "User Updated with ID:", userID)
}

// DONE
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.PathValue("userID"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	query := `DELETE FROM users WHERE userID = ?;`

	if _, err := database.DB.Exec(query, userID); err != nil {
		http.Error(w, "Internal Database Error", http.StatusInternalServerError)
		log.Printf("Database error: %v", err)
		return
	}

	fmt.Fprintln(w, "User Deleted with ID:", userID)
}

// DONE
func GetAllRecordsByUserID(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.PathValue("userID"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	query := `SELECT * FROM emiRecords WHERE ownerID = ?;`

	rows, err := database.DB.Query(query, userID)
	if err != nil {
		http.Error(w, "Internal Database Error", http.StatusInternalServerError)
		log.Printf("Database error: %v\n", err)
		return
	}
	defer rows.Close()

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
		w.WriteHeader(http.StatusNoContent)
		fmt.Println("No records found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(records)
}
