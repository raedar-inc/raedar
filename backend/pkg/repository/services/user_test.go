package services

import (
	"raedar/pkg/repository/models"
	"raedar/pkg/seed"
	"testing"

	"gopkg.in/go-playground/assert.v1"
)

// var userService User{}

// Test create a user
func TestCreateUser(t *testing.T) {
	var user = models.User{
		Username: "Test_user",
		Email:    "testuser@mail.com",
		Password: "testPass12",
	}
	userService := User{}
	createdUser, serviceErr := userService.Save(&user)
	if serviceErr != "" {
		t.Errorf("this is the returned error: %v\n", serviceErr)
		return
	}

	assert.Equal(t, serviceErr, "")
	assert.Equal(t, createdUser.Username, user.Username)
	defer seed.Clean() // drops tables after the function returns
}

// Test find user by email
func TestFindUserByEmail(t *testing.T) {
	defer seed.Clean() // drops tables after the function returns
	user, err := seed.OneUser()
	if err != nil {
		t.Errorf("this is the returned error: %v\n", err)
		return
	}

	userService := User{}
	foundUser, err := userService.FindByEmail(user.Email)

	assert.Equal(t, user.Username, foundUser.Username)
}

// Test find user by UUID
func TestFindByUsername(t *testing.T) {
	defer seed.Clean() // drops tables after the function returns
	user, err := seed.OneUser()
	if err != nil {
		t.Errorf("this is the returned error: %v\n", err)
		return
	}

	userService := User{}
	foundUser, err := userService.FindByUsername(user.Username)

	assert.Equal(t, user.Email, foundUser.Email)
}
