package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/o-zakh/go-http-server-learning/internal/auth"
	"github.com/o-zakh/go-http-server-learning/internal/database"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) createChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Authentication token is missing", err)
		return
	}

	jwtID, err := auth.ValidateJWT(bearerToken, cfg.tokenSecret)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Authentication token is invalid", err)
		return
	}

	const maxChirpLength = 140

	if len(params.Body) > maxChirpLength {
		respondWithError(w, 400, "Chirp is too long", nil)
		return
	} else {
		params.Body = profane_replacement(params.Body)
		dbChirp, err := cfg.dbQueries.CreateChirp(r.Context(), database.CreateChirpParams{
			Body:   params.Body,
			UserID: jwtID,
		})
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't load user parameters", err)
			return
		}
		newChirp := Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body:      dbChirp.Body,
			UserID:    dbChirp.UserID,
		}
		respondWithJSON(w, http.StatusCreated, newChirp)
		return
	}
}

func (cfg *apiConfig) getChirps(w http.ResponseWriter, r *http.Request) {
	author := r.URL.Query().Get("author_id")

	var dbChirps []database.Chirp
	var err error

	if author == "" {
		dbChirps, err = cfg.dbQueries.GetChirps(r.Context())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't load Chirps", err)
			return
		}
	} else {
		author_id, err := uuid.Parse(author)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
			return
		}
		dbChirps, err = cfg.dbQueries.GetAuthorChirps(r.Context(), author_id)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't load Chirps", err)
			return
		}
	}
	allChirps := []Chirp{}
	for _, chirp := range dbChirps {
		newChirp := Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		}
		allChirps = append(allChirps, newChirp)
	}
	respondWithJSON(w, http.StatusOK, allChirps)
}

func (cfg *apiConfig) getChirp(w http.ResponseWriter, r *http.Request) {
	chirpIDString := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}

	chirp, err := cfg.dbQueries.GetChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get chirp", err)
		return
	}

	foundChirp := Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}
	respondWithJSON(w, http.StatusOK, foundChirp)
}

func (cfg *apiConfig) deleteChirp(w http.ResponseWriter, r *http.Request) {

	chirpIDString := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}

	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Authentication token is missing", err)
		return
	}

	jwtID, err := auth.ValidateJWT(bearerToken, cfg.tokenSecret)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Authentication token is invalid", err)
		return
	}

	chirp, err := cfg.dbQueries.GetChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get chirp", err)
		return
	}

	if chirp.UserID == jwtID {
		err = cfg.dbQueries.DeleteChirp(r.Context(), database.DeleteChirpParams{
			ID:     chirp.ID,
			UserID: jwtID,
		})
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't delete Chirp", err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
		return
	} else {
		respondWithError(w, http.StatusForbidden, "Wrong User", err)
		return
	}
}
