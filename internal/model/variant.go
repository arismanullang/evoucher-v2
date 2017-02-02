package model

import (
	"strconv"
	"time"
	//"fmt"
)

type (
	Variant struct {
		ID                 string   `db:"id"`
		CompanyID          string   `db:"company_id"`
		VariantName        string   `db:"variant_name"`
		VariantType        string   `db:"variant_type"`
		VoucherType        string   `db:"voucher_type"`
		PointNeeded        float64  `db:"point_needed"`
		MaxQuantityVoucher float64  `db:"max_quantity_voucher"`
		MaxUsageVoucher    float64  `db:"max_usage_voucher"`
		AllowAccumulative  bool     `db:"allow_accumulative"`
		RedeemtionMethod   string   `db:"redeemtion_method"`
		StartDate          string   `db:"start_date"`
		EndDate            string   `db:"end_date"`
		DiscountValue      float64  `db:"discount_value"`
		ImgUrl             string   `db:"img_url"`
		VariantTnc         string   `db:"variant_tnc"`
		User               string   `db:"created_by"`
		CreatedAt          string   `db:"created_at"`
		Status             string   `db:"status"`
		BlastUsers         []string `db:"-"`
		ValidTenants       []string `db:"-"`
	}
	DeleteVariantRequest struct {
		ID   string `db:"id"`
		User string `db:"deleted_by"`
	}
	VariantResponse struct {
		Status  string
		Message string
		Data    interface{}
	}
	SearchVariant struct {
		ID          string   `db:"id"`
		CompanyID   string   `db:"company_id"`
		VariantName string   `db:"variant_name"`
		VoucherType string   `db:"voucher_type"`
		PointNeeded float64  `db:"point_needed"`
		MaxVoucher  float64  `db:"max_quantity_voucher"`
		StartDate   string   `db:"start_date"`
		EndDate     string   `db:"end_date"`
		ValidUsers  []string `db:"-"`
	}
	UpdateVariantUsersRequest struct {
		ID        string   `db:"id"`
		CompanyID string   `db:"company_id"`
		User      string   `db:"updated_by"`
		Data      []string `db:"-"`
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
			, voucher_type
			, point_needed
			, max_quantity_voucher
			, max_usage_voucher
			, allow_accumulative
			, redeemtion_method
			, start_date
			, end_date
			, discount_value
			, img_url
			, variant_tnc
			, created_by
			, status
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING
			id
	`
	var res []string
	if err := tx.Select(&res, tx.Rebind(q), d.CompanyID, d.VariantName, d.VariantType, d.VoucherType, d.PointNeeded, d.MaxQuantityVoucher, d.MaxUsageVoucher, d.AllowAccumulative, d.RedeemtionMethod, d.StartDate, d.EndDate, d.DiscountValue, d.ImgUrl, d.VariantTnc, d.User, StatusCreated); err != nil {
		return err
	}
	d.ID = res[0]

	if len(d.BlastUsers) > 0 {
		for _, v := range d.BlastUsers {
			q := `
				INSERT INTO broadcast_users(
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
	}

	if len(d.ValidTenants) > 0 {
		for _, v := range d.ValidTenants {
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
			, voucher_type = ?
			, point_needed = ?
			, max_quantity_voucher = ?
			, max_usage_voucher = ?
			, allow_accumulative = ?
			, redeemtion_method = ?
			, start_date = ?
			, end_date = ?
			, discount_value = ?
			, img_url = ?
			, variant_tnc = ?
			, updated_by = ?
			, updated_at = ?
		WHERE
			id = ?
			AND status = ?;
	`

	_, err = tx.Exec(tx.Rebind(q), d.VariantName, d.VariantType, d.VoucherType, d.PointNeeded, d.MaxQuantityVoucher, d.MaxUsageVoucher, d.AllowAccumulative, d.RedeemtionMethod, d.StartDate, d.EndDate, d.DiscountValue, d.ImgUrl, d.VariantTnc, d.User, time.Now(), d.ID, StatusCreated)
	if err != nil {
		return err
	}

	if len(d.BlastUsers) > 0 {
		q = `
			UPDATE broadcast_users
			SET
				status = ?
				, updated_by = ?
				, updated_at = ?
			WHERE
				variant_id = ?
				AND status = ?;
		`
		_, err = tx.Exec(tx.Rebind(q), StatusDeleted, d.User, time.Now(), d.ID, StatusCreated)
		if err != nil {
			return err
		}

		for _, v := range d.BlastUsers {
			q := `
				INSERT INTO broadcast_users (
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
	}

	if len(d.ValidTenants) > 0 {
		q = `
			UPDATE variant_users
			SET
				status = ?
				, updated_by = ?
				, updated_at = ?
			WHERE
				variant_id = ?
				AND status = ?;
		`
		_, err = tx.Exec(tx.Rebind(q), StatusDeleted, d.User, time.Now(), d.ID, StatusCreated)
		if err != nil {
			return err
		}

		for _, v := range d.BlastUsers {
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
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func UpdateBroadcast(user UpdateVariantUsersRequest) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	q := `
		UPDATE broadcast_users
		SET
			status = ?
			, updated_by = ?
			, updated_at = ?
		WHERE
			variant_id = ?
			AND status = ?;
	`
	_, err = tx.Exec(tx.Rebind(q), StatusDeleted, user.User, time.Now(), user.ID, StatusCreated)
	if err != nil {
		return err
	}

	for _, v := range user.Data {
		q := `
			INSERT INTO broadcast_users (
				company_id
				, variant_id
				, account_id
				, created_by
			)
			VALUES (?, ?, ?, ?)
		`

		_, err := tx.Exec(tx.Rebind(q), user.CompanyID, user.ID, v, user.User)
		if err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func UpdateTenant(user UpdateVariantUsersRequest) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	q := `
		UPDATE variant_users
		SET
			status = ?
			, updated_by = ?
			, updated_at = ?
		WHERE
			variant_id = ?
			AND status = ?;
	`
	_, err = tx.Exec(tx.Rebind(q), StatusDeleted, user.User, time.Now(), user.ID, StatusCreated)
	if err != nil {
		return err
	}

	for _, v := range user.Data {
		q := `
			INSERT INTO variant_users (
				company_id
				, variant_id
				, account_id
				, created_by
			)
			VALUES (?, ?, ?, ?)
		`

		_, err := tx.Exec(tx.Rebind(q), user.CompanyID, user.ID, v, user.User)
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
			id = ?
			AND status = ?;
	`

	_, err = tx.Exec(tx.Rebind(q), d.User, time.Now(), StatusDeleted, d.ID, StatusCreated)
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
			variant_id = ?
			AND status = ?;
	`
	_, err = tx.Exec(tx.Rebind(q), d.User, time.Now(), StatusDeleted, d.ID, StatusCreated)
	if err != nil {
		return err
	}

	q = `
		UPDATE broadcast_users
		SET
			updated_by = ?
			, updated_at = ?
			, status = ?
		WHERE
			variant_id = ?
			AND status = ?;
	`
	_, err = tx.Exec(tx.Rebind(q), d.User, time.Now(), StatusDeleted, d.ID, StatusCreated)
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
			, voucher_type
			, point_needed
			, max_quantity_voucher
			, max_usage_voucher
			, allow_accumulative
			, redeemtion_method
			, start_date
			, end_date
			, discount_value
			, img_url
			, variant_tnc
			, created_by
		FROM
			variants
		WHERE
			id = ?
			AND status = ?
	`

	var resv []Variant
	if err := db.Select(&resv, db.Rebind(q), id, StatusCreated); err != nil {
		return VariantResponse{Status: "Error", Message: q, Data: nil}, err
	}
	if len(resv) < 1 {
		return VariantResponse{Status: "404", Message: q, Data: nil}, ErrResourceNotFound
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
		return VariantResponse{Status: "Error", Message: q, Data: nil}, err
	}
	resv[0].ValidTenants = resd

	res := VariantResponse{
		Status:  "200",
		Message: "Ok",
		Data:    resv[0],
	}

	return res, nil
}

func FindVariantByName(name string) (VariantResponse, error) {
	q := `
		SELECT
			id
			, variant_name
			, voucher_type
			, point_needed
			, max_quantity_voucher
			, start_date
			, end_date
		FROM
			variants
		WHERE
			variant_name LIKE ?
			AND status = ?
	`

	var resv []SearchVariant
	if err := db.Select(&resv, db.Rebind(q), name, StatusCreated); err != nil {
		return VariantResponse{Status: "500", Message: "Error at select variant", Data: nil}, err
	}
	if len(resv) < 1 {
		return VariantResponse{Status: "404", Message: "Error at select variant", Data: nil}, ErrResourceNotFound
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
			return VariantResponse{Status: "500", Message: "Error at select user", Data: nil}, err
		}
		resv[i].ValidUsers = resd
	}

	res := VariantResponse{
		Status:  "200",
		Message: "Ok",
		Data:    resv,
	}

	return res, nil
}

func FindVariantByUser(user string) (VariantResponse, error) {
	q := `
		SELECT
			id
			, variant_name
			, voucher_type
			, point_needed
			, max_quantity_voucher
			, start_date
			, end_date
		FROM
			variants
		WHERE
			created_by = ?
			AND status = ?
	`

	var resv []SearchVariant
	if err := db.Select(&resv, db.Rebind(q), user, StatusCreated); err != nil {
		return VariantResponse{Status: "500", Message: "Error at select variant", Data: nil}, err
	}
	if len(resv) < 1 {
		return VariantResponse{Status: "404", Message: "Error at select variant", Data: nil}, ErrResourceNotFound
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
			return VariantResponse{Status: "500", Message: "Error at select user", Data: nil}, err
		}
		resv[i].ValidUsers = resd
	}

	res := VariantResponse{
		Status:  "200",
		Message: "Ok",
		Data:    resv,
	}

	return res, nil
}

func FindVariantByCompanyID(id string) (VariantResponse, error) {
	q := `
		SELECT
			id
			, variant_name
			, voucher_type
			, point_needed
			, max_quantity_voucher
			, start_date
			, end_date
		FROM
			variants
		WHERE
			company_id = ?
			AND status = ?
	`

	var resv []SearchVariant
	if err := db.Select(&resv, db.Rebind(q), id, StatusCreated); err != nil {
		return VariantResponse{Status: "500", Message: "Error at select variant", Data: nil}, err
	}
	if len(resv) < 1 {
		return VariantResponse{Status: "404", Message: "Error at select variant", Data: nil}, ErrResourceNotFound
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
			return VariantResponse{Status: "500", Message: "Error at select user", Data: nil}, err
		}
		resv[i].ValidUsers = resd
	}

	res := VariantResponse{
		Status:  "200",
		Message: "Ok",
		Data:    resv,
	}

	return res, nil
}

func FindVariantByDate(start, end string) (VariantResponse, error) {
	q := `
		SELECT
			id
			, variant_name
			, voucher_type
			, point_needed
			, max_quantity_voucher
			, start_date
			, end_date
		FROM
			variants
		WHERE
			(start_date > ? AND start_date < ?)
			OR (end_date > ? AND end_date < ?)
			AND status = ?
	`

	var resv []SearchVariant
	if err := db.Select(&resv, db.Rebind(q), start, end, start, end, StatusCreated); err != nil {
		return VariantResponse{Status: "500", Message: "Error at select variant", Data: nil}, err
	}
	if len(resv) < 1 {
		return VariantResponse{Status: "404", Message: "Error at select variant", Data: nil}, ErrResourceNotFound
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
			return VariantResponse{Status: "500", Message: "Error at select user", Data: nil}, err
		}
		resv[i].ValidUsers = resd
	}

	res := VariantResponse{
		Status:  "200",
		Message: "Ok",
		Data:    resv,
	}

	return res, nil
}

func FindAllVariants() (VariantResponse, error) {
	q := `
		SELECT
			id
			, variant_name
			, voucher_type
			, point_needed
			, max_quantity_voucher
			, start_date
			, end_date
		FROM
			variants
	`

	var resv []SearchVariant
	if err := db.Select(&resv, db.Rebind(q)); err != nil {
		return VariantResponse{Status: "500", Message: "Error at select variant", Data: nil}, err
	}
	if len(resv) < 1 {
		return VariantResponse{Status: "404", Message: "Error at select variant", Data: nil}, ErrResourceNotFound
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
			return VariantResponse{Status: "500", Message: "Error at select user", Data: nil}, err
		}
		resv[i].ValidUsers = resd
	}

	res := VariantResponse{
		Status:  "200",
		Message: "Ok",
		Data:    resv,
	}

	return res, nil
}

func FindVariantMultipleParam(field, value []string) (VariantResponse, error) {
	q := `
		SELECT
			id
			, variant_name
			, voucher_type
			, point_needed
			, max_quantity_voucher
			, start_date
			, end_date
		FROM
			variants
		WHERE
			status = ?
	`
	for i := 0; i < len(field); i++ {
		if _, err := strconv.Atoi(value[i]); err == nil {
			q += ` AND ` + field[i] + ` = '` + value[i] + `'`
		} else {
			q += ` AND ` + field[i] + ` LIKE '%` + value[i] + `%'`

		}
	}

	var resv []SearchVariant

	if err := db.Select(&resv, db.Rebind(q), StatusCreated); err != nil {
		return VariantResponse{Status: "Error", Message: q, Data: nil}, err
	}
	if len(resv) < 1 {
		return VariantResponse{Status: "404", Message: q, Data: nil}, ErrResourceNotFound
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
			return VariantResponse{Status: "500", Message: "Error at select user", Data: nil}, err
		}
		resv[i].ValidUsers = resd
	}

	res := VariantResponse{
		Status:  "200",
		Message: "Ok",
		Data:    resv,
	}

	return res, nil
}
