package model

import "fmt"

type (
	RegisterUser struct {
		ID        string   `db:"id"`
		AccountID string   `db:"account_id"`
		Username  string   `db:"username"`
		Password  string   `db:"password"`
		Email     string   `db:"email"`
		Phone     string   `db:"phone"`
		Role      []string `db:"-"`
		CreatedBy string   `db:"created_by"`
	}
	User struct {
		ID        string `db:"id"`
		AccountID string `db:"account_id"`
		Username  string `db:"username"`
		Password  string `db:"password"`
		Email     string `db:"email"`
		Phone     string `db:"phone"`
		Role      []Role `db:"-"`
		CreatedBy string `db:"created_by"`
		CreatedAt string `db:"created_at"`
	}
	Role struct {
		Id         string `db:"id"`
		RoleDetail string `db:"role_detail"`
	}

	UserRes struct {
		Id       string `db:"id"`
		Username string `db:"username"`
	}
)

func AddUser(u RegisterUser) error {
	fmt.Println("Add")
	tx, err := db.Beginx()
	if err != nil {
		fmt.Println(err)
		return ErrServerInternal
	}
	defer tx.Rollback()

	username, err := CheckUsername(u.Username)

	if username == "" {
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
		if err := tx.Select(&res, tx.Rebind(q), u.Username, u.Password, u.Email, u.Phone, u.CreatedBy, StatusCreated); err != nil {
			fmt.Println(err)
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

			_, err := tx.Exec(tx.Rebind(q), res[0], v, u.CreatedBy, StatusCreated)
			if err != nil {
				fmt.Println(err)
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

		_, err := tx.Exec(tx.Rebind(q2), res[0], u.AccountID, u.CreatedBy, StatusCreated)
		if err != nil {
			fmt.Println(err)
			return ErrServerInternal
		}

		if err := tx.Commit(); err != nil {
			fmt.Println(err)
			return ErrServerInternal
		}
		return nil
	}

	return ErrDuplicateEntry
}

func CheckUsername(username string) (string, error) {
	q := `
		SELECT id FROM users
		WHERE
			username = ?
			AND status = ?
	`
	var res []string
	if err := db.Select(&res, db.Rebind(q), username, StatusCreated); err != nil {
		fmt.Println(err)
		return "", ErrServerInternal
	}
	if len(res) == 0 {
		return "", nil
	}
	return res[0], nil
}

func FindAllUsers(accountId string) ([]UserRes, error) {
	fmt.Println("Select User " + accountId)
	q := `
		SELECT DISTINCT u.id, u.username FROM users as u
		JOIN user_accounts as ua ON u.id = ua.user_id
		JOIN user_roles as ur ON u.id = ur.user_id
		WHERE ua.account_id = ?
		AND u.status = ?
	`

	var resv []UserRes
	if err := db.Select(&resv, db.Rebind(q), accountId, StatusCreated); err != nil {
		fmt.Println(err)
		return []UserRes{}, ErrServerInternal
	}
	if len(resv) < 1 {
		return []UserRes{}, ErrResourceNotFound
	}

	return resv, nil
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
			AND status = ?
	`
	var res []User
	if err := db.Select(&res, db.Rebind(q), userId, StatusCreated); err != nil {
		fmt.Println(err)
		return User{}, ErrServerInternal
	}

	q1 := `
		SELECT
		  	u.id id,
		  	r.role_detail role_detail
		FROM
		  	users u
		JOIN
			user_roles ur
		ON  	u.id = ur.user_id
		JOIN
			roles r
		    ON  ur.role_id = r.id
		WHERE
		  	u.id = ?
		  AND
		  	u.status = ?

	`
	var role []Role
	if err := db.Select(&role, db.Rebind(q1), userId, StatusCreated); err != nil {
		fmt.Println(err)
		return User{}, ErrServerInternal
	}

	res[0].Role = make([]Role, len(role))
	for k, v := range role {
		res[0].Role[k].Id = v.Id
		res[0].Role[k].RoleDetail = v.RoleDetail
	}

	return res[0], nil
}

func Login(username, password string) (string, error) {
	fmt.Println("Login")

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
	return nil
}

func ChangePassword(id, oldPassword, newPassword string) error {
	fmt.Println("Change Password")
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

	tx, err := db.Beginx()
	if err != nil {
		fmt.Println(err.Error())
		return ErrServerInternal
	}
	defer tx.Rollback()

	q = `
		UPDATE users
		SET
			password = ?
		WHERE
			id = ?
			And password = ?
			AND status = ?
	`

	_, err = tx.Exec(tx.Rebind(q), newPassword, id, oldPassword, StatusCreated)
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
		WHERE
			username = ?
			AND status = ?
	`

	_, err = tx.Exec(tx.Rebind(q), user.Email, user.Phone, user.Username, StatusCreated)
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

// Role -----------------------------------------------------------------------------------------------

func FindAllRole() ([]Role, error) {
	fmt.Println("Select All Role")
	q := `
		SELECT id, role_detail
		FROM roles
		WHERE status = ?
	`

	var resv []Role
	if err := db.Select(&resv, db.Rebind(q), StatusCreated); err != nil {
		fmt.Println(err)
		return []Role{}, ErrServerInternal
	}
	if len(resv) < 1 {
		return []Role{}, ErrResourceNotFound
	}

	return resv, nil
}

// Broadcast User

func InsertBroadcastUser(variantId, user string, target, description []string) error {
	fmt.Println("Add")
	tx, err := db.Beginx()
	if err != nil {
		fmt.Println(err)
		return ErrServerInternal
	}
	defer tx.Rollback()
	q := ""
	for i, v := range target {
		q = q + `
			INSERT INTO broadcast_users(
				variant_id
				, broadcast_target
				, description
				, state
				, created_by
				, status
			)
			VALUES ('` + variantId + `', '` + v + `', '` + description[i] + `', 'pending', '` + user + `', 'created');
		`
	}
	fmt.Println(q)
	_, err = tx.Exec(tx.Rebind(q))
	if err != nil {
		fmt.Println(err)
		return ErrServerInternal
	}

	if err = tx.Commit(); err != nil {

		fmt.Println(err.Error())
		return ErrServerInternal
	}
	return nil
}
