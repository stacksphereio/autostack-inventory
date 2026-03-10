package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/CB-AutoStack/AutoStack/apps/api-inventory/internal/auth"
	"github.com/CB-AutoStack/AutoStack/apps/api-inventory/internal/repository"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// AuthHandler handles authentication requests
type AuthHandler struct {
	repo       *repository.Repository
	jwtManager *auth.JWTManager
	logger     *logrus.Logger
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(repo *repository.Repository, jwtManager *auth.JWTManager, logger *logrus.Logger) *AuthHandler {
	return &AuthHandler{
		repo:       repo,
		jwtManager: jwtManager,
		logger:     logger,
	}
}

// LoginRequest represents a login request
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse represents a login response
type LoginResponse struct {
	Token string      `json:"token"`
	User  interface{} `json:"user"`
}

// HandleLogin handles user login
func (h *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.WithError(err).Warn("Invalid login request")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Find user by email
	user, err := h.repo.GetUserByEmail(req.Email)
	if err != nil {
		h.logger.WithField("email", req.Email).Warn("User not found")
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		h.logger.WithField("email", req.Email).Warn("Invalid password")
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token, err := h.jwtManager.GenerateToken(user.ID, user.Email)
	if err != nil {
		h.logger.WithError(err).Error("Failed to generate token")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Return token and user info (without password)
	userResponse := map[string]interface{}{
		"id":                user.ID,
		"email":             user.Email,
		"name":              user.Name,
		"country":           user.Country,
		"preferredCurrency": user.PreferredCurrency,
		"roles":             user.Roles,
	}

	response := LoginResponse{
		Token: token,
		User:  userResponse,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

	h.logger.WithFields(logrus.Fields{
		"user_id": user.ID,
		"email":   user.Email,
	}).Info("User logged in")
}
