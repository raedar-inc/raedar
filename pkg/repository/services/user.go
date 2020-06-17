package services

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/segmentio/ksuid"
	"golang.org/x/crypto/bcrypt"

	"raedar/pkg/repository/engines"
	"raedar/pkg/repository/models"
	"raedar/tools"
)

func genKsuuid() ksuid.KSUID {
	uuid := ksuid.New()
	return uuid
}

// User struct defines a user domain Model
type User struct{}

//UserServicer service interface
type userServicer interface {
	FindByUID(id int) (*models.User, error)
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

func (u *User) VerifyPassword(hashedPassword, password string) error {
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
	user.UUID = genKsuuid()
	user.IsVerified = false
	user.IsCustomer = true
	return nil
}

func (u *User) prepare(user *models.User) {
	user.Username = html.EscapeString(strings.TrimSpace(user.Username))
	user.Email = html.EscapeString(strings.TrimSpace(user.Email))
}

// validate: validates incoming user (object) map before continuing to store it to the DB.
func (u *User) Validate(action string, user *models.User) string {
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
	case "reset-password":
		if err := checkmail.ValidateFormat(user.Email); err != nil {
			return "Invalid Email provided"
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
func (u *User) Save(user *models.User) (*models.User, string) {
	err := u.Validate("", user)
	if err != "" {
		return nil, err
	}

	errorBS := u.BeforeSave(user)
	if errorBS != nil {
		return nil, "failed to save the user"
	}

	engines.PostgresDB().Create(&user)
	return user, ""
}

// Update method updates a records columns.
func (u *User) Update(user *models.User) (*models.User, string) {
	if user.Password != "" {
		hashedPassword, err := hash(user.Password)
		if err != nil {
			return nil, "Something went wrong"
		}
		user.Password = string(hashedPassword)
	}
	engines.PostgresDB().Select(user.Email).Updates(&user)
	return user, ""
}

// ComparePasswordToConfirmPassword compares a password to the confirm password for equality
func (u *User) ComparePasswordToConfirmPassword(password, confirmPassword string) bool {
	if password == confirmPassword {
		return true
	}
	return false
}

// FindAllUsers returns all users stored in the database
func (u *User) FindAllUsers(db *gorm.DB) (*[]models.User, error) {
	var err error
	var users []models.User
	err = db.Debug().Model(&User{}).Find(&users).Error
	if err != nil {
		return &[]models.User{}, err
	}
	return &users, err
}

// FindUserByID returns a single user record given User Id
func (u *User) FindByUID(uid ksuid.KSUID) (*models.User, error) {
	user := &models.User{}
	engines.PostgresDB().Table("users").Where("UUID = ?", uid).First(&user)
	if user.Email == "" {
		return nil, errors.New("No user found")
	}
	return user, nil
}

// FindByEmail returns a single user record given User email address
func (u *User) FindByEmail(email string) (*models.User, error) {
	user := &models.User{}
	engines.PostgresDB().Table("users").Where("Email = ?", email).First(&user)
	if user.Email == "" {
		return nil, errors.New("No user found")
	}
	return user, nil
}

// FindByUsername returns a single user record given Username
func (u *User) FindByUsername(username string) (*models.User, error) {
	user := &models.User{}
	engines.PostgresDB().Table("users").Where("Username = ?", username).First(&user)
	if user.Username == "" {
		return nil, errors.New("no user found")
	}
	return user, nil
}

// AccessToken creates a user access jwt token
func (u *User) AccessToken(user *models.User) (string, error) {
	var err error
	ecdsaKey := tools.GetPrivEcdsaKey()

	// Set claims
	// This is the information which frontend can use
	// The backend can also decode the token and get admin etc.
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = user.ID
	claims["isAdmin"] = user.IsAdmin
	claims["isClient"] = user.IsClient
	claims["isCustomer"] = user.IsCustomer
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix() //Token expires after 1 hour
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	// create access_token
	accessToken, err := token.SignedString(ecdsaKey)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

// AccessToken creates a user access jwt token
func (u *User) RefreshToken(user *models.User) (string, error) {
	var err error
	ecdsaKey := tools.GetPrivEcdsaKey()

	// Set claims
	// This is the information which frontend can use
	// The backend can also decode the token and get admin etc.
	rfClaims := jwt.MapClaims{}
	rfClaims["sub"] = 1
	rfClaims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	// generate refresh token
	rfToken := jwt.NewWithClaims(jwt.SigningMethodES256, rfClaims)

	refreshToken, err := rfToken.SignedString(ecdsaKey)
	if err != nil {
		return "", err
	}

	engines.PostgresDB().First(&user)
	user.RefreshToken = refreshToken
	engines.PostgresDB().Save(&user)

	return refreshToken, nil
}
