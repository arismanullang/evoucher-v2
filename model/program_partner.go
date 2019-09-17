package model

import (
	"bytes"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/gilkor/evoucher/util"
)

type (
	// ProgramPartner model
	ProgramPartner struct {
		ID        string     `db:"id" json:"id,omitempty"`
		ProgramID string     `db:"program_id" json:"program_id,omitempty"`
		PartnerID string     `db:"partner_id" json:"partner_id,omitempty"`
		CreatedAt *time.Time `db:"created_at" json:"created_at,omitempty"`
		CreatedBy string     `db:"created_by" json:"created_by,omitempty"`
		UpdatedAt *time.Time `db:"updated_at" json:"updated_at,omitempty"`
		UpdatedBy string     `db:"updated_by" json:"updated_by,omitempty"`
		Status    string     `db:"status" json:"status,omitempty"`
		Programs  Programs
		Partners  Programs
	}
	//ProgramPartners model
	ProgramPartners []ProgramPartner
)

// GetPartnerByProgramID :
func GetPartnerByProgramID(programID string, qp *util.QueryParam) (*Partners, bool, error) {
	q, err := qp.GetQueryByDefaultStruct(Partner{})
	if err != nil {
		return &Partners{}, false, err
	}
	q += `
			FROM
				program_partners ProgramPartner
				INNER JOIN partners partner ON ProgramPartner.partner_id = partner.id
				INNER JOIN programs program ON ProgramPartner.program_id = program.id

			WHERE 
				ProgramPartner.status = ?
			AND 
				ProgramPartner.program_id = ?
		`
	q += qp.GetQuerySort()
	q += qp.GetQueryLimit()
	// fmt.Println(q)
	fmt.Println("query struct :", q)
	var resd Partners
	err = db.Select(&resd, db.Rebind(q), StatusCreated, programID)
	if err != nil {
		return &Partners{}, false, err
	}
	if len(resd) < 1 {
		return &Partners{}, false, ErrorResourceNotFound
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

// GetProgramByPartnerID :
func GetProgramByPartnerID(programID string, qp *util.QueryParam) (*Programs, bool, error) {
	q, err := qp.GetQueryByDefaultStruct(Program{})
	if err != nil {
		return &Programs{}, false, err
	}
	q += `
			FROM
				program_partners ProgramPartner
				INNER JOIN partners partner ON ProgramPartner.partner_id = partner.id
				INNER JOIN programs program ON ProgramPartner.program_id = program.id

			WHERE 
				ProgramPartner.status = ?
			AND 
				ProgramPartner.program_id = ?
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

//NewProgramPartners : initiate Program partners obj
func NewProgramPartners(programID string, partners Partners) *ProgramPartners {
	pp := make(ProgramPartners, len(partners))
	for k, v := range partners {
		pp[k].ProgramID = programID
		pp[k].PartnerID = v.ID
	}
	return &pp
}

//Upsert : insert data using upsert/append query
func (pp *ProgramPartners) Upsert(tx *sqlx.Tx) error {
	values := new(bytes.Buffer)
	var args []interface{}
	for _, v := range *pp {
		values.WriteString("(?, ?, ?, ?, ?),")
		args = append(args, v.ProgramID, v.PartnerID, v.CreatedBy, v.UpdatedBy, StatusCreated)
	}

	q := `INSERT INTO 
				program_partners
				( 				
					 program_id
					, partner_id	
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
				( program_id , partner_id ) 
			DO UPDATE 
			SET
				program_id = excluded.program_id
				,partner_id = excluded.partner_id
				,updated_at = now()
				,status = excluded.status
			RETURNING
				id
				, program_id
				, partner_id
				, created_at
				, created_by
				, updated_at
				, updated_by
				, status
	`
	// fmt.Println("pp query : ", q)

	var res ProgramPartners
	err := tx.Select(&res, tx.Rebind(q), args...)
	if err != nil {
		return err
	}

	*pp = res

	return nil
}

// Delete : delete program partners by program ID
func (pp *ProgramPartners) Delete(tx *sqlx.Tx) error {

	q := `UPDATE
				program_partners 
			SET
				updated_at = now()
				, status = ?
			WHERE 
				program_id = ?	
			RETURNING
			id
			, program_id
			, partner_id
			, created_at
			, created_by
			, updated_at
			, updated_by
			, status
	`

	programID := (*pp)[0].ProgramID

	var res ProgramPartners
	err := tx.Select(&res, tx.Rebind(q), StatusDeleted, programID)
	if err != nil {
		return errors.New("Failed when delete program partner ," + ErrorNoDataAffected.Error() + " , " + err.Error())
	}

	*pp = res

	return nil
}
