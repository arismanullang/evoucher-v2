package model

import (
	"time"
)

type (
	Transaction struct {
		ID               string    `db:"id"`
		CompanyID        string    `db:"company_id"`
		MerchantID       string    `db:"pic_merchant"`
		TransactionCode  string    `db:"transaction_code"`
		TotalTransaction float64   `db:"total_transaction"`
		DiscountValue    float64   `db:"discount_value"`
		PaymentType      string    `db:"payment_type"`
		CreatedAt        time.Time `db:"created_at"`
		User             string    `db:"created_by"`
		Vouchers         []string  `db:"-"`
	}
	TransactionsResponse struct {
		Status           string
		Message          string
		TransactionValue []Transaction
	}
	DeleteTransactionRequest struct {
		ID   string `db:"id"`
		User string `db:"deleted_by"`
	}
)

func (d *Transaction) Insert() error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	q := `
		INSERT INTO transactions(
			company_id
			, pic_merchant
			, transaction_code
			, total_transaction
			, discount_value
			, payment_type
			, created_by
		)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		RETURNING
			id
	`
	var res []string
	if err := tx.Select(&res, tx.Rebind(q), d.CompanyID, d.MerchantID, d.TransactionCode, d.TotalTransaction, d.DiscountValue, d.PaymentType, d.User); err != nil {
		return err
	}
	d.ID = res[0]

	for _, v := range d.Vouchers {
		q := `
			INSERT INTO transaction_details(
				company_id
				, transaction_id
				, voucher_id
				, created_by
				, status
			)
			VALUES (?, ?, ?, ?, ?)
		`

		_, err := tx.Exec(tx.Rebind(q), d.CompanyID, d.ID, v, d.User, StatusCreated)
		if err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
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
			, total_transaction = ?
			, discount_value = ?
			, payment_type = ?
			, updated_by = ?
			, updated_at = ?
		WHERE
			id = ?
			AND stauts = ?;
	`

	_, err = tx.Exec(tx.Rebind(q), d.CompanyID, d.MerchantID, d.TransactionCode, d.TotalTransaction, d.DiscountValue, d.PaymentType, d.User, time.Now(), d.ID, StatusCreated)
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
	_, err = tx.Exec(tx.Rebind(q), StatusDeleted, d.User, time.Now(), d.ID, StatusCreated)
	if err != nil {
		return err
	}

	for _, v := range d.Vouchers {
		q := `
			INSERT INTO transaction_details(
				company_id
				, transaction_id
				, voucher_id
				, created_by
			)
			VALUES (?, ?, ?, ?)
			`

		_, err := tx.Exec(tx.Rebind(q), d.CompanyID, d.ID, v, d.User)
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

	_, err = tx.Exec(tx.Rebind(q), d.User, time.Now(), StatusDeleted, d.ID, StatusCreated)
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
			variant_id = ?
			AND status = ?;
	`
	_, err = tx.Exec(tx.Rebind(q), d.User, time.Now(), StatusDeleted, d.ID, StatusCreated)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func FindTransactionByID(id string) (TransactionsResponse, error) {
	q := `
		SELECT
			id
			, company_id
			, pic_merchant
			, transaction_code
			, total_transaction
			, discount_value
			, payment_type
			, created_by
			, created_at
		FROM
			transactions
		WHERE
			id = ?
			AND status = ?
	`

	var resv []Transaction
	if err := db.Select(&resv, db.Rebind(q), id, StatusCreated); err != nil {
		return TransactionsResponse{Status: "500", Message: "Error at select variant", TransactionValue: []Transaction{}}, err
	}
	if len(resv) < 1 {
		return TransactionsResponse{Status: "404", Message: "Error at select variant", TransactionValue: []Transaction{}}, ErrResourceNotFound
	}

	q = `
		SELECT
			voucher_id
		FROM
			transaction_details
		WHERE
			transaction_id = ?
			AND status = ?
	`
	var resd []string
	if err := db.Select(&resd, db.Rebind(q), id, StatusCreated); err != nil {
		return TransactionsResponse{Status: "500", Message: "Error at select user", TransactionValue: []Transaction{}}, err
	}
	if len(resd) < 1 {
		return TransactionsResponse{Status: "404", Message: "Error at select user", TransactionValue: []Transaction{}}, ErrResourceNotFound
	}
	resv[0].Vouchers = resd

	res := TransactionsResponse{
		Status:           "200",
		Message:          "Ok",
		TransactionValue: resv,
	}

	return res, nil
}

func FindTransactionByDate(start, end string) (TransactionsResponse, error) {
	q := `
		SELECT
			transaction_code
			, company_id
			, pic_merchant
			, total_transaction
			, discount_value
			, payment_type
			, created_by
			, created_at
		FROM
			transactions
		WHERE
			(start_date > ? AND start_date < ?)
			OR (end_date > ? AND end_date < ?)
			AND status = ?
	`

	var resv []Transaction
	if err := db.Select(&resv, db.Rebind(q), start, end, start, end, StatusCreated); err != nil {
		return TransactionsResponse{Status: "500", Message: "Error at select variant", TransactionValue: []Transaction{}}, err
	}
	if len(resv) < 1 {
		return TransactionsResponse{Status: "404", Message: "Error at select variant", TransactionValue: []Transaction{}}, ErrResourceNotFound
	}

	res := TransactionsResponse{
		Status:           "200",
		Message:          "Ok",
		TransactionValue: resv,
	}

	return res, nil
}
