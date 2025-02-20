package controllers

import (
	"BillingEngine/dto/dto_in"
	"BillingEngine/models"
	"BillingEngine/service/loan_service"
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func CreateLoan(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var loan dto_in.Loan
		err := ctx.ShouldBindJSON(&loan)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		loanModel := loan_service.CreateLoan(loan)

		var tx *sql.Tx
		defer func() {
			if err != nil {
				if err = tx.Rollback(); err != nil {
					return
				}
			} else {
				if err = tx.Commit(); err != nil {
					return
				}
			}
		}()
		tx, err = db.Begin()
		if err != nil {
			return
		}

		loanOnDB, err := models.GetLoanByCustomerID(loanModel, db)
		if err != nil {
			return
		}

		if loanOnDB.ID.Int64 != 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "You still have outstanding loan"})
			return
		}

		id, err := models.InsertLoan(tx, loanModel)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create loan"})
			return
		}

		loanModel.ID.Int64 = id
		err = models.InsertLoanSchedule(tx, loanModel)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create loan"})
			return
		}

		message := "Loan created successfully with id " + strconv.Itoa(int(id))
		ctx.JSON(http.StatusOK, gin.H{"message": message})
	}
}

func GetOutStanding(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		custId := ctx.Param("cust_id")
		customerId, err := strconv.Atoi(custId)
		if err != nil || customerId == 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid format customer id"})
			return
		}

		userParam := models.LoanModel{
			CustomerID: sql.NullInt64{Int64: int64(customerId)},
		}

		resultOnDb, err := models.GetOutstandingInfo(userParam, db)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get outstanding"})
			return
		}

		if resultOnDb.CustomerID.Int64 == 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "unknown data with this ID"})
			return
		}

		resultDelinquent, err := models.IsDelinquentCustomer(userParam, db)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get is delinquent status "})
			return
		}

		resultOnDb.IsDelinquent = resultDelinquent.IsDelinquent
		outstandingResponse := loan_service.ReformatResponseOutstanding(resultOnDb)

		ctx.JSON(http.StatusOK, gin.H{"outstanding_balance": outstandingResponse})
	}
}

func PaymentLoan(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		custId := ctx.Param("cust_id")
		customerId, err := strconv.Atoi(custId)
		if err != nil || customerId == 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid format customer id"})
			return
		}

		var payment dto_in.Payment
		err = ctx.ShouldBindJSON(&payment)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userParam := models.PaymentModel{
			CustomerID:  sql.NullInt64{Int64: int64(customerId)},
			LoanID:      sql.NullInt64{Int64: payment.LoanID},
			PaymentDate: sql.NullTime{Time: payment.PaymentDate},
			AmountPaid:  sql.NullFloat64{Float64: payment.AmountPaid},
		}

		//loanOnDB, err := models.GetLoanByCustomerID(models.LoanModel{
		//	CustomerID: userParam.CustomerID,
		//}, db)
		//if err != nil {
		//	return
		//}
		//
		//if loanOnDB.ID.Int64 != 0 {
		//	ctx.JSON(http.StatusBadRequest, gin.H{"error": "You still have outstanding loan"})
		//	return
		//}

		resultDelinquent, err := models.IsDelinquentCustomer(models.LoanModel{
			CustomerID: userParam.CustomerID,
		}, db)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get is delinquent status "})
			return
		}

		if resultDelinquent.IsDelinquent.Bool {
			if userParam.AmountPaid.Float64 != resultDelinquent.AmountDue.Float64 {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "Payment amount must match the amount due"})
				return
			}
		}

		var tx *sql.Tx
		defer func() {
			if err != nil {
				if err = tx.Rollback(); err != nil {
					return
				}
			} else {
				if err = tx.Commit(); err != nil {
					return
				}
			}
		}()
		tx, err = db.Begin()
		if err != nil {
			return
		}

	}

}
