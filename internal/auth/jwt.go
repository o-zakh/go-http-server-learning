package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy-access",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn).UTC()),
		Subject:   userID.String(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(tokenSecret))
	return ss, err
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(t *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		err = fmt.Errorf("Error while validating JWT: %v", err)
		return uuid.UUID{}, err
	}
	id_str, err := token.Claims.GetSubject()
	if err != nil {
		err = fmt.Errorf("Error while reading Subject ID: %v", err)
		return uuid.UUID{}, err
	}
	id, err := uuid.Parse(id_str)
	if err != nil {
		err = fmt.Errorf("Error while parsing ID: %v", err)
		return uuid.UUID{}, err
	}
	return id, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	bearerToken := headers.Get("Authorization")
	if bearerToken == "" {
		return "", fmt.Errorf("No 'Authorization' header")
	}
	token_string, ok := strings.CutPrefix(bearerToken, "Bearer ")
	if ok {
		return token_string, nil
	} else {
		return "", fmt.Errorf("Error cutting token")
	}
}
