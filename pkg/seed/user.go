package seed

import (
	"log"
	"raedar/pkg/repository/engines"
	models "raedar/pkg/repository/models"
)

var multipleUsers = []models.User{
	models.User{
		Username: "Test_user1",
		Email:    "testuser1@mail.com",
		Password: "testPass12",
	},
	models.User{
		Username: "Test_user2",
		Email:    "testuser2@mail.com",
		Password: "testPass123",
	},
}

var oneUser = &models.User{
	Username: "Test_user",
	Email:    "testuser@mail.com",
	Password: "testPass12",
}

func OneUser() (models.User, error) {
	err := engines.PostgresDB().Create(oneUser).Error
	if err != nil {
		log.Fatalf("cannot seed users table: %v", err)
		return *oneUser, err
	}

	return *oneUser, nil
}

func MultipleUsers() error {
	for i, _ := range multipleUsers {
		err := engines.PostgresDB().Create(&multipleUsers[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
			return err
		}
	}

	return nil
}
