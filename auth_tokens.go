package main

import (
	"net/http"
	"time"

	"github.com/o-zakh/go-http-server-learning/internal/auth"
)

func (cfg *apiConfig) refreshToken(w http.ResponseWriter, r *http.Request) {
	bearerRefreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Authorization token is missing", err)
		return
	}
	dbRefreshToken, err := cfg.dbQueries.GetUserFromRefreshToken(r.Context(), bearerRefreshToken)
	if err != nil || dbRefreshToken.RevokedAt.Valid || dbRefreshToken.ExpiresAt.Before(time.Now()) {
		respondWithError(w, http.StatusUnauthorized, "Invalid Authorization token", err)
		return
	}
	newToken, err := auth.MakeJWT(dbRefreshToken.UserID, cfg.tokenSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error while creating an Authrization Token", err)
		return
	}
	respondWithJSON(w, http.StatusOK, struct {
		Token string `json:"token"`
	}{
		Token: newToken,
	})
}

func (cfg *apiConfig) revokeToken(w http.ResponseWriter, r *http.Request) {
	bearerRefreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Authorization token is missing", err)
		return
	}
	err = cfg.dbQueries.RevokeRefreshToken(r.Context(), bearerRefreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid Authorization token", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
