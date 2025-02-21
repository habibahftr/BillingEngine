package models

import "database/sql"

func GetCustomerByID(userParam LoanModel, db *sql.DB) (result LoanModel, err error) {
	query :=
		` 
			SELECT 
			    id
			FROM 
			    customer 
			WHERE deleted = FALSE AND id = $1 `

	param := []interface{}{userParam.CustomerID.Int64}
	err = db.QueryRow(query, param...).Scan(
		&result.ID,
	)

	if err != nil && err != sql.ErrNoRows {
		return
	}
	err = nil
	return
}
