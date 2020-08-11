package model

import (
	"time"

	"github.com/gilkor/evoucher-v2/util"
	"github.com/jmoiron/sqlx/types"
)

// [UNDER CONSTRUCTION!!]

type (
	Transaction struct {
		ID                 string             `db:"id" json:"id"`
		CompanyID          string             `db:"company_id" json:"company_id"`
		TransactionCode    string             `db:"transaction_code" json:"transaction_code"`
		TotalAmount        string             `db:"total_amount" json:"total_amount"`
		Holder             string             `db:"holder" json:"holder"`
		Vouchers           Vouchers           `json:"vouchers,omitempty"`
		OutletID           string             `db:"outlet_id" json:"outlet_id"`
		OutletName         string             `db:"outlet_name" json:"outlet_name"`
		OutletDescription  types.JSONText     `db:"outlet_description" json:"outlet_description,omitempty"`
		CreatedBy          string             `db:"created_by" json:"created_by"`
		CreatedAt          *time.Time         `db:"created_at" json:"created_at"`
		UpdatedBy          string             `db:"updated_by" json:"updated_by"`
		UpdatedAt          *time.Time         `db:"updated_at" json:"updated_at"`
		Status             string             `db:"status" json:"status"`
		TransactionDetails TransactionDetails `json:"-"`
		Count              int                `db:"count" json:"-"`
	}
	Transactions      []Transaction
	TransactionDetail struct {
		ID            int        `db:"id" json:"id"`
		TransactionID string     `db:"transaction_id" json:"transaction_id"`
		ProgramID     string     `db:"program_id" json:"program_id"`
		VoucherID     string     `db:"voucher_id" json:"voucher_id"`
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

//GetTransactionByID : get outlet by specified ID
func GetTransactionByID(qp *util.QueryParam, id string) (Transaction, error) {
	// return
	transactions, _, err := getTransactions("id", id, qp)
	if err != nil {
		return Transaction{}, err
	}

	vouchers, err := GetTransactionVouchers(id, qp)
	if err != nil {
		return Transaction{}, ErrorResourceNotFound
	}

	if len(*transactions) > 0 {
		transaction := (*transactions)[0]
		transaction.Vouchers = vouchers
		return transaction, nil
	}

	return Transaction{}, ErrorResourceNotFound
}

//GetTransactionByProgram : get transaction by specified ProgramID
func GetTransactionByProgram(qp *util.QueryParam, val string) (*Transactions, bool, error) {
	return getTransactions("program_id", val, qp)
}

//GetTransactionByOutlet : get transaction by specified OutletID
func GetTransactionByOutlet(qp *util.QueryParam, val string) (*Transactions, bool, error) {
	return getTransactions("outlet_id", val, qp)
}

//GetTransactionDetailByVoucherID : get transaction detail by voucher ID
func GetTransactionDetailByVoucherID(qp *util.QueryParam, voucherID string) (*[]TransactionDetail, error) {
	q, err := qp.GetQueryByDefaultStruct(TransactionDetail{})
	if err != nil {
		return &[]TransactionDetail{}, err
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
		return &[]TransactionDetail{}, err
	}

	if len(resd) < 1 {
		return &[]TransactionDetail{}, ErrorResourceNotFound
	}

	return &resd, nil
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

	q = qp.GetQueryWhereClause(q, qp.Q)
	q = qp.GetQueryWithPagination(q, qp.GetQuerySort(), qp.GetQueryLimit())
	var resd Transactions
	err = db.Select(&resd, db.Rebind(q), StatusCreated, v)
	if err != nil {
		return &Transactions{}, false, err
	}
	if len(resd) < 1 {
		return &Transactions{}, false, nil
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

// GetTransactionVouchers : Get list of vouchers by transaction ID
func GetTransactionVouchers(transactionID string, qp *util.QueryParam) (Vouchers, error) {

	q, err := qp.GetQueryByDefaultStruct(Voucher{})
	if err != nil {
		return Vouchers{}, err
	}

	q += `
			FROM
				m_vouchers voucher
			JOIN transaction_details td ON td.voucher_id = voucher.id
			JOIN transactions t ON t.id = td.transaction_id
			WHERE 
				t.status = ?
			AND voucher.state != ?
			AND t.id = ?`

	q = qp.GetQueryWhereClause(q, qp.Q)
	var resd Vouchers
	err = db.Select(&resd, db.Rebind(q), StatusCreated, StatusDeleted, transactionID)
	if err != nil {
		return Vouchers{}, err
	}
	if len(resd) < 1 {
		return Vouchers{}, nil
	}

	return resd, nil
}

//Insert : transaction data
func (t Transaction) Insert() (*[]Transaction, error) {
	tx, err := db.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	q := `INSERT INTO 
			transactions (company_id, transaction_code, total_amount, holder, outlet_id, created_by, updated_by, status)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)
			RETURNING
				id, company_id, transaction_code, total_amount, holder, outlet_id, created_by, created_at, updated_by, updated_at, status
	`
	var res []Transaction
	err = tx.Select(&res, tx.Rebind(q), t.CompanyID, t.TransactionCode, t.TotalAmount, t.Holder, t.OutletID, t.CreatedBy, t.UpdatedBy, StatusCreated)
	if err != nil {
		return nil, err
	}

	var resd []TransactionDetail
	for _, td := range t.TransactionDetails {
		td.TransactionID = res[0].ID
		q = `INSERT INTO 
			transaction_details (transaction_id, program_id, voucher_id, created_by, updated_by, status)
			VALUES (?, ?, ?, ?, ?, ?)
			RETURNING
				id, transaction_id, program_id, voucher_id, created_by, created_at, updated_by, updated_at, status
	`
		err = tx.Select(&resd, tx.Rebind(q), td.TransactionID, td.ProgramID, td.VoucherID, td.CreatedBy, td.UpdatedBy, StatusCreated)
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
//			transactions (company_id, transaction_code, total_amount, holder, outlet_id, created_by, updated_by, status)
//			VALUES (?, ?, ?, ?, ?, ?, ?, ?)
//			RETURNING
//				id, company_id, transaction_code, total_amount, holder, outlet_id, created_by, created_at, updated_by, updated_at, status
//	`
//	var res []TransactionDetail
//	err = tx.Select(&res, tx.Rebind(q), t.CompanyId, t.TransactionCode, t.TotalAmount, t.Holder, t.OutletId, t.CreatedBy, t.UpdatedBy, StatusCreated)
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
