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
	tx, err := db.Beginx()
	if err != nil {
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
		return ErrServerInternal
	}

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
			return ErrServerInternal
		}
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
		return ErrServerInternal
	}

	if err = tx.Commit(); err != nil {
		return ErrServerInternal
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
			return []User{}, ErrServerInternal
		}
		res[i].Role = role
		account, err := GetAccountDetailByAccountId(accountId)
		if err != nil {
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
		return User{}, ErrServerInternal
	}

	res[0].Role = role
	account, err := GetAccountDetailByUser(userId)
	if err != nil {
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
		return "", ErrServerInternal
	}
	if len(res) == 0 {
		return "", ErrResourceNotFound
	}

	return res[0], nil
}

func UpdatePassword(id, password string) error {
	tx, err := db.Beginx()
	if err != nil {
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
		return ErrServerInternal
	}

	if err := tx.Commit(); err != nil {
		return ErrServerInternal
	}

	return nil
}

func ChangePassword(id, oldPassword, newPassword string) error {
	tx, err := db.Beginx()
	if err != nil {
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
		return ErrServerInternal
	}

	if err := tx.Commit(); err != nil {
		return ErrServerInternal
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

	_, err = FindUserDetail(user.ID)
	if err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}

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
			return ErrServerInternal
		}
	}

	if err := tx.Commit(); err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
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

	_, err = FindUserDetail(user.ID)
	if err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}

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
			return ErrServerInternal
		}
	}

	if err := tx.Commit(); err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
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

	return nil
}

// Broadcast User

func InsertBroadcastUser(variantId, user string, target, description []string) error {
	tx, err := db.Beginx()
	if err != nil {
		return ErrServerInternal
	}
	defer tx.Rollback()

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
	}

	_, err = tx.Exec(tx.Rebind(q))
	if err != nil {
		return ErrServerInternal
	}

	if err = tx.Commit(); err != nil {
		return ErrServerInternal
	}

	return nil
}

// superadmin

func SuperAdminAddUser(u SuperAdminRegisterUser, user string) error {
	tx, err := db.Beginx()
	if err != nil {
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
		return ErrServerInternal
	}

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
			return ErrServerInternal
		}
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
		return ErrServerInternal
	}

	if err := tx.Commit(); err != nil {
		return ErrServerInternal
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
		return []User{}, ErrServerInternal
	}

	for i, v := range res {
		q1 := `
			SELECT
				u.id id,
				r.detail detail
			FROM
				users u
				JOIN user_roles ur ON u.id = ur.user_id
				JOIN roles r ON ur.role_id = r.id
			WHERE
				u.id = ?
		`
		var role []Role
		if err := db.Select(&role, db.Rebind(q1), v.ID); err != nil {
			return []User{}, ErrServerInternal
		}
		res[i].Role = role
		account, err := GetAccountDetailByUser(v.ID)
		if err != nil {
			return []User{}, ErrServerInternal
		}
		res[i].Account = account
	}

	return res, nil
}
