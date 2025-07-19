package models

import "time"

type Subscription struct {
	SubID           int       `json:"subID"`
	OwnerID         int       `json:"ownerID"`
	SubName         string    `json:"subName"`
	TotalAmount     int       `json:"totalAmount"`
	PaidAmount      int       `json:"paidAmount"`
	RemainingAmount int       `json:"remainingAmount"`
	PaymentAmount   int       `json:"paymentAmount"`
	StartDate       time.Time `json:"startDate"` // time format?
	EndDate         time.Time `json:"endDate"`
	DeductDay       int       `json:"deductDate"`
}
