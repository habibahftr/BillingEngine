package dto_out

type OutstandingOut struct {
	CustomerID     int64   `json:"customer_id"`
	CustomerName   string  `json:"customer_name"`
	RemainingWeeks int     `json:"remaining_weeks"`
	Paid           float64 `json:"paid"`
	Outstanding    float64 `json:"outstanding"`
	IsDelinquent   bool    `json:"is_delinquent"`
}
