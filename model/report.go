package model

import (
	"time"

	"github.com/gilkor/evoucher-v2/util"
)

type ReportOutletTransaction struct {
	TransactionDate  *time.Time `db:"transaction_date" json:"transaction_date,omitempty"`
	TransactionMonth int        `db:"transaction_month" json:"transaction_month,omitempty"`
	TransactionYear  int        `db:"transaction_year" json:"transaction_year,omitempty"`
	OutletID         string     `db:"outlet_id" json:"outlet_id,omitempty"`
	OutletName       string     `db:"outlet_name" json:"outlet_name,omitempty"`
	VoucherQty       *int64     `db:"voucher_qty" json:"voucher_qty,omitempty"`
	ReimburseQty     *int64     `db:"reimburse_qty" json:"reimburse_qty,omitempty"`
	TransactionQty   *int64     `db:"transaction_qty" json:"transaction_qty,omitempty"`
	TotalReimburse   *float64   `db:"total_reimburse" json:"total_reimburse,omitempty"`
	TotalTransaction *float64   `db:"total_transaction" json:"total_transaction,omitempty"`
	Count            int        `db:"count" json:"-"`
}

type ReportProgramTransaction struct {
	Date             *time.Time `db:"date" json:"date,omitempty"`
	Month            int        `db:"month" json:"month,omitempty"`
	Year             int        `db:"year" json:"year,omitempty"`
	ProgramID        string     `db:"id" json:"id,omitempty"`
	ProgramName      string     `db:"name" json:"name,omitempty"`
	VoucherQty       *int64     `db:"voucher_qty" json:"voucher_qty,omitempty"`
	ReimburseQty     *int64     `db:"reimburse_qty" json:"reimburse_qty,omitempty"`
	TransactionQty   *int64     `db:"transaction_qty" json:"transaction_qty,omitempty"`
	TotalVoucher     *float64   `db:"total_voucher" json:"total_voucher,omitempty"`
	TotalReimburse   *float64   `db:"total_reimburse" json:"total_reimburse,omitempty"`
	TotalTransaction *float64   `db:"total_transaction" json:"total_transaction,omitempty"`
	Count            int        `db:"count" json:"-"`
}

type ReportReimburse struct {
	Date             *time.Time `db:"date" json:"date,omitempty"`
	Month            int        `db:"month" json:"month,omitempty"`
	Year             int        `db:"year" json:"year,omitempty"`
	OutletID         string     `db:"outlet_idx" json:"outlet_id,omitempty"`
	OutletName       string     `db:"outlet_name" json:"outlet_name,omitempty"`
	ProgramID        string     `db:"program_id" json:"program_id,omitempty"`
	ProgramName      string     `db:"program_name" json:"program_name,omitempty"`
	VoucherQty       *int64     `db:"voucher_qty" json:"voucher_qty,omitempty"`
	TransactionQty   *int64     `db:"transaction_qty" json:"transaction_qty,omitempty"`
	ReimburseQty     *int64     `db:"reimburse_qty" json:"reimburse_qty,omitempty"`
	TotalTransaction float64    `db:"total_transaction" json:"total_transaction,omitempty"`
	TotalReimburse   float64    `db:"total_reimburse" json:"total_reimburse,omitempty"`
	TotalVoucher     float64    `db:"total_voucher" json:"total_voucher,omitempty"`
	Count            int        `db:"count" json:"-"`
}

func GetReportDailyVoucherTransaction(dateFrom, dateTo string, qp *util.QueryParam) ([]ReportReimburse, bool, error) {
	q := `
			SELECT date(v.created_at) as date,
				   sum(case
						   when t.id is not null then 1
						   else 0 end)
									  as transaction_qty,
				   sum(case
						   when c.id is not null then 1
						   else 0 end)
									  as reimburse_qty,
				   sum(case
						   when v.id is not null then 1
						   else 0 end)
									  as voucher_qty,
				   sum(case
						   when t.id is not null then p.value
						   else 0 end)
									  as total_transaction,
				   sum(case
						   when c.id is not null then p.value
						   else 0 end)
									  as total_reimburse,
				   sum(case
						   when v.id is not null then p.value
						   else 0 end)
									  as total_voucher
			FROM programs p,
				 vouchers v
					 LEFT JOIN transaction_details td on td.voucher_id = v.id AND td.status = 'created'
					 LEFT JOIN transactions t on td.transaction_id = t.id AND t.status = 'created'
					 LEFT JOIN cashout_details cd on v.id = cd.voucher_id AND cd.status = 'created'
					 LEFT JOIN cashouts c on cd.cashout_id = c.id AND c.status = 'created'
			WHERE v.created_at BETWEEN '` + dateFrom + ` 00:00:00+07'::timestamp AND '` + dateTo + ` 23:59:59+07'::timestamp 
			  AND v.program_id = p.id
			  AND v.status = 'created'
			  AND p.status = 'created'
			GROUP BY date `

	q = qp.GetQueryWithPagination(q, qp.GetQuerySort(), qp.GetQueryLimit())
	util.DEBUG(q)
	var resd []ReportReimburse
	err := db.Select(&resd, db.Rebind(q))
	if err != nil {
		return []ReportReimburse{}, false, err
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

func GetReportDailyVoucherTransactionWithOutlet(dateFrom, dateTo string, qp *util.QueryParam) ([]ReportReimburse, bool, error) {
	q := `
			SELECT date(v.created_at) as date,
					case
						   when p.id is not null then p.id::varchar
						   else '-' end
									  as outlet_idx,
				   case
						   when p.id is not null then p.name
						   else '-' end
									  as outlet_name,
				   sum(case
						   when t.id is not null then 1
						   else 0 end)
									  as transaction_qty,
				   sum(case
						   when c.id is not null then 1
						   else 0 end)
									  as reimburse_qty,
				   sum(case
						   when v.id is not null then 1
						   else 0 end)
									  as voucher_qty,
				   sum(case
						   when t.id is not null then v.voucher_value
						   else 0 end)
									  as total_transaction,
				   sum(case
						   when c.id is not null then v.voucher_value
						   else 0 end)
									  as total_reimburse,
				   sum(case
						   when v.id is not null then v.voucher_value
						   else 0 end)
									  as total_voucher
			FROM vouchers v
					 LEFT JOIN transaction_details td on td.voucher_id = v.id AND td.status = 'created'
					 LEFT JOIN transactions t on td.transaction_id = t.id AND t.status = 'created'
					 LEFT JOIN outlets p on t.outlet_id = p.id AND p.status = 'created'
					 LEFT JOIN cashout_details cd on v.id = cd.voucher_id AND cd.status = 'created'
					 LEFT JOIN cashouts c on cd.cashout_id = c.id AND c.status = 'created'
			WHERE v.created_at BETWEEN '` + dateFrom + ` 00:00:00+07'::timestamp AND '` + dateTo + ` 23:59:59+07'::timestamp 
			  AND v.status = 'created'
			GROUP BY date, outlet_idx, outlet_name `

	q = qp.GetQueryWithPagination(q, qp.GetQuerySort(), qp.GetQueryLimit())
	util.DEBUG(q)
	var resd []ReportReimburse
	err := db.Select(&resd, db.Rebind(q))
	if err != nil {
		return []ReportReimburse{}, false, err
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

func GetReportDailyOutletTransaction(dateFrom, dateTo string, qp *util.QueryParam) ([]ReportOutletTransaction, bool, error) {
	q := `
			SELECT date                 as transaction_date,
				   id                   as outlet_id,
				   name                 as outlet_name,
				   sum(trans)           as transaction_qty,
				   sum(vouc)            as voucher_qty,
				   sum(amount)          as total_transaction,
				   sum(reimburse_qty)   as reimburse_qty,
				   sum(total_reimburse) as total_reimburse
			FROM (
					 SELECT date(t.created_at) as date,
							p.id,
							p.name,
							sum(1)                          as trans,
							0                               as vouc,
							sum(t.total_amount)             as amount,
							0                               as reimburse_qty,
							0                               as total_reimburse
					 FROM transactions t,
						  outlets p
					 WHERE t.outlet_id = p.id
					   AND t.status = 'created'
					   AND p.status = 'created'
					   AND t.created_at BETWEEN '` + dateFrom + ` 00:00:00+07'::timestamp AND '` + dateTo + ` 23:59:59+07'::timestamp 
					 GROUP BY date, p.id, p.name
					 UNION ALL
					 SELECT date(t.created_at) as date,
							p.id,
							p.name,
							0                               as trans,
							sum(1)                          as vouc,
							0                               as amount,
							sum(case
									when cd.id is not null then 1
									else 0 end)             as reimburse_qty,
							sum(case
									when c.id is not null then c.amount
									else 0 end)             as total_reimburse
					 FROM transactions t,
						  outlets p,
						  transaction_details td
							  LEFT JOIN cashout_details cd ON cd.voucher_id = td.voucher_id AND cd.status = 'created'
							  LEFT JOIN cashouts c ON cd.cashout_id = c.id AND c.status = 'created'
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
	var resd []ReportOutletTransaction
	err := db.Select(&resd, db.Rebind(q))
	if err != nil {
		return []ReportOutletTransaction{}, false, err
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

func GetReportDailyOutletTransactionById(id, dateFrom, dateTo string, qp *util.QueryParam) ([]ReportOutletTransaction, bool, error) {
	q := `
			SELECT date                 as transaction_date,
				   id                   as outlet_id,
				   name                 as outlet_name,
				   sum(trans)           as transaction_qty,
				   sum(vouc)            as voucher_qty,
				   sum(amount)          as total_transaction,
				   sum(reimburse_qty)   as reimburse_qty,
				   sum(total_reimburse) as total_reimburse
			FROM (
					 SELECT date(t.created_at) as date,
							p.id,
							p.name,
							sum(1)                          as trans,
							0                               as vouc,
							sum(t.total_amount)             as amount,
							0                               as reimburse_qty,
							0                               as total_reimburse
					 FROM transactions t,
						  outlets p
					 WHERE p.id = ?
					   AND t.outlet_id = p.id
					   AND t.status = 'created'
					   AND p.status = 'created'
					   AND t.created_at BETWEEN '` + dateFrom + ` 00:00:00+07'::timestamp AND '` + dateTo + ` 23:59:59+07'::timestamp 
					 GROUP BY date, p.id, p.name
					 UNION ALL
					 SELECT date(t.created_at) as date,
							p.id,
							p.name,
							0                               as trans,
							sum(1)                          as vouc,
							0                               as amount,
							sum(case
									when cd.id is not null then 1
									else 0 end)             as reimburse_qty,
							sum(case
									when c.id is not null then c.amount
									else 0 end)             as total_reimburse
					 FROM transactions t,
						  outlets p,
						  transaction_details td
							  LEFT JOIN cashout_details cd ON cd.voucher_id = td.voucher_id AND cd.status = 'created'
							  LEFT JOIN cashouts c ON cd.cashout_id = c.id AND c.status = 'created'
					 WHERE p.id = ?
					   AND t.outlet_id = p.id
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
	var resd []ReportOutletTransaction
	err := db.Select(&resd, db.Rebind(q), id, id)
	if err != nil {
		util.DEBUG(err)
		return []ReportOutletTransaction{}, false, err
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

func GetReportMonthlyOutletTransaction(dateFrom, dateTo string, qp *util.QueryParam) ([]ReportOutletTransaction, bool, error) {
	q := `
			SELECT date as transaction_month, id as outlet_id, name as outlet_name, sum(trans) as transaction_qty, sum(vouc) as voucher_qty, sum(amount) as total_transaction
			FROM (
					 SELECT EXTRACT(MONTH FROM t.created_at) as date, p.id, p.name, sum(1) as trans, 0 as vouc, sum(t.total_amount) as amount
					 FROM transactions t,
						  outlets p
					 WHERE t.outlet_id = p.id
					   AND t.status = 'created'
					   AND p.status = 'created'
					   AND t.created_at BETWEEN '` + dateFrom + ` 00:00:00+07'::timestamp AND '` + dateTo + ` 23:59:59+07'::timestamp 
					 GROUP BY date, p.id, p.name
					 UNION ALL
					 SELECT EXTRACT(MONTH FROM t.created_at) as date, p.id, p.name, 0 as trans, sum(1) as vouc, 0 as amount
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
	var resd []ReportOutletTransaction
	err := db.Select(&resd, db.Rebind(q))
	if err != nil {
		return []ReportOutletTransaction{}, false, err
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

func GetReportYearlyOutletTransaction(dateFrom, dateTo string, qp *util.QueryParam) ([]ReportOutletTransaction, bool, error) {
	q := `
			SELECT date as transaction_year, id as outlet_id, name as outlet_name, sum(trans) as transaction_qty, sum(vouc) as voucher_qty, sum(amount) as total_transaction
			FROM (
					 SELECT EXTRACT(YEAR FROM t.created_at) as date, p.id, p.name, sum(1) as trans, 0 as vouc, sum(t.total_amount) as amount
					 FROM transactions t,
						  outlets p
					 WHERE t.outlet_id = p.id
					   AND t.status = 'created'
					   AND p.status = 'created'
					   AND t.created_at BETWEEN '` + dateFrom + ` 00:00:00+07'::timestamp AND '` + dateTo + ` 23:59:59+07'::timestamp 
					 GROUP BY date, p.id, p.name
					 UNION ALL
					 SELECT EXTRACT(YEAR FROM t.created_at) as date, p.id, p.name, 0 as trans, sum(1) as vouc, 0 as amount
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
	var resd []ReportOutletTransaction
	err := db.Select(&resd, db.Rebind(q))
	if err != nil {
		return []ReportOutletTransaction{}, false, err
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

func GetReportProgramTransaction(dateFrom, dateTo string, qp *util.QueryParam) ([]ReportProgramTransaction, bool, error) {
	q := `
		SELECT p.id,
			   p.name,
			   sum(case
					   when t.id is not null then 1
					   else 0 end)
								  as transaction_qty,
			   sum(case
					   when c.id is not null then 1
					   else 0 end)
								  as reimburse_qty,
			   sum(case
					   when v.id is not null then 1
					   else 0 end)
								  as voucher_qty,
			   sum(case
					   when t.id is not null then p.max_value
					   else 0 end)
								  as total_transaction,
			   sum(case
					   when c.id is not null then p.max_value
					   else 0 end)
								  as total_reimburse,
			   sum(case
					   when v.id is not null then p.max_value
					   else 0 end)
								  as total_voucher
		FROM programs p,
			 vouchers v
				 LEFT JOIN transaction_details td on td.voucher_id = v.id AND td.status = 'created'
				 LEFT JOIN transactions t on td.transaction_id = t.id AND t.status = 'created'
				 LEFT JOIN cashout_details cd on v.id = cd.voucher_id AND cd.status = 'created'
				 LEFT JOIN cashouts c on cd.cashout_id = c.id AND c.status = 'created'
		WHERE v.created_at BETWEEN '` + dateFrom + ` 00:00:00+07'::timestamp AND '` + dateTo + ` 23:59:59+07'::timestamp
		  AND v.program_id = p.id
		  AND v.status = 'created'
		  AND p.status = 'created'
		GROUP BY p.id, p.name `

	q = qp.GetQueryWithPagination(q, qp.GetQuerySort(), qp.GetQueryLimit())
	//ORDER BY year, month `
	util.DEBUG(q)
	var resd []ReportProgramTransaction
	err := db.Select(&resd, db.Rebind(q))
	if err != nil {
		return []ReportProgramTransaction{}, false, err
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

func GetReportProgramTransactionDaily(dateFrom, dateTo string, qp *util.QueryParam) ([]ReportProgramTransaction, bool, error) {
	q := `
		SELECT p.id, p.name, date(v.created_at) as date,
			   sum(case
					   when t.id is not null then 1
					   else 0 end)
								  as transaction_qty,
			   sum(case
					   when c.id is not null then 1
					   else 0 end)
								  as reimburse_qty,
			   sum(case
					   when v.id is not null then 1
					   else 0 end)
								  as voucher_qty,
			   sum(case
					   when t.id is not null then p.max_value
					   else 0 end)
								  as total_transaction,
			   sum(case
					   when c.id is not null then p.max_value
					   else 0 end)
								  as total_reimburse,
			   sum(case
					   when v.id is not null then p.max_value
					   else 0 end)
								  as total_voucher
		FROM programs p,
			 vouchers v
				 LEFT JOIN transaction_details td on td.voucher_id = v.id AND td.status = 'created'
				 LEFT JOIN transactions t on td.transaction_id = t.id AND t.status = 'created'
				 LEFT JOIN cashout_details cd on v.id = cd.voucher_id AND cd.status = 'created'
				 LEFT JOIN cashouts c on cd.cashout_id = c.id AND c.status = 'created'
		WHERE v.created_at BETWEEN '` + dateFrom + ` 00:00:00+07'::timestamp AND '` + dateTo + ` 23:59:59+07'::timestamp 
		  AND v.program_id = p.id
		  AND v.status = 'created'
		  AND p.status = 'created'
		GROUP BY date, p.id, p.name `

	q = qp.GetQueryWithPagination(q, qp.GetQuerySort(), qp.GetQueryLimit())
	//ORDER BY year, month `
	util.DEBUG(q)
	var resd []ReportProgramTransaction
	err := db.Select(&resd, db.Rebind(q))
	if err != nil {
		return []ReportProgramTransaction{}, false, err
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

func GetReportProgramTransactionDailyById(id, dateFrom, dateTo string, qp *util.QueryParam) ([]ReportProgramTransaction, bool, error) {
	q := `
		SELECT p.id, p.name, date(v.created_at) as date,
			   sum(case
					   when t.id is not null then 1
					   else 0 end)
								  as transaction_qty,
			   sum(case
					   when c.id is not null then 1
					   else 0 end)
								  as reimburse_qty,
			   sum(case
					   when v.id is not null then 1
					   else 0 end)
								  as voucher_qty,
			   sum(case
					   when t.id is not null then p.max_value
					   else 0 end)
								  as total_transaction,
			   sum(case
					   when c.id is not null then p.max_value
					   else 0 end)
								  as total_reimburse,
			   sum(case
					   when v.id is not null then p.max_value
					   else 0 end)
								  as total_voucher
		FROM programs p,
			 vouchers v
				 LEFT JOIN transaction_details td on td.voucher_id = v.id AND td.status = 'created'
				 LEFT JOIN transactions t on td.transaction_id = t.id AND t.status = 'created'
				 LEFT JOIN cashout_details cd on v.id = cd.voucher_id AND cd.status = 'created'
				 LEFT JOIN cashouts c on cd.cashout_id = c.id AND c.status = 'created'
		WHERE v.created_at BETWEEN '` + dateFrom + ` 00:00:00+07'::timestamp AND '` + dateTo + ` 23:59:59+07'::timestamp 
		  AND v.program_id = p.id
		  AND p.id = ?
		  AND v.status = 'created'
		  AND p.status = 'created'
		GROUP BY date, p.id, p.name `

	q = qp.GetQueryWithPagination(q, qp.GetQuerySort(), qp.GetQueryLimit())
	//ORDER BY year, month `
	util.DEBUG(q)
	var resd []ReportProgramTransaction
	err := db.Select(&resd, db.Rebind(q), id)
	if err != nil {
		return []ReportProgramTransaction{}, false, err
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

//DEMOGRAFI

type MemberDemografi struct {
	Gender        string `db:"gender" json:"gender,omitempty"`
	MaritalStatus string `db:"marital_status" json:"marital_status,omitempty"`
	LocationName  string `db:"location_name" json:"location_name,omitempty"`
	Age1723       *int   `db:"age1723" json:"age1723,omitempty"`
	Age2433       *int   `db:"age2433" json:"age2433,omitempty"`
	Age3444       *int   `db:"age3444" json:"age3444,omitempty"`
	Age4556       *int   `db:"age4556" json:"age4556,omitempty"`
	Age5770       *int   `db:"age5770" json:"age5770,omitempty"`
	Total         *int   `db:"total" json:"total,omitempty"`
	Count         int    `db:"count" json:"-"`
}

func GetDemografiGender(qp *util.QueryParam) ([]MemberDemografi, bool, error) {
	q := `
		SELECT a.gender, count(*) as total
		FROM accounts a
		GROUP BY a.gender`

	q = qp.GetQueryWithPagination(q, qp.GetQuerySort(), qp.GetQueryLimit())
	util.DEBUG(q)
	var resd []MemberDemografi
	err := db.Select(&resd, db.Rebind(q))
	if err != nil {
		return []MemberDemografi{}, false, err
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

func GetDemografiLocation(qp *util.QueryParam) ([]MemberDemografi, bool, error) {
	q := `
		SELECT a.birthplace as location_name, count(*) as total
		FROM accounts a
		GROUP BY a.birthplace`

	q = qp.GetQueryWithPagination(q, qp.GetQuerySort(), qp.GetQueryLimit())
	util.DEBUG(q)
	var resd []MemberDemografi
	err := db.Select(&resd, db.Rebind(q))
	if err != nil {
		return []MemberDemografi{}, false, err
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

func GetDemografiAge(qp *util.QueryParam) ([]MemberDemografi, bool, error) {
	q := `
		SELECT a.gender,
			sum(case when date_part('years', age(now(), birthdate)) between 17 and 23
			  then 1
				else 0 end) as "age1723",
			sum(case when date_part('years', age(now(), a.birthdate)) between 24 and 33
			  then 1
				else 0 end) as "age2433",
			sum(case when date_part('years', age(now(), a.birthdate)) between 34 and 44
			  then 1
				else 0 end) as "age3444",
			sum(case when date_part('years', age(now(), a.birthdate)) between 45 and 56
			  then 1
				else 0 end) as "age4556",
			sum(case when date_part('years', age(now(), a.birthdate)) between 57 and 70
			  then 1
				else 0 end) as "age5770", count(*) TOTAL
		FROM accounts a
		GROUP BY a.gender`

	q = qp.GetQueryWithPagination(q, qp.GetQuerySort(), qp.GetQueryLimit())
	util.DEBUG(q)
	var resd []MemberDemografi
	err := db.Select(&resd, db.Rebind(q))
	if err != nil {
		return []MemberDemografi{}, false, err
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

func GetDemografiMarital(qp *util.QueryParam) ([]MemberDemografi, bool, error) {
	q := `
		SELECT a.marital_status, count(*) as total
		FROM accounts a
		GROUP BY a.marital_status`

	q = qp.GetQueryWithPagination(q, qp.GetQuerySort(), qp.GetQueryLimit())
	util.DEBUG(q)
	var resd []MemberDemografi
	err := db.Select(&resd, db.Rebind(q))
	if err != nil {
		return []MemberDemografi{}, false, err
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
