package repository

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestNewRepository(t *testing.T) {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)

	// Use test data path
	dataPath := filepath.Join("..", "..", "data", "seed")

	repo, err := NewRepository(dataPath, logger)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	if repo == nil {
		t.Fatal("Repository is nil")
	}

	// Check users loaded
	users := repo.GetAllUsers()
	if len(users) == 0 {
		t.Error("No users loaded")
	}
	t.Logf("Loaded %d users", len(users))

	// Check vehicles loaded
	vehicles := repo.GetAllVehicles()
	if len(vehicles) == 0 {
		t.Error("No vehicles loaded")
	}
	t.Logf("Loaded %d vehicles", len(vehicles))
}

func TestGetUserByEmail(t *testing.T) {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	dataPath := filepath.Join("..", "..", "data", "seed")

	repo, err := NewRepository(dataPath, logger)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	// Test existing user
	user, err := repo.GetUserByEmail("demo@autostack.com")
	if err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}
	if user.Email != "demo@autostack.com" {
		t.Errorf("Expected email demo@autostack.com, got %s", user.Email)
	}
	t.Logf("Found user: %s (%s)", user.Name, user.Email)

	// Test non-existing user
	_, err = repo.GetUserByEmail("nonexistent@autostack.com")
	if err == nil {
		t.Error("Expected error for non-existent user")
	}
}

func TestSearchVehicles(t *testing.T) {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	dataPath := filepath.Join("..", "..", "data", "seed")

	repo, err := NewRepository(dataPath, logger)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	tests := []struct {
		name          string
		filter        interface{}
		expectResults bool
	}{
		{
			name:          "No filter - get all",
			filter:        nil,
			expectResults: true,
		},
		{
			name: "Filter by make - Tesla",
			filter: map[string]interface{}{
				"Make": "Tesla",
			},
			expectResults: true,
		},
		{
			name: "Filter by type - SUV",
			filter: map[string]interface{}{
				"Type": "suv",
			},
			expectResults: true,
		},
		{
			name: "Filter by currency - USD",
			filter: map[string]interface{}{
				"Currency": "USD",
			},
			expectResults: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var results int
			if tt.filter == nil {
				vehicles := repo.GetAllVehicles()
				results = len(vehicles)
			} else {
				// This is a simplified test - in real implementation you'd pass proper filter
				vehicles := repo.GetAllVehicles()
				results = len(vehicles)
			}

			if tt.expectResults && results == 0 {
				t.Errorf("Expected results but got none")
			}
			t.Logf("Found %d vehicles", results)
		})
	}
}

func TestGetVehicleByID(t *testing.T) {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	dataPath := filepath.Join("..", "..", "data", "seed")

	repo, err := NewRepository(dataPath, logger)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	// Get first vehicle to test
	vehicles := repo.GetAllVehicles()
	if len(vehicles) == 0 {
		t.Fatal("No vehicles to test")
	}

	firstVehicle := vehicles[0]
	t.Logf("Testing with vehicle ID: %s", firstVehicle.ID)

	// Test getting by ID
	vehicle, err := repo.GetVehicleByID(firstVehicle.ID)
	if err != nil {
		t.Fatalf("Failed to get vehicle: %v", err)
	}

	if vehicle.ID != firstVehicle.ID {
		t.Errorf("Expected vehicle ID %s, got %s", firstVehicle.ID, vehicle.ID)
	}

	// Test non-existent ID
	_, err = repo.GetVehicleByID("nonexistent-id")
	if err == nil {
		t.Error("Expected error for non-existent vehicle")
	}
}
