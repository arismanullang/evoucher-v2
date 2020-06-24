package model

import (
	"time"

	"github.com/gilkor/evoucher-v2/util"
)

type (
	// ApprovalLog : Struct for approval log
	ApprovalLog struct {
		ID             string    `json:"id,omitempty" db:"id"`
		ObjectID       string    `json:"object_id,omitempty" db:"object_id"`
		ObjectCategory string    `json:"object_category,omitempty" db:"object_category"`
		CreatedBy      string    `json:"created_by,omitempty" db:"created_by"`
		CreatedAt      time.Time `json:"created_at,omitempty" db:"created_at"`
		Status         string    `json:"status,omitempty" db:"status"`
	}
)

// Insert : insert approval log
func (a *ApprovalLog) Insert() error {

	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	q := `INSERT INTO
			approval_log 
			(
				object_id
				, object_category
				, created_by
				, status
			) 
			VALUES (?, ?, ?, ?)
			RETURNING
					id
					, object_id
					, object_category
					, created_at
					, created_by 
					, status
	`
	util.DEBUG("test", q)
	var res []ApprovalLog
	err = tx.Select(&res, tx.Rebind(q), a.ObjectID, a.ObjectCategory, a.CreatedBy, a.Status)
	if err != nil {
		util.DEBUG("err = ", err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil

}
