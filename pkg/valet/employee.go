package valet

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// Employee information
type Employee struct {
	gorm.Model
	User         User                   `json:"user" gorm:"foreignkey:user_id"`
	UserID       uint                   `json:"-"`
	Availability []EmployeeAvailability `json:"availability"`
}

// EmployeeAvailability information
type EmployeeAvailability struct {
	gorm.Model
	Employee   Employee `json:"employee" gorm:"foreignkey:employee_id"`
	EmployeeID uint     `json:"-"`
	Shift      Shift    `json:"shift" gorm:"foreignkey:shift_id"`
	ShiftID    uint     `json:"-"`
	Active     bool     `json:"active"`
}

// TableName will set the custom name for employee_availability
func (EmployeeAvailability) TableName() string {
	return "employee_availability"
}

// GetAvailability returns all active availability shifts
func (e Employee) GetAvailability() []EmployeeAvailability {
	avail := []EmployeeAvailability{}
	for _, ea := range e.Availability {
		if !ea.Active {
			continue
		}
		avail = append(avail, ea)
	}
	return avail
}
