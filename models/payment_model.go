package models

import (
	"database/sql"
)

type PaymentModel struct {
	ID             sql.NullInt64
	CustomerID     sql.NullInt64
	LoanID         sql.NullInt64
	AmountPaid     sql.NullFloat64
	PaymentDate    sql.NullTime
	CreatedAt      sql.NullTime
	LoanScheduleID sql.NullInt16
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

func UpdateLoanStatus(tx *sql.Tx, loanModel LoanModel) (err error) {
	query :=
		`
			UPDATE 
			    loans 
			SET out_standing = $1, is_delinquent = $2 WHERE id = $3 `

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

func UpdateLoanScheduleStatus(tx *sql.Tx, loanModel PaymentModel) (err error) {
	query :=
		`
			UPDATE 
			    loans_schedule
			SET is_paid = TRUE 
			WHERE id = $1 `

	param := []interface{}{
		loanModel.ID.Int64}

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

func GetDelinquentID(db *sql.DB, userParam PaymentModel) (result []int64, err error) {
	query :=
		`
			SELECT 
			    ls.id 
			FROM 
			    loans_schedule ls 
			    LEFT JOIN loans l ON l.id = ls.loan_id 
			WHERE 
			    l.id = $1 AND l.customer_id = $2 AND ls.due_date <= CURRENT_DATE AND ls.is_paid = FALSE
`

	param := []interface{}{userParam.LoanID.Int64, userParam.CustomerID.Int64}
	rows, err := db.Query(query, param...)
	if err != nil {
		return
	}
	if rows != nil {
		defer func() {
			err = rows.Close()
			if err != nil {
				return
			}
		}()
		for rows.Next() {
			var id int64
			err = rows.Scan(
				&id)
			if err != nil {
				return
			}
			result = append(result, id)
		}
	} else {
		return
	}

	return
}
