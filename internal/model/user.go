package model

import (
	"fmt"
)

type (
	AccountDetail struct {
		CompanyID   string `db:"company_id"`
		UserID      string `db:"user_id"`
		AccountRole string `db:"account_role"`
		AssignBy    string `db:"assign_by"`
		CreatedBy   string `db:"created_by"`
	}
	Account struct {
		ID     string `db:"id"`
		UserId string `db:"user_id"`
	}
	UserResponse struct {
		Status  string
		Message string
		Data    interface{}
	}
)

func AddAccount(acc AccountDetail) error {
	fmt.Println("Add")
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	q := `
		INSERT INTO accounts(
			company_id
			, user_id
			, account_role
			, assign_by
			, created_by
			, status
		)
		VALUES (?, ?, ?, ?, ?, ?)
		RETURNING
			id
	`
	var res []string
	if err := tx.Select(&res, tx.Rebind(q), acc.CompanyID, acc.UserID, acc.AccountRole, acc.AssignBy, acc.CreatedBy, StatusCreated); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func FindAllAccount() (UserResponse, error) {
	q := `
		SELECT
			id
			, user_id
		FROM
			accounts
		WHERE
			status = ?
	`

	var resv []Account
	if err := db.Select(&resv, db.Rebind(q), StatusCreated); err != nil {
		return UserResponse{Status: "Error", Message: q, Data: []Account{}}, err
	}
	if len(resv) < 1 {
		return UserResponse{Status: "404", Message: q, Data: []Account{}}, ErrResourceNotFound
	}

	res := UserResponse{
		Status:  "200",
		Message: "Ok",
		Data:    resv,
	}

	return res, nil
}

func FindAccountByRole(role string) (UserResponse, error) {
	q := `
		SELECT
			id
			, user_id
		FROM
			accounts
		WHERE
			account_role = ?
			AND status = ?
	`

	var resv []Account
	if err := db.Select(&resv, db.Rebind(q), role, StatusCreated); err != nil {
		return UserResponse{Status: "Error", Message: q, Data: []Account{}}, err
	}
	if len(resv) < 1 {
		return UserResponse{Status: "404", Message: q, Data: []Account{}}, ErrResourceNotFound
	}

	res := UserResponse{
		Status:  "200",
		Message: "Ok",
		Data:    resv,
	}

	return res, nil
}

func FindAccount(usr map[string]string) (UserResponse, error) {
	q := `
		SELECT
			id
			, user_id
		FROM
			accounts
		WHERE
			status = ?
	`
	for key, value := range usr {
		if key == "q" {
			q = q + `AND ` + key + ` ILIKE '%` + value + `%'`
		}
	}

	var resv []AccountDetail
	if err := db.Select(&resv, db.Rebind(q), StatusCreated); err != nil {
		return UserResponse{Status: "Error", Message: q, Data: []AccountDetail{}}, err
	}
	if len(resv) < 1 {
		return UserResponse{Status: "404", Message: q, Data: []AccountDetail{}}, ErrResourceNotFound
	}

	res := UserResponse{
		Status:  "200",
		Message: "Ok",
		Data:    resv,
	}

	return res, nil
}
