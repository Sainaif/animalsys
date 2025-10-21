package tests

import (
	"animalsys/models"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestUserModel(t *testing.T) {
	user := models.User{
		ID:           primitive.NewObjectID(),
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
		Role:         "user",
	}

	if user.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got '%s'", user.Username)
	}

	if user.Email != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got '%s'", user.Email)
	}

	if user.Role != "user" {
		t.Errorf("Expected role 'user', got '%s'", user.Role)
	}
}
