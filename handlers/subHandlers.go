package handlers

import (
	"fmt"
	"net/http"
	"strconv"
)

func GetAllSubsByUserID(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.PathValue("userID"))
	if err != nil {
		fmt.Fprintln(w, "Invalid ID", http.StatusBadRequest)
	}

	fmt.Fprintln(w, "All Subscriptions List of ID:", userID)
}

func GetSubByID(w http.ResponseWriter, r *http.Request) {
	subID, err := strconv.Atoi(r.PathValue("subID"))
	if err != nil {
		fmt.Fprintln(w, "Invalid ID", http.StatusBadRequest)
	}

	fmt.Fprintln(w, "Subscription details of ID:", subID)
}

func PostSubByUserID(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.PathValue("userID"))
	if err != nil {
		fmt.Fprintln(w, "Invalid ID", http.StatusBadRequest)
	}
	
	fmt.Fprintln(w, "Subscription Added to ID:", userID)
}

func PutSubBySubID(w http.ResponseWriter, r *http.Request) {
	subID, err := strconv.Atoi(r.PathValue("subID"))
	if err != nil {
		fmt.Fprintln(w, "Invalid ID", http.StatusBadRequest)
	}
	
	fmt.Fprintln(w, "Subscription Updated of ID:", subID)
}

func DeleteSubBySubID(w http.ResponseWriter, r *http.Request) {
	subID, err := strconv.Atoi(r.PathValue("subID"))
	if err != nil {
		fmt.Fprintln(w, "Invalid ID", http.StatusBadRequest)
	}
	
	fmt.Fprintln(w, "Subscription Deleted with ID:", subID)
}
