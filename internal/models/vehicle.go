package models

import "time"

// Vehicle represents a vehicle listing in the inventory
type Vehicle struct {
	ID            string    `json:"id"`
	VIN           string    `json:"vin"`
	Year          int       `json:"year"`
	Make          string    `json:"make"`
	Model         string    `json:"model"`
	Trim          string    `json:"trim"`
	Type          string    `json:"type"`
	Condition     string    `json:"condition"`
	Mileage       int       `json:"mileage"`
	Price         float64   `json:"price"`
	Currency      string    `json:"currency"`
	Country       string    `json:"country"`
	Status        string    `json:"status"`
	FuelType      string    `json:"fuelType"`
	Transmission  string    `json:"transmission"`
	Drivetrain    string    `json:"drivetrain"`
	ExteriorColor string    `json:"exteriorColor"`
	InteriorColor string    `json:"interiorColor"`
	Features      []string  `json:"features"`
	Images        []string  `json:"images"`
	DealerRating  float64   `json:"dealerRating"`
	Location      string    `json:"location"`
	ListingDate   time.Time `json:"listingDate"`
}

// VehicleFilter represents filter options for vehicle search
type VehicleFilter struct {
	Make         string   `json:"make,omitempty"`
	Model        string   `json:"model,omitempty"`
	Type         string   `json:"type,omitempty"`
	Condition    string   `json:"condition,omitempty"`
	MinPrice     float64  `json:"minPrice,omitempty"`
	MaxPrice     float64  `json:"maxPrice,omitempty"`
	Currency     string   `json:"currency,omitempty"`
	Country      string   `json:"country,omitempty"`
	MinYear      int      `json:"minYear,omitempty"`
	MaxYear      int      `json:"maxYear,omitempty"`
	FuelType     string   `json:"fuelType,omitempty"`
	Transmission string   `json:"transmission,omitempty"`
	Drivetrain   string   `json:"drivetrain,omitempty"`
	VehicleTypes []string `json:"vehicleTypes,omitempty"`
}
