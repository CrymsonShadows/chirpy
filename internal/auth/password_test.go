package auth

import "testing"

func TestPasswordHashPassesCheck(t *testing.T) {
	password := "Glad tidings"
	hash, _ := HashPassword(password)
	err := CheckPasswordHash(hash, password)
	if err != nil {
		t.Errorf("Password: %s\nResulting Hash: %s\nFailed check", password, hash)
	}
}
