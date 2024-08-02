package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/roman-hushpit/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerTokenRefresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}
	authToken := strings.TrimSpace(strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer"))
	if len(authToken) == 0 {
		respondWithError(w, 401, "Token not present")
		return
	}
	token, err := cfg.DB.GetRefreshToken(authToken)
	if err != nil {
		respondWithError(w, 401, "Token not found")
		return
	}

	if token.ExpiresAt.Before(time.Now()) {
		respondWithError(w, 401, "Token is expired")
		return
	}
	accessToken, err1 := auth.GenerateJwt(token.UserId, cfg.jwsSecter)
	if err1 != nil {
		respondWithError(w, 500, "can not generate access token")
		return
	}
	respondWithJSON(w, 200, response{
		Token : accessToken,
	})
}


func (cfg *apiConfig) handlerTokenRevoke(w http.ResponseWriter, r *http.Request) {
	authToken := strings.TrimSpace(strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer"))
	if len(authToken) == 0 {
		respondWithError(w, 401, "Token not present")
		return
	}
	err := cfg.DB.RevokeRefreshToken(authToken)
	if err != nil {
		respondWithError(w, 500, "can not revoke token")
		return
	}

	respondWithJSON(w, 204, nil)
}