package tests

import (
	"animalsys/utils"
	"testing"
)

func TestPasswordHashAndVerify(t *testing.T) {
	pw := "secret123"
	hash, err := utils.HashPassword(pw)
	if err != nil {
		t.Fatal(err)
	}
	ok, err := utils.VerifyPassword(pw, hash)
	if err != nil || !ok {
		t.Fatal("Password verification failed")
	}
}

func TestJWT(t *testing.T) {
	token, err := utils.GenerateJWT("id1", "user1", "admin", "secret")
	if err != nil {
		t.Fatal(err)
	}
	claims, err := utils.ParseJWT(token, "secret")
	if err != nil {
		t.Fatal(err)
	}
	if claims.Username != "user1" || claims.Role != "admin" {
		t.Fatal("JWT claims mismatch")
	}
}
