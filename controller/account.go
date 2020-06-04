package controller

import (
	"net/http"

	"github.com/gilkor/evoucher-v2/model"
	u "github.com/gilkor/evoucher-v2/util"
	"github.com/go-zoo/bone"
	"github.com/gorilla/schema"
)

type AccountFilter struct {
	ID                string `schema:"id" filter:"array"`
	Name              string `schema:"name" filter:"string"`
	CompanyId         string `schema:"company_id" filter:"string"`
	Gender            string `schema:"gender" filter:"enum"`
	Email             string `schema:"email" filter:"string"`
	MobileCallingCode string `schema:"mobile_calling_code" filter:"string"`
	MobileNo          string `schema:"mobile_no" filter:"string"`
	State             string `schema:"state" filter:"enum"`
	Status            string `schema:"status" filter:"enum"`
	CreatedBy         string `schema:"created_by" filter:"string"`
	CreatedAt         string `schema:"created_at" filter:"date"`
	UpdatedBy         string `schema:"updated_by" filter:"string"`
	UpdatedAt         string `schema:"updated_at" filter:"date"`
}

//GetAccounts : GET list of accounts
func GetAccounts(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()
	qp := u.NewQueryParam(r)

	var decoder = schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)

	var f AccountFilter
	if err := decoder.Decode(&f, r.Form); err != nil {
		res.SetError(JSONErrFatal.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	qp.SetFilterModel(f)

	accounts, next, err := model.GetAccounts(qp)
	if err != nil {
		res.SetError(JSONErrFatal.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	res.SetResponse(accounts)
	res.SetNewPagination(r, qp.Page, next, (*accounts)[0].Count)
	// res.SetCount((*accounts)[0].Count)
	res.JSON(w, res, http.StatusOK)
}

//GetAccountByID : GET
func GetAccountByID(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	qp := u.NewQueryParam(r)
	id := bone.GetValue(r, "id")
	account, _, err := model.GetAccountsByID(qp, id)
	if err != nil {
		res.SetError(JSONErrResourceNotFound)
		res.JSON(w, res, JSONErrResourceNotFound.Status)
		return
	}

	res.SetResponse(account)
	res.JSON(w, res, http.StatusOK)
}
