package model

import (
	"database/sql"
	"fmt"
	"time"
)

type (
	AccountRes struct {
		Id   string `db:"id" json:"id"`
		Name string `db:"name" json:"name"`
	}
	Account struct {
		Id        string         `db:"id" json:"id"`
		Name      string         `db:"name" json:"name"`
		Billing   sql.NullString `db:"billing" json:"billing"`
		Alias     string         `db:"alias" json:"alias"`
		Email     string         `db:"email" json:"email"`
		CreatedAt string         `db:"created_at" json:"created_at"`
		CreatedBy string         `db:"created_by" json:"created_by"`
		UpdatedAt sql.NullString `db:"updated_at" json:"updated_at"`
		UpdatedBy sql.NullString `db:"updated_by" json:"updated_by"`
		Status    string         `db:"status" json:"status"`
	}
	AccountConfig struct {
		AccountId    string `db:"account_id"`
		ConfigDetail string `db:"config_detail"`
		ConfigValue  string `db:"config_value"`
	}
)

func AddAccount(a Account, user string) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	name, err := checkName(a.Name)

	if name != "" {
		return ErrDuplicateEntry
	}

	q := `
			INSERT INTO accounts(
				name
				, billing
				, alias
				, email
				, created_by
				, status
			)
			VALUES (?, ?, ?, ?, ?, ?)
			RETURNING
				id
	`

	var res []string
	if err := tx.Select(&res, tx.Rebind(q), a.Name, a.Billing, a.Alias, a.Email, user, StatusCreated); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}

	err = AddAdmin(res[0], user)
	if err != nil {
		return err
	}
	return nil
}

func checkName(name string) (string, error) {
	q := `
		SELECT id FROM users
		WHERE
			name = ?
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

func UpdateAccount(account Account, userId string) error {
	tx, err := db.Beginx()
	if err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}
	defer tx.Rollback()

	q := `
		UPDATE accounts
		SET
			name = ?
			, alias = ?
			, email = ?
			, updated_by = ?
			, updated_at = ?
		WHERE
			id = ?
	`

	_, err = tx.Exec(tx.Rebind(q), account.Name, account.Alias, account.Email, userId, time.Now(), account.Id)
	if err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}

	if err := tx.Commit(); err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}
	return nil
}

func FindAllAccounts() ([]AccountRes, error) {
	q := `
		SELECT id, name
		FROM accounts
		WHERE status = ?
		AND NOT name = 'suadmin'
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

func FindAllAccountsDetail() ([]Account, error) {
	q := `
		SELECT *
		FROM accounts
		WHERE NOT name = 'suadmin'
	`

	var resv []Account
	if err := db.Select(&resv, db.Rebind(q)); err != nil {
		return []Account{}, err
	}
	if len(resv) < 1 {
		return []Account{}, ErrResourceNotFound
	}

	return resv, nil
}

func GetAccountDetailByUser(userID string) (Account, error) {
	vc, err := db.Beginx()
	if err != nil {
		fmt.Println(err)
		return Account{}, ErrServerInternal
	}
	defer vc.Rollback()

	q := `
		SELECT
			id
			, name
			, billing
			, email
			, alias
			, created_at
		FROM
			accounts
		WHERE
			id = ?
			AND status = ?
	`
	var resd []Account
	if err := db.Select(&resd, db.Rebind(q), userID, StatusCreated); err != nil {
		fmt.Println(err)
		return Account{}, ErrServerInternal
	}
	if len(resd) == 0 {
		return Account{}, ErrResourceNotFound
	}
	return resd[0], nil
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
			, a.name
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

func BlockAccount(accountId, userId string) error {
	tx, err := db.Beginx()
	if err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}
	defer tx.Rollback()

	q := `
		UPDATE accounts
		SET
			status = ?
			, updated_by = ?
			, updated_at = ?
		WHERE
			id = ?
	`

	_, err = tx.Exec(tx.Rebind(q), StatusDeleted, userId, time.Now(), accountId)
	if err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}

	if err := tx.Commit(); err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}
	return nil
}

func ActivateAccount(accountId, userId string) error {
	tx, err := db.Beginx()
	if err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}
	defer tx.Rollback()

	q := `
		UPDATE accounts
		SET
			status = ?
			, updated_by = ?
			, updated_at = ?
		WHERE
			id = ?
	`

	_, err = tx.Exec(tx.Rebind(q), StatusCreated, userId, time.Now(), accountId)
	if err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}

	if err := tx.Commit(); err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}
	return nil
}

func GetAccountConfig() ([]AccountConfig, error) {
	vc, err := db.Beginx()
	if err != nil {
		fmt.Println(err)
		return []AccountConfig{}, ErrServerInternal
	}
	defer vc.Rollback()

	q := `
		SELECT

			ac.account_id, c.value as config_detail, ac.value as config_value
		FROM
			account_configs as ac
		JOIN
			configs as c
		ON
			c.id = ac.config_id
		WHERE
			ac.status = ?
		ORDER BY
			ac.account_id
	`
	var resd []AccountConfig
	if err := db.Select(&resd, db.Rebind(q), StatusCreated); err != nil {
		fmt.Println(err)
		return []AccountConfig{}, ErrServerInternal
	}
	if len(resd) == 0 {
		return []AccountConfig{}, ErrResourceNotFound
	}
	return resd, nil
}
