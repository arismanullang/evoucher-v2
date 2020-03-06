package model

import "time"

// [UNDER CONSTRUCTION!!]

type (
	//Transaction : obj
	Transaction struct {
		ID                 string
		CompanyID          string
		TransactionCode    string
		TotalAmount        float64
		Holder             string
		PartnerID          string
		CreatedBy          string
		CreatedAt          *time.Time
		UpdatedBy          string
		UpdatedAt          *time.Time
		Status             string
		TransactionDetails []TransactionDetail
	}
	//TransactionDetail : obj
	TransactionDetail struct {
		ID            string
		TransactionID string
		VoucherID     string
		CreatedBy     string
		CreatedAt     *time.Time
		UpdatedBy     string
		UpdatedAt     *time.Time
		Status        string
	}
)

//Insert : transaction data
func (t Transaction) Insert() (*[]Transaction, error) {
	tx, err := db.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	q := `INSERT INTO 
			transactions (company_id, transaction_code, total_amount, holder, partner_id, created_by, created_at, updated_by, updated_at, status)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			RETURNING
				id, company_id, transaction_code, total_amount, holder, partner_id, created_by, created_at, updated_by, updated_at, status
	`
	var res []Transaction
	err = tx.Select(&res, tx.Rebind(q))
	if err != nil {
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &res, nil
}

//Update : Transaction
//There is no update for transaction yet
func (t Transaction) Update() error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	q := `INSERT INTO 
				
	`
	var res []Transaction
	err = tx.Select(&res, tx.Rebind(q))
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

//Delete : Soft Delete Transaction data
func (t Transaction) Delete() error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	q := `UPDATE FROM 
				transactions
			SET
			 status = ?
			WHERE
			 id = ?
	`
	var res []Transaction
	err = tx.Select(&res, tx.Rebind(q), StatusDeleted, t.ID)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
