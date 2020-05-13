package model

import (
	"database/sql"
	"github.com/gilkor/evoucher-v2/util"
)

//GetUserAccumulativeVoucher : Getcount voucher by accountID and programID
func GetUserAccumulativeVoucher(accountID, programID string) (int, error) {

	q := ` SELECT COUNT(*) acc FROM vouchers 
			WHERE holder = ? AND program_id = ? AND status = ?`

	var r int
	err := db.QueryRow(db.Rebind(q), accountID, programID, StatusCreated).Scan(&r)
	if err != nil {
		return -1, err
	}
	util.DEBUG("result:", r, accountID, programID, q)

	return r, nil
}

//GetUserAccumulativeTransaction : Getcount voucher by accountID and programID
func GetUserAccumulativeTransaction(accountID, programID string) (int, error) {

	q := ` SELECT count(*) FROM transactions t, transaction_details td
			WHERE
				t.id = td.transaction_id
				AND t.holder = ?
				AND td.program_id = ?
			GROUP BY td.transaction_id`

	var r int
	err := db.QueryRow(db.Rebind(q), accountID, programID).Scan(&r)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return -1, err
	}
	util.DEBUG("result:", r, accountID, programID, q)

	return r, nil
}

//GetUserSpendingTransaction : Get Total Spending transaction by accountID and programID
//TODO need more parameter for filter (periode, program_id, outlet, etc)
func GetUserSpendingTransaction(accountID string) (int, error) {

	q := ` SELECT sum(total_amount) spending 
			FROM transactions
			WHERE holder = ? 
			GROUP BY holder`

	var r int
	err := db.QueryRow(db.Rebind(q), accountID).Scan(&r)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return -1, err
	}
	util.DEBUG("result:", r, accountID, q)

	return r, nil
}

func GetUserUsageVoucher(accountID, programID string) (int, error) {

	q := ` SELECT COUNT(*) acc FROM vouchers 
			WHERE holder = ? AND program_id = ? AND status = ?`

	var r int
	err := db.QueryRow(db.Rebind(q), accountID, programID, VoucherStateUsed).Scan(&r)
	if err != nil {
		return -1, err
	}
	util.DEBUG("result:", r, accountID, programID, q)

	return r, nil
}
