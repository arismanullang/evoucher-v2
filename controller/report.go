package controller

import (
	"github.com/go-zoo/bone"
	"net/http"
	"strings"

	"github.com/gilkor/evoucher-v2/model"
	u "github.com/gilkor/evoucher-v2/util"
	"github.com/gorilla/schema"
)

type (
	ReportFilter struct {
		Date      string `schema:"date" filter:"date"`
		CompanyID string `schema:"company_id" filter:"array"`
	}
)

//GetReportDailyVoucherTransaction :
func GetReportDailyVoucherTransaction(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()
	qp := u.NewQueryParam(r)

	//qp.SetCompanyID(bone.GetValue(r, "company"))
	//date := bone.GetValue(r, "date")
	//dates := strings.Split(date, ",")
	//date_from := bone.GetValue(r, "date_from")

	var decoder = schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)

	var f ReportFilter
	if err := decoder.Decode(&f, r.Form); err != nil {
		res.SetError(JSONErrFatal.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	qp.SetFilterModel(f)

	//temporary o
	if len(f.Date) <= 0 {
		res.SetError(JSONErrFatal.SetArgs("Date not set"))
		res.JSON(w, res, JSONErrBadRequest.Status)
		return
	}
	dates := strings.Split(f.Date, ",")
	data, next, err := model.GetReportDailyVoucherTransaction(dates[0], dates[1], qp)
	if err != nil {
		res.SetError(JSONErrFatal.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	res.SetResponse(data)
	if len(data) > 0 {
		res.SetNewPagination(r, qp.Page, next, data[0].Count)
	}
	res.JSON(w, res, http.StatusOK)
}

//GetReportDailyVoucherTransactionWithOutlet :
func GetReportDailyVoucherTransactionWithOutlet(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()
	qp := u.NewQueryParam(r)

	//qp.SetCompanyID(bone.GetValue(r, "company"))
	//date := bone.GetValue(r, "date")
	//dates := strings.Split(date, ",")
	//date_from := bone.GetValue(r, "date_from")

	var decoder = schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)

	var f ReportFilter
	if err := decoder.Decode(&f, r.Form); err != nil {
		res.SetError(JSONErrFatal.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	qp.SetFilterModel(f)

	//temporary o
	if len(f.Date) <= 0 {
		res.SetError(JSONErrFatal.SetArgs("Date not set"))
		res.JSON(w, res, JSONErrBadRequest.Status)
		return
	}
	dates := strings.Split(f.Date, ",")
	data, next, err := model.GetReportDailyVoucherTransactionWithOutlet(dates[0], dates[1], qp)
	if err != nil {
		res.SetError(JSONErrFatal.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	res.SetResponse(data)
	if len(data) > 0 {
		res.SetNewPagination(r, qp.Page, next, data[0].Count)
	}
	res.JSON(w, res, http.StatusOK)
}

//GetReportDailyOutletTransaction :
func GetReportDailyOutletTransaction(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()
	qp := u.NewQueryParam(r)

	//qp.SetCompanyID(bone.GetValue(r, "company"))
	//date := bone.GetValue(r, "date")
	//dates := strings.Split(date, ",")
	//date_from := bone.GetValue(r, "date_from")

	var decoder = schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)

	var f ReportFilter
	if err := decoder.Decode(&f, r.Form); err != nil {
		res.SetError(JSONErrFatal.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	qp.SetFilterModel(f)

	//temporary o
	if len(f.Date) <= 0 {
		res.SetError(JSONErrFatal.SetArgs("Date not set"))
		res.JSON(w, res, JSONErrBadRequest.Status)
		return
	}
	dates := strings.Split(f.Date, ",")
	data, next, err := model.GetReportDailyOutletTransaction(dates[0], dates[1], qp)
	if err != nil {
		res.SetError(JSONErrFatal.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	res.SetResponse(data)
	if len(data) > 0 {
		res.SetNewPagination(r, qp.Page, next, data[0].Count)
	}
	res.JSON(w, res, http.StatusOK)
}

//GetReportDailyOutletTransactionById
func GetReportDailyOutletTransactionById(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()
	qp := u.NewQueryParam(r)

	//qp.SetCompanyID(bone.GetValue(r, "company"))
	//date := bone.GetValue(r, "date")
	//dates := strings.Split(date, ",")
	//date_from := bone.GetValue(r, "date_from")
	id := bone.GetValue(r, "id")

	var decoder = schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)

	var f ReportFilter
	if err := decoder.Decode(&f, r.Form); err != nil {
		res.SetError(JSONErrFatal.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	qp.SetFilterModel(f)

	//temporary o
	if len(f.Date) <= 0 {
		res.SetError(JSONErrFatal.SetArgs("Date not set"))
		res.JSON(w, res, JSONErrBadRequest.Status)
		return
	}
	dates := strings.Split(f.Date, ",")
	data, next, err := model.GetReportDailyOutletTransactionById(id, dates[0], dates[1], qp)
	if err != nil {
		res.SetError(JSONErrFatal.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	res.SetResponse(data)
	if len(data) > 0 {
		res.SetNewPagination(r, qp.Page, next, data[0].Count)
	}
	res.JSON(w, res, http.StatusOK)
}

//GetReportMonthlyOutletTransaction :
func GetReportMonthlyOutletTransaction(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()
	qp := u.NewQueryParam(r)

	//qp.SetCompanyID(bone.GetValue(r, "company"))
	//date := bone.GetValue(r, "date")
	//dates := strings.Split(date, ",")
	//date_from := bone.GetValue(r, "date_from")

	var decoder = schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)

	var f ReportFilter
	if err := decoder.Decode(&f, r.Form); err != nil {
		res.SetError(JSONErrFatal.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	qp.SetFilterModel(f)

	//temporary o
	if len(f.Date) <= 0 {
		res.SetError(JSONErrFatal.SetArgs("Date not set"))
		res.JSON(w, res, JSONErrBadRequest.Status)
		return
	}
	dates := strings.Split(f.Date, ",")
	data, next, err := model.GetReportMonthlyOutletTransaction(dates[0], dates[1], qp)
	if err != nil {
		res.SetError(JSONErrFatal.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	res.SetResponse(data)
	if len(data) > 0 {
		res.SetNewPagination(r, qp.Page, next, data[0].Count)
	}
	res.JSON(w, res, http.StatusOK)
}

//GetReportYearlyOutletTransaction :
func GetReportYearlyOutletTransaction(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()
	qp := u.NewQueryParam(r)

	//qp.SetCompanyID(bone.GetValue(r, "company"))
	//date := bone.GetValue(r, "date")
	//dates := strings.Split(date, ",")
	//date_from := bone.GetValue(r, "date_from")

	var decoder = schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)

	var f ReportFilter
	if err := decoder.Decode(&f, r.Form); err != nil {
		res.SetError(JSONErrFatal.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	qp.SetFilterModel(f)

	//temporary o
	if len(f.Date) <= 0 {
		res.SetError(JSONErrFatal.SetArgs("Date not set"))
		res.JSON(w, res, JSONErrBadRequest.Status)
		return
	}
	dates := strings.Split(f.Date, ",")
	data, next, err := model.GetReportYearlyOutletTransaction(dates[0], dates[1], qp)
	if err != nil {
		res.SetError(JSONErrFatal.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	res.SetResponse(data)
	if len(data) > 0 {
		res.SetNewPagination(r, qp.Page, next, data[0].Count)
	}
	res.JSON(w, res, http.StatusOK)
}

//GetReportProgramTransactionDaily :
func GetReportProgramTransactionDaily(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()
	qp := u.NewQueryParam(r)

	//qp.SetCompanyID(bone.GetValue(r, "company"))
	//date := bone.GetValue(r, "date")
	//dates := strings.Split(date, ",")
	//date_from := bone.GetValue(r, "date_from")

	var decoder = schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)

	var f ReportFilter
	if err := decoder.Decode(&f, r.Form); err != nil {
		res.SetError(JSONErrFatal.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	qp.SetFilterModel(f)

	//temporary o
	if len(f.Date) <= 0 {
		res.SetError(JSONErrFatal.SetArgs("Date not set"))
		res.JSON(w, res, JSONErrBadRequest.Status)
		return
	}
	dates := strings.Split(f.Date, ",")
	data, next, err := model.GetReportProgramTransactionDaily(dates[0], dates[1], qp)
	if err != nil {
		res.SetError(JSONErrFatal.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	res.SetResponse(data)
	if len(data) > 0 {
		res.SetNewPagination(r, qp.Page, next, data[0].Count)
	}
	res.JSON(w, res, http.StatusOK)
}

//GetReportProgramIdTransactionDaily : Filter by program ID
func GetReportProgramIdTransactionDaily(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()
	qp := u.NewQueryParam(r)

	//qp.SetCompanyID(bone.GetValue(r, "company"))
	//date := bone.GetValue(r, "date")
	//dates := strings.Split(date, ",")
	//date_from := bone.GetValue(r, "date_from")
	id := bone.GetValue(r, "id")

	var decoder = schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)

	var f ReportFilter
	if err := decoder.Decode(&f, r.Form); err != nil {
		res.SetError(JSONErrFatal.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	qp.SetFilterModel(f)

	//temporary o
	if len(f.Date) <= 0 {
		res.SetError(JSONErrFatal.SetArgs("Date not set"))
		res.JSON(w, res, JSONErrBadRequest.Status)
		return
	}
	dates := strings.Split(f.Date, ",")
	data, next, err := model.GetReportProgramTransactionDailyById(id, dates[0], dates[1], qp)
	if err != nil {
		res.SetError(JSONErrFatal.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	res.SetResponse(data)
	if len(data) > 0 {
		res.SetNewPagination(r, qp.Page, next, data[0].Count)
	}
	res.JSON(w, res, http.StatusOK)
}

//GetReportProgramTransaction :
func GetReportProgramTransaction(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()
	qp := u.NewQueryParam(r)

	//qp.SetCompanyID(bone.GetValue(r, "company"))
	//date := bone.GetValue(r, "date")
	//dates := strings.Split(date, ",")
	//date_from := bone.GetValue(r, "date_from")

	var decoder = schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)

	var f ReportFilter
	if err := decoder.Decode(&f, r.Form); err != nil {
		res.SetError(JSONErrFatal.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	qp.SetFilterModel(f)

	//temporary o
	if len(f.Date) <= 0 {
		res.SetError(JSONErrFatal.SetArgs("Date not set"))
		res.JSON(w, res, JSONErrBadRequest.Status)
		return
	}
	dates := strings.Split(f.Date, ",")
	data, next, err := model.GetReportProgramTransaction(dates[0], dates[1], qp)
	if err != nil {
		res.SetError(JSONErrFatal.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	res.SetResponse(data)
	if len(data) > 0 {
		res.SetNewPagination(r, qp.Page, next, data[0].Count)
	}
	res.JSON(w, res, http.StatusOK)
}
