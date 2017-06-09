package model

import (
	"encoding/json"
	"math/rand"
	"time"

	"github.com/garyburd/redigo/redis"
)

type Token struct {
	Token     string    `json:"token"`
	ExpiredAt time.Time `json:"expired_at"`
}

type SessionData struct {
	User	  User
	ExpiredAt time.Time `json:"expired_at"`
}

func GenerateToken(u User) Token {
	DeleteSession(u)

	r := getNewTokenString()
	now := time.Now()

	t := Token{
		Token:     r,
		ExpiredAt: now.Add(time.Duration(TOKENLIFE) * time.Minute),
	}

	c := redisPool.Get()
	defer c.Close()

	if _, err := c.Do("SET", "TOKENS"+u.AccountID+u.ID, t.Token); err != nil {
		c.Close()
		panic(err)
	}
	c.Close()

	setSession(u, t)

	return t
}

func setSession(u User, token Token) {
	c := redisPool.Get()
	defer c.Close()

	sd := SessionData{
		User: u,
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

func IsExistToken(token string) bool {
	c := redisPool.Get()
	defer c.Close()

	exists, _ := redis.Bool(c.Do("EXISTS", "SESSION"+token))
	if dt, err := GetSession(token); err != nil {
		exists = false
	} else if dt.ExpiredAt.Before(time.Now()) {
		exists = false
	}

	c.Close()
	return bool(exists)
}

func getToken(u User) (string, error) {
	c := redisPool.Get()
	defer c.Close()
	t, err := redis.String(c.Do("GET", "TOKENS"+u.AccountID+u.ID))
	if err != nil {
		c.Close()
		return "", ErrTokenNotFound
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
		return SessionData{}, ErrTokenNotFound
	}
	c.Close()

	var data SessionData
	if err := json.Unmarshal([]byte(t), &data); err != nil {
		panic(err)
	}
	return data, nil
}

func getNewTokenString() string {
	// generate Random String
	ln := 64
	rand.Seed(time.Now().UTC().UnixNano())
	chars := ALPHANUMERIC
	result := make([]byte, ln)
	for i := 0; i < ln; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}

	rand := string(result)
	if IsExistToken(rand) {
		return getNewTokenString()
	}
	return rand
}

func UpdateTokenExpireTime(token string) {
	sd, _ := GetSession(token)

	now := time.Now()
	sd.ExpiredAt = now.Add(time.Duration(TOKENLIFE) * time.Minute)

	sds, _ := json.Marshal(sd)

	c := redisPool.Get()
	defer c.Close()

	if _, err := c.Do("SET", "SESSION"+token, string(sds)); err != nil {
		c.Close()
		panic(err)
	}
	c.Close()
}

func DeleteSession(u User) {
	t, _ := getToken(u)

	c := redisPool.Get()
	defer c.Close()

	if _, err := c.Do("DEL", "SESSION"+t); err != nil {
		c.Close()
		panic(err)
	}

	if _, err := c.Do("DEL", "TOKENS"+u.AccountID+u.ID); err != nil {
		c.Close()
		panic(err)
	}
}

//FORGOT PASSWORD TOKEN
//
// func GenerateForgotPasswordToken(AccountID, userId string) Token {
// 	DeleteForgotPasswordSession(AccountID, userId)
//
// 	r := getNewTokenString()
// 	now := time.Now()
//
// 	t := Token{
// 		Token:     r,
// 		ExpiredAt: now.Add(time.Duration(TOKENLIFE) * time.Minute),
// 	}
//
// 	c := redisPool.Get()
// 	defer c.Close()
//
// 	if _, err := c.Do("SET", "FORGOT_PASSWORD_TOKENS"+AccountID+userId, t.Token); err != nil {
// 		c.Close()
// 		panic(err)
// 	}
// 	c.Close()
//
// 	setForgotPasswordSession(AccountID, userId, t)
//
// 	return t
// }

// func setForgotPasswordSession(AccountID string, userId string, token Token) {
// 	c := redisPool.Get()
// 	defer c.Close()
//
// 	sd := SessionData{
// 		UserId:    userId,
// 		AccountID: AccountID,
// 		ExpiredAt: token.ExpiredAt,
// 	}
//
// 	json, _ := json.Marshal(sd)
//
// 	if _, err := c.Do("DEL", "FORGOT_PASSWORD_SESSION"+token.Token); err != nil {
// 		c.Close()
// 		panic(err)
// 	}
//
// 	if _, err := c.Do("SET", "FORGOT_PASSWORD_SESSION"+token.Token, json); err != nil {
// 		c.Close()
// 		panic(err)
// 	}
// }
//
// func AuthenticateForgotPasswordToken(token string) error {
// 	if !isExistForgotPasswordToken(token) {
// 		fmt.Println("err 1")
// 		return ErrTokenNotFound
// 	}
//
// 	sd, err := GetForgotPasswordSession(token)
// 	if err != nil {
// 		fmt.Println("err 2")
// 		return ErrTokenNotFound
// 	}
//
// 	if sd.ExpiredAt.Before(time.Now()) {
// 		fmt.Println("err 3")
// 		return ErrTokenExpired
// 	}
//
// 	return nil
// }
//
// func isExistForgotPasswordToken(token string) bool {
// 	c := redisPool.Get()
// 	defer c.Close()
//
// 	exists, _ := redis.Bool(c.Do("EXISTS", "FORGOT_PASSWORD_SESSION"+token))
//
// 	c.Close()
// 	return bool(exists)
// }
//
// func getForgotPasswordToken(AccountID, userId string) (string, error) {
// 	c := redisPool.Get()
// 	defer c.Close()
// 	t, err := redis.String(c.Do("GET", "FORGOT_PASSWORD_TOKENS"+AccountID+userId))
// 	if err != nil {
// 		c.Close()
// 		return "", ErrTokenNotFound
// 	}
// 	c.Close()
//
// 	return t, nil
// }
//
// func GetForgotPasswordSession(token string) (SessionData, error) {
// 	c := redisPool.Get()
// 	defer c.Close()
// 	t, err := redis.String(c.Do("GET", "FORGOT_PASSWORD_SESSION"+token))
// 	if err != nil {
// 		c.Close()
// 		return SessionData{}, ErrTokenNotFound
// 	}
// 	c.Close()
//
// 	var data SessionData
// 	if err := json.Unmarshal([]byte(t), &data); err != nil {
// 		panic(err)
// 	}
// 	return data, nil
// }
//
// func DeleteForgotPasswordSession(AccountID, userId string) {
// 	t, _ := getForgotPasswordToken(AccountID, userId)
//
// 	c := redisPool.Get()
// 	defer c.Close()
//
// 	if _, err := c.Do("DEL", "FORGOT_PASSWORD_SESSION"+t); err != nil {
// 		c.Close()
// 		panic(err)
// 	}
//
// 	if _, err := c.Do("DEL", "FORGOT_PASSWORD_TOKENS"+AccountID+userId); err != nil {
// 		c.Close()
// 		panic(err)
// 	}
// }
