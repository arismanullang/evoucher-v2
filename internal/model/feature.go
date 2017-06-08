package model

import ()

var (
	UiFeatures  map[string][]string
	ApiFeatures map[string][]string
)

type (
	Feature struct {
		Role     string `db:"role_detail"`
		Category string `db:"feature_category"`
		Detail   string `db:"feature_detail"`
	}
)

func GetAllUiFeatures() ([]Feature, error) {
	q := `
		SELECT
			r.role_detail, f.feature_category, f.feature_detail
		FROM roles AS r
		JOIN role_features AS rf
		ON
			r.id = rf.role_id
		JOIN features AS f
		ON
			f.id = rf.feature_id
	`

	var resv []Feature
	if err := db.Select(&resv, db.Rebind(q)); err != nil {
		return []Feature{}, ErrServerInternal
	}
	if len(resv) < 1 {
		return []Feature{}, ErrResourceNotFound
	}

	return resv, nil
}
