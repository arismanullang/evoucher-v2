package model

import (
	"fmt"
	"time"

	"github.com/gilkor/evoucher/internal/util"
)

type (
	//User :
	User struct {
		ID        string     `db:"id" json:"id"`
		AuthToken string     `db:"auth_token" json:"auth_token"`
		CompanyID string     `db:"company_id" json:"company_id"`
		CreatedAt *time.Time `db:"created_at" json:"created_at,omitempty"`
		CreatedBy string     `db:"created_by" json:"created_by,omitempty"`
		UpdatedAt *time.Time `db:"updated_at" json:"updated_at,omitempty"`
		UpdatedBy string     `db:"updated_by" json:"updated_by,omitempty"`
		Status    string     `db:"status" json:"status,omitempty"`
	}
)

func getUsers(key, value string, f *util.Filter) ([]User, bool, error) {
	q := f.GetQueryByDefaultStruct(Company{})
	q += `
			FROM
				users
			WHERE 
				status = ?			
			AND ` + key + ` = ?`

	q += f.GetQuerySort()
	q += f.GetQueryLimit()
	fmt.Println(q)
	var resd []User
	err := db.Select(&resd, db.Rebind(q), StatusCreated, value)
	if err != nil {
		return []User{}, false, err
	}

	next := false
	if len(resd) > f.Count {
		next = true
	}
	if len(resd) < f.Count {
		f.Count = len(resd)
	}

	return resd, next, nil
}

//Insert : single row inset into table
func (u User) Insert() error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	q := `INSERT INTO 
				users 
				( 
					 auth_token
					, company_id
					, created_by 
					, updated_by
					, status
				)
			VALUES 
				( ?, ?, ?, ?, ?)
			RETURNING
				id
				, name
				, description 
				, alias 
				, client_key
				, client_secret 
				, created_at
				, created_by
				, updated_at
				, updated_by
				, status
	`
	var res []Company
	err = tx.Select(&res, tx.Rebind(q), u.AuthToken, u.CompanyID, u.CreatedBy, u.UpdatedBy, StatusCreated)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
