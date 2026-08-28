package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/o-zakh/go-http-server-learning/internal/auth"
	"github.com/o-zakh/go-http-server-learning/internal/database"
)

func (cfg *apiConfig) webhook(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	apiKey, err := auth.GetAPIKey(r.Header)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid API Key", err)
		return
	}

	if cfg.polkaKey != apiKey {
		respondWithError(w, http.StatusInternalServerError, "Error while extracting API Key", err)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	err = cfg.dbQueries.UserUpgrade(r.Context(), database.UserUpgradeParams{
		ID:          params.Data.UserID,
		IsChirpyRed: true,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't Upgrade User", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
