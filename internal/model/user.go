package model

import (
	"fmt"
	"log"
)

type (
	User struct {
		AccountId string   `db:"account_id"`
		Username  string   `db:"username"`
		Password  string   `db:"password"`
		Email     string   `db:"email"`
		Phone     string   `db:"phone"`
		RoleId    []string `db:"-"`
		CreatedBy string   `db:"created_by"`
	}
	UserRes struct {
		Id       string `db:"id"`
		Username string `db:"username"`
	}
	Role struct {
		Id         string `db:"id"`
		RoleDetail string `db:"role_detail"`
	}
)

func AddUser(u User) error {
	fmt.Println("Add")
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	username, err := CheckUsername(u.Username)

	if username != "" {
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
		if err := tx.Select(&res, tx.Rebind(q), u.Username, u.Password, u.Email, u.Phone, u.CreatedBy, StatusCreated); err != nil {
			return err
		}

		for _, v := range u.RoleId {
			q := `
				INSERT INTO user_roles(
					user_id
					, role_id
					, created_by
					, status
				)
				VALUES (?, ?, ?, ?)
			`

			_, err := tx.Exec(tx.Rebind(q), res[0], v, u.CreatedBy, StatusCreated)
			if err != nil {
				return err
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

		_, err := tx.Exec(tx.Rebind(q2), res[0], u.AccountId, u.CreatedBy, StatusCreated)
		if err != nil {
			return err
		}

		if err := tx.Commit(); err != nil {
			return err
		}
		return nil
	}

	return ErrDuplicateEntry
}

func CheckUsername(username string) (string, error) {
	q := `
		SELECT id FROM users
		WHERE
			username = ?
			AND status = ?
	`
	var res []string
	if err := db.Select(&res, db.Rebind(q), username, StatusCreated); err != nil {
		return "", err
	}

	return res[0], nil
}

func FindAllUser(accountId string) (Response, error) {
	fmt.Println("Select User " + accountId)
	q := `
		SELECT DISTINCT u.id, u.username FROM users as u
		JOIN user_accounts as ua ON u.id = ua.user_id
		JOIN user_roles as ur ON u.id = ur.user_id
		WHERE ua.account_id = ?
		AND u.status = ?
	`

	var resv []UserRes
	if err := db.Select(&resv, db.Rebind(q), accountId, StatusCreated); err != nil {
		return Response{Status: "Error", Message: q, Data: []UserRes{}}, err
	}
	if len(resv) < 1 {
		return Response{Status: "404", Message: q, Data: []UserRes{}}, ErrResourceNotFound
	}

	res := Response{
		Status:  "200",
		Message: "Ok",
		Data:    resv,
	}

	return res, nil
}

func FindUserByRole(role, accountId string) (Response, error) {
	q := `
		SELECT u.id, u.username FROM users AS u
		JOIN user_accounts AS ua ON u.id = ua.user_id
		JOIN user_roles AS ur ON u.id = ur.user_id
		WHERE ua.account_id = ?
		AND ur.role_id = ?
		AND u.status = ?
	`

	var resv []UserRes
	if err := db.Select(&resv, db.Rebind(q), accountId, role, StatusCreated); err != nil {
		return Response{Status: "Error", Message: q, Data: []UserRes{}}, err
	}
	if len(resv) < 1 {
		return Response{Status: "404", Message: q, Data: []UserRes{}}, ErrResourceNotFound
	}

	res := Response{
		Status:  "200",
		Message: "Ok",
		Data:    resv,
	}

	return res, nil
}

func FindUser(usr map[string]string) (Response, error) {
	q := `
		SELECT u.id, u.username FROM users AS u
		JOIN user_accounts AS ua ON u.id = ua.user_id
		JOIN user_roles AS ur ON u.id = ur.user_id
		WHERE
			status = ?
	`

	for key, value := range usr {
		if key == "q" {
			q += `AND (u.username ILIKE '%` + value + `%')`
		} else {
			q += ` AND ` + key + ` LIKE '%` + value + `%'`
		}
	}

	var resv []User
	if err := db.Select(&resv, db.Rebind(q), StatusCreated); err != nil {
		return Response{Status: "Error", Message: q, Data: []User{}}, err
	}
	if len(resv) < 1 {
		return Response{Status: "404", Message: q, Data: []User{}}, ErrResourceNotFound
	}

	res := Response{
		Status:  "200",
		Message: "Ok",
		Data:    resv,
	}

	return res, nil
}

func Login(username, password, accountId string) (string, error) {
	fmt.Println("Login")
	q := `
		SELECT
			u.id
		FROM
			users as u
		JOIN
			user_accounts as ua
		ON
			u.id = ua.user_id
		WHERE
			u.username = ?
			AND u.password = ?
			AND ua.account_id = ?
			AND u.status = ?
	`
	var res []string
	if err := db.Select(&res, db.Rebind(q), username, password, accountId, StatusCreated); err != nil {
		return "", err
	}
	if len(res) == 0 {
		return "", ErrResourceNotFound
	}
	return res[0], nil
}

func GetAccountByUser(userID string) string {
	vc, err := db.Beginx()
	if err != nil {
		return ""
	}
	defer vc.Rollback()

	q := `
		SELECT
			account_id
		FROM
			user_accounts
		WHERE
			user_id = ?
			AND status = ?
	`
	var resd []string
	if err := db.Select(&resd, db.Rebind(q), userID, StatusCreated); err != nil {
		log.Panic(err)
		return ""
	}
	if len(resd) == 0 {
		return ""
	}
	return resd[0]
}

func FindAllRole() (Response, error) {
	fmt.Println("Select All Role")
	q := `
		SELECT id, role_detail
		FROM roles
		WHERE status = ?
	`

	var resv []Role
	if err := db.Select(&resv, db.Rebind(q), StatusCreated); err != nil {
		return Response{Status: "Error", Message: q, Data: []AccountRes{}}, err
	}
	if len(resv) < 1 {
		return Response{Status: "404", Message: q, Data: []AccountRes{}}, ErrResourceNotFound
	}

	res := Response{
		Status:  "200",
		Message: "Ok",
		Data:    resv,
	}

	return res, nil
}
