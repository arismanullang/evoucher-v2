package model

import (
	"fmt"
	"time"
)

type (
	Cashout struct {
		Id            string        `db:"id" json:"id"`
		AccountId     string        `db:"account_id" json:"account_id"`
		CashoutCode   string        `db:"cashout_code" json:"cashout_code"`
		PartnerId     string        `db:"partner_id" json:"partner_id"`
		BankAccount   string        `db:"bank_account" json:"bank_account"`
		TotalCashout  float64       `db:"total_cashout" json:"total_cashout"`
		PaymentMethod string        `db:"payment_method" json:"payment_method"`
		CreatedAt     time.Time     `db:"created_at" json:"created_at"`
		CreatedBy     string        `db:"created_by" json:"created_by"`
		Transactions  []Transaction `db:"-" json:"transactions"`
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
	if err := tx.Select(&res, tx.Rebind(q), d.AccountId, d.CashoutCode, d.PartnerId, d.BankAccount, d.TotalCashout, d.PaymentMethod, d.CreatedBy, StatusCreated); err != nil {
		return "", err
	}
	d.Id = res[0]

	for _, v := range d.Transactions {
		q := `
			INSERT INTO cashout_details(
				cashout_id
				, transaction_id
				, created_by
				, status
			)
			VALUES (?, ?, ?, ?)
		`

		_, err := tx.Exec(tx.Rebind(q), d.Id, v.Id, d.CreatedBy, StatusCreated)
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
			id, cashout_code, partner_id, bank_account, total_cashout, created_at
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
		SELECT
			cd.transaction_id as id, t.transaction_code
		FROM cashout_details as cd
		JOIN
			transactions as t
		ON
			cd.transaction_id = t.id
		WHERE
			cd.status = ?
			AND cd.cashout_id = ?`
	var rest []Transaction
	if err := db.Select(&rest, db.Rebind(q), StatusCreated, res[0].Id); err != nil {
		fmt.Println("cashout detail : " + err.Error())
		return Cashout{}, err
	}
	if len(rest) < 1 {
		fmt.Println("cashout detail : not found")
		return Cashout{}, ErrResourceNotFound
	}

	transactions := []Transaction{}
	for _, v := range rest {
		q := `
		SELECT
			v.id
			, v.voucher_code
			, v.reference_no
			, v.holder
			, v.holder_phone
			, v.holder_email
			, v.holder_description
			, v.program_id
			, v.valid_at
			, v.expired_at
			, v.voucher_value
			, v.state
			, v.created_by
			, v.created_at
			, v.updated_by
			, v.updated_at
			, v.deleted_by
			, v.deleted_at
			, v.status
		FROM vouchers as v
		JOIN transaction_details as dt
		ON
			v.id = dt.voucher_id
		WHERE
			v.status = ?
			AND dt.transaction_id = ?
	`
		//fmt.Println(q)
		var resv []Voucher
		if err := db.Select(&resv, db.Rebind(q), StatusCreated, v.Id); err != nil {
			fmt.Println("voucher : " + err.Error())
			return Cashout{}, err
		}
		if len(resv) < 1 {
			fmt.Println("voucher : not found")
			return Cashout{}, ErrResourceNotFound
		}
		transaction := Transaction{
			Id:              v.Id,
			TransactionCode: v.TransactionCode,
			Vouchers:        resv,
		}
		transactions = append(transactions, transaction)
	}

	res[0].Transactions = transactions
	return res[0], nil
}
