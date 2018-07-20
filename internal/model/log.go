package model

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
	return nil
}

func addLogs(logs []Log) error {
	return nil
}
