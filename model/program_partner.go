package model

import (
	"bytes"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/gilkor/evoucher-v2/util"
)

type (
	// ProgramOutlet model
	ProgramOutlet struct {
		ID        string     `db:"id" json:"id,omitempty"`
		ProgramID string     `db:"program_id" json:"program_id,omitempty"`
		OutletID  string     `db:"outlet_id" json:"outlet_id,omitempty"`
		CreatedAt *time.Time `db:"created_at" json:"created_at,omitempty"`
		CreatedBy string     `db:"created_by" json:"created_by,omitempty"`
		UpdatedAt *time.Time `db:"updated_at" json:"updated_at,omitempty"`
		UpdatedBy string     `db:"updated_by" json:"updated_by,omitempty"`
		Status    string     `db:"status" json:"status,omitempty"`
	}
	//ProgramOutlets model
	ProgramOutlets []ProgramOutlet
)

// GetOutletByProgramID :
func GetOutletByProgramID(programID string, qp *util.QueryParam) (*Outlets, bool, error) {
	q, err := qp.GetQueryByDefaultStruct(Outlet{})
	if err != nil {
		return &Outlets{}, false, err
	}
	q += `
			FROM
				program_outlets ProgramOutlet
				INNER JOIN m_outlets outlet ON ProgramOutlet.outlet_id = outlet.id
				INNER JOIN programs program ON ProgramOutlet.program_id = program.id

			WHERE 
				ProgramOutlet.status = ?
			AND 
				ProgramOutlet.program_id = ?
		`
	// q += qp.GetQuerySort()
	// q += qp.GetQueryLimit()
	var resd Outlets
	err = db.Select(&resd, db.Rebind(q), StatusCreated, programID)
	if err != nil {
		return &Outlets{}, false, err
	}
	if len(resd) < 1 {
		return &Outlets{}, false, ErrorResourceNotFound
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

// GetProgramByOutletID :
func GetProgramByOutletID(programID string, qp *util.QueryParam) (*Programs, bool, error) {
	q, err := qp.GetQueryByDefaultStruct(Program{})
	if err != nil {
		return &Programs{}, false, err
	}
	q += `
			FROM
				program_outlets ProgramOutlet
				INNER JOIN outlets outlet ON ProgramOutlet.outlet_id = outlet.id
				INNER JOIN programs program ON ProgramOutlet.program_id = program.id

			WHERE 
				ProgramOutlet.status = ?
			AND 
				ProgramOutlet.program_id = ?
		`
	q += qp.GetQuerySort()
	q += qp.GetQueryLimit()

	fmt.Println("query struct :", q)
	var resd Programs
	err = db.Select(&resd, db.Rebind(q), StatusCreated, programID)
	if err != nil {
		return &Programs{}, false, err
	}
	if len(resd) < 1 {
		return &Programs{}, false, ErrorResourceNotFound
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

//NewProgramOutlets : initiate Program outlets obj
func NewProgramOutlets(programID string, outlets Outlets) *ProgramOutlets {
	pp := make(ProgramOutlets, len(outlets))
	for k, v := range outlets {
		pp[k].ProgramID = programID
		pp[k].OutletID = v.ID
	}
	return &pp
}

//Upsert : insert data using upsert/append query
func (pp *ProgramOutlets) Upsert(tx *sqlx.Tx) error {
	values := new(bytes.Buffer)
	var args []interface{}
	for _, v := range *pp {
		values.WriteString("(?, ?, ?, ?, ?),")
		args = append(args, v.ProgramID, v.OutletID, v.CreatedBy, v.UpdatedBy, StatusCreated)
	}

	q := `INSERT INTO 
				program_outlets
				( 				
					 program_id
					, outlet_id	
					, created_by
					, updated_by
					, status
				)
			VALUES 
			`
	valuestr := values.String()
	q += valuestr[:len(valuestr)-1]

	q += ` 	
			ON CONFLICT 
				( program_id , outlet_id ) 
			DO UPDATE 
			SET
				program_id = excluded.program_id
				,outlet_id = excluded.outlet_id
				,updated_at = now()
				,status = excluded.status
			RETURNING
				id
				, program_id
				, outlet_id
				, created_at
				, created_by
				, updated_at
				, updated_by
				, status
	`
	// fmt.Println("pp query : ", q)
	util.DEBUG(q, args)

	var res ProgramOutlets
	err := tx.Select(&res, tx.Rebind(q), args...)
	if err != nil {
		return err
	}

	*pp = res

	return nil
}

// Delete : delete program outlets by program ID
func (pp *ProgramOutlets) Delete(tx *sqlx.Tx) error {

	q := `UPDATE
				program_outlets 
			SET
				updated_at = now()
				, status = ?
			WHERE 
				program_id = ?	
			RETURNING
			id
			, program_id
			, outlet_id
			, created_at
			, created_by
			, updated_at
			, updated_by
			, status
	`

	programID := (*pp)[0].ProgramID

	var res ProgramOutlets
	err := tx.Select(&res, tx.Rebind(q), StatusDeleted, programID)
	if err != nil {
		return errors.New("Failed when delete program outlet ," + ErrorNoDataAffected.Error() + " , " + err.Error())
	}

	*pp = res

	return nil
}
