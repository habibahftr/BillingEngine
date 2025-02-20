package dto_in

import "time"

type Payment struct {
	CustomerID  int64     `json:"customer_id"`
	LoanID      int64     `json:"loan_id" binding:"required,min=1"`
	PaymentDate time.Time `json:"payment_date" binding:"required"`
	AmountPaid  float64   `json:"amount_paigitd" binding:"required,min=1"`
}
