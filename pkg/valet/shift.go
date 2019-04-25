package valet

import "time"

// Shift information
type Shift struct {
	Start time.Time `json:"start_time"`
	End   time.Time `json:"end_time"`
}
