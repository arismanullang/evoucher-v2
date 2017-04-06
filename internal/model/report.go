package model

import "fmt"

type (
	TestReport struct {
		Id      string `db:"id"`
		Name    string `db:"name"`
		Total   string `db:"total"`
		Creator string `db:"creator"`
	}
	ReportVariant struct {
		Month       string `db:"month"`
		MonthNumber string `db:"month_number"`
		Total       string `db:"total"`
		Creator     string `db:"created_by"`
	}
)

func MakeReport(id string) ([]TestReport, error) {
	fmt.Println("Select Variant")
	q := `
		select va.id as id,
			va.variant_name as name,
			count(vo.voucher_code) as total,
			va.created_by as creator
		from variants va
		left join vouchers vo
		on va.id = vo.variant_id
		where va.status = ?
		and va.variant_type = 'on-demand'
		and va.created_by = ?
		group by 1, 2, 4
	`

	var resv []TestReport
	if err := db.Select(&resv, db.Rebind(q), StatusCreated, id); err != nil {
		fmt.Println(err.Error())
		return []TestReport{}, ErrServerInternal
	}
	if len(resv) < 1 {
		return []TestReport{}, ErrResourceNotFound
	}

	return resv, nil
}

func MakeReportVariant() ([]ReportVariant, error) {
	fmt.Println("Select Variant")
	q := `
		select to_char(created_at,'Mon') as month,
			EXTRACT(MONTH FROM created_at) as month_number,
			count(variant_name) as total,
			created_by
		from variants
		where status = ?
		group by 1,2,4
		order by created_by, month_number;
	`

	var resv []ReportVariant
	if err := db.Select(&resv, db.Rebind(q), StatusCreated); err != nil {
		fmt.Println(err.Error())
		return []ReportVariant{}, ErrServerInternal
	}
	if len(resv) < 1 {
		return []ReportVariant{}, ErrResourceNotFound
	}

	return resv, nil
}
