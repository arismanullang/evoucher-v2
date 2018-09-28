package model

type (
	TestReport struct {
		Id      string `db:"id"`
		Name    string `db:"name"`
		Total   string `db:"total"`
		Creator string `db:"creator"`
		State   string `db:"state"`
	}
	ReportProgram struct {
		Month       string `db:"month" json:"month"`
		MonthNumber string `db:"month_number" json:"month_number"`
		Total       string `db:"total" json:"total"`
		Username    string `db:"username" json:"username"`
		Creator     string `db:"created_by" json:"created_by"`
	}
	ReportVoucherByUser struct {
		Id      string `db:"id" json:"id"`
		Name    string `db:"name" json:"program_name"`
		Total   string `db:"total" json:"total"`
		Quota   string `db:"quota" json:"quota"`
		Creator string `db:"creator" json:"created_by"`
	}
	CompleteReportVoucherByUser struct {
		Id      string `db:"id" json:"id"`
		Name    string `db:"name" json:"program_name"`
		Total   string `db:"total" json:"total"`
		Quota   int    `db:"quota" json:"quota"`
		Creator string `db:"creator" json:"created_by"`
		State   string `db:"state" json:"state"`
	}
)

func MakeReport(id string) ([]TestReport, error) {
	q := `
		select va.id as id,
		va.program_name as name,
		count(vo.voucher_code) as total,
		va.created_by as creator,
		vo.state
		from programs va
		join vouchers vo
		on va.id = vo.program_id
		where va.status = ?
		and va.program_type = 'on-demand'
		and va.created_by = ?
		group by 1, 2, 4, 5
	`

	var resv []TestReport
	if err := db.Select(&resv, db.Rebind(q), StatusCreated, id); err != nil {
		return []TestReport{}, ErrServerInternal
	}
	if len(resv) < 1 {
		return []TestReport{}, ErrResourceNotFound
	}

	return resv, nil
}

func MakeReportProgram() ([]ReportProgram, error) {
	q := `
		select to_char(v.start_date,'Mon') as month,
		EXTRACT(MONTH FROM v.start_date) as month_number,
		count(v.program_name) as total,
		u.username,
		v.created_by
		from programs as v
		join users as u
		on u.id = v.created_by
		where v.status = ?
		group by 1,2,4,5
		order by u.username, month_number;
	`

	var resv []ReportProgram
	if err := db.Select(&resv, db.Rebind(q), StatusCreated); err != nil {
		return []ReportProgram{}, ErrServerInternal
	}
	if len(resv) < 1 {
		return []ReportProgram{}, ErrResourceNotFound
	}

	return resv, nil
}

func MakeCompleteReportVoucherByUser(id string) ([]CompleteReportVoucherByUser, error) {
	q := `
		select va.id as id,
		va.program_name as name,
		count(vo.voucher_code) as total,
		va.created_by as creator,
		vo.state,
		CAST ((va.max_quantity_voucher - (select count(id) from vouchers where program_id = va.id)) AS INTEGER) as quota
		from programs va
		join vouchers vo
		on va.id = vo.program_id
		join users u
		on u.id = va.created_by
		where va.status = ?
		and va.program_type = 'on-demand'
		and u.username = ?
		group by 1, 2, 4, 5
		order by vo.state
	`

	var resv []CompleteReportVoucherByUser
	if err := db.Select(&resv, db.Rebind(q), StatusCreated, id); err != nil {
		return []CompleteReportVoucherByUser{}, ErrServerInternal
	}
	if len(resv) < 1 {
		return []CompleteReportVoucherByUser{}, ErrResourceNotFound
	}

	return resv, nil
}

func MakeReportVoucherByUser(id string) ([]ReportVoucherByUser, error) {
	q := `
		select va.id as id,
		va.program_name as name,
		count(vo.voucher_code) as total,
		va.created_by as creator,
		CAST ((va.max_quantity_voucher - (select count(id) from vouchers where program_id = va.id)) AS INTEGER) as quota
		from programs va
		join vouchers vo
		on va.id = vo.program_id
		join users u
		on u.id = va.created_by
		where va.status = ?
		and va.program_type = 'on-demand'
		and u.username = ?
		group by 1, 2, 4
	`

	var resv []ReportVoucherByUser
	if err := db.Select(&resv, db.Rebind(q), StatusCreated, id); err != nil {
		return []ReportVoucherByUser{}, ErrServerInternal
	}
	if len(resv) < 1 {
		return []ReportVoucherByUser{}, ErrResourceNotFound
	}

	return resv, nil
}
