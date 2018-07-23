package model

import (
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type (
	Partner struct {
		Id                string         `db:"id" json:"id"`
		AccountId         string         `db:"account_id" json:"acccount_id"`
		Name              string         `db:"name" json:"name"`
		Email             string         `db:"email" json:"email"`
		SerialNumber      sql.NullString `db:"serial_number" json:"serial_number"`
		CreatedBy         sql.NullString `db:"created_by" json:"created_by"`
		ProgramID         string         `db:"program_id" json:"program_id"`
		Tag               sql.NullString `db:"tag" json:"tag"`
		Description       sql.NullString `db:"description" json:"description"`
		BankAccount       BankAccount    `db:"-" json:"bank_account"`
		Address           string         `db:"address" json:"address"`
		City              string         `db:"city" json:"city"`
		Province          string         `db:"province" json:"province"`
		Building          string         `db:"building" json:"building"`
		ZipCode           string         `db:"zip_code" json:"zip_code"`
		CompanyName       string         `db:"company_name" json:"company_name"`
		CompanyPic        string         `db:"company_pic" json:"company_pic"`
		CompanyTelp       string         `db:"company_telp" json:"company_telp"`
		CompanyEmail      string         `db:"company_email" json:"company_email"`
		BankName          string         `db:"bank_name" json:"bank_name"`
		BankBranch        string         `db:"bank_branch" json:"bank_branch"`
		BankAccountNumber string         `db:"bank_account_number" json:"bank_account_number"`
		BankAccountHolder string         `db:"bank_account_holder" json:"bank_account_holder"`
	}
	PartnerUpdateRequest struct {
		Id                string `db:"id" json:"id"`
		AccountId         string `db:"account_id" json:"acccount_id"`
		Name              string `db:"name" json:"name"`
		Email             string `db:"email" json:"email"`
		SerialNumber      string `db:"serial_number" json:"serial_number"`
		CreatedBy         string `db:"created_by" json:"created_by"`
		ProgramID         string `db:"program_id" json:"program_id"`
		Tag               string `db:"tag" json:"tag"`
		Description       string `db:"description" json:"description"`
		Address           string `db:"address" json:"address"`
		City              string `db:"city" json:"city"`
		Province          string `db:"province" json:"province"`
		Building          string `db:"building" json:"building"`
		ZipCode           string `db:"zip_code" json:"zip_code"`
		CompanyName       string `db:"company_name" json:"company_name"`
		CompanyPic        string `db:"company_pic" json:"company_pic"`
		CompanyTelp       string `db:"company_telp" json:"company_telp"`
		CompanyEmail      string `db:"company_email" json:"company_email"`
		BankName          string `db:"bank_name" json:"bank_name"`
		BankBranch        string `db:"bank_branch" json:"bank_branch"`
		BankAccountNumber string `db:"bank_account_number" json:"bank_account_number"`
		BankAccountHolder string `db:"bank_account_holder" json:"bank_account_holder"`
	}
	PartnerProgramSummary struct {
		Id                string         `db:"id" json:"id"`
		AccountId         string         `db:"account_id" json:"acccount_id"`
		Name              string         `db:"name" json:"name"`
		SerialNumber      sql.NullString `db:"serial_number" json:"serial_number"`
		CreatedBy         sql.NullString `db:"created_by" json:"created_by"`
		ProgramID         string         `db:"program_id" json:"program_id"`
		Tag               sql.NullString `db:"tag" json:"tag"`
		Description       sql.NullString `db:"description" json:"description"`
		Transactions      int            `db:"-" json:"transactions"`
		Vouchers          int            `db:"-" json:"vouchers"`
		TransactionValues float32        `db:"-" json:"transaction_values"`
	}
	Tag struct {
		ID        string `db:"id" json:"id"`
		Value     string `db:"value" json:"value"`
		AccountID string `db:"account_id" json:"account_id"`
	}
)

func InsertPartner(p Partner) error {
	fmt.Println("Add")
	tx, err := db.Beginx()
	if err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}
	defer tx.Rollback()

	partner, err := checkPartner(p.Name, p.AccountId)
	if err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}

	if partner != "" {
		return ErrDuplicateEntry
	}

	tags := strings.Split(p.Tag.String, "#")
	err = CheckAndInsertTag(tags, p.CreatedBy.String, p.AccountId)
	if err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}

	q := `
		INSERT INTO partners(
			name
			, account_id
			, serial_number
			, email
			, tag
			, description
			, building
			, address
			, city
			, province
			, zip_code
			, company_name
			, company_pic
			, company_telp
			, company_email
			, bank_name
			, bank_branch
			, bank_account_number
			, bank_account_holder
			, created_by
			, status
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ? ,? ,? ,? ,? ,? ,?)
		RETURNING
			id
	`

	var res []string
	err = tx.Select(&res, tx.Rebind(q), p.Name, p.AccountId, p.SerialNumber, p.Email, p.Tag.String, p.Description, p.Building, p.Address, p.City, p.Province, p.ZipCode, p.CompanyName, p.CompanyPic, p.CompanyTelp, p.CompanyEmail, p.BankName, p.BankBranch, p.BankAccountNumber, p.BankAccountHolder, p.CreatedBy.String, StatusCreated)
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

func checkPartner(name, accountId string) (string, error) {
	q := `
		SELECT id
		FROM partners
		WHERE
			name = ?
			AND account_id = ?
			AND status = ?
	`
	var res []string
	if err := db.Select(&res, db.Rebind(q), name, accountId, StatusCreated); err != nil {
		return "", ErrServerInternal
	}
	if len(res) == 0 {
		return "", nil
	}
	return res[0], nil
}

func FindPartners(param map[string]string) ([]Partner, error) {
	q := `
		SELECT
			id
			, account_id
			, name
			, serial_number
			, email
			, tag
			, description
			, building
			, address
			, city
			, province
			, zip_code
			, bank_name
			, bank_branch
			, bank_account_number
			, bank_account_holder
			, company_name
			, company_pic
			, company_telp
			, company_email
		FROM partners
		WHERE status = ?
	`
	for key, value := range param {
		if key == "q" {
			q += `AND (id ILIKE '%` + value + `%' OR name ILIKE '%` + value + `%' OR serial_number ILIKE '%` + value + `%')`
		} else {
			q += ` AND ` + key + ` ILIKE '%` + value + `%'`
		}
	}
	q += ` ORDER BY name`

	var resv []Partner
	if err := db.Select(&resv, db.Rebind(q), StatusCreated); err != nil {
		fmt.Println(err.Error())
		return []Partner{}, ErrServerInternal
	}
	if len(resv) < 1 {
		return []Partner{}, ErrResourceNotFound
	}

	return resv, nil
}

func FindPartnerById(param string) (Partner, error) {
	q := `
		SELECT
			id
			, account_id
			, name
			, serial_number
			, email
			, tag
			, description
			, building
			, address
			, city
			, province
			, zip_code
		FROM partners
		WHERE status = ?
		AND id = ?
	`

	var resv []Partner
	if err := db.Select(&resv, db.Rebind(q), StatusCreated, param); err != nil {
		fmt.Println(err.Error())
		return Partner{}, ErrServerInternal
	}
	if len(resv) < 1 {
		return Partner{}, ErrResourceNotFound
	}

	return resv[0], nil
}

func FindPartnerByIdUpdateRequest(param string) (PartnerUpdateRequest, error) {
	q := `
		SELECT
			id
			, account_id
			, name
			, serial_number
			, email
			, tag
			, description
			, building
			, address
			, city
			, province
			, zip_code
		FROM
			partners
		WHERE
			status = ?
		 	AND id = ?
	`

	var resv []Partner
	if err := db.Select(&resv, db.Rebind(q), StatusCreated, param); err != nil {
		fmt.Println(err.Error())
		return PartnerUpdateRequest{}, ErrServerInternal
	}
	if len(resv) < 1 {
		return PartnerUpdateRequest{}, ErrResourceNotFound
	}

	result := PartnerUpdateRequest{
		Id:           resv[0].Id,
		AccountId:    resv[0].AccountId,
		Name:         resv[0].Name,
		SerialNumber: resv[0].SerialNumber.String,
		Email:        resv[0].Email,
		Tag:          resv[0].Tag.String,
		Description:  resv[0].Description.String,
		Building:     resv[0].Building,
		Address:      resv[0].Address,
		City:         resv[0].City,
		Province:     resv[0].Province,
		ZipCode:      resv[0].ZipCode,
	}

	return result, nil
}

func FindAllPartners(accountId string) ([]Partner, error) {
	q := `
		SELECT
			id
			, account_id
			, name
			, serial_number
			, tag
			, description
			, bank_name
			, company_name
			, bank_account_number
		FROM partners
		WHERE status = ?
		AND account_id = ?
		ORDER BY name
	`

	var resv []Partner
	if err := db.Select(&resv, db.Rebind(q), StatusCreated, accountId); err != nil {
		fmt.Println(err.Error())
		return []Partner{}, ErrServerInternal
	}
	if len(resv) < 1 {
		return []Partner{}, ErrResourceNotFound
	}

	return resv, nil
}

func UpdatePartner(partner PartnerUpdateRequest, user, account string) error {
	tx, err := db.Beginx()
	if err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}
	defer tx.Rollback()

	tags := strings.Split(partner.Tag, "#")
	err = CheckAndInsertTag(tags, user, account)
	if err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}

	q := `
		UPDATE partners
		SET
			updated_by = ?
			, updated_at = ?
		WHERE
			id = ?
			AND status = ?;
	`
	_, err = tx.Exec(tx.Rebind(q), user, time.Now(), partner.Id, StatusCreated)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	partnerDetail, err := FindPartnerByIdUpdateRequest(partner.Id)
	if err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}

	reflectParam := reflect.ValueOf(&partner)
	dataParam := reflect.Indirect(reflectParam)

	reflectDb := reflect.ValueOf(&partnerDetail).Elem()

	updates := getUpdate(dataParam, reflectDb)

	for k, v := range updates {
		var value = v.String()
		if strings.Contains(value, "<") {
			tempString := strings.Replace(value, "<", "", -1)
			tempString = strings.Replace(tempString, ">", "", -1)
			tempStringArr := strings.Split(tempString, " ")
			if tempStringArr[0] == "int" {
				value = strconv.FormatInt(v.Int(), 64)
			} else if tempStringArr[0] == "float64" {
				value = strconv.FormatFloat(v.Float(), 'f', -1, 64)
			}
		}

		keys := strings.Split(k, ";")
		q = `
			UPDATE partners
			SET
				`
		q += keys[1] + ` = '` + value + `'`
		q += `
			WHERE
				id = ?
				AND status = ?;
		`
		_, err = tx.Exec(tx.Rebind(q), partner.Id, StatusCreated)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println(q)
			return ErrServerInternal
		}
	}

	if err := tx.Commit(); err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}

func DeletePartner(partner, user string) error {
	tx, err := db.Beginx()
	if err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}
	defer tx.Rollback()

	q := `
		UPDATE partners
		SET
			updated_by = ?
			, updated_at = ?
			, status = ?
		WHERE
			id = ?;
	`
	_, err = tx.Exec(tx.Rebind(q), user, time.Now(), StatusDeleted, partner)
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

func FindProgramPartner(param map[string]string) ([]Partner, error) {
	q := `
		SELECT 	b.id
			, b.account_id
			, b.name
			, b.serial_number
			, b.created_by
			, a.program_id
	 	FROM
			program_partners a
		JOIN
		 	partners b
		ON
			a.partner_id = b.id
 		WHERE
			a.status = ?
	`
	for k, v := range param {
		table := "b"
		if k == "program_id" {
			table = "a"
		}

		q += ` AND ` + table + `.` + k + ` = '` + v + `'`

	}
	q += ` ORDER BY b.name`

	var resv []Partner
	if err := db.Select(&resv, db.Rebind(q), StatusCreated); err != nil {
		return []Partner{}, err
	}
	if len(resv) < 1 {
		return []Partner{}, ErrResourceNotFound
	}
	return resv, nil
}

func FindProgramPartners(programId string) ([]Partner, error) {
	q := `
		SELECT 	b.id
			, b.account_id
			, b.name
			, b.serial_number
			, b.created_by
			, a.program_id
	 	FROM
			program_partners a
		JOIN
		 	partners b
		ON
			a.partner_id = b.id
 		WHERE
			a.status = ?
			AND a.program_id = ?
		`
	var resv []Partner
	if err := db.Select(&resv, db.Rebind(q), StatusCreated, programId); err != nil {
		return []Partner{}, err
	}
	if len(resv) < 1 {
		return []Partner{}, ErrResourceNotFound
	}
	return resv, nil
}

func FindProgramPartnerSummary(accountId, programId string) ([]PartnerProgramSummary, error) {
	q := `
		SELECT 	b.id
			, b.account_id
			, b.name
			, b.serial_number
			, b.created_by
			, a.program_id
	 	FROM
			program_partners a
		JOIN
		 	partners b
		ON
			a.partner_id = b.id
 		WHERE
			a.status = ?
			AND a.program_id = ?
		`
	var resv []PartnerProgramSummary
	if err := db.Select(&resv, db.Rebind(q), StatusCreated, programId); err != nil {
		return []PartnerProgramSummary{}, err
	}
	if len(resv) < 1 {
		return []PartnerProgramSummary{}, ErrResourceNotFound
	}

	for i, v := range resv {
		param := make(map[string]string)
		param["t.account_id"] = accountId
		param["va.id"] = programId
		param["p.id"] = v.Id
		transactions, err := FindTransactions(param)
		if err != nil {
			if err != ErrResourceNotFound {
				return []PartnerProgramSummary{}, err
			}
		}

		resv[i].Transactions = len(transactions)

		vouchers := 0
		var voucherValues float32
		for _, vv := range transactions {
			vouchers += len(vv.Voucher)
			voucherValues += vv.VoucherValue
		}

		resv[i].Vouchers = vouchers
		resv[i].TransactionValues = voucherValues
	}
	return resv, nil
}

// ------------------------------------------------------------------------------
// Tag

func FindAllTags(account string) ([]string, error) {
	q := `
		SELECT
			value
		FROM tags
		WHERE status = ?
			AND account_id = ?
	`

	var resv []string
	if err := db.Select(&resv, db.Rebind(q), StatusCreated, account); err != nil {
		fmt.Println(err.Error())
		return []string{}, ErrServerInternal
	}
	if len(resv) < 1 {
		return []string{}, ErrResourceNotFound
	}

	return resv, nil
}

func CheckAndInsertTag(tags []string, user, account string) error {
	tx, err := db.Beginx()
	if err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}
	defer tx.Rollback()

	for _, v := range tags {
		q := `
			SELECT
				id
			FROM
				tags
			WHERE
				status = ?
				AND value = ?
		`

		var resv []string
		if err := db.Select(&resv, db.Rebind(q), StatusCreated, v); err != nil {
			fmt.Println(err.Error())
			return ErrServerInternal
		}

		if len(resv) < 1 && v != "" {
			q := `
				INSERT INTO tags(
					value
					, account_id
					, created_by
					, status
				)
				VALUES (?, ?, ?, ?)
			`

			_, err := tx.Exec(tx.Rebind(q), v, account, user, StatusCreated)
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

func InsertTag(tag, user, account string) error {
	fmt.Println("Add")
	tx, err := db.Beginx()
	if err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}
	defer tx.Rollback()

	q := `
		INSERT INTO tags(
			value
			, account_id
			, created_by
			, status
		)
		VALUES (?, ?, ?, ?)
	`

	_, err = tx.Exec(tx.Rebind(q), tag, account, user, StatusCreated)
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

func DeleteTag(value, user string) error {
	tx, err := db.Beginx()
	if err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}
	defer tx.Rollback()

	q := `
		UPDATE tags
		SET
			updated_by = ?
			, updated_at = ?
			, status = ?
		WHERE
			value = ?;
	`
	_, err = tx.Exec(tx.Rebind(q), user, time.Now(), StatusDeleted, value)
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

func DeleteTagBulk(tagValue []string, user string) error {
	tx, err := db.Beginx()
	if err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}
	defer tx.Rollback()

	q := `
		UPDATE tags
		SET
			updated_by = ?
			, updated_at = ?
			, status = ?
		WHERE
			value = ?
	`
	for i := 1; i < len(tagValue); i++ {
		q += " OR value = '" + tagValue[i] + "'"
	}

	_, err = tx.Exec(tx.Rebind(q), user, time.Now(), StatusDeleted, tagValue[0])
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
