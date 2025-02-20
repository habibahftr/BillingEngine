package dto_in

type Loan struct {
	CustomerID   int64   `json:"customer_id" binding:"required,min=1"`
	LoanAmount   float64 `json:"loan_amount" binding:"required,min=100000"`
	InterestRate float64 `json:"interest_rate" binding:"required,min=1,max=30"`
	TotalWeeks   int     `json:"total_weeks" binding:"required,min=1,max=50"`
}
