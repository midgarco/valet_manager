package valet

import (
	"log"

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
	Token     string `json:"token"`
}

// SetPassword encrypts the password for the user
func (u *User) SetPassword(pass string) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and set it
	u.Password = string(hash)
}

// Create a new user
func (u *User) Create(db *gorm.DB) error {
	return db.Create(u).Error
}

// Save the user data
func (u User) Save(db *gorm.DB) error {
	return db.Save(&u).Error
}
