package model

import (
	"fmt"
	"time"

	"github.com/gilkor/evoucher/internal/util"
)

type (
	//Customer :  represent of customers table model
	Customer struct {
		ID          string     `db:"id" json:"id,omitempty"`
		Name        string     `db:"name" json:"name,omitempty"`
		MobilePhone string     `db:"mobile_phone,null" json:"mobile_phone,omitempty"`
		Email       string     `db:"email,null" json:"email,omitempty"`
		RefID       string     `db:"ref_id,null" json:"ref_id,omitempty"`
		CompanyID   string     `db:"company_id" json:"company_id,omitempty"`
		CreatedAt   *time.Time `db:"created_at" json:"created_at,omitempty"`
		CreatedBy   string     `db:"created_by" json:"created_by,omitempty"`
		UpdatedAt   *time.Time `db:"updated_at" json:"updated_at,omitempty"`
		UpdatedBy   string     `db:"updated_by" json:"updated_by,omitempty"`
		DeletedAt   *time.Time `db:"deleted_at,null" json:"deleted_at,omitempty"`
		DeletedBy   string     `db:"deleted_by,null" json:"deleted_by,omitempty"`
		Status      string     `db:"status" json:"status,omitempty"`
	}
	//Customers :
	Customers []Customer
)

//GetCustomerByCompanyID : get list customers by CompanyID
func GetCustomerByCompanyID(id string, f *util.Filter) (*Customers, bool, error) {
	return getCustomers("company_id", id, f)
}

// GetCustomerByID :  get list customers by ID
func GetCustomerByID(id string, f *util.Filter) (*Customers, bool, error) {
	return getCustomers("id", id, f)
}

// GetCustomers : list customer
func GetCustomers(f *util.Filter) (*Customers, bool, error) {
	return getCustomers("1", "1", f)
}

func getCustomers(key, value string, f *util.Filter) (*Customers, bool, error) {
	q := f.GetQueryByDefaultStruct(Customer{})
	q += `
			FROM
				customers
			WHERE 
				status = ?			
			AND ` + key + ` = ?`

	q += f.GetQuerySort()
	q += f.GetQueryLimit()
	fmt.Println(q)
	var resd Customers
	err := db.Select(&resd, db.Rebind(q), StatusCreated, value)
	if err != nil {
		return &Customers{}, false, err
	}

	next := false
	if len(resd) > f.Count {
		next = true
	}
	if len(resd) < f.Count {
		f.Count = len(resd)
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
					, mobile_pone 
					, email 
					, ref_id 
					, company_id 
					, created_by
					, status)
			VALUES 
				( ?, ?, ?, ?, ?, ?, ?)
			RETURNING
				id
				, name
				, mobile_pone 
				, email 
				, ref_id 
				, company_id 
				, created_at
				, created_by
				, updated_at
				, updated_by
				, deleted_at
				, deleted_by
				, status
	`
	var res []Customer
	err = tx.Select(&res, tx.Rebind(q), c.Name, c.MobilePhone, c.MobilePhone, c.Email, c.RefID, c.CompanyID, c.CreatedBy, StatusCreated)
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
				, mobile_pone 
				, email 
				, ref_id 
				, company_id 
				, created_at
				, created_by
				, updated_at
				, updated_by
				, deleted_at
				, deleted_by
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
				deleted_at = now(),
				deleted_by = ?
				status = ?			
			WHERE 
				id = ?	
			RETURNING
				id
				, name
				, mobile_pone 
				, email 
				, ref_id 
				, company_id 
				, created_at
				, created_by
				, updated_at
				, updated_by
				, deleted_at
				, deleted_by
				, status
	`
	var res []Customer
	err = tx.Select(&res, tx.Rebind(q), c.DeletedBy, StatusDeleted)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
