package routes

import (
	"BillingEngine/controllers"
	"database/sql"
	"github.com/gin-gonic/gin"
)

func RegisterLoanRoutes(router *gin.Engine, db *sql.DB) {
	router.GET("/loan/:cust_id", controllers.GetOutStanding(db))
	router.POST("/loan/create", controllers.CreateLoan(db))
	router.POST("/loan/payment", controllers.PaymentLoan(db))

}
