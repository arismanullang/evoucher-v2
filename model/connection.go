package model

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var db *sqlx.DB

//ConnectDB init connection to database server
func ConnectDB(endpoint string) (err error) {
	db, err = sqlx.Connect("postgres", endpoint)
	return
}
