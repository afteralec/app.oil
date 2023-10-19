package main

import (
	"testing"
)

func TestPassword(t *testing.T) {
	password := "test"

	hash, err := HashPassword(password)
	if err != nil {
		t.Errorf("Password hashing failed: %v", err)
	}

	verified, err := VerifyPassword(password, hash)
	if err != nil {
		t.Errorf("Password verification failed: %v", err)
	}

	if verified != true {
		t.Errorf("Password verification returned false")
	}
}
