package controllers

import (
	"BillingEngine/dto/dto_in"
	"BillingEngine/models"
	"BillingEngine/service/loan_service"
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
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

		custOnDB, err := models.GetCustomerByID(loanModel, db)
		if err != nil {
			return
		}

		if custOnDB.ID.Int64 != 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Unknown data with this customer id"})
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

		if resultDelinquent.IsDelinquent.Bool {
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

			err = models.UpdateLoanStatus(tx, resultOnDb)
			if err != nil {
				return
			}
		}

		outstandingResponse := loan_service.ReformatResponseOutstanding(resultOnDb)

		ctx.JSON(http.StatusOK, gin.H{"outstanding_balance": outstandingResponse})
	}
}

func PaymentLoan(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var payment dto_in.Payment
		err := ctx.ShouldBindJSON(&payment)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
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

		var timeNow = time.Now()
		userParam := models.PaymentModel{
			CustomerID:  sql.NullInt64{Int64: payment.CustomerID},
			LoanID:      sql.NullInt64{Int64: payment.LoanID},
			AmountPaid:  sql.NullFloat64{Float64: payment.AmountPaid},
			CreatedAt:   sql.NullTime{Time: timeNow},
			PaymentDate: sql.NullTime{Time: timeNow},
		}

		loanOnDB, err := models.GetLoanByCustomerAndLoanID(models.LoanModel{
			CustomerID: userParam.CustomerID,
			ID:         userParam.LoanID,
		}, db)
		if err != nil {
			return
		}

		if loanOnDB.ID.Int64 == 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Unknown data with this customer id or loan id"})
			return
		}

		resultDelinquent, err := models.IsDelinquentCustomer(models.LoanModel{
			CustomerID: userParam.CustomerID,
		}, db)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get is delinquent status "})
			return
		}

		var outstanding float64
		if resultDelinquent.IsDelinquent.Bool {
			if userParam.AmountPaid.Float64 != resultDelinquent.AmountDue.Float64 {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "Payment amount must match the amount due"})
				return
			}

			var listID []int64
			listID, err = models.GetDelinquentID(db, userParam)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed get loans schedule delinquent"})
				return
			}

			for i := 0; i < len(listID); i++ {
				err = models.UpdateLoanScheduleStatus(tx, models.PaymentModel{
					ID: sql.NullInt64{Int64: listID[i]},
				})
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed update loan schedule status "})
					return
				}
			}

		} else {
			today := time.Now().Truncate(24 * time.Hour)
			if !today.Equal(loanOnDB.DueDate.Time) {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "You do not have bill to pay"})
				return
			}

			err = models.UpdateLoanScheduleStatus(tx, models.PaymentModel{
				ID: sql.NullInt64{Int64: loanOnDB.LoanScheduleID.Int64},
			})
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed update loan schedule status "})
				return
			}
		}

		_, err = models.InsertPayment(tx, userParam)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed insert into payment "})
			return
		}

		if loanOnDB.Outstanding.Float64 != 0 {
			outstanding = loanOnDB.Outstanding.Float64 - userParam.AmountPaid.Float64
		}
		err = models.UpdateLoanStatus(tx, models.LoanModel{
			Outstanding:  sql.NullFloat64{Float64: outstanding},
			IsDelinquent: sql.NullBool{Bool: false},
			ID:           sql.NullInt64{Int64: userParam.LoanID.Int64},
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed update loan outstanding"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Payment successfully"})
	}
}
