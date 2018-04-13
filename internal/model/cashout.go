package model

import (
	"fmt"
	"time"
)

type (
	Cashout struct {
		Id                   string               `db:"id" json:"id"`
		AccountId            string               `db:"account_id" json:"account_id"`
		CashoutCode          string               `db:"cashout_code" json:"cashout_code"`
		PartnerId            string               `db:"partner_id" json:"partner_id"`
		BankAccount          string               `db:"bank_account" json:"bank_account"`
		BankAccountCompany   string               `db:"bank_account_company" json:"bank_account_company"`
		BankAccountNumber    string               `db:"bank_account_number" json:"bank_account_number"`
		BankAccountRefNumber string               `db:"bank_account_ref_number" json:"bank_account_ref_number"`
		TotalCashout         float64              `db:"total_cashout" json:"total_cashout"`
		PaymentMethod        string               `db:"payment_method" json:"payment_method"`
		CreatedAt            time.Time            `db:"created_at" json:"created_at"`
		CreatedBy            string               `db:"created_by" json:"created_by"`
		Transactions         []CashoutTransaction `db:"-" json:"transactions"`
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
			, bank_account_company
			, bank_account_number
			, bank_account_ref_number
			, total_cashout
			, payment_method
			, created_by
			, created_at
			, status
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING
			id
	`
	var res []string
	if err := tx.Select(&res, tx.Rebind(q), d.AccountId, d.CashoutCode, d.PartnerId, d.BankAccount, d.BankAccountCompany, d.BankAccountNumber, d.BankAccountRefNumber, d.TotalCashout, d.PaymentMethod, d.CreatedBy, time.Now(), StatusCreated); err != nil {
		return "", err
	}
	d.Id = res[0]

	logs := []Log{}
	tempLog := Log{
		TableName:   "cashouts",
		TableNameId: ValueChangeLogNone,
		ColumnName:  ColumnChangeLogInsert,
		Action:      ActionChangeLogInsert,
		Old:         ValueChangeLogNone,
		New:         d.Id,
		CreatedBy:   d.CreatedBy,
	}
	logs = append(logs, tempLog)

	for _, v := range d.Transactions {
		q := `
			INSERT INTO cashout_details(
				cashout_id
				, transaction_id
				, voucher_id
				, created_by
				, created_at
				, status
			)
			VALUES (?, ?, ?, ?, ?, ?)
			RETURNING
				id
		`
		var res1 []string
		err := tx.Select(&res1, tx.Rebind(q), d.Id, v.TransactionId, v.VoucherId, d.CreatedBy, time.Now(), StatusCreated)
		if err != nil {
			return "", err
		}

		tempLog = Log{
			TableName:   "cashout_details",
			TableNameId: d.Id,
			ColumnName:  ColumnChangeLogInsert,
			Action:      ActionChangeLogInsert,
			Old:         ValueChangeLogNone,
			New:         res1[0],
			CreatedBy:   d.CreatedBy,
		}
		logs = append(logs, tempLog)
	}

	if err := tx.Commit(); err != nil {
		return "", err
	}

	err = addLogs(logs)
	if err != nil {
		fmt.Println(err.Error())
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

	logs := []Log{}
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

		tempLog := Log{
			TableName:   "vouchers",
			TableNameId: v,
			ColumnName:  "updated_by",
			Action:      ActionChangeLogUpdate,
			Old:         ValueChangeLogNone,
			New:         user,
			CreatedBy:   user,
		}
		logs = append(logs, tempLog)

		tempLog = Log{
			TableName:   "vouchers",
			TableNameId: v,
			ColumnName:  "updated_at",
			Action:      ActionChangeLogUpdate,
			Old:         ValueChangeLogNone,
			New:         time.Now().String(),
			CreatedBy:   user,
		}
		logs = append(logs, tempLog)

		tempLog = Log{
			TableName:   "vouchers",
			TableNameId: v,
			ColumnName:  "state",
			Action:      ActionChangeLogUpdate,
			Old:         ValueChangeLogNone,
			New:         VoucherStatePaid,
			CreatedBy:   user,
		}
		logs = append(logs, tempLog)
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

	err = addLogs(logs)
	if err != nil {
		fmt.Println(err.Error())
	}

	return nil
}

func PrintCashout(accountId string, cashoutCode string) (Cashout, error) {
	q := `
		SELECT
			id, cashout_code, bank_account, bank_account_company, bank_account_number, bank_account_ref_number, partner_id, total_cashout, created_at
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

	res[0].Transactions = rest

	log := Log{
		TableName:   "cashouts",
		TableNameId: cashoutCode,
		ColumnName:  ColumnChangeLogSelect,
		Action:      ActionChangeLogSelect,
		Old:         ValueChangeLogNone,
		New:         ValueChangeLogNone,
		CreatedBy:   accountId,
	}

	err := addLog(log)
	if err != nil {
		fmt.Println(err.Error())
	}

	return res[0], nil
}

func FindAllReimburse(accountId, user string) ([]Cashout, error) {
	q := `
		SELECT
			c.id
			, p.name as partner_id
			, c.cashout_code
			, c.bank_account
			, c.total_cashout
			, c.bank_account_number
			, c.bank_account_ref_number
			, c.bank_account_company
			, c.created_at
		FROM
			cashouts AS c
		JOIN
			partners AS p
		ON
			c.partner_id = p.id
		WHERE
			c.status = ?
			AND c.account_id = ?
		ORDER BY
		 	c.created_at desc
		`
	var res []Cashout
	if err := db.Select(&res, db.Rebind(q), StatusCreated, accountId); err != nil {
		return []Cashout{}, err
	}
	if len(res) < 1 {
		return []Cashout{}, ErrResourceNotFound
	}

	log := Log{
		TableName:   "cashouts",
		TableNameId: accountId,
		ColumnName:  ColumnChangeLogSelect,
		Action:      ActionChangeLogSelect,
		Old:         ValueChangeLogNone,
		New:         ValueChangeLogNone,
		CreatedBy:   user,
	}

	err := addLog(log)
	if err != nil {
		fmt.Println(err.Error())
	}

	return res, nil
}
