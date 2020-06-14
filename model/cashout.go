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
		PartnerID     string     `db:"partner_id" json:"partner_id,omitempty"`
		PartnerName   string     `db:"partner_name" json:"partner_name,omitempty"`
		UnpaidAmount  float64    `db:"unpaid_amount" json:"unpaid_amount,omitempty"`
		CashoutAmount float64    `db:"cashout_amount" json:"cashout_amount,omitempty"`
		TotalAmount   float64    `db:"total_amount" json:"total_amount,omitempty"`
		VoucherQty    int        `db:"voucher_qty" json:"voucher_aty,omitempty"`
		Count         int        `db:"count" json:"-"`
	}
	CashoutUnpaid struct {
		Date           *time.Time `db:"date" json:"date,omitempty"`
		PartnerID      string     `db:"partner_id" json:"partner_id,omitempty"`
		PartnerName    string     `db:"partner_name" json:"partner_name,omitempty"`
		TransactionQty int64      `db:"transaction_qty" json:"transaction_qty,omitempty"`
		TotalValue     float64    `db:"total_value" json:"total_value,omitempty"`
		Count          int        `db:"count" json:"-"`
	}

	// UnpaidReimburse : unpaid cashout grouped by partner
	UnpaidReimburse struct {
		PartnerID          string         `db:"partner_id" json:"partner_id,omitempty"`
		PartnerName        string         `db:"partner_name" json:"partner_name,omitempty"`
		PartnerDescription types.JSONText `db:"partner_description" json:"partner_description,omitempty"`
		PartnerEmails      *string        `db:"partner_emails" json:"partner_emails,omitempty"`
		CompanyID          string         `db:"company_id" json:"company_id"`
		TransactionQty     int64          `db:"transaction_qty" json:"transaction_qty,omitempty"`
		VouchersQty        int64          `db:"vouchers_qty" json:"vouchers_qty,omitempty"`
		TotalValue         float64        `db:"total_value" json:"total_value,omitempty"`
		Count              int            `db:"count" json:"-"`
	}

	// UnpaidVouchers : unpaid vouchers
	UnpaidVouchers struct {
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
		PartnerID       string         `db:"partner_id" json:"partner_id,omitempty"`
		BankName        string         `db:"bank_name" json:"bank_name,omitempty"`
		BankCompanyName string         `db:"bank_company_name" json:"bank_company_name"`
		BankAccount     string         `db:"bank_account" json:"bank_account,omitempty"`
		ReferenceNo     string         `db:"reference_no" json:"reference_no,omitempty"`
		Amount          float64        `db:"amount" json:"amount,omitempty"`
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
func GetCashoutByID(id string, qp *util.QueryParam) (*Cashouts, bool, error) {
	return getCashouts(qp, "id", id)
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
				cashout
			WHERE 
				status = ?			
			AND ` + key + ` = ?`

	q += qp.GetQuerySort()
	q += qp.GetQueryLimit()
	fmt.Println(q)
	var resd Cashouts
	err = db.Select(&resd, db.Rebind(q), StatusCreated, value)
	if err != nil {
		return &Cashouts{}, false, err
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

// GetCashoutUnpaid : list Cashout Unpaid
func GetCashoutUnpaid(qp *util.QueryParam) ([]CashoutUnpaid, bool, error) {
	return getCashoutUnpaid(qp, "1", "1")
}

func getCashoutUnpaid(qp *util.QueryParam, key, value string) ([]CashoutUnpaid, bool, error) {
	q, err := qp.GetQueryByDefaultStruct(CashoutSummary{})
	if err != nil {
		return []CashoutUnpaid{}, false, err
	}
	q += `
			FROM
				m_cashout_unpaid
			WHERE 		
			 ` + key + ` = ? `

	q += qp.GetQuerySort()
	q += qp.GetQueryLimit()
	fmt.Println(q)
	var resd []CashoutUnpaid
	err = db.Select(&resd, db.Rebind(q), value)
	if err != nil {
		return []CashoutUnpaid{}, false, err
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
					, partner_id 
					, bank_name
					, bank_company_name
					, bank_account
					, reference_no
					, amount 
					, payment_method
					, created_by
					, updated_at
					, updated_by
					, status
				)
			VALUES 
				( ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			RETURNING
				id
				, company_id
				, code 
				, partner_id 
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
	err = tx.Select(&res, tx.Rebind(q), c.CompanyID, c.Code, c.PartnerID, c.BankName, c.BankCompanyName, c.BankAccount, c.ReferenceNo, c.Amount, c.PaymentMethod, c.CreatedBy, c.UpdatedBy, StatusCreated)
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
				, partner_id 
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

// GetUnpaidReimburse : Get list of unpaid reimburse
func GetUnpaidReimburse(qp *util.QueryParam, startDate, endDate string) ([]UnpaidReimburse, bool, error) {

	q := `
		SELECT DISTINCT
			p.id AS partner_id,
			p.name AS partner_name,
			p.description AS partner_description,
			p.emails AS partner_emails,
			p.company_id,
			count(t.id) AS transaction_qty,
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
		FROM partners p
			LEFT JOIN transactions t ON p.id = t.partner_id
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

	var resv []UnpaidReimburse
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
func GetUnpaidVouchersByOutlet(qp *util.QueryParam, partnerID, startDate, endDate string) ([]UnpaidVouchers, bool, error) {

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
			JOIN partners as p ON t.partner_id = p.id
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

	var resv []UnpaidVouchers
	if err := db.Select(&resv, db.Rebind(q), qp.CompanyID, partnerID, startDate, endDate); err != nil {
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
