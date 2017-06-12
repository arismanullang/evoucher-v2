package model

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"
)

type (
	Variant struct {
		Id                 string  `db:"id" json:"id"`
		AccountId          string  `db:"account_id" json:"account_id"`
		VariantName        string  `db:"variant_name" json:"variant_name"`
		VariantType        string  `db:"variant_type" json:"variant_type"`
		VoucherFormat      int     `db:"voucher_format_id" json:"voucher_format"`
		VoucherType        string  `db:"voucher_type" json:"voucher_type"`
		VoucherPrice       float64 `db:"voucher_price" json:"voucher_price"`
		AllowAccumulative  bool    `db:"allow_accumulative" json:"allow_accumulative"`
		StartDate          string  `db:"start_date" json:"start_date"`
		EndDate            string  `db:"end_date" json:"end_date"`
		StartHour          string  `db:"start_hour" json:"start_hour"`
		EndHour            string  `db:"end_hour" json:"end_hour"`
		ValidVoucherStart  string  `db:"valid_voucher_start" json:"valid_voucher_start"`
		ValidVoucherEnd    string  `db:"valid_voucher_end" json:"valid_voucher_end"`
		VoucherLifetime    int     `db:"voucher_lifetime" json:"voucher_lifetime"`
		ValidityDays       string  `db:"validity_days" json:"validity_days"`
		DiscountValue      float64 `db:"discount_value" json:"discount_value"`
		MaxQuantityVoucher float64 `db:"max_quantity_voucher" json:"max_quantity_voucher"`
		MaxUsageVoucher    float64 `db:"max_usage_voucher" json:"max_usage_voucher"`
		RedeemtionMethod   string  `db:"redeemtion_method" json:"redeem_method"`
		ImgUrl             string  `db:"img_url" json:"image_url"`
		VariantTnc         string  `db:"variant_tnc" json:"variant_tnc"`
		VariantDescription string  `db:"variant_description" json:"variant_description"`
		CreatedBy          string  `db:"created_by" json:"created_by"`
		CreatedAt          string  `db:"created_at" json:"created_at"`
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
		StartHour          string   `db:"start_hour"`
		EndHour            string   `db:"end_hour"`
		ValidVoucherStart  string   `db:"valid_voucher_start"`
		ValidVoucherEnd    string   `db:"valid_voucher_end"`
		VoucherLifetime    int      `db:"voucher_lifetime"`
		ValidityDays       string   `db:"validity_days"`
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
		Id      string `db:"id"`
		User    string `db:"deleted_by"`
		Img_url string `db:"img_url"`
	}
	SearchVariant struct {
		Id            string         `db:"id" json:"id"`
		AccountId     string         `db:"account_id" json:"account_id"`
		VariantName   string         `db:"variant_name" json:"variant_name"`
		VoucherType   string         `db:"voucher_type" json:"voucher_type"`
		VoucherPrice  float64        `db:"voucher_price" json:"voucher_price"`
		DiscountValue float64        `db:"discount_value" json:"discount_value"`
		MaxVoucher    float64        `db:"max_quantity_voucher" json:"max_quantity_voucher"`
		ImgUrl        string         `db:"img_url" json:"image_url"`
		StartDate     string         `db:"start_date" json:"start_date"`
		EndDate       string         `db:"end_date" json:"end_date"`
		Voucher       string         `db:"voucher" json:"voucher"`
		State         sql.NullString `db:"state" json:"state"`
		Status        string         `db:"status" json:"status"`
		CreatedAt     string         `db:"created_at" json:"created_at"`
		UpdatedAt     sql.NullString `db:"updated_at" json:"updated_at"`
	}
	UpdateVariantArrayRequest struct {
		VariantId string   `db:"variant_id"`
		User      string   `db:"updated_by"`
		Data      []string `db:"-"`
	}
)

func CustomQuery(q string) (map[int][]map[string]interface{}, error) {
	fmt.Println("Select Database")

	rows, err := db.Query(q)
	if err != nil {
		fmt.Println(err.Error())
		return map[int][]map[string]interface{}{}, ErrServerInternal
	}
	cols, err := rows.Columns()
	if err != nil {
		fmt.Println(err.Error())
		return map[int][]map[string]interface{}{}, ErrServerInternal
	}
	defer rows.Close()

	index := 0
	result := make(map[int][]map[string]interface{})

	for rows.Next() {
		m := make(map[string]interface{})

		pointer := make([]interface{}, len(cols))
		pointers := make([]interface{}, len(cols))

		for i, _ := range cols {
			pointers[i] = &pointer[i]
		}

		err := rows.Scan(pointers...)
		if err != nil {
			log.Fatal(err)
		}

		for i, col := range cols {

			var v interface{}

			val := pointer[i]

			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}

			fmt.Println(col, v)

			v, ok = v.(string)
			if ok {
				m[col] = v.(string)
			} else {
				m[col] = v.(float64)
			}
		}

		result[index] = append(result[index], m)
		index++
	}
	fmt.Println(result)
	return result, nil
}

func InsertVariant(vr VariantReq, fr FormatReq, user string) (string, error) {
	tx, err := db.Beginx()
	if err != nil {
		fmt.Println(err.Error())
		return "", ErrServerInternal
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
		fmt.Println(err.Error(), "(insert voucher format)")
		return "", ErrServerInternal
	}

	fmt.Println(vr)
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
			, start_hour
			, end_hour
			, valid_voucher_start
			, valid_voucher_end
			, voucher_lifetime
			, validity_days
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
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING
			id
	`
	var res2 []string
	if err := tx.Select(&res2, tx.Rebind(q2), vr.AccountId, vr.VariantName, vr.VariantType, res[0], vr.VoucherType, vr.VoucherPrice, vr.AllowAccumulative, vr.StartDate, vr.EndDate, vr.StartHour, vr.EndHour, vr.ValidVoucherStart, vr.ValidVoucherEnd, vr.VoucherLifetime, vr.ValidityDays, vr.DiscountValue, vr.MaxQuantityVoucher, vr.MaxUsageVoucher, vr.RedeemtionMethod, vr.ImgUrl, vr.VariantTnc, vr.VariantDescription, user, StatusCreated); err != nil {
		fmt.Println(err.Error(), "(insert variant)")
		return "", ErrServerInternal
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
				fmt.Println("data :", res2[0], v, user)
				fmt.Println(err.Error(), "(insert variant_partner)")
				return "", ErrServerInternal
			}
		}
	}

	if err := tx.Commit(); err != nil {

		fmt.Println(err.Error())
		return "", ErrServerInternal
	}
	return res2[0], nil
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
			, start_hour = ?
			, end_hour = ?
			, valid_voucher_start = ?
			, valid_voucher_end = ?
			, voucher_lifetime = ?
			, validity_days = ?
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

	_, err = tx.Exec(tx.Rebind(q), d.VariantName, d.VariantType, d.VoucherType, d.VoucherPrice, d.StartDate, d.EndDate, d.StartHour, d.EndHour, d.ValidVoucherStart, d.ValidVoucherEnd, d.VoucherLifetime, d.ValidityDays, d.DiscountValue, d.MaxQuantityVoucher, d.MaxUsageVoucher, d.RedeemtionMethod, d.ImgUrl, d.VariantTnc, d.VariantDescription, d.CreatedBy, time.Now(), d.Id, StatusCreated)
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

func UpdateBulkVariant(id string, voucher int) error {
	tx, err := db.Beginx()
	if err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}
	defer tx.Rollback()

	q := `
		UPDATE variants
		SET
			max_quantity_voucher = ?
		WHERE
			id = ?
			AND status = ?
	`

	_, err = tx.Exec(tx.Rebind(q), voucher, id, StatusCreated)
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

func UpdateVariantBroadcasts(user UpdateVariantArrayRequest) error {
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

func UpdateVariantPartners(param UpdateVariantArrayRequest) error {
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
		UPDATE 	variants
		SET
			deleted_by = ?
			, deleted_at = ?
			, status = ?
		WHERE
			id = ?
			AND status = ?

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

	q = `
		SELECT
			id
			, deleted_by
			, img_url
		FROM
			variants as va
		WHERE
			va.id = ?
	`

	var resv []DeleteVariantRequest
	if err = db.Select(&resv, db.Rebind(q), d.Id); err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}
	if len(resv) < 1 {
		return ErrResourceNotFound
	}
	*d = resv[0]

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
			, va.created_at
			, va.updated_at
			, count (vo.id) as voucher
			, va.status
			, vo.state
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
			va.id, vo.state
		ORDER BY
			va.end_date ASC
	`

	var resv []SearchVariant
	if err := db.Select(&resv, db.Rebind(q), accountId, StatusCreated); err != nil {
		fmt.Println(err.Error())
		return resv, ErrServerInternal
	}
	if len(resv) < 1 {
		return resv, ErrResourceNotFound
	}

	return resv, nil
}

func FindVariantsCustomParam(param map[string]string) ([]SearchVariant, error) {
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
		return []SearchVariant{}, err
	}
	if len(resv) < 1 {
		return []SearchVariant{}, ErrResourceNotFound
	}

	return resv, nil
}

func FindVariantDetailsById(id string) (Variant, error) {
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
			, start_hour
			, end_hour
			, valid_voucher_start
			, valid_voucher_end
			, voucher_lifetime
			, validity_days
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
			status = ?
			AND id = ?
	`

	var resv []Variant
	if err := db.Select(&resv, db.Rebind(q), StatusCreated, id); err != nil {
		fmt.Println(err.Error())
		return Variant{}, ErrServerInternal
	}
	fmt.Println("variant data :", id, StatusCreated, resv)
	if len(resv) < 1 {
		return Variant{}, ErrResourceNotFound
	}

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
			, start_hour
			, end_hour
			, valid_voucher_start
			, valid_voucher_end
			, voucher_lifetime
			, validity_days
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

	return resv, nil
}
