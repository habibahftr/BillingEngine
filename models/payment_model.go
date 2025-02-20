package models

import (
	"database/sql"
)

type PaymentModel struct {
	ID          sql.NullInt64
	CustomerID  sql.NullInt64
	LoanID      sql.NullInt64
	AmountPaid  sql.NullFloat64
	PaymentDate sql.NullTime
	CreatedAt   sql.NullTime
}

func InsertPayment(tx *sql.Tx, paymentModel PaymentModel) (id int64, err error) {
	query :=
		`INSERT INTO
			payments (
		loan_id, payment_date, amount_paid,
		created_at
		)
		VALUES (
		$1, $2, $3,
		$4
		) RETURNING id `

	param := []interface{}{
		paymentModel.LoanID.Int64, paymentModel.PaymentDate.Time, paymentModel.AmountPaid.Float64,
		paymentModel.CreatedAt.Time,
	}

	err = tx.QueryRow(query, param...).Scan(&id)
	if err != nil {
		if err = tx.Rollback(); err != nil {
			return
		}
		return
	}

	return
}

func GetPaymentByWeek(userParam PaymentModel, db *sql.DB) (paymentModel PaymentModel, err error) {
	query :=
		` 
			SELECT 
			    id, loan_id, amount_paid
			FROM 
			    payments WHERE loan_id = $1 `
	param := []interface{}{userParam.LoanID.Int64}
	err = db.QueryRow(query, param...).Scan(
		&paymentModel.ID, &paymentModel.LoanID, &paymentModel.AmountPaid,
	)
	return
}

// UpdateLoanStatus updates the outstanding balance and delinquency status for a loan
func UpdateLoanStatus(tx *sql.Tx, loanModel LoanModel) (err error) {
	query :=
		`
			UPDATE 
			    loans 
			SET outstanding = $1, is_delinquent = $2 WHERE id = $3 `

	param := []interface{}{
		loanModel.Outstanding.Float64, loanModel.IsDelinquent.Bool, loanModel.ID.Int64}

	stmt, err := tx.Prepare(query)
	if err != nil {
		return
	}

	_, err = stmt.Exec(param...)

	if err != nil {
		return
	}
	return
}
