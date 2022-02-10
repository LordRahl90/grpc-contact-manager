package user

import (
	"errors"
	"fmt"

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

	fmt.Printf("Hello! Got here\n")

	result := d.Conn.Create(&user)

	return &user, result.Error
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
