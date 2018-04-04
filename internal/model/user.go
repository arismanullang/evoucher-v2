package model

import (
	"fmt"
	"time"
)

type (
	SuperAdminRegisterUser struct {
		AccountId string   `db:"account_id" json:"account_id"`
		Username  string   `db:"username" json:"username"`
		Password  string   `db:"password" json:"password"`
		Email     string   `db:"email" json:"email"`
		Phone     string   `db:"phone" json:"phone"`
		Role      []string `db:"-" json:"role"`
	}
	SuperAdminUser struct {
		ID        string  `db:"id" json:"id"`
		Account   Account `db:"_" json:"account"`
		Username  string  `db:"username" json:"username"`
		Password  string  `db:"password" json:"password"`
		Email     string  `db:"email" json:"email"`
		Phone     string  `db:"phone" json:"phone"`
		Role      []Role  `db:"-" json:"role"`
		Status    string  `db:"status" json:"status"`
		CreatedBy string  `db:"created_by" json:"created_by"`
		CreatedAt string  `db:"created_at" json:"created_at"`
	}
	RegisterUser struct {
		Username string   `db:"username" json:"username"`
		Password string   `db:"password" json:"password"`
		Email    string   `db:"email" json:"email"`
		Phone    string   `db:"phone" json:"phone"`
		Role     []string `db:"-" json:"role"`
	}
	User struct {
		ID        string  `db:"id" json:"id"`
		Account   Account `db:"_" json:"account"`
		Username  string  `db:"username" json:"username"`
		Password  string  `db:"password" json:"password"`
		Email     string  `db:"email" json:"email"`
		Phone     string  `db:"phone" json:"phone"`
		Role      []Role  `db:"-" json:"role"`
		Status    string  `db:"status" json:"status"`
		CreatedBy string  `db:"created_by" json:"created_by"`
		CreatedAt string  `db:"created_at" json:"created_at"`
	}

	UserRes struct {
		Id       string `db:"id"`
		Username string `db:"username"`
	}
)

func AddUser(u RegisterUser, user, accountId string) error {
	fmt.Println("Add")
	tx, err := db.Beginx()
	if err != nil {
		fmt.Println(err)
		return ErrServerInternal
	}
	defer tx.Rollback()

	username, err := CheckUsername(u.Username, accountId)

	if username != "" {
		return ErrDuplicateEntry
	}

	q := `
		INSERT INTO users(
			username
			, password
			, email
			, phone
			, created_by
			, status
		)
		VALUES (?, ?, ?, ?, ?, ?)
		RETURNING
			id
	`
	var res []string
	if err := tx.Select(&res, tx.Rebind(q), u.Username, u.Password, u.Email, u.Phone, user, StatusCreated); err != nil {
		fmt.Println(err)
		return ErrServerInternal
	}

	logs := []Log{}
	tempLog := Log{
		TableName:   "users",
		TableNameId: ValueChangeLogNone,
		ColumnName:  ColumnChangeLogInsert,
		Action:      ActionChangeLogInsert,
		Old:         ValueChangeLogNone,
		New:         res[0],
		CreatedBy:   user,
	}
	logs = append(logs, tempLog)

	for _, v := range u.Role {
		q := `
			INSERT INTO user_roles(
				user_id
				, role_id
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

		tempLog = Log{
			TableName:   "user_roles",
			TableNameId: res[0],
			ColumnName:  ColumnChangeLogInsert,
			Action:      ActionChangeLogInsert,
			Old:         ValueChangeLogNone,
			New:         v,
			CreatedBy:   user,
		}
		logs = append(logs, tempLog)
	}

	q2 := `
		INSERT INTO user_accounts(
			user_id
			, account_id
			, created_by
			, status
		)
		VALUES (?, ?, ?, ?)
	`

	_, err = tx.Exec(tx.Rebind(q2), res[0], accountId, user, StatusCreated)
	if err != nil {
		fmt.Println(err)
		return ErrServerInternal
	}

	if err = tx.Commit(); err != nil {
		fmt.Println(err)
		return ErrServerInternal
	}

	tempLog = Log{
		TableName:   "user_accounts",
		TableNameId: res[0],
		ColumnName:  ColumnChangeLogInsert,
		Action:      ActionChangeLogInsert,
		Old:         ValueChangeLogNone,
		New:         accountId,
		CreatedBy:   user,
	}
	logs = append(logs, tempLog)

	err = addLogs(logs)
	if err != nil {
		fmt.Println(err.Error())
	}

	return nil
}

func CheckUsername(username, accountId string) (string, error) {
	q := `
		SELECT
			u.id
		FROM
			users as u
		JOIN
			user_accounts as ua
		ON
			u.id = ua.user_id
		WHERE
			u.username = ?
			AND ua.account_id = ?
	`
	var res []string
	if err := db.Select(&res, db.Rebind(q), username, accountId); err != nil {
		fmt.Println(err)
		return "", ErrServerInternal
	}
	if len(res) == 0 {
		return "", nil
	}
	return res[0], nil
}

func GetWebuser() (string, error) {
	q := `
		SELECT
			u.id
		FROM
			users as u
		WHERE
			u.username = 'web'
	`
	var res []string
	if err := db.Select(&res, db.Rebind(q)); err != nil {
		fmt.Println(err)
		return "", ErrServerInternal
	}
	if len(res) == 0 {
		return "", nil
	}
	return res[0], nil
}

func FindAllUsers(accountId string) ([]User, error) {
	q := `
		SELECT
			u.id
			, u.username
			, u.email
			, u.phone
			, u.created_at
			, u.created_by
			, u.status
		FROM
			users as u
		JOIN
			user_accounts as ua
		ON
			u.id = ua.user_id
		WHERE
			ua.account_id = ?
	`
	var res []User
	if err := db.Select(&res, db.Rebind(q), accountId); err != nil {
		fmt.Println(err)
		return []User{}, ErrServerInternal
	}

	for i, v := range res {
		q1 := `
		SELECT
		  	u.id id,
		  	r.detail detail
		FROM
		  	users u
		JOIN
			user_roles ur
		ON
			u.id = ur.user_id
		JOIN
			roles r
	    	ON
	    		ur.role_id = r.id
		WHERE
		  	u.id = ?

	`
		var role []Role
		if err := db.Select(&role, db.Rebind(q1), v.ID); err != nil {
			fmt.Println(err)
			return []User{}, ErrServerInternal
		}
		res[i].Role = role
		account, err := GetAccountDetailByAccountId(accountId)
		if err != nil {
			fmt.Println(err)
			return []User{}, ErrServerInternal
		}
		res[i].Account = account
	}

	return res, nil
}

func FindUsersByRole(role, accountId string) ([]UserRes, error) {
	q := `
		SELECT u.id, u.username FROM users AS u
		JOIN user_accounts AS ua ON u.id = ua.user_id
		JOIN user_roles AS ur ON u.id = ur.user_id
		WHERE ua.account_id = ?
		AND ur.role_id = ?
		AND u.status = ?
	`

	var resv []UserRes
	if err := db.Select(&resv, db.Rebind(q), accountId, role, StatusCreated); err != nil {
		fmt.Println(err)
		return []UserRes{}, ErrServerInternal
	}
	if len(resv) < 1 {
		return []UserRes{}, ErrResourceNotFound
	}

	return resv, nil
}

func FindUsersCustomParam(usr map[string]string) ([]UserRes, error) {
	q := `
		SELECT u.id, u.username FROM users AS u
		JOIN user_accounts AS ua ON u.id = ua.user_id
		JOIN user_roles AS ur ON u.id = ur.user_id
		WHERE
			status = ?
	`

	for key, value := range usr {
		if key == "q" {
			q += `AND (u.username ILIKE '%` + value + `%')`
		} else {
			q += ` AND ` + key + ` LIKE '%` + value + `%'`
		}
	}

	var resv []UserRes
	if err := db.Select(&resv, db.Rebind(q), StatusCreated); err != nil {
		fmt.Println(err)
		return []UserRes{}, ErrServerInternal
	}
	if len(resv) < 1 {
		return []UserRes{}, ErrResourceNotFound
	}

	return resv, nil
}

func FindUserDetail(userId string) (User, error) {
	q := `
		SELECT
			id
			, username
			, email
			, phone
			, created_at
		FROM
			users as u
		WHERE
			id = ?
			OR username = ?
			AND status = ?
	`
	var res []User
	if err := db.Select(&res, db.Rebind(q), userId, userId, StatusCreated); err != nil {
		fmt.Println(err)
		return User{}, ErrServerInternal
	}
	if len(res) == 0 {
		return User{}, ErrResourceNotFound
	}

	q1 := `
		SELECT
		  	r.id id,
		  	r.detail detail
		FROM
		  	users u
		JOIN
			user_roles ur
		ON
			u.id = ur.user_id
		JOIN
			roles r
	    	ON
	    		ur.role_id = r.id
		WHERE
		  	u.id = ?
	  		AND ur.status = ?

	`
	var role []Role
	if err := db.Select(&role, db.Rebind(q1), res[0].ID, StatusCreated); err != nil {
		fmt.Println(err)
		return User{}, ErrServerInternal
	}

	res[0].Role = role
	account, err := GetAccountDetailByUser(userId)
	if err != nil {
		fmt.Println(err)
		return User{}, ErrServerInternal
	}
	res[0].Account = account

	return res[0], nil
}

func Login(username, password string) (string, error) {
	q := `
		SELECT
			id
		FROM
			users
		WHERE
			username = ?
			AND password = ?
			AND status = ?
	`
	var res []string
	if err := db.Select(&res, db.Rebind(q), username, password /*accountId,*/, StatusCreated); err != nil {
		fmt.Println(err)
		return "", ErrServerInternal
	}
	if len(res) == 0 {
		return "", ErrResourceNotFound
	}

	log := Log{
		TableName:   "users",
		TableNameId: res[0],
		ColumnName:  ColumnChangeLogSelect,
		Action:      ActionChangeLogLogin,
		Old:         ValueChangeLogNone,
		New:         ValueChangeLogNone,
		CreatedBy:   res[0],
	}

	err := addLog(log)
	if err != nil {
		fmt.Println(err.Error())
	}

	return res[0], nil
}

func UpdatePassword(id, password string) error {
	tx, err := db.Beginx()
	if err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}
	defer tx.Rollback()

	q := `
		UPDATE users
		SET
			password = ?
		WHERE
			id = ?
			AND status = ?
	`

	_, err = tx.Exec(tx.Rebind(q), password, id, StatusCreated)
	if err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}

	if err := tx.Commit(); err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}

	log := Log{
		TableName:   "users",
		TableNameId: id,
		ColumnName:  "password",
		Action:      ActionChangeLogUpdate,
		Old:         ValueChangeLogNone,
		New:         password,
		CreatedBy:   id,
	}

	err = addLog(log)
	if err != nil {
		fmt.Println(err.Error())
	}

	return nil
}

func ChangePassword(id, oldPassword, newPassword string) error {
	tx, err := db.Beginx()
	if err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}
	defer tx.Rollback()

	q := `
		SELECT
			id
		FROM
			users
		WHERE
			id = ?
			AND password = ?
			AND status = ?
	`
	var res []string
	if err := db.Select(&res, db.Rebind(q), id, oldPassword, StatusCreated); err != nil {
		fmt.Println(err)
		return ErrServerInternal
	}
	if len(res) == 0 {
		return ErrResourceNotFound
	}

	q = `
		UPDATE users
		SET
			password = ?
		WHERE
			id = ?
			AND status = ?
	`

	_, err = tx.Exec(tx.Rebind(q), newPassword, id, StatusCreated)
	if err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}

	if err := tx.Commit(); err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}

	log := Log{
		TableName:   "users",
		TableNameId: id,
		ColumnName:  "password",
		Action:      ActionChangeLogUpdate,
		Old:         oldPassword,
		New:         newPassword,
		CreatedBy:   id,
	}

	err = addLog(log)
	if err != nil {
		fmt.Println(err.Error())
	}

	return nil
}

func ResetPassword(id, newPassword string) error {
	tx, err := db.Beginx()
	if err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}
	defer tx.Rollback()

	q := `
		SELECT
			id
		FROM
			users
		WHERE
			id = ?
			AND status = ?
	`
	var res []string
	if err := db.Select(&res, db.Rebind(q), id, StatusCreated); err != nil {
		fmt.Println(err)
		return ErrServerInternal
	}
	if len(res) == 0 {
		return ErrResourceNotFound
	}

	q = `
		UPDATE users
		SET
			password = ?
		WHERE
			id = ?
			AND status = ?
	`

	_, err = tx.Exec(tx.Rebind(q), newPassword, id, StatusCreated)
	if err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}

	if err := tx.Commit(); err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}

	log := Log{
		TableName:   "users",
		TableNameId: id,
		ColumnName:  "password",
		Action:      ActionChangeLogUpdate,
		Old:         ValueChangeLogNone,
		New:         newPassword,
		CreatedBy:   id,
	}

	err = addLog(log)
	if err != nil {
		fmt.Println(err.Error())
	}

	return nil
}

func UpdateOtherUser(user User) error {
	tx, err := db.Beginx()
	if err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}
	defer tx.Rollback()

	q := `
		UPDATE users
		SET
			email = ?
			, phone = ?
			, updated_by = ?
			, updated_at = ?
		WHERE
			id = ?
			AND status = ?
	`

	_, err = tx.Exec(tx.Rebind(q), user.Email, user.Phone, user.ID, time.Now(), user.ID, StatusCreated)
	if err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}

	userDetail, err := FindUserDetail(user.ID)
	if err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}

	logs := []Log{}
	log := Log{
		TableName:   "users",
		TableNameId: user.ID,
		ColumnName:  "email",
		Action:      ActionChangeLogUpdate,
		Old:         userDetail.Email,
		New:         user.Email,
		CreatedBy:   user.ID,
	}
	logs = append(logs, log)

	log = Log{
		TableName:   "users",
		TableNameId: user.ID,
		ColumnName:  "phone",
		Action:      ActionChangeLogUpdate,
		Old:         userDetail.Phone,
		New:         user.Phone,
		CreatedBy:   user.ID,
	}
	logs = append(logs, log)

	q = `
		UPDATE user_roles
		SET
			status = ?
		WHERE
			user_id = ?
	`

	_, err = tx.Exec(tx.Rebind(q), StatusDeleted, user.ID)
	if err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}

	log = Log{
		TableName:   "user_roles",
		TableNameId: user.ID,
		ColumnName:  ColumnChangeLogDelete,
		Action:      ActionChangeLogUpdate,
		Old:         ValueChangeLogNone,
		New:         StatusDeleted,
		CreatedBy:   user.ID,
	}
	logs = append(logs, log)

	for _, v := range user.Role {
		q := `
				INSERT INTO user_roles(
					user_id
					, role_id
					, created_by
					, status
				)
				VALUES (?, ?, ?, ?)
			`

		_, err := tx.Exec(tx.Rebind(q), user.ID, v.Id, user.CreatedBy, StatusCreated)
		if err != nil {
			fmt.Println(err)
			return ErrServerInternal
		}

		log = Log{
			TableName:   "user_roles",
			TableNameId: user.ID,
			ColumnName:  ColumnChangeLogInsert,
			Action:      ActionChangeLogInsert,
			Old:         ValueChangeLogNone,
			New:         v.Id,
			CreatedBy:   user.ID,
		}
		logs = append(logs, log)
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

func UpdateUser(user User) error {
	tx, err := db.Beginx()
	if err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}
	defer tx.Rollback()

	q := `
			UPDATE users
			SET
				email = ?
				, phone = ?
				, updated_by = ?
				, updated_at = ?
			WHERE
				username = ?
				AND status = ?
		`

	_, err = tx.Exec(tx.Rebind(q), user.Email, user.Phone, user.ID, time.Now(), user.Username, StatusCreated)
	if err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}

	userDetail, err := FindUserDetail(user.ID)
	if err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}

	logs := []Log{}
	log := Log{
		TableName:   "users",
		TableNameId: user.ID,
		ColumnName:  "email",
		Action:      ActionChangeLogUpdate,
		Old:         userDetail.Email,
		New:         user.Email,
		CreatedBy:   user.ID,
	}
	logs = append(logs, log)

	log = Log{
		TableName:   "users",
		TableNameId: user.ID,
		ColumnName:  "phone",
		Action:      ActionChangeLogUpdate,
		Old:         userDetail.Phone,
		New:         user.Phone,
		CreatedBy:   user.ID,
	}
	logs = append(logs, log)

	q = `
			UPDATE user_roles
			SET
				status = ?
			WHERE
				user_id = ?
		`

	_, err = tx.Exec(tx.Rebind(q), StatusDeleted, user.ID)
	if err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}

	log = Log{
		TableName:   "user_roles",
		TableNameId: user.ID,
		ColumnName:  ColumnChangeLogDelete,
		Action:      ActionChangeLogUpdate,
		Old:         ValueChangeLogNone,
		New:         StatusDeleted,
		CreatedBy:   user.ID,
	}
	logs = append(logs, log)

	for _, v := range user.Role {
		q := `
			INSERT INTO user_roles(
				user_id
				, role_id
				, created_by
				, status
			)
			VALUES (?, ?, ?, ?)
		`

		_, err := tx.Exec(tx.Rebind(q), user.ID, v.Id, user.CreatedBy, StatusCreated)
		if err != nil {
			fmt.Println(err)
			return ErrServerInternal
		}

		log = Log{
			TableName:   "user_roles",
			TableNameId: user.ID,
			ColumnName:  ColumnChangeLogInsert,
			Action:      ActionChangeLogInsert,
			Old:         ValueChangeLogNone,
			New:         v.Id,
			CreatedBy:   user.ID,
		}
		logs = append(logs, log)
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

func BlockUser(userId, user string) error {
	tx, err := db.Beginx()
	if err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}
	defer tx.Rollback()

	q := `
		UPDATE users
		SET
			status = ?
			, updated_by = ?
			, updated_at = ?
		WHERE
			id = ?
	`

	_, err = tx.Exec(tx.Rebind(q), StatusDeleted, user, time.Now(), userId)
	if err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}

	if err := tx.Commit(); err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}

	logs := []Log{}
	log := Log{
		TableName:   "user_roles",
		TableNameId: user,
		ColumnName:  "status",
		Action:      ActionChangeLogUpdate,
		Old:         ValueChangeLogNone,
		New:         StatusDeleted,
		CreatedBy:   user,
	}
	logs = append(logs, log)

	log = Log{
		TableName:   "user_roles",
		TableNameId: user,
		ColumnName:  "updated_by",
		Action:      ActionChangeLogUpdate,
		Old:         ValueChangeLogNone,
		New:         user,
		CreatedBy:   user,
	}
	logs = append(logs, log)

	log = Log{
		TableName:   "user_roles",
		TableNameId: user,
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

func ReleaseUser(userId, user string) error {
	tx, err := db.Beginx()
	if err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}
	defer tx.Rollback()

	q := `
		UPDATE users
		SET
			status = ?
			, updated_by = ?
			, updated_at = ?
		WHERE
			id = ?
	`

	_, err = tx.Exec(tx.Rebind(q), StatusCreated, user, time.Now(), userId)
	if err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}

	if err := tx.Commit(); err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}

	logs := []Log{}
	log := Log{
		TableName:   "user_roles",
		TableNameId: user,
		ColumnName:  "status",
		Action:      ActionChangeLogUpdate,
		Old:         ValueChangeLogNone,
		New:         StatusCreated,
		CreatedBy:   user,
	}
	logs = append(logs, log)

	log = Log{
		TableName:   "user_roles",
		TableNameId: user,
		ColumnName:  "updated_by",
		Action:      ActionChangeLogUpdate,
		Old:         ValueChangeLogNone,
		New:         user,
		CreatedBy:   user,
	}
	logs = append(logs, log)

	log = Log{
		TableName:   "user_roles",
		TableNameId: user,
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

// Broadcast User

func InsertBroadcastUser(variantId, user string, target, description []string) error {
	tx, err := db.Beginx()
	if err != nil {
		fmt.Println(err)
		return ErrServerInternal
	}
	defer tx.Rollback()

	logs := []Log{}
	q := ""
	for i, v := range target {
		q = q + `
			INSERT INTO broadcast_users(
				program_id
				, target
				, description
				, state
				, created_by
				, status
			)
			VALUES ('` + variantId + `', '` + v + `', '` + description[i] + `', 'pending', '` + user + `', 'created')
			RETURNING
				id;
		`

		log := Log{
			TableName:   "broadcast_users",
			TableNameId: variantId,
			ColumnName:  ColumnChangeLogInsert,
			Action:      ActionChangeLogInsert,
			Old:         ValueChangeLogNone,
			New:         v,
			CreatedBy:   user,
		}
		logs = append(logs, log)
	}

	_, err = tx.Exec(tx.Rebind(q))
	if err != nil {
		fmt.Println(err)
		return ErrServerInternal
	}

	if err = tx.Commit(); err != nil {

		fmt.Println(err.Error())
		return ErrServerInternal
	}

	err = addLogs(logs)
	if err != nil {
		fmt.Println(err.Error())
	}

	return nil
}

// superadmin

func SuperAdminAddUser(u SuperAdminRegisterUser, user string) error {
	tx, err := db.Beginx()
	if err != nil {
		fmt.Println(err)
		return ErrServerInternal
	}
	defer tx.Rollback()

	username, err := CheckUsername(u.Username, u.AccountId)

	if username != "" {
		return ErrDuplicateEntry
	}

	q := `
		INSERT INTO users(
			username
			, password
			, email
			, phone
			, created_by
			, status
		)
		VALUES (?, ?, ?, ?, ?, ?)
		RETURNING
			id
	`

	var res []string
	if err := tx.Select(&res, tx.Rebind(q), u.Username, u.Password, u.Email, u.Phone, user, StatusCreated); err != nil {
		fmt.Println(err)
		return ErrServerInternal
	}

	logs := []Log{}
	log := Log{
		TableName:   "users",
		TableNameId: ValueChangeLogNone,
		ColumnName:  ColumnChangeLogInsert,
		Action:      ActionChangeLogInsert,
		Old:         ValueChangeLogNone,
		New:         res[0],
		CreatedBy:   user,
	}
	logs = append(logs, log)

	for _, v := range u.Role {
		q := `
			INSERT INTO user_roles(
				user_id
				, role_id
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

		log = Log{
			TableName:   "user_roles",
			TableNameId: res[0],
			ColumnName:  ColumnChangeLogInsert,
			Action:      ActionChangeLogInsert,
			Old:         ValueChangeLogNone,
			New:         v,
			CreatedBy:   user,
		}
		logs = append(logs, log)
	}

	q2 := `
		INSERT INTO user_accounts(
			user_id
			, account_id
			, created_by
			, status
		)
		VALUES (?, ?, ?, ?)
	`

	_, err = tx.Exec(tx.Rebind(q2), res[0], u.AccountId, user, StatusCreated)
	if err != nil {
		fmt.Println(err)
		return ErrServerInternal
	}

	log = Log{
		TableName:   "user_accounts",
		TableNameId: res[0],
		ColumnName:  ColumnChangeLogInsert,
		Action:      ActionChangeLogInsert,
		Old:         ValueChangeLogNone,
		New:         u.AccountId,
		CreatedBy:   user,
	}
	logs = append(logs, log)

	if err := tx.Commit(); err != nil {
		fmt.Println(err)
		return ErrServerInternal
	}

	err = addLogs(logs)
	if err != nil {
		fmt.Println(err.Error())
	}

	return nil
}

func SuperAdminFindAllUsers() ([]User, error) {
	q := `
		SELECT
			u.id
			, u.username
			, u.email
			, u.phone
			, u.created_at
			, u.created_by
			, u.status
		FROM
			users as u
		JOIN
			user_accounts as ua
		ON
			u.id = ua.user_id
		WHERE
			NOT u.username = 'suadmin'
	`
	var res []User
	if err := db.Select(&res, db.Rebind(q)); err != nil {
		fmt.Println("User : " + err.Error())
		return []User{}, ErrServerInternal
	}

	for i, v := range res {
		q1 := `
		SELECT
		  	u.id id,
		  	r.detail detail
		FROM
		  	users u
		JOIN
			user_roles ur
		ON
			u.id = ur.user_id
		JOIN
			roles r
	    	ON
	    		ur.role_id = r.id
		WHERE
		  	u.id = ?
		`

		var role []Role
		if err := db.Select(&role, db.Rebind(q1), v.ID); err != nil {
			fmt.Println("Role : " + err.Error())
			return []User{}, ErrServerInternal
		}
		res[i].Role = role
		account, err := GetAccountDetailByUser(v.ID)
		if err != nil {
			fmt.Println("Account : " + err.Error())
			return []User{}, ErrServerInternal
		}
		res[i].Account = account
	}

	return res, nil
}
