package auth

import (
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestValidToken(t *testing.T) {
	tokenSecret := "secretToken"
	id1 := uuid.New()
	token1, err := MakeJWT(id1, tokenSecret, 5*time.Second)
	if err != nil {
		t.Errorf("making JWT error = %v", err)
	}

	parsedID, err := ValidateJWT(token1, tokenSecret)
	if err != nil {
		t.Errorf("ValidateJWT error: %v", err)
	}
	if parsedID.String() != id1.String() {
		t.Errorf("Expected id: %s\nActual id: %s", id1.String(), parsedID.String())
	}
}

func TestExpiredToken(t *testing.T) {
	tokenSecret := "secretToken"
	id1 := uuid.New()
	token1, err := MakeJWT(id1, tokenSecret, 1*time.Microsecond)
	if err != nil {
		t.Errorf("making JWT error = %v", err)
	}

	time.Sleep(3 * time.Microsecond)

	_, err = ValidateJWT(token1, tokenSecret)
	if !errors.Is(err, jwt.ErrTokenExpired) {
		t.Errorf("ValidateJWT error: %v, expected token to be expired", err)
	}
}

func TestRejectTokenSignedWithWrongSecret(t *testing.T) {
	tokenSecret := "secretToken"
	wrongSecret := "wrongToken"
	id1 := uuid.New()
	token1, err := MakeJWT(id1, tokenSecret, 5*time.Second)
	if err != nil {
		t.Errorf("making JWT error = %v", err)
	}

	_, err = ValidateJWT(token1, wrongSecret)
	if !errors.Is(err, jwt.ErrTokenSignatureInvalid) {
		t.Errorf("ValidateJWT error: %v, expected token signiture to be invalid", err)
	}
}
