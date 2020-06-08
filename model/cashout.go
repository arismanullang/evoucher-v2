package model

import (
	"fmt"
	"time"

	"github.com/gilkor/evoucher-v2/util"
)

type (
	CashoutSummary struct {
		Date          *time.Time `db:"date" json:"date,omitempty"`
		PartnerID     string     `db:"partner_id" json:"partner_id,omitempty"`
		UnpaidAmount  float64    `db:"unpaid_amount" json:"unpaid_amount,omitempty"`
		CashoutAmount float64    `db:"cashout_amount" json:"cashout_amount,omitempty"`
		VoucherQty    int        `db:"voucher_qty" json:"voucher_aty,omitempty"`
		Count         int        `db:"count" json:"-"`
	}
	//Cashout : represent of cashout table model
	Cashout struct {
		ID             string         `db:"id" json:"id,omitempty"`
		AccountID      string         `db:"account_id" json:"account_id,omitempty"`
		Code           string         `db:"code" json:"code,omitempty"`
		PartnerID      string         `db:"partner_id" json:"partner_id,omitempty"`
		BankAccount    string         `db:"bank_account" json:"bank_account,omitempty"`
		Amount         float64        `db:"amount" json:"amount,omitempty"`
		PaymentMethod  string         `db:"payment_method" json:"payment_method,omitempty"`
		CreatedAt      *time.Time     `db:"created_at" json:"created_at,omitempty"`
		CreatedBy      string         `db:"created_by" json:"created_by,omitempty"`
		UpdatedAt      *time.Time     `db:"updated_at" json:"updated_at,omitempty"`
		UpdatedBy      string         `db:"updated_by" json:"updated_by,omitempty"`
		Status         string         `db:"status" json:"status,omitempty"`
		CashoutDetails CashoutDetails `json:"cashout_details,omitempty"`
		Count          int            `db:"count" json:"-"`
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
				v_cashout_summary
			WHERE 
				status = ?			
			AND ` + key + ` = ?`

	q += qp.GetQuerySort()
	q += qp.GetQueryLimit()
	fmt.Println(q)
	var resd []CashoutSummary
	err = db.Select(&resd, db.Rebind(q), StatusCreated, value)
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

	q := `INSERT INTO 
				cashouts 
				( 
					account_id
					, code 
					, partner_id 
					, bank_account
					, amount 
					, payment_method
					, created_by
					, updated_by
					, status
				)
			VALUES 
				( ?, ?, ?, ?, ?, ?, ?, ?)
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
	err = tx.Select(&res, tx.Rebind(q), c.AccountID, c.Code, c.PartnerID, c.BankAccount, c.Amount, c.PaymentMethod, c.CreatedBy, c.UpdatedBy, StatusCreated)
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
		err = tx.Select(&res, tx.Rebind(q), cd.CashoutID, cd.VoucherID, cd.CreatedBy, cd.UpdatedBy, StatusCreated)
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
