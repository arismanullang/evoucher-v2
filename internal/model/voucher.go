package model

import (
	"database/sql"
	"log"
	"time"

	"github.com/lib/pq"
)

type (
	Voucher struct {
		ID            string         `db:"id"`
		VoucherCode   string         `db:"voucher_code"`
		ReferenceNo   string         `db:"reference_no"`
		Holder        string         `db:"holder"`
		VariantID     string         `db:"variant_id"`
		ValidAt       time.Time      `db:"valid_at"`
		ExpiredAt     time.Time      `db:"expired_at"`
		DiscountValue float64        `db:"discount_value"`
		State         string         `db:"state"`
		CreatedBy     string         `db:"created_by"`
		CreatedAt     time.Time      `db:"created_at"`
		UpdatedBy     sql.NullString `db:"updated_by"`
		UpdatedAt     pq.NullTime    `db:"updated_at"`
		DeletedBy     sql.NullString `db:"deleted_by"`
		DeletedAt     pq.NullTime    `db:"deleted_at"`
		Status        string         `db:"status"`
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

func FindVoucher(param map[string]string) (VoucherResponse, error) {
	q := `
		SELECT
			id
			, voucher_code
			, reference_no
			, holder
			, variant_id
			, valid_at
			, expired_at
			, discount_value
			, state
			, created_by
			, created_at
			, updated_by
			, updated_at
			, deleted_by
			, deleted_at
			, status
		FROM
			vouchers
		WHERE	status = ?
	`
	for key, value := range param {
		q += ` AND ` + key + ` = '` + value + `'`
	}

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
			      , variant_id
			      , valid_at
			      , expired_at
			      , discount_value
			      , state
			      , created_by
	      		)
	      	 VALUES ( ?, ?, ?, ?, ?, ?, ?, ?, ?)
	      RETURNING
			      id
			      , voucher_code
			      , reference_no
			      , holder
			      , variant_id
			      , valid_at
			      , expired_at
			      , discount_value
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
	if err := vc.Select(&res, vc.Rebind(q), d.VoucherCode, d.ReferenceNo, d.Holder, d.VariantID, d.ValidAt, d.ExpiredAt, d.DiscountValue, VoucherStateCreated, d.CreatedBy); err != nil {
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
			, variant_id
			, valid_at
			, expired_at
			, discount_value
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
			variant_id = ?
			AND status = ?
	`
	var resd []int
	if err := db.Select(&resd, db.Rebind(q), varID, StatusCreated); err != nil {
		log.Panic(err)
		return 0
	}
	return resd[0]
}
