package model

import (
	"time"
)

type (
	Transaction struct {
		Id              string    `db:"id"`
		AccountId       string    `db:"account_id"`
		PartnerId       string    `db:"partner_id"`
		TransactionCode string    `db:"transaction_code"`
		Token           string    `db:"token"`
		DiscountValue   float64   `db:"discount_value"`
		CreatedAt       time.Time `db:"created_at"`
		User            string    `db:"created_by"`
		Vouchers        []string  `db:"-"`
	}
	DeleteTransactionRequest struct {
		Id   string `db:"id"`
		User string `db:"deleted_by"`
	}
)

func InsertTransaction(d Transaction) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
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
		VALUES (?, ?, ?, ?, ?, ?)
		RETURNING
			id
	`
	var res []string //[]Transaction
	if err := tx.Select(&res, tx.Rebind(q), d.AccountId, d.PartnerId, d.TransactionCode, d.DiscountValue, d.Token, d.User, StatusCreated); err != nil {
		return err
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
			variant_id = ?
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

func FindTransactionByID(id string) ([]Transaction, error) {
	q := `
		SELECT
			id
			, account_id
			, partner_id
			, transaction_code
			, total_transaction
			, discount_value
			, token
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
		return []Transaction{}, err
	}
	if len(resv) < 1 {
		return []Transaction{}, ErrResourceNotFound
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
		return []Transaction{}, err
	}
	if len(resd) < 1 {
		return []Transaction{}, ErrResourceNotFound
	}
	resv[0].Vouchers = resd

	return resv, nil
}

func FindTransactionByDate(start, end string) ([]Transaction, error) {
	q := `
		SELECT
			transaction_code
			, company_id
			, pic_merchant
			, total_transaction
			, discount_value
			, token
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
		return []Transaction{}, err
	}
	if len(resv) < 1 {
		return []Transaction{}, ErrResourceNotFound
	}

	return resv, nil
}
