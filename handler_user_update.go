package main

import (
	"net/http"
	"strings"
	"github.com/roman-hushpit/Chirpy/internal/auth"
	"encoding/json"
	"strconv"
)

func (cfg *apiConfig) handlerUsersUpdate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password 		 string `json:"password"`
		Email    		 string `json:"email"`
	}
	type response struct {
		User
	}
	authToken := strings.TrimSpace(strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer"))
	if len(authToken) == 0 {
		respondWithError(w, 401, "Token not present")
		return
	}

	userId, err := auth.ValidateToken(authToken, cfg.jwsSecter)
	if err != nil {
		respondWithError(w, 401, "Unauthorized")
		return
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	decodeError := decoder.Decode(&params)
	if decodeError != nil {
		respondWithError(w, 404, "Can not parse request")
		return
	}
	id, err := strconv.Atoi(userId)
	if err != nil {
		respondWithError(w, 404, "Invalid id")
		return
	}
	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password")
		return
	}
	existingUser, _ := cfg.DB.GetUser(id)
	existingUser.Email = params.Email
	existingUser.HashedPassword = hashedPassword

	user, err := cfg.DB.UpdateUser(&existingUser)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password")
		return
	}
	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:    user.ID,
			Email: user.Email,
			Token: authToken,
			IsChirpyRed: user.IsChirpyRed,
		},
	})
	
}