package contact

import (
	"errors"

	"grpc-contact-manager/services/user"

	"gorm.io/gorm"
)

var (
	errInvalidUserID  = errors.New("invalid user id")
	errEmptyName      = errors.New("full name must be provided")
	errEmptyPhone     = errors.New("phone number must be provided")
	errEmptyEmail     = errors.New("email must be provided")
	errEmptyAddress   = errors.New("address must be provided")
	errContactExists  = errors.New("contact with this email exists")
	errNotUserContact = errors.New("user has no access to contact")
)

type Contact struct {
	gorm.Model
	UserID   uint `json:"user_id" gorm:"column:user_id;index:idx_user_id"`
	User     user.User
	Fullname string `json:"full_name" gorm:"column:full_name"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
	Email    string `json:"email" gorm:"column:email;index:idx_email"`
	// Email    string `json:"email" gorm:"column:email index:unique"`
	// TODO: Revisit the multiple column index, so a user doesn't add a contact with more than 1 same email
	// UserIDEmail string `json:"-" gorm:"uniqueIndex:idx_user_id_email"`
}

// DB - db connection abstraction
type DB struct {
	Conn *gorm.DB
}

// Create adds a new contact record for the given user.
func (db *DB) Create(contact Contact) (*Contact, error) {
	if err := contact.validate(); err != nil {
		return nil, err
	}
	// check for possible duplicate
	var c Contact
	err := db.Conn.Where("user_id = ? AND email = ?", contact.UserID, contact.Email).First(&c).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if c.ID != 0 {
		return nil, errContactExists
	}
	result := db.Conn.Create(&contact)
	return &contact, result.Error
}

// FindByUserID returns all the contacts for a given user ID
func (db *DB) FindByUserID(userID uint32) ([]Contact, error) {
	var contacts []Contact
	res := db.Conn.Where("user_id = ?", userID).Find(&contacts)
	return contacts, res.Error
}

func (db *DB) FindByID(userID, id uint) (*Contact, error) {
	var contact Contact
	err := db.Conn.First(&contact, id).Error
	if contact.UserID != userID {
		return nil, errNotUserContact
	}
	return &contact, err
}

// Search search the full name and email for the given string
func (db *DB) Search(userID uint32, search string) ([]Contact, error) {
	var contacts []Contact
	res := db.Conn.Where("user_id = ? AND (full_name LIKE ? OR email LIKE ? )", userID, "%"+search+"%", "%"+search+"%").Find(&contacts)
	return contacts, res.Error
}

// Update the value of a contact
func (db *DB) Update(contact *Contact) error {
	return db.Conn.Save(&contact).Error
}

func (c *Contact) validate() error {
	if c.UserID == 0 {
		return errInvalidUserID
	}
	if c.Fullname == "" {
		return errEmptyName
	}
	if c.Email == "" {
		return errEmptyEmail
	}
	if c.Phone == "" {
		return errEmptyPhone
	}
	if c.Address == "" {
		return errEmptyAddress
	}
	return nil
}
