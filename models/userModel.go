package models

type User struct {
	UserID         int    `json:"userID"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	TotalLoaned    int    `json:"totalLoaned"`
	TotalPaid      int    `json:"totalPaid"`
}
