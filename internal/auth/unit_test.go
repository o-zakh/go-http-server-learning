package auth

import (
	"testing"
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
