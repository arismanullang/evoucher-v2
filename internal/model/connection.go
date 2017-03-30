package model

import (
	"github.com/garyburd/redigo/redis"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	redistore "gopkg.in/boj/redistore.v1"
)

var db *sqlx.DB
var redisPool *redis.Pool

func ConnectDB(endpoint string) (err error) {
	db, err = sqlx.Connect("postgres", endpoint)
	return
}

func OpenRedisPool(redisAddress string) (err error) {
	redisPool = &redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000, // max number of connections
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", redisAddress)
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}
	return
}

func OpenRedisStore() (*redistore.RediStore, error) {
	rs, err := redistore.NewRediStoreWithPool(redisPool, []byte("secret-key"))
	if err != nil {
		panic(err)
	}

	return rs, err
}
