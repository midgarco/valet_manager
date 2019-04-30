package valet

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// Client information
type Client struct {
	gorm.Model
	Name string `json:"name"`
}
