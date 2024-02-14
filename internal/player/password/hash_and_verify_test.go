package password

import (
	"testing"
)

func TestHashAndVerify(t *testing.T) {
	password := "test"

	hash, err := Hash(password)
	if err != nil {
		t.Errorf("Password hashing failed: %v", err)
	}

	verified, err := Verify(password, hash)
	if err != nil {
		t.Errorf("Password verification failed: %v", err)
	}

	if verified != true {
		t.Errorf("Password verification returned false")
	}
}
