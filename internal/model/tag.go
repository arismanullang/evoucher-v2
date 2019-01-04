package model

import (
	"time"

	"github.com/gilkor/evoucher/internal/util"
)

type (
	//Tag : represent of tags table model
	Tag struct {
		ID        string     `db:"id" json:"id,omitempty"`
		Name      string     `db:"name" json:"name,omitempty"`
		CompanyID string     `db:"company_id" json:"company_id"`
		CreatedAt *time.Time `db:"created_at" json:"created_at,omitempty"`
		CreatedBy string     `db:"created_by" json:"created_by,omitempty"`
		UpdatedAt *time.Time `db:"updated_at" json:"updated_at,omitempty"`
		UpdatedBy string     `db:"updated_by" json:"updated_by,omitempty"`
		Status    string     `db:"status" json:"status,omitempty"`
	}
	//Tags :
	Tags []Tag
)

//TagFields : default table field
var TagFields = []string{"id", "name", "company_id", "created_at", "created_by", "updated_at", "updated_by", "status"}

//GetTags : get list company by custom filter
func GetTags(f *util.Filter) (*Tags, bool, error) {
	return getTags("1", "1", f)
}

//GetTagByCompanyID : get partner by specified ID
func GetTagByCompanyID(f *util.Filter, id string) (*Tags, bool, error) {
	return getTags("company_id", id, f)
}

//GetTagByID : get partner by specified ID
func GetTagByID(f *util.Filter, id string) (*Tags, bool, error) {
	return getTags("id", id, f)
}

func getTags(k, v string, f *util.Filter) (*Tags, bool, error) {

	// q := f.GetQueryByDefaultStruct(Partner{})
	q := f.GetQueryFields(PartnerFields)
	q += `
			FROM
				tags
			WHERE 
				status = ?
			AND ` + k + ` = ?`

	q += f.GetQuerySort()
	q += f.GetQueryLimit()
	// fmt.Println(q)
	var resd Tags
	err := db.Select(&resd, db.Rebind(q), StatusCreated, v)
	if err != nil {
		return &Tags{}, false, err
	}
	if len(resd) < 1 {
		return &Tags{}, false, ErrorResourceNotFound
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

//Insert : save data to database
func (t *Tag) Insert() error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	q := `INSERT INTO 
				tags ( name, company_id, created_by, status)
			VALUES 
				( ?, ?, ?, ?)
			RETURNING
				id, name, company_id, created_at, created_by, updated_at, updated_by, status
	`
	var res Tags
	err = tx.Select(&res, tx.Rebind(q), t.Name, t.CompanyID, t.CreatedBy, StatusCreated)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	*t = res[0]
	return nil
}

//Update : modify data
func (t *Tag) Update() error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	q := `UPDATE
				tags 
			SET
				name = ?,
				company_id = ?,
				updated_at = now(),
				updated_by = ?				
			WHERE 
				id = ?	
			RETURNING
			id, name, company_id ,created_at, created_by, updated_at, updated_by, status
	`
	var res Tags
	err = tx.Select(&res, tx.Rebind(q), t.Name, t.CompanyID, t.UpdatedBy, t.ID)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	*t = res[0]
	return nil
}

//Delete : soft delated data by updateting row status to "deleted"
func (t *Tag) Delete() error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	q := `UPDATE
				tags 
			SET
				updated_at = now(),
				updated_by = ?
				status = ?			
			WHERE 
				id = ?	
			RETURNING
			id, name, company_id, created_at, created_by, updated_at, updated_by, status
	`
	var res []Partner
	err = tx.Select(&res, tx.Rebind(q), t.UpdatedBy, StatusDeleted, t.ID)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
