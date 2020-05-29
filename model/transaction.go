package model

import (
	"github.com/gilkor/evoucher-v2/util"
	"time"
)

// [UNDER CONSTRUCTION!!]

type (
	Transaction struct {
		ID                 string     `db:"id",json:"id"`
		CompanyId          string     `db:"company_id",json:"company_id"`
		TransactionCode    string     `db:"transaction_code",json:"transaction_code"`
		TotalAmount        string     `db:"total_amount",json:"total_amount"`
		Holder             string     `db:"holder",json:"holder"`
		PartnerId          string     `db:"partner_id",json:"partner_id"`
		CreatedBy          string     `db:"created_by",json:"created_by"`
		CreatedAt          *time.Time `db:"created_at",json:"created_at"`
		UpdatedBy          string     `db:"updated_by",json:"updated_by"`
		UpdatedAt          *time.Time `db:"updated_at",json:"updated_at"`
		Status             string     `db:"status",json:"status"`
		TransactionDetails TransactionDetails
	}
	Transactions      []Transaction
	TransactionDetail struct {
		ID            string     `db:"id",json:"id"`
		TransactionId string     `db:"transaction_id",json:"transaction_id"`
		ProgramId     string     `db:"program_id",json:"program_id"`
		VoucherId     string     `db:"voucher_id",json:"voucher_id"`
		CreatedBy     string     `db:"created_by",json:"created_by"`
		CreatedAt     *time.Time `db:"created_at",json:"created_at"`
		UpdatedBy     string     `db:"updated_by",json:"updated_by"`
		UpdatedAt     *time.Time `db:"updated_at",json:"updated_at"`
		Status        string     `db:"status",json:"status"`
	}
	TransactionDetails []TransactionDetail
)

//GetTransactions :
func GetTransactions(qp *util.QueryParam) (*Transactions, bool, error) {
	return getTransactions("1", "1", qp)
}

//GetTransactionByID : get partner by specified ID
func GetTransactionByID(qp *util.QueryParam, id string) (*Transactions, bool, error) {
	return getTransactions("id", id, qp)
}

//GetTransactionByProgram : get partner by specified ProgramID
func GetTransactionByProgram(qp *util.QueryParam, val string) (*Transactions, bool, error) {
	return getTransactions("program_id", val, qp)
}

//GetTransactionByPartner : get partner by specified PartnerID
func GetTransactionByPartner(qp *util.QueryParam, val string) (*Transactions, bool, error) {
	return getTransactions("partner_id", val, qp)
}

func getTransactions(k, v string, qp *util.QueryParam) (*Transactions, bool, error) {

	q, err := qp.GetQueryByDefaultStruct(Transaction{})
	if err != nil {
		return &Transactions{}, false, err
	}

	q += `
			FROM
				m_transactions transaction
			WHERE 
				status = ?
			AND ` + k + ` = ?`

	q += qp.GetQuerySort()
	q += qp.GetQueryLimit()
	var resd Transactions
	err = db.Select(&resd, db.Rebind(q), StatusCreated, v)
	if err != nil {
		return &Transactions{}, false, err
	}
	if len(resd) < 1 {
		return &Transactions{}, false, ErrorResourceNotFound
	}
	next := false
	if len(resd) > qp.Count {
		next = true
	}
	if len(resd) < qp.Count {
		qp.Count = len(resd)
	}
	return &resd, next, nil
}

//Insert : transaction data
func (t Transaction) Insert() (*[]Transaction, error) {
	tx, err := db.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	q := `INSERT INTO 
			transactions (company_id, transaction_code, total_amount, holder, partner_id, created_by, updated_by, status)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)
			RETURNING
				id, company_id, transaction_code, total_amount, holder, partner_id, created_by, created_at, updated_by, updated_at, status
	`
	var res []Transaction
	err = tx.Select(&res, tx.Rebind(q), t.CompanyId, t.TransactionCode, t.TotalAmount, t.Holder, t.PartnerId, t.CreatedBy, t.UpdatedBy, StatusCreated)
	if err != nil {
		return nil, err
	}

	for _, td := range t.TransactionDetails {
		q = `INSERT INTO 
			transaction_details (transaction_id, program_id, voucher_id, created_by, updated_by, status)
			VALUES (?, ?, ?, ?, ?, ?)
			RETURNING
				id, transaction_id, program_id, voucher_id, created_by, updated_by, status
	`
		var res []Transaction
		err = tx.Select(&res, tx.Rebind(q), td.TransactionId, td.ProgramId, td.VoucherId, td.CreatedBy, td.UpdatedBy, StatusCreated)
		if err != nil {
			return nil, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &res, nil
}

////Insert : transaction data
//func (t TransactionDetail) Insert() (*[]TransactionDetail, error) {
//	tx, err := db.Beginx()
//	if err != nil {
//		return nil, err
//	}
//	defer tx.Rollback()
//
//	q := `INSERT INTO
//			transactions (company_id, transaction_code, total_amount, holder, partner_id, created_by, updated_by, status)
//			VALUES (?, ?, ?, ?, ?, ?, ?, ?)
//			RETURNING
//				id, company_id, transaction_code, total_amount, holder, partner_id, created_by, created_at, updated_by, updated_at, status
//	`
//	var res []TransactionDetail
//	err = tx.Select(&res, tx.Rebind(q), t.CompanyId, t.TransactionCode, t.TotalAmount, t.Holder, t.PartnerId, t.CreatedBy, t.UpdatedBy, StatusCreated)
//	if err != nil {
//		return nil, err
//	}
//
//	for i, i2 := range t.TransactionDetails {
//
//	}
//
//	err = tx.Commit()
//	if err != nil {
//		return nil, err
//	}
//
//	return &res, nil
//}

//Update : Transaction
//There is no update for transaction yet.
// func (t Transaction) Update() error {
// 	tx, err := db.Beginx()
// 	if err != nil {
// 		return err
// 	}
// 	defer tx.Rollback()

// 	q := `INSERT INTO

// 	`
// 	var res []Transaction
// 	err = tx.Select(&res, tx.Rebind(q))
// 	if err != nil {
// 		return err
// 	}
// 	err = tx.Commit()
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

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
