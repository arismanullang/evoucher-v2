package model

import "time"

//"database/sql"
//"fmt"
//"time"

type (
	Log struct {
		Id          string `db:"id" json:"id"`
		TableName   string `db:"table_name" json:"table_name"`
		TableNameId string `db:"table_name_id" json:"table_name_id"`
		ColumnName  string `db:"column_name" json:"column_name"`
		Action      string `db:"action" json:"action"`
		Old         string `db:"old" json:"old"`
		New         string `db:"new" json:"new"`
		CreatedBy   string `db:"created_by" json:"created_by"`
		CreatedAt   string `db:"created_at" json:"created_at"`
		Status      string `db:"status" json:"status"`
	}
)

func addLog(a Log) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	q := `
		INSERT INTO changes_log(
			table_name
			, table_name_id
			, column_name
			, action
			, old
			, new
			, created_by
			, created_at
			, status
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING
			id
	`

	var res []string
	if err := tx.Select(&res, tx.Rebind(q), a.TableName, a.TableNameId, a.ColumnName, a.Action, a.Old, a.New, a.CreatedBy, time.Now(), StatusCreated); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func addLogs(logs []Log) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, a := range logs {
		q := `
			INSERT INTO changes_log(
				table_name
				, table_name_id
				, column_name
				, action
				, old
				, new
				, created_by
				, created_at
				, status
			)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
			RETURNING
				id
		`

		var res []string
		if err := tx.Select(&res, tx.Rebind(q), a.TableName, a.TableNameId, a.ColumnName, a.Action, a.Old, a.New, a.CreatedBy, time.Now(), StatusCreated); err != nil {
			return err
		}
	}
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
