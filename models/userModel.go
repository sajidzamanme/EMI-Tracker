package models

type User struct {
	UserID          int    `json:"userID"`
	Name            string `json:"name"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	TotalLoaned     int    `json:"totalLoaned"`
	TotalPaid       int    `json:"totalPaid"`
	CurrentlyLoaned int    `json:"currentlyLoaned"`
	CurrentlyPaid   int    `json:"currentlyPaid"`
	CompletedEMI    int    `json:"completedEMI"`
}
