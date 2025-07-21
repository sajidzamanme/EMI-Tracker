package models

import "time"

type EMIRecord struct {
	RecordID          int       `json:"recordID"`
	OwnerID           int       `json:"ownerID"`
	Title             string    `json:"title"`
	TotalAmount       int       `json:"totalAmount"`
	PaidAmount        int       `json:"paidAmount"`
	InstallmentAmount int       `json:"installmentAmount"`
	StartDate         time.Time `json:"startDate"`
	EndDate           time.Time `json:"endDate"`
	DeductDay         int       `json:"deductDay"`
}
