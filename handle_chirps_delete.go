package main

import (
	"net/http"
	"strings"
	"strconv"
	"github.com/roman-hushpit/Chirpy/internal/auth"
)


func (cfg *apiConfig) handlerChirpsDelete(w http.ResponseWriter, r *http.Request) {
	authToken := strings.TrimSpace(strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer"))
	chirpIDString := r.PathValue("chirpID")
	chirpID, _ := strconv.Atoi(chirpIDString)
	
	if len(authToken) == 0 {
		respondWithError(w, 401, "Token not present")
		return
	}

	userId, vError := auth.ValidateToken(authToken, cfg.jwsSecter)
	userIntId, _ := strconv.Atoi(userId)
	if vError != nil {
		respondWithError(w, 401, "Unauthorized")
		return
	}

	dbChipr, err := cfg.DB.GetChirp(chirpID)
	if err != nil {
		respondWithError(w, 404, "Not found")
		return
	}
	if dbChipr.AuthorId != userIntId {
		respondWithError(w, 403, "Access not alowed")
		return
	}

	err1 := cfg.DB.DeleteChipr(chirpID)
	if err1 != nil {
		respondWithError(w, 403, "Can not delete")
		return
	}

	respondWithJSON(w, 204, nil)
}
