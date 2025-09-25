package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/tharindulakmal/sl-edu-service/internal/routes"
	// db "github.com/tharindulakmal/sl-edu-service/internal/database"
)

func main() {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
			"http://sl-edu-service-env.eba-f8bzvpsg.us-east-1.elasticbeanstalk.com",
			"https://sl-edu-service-env.eba-f8bzvpsg.us-east-1.elasticbeanstalk.com",
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false, // set true only if you actually use cookies/auth headers
		MaxAge:           12 * time.Hour,
	}))
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env")
	}

	if os.Getenv("SKIP_DB") == "true" {
		log.Println("Skipping database connection — running in placeholder mode")
	} else {
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
		// Register all routes in one place
		routes.RegisterRoutes(r, db)
	}

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	log.Println("Starting server on :8080")
	r.Run(":8080")
}
