package main

import (
	"encoding/json"
	"net/http"
)

func validate_handler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	const maxChirpLength = 140

	if len(params.Body) > maxChirpLength {
		respondWithError(w, 400, "Chirp is too long", nil)
	} else {
		param1 := profane_replacement(params.Body)
		respondWithJSON(w, 200, returnVals{
			CleanedBody: param1,
		})
	}
}
