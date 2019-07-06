package valet

import (
	"errors"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/midgarco/valet_manager/pkg/pagination"
	"golang.org/x/crypto/bcrypt"
)

// User information
type User struct {
	gorm.Model
	FirstName    string  `json:"first_name"`
	LastName     string  `json:"last_name"`
	Email        string  `json:"email" gorm:"unique_index"`
	Password     string  `json:"-"`
	Token        string  `json:"-"`
	Address      Address `json:"address" gorm:"foreignkey:address_id"`
	AddressID    uint    `json:"-"`
	PhoneNumbers []Phone `json:"phone_numbers" gorm:"many2many:user_phones;"`
	Admin        bool    `json:"-"`
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
	// check for an existing user
	ex := User{}
	if !db.Where("email = ?", u.Email).First(&ex).RecordNotFound() {
		return errors.New("email already exists")
	}
	return db.Create(u).Error
}

// Save the user data
func (u *User) Save(db *gorm.DB) error {
	return db.Save(u).Error
}

// FindUser queries the database for the user data
func FindUser(db *gorm.DB, id int) (*User, error) {
	u := User{}
	if err := db.Preload("Address").First(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

// FindUsers queries for all users
func FindUsers(db *gorm.DB, paging pagination.Paging) ([]User, error) {
	users := []User{}
	dbh := db.Preload("Address").Offset(paging.Offset).Limit(paging.Limit)
	for _, ord := range paging.OrderBy {
		dbh = dbh.Order(ord.Field + " " + ord.Direction)
	}
	if err := dbh.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// GetUserByEmail finds the user record with the provided email address
func GetUserByEmail(db *gorm.DB, email string) (*User, error) {
	u := User{}
	if err := db.Preload("Address").Where("email = ?", email).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

// RemoveTestUser will delete ALL data related to a user
// USED FOR TESTING!
func RemoveTestUser(db *gorm.DB, id int) error {
	u := User{}
	if err := db.First(&u, id).Error; err != nil {
		return err
	}
	// delete the phone numbers
	for _, pn := range u.PhoneNumbers {
		if err := db.Unscoped().Delete(&pn).Error; err != nil {
			return err
		}
	}
	// delete the address
	if err := db.Unscoped().Delete(&u.Address).Error; err != nil {
		return err
	}
	// finally just delete the user record
	if err := db.Unscoped().Delete(&u).Error; err != nil {
		return err
	}
	return nil
}

// UserCount returns the number of user records in the db
func UserCount(db *gorm.DB) (int, error) {
	count := 0
	if err := db.Model(&User{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
