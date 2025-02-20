package loan_service

import (
	"BillingEngine/dto/dto_in"
	"BillingEngine/dto/dto_out"
	"BillingEngine/models"
	"database/sql"
	"time"
)

func CreateLoan(loan dto_in.Loan) (loanModel models.LoanModel) {

	totalPayable := loan.LoanAmount + (loan.LoanAmount * (loan.InterestRate / 100))
	weeklyPayment := totalPayable / float64(loan.TotalWeeks)
	startDate := time.Now()
	endDate := startDate.AddDate(0, 0, loan.TotalWeeks*7)

	loanModel = models.LoanModel{
		CustomerID:    sql.NullInt64{Int64: loan.CustomerID},
		LoanAmount:    sql.NullFloat64{Float64: loan.LoanAmount},
		InterestRate:  sql.NullFloat64{Float64: loan.InterestRate},
		TotalWeeks:    sql.NullInt16{Int16: int16(loan.TotalWeeks)},
		TotalPayable:  sql.NullFloat64{Float64: totalPayable},
		WeeklyPayment: sql.NullFloat64{Float64: weeklyPayment},
		Outstanding:   sql.NullFloat64{Float64: totalPayable},
		StartDate:     sql.NullTime{Time: startDate},
		EndDate:       sql.NullTime{Time: endDate},
		CreateAt:      sql.NullTime{Time: startDate},
	}

	return
}

func ReformatResponseOutstanding(model models.LoanModel) (response dto_out.OutstandingOut) {
	response = dto_out.OutstandingOut{
		CustomerID:     model.CustomerID.Int64,
		CustomerName:   model.CustomerName.String,
		RemainingWeeks: int(model.RemainingWeeks.Int16),
		Paid:           model.Paid.Float64,
		Outstanding:    model.Outstanding.Float64,
		IsDelinquent:   model.IsDelinquent.Bool,
	}
	return
}
