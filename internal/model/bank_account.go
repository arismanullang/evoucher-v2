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
	BankAccount struct {
		Id                string         `db:"id" json:"id"`
		CompanyName       string         `db:"company_name" json:"company_name"`
		CompanyPic        string         `db:"company_pic" json:"company_pic"`
		CompanyTelp       string         `db:"company_telp" json:"company_telp"`
		CompanyEmail      string         `db:"company_email" json:"company_email"`
		BankName          string         `db:"bank_name" json:"bank_name"`
		BankBranch        string         `db:"bank_branch" json:"bank_branch"`
		BankAccountNumber string         `db:"bank_account_number" json:"bank_account_number"`
		BankAccountHolder string         `db:"bank_account_holder" json:"bank_account_holder"`
		AccountId         string         `db:"account_id" json:"account_id"`
		CreatedAt         string         `db:"created_at" json:"created_at"`
		CreatedBy         string         `db:"created_by" json:"created_by"`
		UpdatedAt         sql.NullString `db:"updated_at" json:"updated_at"`
		UpdatedBy         sql.NullString `db:"updated_by" json:"updated_by"`
		Status            string         `db:"status" json:"status"`
	}
)

func AddBankAccount(a BankAccount, user User) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	name, err := checkBankAccount(a.BankAccountNumber)

	if name != "" {
		return ErrDuplicateEntry
	}

	q := `
			INSERT INTO bank_accounts(
				company_name
				, company_pic
				, company_telp
				, company_email
				, bank_name
				, bank_branch
				, bank_account_number
				, bank_account_holder
				, account_id
				, created_by
				, status
			)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			RETURNING
				id
	`

	var res []string
	if err := tx.Select(&res, tx.Rebind(q), a.CompanyName, a.CompanyPic, a.CompanyTelp, a.CompanyEmail, a.BankName, a.BankBranch, a.BankAccountNumber, a.BankAccountHolder, user.Account.Id, user.ID, StatusCreated); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func checkBankAccount(name string) (string, error) {
	q := `
		SELECT id FROM bank_accounts
		WHERE
			bank_account_number = ?
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

func UpdateBankAccount(account BankAccount, userId string) error {
	tx, err := db.Beginx()
	if err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}
	defer tx.Rollback()

	bankAccountDetail, err := FindBankAccountById(account.Id)
	if err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}

	reflectParam := reflect.ValueOf(&account)
	dataParam := reflect.Indirect(reflectParam)

	reflectDb := reflect.ValueOf(&bankAccountDetail).Elem()

	updates := getUpdate(dataParam, reflectDb)

	q := `
		UPDATE bank_accounts
		SET
			updated_by = ?
			, updated_at = ?
		WHERE
			id = ?
	`

	_, err = tx.Exec(tx.Rebind(q), userId, time.Now(), account.Id)
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
			UPDATE bank_accounts
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
			fmt.Println(q)
			return ErrServerInternal
		}
	}

	if err := tx.Commit(); err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}

	return nil
}

func FindAllBankAccounts(accountId string) ([]BankAccount, error) {
	q := `
		SELECT
			id
			, company_name
			, company_pic
			, company_telp
			, company_email
			, bank_name
			, bank_branch
			, bank_account_number
			, bank_account_holder
		FROM bank_accounts
		WHERE status = ?
		AND account_id = ?
	`

	var resv []BankAccount
	if err := db.Select(&resv, db.Rebind(q), StatusCreated, accountId); err != nil {
		return []BankAccount{}, err
	}
	if len(resv) < 1 {
		return []BankAccount{}, ErrResourceNotFound
	}

	return resv, nil
}

func FindBankAccount(accountId, number string) (BankAccount, error) {
	q := `
		SELECT
			id
			, company_name
			, company_pic
			, company_telp
			, company_email
			, bank_name
			, bank_branch
			, bank_account_number
			, bank_account_holder
		FROM bank_accounts
		WHERE status = ?
		AND account_id = ?
		AND bank_account_number = ?
	`

	var resv []BankAccount
	if err := db.Select(&resv, db.Rebind(q), StatusCreated, accountId, number); err != nil {
		return BankAccount{}, err
	}
	if len(resv) < 1 {
		return BankAccount{}, ErrResourceNotFound
	}

	return resv[0], nil
}

func FindBankAccountByPartner(accountId, partnerId string) (BankAccount, error) {
	q := `
		SELECT
			ba.id
			, ba.company_name
			, ba.company_pic
			, ba.company_telp
			, ba.company_email
			, ba.bank_name
			, ba.bank_branch
			, ba.bank_account_number
			, ba.bank_account_holder
		FROM bank_accounts as ba
		JOIN partners as p
		ON
			ba.id = p.bank_account_id
		WHERE
			ba.status = ?
			AND ba.account_id = ?
			AND p.id = ?
	`

	var resv []BankAccount
	if err := db.Select(&resv, db.Rebind(q), StatusCreated, accountId, partnerId); err != nil {
		return BankAccount{}, err
	}
	if len(resv) < 1 {
		return BankAccount{}, ErrResourceNotFound
	}

	return resv[0], nil
}

func FindBankAccountById(id string) (BankAccount, error) {
	q := `
		SELECT
			id
			, company_name
			, company_pic
			, company_telp
			, company_email
			, bank_name
			, bank_branch
			, bank_account_number
			, bank_account_holder
		FROM bank_accounts
		WHERE status = ?
		AND id = ?
	`

	var resv []BankAccount
	if err := db.Select(&resv, db.Rebind(q), StatusCreated, id); err != nil {
		return BankAccount{}, err
	}
	if len(resv) < 1 {
		return BankAccount{}, ErrResourceNotFound
	}

	return resv[0], nil
}
