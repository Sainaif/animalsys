package tests

import (
	"animalsys/models"
	"encoding/json"
	"testing"
)

func TestUserJSON(t *testing.T) {
	user := models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Role:     "admin",
	}
	b, err := json.Marshal(user)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) == "" {
		t.Fatal("Empty JSON")
	}
}
