package models

import (
	"github.com/jinzhu/gorm"
)

// User defines a user domain Model
type User struct {
	gorm.Model
	Username string
	Password string
	Email    string
}
