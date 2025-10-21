package tests

import (
	"animalsys/utils"
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "testpassword123"
	hash, err := utils.HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	if !utils.CheckPasswordHash(password, hash) {
		t.Error("Password hash verification failed")
	}

	if utils.CheckPasswordHash("wrongpassword", hash) {
		t.Error("Invalid password should not match")
	}
}

func TestGenerateJWT(t *testing.T) {
	secret := "testsecret"
	userID := "507f1f77bcf86cd799439011"
	username := "testuser"
	role := "user"

	token, err := utils.GenerateJWT(userID, username, role, secret)
	if err != nil {
		t.Fatalf("Failed to generate JWT: %v", err)
	}

	if token == "" {
		t.Error("Generated token is empty")
	}

	claims, err := utils.ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("Failed to validate JWT: %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("Expected UserID %s, got %s", userID, claims.UserID)
	}

	if claims.Username != username {
		t.Errorf("Expected Username %s, got %s", username, claims.Username)
	}

	if claims.Role != role {
		t.Errorf("Expected Role %s, got %s", role, claims.Role)
	}
}

func TestValidateJWT_InvalidToken(t *testing.T) {
	secret := "testsecret"
	invalidToken := "invalid.token.here"

	_, err := utils.ValidateJWT(invalidToken, secret)
	if err == nil {
		t.Error("Expected error for invalid token, got nil")
	}
}

func TestValidateJWT_WrongSecret(t *testing.T) {
	secret := "testsecret"
	wrongSecret := "wrongsecret"
	token, _ := utils.GenerateJWT("123", "test", "user", secret)

	_, err := utils.ValidateJWT(token, wrongSecret)
	if err == nil {
		t.Error("Expected error for wrong secret, got nil")
	}
}
