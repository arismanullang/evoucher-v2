package model

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Token struct {
	Token     string    `json:"token"`
	ExpiredAt time.Time `json:"expired_at"`
}

type SessionData struct {
	User      User
	ExpiredAt time.Time `json:"expired_at"`
}

func GenerateToken(u User) Token {
	now := time.Now()
	exp := time.Now().Add(24 * time.Hour)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: exp.Unix(),
		IssuedAt:  now.Unix(),
		Audience:  u.Account.Id,
		Issuer:    "voucher",
		Subject:   u.ID,
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("AUTH_SECRET_KEY")))
	if err != nil {
		log.Panic(err)
	}

	return Token{
		Token:     tokenString,
		ExpiredAt: exp,
	}
}

func GetSession(tokenString string) (SessionData, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("AUTH_SECRET_KEY")), nil
	})
	if err != nil {
		return SessionData{}, err
	}
	if !token.Valid {
		return SessionData{}, errors.New("invalid token")
	}
	if err := token.Claims.Valid(); err != nil {
		return SessionData{}, err
	}

	user, err := FindUserDetail(token.Claims.(*jwt.StandardClaims).Subject)
	if err != nil {
		return SessionData{}, err
	}

	return SessionData{
		User:      user,
		ExpiredAt: time.Unix(token.Claims.(*jwt.StandardClaims).ExpiresAt, 0),
	}, nil
}
