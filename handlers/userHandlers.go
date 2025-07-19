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

// DONE (convert to json maybe)
func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT * FROM users")
	if err != nil {
		log.Fatalln(err)
	}
	defer rows.Close()

	for rows.Next() {
		var u models.User
		rows.Scan(&u.UserID, &u.Name, &u.Email, &u.Password, &u.TotalLoaned, &u.TotalPaid)
		fmt.Fprintf(w, "User %s (ID: %v) has total loan %v, of which he paid %v. Remaining amount: %v\n", u.Name, u.UserID, u.TotalLoaned, u.TotalPaid, u.TotalLoaned-u.TotalPaid)
	}
}

// DONE (JSON?)
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

	// check if found
	rows.Next()
	var u models.User
	rows.Scan(&u.UserID, &u.Name, &u.Email, &u.Password, &u.TotalLoaned, &u.TotalPaid)

	fmt.Fprintf(w, "User details of ID %v:\n", userID)
	fmt.Fprintf(w,
		"Name: %s\nEmail: %s\nTotal Loan: %v\nTotal Paid: %v\nRemaining Amount: %v",
		u.Name, u.Email, u.TotalLoaned, u.TotalPaid, u.TotalLoaned-u.TotalPaid)
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
