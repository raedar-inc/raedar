package models

import (
	"github.com/jinzhu/gorm"
	"github.com/segmentio/ksuid"
)

// User struct defines a user domain Model
type User struct {
	gorm.Model
	UUID         ksuid.KSUID `json:"uuid"`
	Username     string      `gorm:"type:varchar(255); unique_index; not null" json:"username" binding:"required"`
	Password     string      `gorm:"type:varchar(255); unique_index; not null" json:"password" binding:"required"`
	Email        string      `gorm:"type:varchar(255); unique_index; not null" json:"email" binding:"required"`
	RefreshToken string      `gorm:"type:varchar(255)" json:"refresh_token"`
	IsVerified   bool        `gorm:"default:false" json:"isVerified"`
	IsClient     bool        `gorm:"default:false" json:"isClient"`
	IsAdmin      bool        `gorm:"default:false" json:"isAdmin"`
	IsCustomer   bool        `gorm:"default:false" json:"isCustomer"`
}
