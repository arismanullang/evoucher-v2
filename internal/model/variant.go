package model

import (
	"strings"
	"time"
	//"fmt"
)

type (
	Variant struct {
		ID                string   `db:"id"`
		CompanyID         string   `db:"company_id"`
		VariantName       string   `db:"variant_name"`
		VariantType       string   `db:"variant_type"`
		PointNeeded       float64  `db:"point_needed"`
		MaxVoucher        float64  `db:"max_voucher"`
		AllowAccumulative bool     `db:"allow_accumulative"`
		StartDate         string   `db:"start_date"`
		FinishDate        string   `db:"end_date"`
		DiscountValue     float64  `db:"discount_value"`
		ImgUrl            string   `db:"img_url"`
		VariantTnc        string   `db:"variant_tnc"`
		User              string   `db:"created_by"`
		CreatedAt         string   `db:"created_at"`
		Status            string   `db:"status"`
		ValidUsers        []string `db:"-"`
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
	VariantsResponse struct {
		Status       string
		Message      string
		VariantValue []SearchVariant
	}
	SearchVariant struct {
		ID          string   `db:"id"`
		CompanyID   string   `db:"company_id"`
		VariantName string   `db:"variant_name"`
		PointNeeded float64  `db:"point_needed"`
		MaxVoucher  float64  `db:"max_voucher"`
		StartDate   string   `db:"start_date"`
		FinishDate  string   `db:"end_date"`
		VariantTnc  string   `db:"variant_tnc"`
		ValidUsers  []string `db:"-"`
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
			, start_date
			, end_date
			, discount_value
			, img_url
			, variant_tnc
			, created_by
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ? ,?)
		RETURNING
			id
	`
	var res []string
	if err := tx.Select(&res, tx.Rebind(q), d.CompanyID, d.VariantName, d.VariantType, d.PointNeeded, d.MaxVoucher, d.AllowAccumulative, d.StartDate, d.FinishDate, d.DiscountValue, d.ImgUrl, d.VariantTnc, d.User); err != nil {
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
			variant_name = ?
			, variant_type = ?
			, point_needed = ?
			, max_voucher = ?
			, allow_accumulative = ?
			, start_date = ?
			, end_date = ?
			, discount_value = ?
			, img_url = ?
			, variant_tnc = ?
			, updated_by = ?
			, updated_at = ?
		WHERE
			id = ?;
	`

	_, err = tx.Exec(tx.Rebind(q), d.VariantName, d.VariantType, d.PointNeeded, d.MaxVoucher, d.AllowAccumulative, d.StartDate, d.FinishDate, d.DiscountValue, d.ImgUrl, d.VariantTnc, d.User, time.Now(), d.ID)
	if err != nil {
		return err
	}
	/*
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
	*/
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

	_, err = tx.Exec(tx.Rebind(q), d.User, time.Now(), StatusDeleted, d.ID)
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
	_, err = tx.Exec(tx.Rebind(q), d.User, time.Now(), StatusDeleted, d.ID)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func FindVariantByID(id string) (VariantResponse, error) {
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
			, discount_value
			, img_url
			, variant_tnc
			, created_by
			, created_at
		FROM
			variants
		WHERE
			id = ?
			AND status = ?
	`

	var resv []Variant
	if err := db.Select(&resv, db.Rebind(q), id, StatusCreated); err != nil {
		return VariantResponse{Status: "Error", Message: q, VariantValue: Variant{}}, err
	}
	if len(resv) < 1 {
		return VariantResponse{Status: "404", Message: q, VariantValue: Variant{}}, ErrResourceNotFound
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
	if err := db.Select(&resd, db.Rebind(q), id, StatusCreated); err != nil {
		return VariantResponse{Status: "Error", Message: q, VariantValue: Variant{}}, err
	}
	resv[0].ValidUsers = resd

	res := VariantResponse{
		Status:       "200",
		Message:      "Ok",
		VariantValue: resv[0],
	}

	return res, nil
}

func FindVariantByName(name string) (VariantsResponse, error) {
	q := `
		SELECT
			id
			, variant_name
			, point_needed
			, max_voucher
			, start_date
			, end_date
			, variant_tnc
		FROM
			variants
		WHERE
			variant_name LIKE ?
			AND status = ?
	`

	var resv []SearchVariant
	if err := db.Select(&resv, db.Rebind(q), name, StatusCreated); err != nil {
		return VariantsResponse{Status: "500", Message: "Error at select variant", VariantValue: []SearchVariant{}}, err
	}
	if len(resv) < 1 {
		return VariantsResponse{Status: "404", Message: "Error at select variant", VariantValue: []SearchVariant{}}, ErrResourceNotFound
	}

	for i, v := range resv {
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
		if err := db.Select(&resd, db.Rebind(q), v.ID, StatusCreated); err != nil {
			return VariantsResponse{Status: "500", Message: "Error at select user", VariantValue: []SearchVariant{}}, err
		}
		resv[i].ValidUsers = resd
	}

	res := VariantsResponse{
		Status:       "200",
		Message:      "Ok",
		VariantValue: resv,
	}

	return res, nil
}

func FindVariantByUser(user string) (VariantsResponse, error) {
	q := `
		SELECT
			id
			, variant_name
			, point_needed
			, max_voucher
			, start_date
			, end_date
			, variant_tnc
		FROM
			variants
		WHERE
			created_by = ?
			AND status = ?
	`

	var resv []SearchVariant
	if err := db.Select(&resv, db.Rebind(q), user, StatusCreated); err != nil {
		return VariantsResponse{Status: "500", Message: "Error at select variant", VariantValue: []SearchVariant{}}, err
	}
	if len(resv) < 1 {
		return VariantsResponse{Status: "404", Message: "Error at select variant", VariantValue: []SearchVariant{}}, ErrResourceNotFound
	}

	for i, v := range resv {
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
		if err := db.Select(&resd, db.Rebind(q), v.ID, StatusCreated); err != nil {
			return VariantsResponse{Status: "500", Message: "Error at select user", VariantValue: []SearchVariant{}}, err
		}
		resv[i].ValidUsers = resd
	}

	res := VariantsResponse{
		Status:       "200",
		Message:      "Ok",
		VariantValue: resv,
	}

	return res, nil
}

func FindVariantByCompanyID(id string) (VariantsResponse, error) {
	q := `
		SELECT
			id
			, variant_name
			, point_needed
			, max_voucher
			, start_date
			, end_date
			, variant_tnc
		FROM
			variants
		WHERE
			company_id = ?
			AND status = ?
	`

	var resv []SearchVariant
	if err := db.Select(&resv, db.Rebind(q), id, StatusCreated); err != nil {
		return VariantsResponse{Status: "500", Message: "Error at select variant", VariantValue: []SearchVariant{}}, err
	}
	if len(resv) < 1 {
		return VariantsResponse{Status: "404", Message: "Error at select variant", VariantValue: []SearchVariant{}}, ErrResourceNotFound
	}

	for i, v := range resv {
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
		if err := db.Select(&resd, db.Rebind(q), v.ID, StatusCreated); err != nil {
			return VariantsResponse{Status: "500", Message: "Error at select user", VariantValue: []SearchVariant{}}, err
		}
		resv[i].ValidUsers = resd
	}

	res := VariantsResponse{
		Status:       "200",
		Message:      "Ok",
		VariantValue: resv,
	}

	return res, nil
}

func FindVariantByDate(date string) (VariantsResponse, error) {
	dateSplit := strings.Split(date, ";")
	q := `
		SELECT
			id
			, variant_name
			, point_needed
			, max_voucher
			, start_date
			, end_date
			, variant_tnc
		FROM
			variants
		WHERE
			created_at > ?
			AND created_at < ?
			AND status = ?
	`

	var resv []SearchVariant
	if err := db.Select(&resv, db.Rebind(q), dateSplit[0], dateSplit[1], StatusCreated); err != nil {
		return VariantsResponse{Status: "500", Message: "Error at select variant", VariantValue: []SearchVariant{}}, err
	}
	if len(resv) < 1 {
		return VariantsResponse{Status: "404", Message: "Error at select variant", VariantValue: []SearchVariant{}}, ErrResourceNotFound
	}

	for i, v := range resv {
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
		if err := db.Select(&resd, db.Rebind(q), v.ID, StatusCreated); err != nil {
			return VariantsResponse{Status: "500", Message: "Error at select user", VariantValue: []SearchVariant{}}, err
		}
		resv[i].ValidUsers = resd
	}

	res := VariantsResponse{
		Status:       "200",
		Message:      "Ok",
		VariantValue: resv,
	}

	return res, nil
}

func FindAllVariants() (VariantsResponse, error) {
	q := `
		SELECT
			id
			, variant_name
			, point_needed
			, max_voucher
			, start_date
			, end_date
			, variant_tnc
		FROM
			variants
	`

	var resv []SearchVariant
	if err := db.Select(&resv, db.Rebind(q)); err != nil {
		return VariantsResponse{Status: "500", Message: "Error at select variant", VariantValue: []SearchVariant{}}, err
	}
	if len(resv) < 1 {
		return VariantsResponse{Status: "404", Message: "Error at select variant", VariantValue: []SearchVariant{}}, ErrResourceNotFound
	}

	for i, v := range resv {
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
		if err := db.Select(&resd, db.Rebind(q), v.ID, StatusCreated); err != nil {
			return VariantsResponse{Status: "500", Message: "Error at select user", VariantValue: []SearchVariant{}}, err
		}
		resv[i].ValidUsers = resd
	}

	res := VariantsResponse{
		Status:       "200",
		Message:      "Ok",
		VariantValue: resv,
	}

	return res, nil
}

func FindVariantMultipleParam(field, value []string) (VariantsResponse, error) {
	q := `
		SELECT
			id
			, variant_name
			, point_needed
			, max_voucher
			, start_date
			, end_date
			, variant_tnc
		FROM
			variants
		WHERE
			status = ?
	`
	for i := 0; i < len(field); i++ {
		q += ` AND ` + field[i] + ` = '` + value[i] + `'`
	}

	var resv []SearchVariant
	if err := db.Select(&resv, db.Rebind(q), StatusCreated); err != nil {
		return VariantsResponse{Status: "Error", Message: q, VariantValue: []SearchVariant{}}, err
	}
	if len(resv) < 1 {
		return VariantsResponse{Status: "404", Message: q, VariantValue: []SearchVariant{}}, ErrResourceNotFound
	}

	for i, v := range resv {
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
		if err := db.Select(&resd, db.Rebind(q), v.ID, StatusCreated); err != nil {
			return VariantsResponse{Status: "500", Message: "Error at select user", VariantValue: []SearchVariant{}}, err
		}
		resv[i].ValidUsers = resd
	}

	res := VariantsResponse{
		Status:       "200",
		Message:      "Ok",
		VariantValue: resv,
	}

	return res, nil
}
