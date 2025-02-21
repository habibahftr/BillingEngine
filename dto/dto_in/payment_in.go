package dto_in

type Payment struct {
	CustomerID int64   `json:"customer_id" binding:"required,min=1"`
	LoanID     int64   `json:"loan_id" binding:"required,min=1"`
	AmountPaid float64 `json:"amount_paid" binding:"required,min=1"`
}
