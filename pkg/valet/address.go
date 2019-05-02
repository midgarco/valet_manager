package valet

import "github.com/jinzhu/gorm"

// Address details
type Address struct {
	gorm.Model
	Line1   string  `json:"line1"`
	Line2   string  `json:"line2"`
	City    string  `json:"city"`
	State   string  `json:"state"`
	Zipcode string  `json:"zipcode"`
	Lat     float64 `json:"latitude"`
	Lng     float64 `json:"longitude"`
}
