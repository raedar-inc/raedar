package engines

import (
	"fmt"
	"log"
	"os"
	models "raedar/pkg/repository/models"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // postges dialect for database connection
	"github.com/joho/godotenv"
)

var db *gorm.DB

func init() {
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Print("Failed to load environment variables and trying again:\n", err)
		err = godotenv.Load("../../../.env")
	}
	if err != nil {
		log.Print("Failed to load environment variables and trying the docker.env file :\n", err)
		err = godotenv.Load("../../../docker.env")
	}
	if err != nil {
		log.Print("Failed to load environment variables completely ", err)
	}

	var dbUser = os.Getenv("DB_USER")
	var dbName = os.Getenv("DB_NAME")
	var dbPort = os.Getenv("DB_PORT")
	var dbHost = os.Getenv("DB_HOST")
	var dbPassword = os.Getenv("DB_PASSWORD")

	if mode := os.Getenv("MODE"); mode == "TESTING_MODE" {
		dbPassword = os.Getenv("TEST_DB_PASSWORD")
		dbName = os.Getenv("TEST_DB_NAME")
	}

	// Build connection string
	dbURI := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, dbPort, dbUser, dbName, dbPassword)

	// connect to the database
	conn, err := gorm.Open("postgres", dbURI)
	if err != nil {
		fmt.Printf("Cannot connect to %s database\n ", dbName)
		log.Fatal("This is the error:", err)
	} else {
		fmt.Printf("We are connected to the %v database\n", dbName)
	}

	db = conn
	// Migrate the schema
	db.Debug().AutoMigrate(&models.User{}) // Handles Database migration
}

// PostgresDB returns a handler to the DB object
func PostgresDB() *gorm.DB {
	return db
}
