package model

import (
	"fmt"
	"time"
)

type (
	Role struct {
		Id        string   `db:"id" json:"id"`
		Detail    string   `db:"detail" json:"detail"`
		Features  []string `db:"-" json:"features"`
		AccountId string   `db:"account_id" json:"account_id"`
	}
)

// Role -----------------------------------------------------------------------------------------------

func FindAllRole(accountId string) ([]Role, error) {
	q := `
		SELECT id, detail
		FROM roles
		WHERE status = ?
		AND NOT detail = 'suadmin'
		AND account_id = ?
	`

	var resv []Role
	if err := db.Select(&resv, db.Rebind(q), StatusCreated, accountId); err != nil {
		fmt.Println(err)
		return []Role{}, ErrServerInternal
	}
	if len(resv) < 1 {
		return []Role{}, ErrResourceNotFound
	}

	return resv, nil
}

func AddRole(r Role, user string) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	q := `
		INSERT INTO roles(
			detail
			, account_id
			, created_by
			, status
		)
		VALUES (?, ?, ?, ?)
		RETURNING
			id
	`
	var res []string
	if err := tx.Select(&res, tx.Rebind(q), r.Detail, r.AccountId, user, StatusCreated); err != nil {
		fmt.Println(err)
		return ErrServerInternal
	}

	for _, v := range r.Features {
		q := `
			INSERT INTO role_features(
				role_id
				, feature_id
				, created_by
				, status
			)
			VALUES (?, ?, ?, ?)
		`

		_, err := tx.Exec(tx.Rebind(q), res[0], v, user, StatusCreated)
		if err != nil {
			fmt.Println(err)
			return ErrServerInternal
		}
	}

	if err := tx.Commit(); err != nil {
		fmt.Println(err)
		return ErrServerInternal
	}

	return nil
}

func AddAdmin(account, user string) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	q := `
		INSERT INTO roles(
			detail
			, account_id
			, created_by
			, status
		)
		VALUES (?, ?, ?, ?)
		RETURNING
			id
	`
	var res []string
	if err := tx.Select(&res, tx.Rebind(q), "admin", account, user, StatusCreated); err != nil {
		fmt.Println(err)
		return ErrServerInternal
	}

	features, err := GetAllFeatures()
	if err != nil {
		fmt.Println(err)
		return ErrServerInternal
	}

	for _, v := range features {
		q := `
			INSERT INTO role_features(
				role_id
				, feature_id
				, created_by
				, status
			)
			VALUES (?, ?, ?, ?)
		`

		_, err := tx.Exec(tx.Rebind(q), res[0], v.Id, user, StatusCreated)
		if err != nil {
			fmt.Println(err)
			return ErrServerInternal
		}
	}

	if err := tx.Commit(); err != nil {
		fmt.Println(err)
		return ErrServerInternal
	}

	return nil
}

func GetRoleDetail(id string) (Role, error) {
	q := `
		SELECT
			id
			, detail
			, account_id
		FROM roles
		WHERE
			status = ?
			AND id = ?
	`

	var resv []Role
	if err := db.Select(&resv, db.Rebind(q), StatusCreated, id); err != nil {
		fmt.Println(err)
		return Role{}, ErrServerInternal
	}
	q = `
		SELECT
			f.id
		FROM roles AS r
		JOIN role_features AS rf
		ON
			r.id = rf.role_id
		JOIN features AS f
		ON
			f.id = rf.feature_id
		WHERE
			rf.status = ?
			AND rf.role_id = ?
	`

	var res []string
	if err := db.Select(&res, db.Rebind(q), StatusCreated, id); err != nil {
		return Role{}, ErrServerInternal
	}

	resv[0].Features = res

	return resv[0], nil
}

func UpdateRole(r Role, user string) error {
	tx, err := db.Beginx()
	if err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}
	defer tx.Rollback()

	q := `
		UPDATE role_features
		SET
			status = ?
			, updated_by = ?
			, updated_at = ?
		WHERE
			role_id = ?
			AND status = ?
	`

	_, err = tx.Exec(tx.Rebind(q), StatusDeleted, user, time.Now(), r.Id, StatusCreated)
	if err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}

	for _, v := range r.Features {
		q := `
				INSERT INTO role_features(
					role_id
					, feature_id
					, created_by
					, status
				)
				VALUES (?, ?, ?, ?)
			`

		_, err := tx.Exec(tx.Rebind(q), r.Id, v, user, StatusCreated)
		if err != nil {
			fmt.Println(err)
			return ErrServerInternal
		}
	}

	if err := tx.Commit(); err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}

	return nil
}
