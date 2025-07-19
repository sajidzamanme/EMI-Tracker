package models

import "time"

type Subscription struct {
	SubID           int       `json:"subID"`
	TotalAmount     int       `json:"totalAmount"`
	PaidAmount      int       `json:"paidAmount"`
	RemainingAmount int       `json:"remainingAmount"`
	StartDate       time.Time `json:"startDate"` // time format?
	EndDate         time.Time `json:"endDate"`
	DeductDate      time.Time `json:"deductDate"`
}
