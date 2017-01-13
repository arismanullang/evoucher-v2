package model

import (
	"time"
)

type (
	Transaction struct {
		CompanyID        string   `db:"company_id"`
		MerchantID       string   `db:"pic_merchant"`
		TransactionCode  string   `db:"transaction_code"`
		TotalTransaction float64  `db:"total_transaction"`
		DiscountValue    float64  `db:"discount_value"`
		PaymentType      string   `db:"payment_type"`
		User             string   `db:"created_by"`
		Vouchers         []string `db:"-"`
	}
	TransactionResponse struct {
		Status           string
		Message          string
		TransactionValue Transaction
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
