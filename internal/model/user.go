package model

import "fmt"

type (
	User struct {
		AccountId string   `db:"account_id"`
		Username  string   `db:"username"`
		Password  string   `db:"password"`
		Email     string   `db:"email"`
		Phone     string   `db:"phone"`
		RoleId    []string `db:"-"`
		CreatedBy string   `db:"created_by"`
	}
	UserRes struct {
		Id       string `db:"id"`
		Username string `db:"username"`
	}
	Role struct {
		Id         string `db:"id"`
		RoleDetail string `db:"role_detail"`
	}
)

func AddUser(u User) error {
	fmt.Println("Add")
	tx, err := db.Beginx()
	if err != nil {
		fmt.Println(err)
		return ErrServerInternal
	}
	defer tx.Rollback()

	username, err := CheckUsername(u.Username)

	if username != "" {
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

		for _, v := range u.RoleId {
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

		_, err := tx.Exec(tx.Rebind(q2), res[0], u.AccountId, u.CreatedBy, StatusCreated)
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
			username
			, email
			, phone
		FROM
			users
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
			roles.role_detail
		FROM
			users
		JOIN
			user_roles
		ON
			users.id = user_roles.user_id
		JOIN
			roles
		ON
			user_roles.role_id = roles.id
		WHERE
			users.id = ?
			AND users.status = ?
	`
	var role []string
	if err := db.Select(&role, db.Rebind(q1), userId, StatusCreated); err != nil {
		fmt.Println(err)
		return User{}, ErrServerInternal
	}
	res[0].RoleId = role

	return res[0], nil
}

func Login(username, password, accountId string) (string, error) {
	fmt.Println("Login")
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
			AND u.password = ?
			AND ua.account_id = ?
			AND u.status = ?
	`
	var res []string
	if err := db.Select(&res, db.Rebind(q), username, password, accountId, StatusCreated); err != nil {
		fmt.Println(err)
		return "", ErrServerInternal
	}
	if len(res) == 0 {
		return "", ErrResourceNotFound
	}
	return res[0], nil
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
