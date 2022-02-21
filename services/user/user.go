package user

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	errNoName     = errors.New("name must be provided")
	errNoEmail    = errors.New("email must be provided")
	errNoPassword = errors.New("password must be provided")

	errTokenExpired = errors.New("expired token")
	errInvalidToken = errors.New("invalid token or claims not found")

	signingSecret interface{} = []byte("hello world")
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

// Migrate migrates a new user repository instance.
func Migrate(conn *gorm.DB) error {
	return conn.AutoMigrate(User{})
}

// Create creates a new user
func (d *DB) Create(user User) (*User, error) {
	if err := user.validate(); err != nil {
		return nil, err
	}
	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(password)

	result := d.Conn.Create(&user)

	user.Password = "" //Clear the password before sending it back to user
	return &user, result.Error
}

// Authenticate authenticates the user using the email and password stored in the database
func (d *DB) Authenticate(email, password string) (*User, error) {
	var user User
	if err := d.Conn.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, err
	}
	user.Password = ""
	token, err := generateToken(uint32(user.ID))
	if err != nil {
		return nil, err
	}
	user.Token = token
	return &user, nil
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
