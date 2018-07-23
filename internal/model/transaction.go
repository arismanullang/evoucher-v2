package model

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"
)

type (
	Transaction struct {
		Id              string    `db:"id" json:"id"`
		AccountId       string    `db:"account_id" json:"account_id"`
		PartnerId       string    `db:"partner_id" json:"partner_id"`
		PartnerName     string    `db:"partner_name" json:"partner_name"`
		Holder          string    `db:"holder" json:"holder"`
		Token           string    `db:"token" json:"token"`
		TransactionCode string    `db:"transaction_code" json:"transaction_code"`
		DiscountValue   float64   `db:"discount_value" json:"discount_value"`
		CreatedAt       time.Time `db:"created_at" json:"created_at"`
		User            string    `db:"created_by" json:"user"`
		VoucherIds      []string  `db:"-" json:"voucher_ids"`
		Vouchers        []Voucher `db:"-" json:"vouchers"`
	}
	TransactionCashout struct {
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
		PartnerId       string         `db:"partner_id" json:"partner_id"`
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

	TransactionHistoryList struct {
		TransactionID   string    `db:"id" json:"transaction_id"`
		CreatedAt       time.Time `db:"created_at" json "created_at`
		TransactionCode string    `db:"transaction_code" json:"transaction_code"`
		DiscountValue   float64   `db:"discount_value" json:"discount_value"`
		PartnerID       string    `db:"partner_id" json:"partner_id"`
		PartnerName     string    `db:"partner_name" json:"partner_name"`
	}

	TransactionHistoryDetail struct {
		VoucherID        string  `db:"id" json:"id"`
		VoucherCode      string  `db:"voucher_code" json:"voucher_code"`
		ProgramID        string  `db:"program_id" json:"program_id"`
		ProgramName      string  `db:"program_name" json:"program_name"`
		VoucherValue     float64 `db:"voucher_value" json:"voucher_value"`
		ProgramStartDate string  `db:"program_start_date" json:"program_start_date"`
		ProgramEndDate   string  `db:"program_end_date" json:"program_end_date"`
		ProgramImgUrl    string  `db:"program_img_url" json:"program_img_url"`
	}
)

func InsertTransaction(d Transaction) (Transaction, error) {
	tx, err := db.Beginx()
	if err != nil {
		return d, err
	}
	defer tx.Rollback()

	q := `
		INSERT INTO transactions(
			account_id
			, partner_id
			, holder
			, transaction_code
			, discount_value
			, token
			, created_by
			, created_at
			, status
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING id, account_id, partner_id, token, transaction_code, discount_value,
		created_at, created_by

	`
	var res []Transaction
	if err := tx.Select(&res, tx.Rebind(q), d.AccountId, d.PartnerId, d.Holder, d.TransactionCode, d.DiscountValue, d.Token, d.User, time.Now(), StatusCreated); err != nil {
		return d, err
	}
	fmt.Println("insert transaction response :", res)
	d.Id = res[0].Id

	for _, v := range d.VoucherIds {
		q := `
			INSERT INTO transaction_details(
				transaction_id
				, voucher_id
				, created_by
				, created_at
				, status
			)
			VALUES (?, ?, ?, ?, ?)
		`

		_, err := tx.Exec(tx.Rebind(q), d.Id, v, d.User, res[0].CreatedAt, StatusCreated)
		if err != nil {
			return d, err
		}
	}

	if err := tx.Commit(); err != nil {
		return d, err
	}

	return res[0], nil
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
			transaction_id = ?
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

func FindCashoutTransactionDetails(transactionCode string) (TransactionCashout, error) {
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

	var resv []TransactionCashout
	if err := db.Select(&resv, db.Rebind(q), transactionCode, StatusCreated); err != nil {
		return TransactionCashout{}, err
	}
	if len(resv) < 1 {
		return TransactionCashout{}, ErrResourceNotFound
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
		return TransactionCashout{}, err
	}
	if len(resd) < 1 {
		return TransactionCashout{}, ErrResourceNotFound
	}
	resv[0].Vouchers = resd
	resv[0].State = resd[0].VoucherState

	return resv[0], nil
}

func FindVoucherCycle(accountId, voucherId string) (VoucherTransaction, error) {
	q := `
		SELECT DISTINCT
			 t.id as transaction_id
			, p.name as partner_name
			, va.name as program_name
			, t.transaction_code
			, vo.voucher_value
			, va.created_at as issued
			, t.created_at as redeemed
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

	q = `
		SELECT DISTINCT
			c.created_at as cashout
		FROM cashout_details as cd
		JOIN cashouts as c
		ON
			cd.cashout_id = c.id
		WHERE
			cd.transaction_id = ?`

	var cashoutDate []string
	if err := db.Select(&cashoutDate, db.Rebind(q), resv[0].TransactionId); err != nil {
		fmt.Println(err.Error())
	}
	if len(cashoutDate) > 0 {
		resv[0].CashOut.String = cashoutDate[0]
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

func FindTransactionsByPartner(accountId, partnerId string) ([]TransactionList, error) {
	q := `
		SELECT DISTINCT
			 t.id as transaction_id
			 , p.id as partner_id
			 , p.name as partner_name
			 , va.name as program_name
			 , t.transaction_code
			 , va.voucher_value
			 , vo.created_at as issued
			 , t.created_at as redeemed
			 , u.username
		FROM transactions as t
		JOIN transaction_details as dt
		ON
			t.id = dt.transaction_id
		JOIN vouchers as vo
		ON
			dt.voucher_id = vo.id
		JOIN programs as va
		ON
			va.id = vo.program_id
		JOIN partners as p
		ON
			p.id = t.partner_id
		JOIN users as u
		ON
			u.id = t.created_by
		WHERE
			t.status = ?
			AND t.account_id = ?
	`
	q += `AND p.id LIKE '%` + partnerId + `%' `
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

		q = `
		SELECT
			c.created_at
		FROM cashouts as c
		JOIN cashout_details as ct
		ON
			c.id = ct.cashout_id
		WHERE
			c.status = ?
			AND ct.transaction_id = ?
	`
		//fmt.Println(q)
		var res []string
		if err := db.Select(&res, db.Rebind(q), StatusCreated, v.TransactionId); err != nil {
			return resv, err
		}
		resv[i].CashOut.String = ""
		if len(res) > 0 {
			resv[i].CashOut.String = res[0]
		}
	}

	return resv, nil
}

func FindTransactionsByDate(accountId string, createdAt time.Time) ([]TransactionList, error) {
	q := `
		SELECT DISTINCT
			 t.id as transaction_id
			 , p.id as partner_id
			 , p.name as partner_name
			 , va.name as program_name
			 , t.transaction_code
			 , va.voucher_value
			 , va.created_at as issued
			 , t.created_at as redeemed
			 , u.username
		FROM transactions as t
		JOIN transaction_details as dt
		ON
			t.id = dt.transaction_id
		JOIN vouchers as vo
		ON
			dt.voucher_id = vo.id
		JOIN programs as va
		ON
			va.id = vo.program_id
		JOIN partners as p
		ON
			p.id = t.partner_id
		JOIN users as u
		ON
			u.id = t.created_by
		WHERE
			t.status = ?
			AND t.account_id = ?
			AND t.created_at BETWEEN ? AND ?
			AND vo.state = 'used'
	`
	q += `ORDER BY t.created_at DESC;`
	//fmt.Println(q)
	var resv []TransactionList
	if err := db.Select(&resv, db.Rebind(q), StatusCreated, accountId, createdAt, createdAt.AddDate(0, 0, 1)); err != nil {
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

		q = `
		SELECT
			c.created_at
		FROM cashouts as c
		JOIN cashout_details as ct
		ON
			c.id = ct.cashout_id
		WHERE
			c.status = ?
			AND ct.transaction_id = ?
	`
		//fmt.Println(q)
		var res []string
		if err := db.Select(&res, db.Rebind(q), StatusCreated, v.TransactionId); err != nil {
			return resv, err
		}
		resv[i].CashOut.String = ""
		if len(res) > 0 {
			resv[i].CashOut.String = res[0]
		}
	}

	return resv, nil
}

func FindTransactionsByPartnerDate(accountId, partnerId string, createdAt time.Time) ([]TransactionList, error) {
	q := `
		SELECT DISTINCT
			 t.id as transaction_id
			 , p.id as partner_id
			 , p.name as partner_name
			 , va.name as program_name
			 , t.transaction_code
			 , va.voucher_value
			 , va.created_at as issued
			 , t.created_at as redeemed
			 , u.username
		FROM transactions as t
		JOIN transaction_details as dt
		ON
			t.id = dt.transaction_id
		JOIN vouchers as vo
		ON
			dt.voucher_id = vo.id
		JOIN programs as va
		ON
			va.id = vo.program_id
		JOIN partners as p
		ON
			p.id = t.partner_id
		JOIN users as u
		ON
			u.id = t.created_by
		WHERE
			t.status = ?
			AND t.account_id = ?
			AND p.id = ?
			AND t.created_at BETWEEN ? AND ?
	`
	q += `ORDER BY t.created_at DESC;`
	//fmt.Println(q)
	var resv []TransactionList
	if err := db.Select(&resv, db.Rebind(q), StatusCreated, accountId, partnerId, createdAt, createdAt.AddDate(0, 0, 1)); err != nil {
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

		q = `
		SELECT
			c.created_at
		FROM cashouts as c
		JOIN cashout_details as ct
		ON
			c.id = ct.cashout_id
		WHERE
			c.status = ?
			AND ct.transaction_id = ?
	`
		//fmt.Println(q)
		var res []string
		if err := db.Select(&res, db.Rebind(q), StatusCreated, v.TransactionId); err != nil {
			return resv, err
		}
		resv[i].CashOut.String = ""
		if len(res) > 0 {
			resv[i].CashOut.String = res[0]
		}
	}

	return resv, nil
}

func FindTransactions(param map[string]string) ([]TransactionList, error) {
	q := `
		SELECT DISTINCT
			 t.id as transaction_id
			 , p.id as partner_id
			 , p.name as partner_name
			 , va.name as program_name
			 , t.transaction_code
			 , t.discount_value as voucher_value
			 , va.created_at as issued
			 , t.created_at as redeemed
			 , u.username
		FROM transactions as t
		JOIN transaction_details as dt
		ON
			t.id = dt.transaction_id
		JOIN vouchers as vo
		ON
			dt.voucher_id = vo.id
		JOIN programs as va
		ON
			va.id = vo.program_id
		JOIN partners as p
		ON
			p.id = t.partner_id
		JOIN users as u
		ON
			u.id = t.created_by
		WHERE
			t.status = ?
	`
	for key, value := range param {
		if _, err := strconv.Atoi(value); err == nil {
			q += ` AND ` + key + ` = '` + value + `'`
		} else {
			q += ` AND ` + key + ` LIKE '%` + value + `%'`

		}
	}
	q += `ORDER BY t.created_at DESC;`
	//fmt.Println(q)
	var resv []TransactionList
	if err := db.Select(&resv, db.Rebind(q), StatusCreated); err != nil {
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
		SELECT DISTINCT
			 t.id as transaction_id
			 , p.name as partner_name
			 , va.name as program_name
			 , t.transaction_code
			 , va.voucher_value
			 , va.created_at as issued
			 , t.created_at as redeemed
			 , c.created_at as cashout
			 , u.username
		FROM transactions as t
		JOIN transaction_details as dt
		ON
			t.id = dt.transaction_id
		JOIN vouchers as vo
		ON
			dt.voucher_id = vo.id
		JOIN programs as va
		ON
			va.id = vo.program_id
		JOIN partners as p
		ON
			p.id = t.partner_id
		JOIN cashout_details as cd
		ON
			t.id = cd.transaction_id
		JOIN cashouts as c
		ON
			cd.cashout_id = c.id
		JOIN users as u
		ON
			u.id = c.created_by
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

//GetTransactionListByHolder get list of transaction ids by holder
func GetTransactionListByHolder(holder string) ([]TransactionHistoryList, error) {
	q := `
		SELECT DISTINCT
			t.id,
			t.created_at,
			t.discount_value,
			t.transaction_code,
			pt.id as partner_id,
			pt.name	as partner_name
		FROM transactions as t
		JOIN partners as pt
			ON t.partner_id = pt.id
		WHERE t.status = ?
		AND t.holder = ?
		ORDER BY t.created_at DESC`

	var listTransactionHistory []TransactionHistoryList
	if err := db.Select(&listTransactionHistory, db.Rebind(q), StatusCreated, holder); err != nil {
		fmt.Println(err.Error())
		return listTransactionHistory, err
	}
	if len(listTransactionHistory) < 1 {
		return listTransactionHistory, ErrResourceNotFound
	}

	return listTransactionHistory, nil
}

//GetVoucherByTransaction Get Transaction History Detail
func GetVoucherByTransaction(transactionID string) ([]TransactionHistoryDetail, error) {
	q := `
	SELECT
		v.id,
		v.voucher_code,
		p.id as program_id,
		p.name as program_name,
		p.voucher_value,
		p.start_date as program_start_date,
		p.end_date as program_end_date,
		p.img_url as program_img_url
	FROM transaction_details as td
	JOIN
		vouchers as v ON td.voucher_id = v.id
	JOIN
		programs as p ON v.program_id = p.id
	WHERE
		td.status = ?
		AND td.transaction_id = ?
	ORDER BY v.voucher_code DESC`

	var transactionHistoryDetail []TransactionHistoryDetail
	if err := db.Select(&transactionHistoryDetail, db.Rebind(q), StatusCreated, transactionID); err != nil {
		fmt.Println(err.Error())
		return transactionHistoryDetail, err
	}
	if len(transactionHistoryDetail) < 1 {
		return transactionHistoryDetail, ErrResourceNotFound
	}

	return transactionHistoryDetail, nil
}
