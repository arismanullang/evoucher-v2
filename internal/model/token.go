package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/garyburd/redigo/redis"
)

var (
	letterRunes     = []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has been expired")
	TokenLife       int
)

type Token struct {
	Token     string    `json:"token"`
	ExpiredAt time.Time `json:"expired_at"`
}

type SessionData struct {
	UserId    string    `json:"user_id"`
	CompanyId string    `json:"company_id"`
	ExpiredAt time.Time `json:"expired_at"`
}

func GenerateToken(companyId, userId string) Token {
	DeleteSession(companyId, userId)

	r := getNewTokenString()
	now := time.Now()

	t := Token{
		Token:     r,
		ExpiredAt: now.Add(time.Duration(TokenLife) * time.Minute),
	}

	c := redisPool.Get()
	defer c.Close()

	if _, err := c.Do("SET", "TOKENS"+companyId+userId, t.Token); err != nil {
		c.Close()
		panic(err)
	}
	c.Close()

	setSession(companyId, userId, t)

	return t
}

func setSession(companyId string, userId string, token Token) {
	c := redisPool.Get()
	defer c.Close()

	sd := SessionData{
		UserId:    userId,
		CompanyId: companyId,
		ExpiredAt: token.ExpiredAt,
	}

	json, _ := json.Marshal(sd)

	if _, err := c.Do("DEL", "SESSION"+token.Token); err != nil {
		c.Close()
		panic(err)
	}

	if _, err := c.Do("SET", "SESSION"+token.Token, json); err != nil {
		c.Close()
		panic(err)
	}
}

func AuthenticateToken(token string) error {
	if !isExistToken(token) {
		return ErrInvalidToken
	}

	sd, err := GetSession(token)
	if err != nil {
		return ErrInvalidToken
	}

	if sd.ExpiredAt.Before(time.Now()) {
		return ErrExpiredToken
	}

	return nil
}

func isExistToken(token string) bool {
	c := redisPool.Get()
	defer c.Close()

	exists, _ := redis.Bool(c.Do("EXISTS", "SESSION"+token))

	c.Close()
	return bool(exists)
}

func getToken(companyId, userId string) (string, error) {
	c := redisPool.Get()
	defer c.Close()
	t, err := redis.String(c.Do("GET", "TOKENS"+companyId+userId))
	if err != nil {
		c.Close()
		return "", ErrInvalidToken
	}
	c.Close()

	return t, nil
}

func GetSession(token string) (SessionData, error) {
	c := redisPool.Get()
	defer c.Close()
	t, err := redis.String(c.Do("GET", "SESSION"+token))
	if err != nil {
		c.Close()
		return SessionData{}, ErrInvalidToken
	}
	c.Close()

	var data SessionData
	if err := json.Unmarshal([]byte(t), &data); err != nil {
		panic(err)
	}
	return data, nil
}

func getNewTokenString() string {
	rand := RandStringRunes(39, letterRunes)
	if isExistToken(rand) {
		return getNewTokenString()
	}
	return rand
}

func UpdateTokenExpireTime(token string) {
	sd, _ := GetSession(token)

	now := time.Now()
	sd.ExpiredAt = now.Add(time.Duration(TokenLife) * time.Minute)

	sds, _ := json.Marshal(sd)

	c := redisPool.Get()
	defer c.Close()

	if _, err := c.Do("SET", "SESSION"+token, string(sds)); err != nil {
		c.Close()
		panic(err)
	}
	c.Close()
}

func DeleteSession(companyId, userId string) {
	t, _ := getToken(companyId, userId)

	c := redisPool.Get()
	defer c.Close()

	if _, err := c.Do("DEL", "SESSION"+t); err != nil {
		c.Close()
		panic(err)
	}

	if _, err := c.Do("DEL", "TOKENS"+companyId+userId); err != nil {
		c.Close()
		panic(err)
	}
}

//FORGOT PASSWORD TOKEN

func GenerateForgotPasswordToken(companyId, userId string) Token {
	DeleteForgotPasswordSession(companyId, userId)

	r := getNewTokenString()
	now := time.Now()

	t := Token{
		Token:     r,
		ExpiredAt: now.Add(time.Duration(TokenLife) * time.Minute),
	}

	c := redisPool.Get()
	defer c.Close()

	if _, err := c.Do("SET", "FORGOT_PASSWORD_TOKENS"+companyId+userId, t.Token); err != nil {
		c.Close()
		panic(err)
	}
	c.Close()

	setForgotPasswordSession(companyId, userId, t)

	return t
}

func setForgotPasswordSession(companyId string, userId string, token Token) {
	c := redisPool.Get()
	defer c.Close()

	sd := SessionData{
		UserId:    userId,
		CompanyId: companyId,
		ExpiredAt: token.ExpiredAt,
	}

	json, _ := json.Marshal(sd)

	if _, err := c.Do("DEL", "FORGOT_PASSWORD_SESSION"+token.Token); err != nil {
		c.Close()
		panic(err)
	}

	if _, err := c.Do("SET", "FORGOT_PASSWORD_SESSION"+token.Token, json); err != nil {
		c.Close()
		panic(err)
	}
}

func AuthenticateForgotPasswordToken(token string) error {
	if !isExistForgotPasswordToken(token) {
		fmt.Println("err 1")
		return ErrInvalidToken
	}

	sd, err := GetForgotPasswordSession(token)
	if err != nil {
		fmt.Println("err 2")
		return ErrInvalidToken
	}

	if sd.ExpiredAt.Before(time.Now()) {
		fmt.Println("err 3")
		return ErrExpiredToken
	}

	return nil
}

func isExistForgotPasswordToken(token string) bool {
	c := redisPool.Get()
	defer c.Close()

	exists, _ := redis.Bool(c.Do("EXISTS", "FORGOT_PASSWORD_SESSION"+token))

	c.Close()
	return bool(exists)
}

func getForgotPasswordToken(companyId, userId string) (string, error) {
	c := redisPool.Get()
	defer c.Close()
	t, err := redis.String(c.Do("GET", "FORGOT_PASSWORD_TOKENS"+companyId+userId))
	if err != nil {
		c.Close()
		return "", ErrInvalidToken
	}
	c.Close()

	return t, nil
}

func GetForgotPasswordSession(token string) (SessionData, error) {
	c := redisPool.Get()
	defer c.Close()
	t, err := redis.String(c.Do("GET", "FORGOT_PASSWORD_SESSION"+token))
	if err != nil {
		c.Close()
		return SessionData{}, ErrInvalidToken
	}
	c.Close()

	var data SessionData
	if err := json.Unmarshal([]byte(t), &data); err != nil {
		panic(err)
	}
	return data, nil
}

func DeleteForgotPasswordSession(companyId, userId string) {
	t, _ := getForgotPasswordToken(companyId, userId)

	c := redisPool.Get()
	defer c.Close()

	if _, err := c.Do("DEL", "FORGOT_PASSWORD_SESSION"+t); err != nil {
		c.Close()
		panic(err)
	}

	if _, err := c.Do("DEL", "FORGOT_PASSWORD_TOKENS"+companyId+userId); err != nil {
		c.Close()
		panic(err)
	}
}

func RandStringRunes(n int, letterRunes []rune) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
