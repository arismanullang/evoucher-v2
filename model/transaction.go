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

	q := `
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

	q := `
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
