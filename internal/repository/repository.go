package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/CB-AutoStack/AutoStack/apps/api-inventory/internal/models"
	"github.com/sirupsen/logrus"
)

// Repository provides data access for users and vehicles
type Repository struct {
	users    map[string]*models.User
	vehicles map[string]*models.Vehicle
	mu       sync.RWMutex
	logger   *logrus.Logger
}

// NewRepository creates a new repository and loads data from JSON files
func NewRepository(dataPath string, logger *logrus.Logger) (*Repository, error) {
	repo := &Repository{
		users:    make(map[string]*models.User),
		vehicles: make(map[string]*models.Vehicle),
		logger:   logger,
	}

	// Load users
	if err := repo.loadUsers(filepath.Join(dataPath, "users.json")); err != nil {
		return nil, fmt.Errorf("failed to load users: %w", err)
	}

	// Load vehicles
	if err := repo.loadVehicles(filepath.Join(dataPath, "vehicles.json")); err != nil {
		return nil, fmt.Errorf("failed to load vehicles: %w", err)
	}

	logger.Infof("Loaded %d users and %d vehicles from %s", len(repo.users), len(repo.vehicles), dataPath)

	return repo, nil
}

// loadUsers loads users from a JSON file
func (r *Repository) loadUsers(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var users []*models.User
	if err := json.Unmarshal(data, &users); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	for _, user := range users {
		r.users[user.ID] = user
	}

	return nil
}

// loadVehicles loads vehicles from a JSON file
func (r *Repository) loadVehicles(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var vehicles []*models.Vehicle
	if err := json.Unmarshal(data, &vehicles); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	for _, vehicle := range vehicles {
		r.vehicles[vehicle.ID] = vehicle
	}

	return nil
}

// GetUserByID retrieves a user by ID
func (r *Repository) GetUserByID(userID string) (*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[userID]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}

// GetUserByEmail retrieves a user by email address
func (r *Repository) GetUserByEmail(email string) (*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}

	return nil, fmt.Errorf("user not found")
}

// GetAllUsers returns all users
func (r *Repository) GetAllUsers() []*models.User {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]*models.User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}

	return users
}

// GetVehicleByID retrieves a vehicle by ID
func (r *Repository) GetVehicleByID(vehicleID string) (*models.Vehicle, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	vehicle, exists := r.vehicles[vehicleID]
	if !exists {
		return nil, fmt.Errorf("vehicle not found")
	}

	return vehicle, nil
}

// GetAllVehicles returns all vehicles
func (r *Repository) GetAllVehicles() []*models.Vehicle {
	r.mu.RLock()
	defer r.mu.RUnlock()

	vehicles := make([]*models.Vehicle, 0, len(r.vehicles))
	for _, vehicle := range r.vehicles {
		vehicles = append(vehicles, vehicle)
	}

	return vehicles
}

// SearchVehicles searches for vehicles based on filter criteria
func (r *Repository) SearchVehicles(filter *models.VehicleFilter) []*models.Vehicle {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var results []*models.Vehicle

	for _, vehicle := range r.vehicles {
		if r.matchesFilter(vehicle, filter) {
			results = append(results, vehicle)
		}
	}

	return results
}

// matchesFilter checks if a vehicle matches the given filter
func (r *Repository) matchesFilter(vehicle *models.Vehicle, filter *models.VehicleFilter) bool {
	if filter == nil {
		return true
	}

	// Make filter
	if filter.Make != "" && !strings.EqualFold(vehicle.Make, filter.Make) {
		return false
	}

	// Model filter
	if filter.Model != "" && !strings.EqualFold(vehicle.Model, filter.Model) {
		return false
	}

	// Type filter
	if filter.Type != "" && !strings.EqualFold(vehicle.Type, filter.Type) {
		return false
	}

	// Condition filter
	if filter.Condition != "" && !strings.EqualFold(vehicle.Condition, filter.Condition) {
		return false
	}

	// Price range filter
	if filter.MinPrice > 0 && vehicle.Price < filter.MinPrice {
		return false
	}
	if filter.MaxPrice > 0 && vehicle.Price > filter.MaxPrice {
		return false
	}

	// Currency filter
	if filter.Currency != "" && !strings.EqualFold(vehicle.Currency, filter.Currency) {
		return false
	}

	// Country filter
	if filter.Country != "" && !strings.EqualFold(vehicle.Country, filter.Country) {
		return false
	}

	// Year range filter
	if filter.MinYear > 0 && vehicle.Year < filter.MinYear {
		return false
	}
	if filter.MaxYear > 0 && vehicle.Year > filter.MaxYear {
		return false
	}

	// Fuel type filter
	if filter.FuelType != "" && !strings.EqualFold(vehicle.FuelType, filter.FuelType) {
		return false
	}

	// Transmission filter
	if filter.Transmission != "" && !strings.EqualFold(vehicle.Transmission, filter.Transmission) {
		return false
	}

	// Drivetrain filter
	if filter.Drivetrain != "" && !strings.EqualFold(vehicle.Drivetrain, filter.Drivetrain) {
		return false
	}

	// Vehicle types filter (array)
	if len(filter.VehicleTypes) > 0 {
		found := false
		for _, vtype := range filter.VehicleTypes {
			if strings.EqualFold(vehicle.Type, vtype) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	return true
}
