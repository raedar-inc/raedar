package services

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"

	"raedar/pkg/repository/engines"
	"raedar/pkg/repository/models"
)

// Token for user - JWT claims struct
type Token struct {
	UserID uint
	jwt.StandardClaims
}

// User struct defines a user domain Model
type User struct{}

//UserServicer service interface
type userServicer interface {
	FindById(id int) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	FindByChangePasswordHash(hash string) (*models.User, error)
	FindByValidationHash(hash string) (*models.User, error)
	FindAllUsers() ([]*models.User, error)
	Update(user *models.User) error
	Delete(id int) error
	Save(user *models.User) (*models.User, error)
}

func hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func verifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// BeforeSave This method will be triggered before creating a new user
func (u *User) BeforeSave(user *models.User) error {
	hashedPassword, err := hash(user.Password)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.IsClient = false
	user.IsVerified = false
	user.IsCustomer = true
	return nil
}

func (u *User) prepare(user *models.User) {
	user.Username = html.EscapeString(strings.TrimSpace(user.Username))
	user.Email = html.EscapeString(strings.TrimSpace(user.Email))
}

// validate: validates incoming user (object) map before continuing to store it to the DB.
func (u *User) validate(action string, user *models.User) string {
	switch strings.ToLower(action) {
	case "update":
		if user.Password == "" {
			return "Required Password"
		}
		if user.Email == "" {
			return "Required Email"
		}
		if err := checkmail.ValidateFormat(user.Email); err != nil {
			return "Invalid Email"
		}

		return ""
	case "login":
		if user.Password == "" {
			return "Required Password"
		}
		if user.Email == "" {
			return "Required Email"
		}
		if err := checkmail.ValidateFormat(user.Email); err != nil {
			return "Invalid Email"
		}
		return ""

	default:
		if user.Password == "" {
			return "Password is required"
		}
		if user.Username == "" {
			return "Username is required"
		}
		if user.Email == "" {
			return "Required Email"
		}
		if err := checkmail.ValidateFormat(user.Email); err != nil {
			return "Invalid Email"
		}
	}

	//Email must be unique
	temp := &models.User{}

	//check for errors and duplicate emails
	err := engines.PostgresDB().Table("users").Where("email = ?", user.Email).First(temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return "Connection error. Please retry"
	}
	if temp.Email != "" {
		return "email address is taken"
	} else if temp.Username != "" {
		return "username is taken"
	}

	return ""
}

// Save method saves a new user to the database.
func (u *User) Save(user *models.User) (*User, string) {
	err := u.validate("", user)
	if err != "" {
		return nil, err
	}
	engines.PostgresDB().Create(&user)
	return u, ""
}

// GetUser Returns a single user
func (u *User) GetUser(id uint) *models.User {
	user := &models.User{}
	engines.PostgresDB().Table("users").Where("id = ?", id).First(user)
	if user.Email == "" { //User not found!
		return nil
	}

	user.Password = ""
	return user
}

// FindAllUsers returns all users stored in the database
func (u *User) FindAllUsers(db *gorm.DB) (*[]models.User, error) {
	var err error
	users := []models.User{}
	err = db.Debug().Model(&User{}).Find(&users).Error
	if err != nil {
		return &[]models.User{}, err
	}
	return &users, err
}

// FindUserByID returns a single user record given User Id
func (u *User) FindUserByID(uid uint32) (*models.User, error) {
	user := &models.User{}
	engines.PostgresDB().Find(&user)
	if user.Email == "" {
		return nil, errors.New("No user found")
	}
	return user, nil
}
