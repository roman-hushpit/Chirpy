package database

import (
	"time"
	"errors"
)

type Token struct {
	ID 		 string     `json:"id"`
	UserId   int   	    `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
}

func (db *DB) CreateRefreshToken(refreshToken string, user *User) (Token, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Token{}, err
	}
	now := time.Now()
	expirationTime := now.Add(time.Duration(60 * 24) * time.Hour)
	token := Token{
		ID:             refreshToken,
		UserId:          user.ID,
		ExpiresAt:   expirationTime,
	}
	dbStructure.Tokens[refreshToken] = token

	err = db.writeDB(dbStructure)
	if err != nil {
		return Token{}, err
	}

	return token, nil
}

func (db *DB) GetRefreshToken(refreshToken string) (Token, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Token{}, err
	}
	token, ok := dbStructure.Tokens[refreshToken] 
	if !ok  {
		return Token{}, errors.New("Token not found")
	}
	return token, nil
}

func (db *DB) RevokeRefreshToken(refreshToken string) (error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}
	delete(dbStructure.Tokens, refreshToken)
	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}
	return nil
}