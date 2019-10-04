package model

import (
	"time"

	"github.com/gilkor/evoucher-v2/util"
)

type (
	//Customer :  represent of customers table model
	Customer struct {
		ID          string     `db:"id" json:"id,omitempty"`
		Name        string     `db:"name" json:"name,omitempty"`
		MobilePhone string     `db:"mobile_phone" json:"mobile_phone,omitempty"`
		Email       string     `db:"email" json:"email,omitempty"`
		RefID       string     `db:"ref_id" json:"ref_id,omitempty"`
		CompanyID   string     `db:"company_id" json:"company_id,omitempty"`
		CreatedAt   *time.Time `db:"created_at" json:"created_at,omitempty"`
		CreatedBy   string     `db:"created_by" json:"created_by,omitempty"`
		UpdatedAt   *time.Time `db:"updated_at" json:"updated_at,omitempty"`
		UpdatedBy   string     `db:"updated_by" json:"updated_by,omitempty"`
		Status      string     `db:"status" json:"status,omitempty"`
		Tags        Tags       `json:"tags"`
	}
	//Customers :
	Customers []Customer
)

//GetCustomerByCompanyID : get list customers by CompanyID
func GetCustomerByCompanyID(id string, qp *util.QueryParam) (*Customers, bool, error) {
	return getCustomers("company_id", id, qp)
}

// GetCustomerByID :  get list customers by ID
func GetCustomerByID(id string, qp *util.QueryParam) (*Customers, bool, error) {
	return getCustomers("id", id, qp)
}

// GetCustomers : list customer
func GetCustomers(qp *util.QueryParam) (*Customers, bool, error) {
	return getCustomers("1", "1", qp)
}

func getCustomers(key, value string, qp *util.QueryParam) (*Customers, bool, error) {
	q, err := qp.GetQueryByDefaultStruct(Customer{})
	if err != nil {
		return &Customers{}, false, err
	}
	q += `
			FROM
				customers customer
			WHERE 
				status = ?			
			AND ` + key + ` = ?`

	q += qp.GetQuerySort()
	q += qp.GetQueryLimit()

	util.DEBUG(q)
	var resd Customers
	err = db.Select(&resd, db.Rebind(q), StatusCreated, value)
	if err != nil {
		return &Customers{}, false, err
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

//Insert : single row inset into table
func (c Customer) Insert() error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	q := `INSERT INTO 
				customers 
				( 
					name
					, mobile_phone 
					, email 
					, ref_id 
					, company_id 
					, created_by
					, updated_by
					, status)
			VALUES 
				( ?, ?, ?, ?, ?, ?, ?, ?)
			RETURNING
				id
				, name
				, mobile_phone 
				, email 
				, ref_id 
				, company_id 
				, created_at
				, created_by
				, updated_at
				, updated_by
				, status
	`
	var res []Customer
	err = tx.Select(&res, tx.Rebind(q), c.Name, c.MobilePhone, c.Email, c.RefID, c.CompanyID, c.CreatedBy, c.CreatedBy, StatusCreated)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

//Update : update customer
func (c *Customer) Update() error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	q := `UPDATE
				customers 
			SET
				name = ?,
				mobule_phine = ?,
				email = ?,
				ref_id = ?,
				company_id = ?,
				updated_at = now(),
				updated_by = ?,					
				status = ?
			WHERE 
				id = ?
			RETURNING
				id
				, name
				, mobile_phone 
				, email 
				, ref_id 
				, company_id 
				, created_at
				, created_by
				, updated_at
				, updated_by
				, status
	`
	var res []Customer
	err = tx.Select(&res, tx.Rebind(q), c.Name, c.UpdatedBy, c.ID)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

//Delete : soft delated data by updateting row status to "deleted"
func (c *Customer) Delete() error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	q := `UPDATE
				customers 
			SET
				updated_at = now(),
				updated_by = ?
				status = ?			
			WHERE 
				id = ?	
			RETURNING
				id
				, name
				, mobile_phone 
				, email 
				, ref_id 
				, company_id 
				, created_at
				, created_by
				, updated_at
				, updated_by
				, status
	`
	var res []Customer
	err = tx.Select(&res, tx.Rebind(q), c.UpdatedBy, StatusDeleted)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
