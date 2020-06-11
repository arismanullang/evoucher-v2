package model

import (
	"encoding/json"
	"time"

	"github.com/gilkor/evoucher-v2/util"
	"github.com/jmoiron/sqlx/types"
)

// [UNDER CONSTRUCTION!!]

type (
	Transaction struct {
		ID                   string             `db:"id" json:"id"`
		CompanyId            string             `db:"company_id" json:"company_id"`
		TransactionCode      string             `db:"transaction_code" json:"transaction_code"`
		TotalAmount          string             `db:"total_amount" json:"total_amount"`
		Holder               string             `db:"holder" json:"holder"`
		Vouchers             Vouchers           `json:"-"`
		Programs             Programs           `json:"programs,omitempty"`
		Partner              Partner            `json:"partner,omitempty"`
		PartnerId            string             `db:"partner_id" json:"partner_id"`
		CreatedBy            string             `db:"created_by" json:"created_by"`
		CreatedAt            *time.Time         `db:"created_at" json:"created_at"`
		UpdatedBy            string             `db:"updated_by" json:"updated_by"`
		UpdatedAt            *time.Time         `db:"updated_at" json:"updated_at"`
		Status               string             `db:"status" json:"status"`
		TransactionDetailsDB types.JSONText     `db:"transaction_details" json:"-"`
		TransactionDetails   TransactionDetails `json:"transaction_details,omitempty"`
	}
	Transactions      []Transaction
	TransactionDetail struct {
		ID            int        `db:"id" json:"id"`
		TransactionId string     `db:"transaction_id" json:"transaction_id"`
		ProgramId     string     `db:"program_id" json:"program_id"`
		VoucherId     string     `db:"voucher_id" json:"voucher_id"`
		CreatedBy     string     `db:"created_by" json:"created_by"`
		CreatedAt     *time.Time `db:"created_at" json:"created_at"`
		UpdatedBy     string     `db:"updated_by" json:"updated_by"`
		UpdatedAt     *time.Time `db:"updated_at" json:"updated_at"`
		Status        string     `db:"status" json:"status"`
	}
	TransactionDetails []TransactionDetail
)

//GetTransactions :
func GetTransactions(qp *util.QueryParam) (*Transactions, bool, error) {
	return getTransactions("1", "1", qp)
}

//GetTransactionByID : get partner by specified ID
func GetTransactionByID(qp *util.QueryParam, id string) (Transaction, bool, error) {
	// return
	transactions, _, err := getTransactions("id", id, qp)
	if err != nil {
		return Transaction{}, false, err
	}

	if len(*transactions) > 0 {
		transaction := (*transactions)[0]
		return transaction, false, nil
	}

	return Transaction{}, false, ErrorResourceNotFound
}

//GetTransactionByHolder : get transaction by specified Holder
func GetTransactionByHolder(qp *util.QueryParam, val string) (*Transactions, bool, error) {
	// SELECT DISTINCT
	// 		t.id,
	// 		t.created_at,
	// 		t.discount_value,
	// 		t.transaction_code,
	// 		pt.id as partner_id,
	// 		pt.name	as partner_name
	// 	FROM transactions as t
	// 	JOIN partners as pt
	// 		ON t.partner_id = pt.id
	// 	WHERE t.status = ?
	// 	AND t.holder = ?
	// 	ORDER BY t.created_at DESC`
	return getTransactions("holder", val, qp)
}

//GetTransactionByProgram : get transaction by specified ProgramID
func GetTransactionByProgram(qp *util.QueryParam, val string) (*Transactions, bool, error) {
	return getTransactions("program_id", val, qp)
}

//GetTransactionByPartner : get transaction by specified PartnerID
func GetTransactionByPartner(qp *util.QueryParam, val string) (*Transactions, bool, error) {
	return getTransactions("partner_id", val, qp)
}

//GetTransactionDetailByVoucherID : get transaction detail by voucher ID
func GetTransactionDetailByVoucherID(qp *util.QueryParam, voucherID string) (*TransactionDetail, error) {
	q, err := qp.GetQueryByDefaultStruct(TransactionDetail{})
	if err != nil {
		return &TransactionDetail{}, err
	}

	q += `
			FROM
				transaction_details as TransactionDetail
			WHERE 
				status = ?
			AND voucher_id = ?`

	var resd []TransactionDetail
	err = db.Select(&resd, db.Rebind(q), StatusCreated, voucherID)
	if err != nil {
		return &TransactionDetail{}, err
	}

	if len(resd) < 1 {
		return &TransactionDetail{}, ErrorResourceNotFound
	}

	return &resd[0], nil
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

	for i := range resd {
		err = json.Unmarshal([]byte(resd[i].TransactionDetailsDB), &resd[i].TransactionDetails)
		if err != nil {
			return &Transactions{}, false, err
		}
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

	var resd []TransactionDetail
	for _, td := range t.TransactionDetails {
		td.TransactionId = res[0].ID
		q = `INSERT INTO 
			transaction_details (transaction_id, program_id, voucher_id, created_by, updated_by, status)
			VALUES (?, ?, ?, ?, ?, ?)
			RETURNING
				id, transaction_id, program_id, voucher_id, created_by, created_at, updated_by, updated_at, status
	`
		err = tx.Select(&resd, tx.Rebind(q), td.TransactionId, td.ProgramId, td.VoucherId, td.CreatedBy, td.UpdatedBy, StatusCreated)
		if err != nil {
			return nil, err
		}
	}
	res[0].TransactionDetails = resd

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
