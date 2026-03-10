package models

import "time"

// User represents a user in the system
type User struct {
	ID                string       `json:"id"`
	Email             string       `json:"email"`
	Password          string       `json:"password"`
	Name              string       `json:"name"`
	Country           string       `json:"country"`
	PreferredCurrency string       `json:"preferredCurrency"`
	Roles             []string     `json:"roles"`
	Preferences       *Preferences `json:"preferences,omitempty"`
	CreatedAt         time.Time    `json:"createdAt"`
}

// Preferences represents user preferences for vehicle search
type Preferences struct {
	PriceRange   []int    `json:"priceRange,omitempty"`
	VehicleTypes []string `json:"vehicleTypes,omitempty"`
}
