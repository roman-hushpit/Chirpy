package main

import (
	"encoding/json"
	"net/http"
	"strings"
)


func (cfg *apiConfig) handlerWebhook(w http.ResponseWriter, r *http.Request) {
	authToken := strings.TrimSpace(strings.TrimPrefix(r.Header.Get("Authorization"), "ApiKey"))
	if len(authToken) == 0 {
		respondWithError(w, 401, "Token not present")
		return
	}
	if authToken != cfg.polkaApiKey {
		respondWithError(w, 401, "Not authorized")
		return
	}


	type data struct {
		UserId int `json:"user_id"`
	}
	type parameters struct {
		Event string `json:"event"`
		Data data `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}
	if params.Event != "user.upgraded" {
		respondWithJSON(w, 204, nil)
		return
	}

	userId := params.Data.UserId
	user, err := cfg.DB.GetUser(userId)
	if err != nil {
		respondWithError(w, 404, "User not found")
		return
	}
	user.IsChirpyRed = true
	_, error := cfg.DB.UpdateUser(&user)
	if error != nil {
		respondWithError(w, 500, "Can not update user")
		return
	}
	respondWithJSON(w, 204, nil)
}
