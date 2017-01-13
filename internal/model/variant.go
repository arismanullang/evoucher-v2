package model

import (
	"time"
)

type (
	Variant struct {
		ID                string    `db:"id"`
		CompanyID         string    `db:"company_id"`
		VariantName       string    `db:"variant_name"`
		VariantType       string    `db:"variant_type"`
		PointNeeded       float64   `db:"point_needed"`
		MaxVoucher        float64   `db:"max_voucher"`
		AllowAccumulative bool      `db:"allow_accumulative"`
		StartDate         time.Time `db:"start_date"`
		FinishDate        time.Time `db:"end_date"`
		ImgUrl            string    `db:"img_url"`
		VariantTnc        string    `db:"variant_tnc"`
		User              string    `db:"created_by"`
		CreatedAt         string    `db:"created_at"`
		Status            string    `db:"status"`
		ValidUsers        []string  `db:"-"`
	}
	VariantResponse struct {
		Status       string
		Message      string
		VariantValue Variant
	}
	DeleteVariantRequest struct {
		ID   string `db:"id"`
		User string `db:"deleted_by"`
	}
)

func (d *Variant) Insert() error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	q := `
		INSERT INTO variants(
			company_id
			, variant_name
			, variant_type
			, point_needed
			, max_voucher
			, allow_accumulative
			, img_url
			, variant_tnc
			, created_by
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING
			id
	`
	var res []string
	if err := tx.Select(&res, tx.Rebind(q), d.CompanyID, d.VariantName, d.VariantType, d.PointNeeded, d.MaxVoucher, d.AllowAccumulative, d.ImgUrl, d.VariantTnc, d.User); err != nil {
		return err
	}
	d.ID = res[0]

	for _, v := range d.ValidUsers {
		q := `
			INSERT INTO variant_users(
				company_id
				, variant_id
				, account_id
				, created_by
			)
			VALUES (?, ?, ?, ?)
		`

		_, err := tx.Exec(tx.Rebind(q), d.CompanyID, d.ID, v, d.User)
		if err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (d *Variant) Update() error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	q := `
		UPDATE variants
		SET
			company_id = ?
			, variant_name = ?
			, variant_type = ?
			, point_needed = ?
			, max_voucher = ?
			, allow_accumulative = ?
			, img_url = ?
			, variant_tnc = ?
			, updated_by = ?
			, updated_at = ?
		WHERE
			id = ?;
	`

	_, err = tx.Exec(tx.Rebind(q), d.CompanyID, d.VariantName, d.VariantType, d.PointNeeded, d.MaxVoucher, d.AllowAccumulative, d.ImgUrl, d.VariantTnc, d.User, time.Now(), d.ID)
	if err != nil {
		return err
	}

	q = `
		DELETE FROM variant_users
		WHERE
			variant_id = ?;
	`
	_, err = tx.Exec(tx.Rebind(q), d.ID)
	if err != nil {
		return err
	}

	for _, v := range d.ValidUsers {
		q := `
			INSERT INTO variant_users (
				company_id
				, variant_id
				, account_id
				, created_by
			)
			VALUES (?, ?, ?, ?)
		`

		_, err := tx.Exec(tx.Rebind(q), d.CompanyID, d.ID, v, d.User)
		if err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (d *DeleteVariantRequest) Delete() error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	q := `
		UPDATE variants
		SET
			deleted_by = ?
			, deleted_at = ?
			, status = ?
		WHERE
			id = ?;
	`

	_, err = tx.Exec(tx.Rebind(q), d.User, time.Now(), "deleted", d.ID)
	if err != nil {
		return err
	}

	q = `
		UPDATE variant_users
		SET
			updated_by = ?
			, updated_at = ?
			, status = ?
		WHERE
			variant_id = ?;
	`
	_, err = tx.Exec(tx.Rebind(q), d.User, time.Now(), "deleted", d.ID)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func FindVariantByID(id string) (VariantResponse, error) {
	return findVariant("id", id)
}

func FindVariantByName(name string) (VariantResponse, error) {
	return findVariant("variant_name", name)
}

func findVariant(field, value string) (VariantResponse, error) {
	q := `
		SELECT
			id
			, company_id
			, variant_name
			, variant_type
			, point_needed
			, max_voucher
			, allow_accumulative
			, start_date
			, end_date
			, img_url
			, variant_tnc
			, created_by
			, created_at
		FROM
			variants
		WHERE
			` + field + ` = ?
			AND status = ?
	`

	var resv []Variant
	if err := db.Select(&resv, db.Rebind(q), value, StatusCreated); err != nil {
		return VariantResponse{Status: "Error", Message: "Error at select variant", VariantValue: Variant{}}, err
	}
	if len(resv) < 1 {
		return VariantResponse{Status: "404", Message: "Error at select variant", VariantValue: Variant{}}, ErrResourceNotFound
	}

	q = `
		SELECT
			account_id
		FROM
			variant_users
		WHERE
			variant_id = ?
			AND status = ?
	`
	var resd []string
	if err := db.Select(&resd, db.Rebind(q), value, StatusCreated); err != nil {
		return VariantResponse{Status: "Error", Message: "Error at select user", VariantValue: Variant{}}, err
	}
	if len(resd) < 1 {
		return VariantResponse{Status: "404", Message: "Error at select user", VariantValue: Variant{}}, ErrResourceNotFound
	}
	resv[0].ValidUsers = resd

	res := VariantResponse{
		Status:       "Ok",
		Message:      "Ok",
		VariantValue: resv[0],
	}

	return res, nil
}
