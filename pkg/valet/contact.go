package valet

import (
	"github.com/jinzhu/gorm"
	// mysql driver
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// Contact information
type Contact struct {
	gorm.Model
	Name        string        `json:"name"`
	Address     Address       `json:"address" gorm:"foreignkey:address_id"`
	AddressID   uint          `json:"-"`
	PhoneNumber []PhoneNumber `json:"phone_numbers" gorm:"many2many:contact_phone_numbers;"`
}

// PhoneNumber information
type PhoneNumber struct {
	ID    int    `json:"id"`
	Type  string `json:"type"` // mobile, home, work, etc
	Value string `json:"value"`
}

// String prints out the value
func (pn PhoneNumber) String() string {
	return pn.Value
}
