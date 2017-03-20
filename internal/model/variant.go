package model

import (
	"fmt"
	"strconv"
	"time"
)

type (
	Variant struct {
		Id                 string   `db:"id"`
		AccountId          string   `db:"account_id"`
		VariantName        string   `db:"variant_name"`
		VariantType        string   `db:"variant_type"`
		VoucherFormat      int      `db:"voucher_format_id"`
		VoucherType        string   `db:"voucher_type"`
		VoucherPrice       float64  `db:"voucher_price"`
		AllowAccumulative  bool     `db:"allow_accumulative"`
		StartDate          string   `db:"start_date"`
		EndDate            string   `db:"end_date"`
		DiscountValue      float64  `db:"discount_value"`
		MaxQuantityVoucher float64  `db:"max_quantity_voucher"`
		MaxUsageVoucher    float64  `db:"max_usage_voucher"`
		RedeemtionMethod   string   `db:"redeemtion_method"`
		ImgUrl             string   `db:"img_url"`
		VariantTnc         string   `db:"variant_tnc"`
		VariantDescription string   `db:"variant_description"`
		CreatedBy          string   `db:"created_by"`
		CreatedAt          string   `db:"created_at"`
		ValidPartners      []string `db:"-"`
		Voucher            []string `db:"-"`
	}
	VariantReq struct {
		AccountId          string   `db:"account_id"`
		VariantName        string   `db:"variant_name"`
		VariantType        string   `db:"variant_type"`
		VoucherType        string   `db:"voucher_type"`
		VoucherPrice       float64  `db:"voucher_price"`
		AllowAccumulative  bool     `db:"allow_accumulative"`
		StartDate          string   `db:"start_date"`
		EndDate            string   `db:"end_date"`
		DiscountValue      float64  `db:"discount_value"`
		MaxQuantityVoucher float64  `db:"max_quantity_voucher"`
		MaxUsageVoucher    float64  `db:"max_usage_voucher"`
		RedeemtionMethod   string   `db:"redeemtion_method"`
		ImgUrl             string   `db:"img_url"`
		VariantTnc         string   `db:"variant_tnc"`
		VariantDescription string   `db:"variant_description"`
		ValidPartners      []string `db:"valid_partners"`
	}
	FormatReq struct {
		Prefix     string `db:"prefix"`
		Postfix    string `db:"postfix"`
		Body       string `db:"body"`
		FormatType string `db:"format_type"`
		Length     int    `db:"length"`
	}
	DeleteVariantRequest struct {
		Id   string `db:"id"`
		User string `db:"deleted_by"`
	}
	SearchVariant struct {
		Id            string  `db:"id"`
		AccountId     string  `db:"account_id"`
		VariantName   string  `db:"variant_name"`
		VoucherType   string  `db:"voucher_type"`
		VoucherPrice  float64 `db:"voucher_price"`
		DiscountValue float64 `db:"discount_value"`
		MaxVoucher    float64 `db:"max_quantity_voucher"`
		ImgUrl        string  `db:"img_url"`
		StartDate     string  `db:"start_date"`
		EndDate       string  `db:"end_date"`
		Voucher       string  `db:"voucher"`
	}
	UpdateVariantUsersRequest struct {
		VariantId string   `db:"variant_id"`
		User      string   `db:"updated_by"`
		Data      []string `db:"-"`
	}
)

func InsertVariant(vr VariantReq, fr FormatReq, user string) error {
	tx, err := db.Beginx()
	if err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}
	defer tx.Rollback()

	q := `
		INSERT INTO voucher_formats(
			prefix
			, postfix
			, body
			, format_type
			, length
			, created_by
			, status
		)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		RETURNING
			id
	`
	var res []string
	if err := tx.Select(&res, tx.Rebind(q), fr.Prefix, fr.Postfix, fr.Body, fr.FormatType, fr.Length, user, StatusCreated); err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}

	q2 := `
		INSERT INTO variants(
			account_id
			, variant_name
			, variant_type
			, voucher_format_id
			, voucher_type
			, voucher_price
			, allow_accumulative
			, start_date
			, end_date
			, discount_value
			, max_quantity_voucher
			, max_usage_voucher
			, redeemtion_method
			, img_url
			, variant_tnc
			, variant_description
			, created_by
			, status
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING
			id
	`
	var res2 []string
	if err := tx.Select(&res2, tx.Rebind(q2), vr.AccountId, vr.VariantName, vr.VariantType, res[0], vr.VoucherType, vr.VoucherPrice, vr.AllowAccumulative, vr.StartDate, vr.EndDate, vr.DiscountValue, vr.MaxQuantityVoucher, vr.MaxUsageVoucher, vr.RedeemtionMethod, vr.ImgUrl, vr.VariantTnc, vr.VariantDescription, user, StatusCreated); err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}

	if len(vr.ValidPartners) > 0 {
		for _, v := range vr.ValidPartners {
			q := `
				INSERT INTO variant_partners(
					variant_id
					, partner_id
					, created_by
					, status
				)
				VALUES (?, ?, ?, ?)
			`

			_, err := tx.Exec(tx.Rebind(q), res2[0], v, user, StatusCreated)
			if err != nil {
				fmt.Println(err.Error())
				return ErrServerInternal
			}
		}
	}

	if err := tx.Commit(); err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}
	return nil
}

func UpdateVariant(d Variant) error {
	tx, err := db.Beginx()
	if err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}
	defer tx.Rollback()

	q := `
		UPDATE variants
		SET
			variant_name = ?
			, variant_type = ?
			, voucher_type = ?
			, voucher_price = ?
			, start_date = ?
			, end_date = ?
			, discount_value = ?
			, max_quantity_voucher = ?
			, max_usage_voucher = ?
			, redeemtion_method = ?
			, img_url = ?
			, variant_tnc = ?
			, variant_description = ?
			, updated_by = ?
			, updated_at = ?
		WHERE
			id = ?
			AND status = ?
	`

	_, err = tx.Exec(tx.Rebind(q), d.VariantName, d.VariantType, d.VoucherType, d.VoucherPrice, d.StartDate, d.EndDate, d.DiscountValue, d.MaxQuantityVoucher, d.MaxUsageVoucher, d.RedeemtionMethod, d.ImgUrl, d.VariantTnc, d.VariantDescription, d.CreatedBy, time.Now(), d.Id, StatusCreated)
	if err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}

	if err := tx.Commit(); err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}
	return nil
}

func UpdateVariantBroadcasts(user UpdateVariantUsersRequest) error {
	tx, err := db.Beginx()
	if err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
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
	_, err = tx.Exec(tx.Rebind(q), StatusDeleted, user.User, time.Now(), user.VariantId, StatusCreated)
	if err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}

	for _, v := range user.Data {
		q := `
			INSERT INTO broadcast_users (
				variant_id
				, account_id
				, created_by
				, status
			)
			VALUES (?, ?, ?, ?)
		`

		_, err := tx.Exec(tx.Rebind(q), user.VariantId, v, user.User, StatusCreated)
		if err != nil {
			fmt.Println(err.Error())
			return ErrServerInternal
		}
	}

	if err := tx.Commit(); err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}
	return nil
}

func UpdateVariantPartners(param UpdateVariantUsersRequest) error {
	tx, err := db.Beginx()
	if err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}
	defer tx.Rollback()

	q := `
		UPDATE variant_partners
		SET
			status = ?
			, updated_by = ?
			, updated_at = ?
		WHERE
			variant_id = ?
			AND status = ?;
	`
	_, err = tx.Exec(tx.Rebind(q), StatusDeleted, param.User, time.Now(), param.VariantId, StatusCreated)
	if err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}

	for _, v := range param.Data {
		q := `
			INSERT INTO variant_partners (
				variant_id
				, partner_id
				, created_by
				, status
			)
			VALUES (?, ?, ?, ?)
		`

		_, err := tx.Exec(tx.Rebind(q), param.VariantId, v, param.User, StatusCreated)
		if err != nil {
			fmt.Println(err.Error())
			return ErrServerInternal
		}
	}

	if err := tx.Commit(); err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}
	return nil
}

func (d *DeleteVariantRequest) Delete() error {
	tx, err := db.Beginx()
	if err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
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

	_, err = tx.Exec(tx.Rebind(q), d.User, time.Now(), StatusDeleted, d.Id, StatusCreated)
	if err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}

	q = `
		UPDATE variant_partners
		SET
			updated_by = ?
			, updated_at = ?
			, status = ?
		WHERE
			variant_id = ?
			AND status = ?;
	`
	_, err = tx.Exec(tx.Rebind(q), d.User, time.Now(), StatusDeleted, d.Id, StatusCreated)
	if err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}

	if err := tx.Commit(); err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}
	return nil
}

func FindVariantsByDate(start, end, accountId string) ([]SearchVariant, error) {
	fmt.Println("Select By Date " + start)
	q := `
		SELECT
			va.id
			, va.account_id
			, va.variant_name
			, va.voucher_type
			, va.voucher_price
			, va.discount_value
			, va.max_quantity_voucher
			, va.img_url
			, va.start_date
			, va.end_date
			, count (vo.id) as voucher
		FROM
			variants as va
		LEFT JOIN
			vouchers as vo
		ON
			va.id = vo.variant_id
		WHERE
			(start_date > ? AND start_date < ?)
			OR (end_date > ? AND end_date < ?)
			AND account_id = ?
			AND status = ?
		GROUP BY
			va.id
		ORDER BY
			va.start_date DESC
	`

	var resv []SearchVariant
	if err := db.Select(&resv, db.Rebind(q), start, end, start, end, accountId, StatusCreated); err != nil {
		return []SearchVariant{}, ErrServerInternal
	}
	if len(resv) < 1 {
		return []SearchVariant{}, ErrResourceNotFound
	}

	return resv, nil
}

func FindAllVariants(accountId string) ([]SearchVariant, error) {
	q := `
		SELECT
			va.id
			, va.account_id
			, va.variant_name
			, va.voucher_type
			, va.voucher_price
			, va.discount_value
			, va.max_quantity_voucher
			, va.img_url
			, va.start_date
			, va.end_date
			, count (vo.id) as voucher
		FROM
			variants as va
		LEFT JOIN
			vouchers as vo
		ON
			va.id = vo.variant_id
		WHERE
			va.account_id = ?
			AND va.status = ?
		GROUP BY
			va.id
		ORDER BY
			va.start_date DESC
	`

	var resv []SearchVariant
	if err := db.Select(&resv, db.Rebind(q), accountId, StatusCreated); err != nil {
		return resv, ErrServerInternal
	}
	if len(resv) < 1 {
		return resv, ErrResourceNotFound
	}

	return resv, nil
}

func FindVariantsCustomParam(param map[string]string) ([]SearchVariant, error) {
	fmt.Println("Query start")
	q := `
		SELECT
			va.id
			, va.account_id
			, va.variant_name
			, va.voucher_type
			, va.voucher_price
			, va.discount_value
			, va.max_quantity_voucher
			, va.img_url
			, va.start_date
			, va.end_date
			, count (vo.id) as voucher
		FROM
			variants as va
		LEFT JOIN
			vouchers as vo
		ON
			va.id = vo.variant_id
		WHERE
			va.status = ?
	`
	for key, value := range param {
		if key == "q" {
			q += `AND (va.variant_name ILIKE '%` + value + `%' OR va.account_id ILIKE '%` + value + `%' OR va.voucher_type ILIKE '%` + value + `%')`
		} else {
			if _, err := strconv.Atoi(value); err == nil {
				q += ` AND va.` + key + ` = '` + value + `'`
			} else {
				q += ` AND va.` + key + ` LIKE '%` + value + `%'`

			}
		}
	}
	q += `
		GROUP BY
			va.id
	`

	var resv []SearchVariant
	if err := db.Select(&resv, db.Rebind(q), StatusCreated); err != nil {
		fmt.Println(err.Error())
		return []SearchVariant{}, ErrServerInternal
	}
	if len(resv) < 1 {
		return []SearchVariant{}, ErrResourceNotFound
	}

	return resv, nil
}

func FindVariantDetailsById(id string) (Variant, error) {
	q := `
		SELECT
			va.id
			, va.account_id
			, va.variant_name
			, va.variant_type
			, va.voucher_format_id
			, va.voucher_type
			, va.voucher_price
			, va.allow_accumulative
			, va.start_date
			, va.end_date
			, va.discount_value
			, va.max_quantity_voucher
			, va.max_usage_voucher
			, va.redeemtion_method
			, va.img_url
			, va.variant_tnc
			, va.variant_description
			, va.created_by
			, va.created_at
		FROM
			variants as va
		LEFT JOIN
			vouchers as vo
		ON
			va.id = vo.variant_id
		WHERE
			va.id = ?
			AND va.status = ?
		GROUP BY
			va.id
	`

	var resv []Variant
	if err := db.Select(&resv, db.Rebind(q), id, StatusCreated); err != nil {
		fmt.Println(err.Error())
		return Variant{}, ErrServerInternal
	}
	if len(resv) < 1 {
		return Variant{}, ErrResourceNotFound
	}

	q = `
		SELECT
			partner_id
		FROM
			variant_partners
		WHERE
			variant_id = ?
			AND status = ?
	`
	var resd []string
	if err := db.Select(&resd, db.Rebind(q), id, StatusCreated); err != nil {
		fmt.Println(err.Error())
		return Variant{}, ErrServerInternal
	}
	resv[0].ValidPartners = resd

	q = `
		SELECT
			voucher_code
		FROM
			vouchers
		WHERE
			variant_id = ?
			AND status = ?
	`
	var reso []string
	if err := db.Select(&reso, db.Rebind(q), id, StatusCreated); err != nil {
		fmt.Println(err.Error())
		return Variant{}, ErrServerInternal
	}
	resv[0].Voucher = reso

	return resv[0], nil
}

func FindVariantDetailsCustomParam(param map[string]string) ([]Variant, error) {
	q := `
		SELECT
			id
			, account_id
			, variant_name
			, variant_type
			, voucher_format_id
			, voucher_type
			, voucher_price
			, allow_accumulative
			, start_date
			, end_date
			, discount_value
			, max_quantity_voucher
			, max_usage_voucher
			, redeemtion_method
			, img_url
			, variant_tnc
			, variant_description
			, created_by
			, created_at
		FROM
			variants
		WHERE
			AND status = ?
	`
	for key, value := range param {
		if key == "q" {
			q += `AND (variant_name ILIKE '%` + value + `%' OR account_id ILIKE '%` + value + `%' OR voucher_type ILIKE '%` + value + `%')`
		} else {
			if _, err := strconv.Atoi(value); err == nil {
				q += ` AND ` + key + ` = '` + value + `'`
			} else {
				q += ` AND ` + key + ` LIKE '%` + value + `%'`

			}
		}
	}

	var resv []Variant
	if err := db.Select(&resv, db.Rebind(q), StatusCreated); err != nil {
		return []Variant{}, ErrServerInternal
	}
	if len(resv) < 1 {
		return []Variant{}, ErrResourceNotFound
	}

	for i, v := range resv {
		q = `
			SELECT
				partner_id
			FROM
				variant_partners
			WHERE
				variant_id = ?
				AND status = ?
		`
		var resd []string
		if err := db.Select(&resd, db.Rebind(q), v.Id, StatusCreated); err != nil {
			return []Variant{}, err
		}
		resv[i].ValidPartners = resd
	}

	return resv, nil
}
