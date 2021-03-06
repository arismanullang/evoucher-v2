package model

import (
	"time"

	"github.com/gilkor/evoucher-v2/util"
)

type (
	Account struct {
		ID                string      `db:"id" json:"id"`
		Name              string      `db:"name" json:"name"`
		CompanyID         string      `db:"company_id" json:"company_id"`
		Gender            string      `db:"gender" json:"gender"`
		Email             string      `db:"email" json:"email"`
		MobileCallingCode string      `db:"mobile_calling_code" json:"mobile_calling_code"`
		MobileNo          string      `db:"mobile_no" json:"mobile_no"`
		State             string      `db:"state" json:"state"`
		Status            string      `db:"status" josn:"status"`
		CreatedBy         string      `db:"created_by" json:"created_by"`
		CreatedAt         time.Time   `db:"created_at" json:"created_at"`
		UpdatedBy         string      `db:"updated_by" json:"updated_by"`
		UpdatedAt         time.Time   `db:"updated_at" json:"updated_at"`
		DeletedBy         string      `db:"deleted_by" json:"deleted_by"`
		DeletedAt         interface{} `db:"deleted_at" json:"deleted_at"`
		Count             int         `db:"count" json:"-"`
	}

	Accounts []Account
)

//AccountFields : default table field
var AccountFields = []string{"id", "name", "company_id", "gender", "email", "mobile_calling_code", "mobile_no", "state", "status", "created_at", "created_by", "updated_at", "updated_by", "deleted_at", "deleted_by"}

//GetAccounts : get list company by custom filter
func GetAccounts(qp *util.QueryParam) (*Accounts, bool, error) {
	return getAccounts("1", "1", qp)
}

//GetAccountByID : get account by specified ID
func GetAccountsByID(qp *util.QueryParam, id string) (*Accounts, bool, error) {
	return getAccounts("id", id, qp)
}

// GetAccountByID :  get account
func GetAccountByID(qp *util.QueryParam, id string) (*Account, error) {

	accounts, _, err := getAccounts("id", id, qp)
	if err != nil {
		return &Account{}, err
	}
	account := &(*accounts)[0]

	return account, nil
}

func VerifyAccountToken(accountToken string) (*JWTJunoClaims, error) {

	// TODO: token not validated yet

	token, err := VerifyJWT(accountToken)
	if err != nil {
		if err == ErrorForbidden {
			return nil, err
		} else if err == ErrorInternalServer {
			return nil, err
		}
	}

	claims, ok := token.Claims.(*JWTJunoClaims)
	if ok && token.Valid {
		return claims, nil
	}

	return nil, err
}

func getAccounts(k, v string, qp *util.QueryParam) (*Accounts, bool, error) {

	q, err := qp.GetQueryByDefaultStruct(Account{})
	if err != nil {
		return &Accounts{}, false, err
	}
	// q := qp.GetQueryFields(AccountFields)

	q += `
			FROM
				accounts account
			WHERE 
				account.status = ?
			AND ` + k + ` = ?`

	q = qp.GetQueryWhereClause(q, qp.Q)
	q = qp.GetQueryWithPagination(q, qp.GetQuerySort(), qp.GetQueryLimit())

	var resd Accounts
	err = db.Select(&resd, db.Rebind(q), StatusCreated, v)
	if err != nil {
		return &Accounts{}, false, err
	}
	if len(resd) < 1 {
		return &Accounts{}, false, ErrorResourceNotFound
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

//Insert : save data to database
func (p *Account) Insert() (*Accounts, error) {
	tx, err := db.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	q := `INSERT INTO 
			accounts (
					id,
					name,
					company_id,
					gender,
					email,
					mobile_calling_code,
					mobile_no,
					state,
					status,
					created_at,
					created_by,
					updated_at,
					updated_by,
					deleted_at,
					deleted_by
			)
			VALUES (
				?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
			)
			RETURNING
				id,
				name,
				company_id,
				gender,
				email, 
				mobile_calling_code,
				mobile_no,
				state,
				status,
				created_at,
				created_by,
				updated_at,
				updated_by,
				deleted_at,
				deleted_by
		`
	if err != nil {
		return nil, err
	}
	var res Accounts
	err = tx.Select(&res, tx.Rebind(q), p.ID, p.Name, p.CompanyID, p.Gender, p.Email, p.MobileCallingCode, p.MobileNo, p.State, p.Status, p.CreatedAt, p.CreatedBy, p.UpdatedAt, p.UpdatedBy, p.DeletedAt, p.DeletedBy)
	if err != nil {
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	*p = res[0]
	return &res, nil
}

//Update : modify data
func (p *Account) Update() error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	q := `UPDATE
				accounts 
			SET
				name = ?,
				company_id = ?,
				gender = ?,
				email = ?,
				mobile_calling_code = ?,
				mobile_no = ?,
				status = ?,
				state = ?,
				created_at = ?,
				created_by = ?,
				updated_at = ?,
				updated_by = ?,
				deleted_at = ?,
				deleted_by =?			
			WHERE 
				id = ?	
			RETURNING
				id,
				name,
				company_id,
				gender,
				email, 
				mobile_calling_code,
				mobile_no,
				state,
				status,
				created_at,
				created_by,
				updated_at,
				updated_by,
				deleted_at,
				deleted_by
	`
	var res []Account
	err = tx.Select(&res, tx.Rebind(q), p.Name, p.CompanyID, p.Gender, p.Email, p.MobileCallingCode, p.MobileNo, p.Status, p.State, p.CreatedAt, p.CreatedBy, p.UpdatedAt, p.UpdatedBy, p.DeletedAt, p.DeletedBy, p.ID)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	if len(res) < 1 {
		return ErrorResourceNotFound
	}

	*p = res[0]
	return nil
}
