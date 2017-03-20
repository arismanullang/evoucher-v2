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

	if accountName == "" {
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

func checkAccountName(name string) (string, error) {
	q := `
		SELECT id FROM users
		WHERE
			account_name = ?
			AND status = ?
	`
	var res []string
	if err := db.Select(&res, db.Rebind(q), name, StatusCreated); err != nil {
		return "", err
	}
	if len(res) == 0 {
		return "", nil
	}
	return res[0], nil
}

func FindAllAccounts() ([]AccountRes, error) {
	fmt.Println("Select All Account")
	q := `
		SELECT id, account_name
		FROM accounts
		WHERE status = ?
	`

	var resv []AccountRes
	if err := db.Select(&resv, db.Rebind(q), StatusCreated); err != nil {
		return []AccountRes{}, err
	}
	if len(resv) < 1 {
		return []AccountRes{}, ErrResourceNotFound
	}

	return resv, nil
}

func GetAccountByUser(userID string) ([]AccountRes, error) {
	vc, err := db.Beginx()
	if err != nil {
		fmt.Println(err)
		return []AccountRes{}, ErrServerInternal
	}
	defer vc.Rollback()

	q := `
		SELECT
			a.id
			, a.account_name
		FROM
			accounts as a
		JOIN
			user_accounts as ua
		WHERE
			ua.user_id = ?
			AND a.status = ?
	`
	var resd []AccountRes
	if err := db.Select(&resd, db.Rebind(q), userID, StatusCreated); err != nil {
		fmt.Println(err)
		return []AccountRes{}, ErrServerInternal
	}
	if len(resd) == 0 {
		return []AccountRes{}, ErrResourceNotFound
	}
	return resd, nil
}