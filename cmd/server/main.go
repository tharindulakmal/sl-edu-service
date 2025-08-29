package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"

	// db "github.com/tharindulakmal/sl-edu-service/internal/database"
	"github.com/tharindulakmal/sl-edu-service/internal/handlers"
	"github.com/tharindulakmal/sl-edu-service/internal/repository"
)

func main() {

	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env")
	}

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		dbUser, dbPass, dbHost, dbPort, dbName,
	)

	db, err := sql.Open("mysql", dsn)

	if err != nil {
		log.Fatalf("could not connect to DB: %v", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("could not ping DB: %v", err)
	}
	gradeRepo := repository.NewGradeRepository(db)
	gradeHandler := handlers.NewGradeHandler(gradeRepo)

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	r.GET("/api/v1/grades", gradeHandler.GetGrades)

	log.Println("Starting server on :8080")
	r.Run(":8080")
}
