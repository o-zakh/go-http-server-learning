package auth

import (
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestUser1(t *testing.T) {
	password := "123456789"
	hashedPassword, err := HashPassword(password)

	if err != nil {
		t.Errorf(`Error while hashing passwords = %v, expected nil`, err)
	}

	match, err := CheckPasswordHash(password, hashedPassword)

	want := true

	if want != match || err != nil {
		t.Errorf(`CheckPasswordHash = %t, %v, want match for %t, nil`, match, err, want)
	}
}

func TestUser2(t *testing.T) {
	password := "22@wqtW9pS?8"
	hashedPassword, err := HashPassword(password)

	if err != nil {
		t.Errorf(`Error while hashing passwords = %v, expected nil`, err)
	}

	match, err := CheckPasswordHash(password, hashedPassword)

	want := true

	if want != match || err != nil {
		t.Errorf(`CheckPasswordHash = %t, %v, want match for %t, nil`, match, err, want)
	}
}

func TestUser3(t *testing.T) {
	password := "1PFC8A1or82i"
	hashedPassword, err := HashPassword(password)

	if err != nil {
		t.Errorf(`Error while hashing passwords = %v, expected nil`, err)
	}

	match, err := CheckPasswordHash(password, hashedPassword)

	want := true

	if want != match || err != nil {
		t.Errorf(`CheckPasswordHash = %t, %v, want match for %t, nil`, match, err, want)
	}
}

func TestUser_WrongPassword(t *testing.T) {
	password := "correct"

	hashedPassword, err := HashPassword(password)

	if err != nil {
		t.Errorf(`Error while hashing passwords = %v, expected nil`, err)
	}

	inputPassword := "wrong"
	match, err := CheckPasswordHash(inputPassword, hashedPassword)

	want := false

	if want != match || err != nil {
		t.Errorf(`CheckPasswordHash = %t, %v, want match for %t, nil`, match, err, want)
	}
}

// JWT Tests:

func TestJWT1(t *testing.T) {
	originalID := uuid.New()
	tokenSecret := "tokenSecret"
	duration, err := time.ParseDuration("1h")
	if err != nil {
		t.Errorf(`Error while setting duration = %v, expected nil`, err)
	}
	newJWT, err := MakeJWT(originalID, tokenSecret, duration)
	if err != nil {
		t.Errorf(`Error while signing a JWT = %v, expected nil`, err)
	}
	resultID, err := ValidateJWT(newJWT, tokenSecret)
	if err != nil || resultID != originalID {
		t.Errorf(`Error while validating a JWT. Original User ID = %v; Resulting ID = %v. Error: %v`, originalID, resultID, err)
	}
}

func TestJWT2(t *testing.T) {
	originalID := uuid.New()
	tokenSecret := "newSecret"
	duration, err := time.ParseDuration("20h")
	if err != nil {
		t.Errorf(`Error while setting duration = %v, expected nil`, err)
	}
	newJWT, err := MakeJWT(originalID, tokenSecret, duration)
	if err != nil {
		t.Errorf(`Error while signing a JWT = %v, expected nil`, err)
	}
	resultID, err := ValidateJWT(newJWT, tokenSecret)
	if err != nil || resultID != originalID {
		t.Errorf(`Error while validating a JWT. Original User ID = %v; Resulting ID = %v. Error: %v`, originalID, resultID, err)
	}
}

func TestJWT3(t *testing.T) {
	originalID := uuid.New()
	tokenSecret := "anotherSecret"
	duration, err := time.ParseDuration("1h")
	if err != nil {
		t.Errorf(`Error while setting duration = %v, expected nil`, err)
	}
	newJWT, err := MakeJWT(originalID, tokenSecret, duration)
	if err != nil {
		t.Errorf(`Error while signing a JWT = %v, expected nil`, err)
	}
	resultID, err := ValidateJWT(newJWT, tokenSecret)
	if err != nil || resultID != originalID {
		t.Errorf(`Error while validating a JWT. Original User ID = %v; Resulting ID = %v. Error: %v`, originalID, resultID, err)
	}
}

func TestJWT_ExpiredToken(t *testing.T) {
	originalID := uuid.New()
	tokenSecret := "anotherSecret"
	duration, err := time.ParseDuration("-1h")
	if err != nil {
		t.Errorf(`Error while setting duration = %v, expected nil`, err)
	}
	newJWT, err := MakeJWT(originalID, tokenSecret, duration)
	if err != nil {
		t.Errorf(`Error while signing a JWT = %v, expected nil`, err)
	}
	resultID, err := ValidateJWT(newJWT, tokenSecret)
	if !strings.Contains(err.Error(), "token is expired") || resultID != (uuid.UUID{}) {
		t.Errorf(`Missed expired token detection. Original User ID = %v; Resulting ID = %v. Error: %v`, originalID, resultID, err)
	}
}

func TestJWT_WrongToken(t *testing.T) {
	originalID := uuid.New()
	tokenSecret := "anotherSecret"
	duration, err := time.ParseDuration("1h")
	if err != nil {
		t.Errorf(`Error while setting duration = %v, expected nil`, err)
	}
	newJWT, err := MakeJWT(originalID, tokenSecret, duration)
	if err != nil {
		t.Errorf(`Error while signing a JWT = %v, expected nil`, err)
	}
	fakeTokenSecret := "faketoken"
	resultID, err := ValidateJWT(newJWT, fakeTokenSecret)
	if !strings.Contains(err.Error(), "token signature is invalid") || resultID != (uuid.UUID{}) {
		t.Errorf(`Missed fake token signature. Original Token Secret = %v; Resulting Token Secret = %v. Error: %v`, tokenSecret, fakeTokenSecret, err)
	}
}

func TestValidAuthHeader(t *testing.T) {
	header := http.Header{}
	authToken := "Bearer 1234"
	header.Add("Accept", "value")
	header.Add("Authorization", authToken)
	token, err := GetBearerToken(header)
	want := "1234"
	if token != want || err != nil {
		t.Errorf(`Error getting or cutting a JWT. Expected Token: %v, Result: %v, Error: %v`, want, token, err)
	}
}

func TestNoAuthHeader(t *testing.T) {
	header := http.Header{}
	authToken := "Bearer 1234"
	header.Add("Accept", authToken)
	token, err := GetBearerToken(header)
	want := ""
	if token != want || !strings.Contains(err.Error(), "No 'Authorization'") {
		t.Errorf(`Error getting or cutting a JWT. Expected Token: %v, Result: %v, Error: %v`, want, token, err)
	}
}
