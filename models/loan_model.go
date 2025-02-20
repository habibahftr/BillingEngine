package models

import (
	"database/sql"
	"fmt"
)

type LoanModel struct {
	ID             sql.NullInt64
	CustomerID     sql.NullInt64
	LoanAmount     sql.NullFloat64
	InterestRate   sql.NullFloat64
	TotalWeeks     sql.NullInt16
	TotalPayable   sql.NullFloat64
	WeeklyPayment  sql.NullFloat64
	Outstanding    sql.NullFloat64
	StartDate      sql.NullTime
	EndDate        sql.NullTime
	CreateAt       sql.NullTime
	IsDelinquent   sql.NullBool
	Paid           sql.NullFloat64
	CustomerName   sql.NullString
	RemainingWeeks sql.NullInt16
	AmountDue      sql.NullFloat64
}

func InsertLoan(tx *sql.Tx, loanModel LoanModel) (id int64, err error) {
	query :=
		`INSERT INTO
			loans (
		customer_id, loan_amount, interest_rate,
		total_weeks, total_payable, weekly_payment,
		out_standing, start_date, end_date,
		created_at
		)
		VALUES (
		$1, $2, $3,
		$4, $5, $6,
		$7, $8, $9,
		$10
		) RETURNING id `

	param := []interface{}{
		loanModel.CustomerID.Int64, loanModel.LoanAmount.Float64, loanModel.InterestRate.Float64,
		loanModel.TotalWeeks.Int16, loanModel.TotalPayable.Float64, loanModel.WeeklyPayment.Float64,
		loanModel.Outstanding.Float64, loanModel.StartDate.Time, loanModel.EndDate.Time,
		loanModel.CreateAt.Time,
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

func InsertLoanSchedule(tx *sql.Tx, loanModel LoanModel) (err error) {
	query := `
		INSERT INTO 
		    loans_schedule
		(
		 	loan_id, week_number, due_date, 
		 	amount_due
		)
		VALUES (
		        $1, $2, $3, 
		        $4
		)`

	for week := 1; week <= int(loanModel.TotalWeeks.Int16); week++ {
		dueDate := loanModel.StartDate.Time.AddDate(0, 0, 7*week)

		param := []interface{}{
			loanModel.ID.Int64, week, dueDate,
			loanModel.WeeklyPayment.Float64,
		}

		stmt, errS := tx.Prepare(query)
		if errS != nil {
			err = errS
			_ = tx.Rollback()
			return
		}

		_, err = stmt.Exec(param...)
		if err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("failed to insert schedule for week %d: %v", week, err)
		}

	}

	return nil
}

func GetLoanByCustomerID(userParam LoanModel, db *sql.DB) (result LoanModel, err error) {
	query :=
		` 
			SELECT 
			    id
			FROM 
			    loans 
			WHERE out_standing != 0 AND deleted = FALSE AND customer_id = $1 `

	param := []interface{}{userParam.CustomerID.Int64}
	err = db.QueryRow(query, param...).Scan(
		&result.ID,
	)
	return
}

func GetOutstandingInfo(userParam LoanModel, db *sql.DB) (result LoanModel, err error) {
	query :=
		` 
			SELECT 
				COUNT(ls.id) AS remaining_weeks, (l.total_payable - l.out_standing)as paid, l.out_standing,
				l.customer_id, c.name
			FROM loans_schedule ls
			LEFT JOIN loans l ON ls.loan_id = l.id AND ls.is_paid = FALSE
			LEFT JOIN customer c ON c.id = l.customer_id
			WHERE l.customer_id = $1
			GROUP BY l.out_standing, l.total_payable, l.customer_id, c.name `

	param := []interface{}{userParam.CustomerID.Int64}
	err = db.QueryRow(query, param...).Scan(
		&result.RemainingWeeks, &result.Paid, &result.Outstanding,
		&result.CustomerID, &result.CustomerName,
	)
	return
}

func IsDelinquentCustomer(userParam LoanModel, db *sql.DB) (result LoanModel, err error) {
	query :=
		` 
		SELECT 
			CASE WHEN COUNT(ls.id) >= 2 THEN TRUE
			ELSE FALSE
			END, 
		    SUM(ls.amount_due)
		FROM 
			loans_schedule ls
			LEFT JOIN loans l ON ls.loan_id = l.id AND ls.is_paid = FALSE
			LEFT JOIN customer c ON c.id = l.customer_id
		WHERE ls.due_date < CURRENT_DATE  AND l.deleted = FALSE AND l.customer_id = 1 `

	param := []interface{}{userParam.CustomerID.Int64}
	err = db.QueryRow(query, param...).Scan(
		&result.IsDelinquent, &result.AmountDue,
	)
	return
}
