package model

import (
	"fmt"
)

type (
	AccountRes struct {
		Id          string `db:"id"`
		AccountName string `db:"account_name"`
	}
	Account struct {
		Id          string `db:"id"`
		AccountName string `db:"account_name"`
		Billing     string `db:"billing"`
		CreatedBy   string `db:"created_by"`
	}
)

func AddAccount(a Account) error {
	fmt.Println("Add")
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	accountName, err := checkAccountName(a.AccountName)

	if accountName == 0 {
		q := `
			INSERT INTO accounts(
				account_name
				, billing
				, created_by
				, status
			)
			VALUES (?, ?, ?, ?)
		`
		if _, err := tx.Exec(tx.Rebind(q), a.AccountName, a.Billing, a.CreatedBy, StatusCreated); err != nil {
			return err
		}

		if err := tx.Commit(); err != nil {
			return err
		}
		return nil
	} else {
		return ErrDuplicateEntry
	}
}

func checkAccountName(name string) (int, error) {
	q := `
		SELECT id FROM users
		WHERE
			account_name = ?
			AND status = ?
	`
	var res []string
	if err := db.Select(&res, db.Rebind(q), name, StatusCreated); err != nil {
		return 0, err
	}

	return len(res), nil
}

func FindAllAccount() (Response, error) {
	fmt.Println("Select All Account")
	q := `
		SELECT id, account_name
		FROM accounts
		WHERE status = ?
	`

	var resv []AccountRes
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
