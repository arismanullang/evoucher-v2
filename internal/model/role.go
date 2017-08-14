package model

import (
	"database/sql"
	"fmt"
	"time"
)

type (
	Role struct {
		Id       string   `db:"id" json:"id"`
		Detail   string   `db:"detail" json:"detail"`
		Featrues []string `db:"-" json:"features"`
	}
)

// Role -----------------------------------------------------------------------------------------------

func FindAllRole() ([]Role, error) {
	q := `
		SELECT id, detail
		FROM roles
		WHERE status = ?
		AND NOT detail = 'suadmin'
	`

	var resv []Role
	if err := db.Select(&resv, db.Rebind(q), StatusCreated); err != nil {
		fmt.Println(err)
		return []Role{}, ErrServerInternal
	}
	if len(resv) < 1 {
		return []Role{}, ErrResourceNotFound
	}

	return resv, nil
}

func AddRole(r Role, user string) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	q := `
			INSERT INTO users(
				username
				, password
				, email
				, phone
				, created_by
				, status
			)
			VALUES (?, ?, ?, ?, ?, ?)
			RETURNING
				id
		`
	var res []string
	if err := tx.Select(&res, tx.Rebind(q), u.Username, u.Password, u.Email, u.Phone, user, StatusCreated); err != nil {
		fmt.Println(err)
		return ErrServerInternal
	}

	for _, v := range u.Role {
		q := `
				INSERT INTO user_roles(
					user_id
					, role_id
					, created_by
					, status
				)
				VALUES (?, ?, ?, ?)
			`

		_, err := tx.Exec(tx.Rebind(q), res[0], v, user, StatusCreated)
		if err != nil {
			fmt.Println(err)
			return ErrServerInternal
		}
	}

	q2 := `
			INSERT INTO user_accounts(
				user_id
				, account_id
				, created_by
				, status
			)
			VALUES (?, ?, ?, ?)
		`

	_, err := tx.Exec(tx.Rebind(q2), res[0], accountId, user, StatusCreated)
	if err != nil {
		fmt.Println(err)
		return ErrServerInternal
	}

	if err := tx.Commit(); err != nil {
		fmt.Println(err)
		return ErrServerInternal
	}
	return nil

}
