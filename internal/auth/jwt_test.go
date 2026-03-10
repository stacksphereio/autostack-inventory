package auth

import (
	"testing"
	"time"
)

func TestGenerateToken(t *testing.T) {
	manager := NewJWTManager("test-secret", 24*time.Hour)

	token, err := manager.GenerateToken("user-001", "test@autostack.com")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	if token == "" {
		t.Error("Generated token is empty")
	}

	t.Logf("Generated token: %s", token[:50]+"...")
}

func TestValidateToken(t *testing.T) {
	manager := NewJWTManager("test-secret", 24*time.Hour)

	// Generate a token
	token, err := manager.GenerateToken("user-001", "test@autostack.com")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Validate the token
	claims, err := manager.ValidateToken(token)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}

	if claims.UserID != "user-001" {
		t.Errorf("Expected UserID user-001, got %s", claims.UserID)
	}

	if claims.Email != "test@autostack.com" {
		t.Errorf("Expected Email test@autostack.com, got %s", claims.Email)
	}

	t.Logf("Validated token for user: %s (%s)", claims.UserID, claims.Email)
}

func TestValidateInvalidToken(t *testing.T) {
	manager := NewJWTManager("test-secret", 24*time.Hour)

	// Test with invalid token
	_, err := manager.ValidateToken("invalid.token.here")
	if err == nil {
		t.Error("Expected error for invalid token")
	}
	t.Logf("Correctly rejected invalid token: %v", err)
}

func TestValidateExpiredToken(t *testing.T) {
	// Create manager with very short duration
	manager := NewJWTManager("test-secret", 1*time.Millisecond)

	token, err := manager.GenerateToken("user-001", "test@autostack.com")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Wait for token to expire
	time.Sleep(100 * time.Millisecond)

	// Try to validate expired token
	_, err = manager.ValidateToken(token)
	if err == nil {
		t.Error("Expected error for expired token")
	}
	t.Logf("Correctly rejected expired token: %v", err)
}

func TestTokenWithDifferentSecrets(t *testing.T) {
	manager1 := NewJWTManager("secret-one", 24*time.Hour)
	manager2 := NewJWTManager("secret-two", 24*time.Hour)

	// Generate token with manager1
	token, err := manager1.GenerateToken("user-001", "test@autostack.com")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Try to validate with manager2 (different secret)
	_, err = manager2.ValidateToken(token)
	if err == nil {
		t.Error("Expected error when validating token with different secret")
	}
	t.Logf("Correctly rejected token signed with different secret: %v", err)
}
