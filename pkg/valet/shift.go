package valet

import (
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// Shift information
type Shift struct {
	gorm.Model
	Start time.Time `json:"start_time"`
	End   time.Time `json:"end_time"`
}
