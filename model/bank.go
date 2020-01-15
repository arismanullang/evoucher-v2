package model

import (
	"time"

	"github.com/gilkor/evoucher-v2/util"
)

type (
	//Banks : array of bank
	Banks []Bank
	//Bank :
	Bank struct {
		// Bank            string `json:"bank,omitempty"`
		ID              string     `db:"id" json:"id,omitempty"`
		PartnerID       string     `db:"partner_id" json:"partner_id"`
		BankName        string     `db:"bank_name" json:"bank_name,omitempty"`
		BankBranch      string     `db:"bank_branch" json:"bank_branch,omitempty"`
		BankAccount     string     `db:"bank_account" json:"bank_account,omitempty"`
		BankAccountName string     `db:"bank_account_name" json:"bank_account_name,omitempty"`
		CompanyName     string     `db:"company_name" json:"company_name,omitempty"`
		Name            string     `db:"name" json:"name,omitempty"`
		Phone           string     `db:"phone" json:"phone,omitempty"`
		Email           string     `db:"email" json:"email,omitempty"`
		CreatedAt       *time.Time `db:"created_at" json:"created_at,omitempty"`
		CreatedBy       string     `db:"created_by" json:"created_by,omitempty"`
		UpdatedAt       *time.Time `db:"updated_at" json:"updated_at,omitempty"`
		UpdatedBy       string     `db:"updated_by" json:"updated_by,omitempty"`
		Status          string     `db:"status" json:"status,omitempty"`
	}
)

//BankFields : default table field
// var BankFields = []string{"id", "name", "description", "is_super", "created_at", "created_by", "updated_at", "updated_by", "status"}

//GetBanks : get list company by custom filter
func GetBanks(qp *util.QueryParam) (*Banks, bool, error) {
	return getBanks("1", "1", qp)
}

//GetBankByID : get bank by specified ID
func GetBankByID(qp *util.QueryParam, id string) (*Banks, bool, error) {
	return getBanks("id", id, qp)
}

//GetBankByPartnerID : get bank by specified partner ID
func GetBankByPartnerID(qp *util.QueryParam, id string) (*Banks, bool, error) {
	return getBanks("partner_id", id, qp)
}

func getBanks(k, v string, qp *util.QueryParam) (*Banks, bool, error) {

	q, err := qp.GetQueryByDefaultStruct(Bank{})
	if err != nil {
		return &Banks{}, false, err
	}
	// q := qp.GetQueryFields(BankFields)

	q += `
			FROM
				partner_banks bank
			WHERE 
				status = ?
			AND ` + k + ` = ?`

	q += qp.GetQuerySort()
	q += qp.GetQueryLimit()
	var resd Banks
	util.DEBUG(q)
	err = db.Select(&resd, db.Rebind(q), StatusCreated, v)
	if err != nil {
		return &Banks{}, false, err
	}
	if len(resd) < 1 {
		return &Banks{}, false, ErrorResourceNotFound
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

//Insert : save data to database
func (p *Bank) Insert() (*Banks, error) {
	tx, err := db.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	q := `INSERT INTO 
				partner_banks ( partner_id
					, bank_name
					, bank_branch
					, bank_account
					, bank_account_name
					, company_name
					, name
					, phone
					, email
					, created_at
					, created_by
					, updated_at
					, updated_by
					, status)
			VALUES 
				( ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			RETURNING
				id
				, partner_id
				, bank_name
				, bank_branch
				, bank_account
				, bank_account_name
				, company_name
				, name
				, phone
				, email
				, created_at
				, created_by
				, updated_at
				, updated_by
				, status
	`
	// bank, err := json.Marshal(p.Bank)
	if err != nil {
		return nil, err
	}
	var res Banks
	// util.DEBUG(p.Bank)
	err = tx.Select(&res, tx.Rebind(q), p.PartnerID, p.BankName, p.BankBranch, p.BankAccount, p.BankAccountName,
		p.CompanyName, p.Name, p.Phone, p.Email, p.CreatedAt, p.CreatedBy, p.CreatedAt, p.CreatedBy, StatusCreated)
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
func (p *Bank) Update() error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	q := `UPDATE
				partner_banks 
			SET
			bank_name = ?
			, bank_branch = ?
			, bank_account = ?
			, bank_account_name = ?
			, company_name = ?
			, name = ?
			, phone = ?
			, email = ?
			, updated_at = ?
			, updated_by = ?
			, status = ?				
			WHERE 
				partner_id = ?	
			RETURNING
			id
			, partner_id
			, bank_name
			, bank_branch
			, bank_account
			, bank_account_name
			, company_name
			, name
			, phone
			, email
			, created_at
			, created_by
			, updated_at
			, updated_by
			, status
	`
	var res []Bank
	err = tx.Select(&res, tx.Rebind(q), p.BankName, p.BankBranch, p.BankAccount, p.BankAccountName,
		p.CompanyName, p.Name, p.Phone, p.Email, p.UpdatedAt, p.UpdatedBy, p.Status, p.PartnerID)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	*p = res[0]
	return nil
}

//Delete : soft delated data by updateting row status to "deleted"
func (p *Bank) Delete() error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	q := `UPDATE
				partner_banks 
			SET
				updated_at = now(),
				updated_by = ?,
				status = ?			
			WHERE 
				id = ?	
			RETURNING
			id
			, partner_id
			, bank_name
			, bank_branch
			, bank_account
			, bank_account_name
			, company_name
			, name
			, phone
			, email
			, created_at
			, created_by
			, updated_at
			, updated_by
			, status
	`
	var res []Bank
	err = tx.Select(&res, tx.Rebind(q), p.UpdatedBy, StatusDeleted, p.ID)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
