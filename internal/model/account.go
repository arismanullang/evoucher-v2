package model

import (
	"database/sql"
	"fmt"
)

type (
	AccountRes struct {
		Id          string `db:"id"`
		AccountName string `db:"account_name"`
	}
	Account struct {
		Id          string         `db:"id" json:"id"`
		AccountName string         `db:"account_name" json:"account_name"`
		Billing     sql.NullString `db:"billing" json:"billing"`
		Alias       string         `db:"alias" json:"alias"`
		CreatedAt   string         `db:"created_at" json:"created_at"`
	}
)

func AddAccount(a Account, user string) error {
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
				, alias
				, created_by
				, status
			)
			VALUES (?, ?, ?, ?, ?)
		`
		if _, err := tx.Exec(tx.Rebind(q), a.AccountName, a.Billing, a.Alias, user, StatusCreated); err != nil {
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

func GetAccountDetailByUser(userID string) ([]Account, error) {
	vc, err := db.Beginx()
	if err != nil {
		fmt.Println(err)
		return []Account{}, ErrServerInternal
	}
	defer vc.Rollback()

	q := `
		SELECT
			a.id
			, a.account_name
			, a.billing
			, a.alias
			, a.created_at
		FROM
			accounts as a
		JOIN
			user_accounts as ua
		ON
			a.id = ua.account_id
		WHERE
			ua.user_id = ?
			AND a.status = ?
	`
	var resd []Account
	if err := db.Select(&resd, db.Rebind(q), userID, StatusCreated); err != nil {
		fmt.Println(err)
		return []Account{}, ErrServerInternal
	}
	if len(resd) == 0 {
		return []Account{}, ErrResourceNotFound
	}
	return resd, nil
}

func GetAccountDetailByAccountId(accountId string) ([]Account, error) {
	vc, err := db.Beginx()
	if err != nil {
		fmt.Println(err)
		return []Account{}, ErrServerInternal
	}
	defer vc.Rollback()

	q := `
		SELECT
			a.id
			, a.account_name
			, a.billing
			, a.alias
			, a.created_at
		FROM
			accounts as a
		JOIN
			user_accounts as ua
		ON
			a.id = ua.account_id
		WHERE
			a.id = ?
			AND a.status = ?
	`
	var resd []Account
	if err := db.Select(&resd, db.Rebind(q), accountId, StatusCreated); err != nil {
		fmt.Println(err)
		return []Account{}, ErrServerInternal
	}
	if len(resd) == 0 {
		return []Account{}, ErrResourceNotFound
	}
	return resd, nil
}

func GetAccountsByUser(userID string) ([]string, error) {
	vc, err := db.Beginx()
	if err != nil {
		fmt.Println(err)
		return []string{}, ErrServerInternal
	}
	defer vc.Rollback()

	q := `
		SELECT
			a.id
		FROM
			accounts as a
		JOIN
			user_accounts as ua
		ON
			a.id = ua.account_id
		WHERE
			ua.user_id = ?
			AND a.status = ?
	`
	var resd []string
	if err := db.Select(&resd, db.Rebind(q), userID, StatusCreated); err != nil {
		fmt.Println(err)
		return []string{}, ErrServerInternal
	}
	if len(resd) == 0 {
		return []string{}, ErrResourceNotFound
	}
	return resd, nil
}
