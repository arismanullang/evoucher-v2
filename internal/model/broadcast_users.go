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
	EmailUser struct {
		ID        string    `db:"id" json:"id"`
		Name      string    `db:"name" json:"name"`
		Email     string    `db:"email" json:"email"`
		AccountID string    `db:"account_id" json:"account_id"`
		CreatedBy string    `db:"created_by" json:"created_by"`
		CreatedAt time.Time `db:"created_at" json:"created_at"`
		UpdatedBy string    `db:"updated_by" json:"updated_by"`
		UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	}
	ListEmailUser struct {
		ID         string      `db:"id" json:"id"`
		Name       string      `db:"name" json:"name"`
		AccountID  string      `db:"account_id" json:"account_id"`
		EmailUsers []EmailUser `db:"-" json:"email_users"`
		CreatedBy  string      `db:"created_by" json:"created_by"`
		CreatedAt  time.Time   `db:"created_at" json:"created_at"`
		UpdatedBy  string      `db:"updated_by" json:"updated_by"`
		UpdatedAt  time.Time   `db:"updated_at" json:"updated_at"`
	}
	InsertCampaignUserRequest struct {
		CampaignID string    `db:"campaign_id" json:"campaign_id"`
		EmailUsers []string  `db:"email_user_id" json:"emails"`
		CreatedBy  string    `db:"created_by" json:"created_by"`
		CreatedAt  time.Time `db:"created_at" json:"created_at"`
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

func UpdateBroadcastUserState(programId, email, user string) error {
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
			, updated_by = ?
			, updated_at = ?
		WHERE
			program_id = ?
			AND target = ?
			AND status = ?
	`

	_, err = tx.Exec(tx.Rebind(q), EmailSend, user, time.Now(), programId, email, StatusCreated)
	if err != nil {
		fmt.Println("Update User State | " + err.Error())
		return ErrServerInternal
	}

	if err := tx.Commit(); err != nil {
		fmt.Println("Update User State | " + err.Error())
		return ErrServerInternal
	}

	logs := []Log{}
	log := Log{
		TableName:   "broadcast_users",
		TableNameId: programId,
		ColumnName:  "state",
		Action:      ActionChangeLogUpdate,
		Old:         ValueChangeLogNone,
		New:         EmailSend,
		CreatedBy:   user,
	}
	logs = append(logs, log)

	log = Log{
		TableName:   "broadcast_users",
		TableNameId: programId,
		ColumnName:  "updated_by",
		Action:      ActionChangeLogUpdate,
		Old:         ValueChangeLogNone,
		New:         user,
		CreatedBy:   user,
	}
	logs = append(logs, log)

	log = Log{
		TableName:   "broadcast_users",
		TableNameId: programId,
		ColumnName:  "updated_at",
		Action:      ActionChangeLogUpdate,
		Old:         ValueChangeLogNone,
		New:         time.Now().String(),
		CreatedBy:   user,
	}
	logs = append(logs, log)

	err = addLogs(logs)
	if err != nil {
		fmt.Println(err.Error())
	}

	return nil
}

// Email user ------------------------------------------------------------------------------------------------------------
func InsertEmailUser(param EmailUser) (string, error) {
	tx, err := db.Beginx()
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	ex, err := checkEmailUser(param)
	if ex != "" && err == nil {
		return "", ErrDuplicateEntry
	}

	q := `
		INSERT INTO email_users(
			name
			, email
			, account_id
			, created_by
			, created_at
			, status
		)
		VALUES (?, ?, ?, ?, ?, ?)
		RETURNING
			id
	`

	var res []string
	if err := tx.Select(&res, tx.Rebind(q), param.Name, param.Email, param.AccountID, param.CreatedBy, time.Now(), StatusCreated); err != nil {
		return "", err
	}
	if err := tx.Commit(); err != nil {
		return "", err
	}

	return res[0], nil
}

func InsertEmailUsers(param []EmailUser) ([]string, error) {
	tx, err := db.Beginx()
	if err != nil {
		return []string{}, err
	}
	defer tx.Rollback()

	ids := []string{}
	for _, v := range param {
		q := `
			INSERT INTO email_users(
				name
				, email
				, account_id
				, created_by
				, created_at
				, status
			)
			VALUES (?, ?, ?, ?, ?, ?)
			RETURNING
				id
		`

		var res []string
		if err := tx.Select(&res, tx.Rebind(q), v.Name, v.Email, v.AccountID, v.CreatedBy, time.Now(), StatusCreated); err != nil {
			return []string{}, err
		}

		ids = append(ids, res[0])
	}
	if err := tx.Commit(); err != nil {
		return []string{}, err
	}

	return ids, nil
}

func checkEmailUser(param EmailUser) (string, error) {
	q := `
		SELECT
			id
		FROM
			email_users
		WHERE
			email = ?
			AND account_id = ?
			AND status = ?
	`
	var res []string
	if err := db.Select(&res, db.Rebind(q), param.Email, param.AccountID, StatusCreated); err != nil {
		fmt.Println(err)
		return "", ErrServerInternal
	}
	if len(res) == 0 {
		return "", nil
	}
	return res[0], nil
}

func GetAllEmailUser(accountId string) ([]EmailUser, error) {
	q := `
		SELECT
			id
			, name
			, email
		FROM
			email_users
		WHERE
			account_id = ?
			AND status = ?
	`
	var res []EmailUser
	if err := db.Select(&res, db.Rebind(q), accountId, StatusCreated); err != nil {
		fmt.Println(err)
		return []EmailUser{}, ErrServerInternal
	}
	if len(res) == 0 {
		return []EmailUser{}, nil
	}
	return res, nil
}

func GetEmailUser(param map[string]string) ([]EmailUser, error) {
	q := `
		SELECT
			id
			, name
			, email
		FROM
			email_users
		WHERE
			status = ?
	`
	for key, value := range param {
		q += ` AND ` + key + ` LIKE '%` + value + `%'`
	}

	var res []EmailUser
	if err := db.Select(&res, db.Rebind(q), StatusCreated); err != nil {
		fmt.Println(err)
		return []EmailUser{}, ErrServerInternal
	}
	if len(res) == 0 {
		return []EmailUser{}, nil
	}
	return res, nil
}

func GetEmailUserByIDs(id []string) ([]EmailUser, error) {
	q := `
		SELECT DISTINCT
			eu.id
			, eu.name
			, eu.email
		FROM
			email_users as eu
		JOIN
			list_users as lu
		ON
			eu.id = lu.email_user_id
		WHERE
			lu.email_user_id = ?
	`
	if len(id) > 1 {
		for _, value := range id {
			q += ` OR lu.email_user_id LIKE '%` + value + `%'`
		}
	}

	var res []EmailUser
	if err := db.Select(&res, db.Rebind(q), id[0]); err != nil {
		fmt.Println(err)
		return []EmailUser{}, ErrServerInternal
	}
	if len(res) == 0 {
		return []EmailUser{}, nil
	}
	return res, nil
}

func GetEmailUserByListIDs(id []string) ([]EmailUser, error) {
	q := `
		SELECT DISTINCT
			eu.id
			, eu.name
			, eu.email
		FROM
			email_users as eu
		JOIN
			list_users as lu
		ON
			eu.id = lu.email_user_id
		WHERE
			eu.status = ?
			AND lu.status = ?
			AND (lu.list_email_user_id = ?
	`
	if len(id) > 1 {
		for _, value := range id {
			q += ` OR lu.list_email_user_id LIKE '%` + value + `%'`
		}
	}
	q += `)`

	var res []EmailUser
	if err := db.Select(&res, db.Rebind(q), StatusCreated, StatusCreated, id[0]); err != nil {
		fmt.Println(err)
		return []EmailUser{}, ErrServerInternal
	}
	if len(res) == 0 {
		return []EmailUser{}, nil
	}
	return res, nil
}

func DeleteEmailUser(emailUser, user string) error {
	tx, err := db.Beginx()
	if err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}
	defer tx.Rollback()

	q := `
		UPDATE email_users
		SET
			status = ?
			, updated_by = ?
			, updated_at = ?
		WHERE
			id = ?
			AND status = ?
	`

	_, err = tx.Exec(tx.Rebind(q), StatusDeleted, user, time.Now(), emailUser, StatusCreated)
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

// List email----------------------------------------------------------------------------------------------------------------
func InsertListEmail(param ListEmailUser) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	q := `
		INSERT INTO list_email_users(
			name
			, account_id
			, created_by
			, created_at
			, status
		)
		VALUES (?, ?, ?, ?, ?)
		RETURNING
			id
	`

	var resListEmailUser []string
	if err := tx.Select(&resListEmailUser, tx.Rebind(q), param.Name, param.AccountID, param.CreatedBy, time.Now(), StatusCreated); err != nil {
		return err
	}

	for _, v := range param.EmailUsers {
		q = `
			INSERT INTO list_users(
				list_email_user_id
				, email_user_id
				, created_by
				, created_at
				, status
			)
			VALUES (?, ?, ?, ?, ?)
			RETURNING
				id
		`

		var resEmailUser []string
		if err := tx.Select(&resEmailUser, tx.Rebind(q), resListEmailUser[0], v.ID, param.CreatedBy, time.Now(), StatusCreated); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func AddNewEmailUserToList(param EmailUser, listId string) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	emailUserId, err := InsertEmailUser(param)
	if err != nil {
		return err
	}

	q := `
		INSERT INTO list_users(
			list_email_user_id
			, email_user_id
			, created_by
			, created_at
			, status
		)
		VALUES (?, ?, ?, ?, ?)
		RETURNING
			id
	`

	var res []string
	if err := tx.Select(&res, tx.Rebind(q), listId, emailUserId, param.CreatedBy, time.Now(), StatusCreated); err != nil {
		DeleteEmailUser(emailUserId, param.CreatedBy)
		return err
	}

	if err := tx.Commit(); err != nil {
		DeleteEmailUser(emailUserId, param.CreatedBy)
		return err
	}

	return nil
}

func AddEmailUserToList(emailUserId []string, listId, accountId, user string) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	q := ``
	for _, v := range emailUserId {
		q = `
			INSERT INTO list_users(
				list_email_user_id
				, email_user_id
				, created_by
				, created_at
				, status
			)
			VALUES (?, ?, ?, ?, ?)
			RETURNING
				id
		`

		var res []string
		if err := tx.Select(&res, tx.Rebind(q), listId, v, user, time.Now(), StatusCreated); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func RemoveEmailUserFromList(userId, listId, user string) error {
	tx, err := db.Beginx()
	if err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}
	defer tx.Rollback()

	q := `
		UPDATE list_users
		SET
			status = ?
			, updated_by = ?
			, updated_at = ?
		WHERE
			list_email_user_id = ?
			AND email_user_id = ?
			AND status = ?
	`

	_, err = tx.Exec(tx.Rebind(q), StatusDeleted, userId, time.Now(), listId, userId, StatusCreated)
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

func GetAllListEmailUser(accountId string) ([]ListEmailUser, error) {
	q := `
		SELECT
			id
			, name
		FROM
			list_email_users
		WHERE
			account_id = ?
			AND status = ?
	`
	var res []ListEmailUser
	if err := db.Select(&res, db.Rebind(q), accountId, StatusCreated); err != nil {
		fmt.Println(err)
		return []ListEmailUser{}, ErrServerInternal
	}
	if len(res) == 0 {
		return []ListEmailUser{}, nil
	}

	for i, v := range res {
		q = `
			SELECT
				eu.id
				, eu.name
				, eu.email
			FROM
				email_users as eu
			JOIN
				list_users as lu
			ON
				lu.email_user_id = eu.id
			WHERE
				lu.list_email_user_id = ?
				AND eu.status = ?
				AND lu.status = ?
		`
		var resE []EmailUser
		if err := db.Select(&resE, db.Rebind(q), v.ID, StatusCreated, StatusCreated); err != nil {
			fmt.Println(err)
		}

		res[i].EmailUsers = resE
	}
	return res, nil
}

func GetListEmailUserById(id, accountId string) (ListEmailUser, error) {
	q := `
		SELECT
			id
			, name
		FROM
			list_email_users
		WHERE
			account_id = ?
			AND id = ?
			AND status = ?
	`
	var res []ListEmailUser
	if err := db.Select(&res, db.Rebind(q), accountId, id, StatusCreated); err != nil {
		fmt.Println(err)
		return ListEmailUser{}, ErrServerInternal
	}
	if len(res) == 0 {
		return ListEmailUser{}, nil
	}

	q = `
		SELECT
			eu.id
			, eu.name
			, eu.email
		FROM
			email_users as eu
		JOIN
			list_users as lu
		ON
			lu.email_user_id = eu.id
		WHERE
			lu.list_email_user_id = ?
			AND eu.status = ?
			AND lu.status = ?
	`
	var resE []EmailUser
	if err := db.Select(&resE, db.Rebind(q), res[0].ID, StatusCreated, StatusCreated); err != nil {
		fmt.Println(err)
	}

	res[0].EmailUsers = resE

	return res[0], nil
}

func GetListEmailUserByIds(id []string, accountId string) ([]ListEmailUser, error) {
	q := `
		SELECT
			id
			, name
		FROM
			list_email_users
		WHERE
			account_id = ?
			AND status = ?
			AND (
	`

	q += ` id LIKE '%` + id[0] + `%'`
	if len(id) > 1 {
		for _, v := range id {
			q += ` OR id LIKE '%` + v + `%'`
		}
	}
	q += `)`

	var res []ListEmailUser
	if err := db.Select(&res, db.Rebind(q), accountId, StatusCreated); err != nil {
		fmt.Println(err)
		return []ListEmailUser{}, ErrServerInternal
	}
	if len(res) == 0 {
		return []ListEmailUser{}, nil
	}

	return res, nil
}

func DeleteListUser(emailUser, user string) error {
	tx, err := db.Beginx()
	if err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}
	defer tx.Rollback()

	q := `
		UPDATE list_email_users
		SET
			status = ?
			, updated_by = ?
			, updated_at = ?
		WHERE
			id = ?
			AND status = ?
	`

	_, err = tx.Exec(tx.Rebind(q), StatusDeleted, user, time.Now(), emailUser, StatusCreated)
	if err != nil {
		fmt.Println("Update List | " + err.Error())
		return ErrServerInternal
	}

	if err := tx.Commit(); err != nil {
		fmt.Println("Update List | " + err.Error())
		return ErrServerInternal
	}

	return nil
}

// Campaign user
func InsertCampaignUser(param InsertCampaignUserRequest) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, v := range param.EmailUsers {
		q := `
			INSERT INTO campaign_users(
				campaign_id
				, email_user_id
				, created_by
				, created_at
				, status
				, state
			)
			VALUES (?, ?, ?, ?, ?, ?)
			RETURNING
				id
		`

		var res []string
		if err := tx.Select(&res, tx.Rebind(q), param.CampaignID, v, param.CreatedBy, time.Now(), StatusCreated, VoucherStateSend); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
