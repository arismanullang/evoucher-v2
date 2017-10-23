package model

import (
	"database/sql"
	"log"
	"time"

	"github.com/lib/pq"
)

type (
	Voucher struct {
		ID                string         `db:"id" json:"id"`
		VoucherCode       string         `db:"voucher_code" json:"voucher_code"`
		ReferenceNo       string         `db:"reference_no" json:"reference_no"`
		Holder            sql.NullString `db:"holder" json:"holder"`
		HolderPhone       sql.NullString `db:"holder_phone" json:"holder_phone"`
		HolderEmail       sql.NullString `db:"holder_email" json:"holder_email"`
		HolderDescription sql.NullString `db:"holder_description" json:"holder_description"`
		ProgramID         string         `db:"program_id" json:"program_id"`
		ProgramName       string         `db:"program_name" json:"program_name"`
		ValidAt           time.Time      `db:"valid_at" json:"valid_at"`
		ExpiredAt         time.Time      `db:"expired_at" json:"expired_at"`
		VoucherValue      float64        `db:"voucher_value" json:"voucher_value"`
		State             string         `db:"state" json:"state"`
		CreatedBy         string         `db:"created_by" json:"created_by"`
		CreatedAt         time.Time      `db:"created_at" json:"created_at"`
		UpdatedBy         sql.NullString `db:"updated_by" json:"updated_by"`
		UpdatedAt         pq.NullTime    `db:"updated_at" json:"updated_at"`
		DeletedBy         sql.NullString `db:"deleted_by" json:"deleted_by"`
		DeletedAt         pq.NullTime    `db:"deleted_at" json:"deleted_at"`
		Status            string         `db:"status" json:"status"`
	}

	VoucherResponse struct {
		Status      string
		Message     string
		VoucherData []Voucher
	}
	UpdateDeleteRequest struct {
		ID          string `db:"id"`
		VoucherCode string `db:"voucher_code"`
		State       string `db:"state"`
		User        string `db:"created_by"`
	}
	VoucherCodeFormat struct {
		Prefix     sql.NullString `db:"prefix"`
		Postfix    sql.NullString `db:"postfix"`
		Body       sql.NullString `db:"body"`
		FormatType string         `db:"format_type"`
		Length     int            `db:"length"`
	}
)

func FindAvailableVoucher(accountId string, param map[string]string) (VoucherResponse, error) {
	q := `
		SELECT
			v.id
			, v.voucher_code
			, v.reference_no
			, v.holder
			, v.holder_phone
			, v.holder_email
			, v.holder_description
			, v.program_id
			, v.valid_at
			, v.expired_at
			, v.voucher_value
			, v.state
			, v.created_by
			, v.created_at
			, v.updated_by
			, v.updated_at
			, v.deleted_by
			, v.deleted_at
			, v.status
		FROM
			vouchers as v
		JOIN
			programs as p
		ON
			v.program_id = p.id
		WHERE
			v.status = ?
			AND p.account_id = ?
			AND v.expired_at > now()
			AND v.valid_at < now()
			AND v.state = 'created'
	`
	for key, value := range param {
		q += ` AND v.` + key + ` = '` + value + `'`
	}
	q += ` ORDER BY v.state DESC`

	var resd []Voucher
	if err := db.Select(&resd, db.Rebind(q), StatusCreated, accountId); err != nil {
		return VoucherResponse{Status: ResponseStateNok, Message: err.Error(), VoucherData: resd}, err
	}
	if len(resd) < 1 {
		return VoucherResponse{Status: ResponseStateNok, Message: ErrResourceNotFound.Error(), VoucherData: resd}, ErrResourceNotFound
	} else if resd[0].State != VoucherStateActived && resd[0].State != VoucherStateCreated {
		return VoucherResponse{Status: ErrCodeVoucherDisabled, Message: ErrMessageVoucherDisabled, VoucherData: resd}, nil
	} else if resd[0].ValidAt.After(time.Now()) {
		return VoucherResponse{Status: ErrCodeVoucherNotActive, Message: ErrMessageVoucherNotActive, VoucherData: resd}, nil
	} else if resd[0].ExpiredAt.Before(time.Now()) {
		return VoucherResponse{Status: ErrCodeVoucherExpired, Message: ErrMessageVoucherExpired, VoucherData: resd}, nil
	}

	res := VoucherResponse{Status: ResponseStateOk, Message: "success", VoucherData: resd}
	return res, nil
}

func FindVoucher(param map[string]string) (VoucherResponse, error) {
	q := `
		SELECT
			v.id
			, v.voucher_code
			, v.reference_no
			, v.holder
			, v.holder_phone
			, v.holder_email
			, v.holder_description
			, v.program_id
			, p.name as program_name
			, v.valid_at
			, v.expired_at
			, v.voucher_value
			, v.state
			, v.created_by
			, v.created_at
			, v.updated_by
			, v.updated_at
			, v.deleted_by
			, v.deleted_at
			, v.status
		FROM
			vouchers as v
		JOIN
			programs as p
		ON
			v.program_id = p.id
		WHERE
			v.status = ?
	`
	for key, value := range param {
		if key == "holder" {
			q += ` AND (LOWER(v.holder) LIKE LOWER('%` + value + `%')`
			q += ` OR LOWER(v.holder_description) LIKE LOWER('%` + value + `%'))`
		} else {
			q += ` AND v.` + key + ` = '` + value + `'`
		}
	}
	q += ` ORDER BY state DESC`

	var resd []Voucher
	if err := db.Select(&resd, db.Rebind(q), StatusCreated); err != nil {
		return VoucherResponse{Status: ResponseStateNok, Message: err.Error(), VoucherData: resd}, err
	}
	if len(resd) < 1 {
		return VoucherResponse{Status: ResponseStateNok, Message: ErrResourceNotFound.Error(), VoucherData: resd}, ErrResourceNotFound
	} else if resd[0].State != VoucherStateActived && resd[0].State != VoucherStateCreated {
		return VoucherResponse{Status: ErrCodeVoucherDisabled, Message: ErrMessageVoucherDisabled, VoucherData: resd}, nil
	} else if resd[0].ValidAt.After(time.Now()) {
		return VoucherResponse{Status: ErrCodeVoucherNotActive, Message: ErrMessageVoucherNotActive, VoucherData: resd}, nil
	} else if resd[0].ExpiredAt.Before(time.Now()) {
		return VoucherResponse{Status: ErrCodeVoucherExpired, Message: ErrMessageVoucherExpired, VoucherData: resd}, nil
	}

	res := VoucherResponse{Status: ResponseStateOk, Message: "success", VoucherData: resd}
	return res, nil
}

func FindVouchers(param map[string]string) (VoucherResponse, error) {
	q := `
		SELECT DISTINCT
			v.id
			, v.voucher_code
			, v.reference_no
			, v.holder
			, v.holder_phone
			, v.holder_email
			, v.holder_description
			, v.program_id
			, pr.name as program_name
			, v.valid_at
			, v.expired_at
			, v.voucher_value
			, v.state
			, v.created_by
			, v.created_at
			, v.updated_by
			, v.updated_at
			, v.deleted_by
			, v.deleted_at
			, v.status
		FROM
			vouchers as v
		JOIN
			programs as pr
		ON
			v.program_id = pr.id
		JOIN
			program_partners as pp
		ON
			pp.program_id = pr.id
		JOIN
			partners as pa
		ON
			pp.partner_id = pa.id
		JOIN
			transactions as t
		ON
			t.partner_id = pa.id
		WHERE
			v.status = ?
	`
	for key, value := range param {
		if key == "holder" {
			q += ` AND (LOWER(v.holder) LIKE LOWER('%` + value + `%')`
			q += ` OR LOWER(v.holder_description) LIKE LOWER('%` + value + `%'))`
		} else {
			q += ` AND ` + key + ` = '` + value + `'`
		}
	}
	q += ` ORDER BY state DESC`

	var resd []Voucher
	if err := db.Select(&resd, db.Rebind(q), StatusCreated); err != nil {
		return VoucherResponse{Status: ResponseStateNok, Message: err.Error(), VoucherData: resd}, err
	}
	if len(resd) < 1 {
		return VoucherResponse{Status: ResponseStateNok, Message: ErrResourceNotFound.Error(), VoucherData: resd}, ErrResourceNotFound
	} else if resd[0].State != VoucherStateActived && resd[0].State != VoucherStateCreated {
		return VoucherResponse{Status: ErrCodeVoucherDisabled, Message: ErrMessageVoucherDisabled, VoucherData: resd}, nil
	} else if resd[0].ValidAt.After(time.Now()) {
		return VoucherResponse{Status: ErrCodeVoucherNotActive, Message: ErrMessageVoucherNotActive, VoucherData: resd}, nil
	} else if resd[0].ExpiredAt.Before(time.Now()) {
		return VoucherResponse{Status: ErrCodeVoucherExpired, Message: ErrMessageVoucherExpired, VoucherData: resd}, nil
	}

	res := VoucherResponse{Status: ResponseStateOk, Message: "success", VoucherData: resd}
	return res, nil
}

func (d *Voucher) InsertVc() error {
	vc, err := db.Beginx()
	if err != nil {
		return err
	}
	defer vc.Rollback()

	q := `
	      INSERT INTO
	      	vouchers (
			      voucher_code
			      , reference_no
			      , holder
			      , holder_phone
			      , holder_email
			      , holder_description
			      , program_id
			      , valid_at
			      , expired_at
			      , voucher_value
			      , state
			      , created_by
	      		)
	      	 VALUES ( ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	      RETURNING
			      id
			      , voucher_code
			      , reference_no
			      , holder
			      , holder_phone
			      , holder_email
			      , holder_description
			      , program_id
			      , valid_at
			      , expired_at
			      , voucher_value
			      , state
			      , created_by
			      , created_at
			      , updated_by
			      , updated_at
			      , deleted_by
			      , deleted_at
			      , status
      `
	var res []Voucher
	// fmt.Println("insert data =>", d)
	if err := vc.Select(&res, vc.Rebind(q), d.VoucherCode, d.ReferenceNo, d.Holder, d.HolderPhone, d.HolderEmail, d.HolderDescription, d.ProgramID, d.ValidAt, d.ExpiredAt, d.VoucherValue, VoucherStateCreated, d.CreatedBy); err != nil {
		return err
	}

	*d = res[0]

	if err := vc.Commit(); err != nil {
		return err
	}
	return nil
}

func (d *UpdateDeleteRequest) DeleteVc() error {
	vc, err := db.Beginx()
	if err != nil {
		return err
	}
	defer vc.Rollback()

	q := `
		UPDATE 	vouchers
		SET
			state = ?
			, status = ?
			, deleted_by = ?
			, deleted_at = ?
		WHERE
			id = ?
		AND status = ?
		RETURNING id
      `
	var result []string
	if err := vc.Select(&result, vc.Rebind(q), d.State, StatusDeleted, d.User, time.Now(), d.ID, StatusCreated); err != nil {
		return err
	}

	if len(result) < 1 {
		return ErrNotModified
	}

	if err := vc.Commit(); err != nil {
		return err
	}
	return nil
}

func (d *UpdateDeleteRequest) UpdateVc() (Voucher, error) {
	vc, err := db.Beginx()
	if err != nil {
		return Voucher{}, err
	}
	defer vc.Rollback()

	q := `
		UPDATE 	vouchers
		SET
			state = ?
			, updated_by = ?
			, updated_at = ?
		WHERE
			id = ?
		RETURNING
			id
			, voucher_code
			, reference_no
			, holder
			, program_id
			, valid_at
			, expired_at
			, voucher_value
			, state
			, created_by
			, created_at
			, updated_by
			, updated_at
			, deleted_by
			, deleted_at
			, status
	`

	var result []Voucher
	if err := vc.Select(&result, vc.Rebind(q), d.State, d.User, time.Now(), d.ID); err != nil {
		return Voucher{}, err
	}

	if err := vc.Commit(); err != nil {
		return Voucher{}, err
	}
	return result[0], nil
}

func GetVoucherCodeFormat(id int) (VoucherCodeFormat, error) {
	vc, err := db.Beginx()
	if err != nil {
		return VoucherCodeFormat{}, err
	}
	defer vc.Rollback()

	q := `
		SELECT
			prefix
			, postfix
			, body
			, format_type
			, length
		FROM
			voucher_formats
		WHERE
			id = ?
			AND status = ?
	`

	var resd []VoucherCodeFormat
	if err := db.Select(&resd, db.Rebind(q), id, StatusCreated); err != nil {
		log.Panic(err)
		return VoucherCodeFormat{}, err
	}
	return resd[0], nil
}

func CountVoucher(varID string) int {
	vc, err := db.Beginx()
	if err != nil {
		return 0
	}
	defer vc.Rollback()

	q := `
		SELECT
			count(1)
		FROM
			vouchers
		WHERE
			program_id = ?
			AND status = ?
	`
	var resd []int
	if err := db.Select(&resd, db.Rebind(q), varID, StatusCreated); err != nil {
		log.Panic(err)
		return 0
	}
	return resd[0]
}

func CountHolderVoucher(programId, holder string) int {
	vc, err := db.Beginx()
	if err != nil {
		return 0
	}
	defer vc.Rollback()

	q := `
		SELECT
			count(1)
		FROM
			vouchers
		WHERE
			program_id = ?
			AND holder = ?
			AND status = ?
	`
	var resd []int
	if err := db.Select(&resd, db.Rebind(q), programId, holder, StatusCreated); err != nil {
		log.Panic(err)
		return 0
	}
	return resd[0]
}

func HardDelete(program string) error {
	vc, err := db.Beginx()
	if err != nil {
		return err
	}
	defer vc.Rollback()

	q := `
		DELETE 	FROM
			vouchers
		WHERE
			program_id = ?
		AND
			status = ?
		RETURNING id
      `
	var result []string
	if err := vc.Select(&result, vc.Rebind(q), program, StatusCreated); err != nil {
		return err
	}

	if len(result) < 1 {
		return ErrNotModified
	}

	if err := vc.Commit(); err != nil {
		return err
	}
	return nil
}

func RollbackVoucher(vcid string) error {
	vc, err := db.Beginx()
	if err != nil {
		return err
	}
	defer vc.Rollback()

	q := `
		DELETE FROM
			vouchers
		WHERE
			id = ?
		AND
			status = ?
		RETURNING id
      `
	var result []string
	if err := vc.Select(&result, vc.Rebind(q), vcid, StatusCreated); err != nil {
		return err
	}

	if len(result) < 1 {
		return ErrNotModified
	}

	if err := vc.Commit(); err != nil {
		return err
	}
	return nil
}
