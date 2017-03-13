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
		ValidPartners      []string `db:"valid_partners"`
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
		Id string `db:"id"`
		//AccountId     string   `db:"account_id"`
		VariantName string `db:"variant_name"`
		//VoucherType   string   `db:"voucher_type"`
		VoucherPrice  float64 `db:"voucher_price"`
		DiscountValue float64 `db:"discount_value"`
		MaxVoucher    float64 `db:"max_quantity_voucher"`
		// StartDate     string   `db:"start_date"`
		// EndDate       string   `db:"end_date"`
		// ValidUsers    []string `db:"-"`
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
		return err
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
		return err
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
		return err
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
				return err
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func UpdateVariant(d Variant) error {
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
		return err
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
	_, err = tx.Exec(tx.Rebind(q), StatusDeleted, user.User, time.Now(), user.VariantId, StatusCreated)
	if err != nil {
		return err
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
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func UpdatePartner(param UpdateVariantUsersRequest) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
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
		return err
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

	_, err = tx.Exec(tx.Rebind(q), d.User, time.Now(), StatusDeleted, d.Id, StatusCreated)
	if err != nil {
		return err
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
	_, err = tx.Exec(tx.Rebind(q), d.User, time.Now(), StatusDeleted, d.Id, StatusCreated)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func FindVariantByDate(start, end string) (Response, error) {
	fmt.Println("Select By Date " + start)
	q := `
		SELECT
			id
			, variant_name
			, account_id
			, voucher_type
			, voucher_price
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
		return Response{Status: "500", Message: "Error at select variant", Data: nil}, err
	}
	if len(resv) < 1 {
		return Response{Status: "404", Message: "Error at select variant", Data: nil}, ErrResourceNotFound
	}

	// for i, v := range resv {
	// 	q = `
	// 		SELECT
	// 			partner_id
	// 		FROM
	// 			variant_partners
	// 		WHERE
	// 			variant_id = ?
	// 			AND status = ?
	// 	`
	// 	var resd []string
	// 	if err := db.Select(&resd, db.Rebind(q), v.Id, StatusCreated); err != nil {
	// 		return Response{Status: "500", Message: "Error at select user", Data: nil}, err
	// 	}
	// 	resv[i].ValidUsers = resd
	// }

	res := Response{
		Status:  "200",
		Message: "Ok",
		Data:    resv,
	}

	return res, nil
}

func FindAllVariants(accountId string) (Response, error) {
	q := `
		SELECT
			id
			, variant_name
			, voucher_price
			, discount_value
			, max_quantity_voucher
		FROM
			variants
		WHERE
			account_id = ?
			AND status = ?
	`

	var resv []SearchVariant
	if err := db.Select(&resv, db.Rebind(q), accountId, StatusCreated); err != nil {
		return Response{Status: "500", Message: "Error at select variant", Data: nil}, err
	}
	if len(resv) < 1 {
		return Response{Status: "404", Message: ErrMessageResourceNotFound, Data: nil}, ErrResourceNotFound
	}

	// for i, v := range resv {
	// 	q = `
	// 		SELECT
	// 			partner_id
	// 		FROM
	// 			variant_partners
	// 		WHERE
	// 			variant_id = ?
	// 			AND status = ?
	// 	`
	// 	var resd []string
	// 	if err := db.Select(&resd, db.Rebind(q), v.Id, StatusCreated); err != nil {
	// 		return Response{Status: "500", Message: "Error at select user", Data: nil}, err
	// 	}
	// 	resv[i].ValidUsers = resd
	// }

	res := Response{
		Status:  "200",
		Message: "Ok",
		Data:    resv,
	}

	return res, nil
}

func FindVariantMultipleParam(param map[string]string) (Response, error) {
	fmt.Println("Query start")
	q := `
		SELECT
			id
			, variant_name
			, account_id
			, voucher_type
			, voucher_price
			, max_quantity_voucher
			, start_date
			, end_date
		FROM
			variants
		WHERE
			status = ?
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

	var resv []SearchVariant
	if err := db.Select(&resv, db.Rebind(q), StatusCreated); err != nil {
		return Response{Status: "Error", Message: q, Data: nil}, err
	}
	if len(resv) < 1 {
		return Response{Status: "404", Message: q, Data: nil}, ErrResourceNotFound
	}

	// for i, v := range resv {
	// 	q = `
	// 		SELECT
	// 			partner_id
	// 		FROM
	// 			variant_partners
	// 		WHERE
	// 			variant_id = ?
	// 			AND status = ?
	// 	`
	// 	var resd []string
	// 	if err := db.Select(&resd, db.Rebind(q), v.Id, StatusCreated); err != nil {
	// 		return Response{Status: "500", Message: "Error at select user", Data: nil}, err
	// 	}
	// 	resv[i].ValidUsers = resd
	// }

	res := Response{
		Status:  "200",
		Message: "Ok",
		Data:    resv,
	}

	return res, nil
}

func FindVariantById(id string) (Response, error) {
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
			id = ?
			AND status = ?
	`

	var resv []Variant
	if err := db.Select(&resv, db.Rebind(q), id, StatusCreated); err != nil {
		return Response{Status: "Error", Message: q, Data: Variant{}}, err
	}
	if len(resv) < 1 {
		return Response{Status: "404", Message: q, Data: Variant{}}, ErrResourceNotFound
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
		return Response{Status: "Error", Message: q, Data: Variant{}}, err
	}
	resv[0].ValidPartners = resd

	res := Response{
		Status:  "200",
		Message: "Ok",
		Data:    resv[0],
	}

	return res, nil
}
