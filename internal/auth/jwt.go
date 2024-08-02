
package auth

import (
	"time"
	"github.com/golang-jwt/jwt/v5"
	"strconv"
	"log"
	"errors"
	"crypto/rand"
	"encoding/hex"
)

func GenerateJwt(id int, apiKey string) (string, error) {
	now := time.Now()
	expirationTime := now.Add(time.Duration(1) * time.Hour)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer : "chirpy",
		IssuedAt : jwt.NewNumericDate(now),
		ExpiresAt : jwt.NewNumericDate(expirationTime), 
		Subject : strconv.Itoa(id),
	})
	signedToken, err := token.SignedString([]byte(apiKey))
	if err != nil {
		log.Printf("Can not create jwt: %v\n", err)
		return "", err
	}

	return signedToken, nil
}

func ValidateToken(tokenString string, apiKey string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(apiKey), nil
	});
	if err != nil {
		log.Print(err)
		return "", err
	} else if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok {
		return claims.Subject, nil
	} else {
		return "", errors.New("can not parse claims")
	}
}

func GenerateRefreshToken() string {
	c := 10
	b := make([]byte, c)
	rand.Read(b)
	return hex.EncodeToString(b)
}