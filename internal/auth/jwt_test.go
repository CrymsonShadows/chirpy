package auth

import (
	"errors"
	"fmt"
	"net/http"
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

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	validToken, _ := MakeJWT(userID, "secret", time.Hour)

	tests := []struct {
		name        string
		tokenString string
		tokenSecret string
		wantUserID  uuid.UUID
		wantErr     bool
	}{
		{
			name:        "Valid token",
			tokenString: validToken,
			tokenSecret: "secret",
			wantUserID:  userID,
			wantErr:     false,
		},
		{
			name:        "Invalid token",
			tokenString: "invalid.token.string",
			tokenSecret: "secret",
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
		{
			name:        "Wrong secret",
			tokenString: validToken,
			tokenSecret: "wrong_secret",
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUserID, err := ValidateJWT(tt.tokenString, tt.tokenSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotUserID != tt.wantUserID {
				t.Errorf("ValidateJWT() gotUserID = %v, want %v", gotUserID, tt.wantUserID)
			}
		})
	}
}

func TestGetBearerToken(t *testing.T) {
	validToken := "validToken"

	tests := []struct {
		name        string
		tokenString string
		header      http.Header
		wantErr     bool
	}{
		{
			name:        "No Authorization header",
			tokenString: validToken,
			header:      http.Header{},
			wantErr:     true,
		},
		{
			name:        "Authorization header has token",
			tokenString: validToken,
			header:      http.Header{"Authorization": []string{fmt.Sprintf("Bearer %s", validToken)}},
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GetBearerToken(tt.header)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBearerToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (token == "") != tt.wantErr {
				t.Errorf("GetBearerToken() returned = %s, want %s, err = %v", token, validToken, err)
			}
		})
	}
}
