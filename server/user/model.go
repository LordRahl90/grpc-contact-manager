package user

import (
	"errors"

	"gorm.io/gorm"
)

var (
	errNoName     = errors.New("name must be provided")
	errNoEmail    = errors.New("email must be provided")
	errNoPassword = errors.New("password must be provided")
)

// User model
type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

type DB struct {
	Conn *gorm.DB
}

// Create creates a new user
func (d *DB) Create(user User) (*User, error) {
	if err := user.validate(); err != nil {
		return nil, err
	}

	result := d.Conn.Create(&user)

	return &user, result.Error
}

// Authenticate authenticates the user using the email and password stored in the database
func (d *DB) Authenticate(email, password string) (*User, error) {
	return nil, nil
}

// GenerateToken generates the JWT token for the given user
func (u *User) GenerateToken() (string, error) {
	return "", nil
}

func ValidateToken(token string) bool {
	return false
}

func (u User) validate() error {
	if u.Name == "" {
		return errNoName
	}
	if u.Email == "" {
		return errNoEmail
	}

	if u.Password == "" {
		return errNoPassword
	}

	return nil
}
