package model

import (
	"bytes"
	"github.com/gilkor/evoucher/util"
	"time"
)

type (
	//Voucher model
	Voucher struct {
		ID           string     `json:"id,omitempty" db:"id"`
		Code         string     `json:"code,omitempty" db:"code"`
		ReferenceNo  string     `json:"reference_no,omitempty" db:"reference_no"`
		Holder       *string    `json:"holder,omitempty" db:"holder"`
		HolderDetail JSONExpr   `json:"holder_detail,omitempty" db:"holder_detail"`
		ProgramID    string     `json:"program_id,omitempty" db:"program_id"`
		ValidAt      *time.Time `json:"valid_at,omitempty" db:"valid_at"`
		ExpiredAt    *time.Time `json:"expired_at,omitempty" db:"expired_at"`
		State        string     `json:"state,omitempty" db:"state"`
		CreatedBy    string     `json:"created_by,omitempty" db:"created_by"`
		CreatedAt    *time.Time `json:"created_at,omitempty" db:"created_at"`
		UpdatedBy    *string    `json:"updated_by,omitempty" db:"updated_by"`
		UpdatedAt    *time.Time `json:"updated_at,omitempty" db:"updated_at"`
		Status       string     `json:"status,omitempty" db:"status"`
	}
	//Vouchers :
	Vouchers []Voucher
	//HolderDetail :type struct Voucher.JSONExpr.Unmarshal(&HolderDetail)
	HolderDetail struct {
		Name        string `json:"holder_name,omitempty"`
		Phone       string `json:"holder_phone,omitempty"`
		Email       string `json:"holder_email,omitempty"`
		Description string `json:"holder_description,omitempty"`
	}
)

// GetVoucherByID :  get list vouchers by ID
func GetVoucherByID(id string, qp *util.QueryParam) (*Vouchers, bool, error) {
	return getVouchers("id", id, qp)
}

// GetVouchers : list voucher
func GetVouchers(qp *util.QueryParam) (*Vouchers, bool, error) {
	return getVouchers("1", "1", qp)
}

// GetVouchersByProgramID : get list vouchers by program.ID
func GetVouchersByProgramID(programID string, qp *util.QueryParam) (*Vouchers, bool, error) {
	return getVouchers("program_id", programID, qp)
}

func getVouchers(key, value string, qp *util.QueryParam) (*Vouchers, bool, error) {
	q, err := qp.GetQueryByDefaultStruct(Voucher{})
	if err != nil {
		return &Vouchers{}, false, err
	}
	q += `
			FROM
				vouchers voucher
			WHERE 
				status = ?			
			AND ` + key + ` = ?`

	q += qp.GetQuerySort()
	q += qp.GetQueryLimit()
	var resd Vouchers
	err = db.Select(&resd, db.Rebind(q), StatusCreated, value)
	if err != nil {
		return &Vouchers{}, false, err
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

//Insert : single row insert into table
func (v Voucher) Insert() error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	q := `INSERT INTO 
				vouchers 
				( 
					code  			= ?
					, reference_no  = ?
					, holder        = ?
					, holder_detail = ?
					, program_id    = ?
					, valid_at      = ?
					, expired_at    = ?
					, state         = ?
					, created_by    = ?
					, created_at    = ?
					, updated_by    = ?
					, updated_at    = ?
					, status        = ?
				)
			RETURNING
			id
			, code
			, reference_no
			, holder
			, holder_detail
			, program_id
			, valid_at
			, expired_at
			, state
			, created_by
			, created_at
			, updated_by
			, updated_at
			, status
	`
	var res Voucher
	t1 := time.Now()
	err = tx.Select(&res, tx.Rebind(q),
		v.Code,
		v.ReferenceNo,
		v.Holder,
		v.HolderDetail,
		v.ProgramID,
		v.ValidAt,
		v.ExpiredAt,
		v.State,
		v.CreatedBy,
		t1,
		v.UpdatedBy,
		t1,
		StatusCreated,
	)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

//Update : update voucher
func (v *Voucher) Update() error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	t1 := time.Now()
	q := `UPDATE
				vouchers 
			SET
				holder = ?,
				holder_detail = ?,
				state = ?,
				updated_at = ?,
				updated_by = ?,					
				status = ?
			WHERE 
				id = ?
			RETURNING
				id
				, code
				, reference_no
				, holder
				, holder_detail
				, program_id
				, valid_at
				, expired_at
				, state
				, created_by
				, created_at
				, updated_by
				, updated_at
				, status
	`
	var res []Voucher
	err = tx.Select(&res, tx.Rebind(q), v.Holder, v.HolderDetail, v.State, t1, v.UpdatedBy, v.Status)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

//Delete : soft deleted data by updating row status to "deleted"
func (v *Voucher) Delete() error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	t1 := time.Now()
	q := `UPDATE
				vouchers 
			SET
				updated_at = ?,
				updated_by = ?
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
				, status
	`
	var res []Voucher
	err = tx.Select(&res, tx.Rebind(q), t1, v.UpdatedBy, StatusDeleted)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

//Insert : insert data, build query using string append
func (vs *Vouchers) Insert() error {
	tx, err := db.Beginx()
	values := new(bytes.Buffer)
	var args []interface{}
	for _, v := range *vs {
		values.WriteString("(?, ?, ?, ?, ?, ?),")
		args = append(args, v.Code, v.ReferenceNo, VoucherStateCreated, v.CreatedBy, v.UpdatedBy, StatusCreated)
	}

	q := `INSERT INTO 
				vouchers
				( 				
					 code
					, reference_no
					, state
					, created_by
					, updated_by
					, status
				)
			VALUES 
			`
	valuestr := values.String()
	q += valuestr[:len(valuestr)-1]

	q += `
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
	var res Vouchers
	err = tx.Select(&res, tx.Rebind(q), args...)
	if err != nil {
		return err
	}

	*vs = res
	return nil
}
