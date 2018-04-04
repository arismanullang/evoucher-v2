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
		Id           string         `db:"id" json:"id"`
		AccountId    string         `db:"account_id" json:"acccount_id"`
		Name         string         `db:"name" json:"name"`
		Email        string         `db:"email" json:"email"`
		SerialNumber sql.NullString `db:"serial_number" json:"serial_number"`
		CreatedBy    sql.NullString `db:"created_by" json:"created_by"`
		ProgramID    string         `db:"program_id" json:"program_id"`
		Tag          sql.NullString `db:"tag" json:"tag"`
		Description  sql.NullString `db:"description" json:"description"`
		BankAccount  BankAccount    `db:"-" json:"bank_account"`
		Address      string         `db:"address" json:"address"`
		City         string         `db:"city" json:"city"`
		Province     string         `db:"province" json:"province"`
		Building     string         `db:"building" json:"building"`
		ZipCode      string         `db:"zip_code" json:"zip_code"`
	}
	PartnerUpdateRequest struct {
		Id           string `db:"id" json:"id"`
		AccountId    string `db:"account_id" json:"acccount_id"`
		Name         string `db:"name" json:"name"`
		Email        string `db:"email" json:"email"`
		SerialNumber string `db:"serial_number" json:"serial_number"`
		CreatedBy    string `db:"created_by" json:"created_by"`
		ProgramID    string `db:"program_id" json:"program_id"`
		Tag          string `db:"tag" json:"tag"`
		Description  string `db:"description" json:"description"`
		BankAccount  string `db:"-" json:"bank_account_id"`
		Address      string `db:"address" json:"address"`
		City         string `db:"city" json:"city"`
		Province     string `db:"province" json:"province"`
		Building     string `db:"building" json:"building"`
		ZipCode      string `db:"zip_code" json:"zip_code"`
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
		Value string `db:"value" json:"value"`
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
	err = CheckAndInsertTag(tags, p.CreatedBy.String)
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
			, bank_account_id
			, building
			, address
			, city
			, province
			, zip_code
			, created_by
			, created_at
			, status
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING
			id
	`

	var res []string
	err = tx.Select(&res, tx.Rebind(q), p.Name, p.AccountId, p.SerialNumber, p.Email, p.Tag.String, p.Description, p.BankAccount.Id, p.Building, p.Address, p.City, p.Province, p.ZipCode, p.CreatedBy.String, timeNow(), StatusCreated)
	if err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}

	if err := tx.Commit(); err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}

	log := Log{
		TableName:   "partners",
		TableNameId: ValueChangeLogNone,
		ColumnName:  ColumnChangeLogInsert,
		Action:      ActionChangeLogInsert,
		Old:         ValueChangeLogNone,
		New:         res[0],
		CreatedBy:   p.CreatedBy.String,
	}

	err = addLog(log)
	if err != nil {
		fmt.Println(err.Error())
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

	var resv []Partner
	if err := db.Select(&resv, db.Rebind(q), StatusCreated); err != nil {
		fmt.Println(err.Error())
		return []Partner{}, ErrServerInternal
	}
	if len(resv) < 1 {
		return []Partner{}, ErrResourceNotFound
	}

	for i := 0; i < len(resv); i++ {
		tempBank, err := FindBankAccountByPartner(resv[i].AccountId, resv[i].Id)
		if err != nil {
			fmt.Println(err.Error())
			return []Partner{}, ErrServerInternal
		}
		if len(resv) < 1 {
			return []Partner{}, ErrResourceNotFound
		}

		resv[i].BankAccount = tempBank
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

	bank, _ := FindBankAccountByPartner(resv[0].AccountId, resv[0].Id)
	result := PartnerUpdateRequest{
		Id:           resv[0].Id,
		AccountId:    resv[0].AccountId,
		BankAccount:  bank.Id,
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
		FROM partners
		WHERE status = ?
		AND account_id = ?
	`

	var resv []Partner
	if err := db.Select(&resv, db.Rebind(q), StatusCreated, accountId); err != nil {
		fmt.Println(err.Error())
		return []Partner{}, ErrServerInternal
	}
	if len(resv) < 1 {
		return []Partner{}, ErrResourceNotFound
	}

	for i := 0; i < len(resv); i++ {
		tempBank, err := FindBankAccountByPartner(resv[i].AccountId, resv[i].Id)
		if err != nil {
			fmt.Println(err.Error())
			return []Partner{}, ErrServerInternal
		}

		resv[i].BankAccount = tempBank
	}
	return resv, nil
}

func UpdatePartner(partner PartnerUpdateRequest, user string) error {
	tx, err := db.Beginx()
	if err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}
	defer tx.Rollback()

	tags := strings.Split(partner.Tag, "#")
	err = CheckAndInsertTag(tags, user)
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
	dataDb := reflect.Indirect(reflectDb)

	updates := getUpdate(dataParam, reflectDb)

	logs := []Log{}
	tempLog := Log{
		TableName:   "partners",
		TableNameId: partner.Id,
		ColumnName:  "updated_by",
		Action:      ActionChangeLogUpdate,
		Old:         ValueChangeLogNone,
		New:         user,
		CreatedBy:   user,
	}
	logs = append(logs, tempLog)

	tempLog = Log{
		TableName:   "partners",
		TableNameId: partner.Id,
		ColumnName:  "updated_at",
		Action:      ActionChangeLogUpdate,
		Old:         ValueChangeLogNone,
		New:         time.Now().String(),
		CreatedBy:   user,
	}
	logs = append(logs, tempLog)

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

		tempLog = Log{
			TableName:   "partners",
			TableNameId: partner.Id,
			ColumnName:  keys[1],
			Action:      ActionChangeLogUpdate,
			Old:         dataDb.FieldByName(keys[0]).Interface().(string),
			New:         value,
			CreatedBy:   user,
		}
		logs = append(logs, tempLog)
	}

	if err := tx.Commit(); err != nil {
		fmt.Println(err.Error())
		return err
	}

	err = addLogs(logs)
	if err != nil {
		fmt.Println(err.Error())
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

	logs := []Log{}
	tempLog := Log{
		TableName:   "partners",
		TableNameId: partner,
		ColumnName:  "updated_by",
		Action:      ActionChangeLogUpdate,
		Old:         ValueChangeLogNone,
		New:         user,
		CreatedBy:   user,
	}
	logs = append(logs, tempLog)

	tempLog = Log{
		TableName:   "partners",
		TableNameId: partner,
		ColumnName:  "updated_at",
		Action:      ActionChangeLogUpdate,
		Old:         ValueChangeLogNone,
		New:         time.Now().String(),
		CreatedBy:   user,
	}
	logs = append(logs, tempLog)

	tempLog = Log{
		TableName:   "partners",
		TableNameId: partner,
		ColumnName:  "status",
		Action:      ActionChangeLogUpdate,
		Old:         ValueChangeLogNone,
		New:         StatusDeleted,
		CreatedBy:   user,
	}
	logs = append(logs, tempLog)

	err = addLogs(logs)
	if err != nil {
		fmt.Println(err.Error())
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

func FindAllTags() ([]string, error) {
	q := `
		SELECT
			value
		FROM tags
		WHERE status = ?
	`

	var resv []string
	if err := db.Select(&resv, db.Rebind(q), StatusCreated); err != nil {
		fmt.Println(err.Error())
		return []string{}, ErrServerInternal
	}
	if len(resv) < 1 {
		return []string{}, ErrResourceNotFound
	}

	return resv, nil
}

func CheckAndInsertTag(tags []string, user string) error {
	tx, err := db.Beginx()
	if err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}
	defer tx.Rollback()

	logs := []Log{}

	for _, v := range tags {
		q := `
		SELECT
			id
		FROM tags
		WHERE status = ?
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
					, created_by
					, status
				)
				VALUES (?, ?, ?)
			`

			_, err := tx.Exec(tx.Rebind(q), v, user, StatusCreated)
			if err != nil {
				fmt.Println(err.Error())
				return ErrServerInternal
			}

			tempLog := Log{
				TableName:   "tags",
				TableNameId: ValueChangeLogNone,
				ColumnName:  ColumnChangeLogInsert,
				Action:      ActionChangeLogInsert,
				Old:         ValueChangeLogNone,
				New:         v,
				CreatedBy:   user,
			}
			logs = append(logs, tempLog)
		}
	}

	if err := tx.Commit(); err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}

	err = addLogs(logs)
	if err != nil {
		fmt.Println(err.Error())
	}

	return nil
}

func InsertTag(tag, user string) error {
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
			, created_by
			, status
		)
		VALUES (?, ?, ?)
	`

	_, err = tx.Exec(tx.Rebind(q), tag, user, StatusCreated)
	if err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}

	if err := tx.Commit(); err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}

	log := Log{
		TableName:   "tags",
		TableNameId: ValueChangeLogNone,
		ColumnName:  ColumnChangeLogInsert,
		Action:      ActionChangeLogInsert,
		Old:         ValueChangeLogNone,
		New:         tag,
		CreatedBy:   user,
	}

	err = addLog(log)
	if err != nil {
		fmt.Println(err.Error())
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

	logs := []Log{}
	tempLog := Log{
		TableName:   "tags",
		TableNameId: value,
		ColumnName:  "updated_by",
		Action:      ActionChangeLogUpdate,
		Old:         ValueChangeLogNone,
		New:         user,
		CreatedBy:   user,
	}
	logs = append(logs, tempLog)

	tempLog = Log{
		TableName:   "tags",
		TableNameId: value,
		ColumnName:  "updated_at",
		Action:      ActionChangeLogUpdate,
		Old:         ValueChangeLogNone,
		New:         time.Now().String(),
		CreatedBy:   user,
	}
	logs = append(logs, tempLog)

	tempLog = Log{
		TableName:   "tags",
		TableNameId: value,
		ColumnName:  "status",
		Action:      ActionChangeLogUpdate,
		Old:         ValueChangeLogNone,
		New:         StatusDeleted,
		CreatedBy:   user,
	}
	logs = append(logs, tempLog)

	err = addLogs(logs)
	if err != nil {
		fmt.Println(err.Error())
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

	logs := []Log{}
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

		tempLog := Log{
			TableName:   "tags",
			TableNameId: tagValue[i],
			ColumnName:  "updated_by",
			Action:      ActionChangeLogUpdate,
			Old:         ValueChangeLogNone,
			New:         user,
			CreatedBy:   user,
		}
		logs = append(logs, tempLog)

		tempLog = Log{
			TableName:   "tags",
			TableNameId: tagValue[i],
			ColumnName:  "updated_at",
			Action:      ActionChangeLogUpdate,
			Old:         ValueChangeLogNone,
			New:         time.Now().String(),
			CreatedBy:   user,
		}
		logs = append(logs, tempLog)

		tempLog = Log{
			TableName:   "tags",
			TableNameId: tagValue[i],
			ColumnName:  "status",
			Action:      ActionChangeLogUpdate,
			Old:         ValueChangeLogNone,
			New:         StatusDeleted,
			CreatedBy:   user,
		}
		logs = append(logs, tempLog)
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

	err = addLogs(logs)
	if err != nil {
		fmt.Println(err.Error())
	}

	return nil
}
