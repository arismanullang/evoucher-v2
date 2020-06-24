package model

import (
	"fmt"
	"time"

	"github.com/gilkor/evoucher-v2/util"
	"github.com/jmoiron/sqlx/types"
)

type (
	CashoutSummary struct {
		Date          *time.Time `db:"date" json:"date,omitempty"`
		OutletID      string     `db:"outlet_id" json:"outlet_id,omitempty"`
		OutletName    string     `db:"outlet_name" json:"outlet_name,omitempty"`
		UnpaidAmount  float64    `db:"unpaid_amount" json:"unpaid_amount,omitempty"`
		CashoutAmount float64    `db:"cashout_amount" json:"cashout_amount,omitempty"`
		TotalAmount   float64    `db:"total_amount" json:"total_amount,omitempty"`
		VoucherQty    int        `db:"voucher_qty" json:"voucher_aty,omitempty"`
		Count         int        `db:"count" json:"-"`
	}

	// UnpaidCashout : unpaid cashout grouped by outlet
	UnpaidCashout struct {
		OutletID          string         `db:"outlet_id" json:"outlet_id,omitempty"`
		OutletName        string         `db:"outlet_name" json:"outlet_name,omitempty"`
		OutletDescription types.JSONText `db:"outlet_description" json:"outlet_description,omitempty"`
		OutletEmails      *string        `db:"outlet_emails" json:"outlet_emails,omitempty"`
		CompanyID         string         `db:"company_id" json:"company_id"`
		TransactionQty    int64          `db:"transaction_qty" json:"transaction_qty,omitempty"`
		VouchersQty       int64          `db:"vouchers_qty" json:"vouchers_qty,omitempty"`
		TotalValue        float64        `db:"total_value" json:"total_value,omitempty"`
		Count             int            `db:"count" json:"-"`
	}

	// VoucherTransaction : voucher with transaction detail
	VoucherTransaction struct {
		VoucherID       string    `db:"voucher_id" json:"voucher_id,omitempty"`
		VoucherCode     string    `db:"voucher_code" json:"voucher_code,omitempty"`
		TransactionID   string    `db:"transaction_id" json:"transaction_id,omitempty"`
		TransactionCode string    `db:"transaction_code" json:"transaction_code,omitempty"`
		ClaimedAt       time.Time `db:"claimed_at" json:"claimed_at,omitempty"`
		UsedAt          time.Time `db:"used_at" json:"used_at,omitempty"`
		ProgramName     string    `db:"program_name" json:"program_name,omitempty"`
		ProgramValue    float64   `db:"program_value" json:"program_value,omitempty"`
		ProgramMaxValue float64   `db:"program_max_value" json:"program_max_value,omitempty"`
		Count           int       `db:"count" json:"-"`
	}

	//Cashout : represent of cashout table model
	Cashout struct {
		ID              string         `db:"id" json:"id,omitempty"`
		CompanyID       string         `db:"company_id" json:"company_id,omitempty"`
		Code            string         `db:"code" json:"code,omitempty"`
		OutletName      string         `db:"outlet_name" json:"outlet_name,omitempty"`
		OutletID        string         `db:"outlet_id" json:"outlet_id,omitempty"`
		BankName        string         `db:"bank_name" json:"bank_name,omitempty"`
		BankCompanyName string         `db:"bank_company_name" json:"bank_company_name"`
		BankAccount     string         `db:"bank_account" json:"bank_account,omitempty"`
		ReferenceNo     string         `db:"reference_no" json:"reference_no,omitempty"`
		Amount          float64        `db:"amount" json:"amount,omitempty"`
		AttachmentUrl   string         `db:"attachment_url" json:"attachment_url,omitempty"`
		PaymentMethod   string         `db:"payment_method" json:"payment_method,omitempty"`
		CreatedAt       *time.Time     `db:"created_at" json:"created_at,omitempty"`
		CreatedBy       string         `db:"created_by" json:"created_by,omitempty"`
		UpdatedAt       *time.Time     `db:"updated_at" json:"updated_at,omitempty"`
		UpdatedBy       string         `db:"updated_by" json:"updated_by,omitempty"`
		Status          string         `db:"status" json:"status,omitempty"`
		CashoutDetails  CashoutDetails `json:"cashout_details,omitempty"`
		Count           int            `db:"count" json:"-"`
	}
	Cashouts      []Cashout
	CashoutDetail struct {
		ID        int        `db:"id" json:"id,omitempty"`
		CashoutID string     `db:"cashout_id" json:"cashout_id,omitempty"`
		VoucherID string     `db:"voucher_id" json:"voucher_id,omitempty"`
		CreatedAt *time.Time `db:"created_at" json:"created_at,omitempty"`
		CreatedBy string     `db:"created_by" json:"created_by,omitempty"`
		UpdatedAt *time.Time `db:"updated_at" json:"updated_at,omitempty"`
		UpdatedBy string     `db:"updated_by" json:"updated_by,omitempty"`
		Status    string     `db:"status" json:"status,omitempty"`
	}
	CashoutDetails []CashoutDetail
)

// GetCashoutByID :  get list Cashouts by ID
func GetCashoutByID(id string, qp *util.QueryParam) (*Cashout, error) {
	cashouts, _, err := getCashouts(qp, "id", id)
	if err != nil {
		return &Cashout{}, err
	}

	if len(*cashouts) > 0 {
		cashout := (*cashouts)[0]
		return &cashout, nil
	}

	return &Cashout{}, ErrorResourceNotFound
}

// GetCashouts : list Cashout
func GetCashouts(qp *util.QueryParam) (*Cashouts, bool, error) {
	return getCashouts(qp, "1", "1")
}

func getCashouts(qp *util.QueryParam, key, value string) (*Cashouts, bool, error) {
	q, err := qp.GetQueryByDefaultStruct(Cashout{})
	if err != nil {
		return &Cashouts{}, false, err
	}
	q += `
			FROM
				m_cashouts cashout
			WHERE ` + key + ` = ?`

	q = qp.GetQueryWhereClause(q, qp.Q)
	q = qp.GetQueryWithPagination(q, qp.GetQuerySort(), qp.GetQueryLimit())

	fmt.Println(q)
	var resd Cashouts
	err = db.Select(&resd, db.Rebind(q), value)
	if err != nil {
		return &Cashouts{}, false, err
	}

	if len(resd) < 1 {
		return &Cashouts{}, false, nil
	}

	next := false
	if len(resd) > qp.Count {
		next = true
		resd = resd[:qp.Count]
	}
	if len(resd) < qp.Count {
		qp.Count = len(resd)
	}

	return &resd, next, nil
}

// GetCashouts : list Cashout Summary
func GetCashoutSummary(qp *util.QueryParam) ([]CashoutSummary, bool, error) {
	return getCashoutSummary(qp, "1", "1")
}

func getCashoutSummary(qp *util.QueryParam, key, value string) ([]CashoutSummary, bool, error) {
	q, err := qp.GetQueryByDefaultStruct(CashoutSummary{})
	if err != nil {
		return []CashoutSummary{}, false, err
	}
	q += `
			FROM
				m_cashout_summary
			WHERE 		
			 ` + key + ` = ? `

	q += qp.GetQuerySort()
	q += qp.GetQueryLimit()
	fmt.Println(q)
	var resd []CashoutSummary
	err = db.Select(&resd, db.Rebind(q), value)
	if err != nil {
		return []CashoutSummary{}, false, err
	}

	next := false
	if len(resd) > qp.Count {
		next = true
	}
	if len(resd) < qp.Count {
		qp.Count = len(resd)
	}

	return resd, next, nil
}

//Insert : single row inset into table
func (c *Cashout) Insert() (*[]Cashout, error) {
	tx, err := db.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	q := `INSERT INTO cashouts ( 
					company_id
					, code 
					, outlet_id 
					, bank_name
					, bank_company_name
					, bank_account
					, reference_no
					, amount 
					, payment_method
					, created_by
					, updated_by
					, status
				)
			VALUES 
				( ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			RETURNING
				id
				, company_id
				, code 
				, outlet_id 
				, bank_name
				, bank_company_name
				, bank_account
				, reference_no
				, amount 
				, payment_method
				, created_by
				, updated_by
				, status
	`

	var res []Cashout
	err = tx.Select(&res, tx.Rebind(q), c.CompanyID, c.Code, c.OutletID, c.BankName, c.BankCompanyName, c.BankAccount, c.ReferenceNo, c.Amount, c.PaymentMethod, c.CreatedBy, c.UpdatedBy, StatusCreated)
	if err != nil {
		return nil, err
	}

	var resd []CashoutDetail
	for _, cd := range c.CashoutDetails {
		q = `INSERT INTO 
			cashout_details (cashout_id, voucher_id, created_by, updated_by, status)
			VALUES (?, ?, ?, ?, ?)
			RETURNING
				id, cashout_id, voucher_id, created_by, updated_by, status
	`
		err = tx.Select(&resd, tx.Rebind(q), res[0].ID, cd.VoucherID, cd.CreatedBy, cd.UpdatedBy, StatusCreated)
		if err != nil {
			return nil, err
		}
	}
	res[0].CashoutDetails = resd

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &res, nil
}

//Delete : soft delated data by updateting row status to "deleted"
func (c *Cashout) Delete() error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	q := `UPDATE
				cashout 
			SET
				deleted_at = now(),
				deleted_by = ?
				status = ?			
			WHERE 
				id = ?	
			RETURNING
				id
				, account_id
				, code 
				, outlet_id 
				, bank_account
				, amount 
				, payment_method
				, created_by
				, updated_by
				, status
	`
	var res []Cashout
	err = tx.Select(&res, tx.Rebind(q), c.UpdatedBy, StatusDeleted, c.ID)

	q = `UPDATE
				cashouts_details
			SET
				deleted_at = now(),
				deleted_by = ?
				status = ?			
			WHERE 
				cashout_id = ?	
			RETURNING
				id
	`
	var resd []CashoutDetail
	err = tx.Select(&res, tx.Rebind(q), c.UpdatedBy, StatusDeleted, c.ID)

	res[0].CashoutDetails = resd
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// GetUnpaidCashout : Get list of unpaid cashout
func GetUnpaidCashout(qp *util.QueryParam, startDate, endDate string) ([]UnpaidCashout, bool, error) {

	q := `
		SELECT DISTINCT
			p.id AS outlet_id,
			p.name AS outlet_name,
			p.description AS outlet_description,
			p.emails AS outlet_emails,
			p.company_id,
			count(DISTINCT t.id) AS transaction_qty,
			COALESCE(sum(
				CASE
					WHEN v.state = 'used'::voucher_state THEN 1
					ELSE 0
				END), 0::bigint) AS vouchers_qty,
			COALESCE(sum(
				CASE
					WHEN v.state = 'used'::voucher_state THEN 1 * pr.max_value
					ELSE 0
				END), 0::bigint) AS total_value
		FROM outlets p
			LEFT JOIN transactions t ON p.id = t.outlet_id
			LEFT JOIN transaction_details td ON t.id = td.transaction_id
			LEFT JOIN vouchers v ON v.id = td.voucher_id
			LEFT JOIN programs pr ON pr.id = v.program_id
			WHERE v.state = 'used'
			AND pr.is_reimburse = true
			AND p.company_id = ?
			AND t.created_at BETWEEN ? AND ?
		GROUP BY p.id`

	q = qp.GetQueryWithPagination(q, qp.GetQuerySort(), qp.GetQueryLimit())

	util.DEBUG("query struct :", q)

	var resv []UnpaidCashout
	if err := db.Select(&resv, db.Rebind(q), qp.CompanyID, startDate, endDate); err != nil {
		fmt.Println(err.Error())
		return resv, false, err
	}

	if len(resv) < 1 {
		return resv, false, ErrorResourceNotFound
	}

	next := false
	if len(resv) > qp.Count {
		next = true
		resv = resv[:qp.Count]
	}
	if len(resv) < qp.Count {
		qp.Count = len(resv)
	}
	return resv, next, nil
}

// GetUnpaidVouchersByOutlet : Get list of unpaid vouchers by outlet transaction
func GetUnpaidVouchersByOutlet(qp *util.QueryParam, outletID, startDate, endDate string) ([]VoucherTransaction, bool, error) {

	q := `
		SELECT DISTINCT
			v.id as voucher_id
			, v.code as voucher_code
			, t.id as transaction_id
			, t.transaction_code as transaction_code
			, v.created_at as claimed_at
			, t.created_at as used_at
			, pr.name as program_name
			, pr.value as program_value
			, pr.max_value as program_max_value
		FROM vouchers as v
			JOIN transaction_details as td ON td.voucher_id = v.id
			JOIN transactions as t ON td.transaction_id = t.id
			JOIN outlets as p ON t.outlet_id = p.id
			JOIN programs as pr ON v.program_id = pr.id
		WHERE
			v.state = 'used'
			AND pr.company_id = ?
			AND pr.is_reimburse = true
			AND p.id = ?
			AND t.created_at BETWEEN ? AND ?
`

	q = qp.GetQueryWithPagination(q, qp.GetQuerySort(), qp.GetQueryLimit())

	util.DEBUG("query struct :", q)

	var resv []VoucherTransaction
	if err := db.Select(&resv, db.Rebind(q), qp.CompanyID, outletID, startDate, endDate); err != nil {
		fmt.Println(err.Error())
		return resv, false, err
	}

	if len(resv) < 1 {
		return resv, false, ErrorResourceNotFound
	}

	next := false
	if len(resv) > qp.Count {
		next = true
		resv = resv[:qp.Count]
	}
	if len(resv) < qp.Count {
		qp.Count = len(resv)
	}
	return resv, next, nil
}

// GetCashoutVouchers : Get list of paid vouchers by cashoutID
func GetCashoutVouchers(qp *util.QueryParam, cashoutID string) ([]VoucherTransaction, bool, error) {

	q := `
		SELECT
			v.id as voucher_id
			, v.code as voucher_code
			, p.name as program_name
			, p.value as program_value
			, p.max_value as program_max_value
			, v.created_at as claimed_at
			, t.created_at as used_at
			, t.id as transaction_id
			, t.transaction_code as transaction_code
		FROM vouchers v
			JOIN transaction_details td ON v.id = td.voucher_id
			JOIN transactions t ON t.id = td.transaction_id
			JOIN programs p ON v.program_id = p.id
			JOIN cashout_details cd ON cd.voucher_id = v.id
			JOIN cashouts c ON c.id = cd.cashout_id
		WHERE v.state = 'paid'
			AND p.company_id = ?
			AND c.id = ?
			`

	q = qp.GetQueryWithPagination(q, qp.GetQuerySort(), qp.GetQueryLimit())

	util.DEBUG("query struct :", q)

	var resv []VoucherTransaction
	if err := db.Select(&resv, db.Rebind(q), qp.CompanyID, cashoutID); err != nil {
		fmt.Println(err.Error())
		return resv, false, err
	}

	if len(resv) < 1 {
		return resv, false, ErrorResourceNotFound
	}

	next := false
	if len(resv) > qp.Count {
		next = true
		resv = resv[:qp.Count]
	}
	if len(resv) < qp.Count {
		qp.Count = len(resv)
	}
	return resv, next, nil
}

func (c *Cashout) Update() (*Cashout, error) {
	tx, err := db.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	q := `UPDATE
				cashouts 
			SET
				attachment_url = ?,
				updated_at = now(),
				updated_by = ?,					
				status = ?
			WHERE 
				id = ?
			RETURNING
				id
				, code
				, bank_name
				, bank_company_name
				, bank_account
				, reference_no
				, amount
				, attachment_url
				, payment_method
				, created_at
				, created_by
				, updated_at
				, updated_by
				, status
				
	`
	var res []Cashout
	err = tx.Select(&res, tx.Rebind(q), c.AttachmentUrl, c.UpdatedBy, c.Status, c.ID)
	if err != nil {
		return nil, err
	}

	if c.Status == StatusApproved {
		//insert approval_log
		approvalLog := ApprovalLog{
			ObjectID:       c.ID,
			ObjectCategory: "cashout",
			CreatedBy:      c.UpdatedBy,
			Status:         StatusCreated,
		}

		fmt.Println("approval log = ", approvalLog)
		err = approvalLog.Insert()
		if err != nil {
			return nil, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &res[0], nil
}
