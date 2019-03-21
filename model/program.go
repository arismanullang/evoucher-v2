package model

import (
	"fmt"
	"time"

	"github.com/gilkor/evoucher/util"
)

type (
	//Program : base model
	Program struct {
		ID          string      `db:"id" json:"id"`
		CompanyID   string      `db:"company_id" json:"company_id"`
		Name        string      `db:"name" json:"name,omitempty"`
		Type        string      `db:"type" json:"type,omitempty"`
		Value       float64     `db:"value" json:"value,omitempty"`
		MaxValue    float64     `db:"max_value" json:"max_value,omitempty"`
		StartDate   time.Time   `db:"start_date" json:"start_date,omitempty"`
		EndDate     time.Time   `db:"end_date" json:"end_date,omitempty"`
		Description interface{} `db:"description" json:"description,omitempty"`
		ImageURL    string      `db:"image_url" json:"image_url,omitempty"`
		Template    string      `db:"template" json:"template,omitempty"`
		Rule        string      `db:"rule" json:"rule"`
		State       string      `db:"state" json:"state"`
		Stock       int64       `db:"stock" json:"stock"`
		CreatedAt   *time.Time  `db:"created_at" json:"created_at,omitempty"`
		CreatedBy   string      `db:"created_by" json:"created_by,omitempty"`
		UpdatedAt   *time.Time  `db:"updated_at" json:"updated_at,omitempty"`
		UpdatedBy   string      `db:"updated_by" json:"updated_by,omitempty"`
		Status      string      `db:"status" json:"status,omitempty"`
		Partners    Partners    `json:"partners"`
	}
	// Programs : base model
	Programs []Program
)

// GetProgramByID :  program details
func GetProgramByID(id string, qp *util.QueryParam) (*Programs, bool, error) {
	return getPrograms("id", id, qp)
}

// GetPrograms : get program list
func GetPrograms(qp *util.QueryParam) (*Programs, bool, error) {
	return getPrograms("1", "1", qp)
}

func getPrograms(key, value string, qp *util.QueryParam) (*Programs, bool, error) {
	q, err := qp.GetQueryByDefaultStruct(Programs{})
	if err != nil {
		return &Programs{}, false, err
	}
	q += `
			FROM
				programs
			WHERE 
				status = ?			
			AND ` + key + ` = ?`

	q += qp.GetQuerySort()
	q += qp.GetQueryLimit()
	fmt.Println(q)
	var resd Programs
	err = db.Select(&resd, db.Rebind(q), StatusCreated, value)
	if err != nil {
		return &Programs{}, false, err
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
func (p Program) Insert() error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	q := `INSERT INTO 
				programs 
				( 
					company_id
					, name
					, type
					, value
					, max_value
					, start_date
					, end_date
					, description
					, image_url
					, template
					, rule
					, state
					, description
					, stock
					, created_by
					, updated_by					
					, status
				)
			VALUES 
				( ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			RETURNING
			id
			, company_id
			, name
			, type
			, value
			, max_value
			, start_date
			, end_date
			, description
			, image_url
			, template
			, rule
			, state
			, description
			, stock
			, created_at
			, created_by
			, updated_at
			, updated_by					
			, status
	`
	var res []Customer
	err = tx.Select(&res, tx.Rebind(q), p.CompanyID, p.Name, p.Type, p.Value, p.MaxValue, p.StartDate, p.EndDate,
		p.Description, p.ImageURL, p.Template, p.Rule, p.State, p.Description, p.Stock, p.CreatedBy, p.UpdatedBy, StatusCreated)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

//Update : update program
func (p Program) Update() error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	q := `UPDATE
				programs 
			SET
					company_id = ?
					, name = ?
					, type = ?
					, value = ?
					, max_value = ?
					, start_date = ?
					, end_date = ?
					, description = ?
					, image_url = ?
					, template = ?
					, rule = ?
					, state = ?
					, description = ?
					, stock = ?
					, updated_by = ?		
			WHERE 
				id = ?
			RETURNING
				id
				, company_id
				, name
				, type
				, value
				, max_value
				, start_date
				, end_date
				, description
				, image_url
				, template
				, rule
				, state
				, description
				, stock
				, created_at
				, created_by
				, updated_at
				, updated_by					
				, status
	`
	var res []Customer
	err = tx.Select(&res, tx.Rebind(q), p.CompanyID, p.Name, p.Type, p.Value, p.MaxValue, p.StartDate, p.EndDate,
		p.Description, p.ImageURL, p.Template, p.Rule, p.State, p.Description, p.Stock, p.UpdatedBy, p.ID)
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
func (p Program) Delete() error {
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
				, company_id
				, name
				, type
				, value
				, max_value
				, start_date
				, end_date
				, description
				, image_url
				, template
				, rule
				, state
				, description
				, stock
				, created_at
				, created_by
				, updated_at
				, updated_by					
				, status
	`
	var res []Customer
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
