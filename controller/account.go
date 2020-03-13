package controller

import (
	"net/http"

	"github.com/gilkor/evoucher-v2/model"
	u "github.com/gilkor/evoucher-v2/util"
	"github.com/go-zoo/bone"
)

//GetAccounts : GET list of accounts
func GetAccounts(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	qp := u.NewQueryParam(r)

	accounts, next, err := model.GetAccounts(qp)
	if err != nil {
		res.SetError(JSONErrFatal.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	res.SetResponse(accounts)
	res.SetPagination(r, qp.Page, next)
	res.JSON(w, res, http.StatusOK)
}

//GetAccountByID : GET
func GetAccountByID(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	qp := u.NewQueryParam(r)
	id := bone.GetValue(r, "id")
	account, _, err := model.GetAccountByID(qp, id)
	if err != nil {
		res.SetError(JSONErrResourceNotFound)
		res.JSON(w, res, JSONErrResourceNotFound.Status)
		return
	}

	res.SetResponse(account)
	res.JSON(w, res, http.StatusOK)
}
