package main

import (
	"BillingEngine/config"
	"BillingEngine/middleware"
	"BillingEngine/routes"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/packr/v2"
	"github.com/joho/godotenv"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/tpkeeper/gin-dump"
	"io"
	"log"
	"os"
)

// TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func SetupLogOutput() {
	f, _ := os.Create("BillingEngine.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	db := config.ConnectDB()
	defer db.Close()
	dbMigrate(db)

	SetupLogOutput()
	app := gin.Default()
	app.Use(
		gin.Recovery(),
		middleware.Logger(),
		middleware.BasicAuth(
			os.Getenv("USERNAME"),
			os.Getenv("PASSWORD")),
		gindump.Dump())

	routes.RegisterLoanRoutes(app, db)
	if err = app.Run(os.Getenv("PORT")); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func dbMigrate(db *sql.DB) {
	migrations := &migrate.PackrMigrationSource{
		Box: packr.New("migrations", "./sql_migrations"),
	}
	if db != nil {
		_, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
		if err != nil {
			log.Fatalf("Failed to migrate DB: %v", err)
		} else {
			fmt.Println("success migrate db ")
		}
	} else {
		os.Exit(3)
	}
}
