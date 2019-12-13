package engines

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // postges dialect for database connection
	"github.com/joho/godotenv"
	"log"
	"os"
	models "raedar/pkg/repository/models"
)

var db *gorm.DB

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Print("Failed to load environment variables", err)
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	dbHost := os.Getenv("DB_HOST")

	// Build connection string
	dbURI := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, dbPort, dbUser, dbName, dbPassword)

	// connect to the database
	conn, err := gorm.Open("postgres", dbURI)
	if err != nil {
		fmt.Printf("Cannot connect to %s database", "postgres")
		log.Fatal("This is the error:", err)
	} else {
		fmt.Printf("We are connected to the %s database", "postgres")
	}

	db = conn
	// Migrate the schema
	db.Debug().AutoMigrate(&models.User{}) // Handles Database migration
}

// PostgresDB returns a handler to the DB object
func PostgresDB() *gorm.DB {
	return db
}
