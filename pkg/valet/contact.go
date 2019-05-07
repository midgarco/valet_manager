package valet

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// Contact information
type Contact struct {
	gorm.Model
	Name      string  `json:"name"`
	Address   Address `json:"address" gorm:"foreignkey:address_id"`
	AddressID uint    `json:"-"`
	Phone     []Phone `json:"phone" gorm:"many2many:contact_phones;"`
}

// Phone information
type Phone struct {
	ID     int    `json:"id"`
	Type   string `json:"type"` // mobile, home, work, etc
	Number string `json:"number"`
}

// String prints out the value
func (p Phone) String() string {
	return p.Number
}
