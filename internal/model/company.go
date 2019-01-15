package model

import (
	"fmt"
	"time"

	"github.com/gilkor/evoucher/internal/util"
)

type (
	//Company : represent of partners table model
	Company struct {
		ID           string     `db:"id" json:"id,omitempty"`
		Name         string     `db:"name" json:"name,omitempty"`
		Description  string     `db:"description" json:"description,omitempty"`
		Alias        string     `db:"alias" json:"alias,omitempty"`
		ClientKey    string     `db:"client_key" json:"client_key,omitempty"`
		ClientSecret string     `db:"client_secret" json:"client_secret,omitempty"`
		CreatedAt    *time.Time `db:"created_at" json:"created_at,omitempty"`
		CreatedBy    string     `db:"created_by" json:"created_by,omitempty"`
		UpdatedAt    *time.Time `db:"updated_at" json:"updated_at,omitempty"`
		UpdatedBy    string     `db:"updated_by" json:"updated_by,omitempty"`
		Status       string     `db:"status" json:"status,omitempty"`
	}
)

// GetCompanyByAlias :  get list Companies by Alias
func GetCompanyByAlias(v string, f *util.Filter) ([]Company, bool, error) {
	return getCompanies("alias", v, f)
}

// GetCompanyByID :  get list Companies by ID
func GetCompanyByID(id string, f *util.Filter) ([]Company, bool, error) {
	return getCompanies("id", id, f)
}

// GetCompanies : list Company
func GetCompanies(f *util.Filter) ([]Company, bool, error) {
	return getCompanies("1", "1", f)
}

func getCompanies(key, value string, f *util.Filter) ([]Company, bool, error) {
	q := f.GetQueryByDefaultStruct(Company{})
	q += `
			FROM
				Companies
			WHERE 
				status = ?			
			AND ` + key + ` = ?`

	q += f.GetQuerySort()
	q += f.GetQueryLimit()
	fmt.Println(q)
	var resd []Company
	err := db.Select(&resd, db.Rebind(q), StatusCreated, value)
	if err != nil {
		return []Company{}, false, err
	}

	next := false
	if len(resd) > f.Count {
		next = true
	}
	if len(resd) < f.Count {
		f.Count = len(resd)
	}

	return resd, next, nil
}

//Insert : single row inset into table
func (c *Company) Insert() error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	q := `INSERT INTO 
				Companies 
				( 
					name
					, description 
					, alias 
					, client_key
					, client_secret 
					, created_by
					, updated_by
					, status
				)
			VALUES 
				( ?, ?, ?, ?, ?, ?, ?)
			RETURNING
				id
				, name
				, description 
				, alias 
				, client_key
				, client_secret 
				, created_at
				, created_by
				, updated_at
				, updated_by
				, status
	`
	var res []Company
	err = tx.Select(&res, tx.Rebind(q), c.Name, c.Description, c.Alias, c.ClientKey, c.ClientSecret, c.CreatedBy, c.UpdatedBy, StatusCreated)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

//Update : update Company
func (c *Company) Update() error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	q := `UPDATE
				Companies 
			SET
				name = ?
				, description = ?
				, alias = ?
				, client_key = ?
				, client_secret =?
				, updated_at = now()
				, updated_by = ?				
				, status = ?
			WHERE 
				id = ?
			RETURNING
				id
				, name
				, description 
				, alias 
				, client_key
				, client_secret 
				, created_at
				, created_by
				, updated_at
				, updated_by
				, status
	`
	var res []Company
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
func (c *Company) Delete() error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	q := `UPDATE
				Companies 
			SET
				deleted_at = now(),
				deleted_by = ?
				status = ?			
			WHERE 
				id = ?	
			RETURNING
				id
				, name
				, description 
				, alias 
				, client_key
				, client_secret 
				, created_at
				, created_by
				, updated_at
				, updated_by
				, status
	`
	var res []Company
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
