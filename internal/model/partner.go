package model

import (
	"fmt"
	"time"
)

type (
	Partner struct {
		Id           string `db:"id"`
		PartnerName  string `db:"partner_name"`
		SerialNumber string `db:"serial_number"`
		CreatedBy    string `db:"created_by"`
	}
)

func InsertPartner(p Partner) error {
	fmt.Println("Add")
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	partner, err := checkPartner(p.PartnerName)

	if partner == 0 {
		q := `
			INSERT INTO partners(
				partner_name
				, serial_number
				, created_by
				, status
			)
			VALUES (?, ?, ?, ?)
		`

		_, err := tx.Exec(tx.Rebind(q), p.PartnerName, p.SerialNumber, p.CreatedBy, StatusCreated)
		if err != nil {
			return err
		}

		if err := tx.Commit(); err != nil {
			return err
		}
		return nil
	} else {
		return ErrDuplicateEntry
	}
}

func FindAllPartner() (Response, error) {
	fmt.Println("Select partner")
	q := `
		SELECT id, partner_name, serial_number
		FROM partners
		WHERE status = ?
	`

	var resv []Partner
	if err := db.Select(&resv, db.Rebind(q), StatusCreated); err != nil {
		return Response{Status: "Error", Message: q, Data: []Partner{}}, err
	}
	if len(resv) < 1 {
		return Response{Status: "404", Message: q, Data: []Partner{}}, ErrResourceNotFound
	}

	res := Response{
		Status:  "200",
		Message: "Ok",
		Data:    resv,
	}

	return res, nil
}

func FindPartnerSerialNumber(param string) (Response, error) {
	fmt.Println("Select partner")
	q := `
		SELECT serial_number
		FROM partners
		WHERE (id ILIKE '%?%'
		OR partner_name ILIKE '%?%' )
		AND status = ?
	`

	var resv []Partner
	if err := db.Select(&resv, db.Rebind(q), param, param, StatusCreated); err != nil {
		return Response{Status: "Error", Message: q, Data: []Partner{}}, err
	}
	if len(resv) < 1 {
		return Response{Status: "404", Message: q, Data: []Partner{}}, ErrResourceNotFound
	}

	res := Response{
		Status:  "200",
		Message: "Ok",
		Data:    resv,
	}

	return res, nil
}

func DeletePartner(partnerId, userId string) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	q := `
		UPDATE partners
		SET
			updated_by = ?
			, updated_at = ?
			, status = ?
		WHERE
			id = ?
			AND status = ?;
	`
	_, err = tx.Exec(tx.Rebind(q), userId, time.Now(), StatusDeleted, partnerId, StatusCreated)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func checkPartner(name string) (int, error) {
	q := `
		SELECT id
		FROM partners
		WHERE
			partner_name = ?
			AND status = ?
	`
	var res []string
	if err := db.Select(&res, db.Rebind(q), name, StatusCreated); err != nil {
		return 0, err
	}

	return len(res), nil
}
