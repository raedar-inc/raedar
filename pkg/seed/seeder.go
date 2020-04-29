package seed

import (
	"log"
	"raedar/pkg/repository/engines"
	models "raedar/pkg/repository/models"
)

// Clean drops and recreates all tables in the Database.
func Clean() {
	err := engines.PostgresDB().DropTableIfExists(&models.User{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}
	err = engines.PostgresDB().AutoMigrate(&models.User{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}
}
