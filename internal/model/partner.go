package model

import (
	"database/sql"
	"fmt"
	"time"
)

type (
	Partner struct {
		Id           string         `db:"id" json:"id"`
		PartnerName  string         `db:"partner_name" json:"partner_name"`
		SerialNumber sql.NullString `db:"serial_number" json:"serial_number"`
		CreatedBy    sql.NullString `db:"created_by" json:"created_by"`
		VariantID    string         `db:"variant_id" json:"variant_id"`
		Tag          sql.NullString `db:"tag" json:"tag"`
		Description  sql.NullString `db:"description" json:"description"`
	}

	Tag struct {
		Value string `db:"tag_value" json:"tag_value"`
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

	partner, err := checkPartner(p.PartnerName)
	if err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}

	if partner == "" {
		q := `
			INSERT INTO partners(
				partner_name
				, serial_number
				, tag
				, description
				, created_by
				, status
			)
			VALUES (?, ?, ?, ?, ?, ?)
		`

		_, err := tx.Exec(tx.Rebind(q), p.PartnerName, p.SerialNumber, p.Tag, p.Description, p.CreatedBy, StatusCreated)
		if err != nil {
			fmt.Println(err.Error())
			return ErrServerInternal
		}

		if err := tx.Commit(); err != nil {
			fmt.Println(err.Error())
			return ErrServerInternal
		}
		return nil
	} else {
		return ErrDuplicateEntry
	}
}

func checkPartner(name string) (string, error) {
	q := `
		SELECT id
		FROM partners
		WHERE
			partner_name = ?
			AND status = ?
	`
	var res []string
	if err := db.Select(&res, db.Rebind(q), name, StatusCreated); err != nil {
		return "", ErrServerInternal
	}
	if len(res) == 0 {
		return "", nil
	}
	return res[0], nil
}

func FindPartners(param map[string]string) ([]Partner, error) {
	fmt.Println("Select partner")
	q := `
		SELECT
			id
			, partner_name
			, serial_number
			, tag
			, description
		FROM partners
		WHERE status = ?
	`
	for key, value := range param {
		if key == "q" {
			q += `AND (id ILIKE '%` + value + `%' OR partner_name ILIKE '%` + value + `%' OR serial_number ILIKE '%` + value + `%')`
		} else {
			q += ` AND ` + key + ` = '` + value + `'`
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

	return resv, nil
}

func FindAllPartners() ([]Partner, error) {
	fmt.Println("Select partner")
	q := `
		SELECT
			id
			, partner_name
			, serial_number
			, tag
			, description
		FROM partners
		WHERE status = ?
	`

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

func UpdatePartner(partnerId, serialNumber, user string) error {
	tx, err := db.Beginx()
	if err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}
	defer tx.Rollback()

	q := `
		UPDATE partners
		SET
			serial_number = ?
			, updated_by = ?
			, updated_at = ?
		WHERE
			id = ?
			AND status = ?;
	`
	_, err = tx.Exec(tx.Rebind(q), serialNumber, user, time.Now(), partnerId, StatusCreated)
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

func DeletePartner(partnerId, userId string) error {
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
	_, err = tx.Exec(tx.Rebind(q), userId, time.Now(), StatusDeleted, partnerId)
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

func FindVariantPartner(param map[string]string) ([]Partner, error) {
	q := `
		SELECT 	b.id
			, b.partner_name
			, b.serial_number
			, b.created_by
			, a.variant_id
	 	FROM
			variant_partners a
		JOIN
		 	partners b
		ON
			a.partner_id = b.id
 		WHERE
			a.status = ?
	`
	for k, v := range param {
		table := "b"
		if k == "variant_id" {
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
	fmt.Println(resv)
	return resv, nil
}

// ------------------------------------------------------------------------------
// Tag

func FindAllTags() ([]string, error) {
	fmt.Println("Select partner")
	q := `
		SELECT
			tag_value
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
			tag_value
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
	return nil
}

func DeleteTag(tagValue, user string) error {
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
			tag_value = ?;
	`
	_, err = tx.Exec(tx.Rebind(q), user, time.Now(), StatusDeleted, tagValue)
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
			tag_value = ?
	`
	for i := 1; i < len(tagValue); i++ {
		q += " OR tag_value = '" + tagValue[i] + "'"
	}
	fmt.Println(q)
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
