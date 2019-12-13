package models

import (
	"github.com/jinzhu/gorm"
)

// User struct defines a user domain Model
type User struct {
	gorm.Model
	// ID          int        `gorm:"primary_key" json:"id"`
	Username    string `gorm:"type:varchar(255); unique_index; not null" json:"username" binding:"required"`
	Password    string `gorm:"type:varchar(255); unique_index; not null" json:"password" binding:"required"`
	Email       string `gorm:"type:varchar(255); unique_index; not null" json:"email" binding:"required"`
	AccessToken string `gorm:"type:varchar(255)" json:"access_token"`
	IsVerified  bool   `gorm:"default:false" json:"isVerified"`
	IsClient    bool   `gorm:"default:false" json:"isClient"`
	IsCustomer  bool   `gorm:"default:false" json:"isCustomer"`
}
