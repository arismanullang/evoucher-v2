package model

import (
	"fmt"
	"time"
)

type (
	BroadcastUser struct {
		ID          int       `db:"id"`
		State       string    `db:"state"`
		ProgramID   string    `db:"program_id"`
		Target      string    `db:"target"`
		Description string    `db:"description"`
		CreatedBy   string    `db:"created_by"`
		CreatedAt   time.Time `db:"created_at"`
	}
)

func FindBroadcastUser(param map[string]string) ([]BroadcastUser, error) {
	q := `
		SELECT
			id
			, state
			, program_id
			, target
			, description
		FROM
			broadcast_users
		WHERE
			status = ?
	`
	for key, value := range param {
		q += ` AND ` + key + ` = '` + value + `'`
	}
	var resd []BroadcastUser
	if err := db.Select(&resd, db.Rebind(q), StatusCreated); err != nil {
		return []BroadcastUser{}, err
	} else if len(resd) < 1 {
		return []BroadcastUser{}, ErrResourceNotFound
	}
	return resd, nil
}
func UpdateBroadcastUserState(programId, email string) error {
	tx, err := db.Beginx()
	if err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}
	defer tx.Rollback()

	q := `
		UPDATE broadcast_users
		SET
			state = ?
		WHERE
			program_id = ?
			AND target = ?
			AND status = ?
	`

	_, err = tx.Exec(tx.Rebind(q), EmailSend, programId, email, StatusCreated)
	if err != nil {
		fmt.Println("Update User State | " + err.Error())
		return ErrServerInternal
	}

	if err := tx.Commit(); err != nil {
		fmt.Println("Update User State | " + err.Error())
		return ErrServerInternal
	}
	return nil
}
