package model

import (
	"fmt"
	"time"
)

type (
	Cashout struct {
		Id            string               `db:"id" json:"id"`
		AccountId     string               `db:"account_id" json:"account_id"`
		CashoutCode   string               `db:"cashout_code" json:"cashout_code"`
		PartnerId     string               `db:"partner_id" json:"partner_id"`
		BankAccount   BankAccount          `db:"-" json:"bank_account"`
		TotalCashout  float64              `db:"total_cashout" json:"total_cashout"`
		PaymentMethod string               `db:"payment_method" json:"payment_method"`
		CreatedAt     time.Time            `db:"created_at" json:"created_at"`
		CreatedBy     string               `db:"created_by" json:"created_by"`
		Transactions  []CashoutTransaction `db:"-" json:"transactions"`
	}
	CashoutTransaction struct {
		TransactionId string `db:"transaction_id" json:"transaction_id"`
		VoucherId     string `db:"voucher_id" json:"voucher_id"`
		VoucherValue  string `db:"voucher_value" json:"voucher_value"`
		CreatedAt     string `db:"created_at" json:"created_at"`
	}
)

func InsertCashout(d Cashout) (string, error) {
	tx, err := db.Beginx()
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	q := `
		INSERT INTO cashouts(
			account_id
			, cashout_code
			, partner_id
			, bank_account
			, total_cashout
			, payment_method
			, created_by
			, status
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING
			id
	`
	var res []string
	if err := tx.Select(&res, tx.Rebind(q), d.AccountId, d.CashoutCode, d.PartnerId, d.BankAccount.Id, d.TotalCashout, d.PaymentMethod, d.CreatedBy, StatusCreated); err != nil {
		return "", err
	}
	d.Id = res[0]

	for _, v := range d.Transactions {
		q := `
			INSERT INTO cashout_details(
				cashout_id
				, transaction_id
				, voucher_id
				, created_by
				, status
			)
			VALUES (?, ?, ?, ?, ?)
		`

		_, err := tx.Exec(tx.Rebind(q), d.Id, v.TransactionId, v.VoucherId, d.CreatedBy, StatusCreated)
		if err != nil {
			return "", err
		}
	}

	if err := tx.Commit(); err != nil {
		return "", err
	}
	return d.Id, nil
}

func UpdateCashoutTransactions(transactionId []string, user string) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	q := `
		SELECT v.id
		FROM vouchers as v
		JOIN transaction_details as td
		ON
			v.id = td.voucher_id
		JOIN transactions as t
		ON
			t.id = td.transaction_id
		WHERE
			t.status = ?
			AND (
	`

	for i, v := range transactionId {
		if i != 0 {
			q += ` OR `
		}
		q += `t.id = '` + v + `'`
	}
	q += `)`
	//fmt.Println(q)
	var resv []string
	if err := db.Select(&resv, db.Rebind(q), StatusCreated); err != nil {
		return err
	}
	if len(resv) < 1 {
		return ErrResourceNotFound
	}

	q = `
		UPDATE vouchers
		SET
			state = ?
			, updated_by = ?
			, updated_at = ?
		WHERE
			status = ?
			AND (
	`
	for i, v := range resv {
		if i != 0 {
			q += ` OR `
		}
		q += `id = '` + v + `'`
	}
	q += `)`
	//fmt.Println(q)
	_, err = tx.Exec(tx.Rebind(q), VoucherStatePaid, user, time.Now(), StatusCreated)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func PrintCashout(accountId string, cashoutCode string) (Cashout, error) {
	q := `
		SELECT
			id, cashout_code, partner_id, total_cashout, created_at
		FROM cashouts
		WHERE
			status = ?
			AND account_id = ?
			AND id = ?`
	var res []Cashout
	if err := db.Select(&res, db.Rebind(q), StatusCreated, accountId, cashoutCode); err != nil {
		fmt.Println("cashout : " + err.Error())
		return Cashout{}, err
	}
	if len(res) < 1 {
		fmt.Println("cashout : not found")
		return Cashout{}, ErrResourceNotFound
	}

	q = `
		SELECT DISTINCT
			t.transaction_code as transaction_id, v.voucher_code as voucher_id, v.voucher_value, t.created_at
		FROM cashout_details as cd
		JOIN
			transactions as t
		ON
			cd.transaction_id = t.id
		JOIN
			vouchers as v
		ON
			cd.voucher_id = v.id
		WHERE
			cd.status = ?
			AND cd.cashout_id = ?`
	var rest []CashoutTransaction
	if err := db.Select(&rest, db.Rebind(q), StatusCreated, res[0].Id); err != nil {
		fmt.Println("cashout detail : " + err.Error())
		return Cashout{}, err
	}
	if len(rest) < 1 {
		fmt.Println("cashout detail : not found")
		return Cashout{}, ErrResourceNotFound
	}

	bank, err := FindBankAccountByPartner(accountId, res[0].PartnerId)
	if err != nil {
		fmt.Println("cashout detail : " + err.Error())
		return Cashout{}, err
	}
	if len(rest) < 1 {
		fmt.Println("cashout detail : not found")
		return Cashout{}, ErrResourceNotFound
	}

	res[0].Transactions = rest
	res[0].BankAccount = bank
	return res[0], nil
}
