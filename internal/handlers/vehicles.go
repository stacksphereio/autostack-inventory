package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/CB-AutoStack/AutoStack/apps/api-inventory/internal/models"
	"github.com/CB-AutoStack/AutoStack/apps/api-inventory/internal/repository"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// VehicleHandler handles vehicle-related requests
type VehicleHandler struct {
	repo   *repository.Repository
	logger *logrus.Logger
}

// NewVehicleHandler creates a new vehicle handler
func NewVehicleHandler(repo *repository.Repository, logger *logrus.Logger) *VehicleHandler {
	return &VehicleHandler{
		repo:   repo,
		logger: logger,
	}
}

// HandleListVehicles returns all vehicles or filtered results
func (h *VehicleHandler) HandleListVehicles(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	// Build filter from query parameters
	filter := &models.VehicleFilter{}

	if make := query.Get("make"); make != "" {
		filter.Make = make
	}
	if model := query.Get("model"); model != "" {
		filter.Model = model
	}
	if vtype := query.Get("type"); vtype != "" {
		filter.Type = vtype
	}
	if condition := query.Get("condition"); condition != "" {
		filter.Condition = condition
	}
	if currency := query.Get("currency"); currency != "" {
		filter.Currency = currency
	}
	if country := query.Get("country"); country != "" {
		filter.Country = country
	}
	if fuelType := query.Get("fuelType"); fuelType != "" {
		filter.FuelType = fuelType
	}
	if transmission := query.Get("transmission"); transmission != "" {
		filter.Transmission = transmission
	}
	if drivetrain := query.Get("drivetrain"); drivetrain != "" {
		filter.Drivetrain = drivetrain
	}

	// Parse numeric filters
	if minPrice := query.Get("minPrice"); minPrice != "" {
		if val, err := strconv.ParseFloat(minPrice, 64); err == nil {
			filter.MinPrice = val
		}
	}
	if maxPrice := query.Get("maxPrice"); maxPrice != "" {
		if val, err := strconv.ParseFloat(maxPrice, 64); err == nil {
			filter.MaxPrice = val
		}
	}
	if minYear := query.Get("minYear"); minYear != "" {
		if val, err := strconv.Atoi(minYear); err == nil {
			filter.MinYear = val
		}
	}
	if maxYear := query.Get("maxYear"); maxYear != "" {
		if val, err := strconv.Atoi(maxYear); err == nil {
			filter.MaxYear = val
		}
	}

	// Get vehicles
	var vehicles []*models.Vehicle
	if filter.Make != "" || filter.Model != "" || filter.Type != "" || filter.MinPrice > 0 || filter.Currency != "" {
		vehicles = h.repo.SearchVehicles(filter)
	} else {
		vehicles = h.repo.GetAllVehicles()
	}

	response := map[string]interface{}{
		"data":  vehicles,
		"count": len(vehicles),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// HandleGetVehicle returns a single vehicle by ID
func (h *VehicleHandler) HandleGetVehicle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	vehicleID := vars["id"]

	vehicle, err := h.repo.GetVehicleByID(vehicleID)
	if err != nil {
		h.logger.WithField("vehicle_id", vehicleID).Warn("Vehicle not found")
		http.Error(w, "Vehicle not found", http.StatusNotFound)
		return
	}

	response := map[string]interface{}{
		"data": vehicle,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// HandleSearchVehicles handles POST requests for vehicle search
func (h *VehicleHandler) HandleSearchVehicles(w http.ResponseWriter, r *http.Request) {
	var filter models.VehicleFilter
	if err := json.NewDecoder(r.Body).Decode(&filter); err != nil {
		h.logger.WithError(err).Warn("Invalid search request")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	vehicles := h.repo.SearchVehicles(&filter)

	response := map[string]interface{}{
		"data":  vehicles,
		"count": len(vehicles),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
