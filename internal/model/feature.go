package model

var (
	UiFeatures  map[string][]string
	ApiFeatures map[string][]string
)

type (
	Feature struct {
		Id       string `db:"id" json:"id"`
		Role     string `db:"role" json:"role"`
		Category string `db:"category" json:"category"`
		Detail   string `db:"detail" json:"detail"`
		Type     string `db:"type" json:"type"`
	}
)

func GetAllFeatures() ([]Feature, error) {
	q := `
		SELECT
			f.id, f.type, f.category, f.detail
		FROM 	features AS f
		WHERE
			f.status = ?
			AND NOT f.category = 'sa'
	`

	var resv []Feature
	if err := db.Select(&resv, db.Rebind(q), StatusCreated); err != nil {
		return []Feature{}, ErrServerInternal
	}
	if len(resv) < 1 {
		return []Feature{}, ErrResourceNotFound
	}

	return resv, nil
}

func GetAllUiFeatures() ([]Feature, error) {
	q := `
		SELECT
			f.id, f.type, r.detail as role, f.category, f.detail
		FROM roles AS r
		JOIN role_features AS rf
		ON
			r.id = rf.role_id
		JOIN features AS f
		ON
			f.id = rf.feature_id
		WHERE
			f.type = 'ui'
			AND f.status = ?
	`

	var resv []Feature
	if err := db.Select(&resv, db.Rebind(q), StatusCreated); err != nil {
		return []Feature{}, ErrServerInternal
	}
	if len(resv) < 1 {
		return []Feature{}, ErrResourceNotFound
	}

	return resv, nil
}

func GetUiFeatures(roleId string) ([]Feature, error) {
	q := `
		SELECT
			f.id, f.type, r.detail as role, f.category, f.detail
		FROM roles AS r
		JOIN role_features AS rf
		ON
			r.id = rf.role_id
		JOIN features AS f
		ON
			f.id = rf.feature_id
		WHERE
			f.type = 'ui'
			AND rf.status = ?
			AND r.id = ?
	`

	var resv []Feature
	if err := db.Select(&resv, db.Rebind(q), StatusCreated, roleId); err != nil {
		return []Feature{}, ErrServerInternal
	}
	if len(resv) < 1 {
		return []Feature{}, ErrResourceNotFound
	}

	return resv, nil
}

func GetApiFeatures(roleId string) ([]Feature, error) {
	q := `
		SELECT
			f.id, f.type, r.detail as role, f.category, f.detail
		FROM roles AS r
		JOIN role_features AS rf
		ON
			r.id = rf.role_id
		JOIN features AS f
		ON
			f.id = rf.feature_id
		WHERE
			f.type = 'api'
			AND rf.status = ?
			AND rf.role_id = ?
	`

	var resv []Feature
	if err := db.Select(&resv, db.Rebind(q), StatusCreated, roleId); err != nil {
		return []Feature{}, ErrServerInternal
	}
	if len(resv) < 1 {
		return []Feature{}, ErrResourceNotFound
	}

	return resv, nil
}

func GetAllApiFeatures() ([]Feature, error) {
	q := `
		SELECT
			f.id, f.type, r.detail as role, f.category, f.detail
		FROM roles AS r
		JOIN role_features AS rf
		ON
			r.id = rf.role_id
		JOIN features AS f
		ON
			f.id = rf.feature_id
		WHERE
			f.type = 'api'
			AND f.status = ?
	`

	var resv []Feature
	if err := db.Select(&resv, db.Rebind(q), StatusCreated); err != nil {
		return []Feature{}, ErrServerInternal
	}
	if len(resv) < 1 {
		return []Feature{}, ErrResourceNotFound
	}

	return resv, nil
}
