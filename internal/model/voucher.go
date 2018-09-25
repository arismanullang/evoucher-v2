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

	AssignGiftRequest struct {
		HolderID          string           `json:"holder_id" valid:"required~holder_id is required"`
		HolderEmail       string           `json:"holder_email" valid:"required~holder_email is required"`
		HolderPhone       string           `json:"holder_phone" valid:"required~holder_phone is required"`
		HolderDescription string           `json:"holder_description" valid:"required~holder_description is required"`
		ReferenceNo       string           `json:"reference_no" valid:"required~reference_no is required"`
		Data              []AssignGiftData `json:"data" valid:"required~data is required"`
		User              string           `json:"updated_by"`
	}

	UnassignGiftRequest struct {
		VoucherID string `json:"voucher_id"`
	}

	AssignGiftData struct {
		ProgramID  string    `json:"program_id" valid:"required~program_id is required`
		VoucherIDs []string  `json:"voucher_ids" valid:"required~voucher_ids is required`
		ValidAt    time.Time `json:"valid_at"`
		ExpiredAt  time.Time `json:"expired_at"`
	}

	VoucherCodeFormat struct {
		Prefix     sql.NullString `db:"prefix"`
		Postfix    sql.NullString `db:"postfix"`
		Body       sql.NullString `db:"body"`
		FormatType string         `db:"format_type"`
		Length     int            `db:"length"`
	}
	VoucherState struct {
		ID    string `db:"id"`
		State string `db:"state"`
	}

	GeneratePrivilegeRequest struct {
		CompanyID  string `json:"company_id"`
		MemberID   string `json:"member_id"`
		MemberName string `json:"member_name"`
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
	q += ` ORDER BY
		CASE p.voucher_type
			WHEN 'privilege' THEN 1
			ELSE 2
		END, v.expired_at ASC`

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

func FindVouchersById(param []string) (VoucherResponse, error) {
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
			AND (
	`
	for key, value := range param {
		if key != 0 {
			q += `OR `
		}
		q += `v.id ='` + value + `'`
	}
	q += `) ORDER BY state DESC`

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

func FindVouchersState(param []string) ([]VoucherState, error) {
	q := `
		SELECT
			v.id
			, v.state
		FROM
			vouchers as v
		WHERE
			v.status = ?
			AND (
	`
	for key, value := range param {
		if key != 0 {
			q += `OR `
		}
		q += `v.id ='` + value + `'`
	}
	q += `) ORDER BY state DESC`

	var resd []VoucherState
	if err := db.Select(&resd, db.Rebind(q), StatusCreated); err != nil {
		return []VoucherState{}, err
	}
	if len(resd) < 1 {
		return []VoucherState{}, ErrResourceNotFound
	}

	return resd, nil
}

func FindVouchersByPartner(partnerId string) ([]Voucher, error) {
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
		JOIN transaction_details as td ON v.id = td.voucher_id
		JOIN transactions as t ON td.transaction_id = t.id
		JOIN programs as p ON v.program_id = p.id
		WHERE
			v.status = ?
		AND 	t.partner_id = ?
		ORDER BY state DESC
	`

	var resd []Voucher
	if err := db.Select(&resd, db.Rebind(q), StatusCreated, partnerId); err != nil {
		return []Voucher{}, err
	}
	if len(resd) < 1 {
		return []Voucher{}, ErrResourceNotFound
	}

	return resd, nil
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

func FindsVouchers(param map[string]string) ([]Voucher, error) {
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
		return []Voucher{}, err
	}
	if len(resd) < 1 {
		return []Voucher{}, ErrResourceNotFound
	}

	return resd, nil
}

func GetGiftVouchers(giftType string, programID string) ([]Voucher, error) {
	q := `
		SELECT
			v.id
			, v.voucher_code
			, v.reference_no
			, v.holder
			, v.holder_phone
			, v.holder_email
			, v.holder_description
			, v.valid_at
			, v.expired_at
			, v.voucher_value
			, v.state
			, v.created_by
			, v.created_at
			, v.updated_by
			, v.updated_at
			, v.status
		FROM vouchers v
		WHERE program_id = ?
			AND status = ?
	`
	if giftType == VouchersGiftAssigned {
		q += `AND holder != ''`
	} else if giftType == VouchersGiftUnassigned {
		q += `AND holder = ''`
	}

	q += ` ORDER BY updated_at DESC`

	var resd []Voucher
	if err := db.Select(&resd, db.Rebind(q), programID, StatusCreated); err != nil {
		return []Voucher{}, err
	}
	if len(resd) < 1 {
		return []Voucher{}, ErrResourceNotFound
	}

	return resd, nil
}

func FindTodayVouchers(param map[string]string) ([]Voucher, error) {
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
			AND t.created_at = ?
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
	if err := db.Select(&resd, db.Rebind(q), StatusCreated, time.Now()); err != nil {
		return []Voucher{}, err
	}
	if len(resd) < 1 {
		return []Voucher{}, ErrResourceNotFound
	}

	return resd, nil
}

func UnassignVoucher(user, voucherID string) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	q :=
		`
		UPDATE vouchers
		SET
			holder = ''
			, holder_email =  ''
			, holder_phone =  ''
			, holder_description =  ''
			, reference_no =  ''
			, valid_at = '0001-01-01 00:00:00+00:00'
			, expired_at = '0001-01-01 00:00:00+00:00'
			, updated_by = ?
			, updated_at = ?
		WHERE
			id = ?
			AND state = 'created'
		RETURNING id
	`

	var result []string
	if err := tx.Select(&result, tx.Rebind(q), user, time.Now(), voucherID); err != nil {
		return err
	}

	if len(result) == 0 {
		return ErrResourceNotFound
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (d *AssignGiftRequest) AssignVoucher() (string, error) {
	tx, err := db.Beginx()
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	q := `
		UPDATE vouchers
		SET
			holder = '` + d.HolderID + `'
			, holder_email =  '` + d.HolderEmail + `'
			, holder_phone =  '` + d.HolderPhone + `'
			, holder_description =  '` + d.HolderDescription + `'
			, reference_no =  '` + d.ReferenceNo + `'
			, valid_at = ?
			, expired_at = ?
			, updated_by = ?
			, updated_at = ?
		WHERE
			id = ?
			AND holder = ''
			AND program_id = ?
		RETURNING id
	`

	var finalResult []string
	var voucherIDs []string

	for _, assignData := range d.Data {
		for _, voucherID := range assignData.VoucherIDs {
			var result []string
			//add voucherID to list voucherIDs
			voucherIDs = append(voucherIDs, voucherID)
			if err := tx.Select(&result, tx.Rebind(q), assignData.ValidAt, assignData.ExpiredAt, d.User, time.Now(), voucherID, assignData.ProgramID); err != nil {
				return "", err
			}

			//add result to finalResult for checking purpose
			if len(result) > 0 {
				finalResult = append(finalResult, result[0])
			} else if len(result) == 0 {
				return voucherID + "-", ErrResourceNotFound
			}
		}
	}

	if len(finalResult) != len(voucherIDs) {
		return "", ErrNotModified
	}

	if err := tx.Commit(); err != nil {
		return "", err
	}

	return "", nil
}

func (d *GeneratePrivilegeRequest) InsertPrivilegeVc() error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	q := `
		INSERT INTO vouchers(
			voucher_code
			, reference_no
			, holder
			, holder_email
			, holder_phone
			, holder_description
			, program_id
			, valid_at
			, expired_at
			, voucher_value
			, state
			, created_by
			, created_at
			, status
		)
		SELECT
			'PRIVILEGE'
			, 'privilege'
			, ?
			, ''
			, ''
			, ?
			, id
			, start_date
			, end_date
			, 0
			, 'created'
			, 'system'
			, ?
			, 'created'
		FROM
			programs
		WHERE
			account_id = (SELECT id FROM accounts WHERE company_id = ?)
			AND type = 'privilege'
			AND voucher_type = 'privilege'
			AND start_date< current_timestamp
			AND end_date > current_timestamp
			AND status = 'created'
			AND id NOT IN (
				SELECT program_id
				FROM
					vouchers
				WHERE
					holder = ?
			)
	`

	if _, err := tx.Exec(tx.Rebind(q), d.MemberID, d.MemberName, time.Now(), d.CompanyID, d.MemberID); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (d *Voucher) InsertVc() error {
	vc, err := db.Beginx()
	if err != nil {
		return err
	}
	defer vc.Rollback()

	q := `
		INSERT INTO vouchers (
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
			, created_at
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
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
	if err := vc.Select(&res, vc.Rebind(q), d.VoucherCode, d.ReferenceNo, d.Holder, d.HolderPhone, d.HolderEmail, d.HolderDescription, d.ProgramID, d.ValidAt, d.ExpiredAt, d.VoucherValue, VoucherStateCreated, d.CreatedBy, time.Now()); err != nil {
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

func RollbackVoucher(vcid, user string) error {
	vc, err := db.Beginx()
	if err != nil {
		return err
	}
	defer vc.Rollback()

	q := `
		UPDATE
			vouchers
		SET
			state = ?
			, status = ?
		WHERE
			id = ?
		AND
			status = ?
		RETURNING id
	`
	var result []string
	if err := vc.Select(&result, vc.Rebind(q), VoucherStateRollback, StatusDeleted, vcid, StatusCreated); err != nil {
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
