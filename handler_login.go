package main

import (
	"encoding/json"
	"net/http"

	"github.com/roman-hushpit/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
		ExpiresInSeconds uint 	`json:"expires_in_seconds"`
	}
	type response struct {
		User
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	user, err := cfg.DB.GetUserByEmail(params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get user")
		return
	}

	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid password")
		return
	}

	jwtToken, err := auth.GenerateJwt(user.ID, cfg.jwsSecter)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create auth token")
		return
	}

	refreshToken := auth.GenerateRefreshToken()

	refreshTokenDto, err := cfg.DB.CreateRefreshToken(refreshToken, &user)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create auth token")
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:    user.ID,
			Email: user.Email,
			Token: jwtToken,
			RefreshToken: refreshTokenDto.ID,
			IsChirpyRed: user.IsChirpyRed,
		},
	})
}
