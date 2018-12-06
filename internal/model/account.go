package model

import (
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type (
	//AccountRes Response Account
	AccountRes struct {
		Id   string `db:"id" json:"id"`
		Name string `db:"name" json:"name"`
	}
	//Account Account object
	Account struct {
		Id          string         `db:"id" json:"id"`
		Name        string         `db:"name" json:"name"`
		Billing     sql.NullString `db:"billing" json:"billing"`
		Alias       string         `db:"alias" json:"alias"`
		Email       string         `db:"email" json:"email"`
		Address     string         `db:"address" json:"address"`
		City        string         `db:"city" json:"city"`
		Province    string         `db:"province" json:"province"`
		Building    string         `db:"building" json:"building"`
		ZipCode     string         `db:"zip_code" json:"zip_code"`
		CreatedAt   time.Time      `db:"created_at" json:"created_at"`
		CreatedBy   string         `db:"created_by" json:"created_by"`
		UpdatedAt   sql.NullString `db:"updated_at" json:"updated_at"`
		UpdatedBy   sql.NullString `db:"updated_by" json:"updated_by"`
		Status      string         `db:"status" json:"status"`
		SenderEmail string         `db:"sender_mail" json:"sender_mail"`
	}
	//AccountConfig Config account
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
				, created_at
				, status
			)
			VALUES (?, ?, ?, ?, ?, ?, ?)
			RETURNING
				id
	`

	var res []string
	if err := tx.Select(&res, tx.Rebind(q), a.Name, a.Billing, a.Alias, a.Email, user, time.Now(), StatusCreated); err != nil {
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

func UpdateAccount(account Account, user string) error {
	tx, err := db.Beginx()
	if err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}
	defer tx.Rollback()

	accountDetail, err := GetAccountDetailByAccountId(account.Id)
	if err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}

	reflectParam := reflect.ValueOf(&account)
	dataParam := reflect.Indirect(reflectParam)

	reflectDb := reflect.ValueOf(&accountDetail).Elem()

	updates := getUpdate(dataParam, reflectDb)

	q := `
		UPDATE accounts
		SET
			updated_by = ?
			, updated_at = ?
		WHERE
			id = ?
	`

	_, err = tx.Exec(tx.Rebind(q), user, time.Now(), account.Id)
	if err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}

	for k, v := range updates {
		var value = v.String()
		if strings.Contains(value, "<") {
			tempString := strings.Replace(value, "<", "", -1)
			tempString = strings.Replace(tempString, ">", "", -1)
			tempStringArr := strings.Split(tempString, " ")
			if tempStringArr[0] == "int" {
				value = strconv.FormatInt(v.Int(), 64)
			} else if tempStringArr[0] == "float64" {
				value = strconv.FormatFloat(v.Float(), 'f', -1, 64)
			}
		}

		keys := strings.Split(k, ";")
		q = `
			UPDATE accounts
			SET
				`
		q += keys[1] + ` = '` + value + `'`
		q += `
			WHERE
				id = ?
				AND status = ?;
		`
		_, err = tx.Exec(tx.Rebind(q), account.Id, StatusCreated)
		if err != nil {
			fmt.Println(err.Error())
			return ErrServerInternal
		}
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
		SELECT
			id
			, name
			, billing
			, alias
			, email
			, address
			, city
			, province
			, building
			, zip_code
			, created_at
			, created_by
			, updated_at
			, updated_by
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
			a.id
			, a.name
			, a.billing
			, a.alias
			, a.created_at
			, a.address
			, a.city
			, a.province
			, a.building
			, a.zip_code
			, a.sender_mail
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
		return Account{}, ErrServerInternal
	}
	if len(resd) == 0 {
		return Account{}, ErrResourceNotFound
	}
	return resd[0], nil
}

func GetAccountDetailByAccountId(accountId string) (Account, error) {
	vc, err := db.Beginx()
	if err != nil {
		fmt.Println(err)
		return Account{}, ErrServerInternal
	}
	defer vc.Rollback()

	q := `
		SELECT
			a.id
			, a.name
			, a.billing
			, a.alias
			, a.email
			, a.created_at
			, a.address
			, a.city
			, a.province
			, a.building
			, a.zip_code
		FROM
			accounts as a
		WHERE
			a.id = ?
			AND a.status = ?
	`
	var resd []Account
	if err := db.Select(&resd, db.Rebind(q), accountId, StatusCreated); err != nil {
		fmt.Println(err)
		return Account{}, ErrServerInternal
	}
	if len(resd) == 0 {
		return Account{}, ErrResourceNotFound
	}
	return resd[0], nil
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
