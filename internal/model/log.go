package model

import (
//"database/sql"
//"fmt"
//"time"
)

type (
	Log struct {
		Id   string `db:"id" json:"id"`
		Name string `db:"name" json:"name"`
	}
)
