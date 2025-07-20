package models

import "time"

type Subscription struct {
	SubID         int       `json:"subID"`
	OwnerID       int       `json:"ownerID"`
	SubName       string    `json:"subName"`
	TotalAmount   int       `json:"totalAmount"`
	PaidAmount    int       `json:"paidAmount"`
	PaymentAmount int       `json:"paymentAmount"`
	StartDate     time.Time `json:"startDate"`
	EndDate       time.Time `json:"endDate"`
	DeductDay     int       `json:"deductDay"`
}
