package model

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var (
	ProgramType []string
)

type (
	Program struct {
		Id                 string  `db:"id" json:"id"`
		AccountId          string  `db:"account_id" json:"account_id"`
		Name               string  `db:"name" json:"name"`
		Type               string  `db:"type" json:"type"`
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
		VoucherValue       float64 `db:"voucher_value" json:"voucher_value"`
		MaxQuantityVoucher float64 `db:"max_quantity_voucher" json:"max_quantity_voucher"`
		MaxGenerateVoucher float64 `db:"max_generate_voucher" json:"max_generate_voucher"`
		MaxRedeemVoucher   float64 `db:"max_redeem_voucher" json:"max_redeem_voucher"`
		RedemptionMethod   string  `db:"redemption_method" json:"redeem_method"`
		ImgUrl             string  `db:"img_url" json:"image_url"`
		Tnc                string  `db:"tnc" json:"tnc"`
		Description        string  `db:"description" json:"description"`
		Visibility         bool    `db:"visibility" json:"visibility"`
		CreatedBy          string  `db:"created_by" json:"created_by"`
		CreatedAt          string  `db:"created_at" json:"created_at"`
	}
	ProgramReq struct {
		AccountId          string   `db:"account_id"`
		Name               string   `db:"name"`
		Type               string   `db:"type"`
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
		VoucherValue       float64  `db:"voucher_value"`
		MaxQuantityVoucher float64  `db:"max_quantity_voucher"`
		MaxGenerateVoucher float64  `db:"max_generate_voucher"`
		MaxRedeemVoucher   float64  `db:"max_redeem_voucher"`
		RedemptionMethod   string   `db:"redemption_method"`
		ImgUrl             string   `db:"img_url"`
		Tnc                string   `db:"tnc"`
		Description        string   `db:"description"`
		ValidPartners      []string `db:"valid_partners"`
	}
	FormatReq struct {
		Prefix     string `db:"prefix"`
		Postfix    string `db:"postfix"`
		Body       string `db:"body"`
		FormatType string `db:"format_type"`
		Length     int    `db:"length"`
	}
	DeleteProgramRequest struct {
		Id      string `db:"id"`
		User    string `db:"deleted_by"`
		Img_url string `db:"img_url"`
	}
	SearchProgram struct {
		Id                string                 `db:"id" json:"id"`
		AccountId         string                 `db:"account_id" json:"account_id"`
		Name              string                 `db:"name" json:"name"`
		Type              string                 `db:"type" json:"type"`
		VoucherType       string                 `db:"voucher_type" json:"voucher_type"`
		VoucherPrice      float64                `db:"voucher_price" json:"voucher_price"`
		VoucherValue      float64                `db:"voucher_value" json:"voucher_value"`
		AllowAccumulative bool                   `db:"allow_accumulative" json:"allow_accumulative"`
		MaxVoucher        float64                `db:"max_quantity_voucher" json:"max_quantity_voucher"`
		ImgUrl            string                 `db:"img_url" json:"image_url"`
		StartDate         string                 `db:"start_date" json:"start_date"`
		EndDate           string                 `db:"end_date" json:"end_date"`
		Vouchers          []SearchProgramVoucher `db:"-" json:"vouchers"`
		Voucher           string                 `db:"voucher" json:"voucher"`
		State             sql.NullString         `db:"state" json:"state"`
		Status            string                 `db:"status" json:"status"`
		CreatedAt         string                 `db:"created_at" json:"created_at"`
		UpdatedAt         sql.NullString         `db:"updated_at" json:"updated_at"`
	}
	SearchProgramVoucher struct {
		Voucher     string         `db:"voucher" json:"voucher"`
		State       string         `db:"state" json:"state"`
		PartnerName sql.NullString `db:"partner_name" json:"partner_name"`
		PartnerId   sql.NullString `db:"partner_id" json:"partner_id"`
	}
	UpdateProgramArrayRequest struct {
		ProgramId string   `db:"program_id"`
		User      string   `db:"updated_by"`
		Data      []string `db:"-"`
	}
	ProgramUpdateRequest struct {
		Id                 string `db:"id" json:"id"`
		AccountId          string `db:"account_id" json:"account_id"`
		Name               string `db:"name" json:"name"`
		Type               string `db:"type" json:"type"`
		VoucherFormat      string `db:"voucher_format_id" json:"voucher_format"`
		VoucherType        string `db:"voucher_type" json:"voucher_type"`
		VoucherPrice       string `db:"voucher_price" json:"voucher_price"`
		AllowAccumulative  string `db:"allow_accumulative" json:"allow_accumulative"`
		StartDate          string `db:"start_date" json:"start_date"`
		EndDate            string `db:"end_date" json:"end_date"`
		StartHour          string `db:"start_hour" json:"start_hour"`
		EndHour            string `db:"end_hour" json:"end_hour"`
		ValidVoucherStart  string `db:"valid_voucher_start" json:"valid_voucher_start"`
		ValidVoucherEnd    string `db:"valid_voucher_end" json:"valid_voucher_end"`
		VoucherLifetime    string `db:"voucher_lifetime" json:"voucher_lifetime"`
		ValidityDays       string `db:"validity_days" json:"validity_days"`
		VoucherValue       string `db:"voucher_value" json:"voucher_value"`
		MaxQuantityVoucher string `db:"max_quantity_voucher" json:"max_quantity_voucher"`
		MaxGenerateVoucher string `db:"max_generate_voucher" json:"max_generate_voucher"`
		MaxRedeemVoucher   string `db:"max_redeem_voucher" json:"max_redeem_voucher"`
		RedemptionMethod   string `db:"redemption_method" json:"redeem_method"`
		ImgUrl             string `db:"img_url" json:"image_url"`
		Tnc                string `db:"tnc" json:"tnc"`
		Description        string `db:"description" json:"description"`
		Visibility         string `db:"visibility" json:"visibility"`
		CreatedBy          string `db:"created_by" json:"created_by"`
		CreatedAt          string `db:"created_at" json:"created_at"`
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

func CheckValueProgram(column, id string) (string, error) {
	q := `
		SELECT ` + column + `
		FROM
			programs
		WHERE
			id = ?
	`
	var res []string
	if err := db.Select(&res, db.Rebind(q), id); err != nil {
		fmt.Println(err)
		return "", ErrServerInternal
	}
	if len(res) == 0 {
		return "", nil
	}
	return res[0], nil
}

func InsertProgram(vr ProgramReq, fr FormatReq, user string) (string, error) {
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
			, created_at
			, status
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING
			id
	`
	var res []string
	if err := tx.Select(&res, tx.Rebind(q), fr.Prefix, fr.Postfix, fr.Body, fr.FormatType, fr.Length, user, time.Now(), StatusCreated); err != nil {
		fmt.Println(err.Error(), "(insert voucher format)")
		return "", ErrServerInternal
	}

	q2 := `
		INSERT INTO programs(
			account_id
			, name
			, type
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
			, voucher_value
			, max_quantity_voucher
			, max_redeem_voucher
			, max_generate_voucher
			, redemption_method
			, img_url
			, tnc
			, description
			, created_by
			, created_at
			, status
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING
			id
	`
	var res2 []string
	if err := tx.Select(&res2, tx.Rebind(q2), vr.AccountId, vr.Name, vr.Type, res[0], vr.VoucherType, vr.VoucherPrice, vr.AllowAccumulative, vr.StartDate, vr.EndDate, vr.StartHour, vr.EndHour, vr.ValidVoucherStart, vr.ValidVoucherEnd, vr.VoucherLifetime, vr.ValidityDays, vr.VoucherValue, vr.MaxQuantityVoucher, vr.MaxRedeemVoucher, vr.MaxGenerateVoucher, vr.RedemptionMethod, vr.ImgUrl, vr.Tnc, vr.Description, user, time.Now(), StatusCreated); err != nil {
		fmt.Println(err.Error(), "(insert program)")
		return "", ErrServerInternal
	}

	if len(vr.ValidPartners) == 1 && vr.ValidPartners[0] == "all" {
		for _, v := range vr.ValidPartners {
			q := `
				INSERT INTO program_partners(
					program_id
					, partner_id
					, created_by
					, created_at
					, status
				)
				VALUES (?, ?, ?, ?, ?)
			`

			_, err := tx.Exec(tx.Rebind(q), res2[0], v, user, time.Now(), StatusCreated)
			if err != nil {
				fmt.Println("data :", res2[0], v, user)
				fmt.Println(err.Error(), "(insert program_partner)")
				return "", ErrServerInternal
			}
		}
	} else if len(vr.ValidPartners) > 0 {
		for _, v := range vr.ValidPartners {
			q := `
				INSERT INTO program_partners(
					program_id
					, partner_id
					, created_by
					, created_at
					, status
				)
				VALUES (?, ?, ?, ?, ?)
			`

			_, err := tx.Exec(tx.Rebind(q), res2[0], v, user, time.Now(), StatusCreated)
			if err != nil {
				fmt.Println("data :", res2[0], v, user)
				fmt.Println(err.Error(), "(insert program_partner)")
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

func UpdateProgram(d Program) error {
	tx, err := db.Beginx()
	if err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}
	defer tx.Rollback()

	programDetail, err := FindProgramDetailsByIdUpdateRequest(d.Id)
	if err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}

	reflectParam := reflect.ValueOf(&d)
	dataParam := reflect.Indirect(reflectParam)

	reflectDb := reflect.ValueOf(&programDetail).Elem()

	updates := getUpdate(dataParam, reflectDb)

	q := `
		UPDATE programs
		SET
			updated_by = ?
		WHERE
			id = ?
			AND status = ?;
		`
	_, err = tx.Exec(tx.Rebind(q), d.CreatedBy, d.Id, StatusCreated)
	if err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}

	q = `
		UPDATE programs
		SET
			updated_at = ?
		WHERE
			id = ?
			AND status = ?;
		`
	_, err = tx.Exec(tx.Rebind(q), time.Now(), d.Id, StatusCreated)
	if err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}

	for k, v := range updates {
		var value = v.String()
		if strings.Contains(value, "<") {
			tempString := strings.Replace(value, "<", "", -1)
			tempString = strings.Replace(tempString, ">", "", -1)
			tempStringArr := strings.Split(tempString, " ")

			if tempStringArr[0] == "int" {
				value = strconv.FormatInt(v.Int(), 10)
			} else if tempStringArr[0] == "float64" {
				value = strconv.FormatFloat(v.Float(), 'f', -1, 64)
			} else if tempStringArr[0] == "bool" {
				value = strconv.FormatBool(v.Bool())
			}
		}

		keys := strings.Split(k, ";")
		q = `
			UPDATE programs
			SET
				`
		q += keys[1] + ` = '` + value + `'`
		q += `
			WHERE
				id = ?
				AND status = ?;
		`
		_, err = tx.Exec(tx.Rebind(q), d.Id, StatusCreated)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println(q)
			return ErrServerInternal
		}
	}

	if err := tx.Commit(); err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}

	return nil
}

func UpdateBulkProgram(id, user string, voucher int) error {
	tx, err := db.Beginx()
	if err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}
	defer tx.Rollback()

	_, err = CheckValueProgram("max_quantity_voucher", id)

	q := `
		UPDATE programs
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

func UpdateProgramBroadcasts(user UpdateProgramArrayRequest) error {
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
			program_id = ?
			AND status = ?;
	`
	_, err = tx.Exec(tx.Rebind(q), StatusDeleted, user.User, time.Now(), user.ProgramId, StatusCreated)
	if err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}

	for _, v := range user.Data {
		q := `
			INSERT INTO broadcast_users (
				program_id
				, account_id
				, created_by
				, created_at
				, status
			)
			VALUES (?, ?, ?, ?, ?)
		`

		_, err := tx.Exec(tx.Rebind(q), user.ProgramId, v, user.User, time.Now(), StatusCreated)
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

func UpdateProgramPartners(param UpdateProgramArrayRequest) error {
	tx, err := db.Beginx()
	if err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}
	defer tx.Rollback()

	q := `
		UPDATE program_partners
		SET
			status = ?
			, updated_by = ?
			, updated_at = ?
		WHERE
			program_id = ?
			AND status = ?;
	`
	_, err = tx.Exec(tx.Rebind(q), StatusDeleted, param.User, time.Now(), param.ProgramId, StatusCreated)
	if err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}

	for _, v := range param.Data {
		q := `
			INSERT INTO program_partners (
				program_id
				, partner_id
				, created_by
				, created_at
				, status
			)
			VALUES (?, ?, ?, ?, ?)
		`

		_, err := tx.Exec(tx.Rebind(q), param.ProgramId, v, param.User, time.Now(), StatusCreated)
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

func (d *DeleteProgramRequest) Delete() error {
	tx, err := db.Beginx()
	if err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}
	defer tx.Rollback()

	voucher := CountVoucher(d.Id)
	if voucher > 0 {
		return ErrProgramNotNull
	}

	q := `
		UPDATE 	programs
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
		UPDATE program_partners
		SET
			updated_by = ?
			, updated_at = ?
			, status = ?
		WHERE
			program_id = ?
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
			programs as va
		WHERE
			va.id = ?
	`

	var resv []DeleteProgramRequest
	if err = db.Select(&resv, db.Rebind(q), d.Id); err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}
	if len(resv) < 1 {
		return ErrResourceNotFound
	}

	ctx := context.Background()
	if resv[0].Img_url != "https://storage.googleapis.com/e-voucher/L1LXN5bpMphnvG6Ce8eUbBSYDW5G3MaH.jpg" {
		o := StorageBucket.Object(strings.Split(resv[0].Img_url, "/")[4])
		if err := o.Delete(ctx); err != nil {
			return ErrServerInternal
		}
	}

	return nil
}

func VisibilityProgram(d DeleteProgramRequest, status bool) error {
	tx, err := db.Beginx()
	if err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}
	defer tx.Rollback()

	q := `
		UPDATE 	programs
		SET
			updated_by = ?
			, updated_at = ?
			, visibility = ?
		WHERE
			id = ?
			AND status = ?

	`
	_, err = tx.Exec(tx.Rebind(q), d.User, time.Now(), status, d.Id, StatusCreated)
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

func FindAllPrograms(accountId string) ([]SearchProgram, error) {
	q := `
		SELECT
			id
			, account_id
			, type
			, name
			, voucher_type
			, voucher_price
			, voucher_value
			, max_quantity_voucher
			, img_url
			, start_date
			, end_date
			, created_at
			, updated_at
			, status
		FROM
			programs
		WHERE
			account_id = ?
			AND status = ?
		ORDER BY
			end_date ASC
	`

	var resv []SearchProgram
	if err := db.Select(&resv, db.Rebind(q), accountId, StatusCreated); err != nil {
		fmt.Println(err.Error())
		return resv, ErrServerInternal
	}
	if len(resv) < 1 {
		return resv, ErrResourceNotFound
	}

	for index, value := range resv {
		q = `
		SELECT
			COUNT(id) as voucher, state
		FROM
			vouchers
		WHERE
			status = ?
			AND program_id = ?
		GROUP BY
			state
	`

		var resvo []SearchProgramVoucher
		if err := db.Select(&resvo, db.Rebind(q), StatusCreated, value.Id); err != nil {
			fmt.Println(err.Error())
			return []SearchProgram{}, ErrServerInternal
		}
		if len(resv) < 1 {
			return []SearchProgram{}, ErrResourceNotFound
		}

		resv[index].Vouchers = resvo
	}
	return resv, nil
}

func FindAllProgramsCustom(param map[string]string) ([]SearchProgram, error) {
	q := `
		SELECT
			id
			, account_id
			, type
			, name
			, voucher_type
			, voucher_price
			, voucher_value
			, max_quantity_voucher
			, img_url
			, start_date
			, end_date
			, created_at
			, updated_at
			, status
		FROM
			programs
		WHERE
			status = ?
	`
	for key, value := range param {
		if _, err := strconv.Atoi(value); err == nil {
			q += ` AND ` + key + ` = '` + value + `'`
		} else if strings.Contains(key, "date") {
			q += ` AND ` + key + ` ` + value
		} else {
			q += ` AND ` + key + ` LIKE '%` + value + `%'`

		}
	}
	q += `
		ORDER BY
			end_date ASC
	`
	//fmt.Println(q)
	var resv []SearchProgram
	if err := db.Select(&resv, db.Rebind(q), StatusCreated); err != nil {
		fmt.Println(err.Error())
		return resv, ErrServerInternal
	}
	if len(resv) < 1 {
		return resv, ErrResourceNotFound
	}

	for index, value := range resv {
		q = `SELECT
			COUNT(v.id) as voucher, v.state, p.name as partner_name, p.id as partner_id
		FROM
			vouchers as v
		LEFT JOIN
			transaction_details as td
		ON
			td.voucher_id = v.id
		LEFT JOIN
			transactions as t
		ON
			td.transaction_id = t.id
		LEFT JOIN
			partners as p
		ON
			t.partner_id = p.id
		WHERE
			v.status = ?
			AND v.program_id = ?
		GROUP BY
			v.state, p.name, p.id
		ORDER BY
			voucher DESC
		`

		var resvo []SearchProgramVoucher
		if err := db.Select(&resvo, db.Rebind(q), StatusCreated, value.Id); err != nil {
			fmt.Println(err.Error())
			return []SearchProgram{}, ErrServerInternal
		}
		if len(resv) < 1 {
			return []SearchProgram{}, ErrResourceNotFound
		}

		resv[index].Vouchers = resvo
	}
	return resv, nil
}

func FindAvailablePrograms(param map[string]string) ([]SearchProgram, error) {
	q := `
		SELECT
			va.id
			, va.account_id
			, va.name
			, va.type
			, va.voucher_type
			, va.voucher_price
			, va.voucher_value
			, va.allow_accumulative
			, va.max_quantity_voucher
			, va.img_url
			, va.start_date
			, va.end_date
			, count (vo.id) as voucher
		FROM
			programs as va
		LEFT JOIN
			vouchers as vo
		ON
			va.id = vo.program_id
		WHERE
			va.status = ?
			AND visibility = ?
	`
	for key, value := range param {
		if strings.Contains(key, "date") {
			q += ` AND va.` + key + ` ` + value
		} else {
			q += ` AND va.` + key + ` = '` + value + `'`
		}
	}

	q += `
		GROUP BY
			va.id
		ORDER BY
			va.created_at DESC
	`

	var resv []SearchProgram
	if err := db.Select(&resv, db.Rebind(q), StatusCreated, true); err != nil {
		fmt.Println(err.Error())
		return []SearchProgram{}, err
	}
	if len(resv) < 1 {
		return []SearchProgram{}, ErrResourceNotFound
	}

	return resv, nil
}

func FindProgramsCustomParam(param map[string]string) ([]SearchProgram, error) {
	q := `
		SELECT
			va.id
			, va.account_id
			, va.name
			, va.voucher_type
			, va.voucher_price
			, va.voucher_value
			, va.max_quantity_voucher
			, va.img_url
			, va.start_date
			, va.end_date
			, count (vo.id) as voucher
		FROM
			programs as va
		LEFT JOIN
			vouchers as vo
		ON
			va.id = vo.program_id
		WHERE
			va.status = ?
	`
	for key, value := range param {
		if key == "q" {
			q += `AND (va.name ILIKE '%` + value + `%' OR va.account_id ILIKE '%` + value + `%' OR va.voucher_type ILIKE '%` + value + `%')`
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

	var resv []SearchProgram
	if err := db.Select(&resv, db.Rebind(q), StatusCreated); err != nil {
		fmt.Println(err.Error())
		return []SearchProgram{}, err
	}
	if len(resv) < 1 {
		return []SearchProgram{}, ErrResourceNotFound
	}

	return resv, nil
}

func FindProgramDetailsById(id string) (Program, error) {
	q := `
		SELECT
			id
			, account_id
			, name
			, type
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
			, voucher_value
			, max_quantity_voucher
			, max_redeem_voucher
			, max_generate_voucher
			, redemption_method
			, img_url
			, tnc
			, description
			, created_by
			, created_at
		FROM
			programs
		WHERE
			status = ?
			AND id = ?
	`

	var resv []Program
	if err := db.Select(&resv, db.Rebind(q), StatusCreated, id); err != nil {
		fmt.Println(err.Error())
		return Program{}, ErrServerInternal
	}
	//fmt.Println("program data :", id, StatusCreated, resv)
	if len(resv) < 1 {
		return Program{}, ErrResourceNotFound
	}

	return resv[0], nil
}

func FindProgramDetailsByIdUpdateRequest(id string) (ProgramUpdateRequest, error) {
	q := `
		SELECT
			id
			, account_id
			, name
			, type
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
			, voucher_value
			, max_quantity_voucher
			, max_redeem_voucher
			, max_generate_voucher
			, redemption_method
			, img_url
			, tnc
			, description
			, created_by
			, created_at
		FROM
			programs
		WHERE
			status = ?
			AND id = ?
	`

	var resv []ProgramUpdateRequest
	if err := db.Select(&resv, db.Rebind(q), StatusCreated, id); err != nil {
		fmt.Println(err.Error())
		return ProgramUpdateRequest{}, ErrServerInternal
	}
	//fmt.Println("program data :", id, StatusCreated, resv)
	if len(resv) < 1 {
		return ProgramUpdateRequest{}, ErrResourceNotFound
	}

	return resv[0], nil
}

func FindProgramDetailsCustomParam(param map[string]string) ([]Program, error) {
	q := `
		SELECT
			id
			, account_id
			, name
			, type
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
			, voucher_value
			, max_quantity_voucher
			, max_redeem_voucher
			, max_generate_voucher
			, redemption_method
			, img_url
			, tnc
			, description
			, visibility
			, created_by
			, created_at
		FROM
			programs
		WHERE
			status = ?
	`
	for key, value := range param {
		if key == "q" {
			q += `AND (name ILIKE '%` + value + `%' OR account_id ILIKE '%` + value + `%' OR voucher_type ILIKE '%` + value + `%')`
		} else {
			if _, err := strconv.Atoi(value); err == nil {
				q += ` AND ` + key + ` = '` + value + `'`
			} else {
				q += ` AND ` + key + ` LIKE '%` + value + `%'`

			}
		}
	}

	var resv []Program
	if err := db.Select(&resv, db.Rebind(q), StatusCreated); err != nil {
		return []Program{}, ErrServerInternal
	}
	if len(resv) < 1 {
		return []Program{}, ErrResourceNotFound
	}

	return resv, nil
}

func FindProgramsPartner(partnerId, accountId string) ([]SearchProgram, error) {
	q := `
		SELECT
			DISTINCT
			p.id
			, p.account_id
			, p.name
			, p.type
			, p.voucher_type
			, p.voucher_price
			, p.voucher_value
			, p.max_quantity_voucher
			, p.img_url
			, p.start_date
			, p.end_date
		FROM
			programs as p
		JOIN
			program_partners as pp
		ON
			p.id = pp.program_id
		WHERE
			p.account_id = ?
			AND pp.partner_id = ?
			AND p.status = ?
	`

	var resv []SearchProgram
	if err := db.Select(&resv, db.Rebind(q), accountId, partnerId, StatusCreated); err != nil {
		return []SearchProgram{}, err
	}
	if len(resv) < 1 {
		return []SearchProgram{}, ErrResourceNotFound
	}

	for i, v := range resv {
		q = `
			SELECT
				count (DISTINCT td.voucher_id) as voucher
			FROM
				transactions as t
			JOIN
				transaction_details as td
			ON
				t.id = td.transaction_id
			JOIN
				vouchers as v
			ON
				v.id = td.voucher_id
			JOIN
				programs as pr
			ON
				pr.id = v.program_id
			WHERE
				t.partner_id = ?
				AND pr.id = ?
				AND pr.account_id = ?
				AND pr.status = ?
		`

		var vo []string
		if err := db.Select(&vo, db.Rebind(q), partnerId, v.Id, accountId, StatusCreated); err != nil {
			return []SearchProgram{}, err
		}

		resv[i].Voucher = vo[0]
	}

	return resv, nil
}

func FindTodayProgramsPartner(parterId, accountId string) ([]SearchProgram, error) {
	q := `
		SELECT
			va.id
			, va.account_id
			, va.name
			, va.voucher_type
			, va.voucher_price
			, va.voucher_value
			, va.max_quantity_voucher
			, va.img_url
			, va.start_date
			, va.end_date
			, count (vo.id) as voucher
		FROM
			programs as va
		LEFT JOIN
			vouchers as vo
		ON
			va.id = vo.program_id
		JOIN
			program_partners as vp
		ON
			va.id = vp.program_id
		WHERE
			vp.partner_id = ?
			AND va.account_id = ?
			AND va.status = ?
			AND va.start_date <= ?
			AND va.end_date >= ?
		GROUP BY
			va.id
		ORDER BY
			va.start_date DESC
	`

	var resv []SearchProgram
	if err := db.Select(&resv, db.Rebind(q), parterId, accountId, StatusCreated, time.Now(), time.Now()); err != nil {
		return []SearchProgram{}, err
	}
	if len(resv) < 1 {
		return []SearchProgram{}, ErrResourceNotFound
	}

	return resv, nil
}

func GetProgramTypes() error {
	q := `
		SELECT
			value
		FROM
			account_configs
		WHERE
			config_id = 4
			AND account_id = 'unknown'
			AND status = ?
	`

	var resv []string
	if err := db.Select(&resv, db.Rebind(q), StatusCreated); err != nil {
		return ErrServerInternal
	}
	if len(resv) < 1 {
		return ErrResourceNotFound
	}
	ProgramType = resv
	return nil
}
