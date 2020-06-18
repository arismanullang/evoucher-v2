package model

import (
	"time"

	"github.com/gilkor/evoucher-v2/util"
)

type (
	DashboardVoucherUsage struct {
		Year   int64   `json:"year,omitempty" db:"year"`
		Month  int64   `json:"month,omitempty" db:"month"`
		Amount float64 `json:"amount,omitempty" db:"amount"`
		Count  int     `db:"count" json:"-"`
	}
	DashboardTopProgram struct {
		ProgramID         string     `db:"program_id" json:"program_id,omitempty"`
		ProgramName       string     `db:"program_name" json:"program_name,omitempty"`
		ProgramdStartDate *time.Time `db:"program_start_date" json:"program_start_date,omitempty"`
		ProgramdEndDate   *time.Time `db:"program_end_date" json:"program_end_date,omitempty"`
		VoucherClaim      int64      `db:"voucher_claim" json:"voucher_claim,omitempty"`
		VoucherUsed       int64      `db:"voucher_used" json:"voucher_used,omitempty"`
		Count             int        `db:"count" json:"-"`
	}
	DashboardTopOutlet struct {
		TransactionDate *time.Time `db:"transaction_date" json:"transaction_date,omitempty"`
		OutletID        string     `db:"outlet_id" json:"outlet_id,omitempty"`
		OutletName      string     `db:"outlet_name" json:"outlet_name,omitempty"`
		VoucherQty      int64      `db:"voucher_qty" json:"voucher_qty,omitempty"`
		TransactionQty  int64      `db:"transaction_qty" json:"transaction_qty,omitempty"`
		TotalAmount     float64    `db:"total_amount" json:"total_amount,omitempty"`
		Count           int        `db:"count" json:"-"`
	}
)

func GetDashboardVoucherUsage(dateFrom, dateTo string, qp *util.QueryParam) ([]DashboardVoucherUsage, bool, error) {
	q := `
			SELECT 
				EXTRACT(YEAR FROM created_at) as year, EXTRACT(MONTH FROM created_at) as month, sum(transactions.discount_value) as amount 
			FROM transactions
			WHERE 
				status = ?  
				AND created_at BETWEEN '` + dateFrom + ` 00:00:00+07'::timestamp AND '` + dateTo + ` 23:59:59+07'::timestamp `

	//q = qp.GetQueryWhereClause(q, qp.Q)

	q += ` GROUP BY year, month `

	q = qp.GetQueryWithPagination(q, qp.GetQuerySort(), qp.GetQueryLimit())
	//ORDER BY year, month `
	util.DEBUG(q)
	var resd []DashboardVoucherUsage
	err := db.Select(&resd, db.Rebind(q), StatusCreated)
	if err != nil {
		return []DashboardVoucherUsage{}, false, err
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

func GetDashboardTopProgram(dateFrom, dateTo string, qp *util.QueryParam) ([]DashboardTopProgram, bool, error) {
	q := `
			SELECT
				p.id, p.name, p.created_at, p.end_date, sum(case
				when v.id is not null then 1
				else 0 end)
					as claim,
				   sum(case
				when t.id is not null then 1
				else 0 end) as used
			FROM
				programs p
				LEFT JOIN vouchers v ON p.id = v.program_id AND v.status = 'created'
				LEFT JOIN transaction_details td ON td.voucher_id = v.id AND p.status = 'created'
				LEFT JOIN transactions t ON t.id = td.transaction_id AND t.status = 'created'
			WHERE p.status = 'created'
				AND t.created_at BETWEEN '` + dateFrom + ` 00:00:00+07'::timestamp AND '` + dateTo + ` 23:59:59+07'::timestamp `

	q += ` GROUP BY p.id, p.name `
	q = qp.GetQueryWithPagination(q, qp.GetQuerySort(), qp.GetQueryLimit())
	//ORDER BY year, month `
	util.DEBUG(q)
	var resd []DashboardTopProgram
	err := db.Select(&resd, db.Rebind(q), StatusCreated)
	if err != nil {
		return []DashboardTopProgram{}, false, err
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

func GetDashboardTopOutlet(dateFrom, dateTo string, qp *util.QueryParam) ([]DashboardTopOutlet, bool, error) {
	q := `
			SELECT date as transaction_date, id as outlet_id, name as outlet_name, sum(trans) as transaction_qty, sum(vouc) as voucher_qty, sum(amount) as total_amount
			FROM (
					 SELECT date(t.created_at) as date, p.id, p.name, sum(1) as trans, 0 as vouc, sum(t.discount_value) as amount
					 FROM transactions t,
						  outlets p
					 WHERE t.outlet_id = p.id
					   AND t.status = 'created'
					   AND p.status = 'created'
					   AND t.created_at BETWEEN '` + dateFrom + ` 00:00:00+07'::timestamp AND '` + dateTo + ` 23:59:59+07'::timestamp 
					 GROUP BY date, p.id, p.name
					 UNION ALL
					 SELECT date(t.created_at) as date, p.id, p.name, 0 as trans, sum(1) as vouc, 0 as amount
					 FROM transactions t,
						  transaction_details td,
						  outlets p
					 WHERE t.outlet_id = p.id
					   AND t.id = td.transaction_id
					   AND t.status = 'created'
					   AND td.status = 'created'
					   AND p.status = 'created'
					   AND t.created_at BETWEEN '` + dateFrom + ` 00:00:00+07'::timestamp AND '` + dateTo + ` 23:59:59+07'::timestamp 
					 GROUP BY date, p.id, p.name
				 ) as x
			GROUP BY date, id, name `

	q = qp.GetQueryWithPagination(q, qp.GetQuerySort(), qp.GetQueryLimit())
	//ORDER BY year, month `
	util.DEBUG(q)
	var resd []DashboardTopOutlet
	err := db.Select(&resd, db.Rebind(q), StatusCreated)
	if err != nil {
		return []DashboardTopOutlet{}, false, err
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
