package model

import (
	"database/sql"
	"fmt"
	"time"
)

type (
	Transaction struct {
		Id              string    `db:"id" json:"id"`
		AccountId       string    `db:"account_id" json:"account_id"`
		PartnerId       string    `db:"partner_id" json:"partner_id"`
		PartnerName     string    `db:"partner_name" json:"partner_name"`
		Token           string    `db:"token" json:"token"`
		TransactionCode string    `db:"transaction_code" json:"transaction_code"`
		DiscountValue   float64   `db:"discount_value" json:"discount_value"`
		CreatedAt       time.Time `db:"created_at" json:"created_at"`
		User            string    `db:"created_by" json:"user"`
		Vouchers        []string  `db:"-" json:"vouchers"`
	}
	Cashout struct {
		Id              string           `db:"id" json:"id"`
		PartnerId       string           `db:"partner_id" json:"partner_id"`
		PartnerName     string           `db:"partner_name" json:"partner_name"`
		TransactionCode string           `db:"transaction_code" json:"transaction_code"`
		State           string           `db:"state" json:"state"`
		DiscountValue   float64          `db:"discount_value" json:"discount_value"`
		CreatedAt       time.Time        `db:"created_at" json:"created_at"`
		Vouchers        []CashoutVoucher `db:"-" json:"vouchers"`
	}
	CashoutVoucher struct {
		Id           string `db:"id" json:"id"`
		VoucherCode  string `db:"voucher_code" json:"voucher_code"`
		VoucherState string `db:"state" json:"state"`
		Holder       string `db:"holder" json:"holder"`
	}
	DeleteTransactionRequest struct {
		Id   string `db:"id"`
		User string `db:"deleted_by"`
	}
	TransactionList struct {
		PartnerName     string         `db:"partner_name" json:"partner_name"`
		TransactionId   string         `db:"transaction_id" json:"transaction_id"`
		TransactionCode string         `db:"transaction_code" json:"transaction_code"`
		ProgramName     string         `db:"program_name" json:"program_name"`
		Voucher         []Voucher      `db:"-" json:"vouchers"`
		VoucherValue    float32        `db:"voucher_value" json:"voucher_value"`
		Issued          string         `db:"issued" json:"issued"`
		Redeem          string         `db:"redeemed" json:"redeemed"`
		CashOut         sql.NullString `db:"cashout" json:"cashout"`
		Username        sql.NullString `db:"username" json:"username"`
		State           string         `db:"state" json:"state"`
	}
	VoucherTransaction struct {
		PartnerName     string         `db:"partner_name" json:"partner_name"`
		TransactionId   string         `db:"transaction_id" json:"transaction_id"`
		TransactionCode string         `db:"transaction_code" json:"transaction_code"`
		ProgramName     string         `db:"program_name" json:"program_name"`
		Voucher         Voucher        `db:"-" json:"voucher"`
		VoucherValue    float32        `db:"voucher_value" json:"voucher_value"`
		Issued          string         `db:"issued" json:"issued"`
		Redeem          string         `db:"redeemed" json:"redeemed"`
		CashOut         sql.NullString `db:"cashout" json:"cashout"`
		Username        sql.NullString `db:"username" json:"username"`
		State           string         `db:"state" json:"state"`
	}
)

func InsertTransaction(d Transaction) (string, error) {
	tx, err := db.Beginx()
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	q := `
		INSERT INTO transactions(
			account_id
			, partner_id
			, transaction_code
			, discount_value
			, token
			, created_by
			, status
		)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		RETURNING
			id
	`
	var res []string //[]Transaction
	if err := tx.Select(&res, tx.Rebind(q), d.AccountId, d.PartnerId, d.TransactionCode, d.DiscountValue, d.Token, d.User, StatusCreated); err != nil {
		return "", err
	}
	d.Id = res[0]

	for _, v := range d.Vouchers {
		q := `
			INSERT INTO transaction_details(
				transaction_id
				, voucher_id
				, created_by
				, status
			)
			VALUES (?, ?, ?, ?)
		`

		_, err := tx.Exec(tx.Rebind(q), d.Id, v, d.User, StatusCreated)
		if err != nil {
			return "", err
		}
	}

	if err := tx.Commit(); err != nil {
		return "", err
	}
	return d.Id, nil
}

func (d *Transaction) Update() error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	q := `
		UPDATE transactions
		SET
			company_id = ?
			, pic_merchant = ?
			, transaction_code = ?
			, discount_value = ?
			, updated_by = ?
			, updated_at = ?
		WHERE
			id = ?
			AND status = ?;
	`

	_, err = tx.Exec(tx.Rebind(q), d.AccountId, d.PartnerId, d.TransactionCode, d.DiscountValue, d.User, time.Now(), d.Id, StatusCreated)
	if err != nil {
		return err
	}

	q = `
		UPDATE transaction_details
		SET
			status = ?
			, updated_by = ?
			, updated_at = ?
		WHERE
			transaction_id = ?
			AND status = ?;
	`
	_, err = tx.Exec(tx.Rebind(q), StatusDeleted, d.User, time.Now(), d.Id, StatusCreated)
	if err != nil {
		return err
	}

	for _, v := range d.Vouchers {
		q := `
			INSERT INTO transaction_details(
				transaction_id
				, voucher_id
				, created_by
				, status
			)
			VALUES (?, ?, ?, ?)
			`

		_, err := tx.Exec(tx.Rebind(q), d.Id, v, d.User, StatusCreated)
		if err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (d *DeleteTransactionRequest) Delete() error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	q := `
		UPDATE transactions
		SET
			deleted_by = ?
			, deleted_at = ?
			, status = ?
		WHERE
			id = ?
			AND status = ?;
	`

	_, err = tx.Exec(tx.Rebind(q), d.User, time.Now(), StatusDeleted, d.Id, StatusCreated)
	if err != nil {
		return err
	}

	q = `
		UPDATE transaction_details
		SET
			updated_by = ?
			, updated_at = ?
			, status = ?
		WHERE
			program_id = ?
			AND status = ?;
	`
	_, err = tx.Exec(tx.Rebind(q), d.User, time.Now(), StatusDeleted, d.Id, StatusCreated)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func FindCashoutTransactionDetails(transactionCode string) (Cashout, error) {
	q := `
		SELECT
			t.id
			, p.name as partner_name
			, t.transaction_code
			, t.discount_value
			, t.created_at
		FROM transactions as t
		JOIN partners as p
		ON
			p.id = t.partner_id
		WHERE
			t.transaction_code = ?
			AND t.status = ?
	`

	var resv []Cashout
	if err := db.Select(&resv, db.Rebind(q), transactionCode, StatusCreated); err != nil {
		return Cashout{}, err
	}
	if len(resv) < 1 {
		return Cashout{}, ErrResourceNotFound
	}

	q = `
		SELECT
			v.id, v.voucher_code, v.state, v.holder_description as holder
		FROM
			transaction_details as td
		JOIN
			vouchers as v
		ON
			v.id = td.voucher_id
		WHERE
			td.transaction_id = ?
			AND td.status = ?
	`
	var resd []CashoutVoucher
	if err := db.Select(&resd, db.Rebind(q), resv[0].Id, StatusCreated); err != nil {
		return Cashout{}, err
	}
	if len(resd) < 1 {
		return Cashout{}, ErrResourceNotFound
	}
	resv[0].Vouchers = resd
	resv[0].State = resd[0].VoucherState

	return resv[0], nil
}

func FindVoucherCycle(accountId, voucherId string) (VoucherTransaction, error) {
	q := `
		SELECT
			 t.id as transaction_id, p.name as partner_name, va.name as program_name, t.transaction_code, vo.voucher_value, va.created_at as issued, t.created_at as redeemed, vo.updated_at as cashout, u.username, vo.state
		FROM transactions as t
		JOIN transaction_details as dt
		ON
			t.id = dt.transaction_id
		JOIN vouchers as vo
		ON
			dt.voucher_id = vo.id
		JOIN users as u
		ON
			vo.updated_by = u.id
		JOIN programs as va
		ON
			va.id = vo.program_id
		JOIN partners as p
		ON
			p.id = t.partner_id
		WHERE
			t.status = ?
			AND t.account_id = ?
			AND vo.id = ?
		ORDER BY t.created_at DESC`

	var resv []VoucherTransaction
	if err := db.Select(&resv, db.Rebind(q), StatusCreated, accountId, voucherId); err != nil {
		fmt.Println(err.Error())
		return VoucherTransaction{}, err
	}
	if len(resv) < 1 {
		return VoucherTransaction{}, ErrResourceNotFound
	}

	q1 := `
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
	WHERE
		v.status = ?
		AND v.id = ?
`
	//fmt.Println(q)
	var resv1 []Voucher
	if err := db.Select(&resv1, db.Rebind(q1), StatusCreated, voucherId); err != nil {
		return VoucherTransaction{}, err
	}
	if len(resv) < 1 {
		return VoucherTransaction{}, ErrResourceNotFound
	}
	resv[0].Voucher = resv1[0]

	return resv[0], nil

}

func FindAllTransactionByPartner(accountId, partnerId string) ([]TransactionList, error) {
	q := `
		SELECT
			 t.id as transaction_id
			 , p.name as partner_name
			 , va.name as program_name
			 , t.transaction_code
			 , vo.voucher_value
			 , va.created_at as issued
			 , t.created_at as redeemed
			 , vo.updated_at as cashout
			 , u.username
			 , vo.state
		FROM transactions as t
		JOIN transaction_details as dt
		ON
			t.id = dt.transaction_id
		JOIN vouchers as vo
		ON
			dt.voucher_id = vo.id
		JOIN users as u
		ON
			vo.updated_by = u.id
		JOIN programs as va
		ON
			va.id = vo.program_id
		JOIN partners as p
		ON
			p.id = t.partner_id
		WHERE
			t.status = ?
			AND t.account_id = ?
	`
	q += `AND p.id LIKE '%` + partnerId + `%'`
	q += `ORDER BY t.created_at DESC;`
	//fmt.Println(q)
	var resv []TransactionList
	if err := db.Select(&resv, db.Rebind(q), StatusCreated, accountId); err != nil {
		fmt.Println(err.Error())
		return resv, err
	}
	if len(resv) < 1 {
		return resv, ErrResourceNotFound
	}

	for i, v := range resv {
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
		var resv1 []Voucher
		if err := db.Select(&resv1, db.Rebind(q), StatusCreated, v.TransactionId); err != nil {
			return resv, err
		}
		if len(resv) < 1 {
			return resv, ErrResourceNotFound
		}
		resv[i].Voucher = resv1
	}

	return resv, nil
}

func FindTodayTransactionByPartner(accountId, partnerId string) ([]TransactionList, error) {
	q := `
		SELECT
			 t.id as transaction_id
			 , p.name as partner_name
			 , va.name as program_name
			 , t.transaction_code
			 , vo.voucher_value
			 , va.created_at as issued
			 , t.created_at as redeemed
			 , vo.updated_at as cashout
			 , u.username
			 , vo.state
		FROM transactions as t
		JOIN transaction_details as dt
		ON
			t.id = dt.transaction_id
		JOIN vouchers as vo
		ON
			dt.voucher_id = vo.id
		JOIN users as u
		ON
			vo.updated_by = u.id
		JOIN programs as va
		ON
			va.id = vo.program_id
		JOIN partners as p
		ON
			p.id = t.partner_id
		WHERE
			t.status = ?
			AND t.account_id = ?
	`
	q += `AND p.id LIKE '%` + partnerId + `%'`
	q += `AND t.created_at = ?
		ORDER BY t.created_at DESC;`
	//fmt.Println(q)
	var resv []TransactionList
	if err := db.Select(&resv, db.Rebind(q), StatusCreated, accountId, time.Now()); err != nil {
		fmt.Println(err.Error())
		return resv, err
	}
	if len(resv) < 1 {
		return resv, ErrResourceNotFound
	}

	for i, v := range resv {
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
		var resv1 []Voucher
		if err := db.Select(&resv1, db.Rebind(q), StatusCreated, v.TransactionId); err != nil {
			return resv, err
		}
		if len(resv) < 1 {
			return resv, ErrResourceNotFound
		}
		resv[i].Voucher = resv1
	}

	return resv, nil
}

func UpdateCashoutTransaction(transactionCode, user string) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	q := `
		UPDATE vouchers
		SET
			state = ?
			, updated_by = ?
			, updated_at = ?
		WHERE
		  id = (SELECT v.id
			FROM vouchers as v
			JOIN transaction_details as td
			ON
				v.id = td.voucher_id
			JOIN transactions as t
			ON
				t.id = td.transaction_id
			WHERE
				t.transaction_code = ?
				AND t.status = ?
			)
	`

	_, err = tx.Exec(tx.Rebind(q), VoucherStatePaid, user, time.Now(), transactionCode, StatusCreated)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func UpdateCashoutTransactions(transactionCode []string, user string) error {
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

	for i, v := range transactionCode {
		if i != 0 {
			q += ` OR `
		}
		q += `t.transaction_code = '` + v + `'`
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

func PrintCashout(accountId string, transactionCode []string) ([]TransactionList, error) {
	q := `
		SELECT
			t.id as transaction_id, p.name as partner_name, t.transaction_code, vo.voucher_value, va.created_at as issued, t.created_at as redeemed, vo.updated_at as cashout, u.username, vo.state
		FROM transactions as t
		JOIN transaction_details as dt
		ON
			t.id = dt.transaction_id
		JOIN vouchers as vo
		ON
			dt.voucher_id = vo.id
		JOIN users as u
		ON
			vo.updated_by = u.id
		JOIN programs as va
		ON
			va.id = vo.program_id
		JOIN partners as p
		ON
			p.id = t.partner_id
		WHERE
			t.status = ?
			AND t.account_id = ?
			AND vo.state = ?
			AND (
	`

	for i, v := range transactionCode {
		if i != 0 {
			q += ` OR `
		}
		q += `t.transaction_code = '` + v + `'`
	}
	q += `) ORDER BY t.created_at DESC;`
	var resv []TransactionList
	if err := db.Select(&resv, db.Rebind(q), StatusCreated, accountId, VoucherStatePaid); err != nil {
		fmt.Println(err.Error())
		return resv, err
	}
	if len(resv) < 1 {
		return resv, ErrResourceNotFound
	}

	for i, v := range resv {
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
		var resv1 []Voucher
		if err := db.Select(&resv1, db.Rebind(q), StatusCreated, v.TransactionId); err != nil {
			return resv, err
		}
		if len(resv) < 1 {
			return resv, ErrResourceNotFound
		}
		resv[i].Voucher = resv1
	}

	return resv, nil
}
