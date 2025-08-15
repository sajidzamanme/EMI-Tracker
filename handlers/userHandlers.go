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

// JSON Response with all Users
func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := repo.GetAllUsers()
	if err != nil {
		fmt.Fprintln(w, err.Error())
		return
	}

	err = utils.EncodeJson(w, users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// JSON Response with User (through userID)
func GetUserByID(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.PathValue("userID"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var u models.User
	err = repo.FindUserByUserID(userID, &u)
	if errors.Is(err, repo.ErrorUserNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = utils.EncodeJson(w, u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Add User to Database
func InsertUser(w http.ResponseWriter, r *http.Request) {
	// Save User from request body to u
	var u models.User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("JSON Decoding failed")
		return
	}
	defer r.Body.Close()

	// Hash the password
	u.Password, err = utils.HashPassword(u.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set default values
	u.TotalLoaned = 0
	u.TotalPaid = 0
	u.ActiveEMI = 0
	u.CompletedEMI = 0

	// Insert user to database
	id, err := repo.InsertUser(u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "User Added. ID:", id)
}

// Update User in Database
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.PathValue("userID"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var u models.User
	err = repo.FindUserByUserID(userID, &u)
	if errors.Is(err, repo.ErrorUserNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Overwrite the new info
	if err = json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("JSON Decoding failed")
		return
	}
	defer r.Body.Close()

	// Update entry in database
	err = repo.UpdateUser(u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

	err = repo.DeleteUser(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
	records, err := repo.GetAllEMIRecordByUserID(userID)
	if err != nil {
		fmt.Fprintln(w, err.Error())
		return
	}

	err = utils.EncodeJson(w, records)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Takes email and password, verifies with hashed password for login
func UserLogin(w http.ResponseWriter, r *http.Request) {
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
	hashedPassword, err := repo.GetHashedPasswordByEmail(inputU.Email)
	if errors.Is(err, repo.ErrorUserNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
