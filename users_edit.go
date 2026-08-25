package main

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/o-zakh/go-http-server-learning/internal/auth"
	"github.com/o-zakh/go-http-server-learning/internal/database"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
}

func (cfg *apiConfig) createUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash new user password", err)
		return
	}

	dbUser, err := cfg.dbQueries.CreateUser(r.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPassword,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't load user parameters", err)
		return
	}
	newUser := User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
	}
	respondWithJSON(w, http.StatusCreated, newUser)
}

func (cfg *apiConfig) loginUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password         string `json:"password"`
		Email            string `json:"email"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	dbUser, err := cfg.dbQueries.GetUser(r.Context(), params.Email)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "User not found", err)
		return
	}

	pswd_match, err := auth.CheckPasswordHash(params.Password, dbUser.HashedPassword)

	if pswd_match {
		if params.ExpiresInSeconds > 3600 || params.ExpiresInSeconds == 0 {
			params.ExpiresInSeconds = 3600
		}
		newToken, err := auth.MakeJWT(dbUser.ID, cfg.tokenSecret, time.Duration(params.ExpiresInSeconds)*time.Second)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error generating a JWT", err)
			return
		}
		respondWithJSON(w, http.StatusOK, User{
			ID:        dbUser.ID,
			CreatedAt: dbUser.CreatedAt,
			UpdatedAt: dbUser.UpdatedAt,
			Email:     dbUser.Email,
			Token:     newToken,
		})
	} else {
		respondWithError(w, http.StatusUnauthorized, "User not found", err)
		return
	}
}

func (cfg *apiConfig) deleteUser(w http.ResponseWriter, r *http.Request) {
	if os.Getenv("PLATFORM") != "dev" {
		respondWithError(w, http.StatusForbidden, "Only accessible via an admin device", nil)
		return
	}
	cfg.dbQueries.DeleteUsers(r.Context())
}
