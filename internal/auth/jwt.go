package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	signingKey := []byte(tokenSecret)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   userID.String(),
	})
	return token.SignedString(signingKey)
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	parsedToken, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("error parsing token: %w", err)
	}

	if !parsedToken.Valid {
		return uuid.Nil, fmt.Errorf("token not valid: %w", err)
	}

	idString, err := parsedToken.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, fmt.Errorf("error getting subject of token: %w", err)
	}

	issuer, err := parsedToken.Claims.GetIssuer()
	if err != nil {
		return uuid.Nil, err
	}
	if issuer != string("chirpy") {
		return uuid.Nil, errors.New("invalid issuer")
	}

	id, err := uuid.Parse(idString)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user id: %w", err)
	}

	return id, nil
}
