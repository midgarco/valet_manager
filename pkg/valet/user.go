package valet

import (
	"errors"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"golang.org/x/crypto/bcrypt"
)

// User information
type User struct {
	gorm.Model
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	Token     string `json:"-"`
}

// Hook Functions

// func (u *User) AfterFind() error {
// 	return nil
// }

// func (u *User) BeforeDelete() error {
// 	return nil
// }

// func (u *User) AfterDelete() error {
// 	return nil
// }

// func (u *User) BeforeSave() error {
// 	return nil
// }

// func (u *User) AfterSave() error {
// 	return nil
// }

// func (u *User) AfterCreate() error {
// 	return nil
// }

// BeforeCreate hooks in logic before the user is created in the db
func (u *User) BeforeCreate() error {
	// hash the password
	err := u.SetPassword(u.Password)
	if err != nil {
		return err
	}

	return nil
}

// Public Methods

// SetPassword encrypts the password for the user
func (u *User) SetPassword(pass string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.MinCost)
	if err != nil {
		return errors.New("generating password hash")
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and set it
	u.Password = string(hash)

	return nil
}

// Create a new user
func (u *User) Create(db *gorm.DB) error {
	return db.Create(u).Error
}

// Save the user data
func (u *User) Save(db *gorm.DB) error {
	return db.Save(u).Error
}

// FindUser queries the database for the user data
func FindUser(db *gorm.DB, id int) (*User, error) {
	u := User{}
	if err := db.First(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

// FindUsers queries for all users
func FindUsers(db *gorm.DB) ([]User, error) {
	users := []User{}
	if err := db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
