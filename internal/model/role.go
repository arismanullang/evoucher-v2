package model

import (
	"fmt"
	"time"

	"github.com/gilkor/evoucher/internal/util"
)

type (
	//Role :
	Role struct {
		ID        string     `db:"id" json:"id"`
		Name      string     `db:"name" json:"name"`
		CompanyID string     `db:"company_id" json:"company_id"`
		CreatedAt *time.Time `db:"created_at" json:"created_at,omitempty"`
		CreatedBy string     `db:"created_by" json:"created_by,omitempty"`
		UpdatedAt *time.Time `db:"updated_at" json:"updated_at,omitempty"`
		UpdatedBy string     `db:"updated_by" json:"updated_by,omitempty"`
		Status    string     `db:"status" json:"status,omitempty"`
	}
	//Roles :
	Roles []Role
)

//GetRoleByCompanyID : get list Roles by CompanyID
func GetRoleByCompanyID(id string, f *util.Filter) (*Roles, bool, error) {
	return getRoles("company_id", id, f)
}

// GetRoleByID :  get list Roles by ID
func GetRoleByID(id string, f *util.Filter) (*Roles, bool, error) {
	return getRoles("id", id, f)
}

// GetRoles : list Role
func GetRoles(f *util.Filter) (*Roles, bool, error) {
	return getRoles("1", "1", f)
}

func getRoles(key, value string, f *util.Filter) (*Roles, bool, error) {
	q := f.GetQueryByDefaultStruct(Role{})
	q += `
			FROM
				Roles
			WHERE 
				status = ?			
			AND ` + key + ` = ?`

	q += f.GetQuerySort()
	q += f.GetQueryLimit()
	fmt.Println(q)
	var resd Roles
	err := db.Select(&resd, db.Rebind(q), StatusCreated, value)
	if err != nil {
		return &Roles{}, false, err
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
func (c Role) Insert() error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	q := `INSERT INTO 
				Roles 
				( 
					name
					, company_id 
					, created_by
					, status)
			VALUES 
				( ?, ?, ?, ?)
			RETURNING
				id
				, name
				, company_id 
				, created_at
				, created_by
				, updated_at
				, updated_by
				, status
	`
	var res []Role
	err = tx.Select(&res, tx.Rebind(q), c.Name, c.CompanyID, c.CreatedBy, StatusCreated)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

//Update : update Role
func (c *Role) Update() error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	q := `UPDATE
				Roles 
			SET
				name = ?,
				company_id = ?,
				updated_at = now(),
				updated_by = ?,					
				status = ?
			WHERE 
				id = ?
			RETURNING
				id
				, company_id 
				, created_at
				, created_by
				, updated_at
				, updated_by
				, status
	`
	var res []Role
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
func (c *Role) Delete() error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	q := `UPDATE
				Roles 
			SET
				updated_at = now(),
				updated_by = ?
				status = ?			
			WHERE 
				id = ?	
			RETURNING
				id
				, name
				, company_id 
				, created_at
				, created_by
				, updated_at
				, updated_by
				, status
	`
	var res []Role
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
