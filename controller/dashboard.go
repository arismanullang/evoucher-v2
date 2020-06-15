package controller

import (
	"github.com/gilkor/evoucher-v2/model"
	u "github.com/gilkor/evoucher-v2/util"
	"github.com/gorilla/schema"
	"net/http"
	"strings"
)

type (
	DashboardFilter struct {
		Date      string `schema:"date" filter:"date"`
		CompanyID string `schema:"company_id" filter:"array"`
	}
)

//GetDashboardVoucherUsage :
func GetDashboardVoucherUsage(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()
	qp := u.NewQueryParam(r)

	//qp.SetCompanyID(bone.GetValue(r, "company"))
	//date := bone.GetValue(r, "date")
	//dates := strings.Split(date, ",")
	//date_from := bone.GetValue(r, "date_from")

	var decoder = schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)

	var f DashboardFilter
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
	data, _, err := model.GetDashboardVoucherUsage(dates[0], dates[1], qp)
	if err != nil {
		res.SetError(JSONErrFatal.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	res.SetResponse(data)
	//res.SetNewPagination(r, qp.Page, next, data[0].Count)
	res.JSON(w, res, http.StatusOK)
}

//GetDashboardTopProgram :
func GetDashboardTopProgram(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()
	qp := u.NewQueryParam(r)

	//qp.SetCompanyID(bone.GetValue(r, "company"))
	//date := bone.GetValue(r, "date")

	//date_from := bone.GetValue(r, "date_from")

	var decoder = schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)

	var f DashboardFilter
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
	data, _, err := model.GetDashboardTopProgram(dates[0], dates[1], qp)
	if err != nil {
		res.SetError(JSONErrFatal.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	res.SetResponse(data)
	//res.SetNewPagination(r, qp.Page, next, data[0].Count)
	res.JSON(w, res, http.StatusOK)
}

//GetDashboardTopOutlet :
func GetDashboardTopOutlet(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()
	qp := u.NewQueryParam(r)

	//qp.SetCompanyID(bone.GetValue(r, "company"))
	//date := bone.GetValue(r, "date")
	//dates := strings.Split(date, ",")
	//date_from := bone.GetValue(r, "date_from")

	var decoder = schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)

	var f DashboardFilter
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
	data, _, err := model.GetDashboardTopOutlet(dates[0], dates[1], qp)
	if err != nil {
		res.SetError(JSONErrFatal.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	res.SetResponse(data)
	//res.SetNewPagination(r, qp.Page, next, data[0].Count)
	res.JSON(w, res, http.StatusOK)
}
