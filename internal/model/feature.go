package model

import ()

var (
	UiFeatures  map[string][]string
	ApiFeatures map[string][]string
)

type (
	Feature struct {
		Id       string `db:"id"`
		Role     string `db:"role"`
		Category string `db:"category"`
		Detail   string `db:"detail"`
		Type     string `db:"type"`
	}
)

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
