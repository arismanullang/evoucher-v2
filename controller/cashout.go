package controller

import (
	"encoding/json"
	"github.com/gilkor/evoucher-v2/model"
	u "github.com/gilkor/evoucher-v2/util"
	"github.com/go-zoo/bone"
	"github.com/gorilla/schema"
	"net/http"
	"time"
)

type CashoutFilter struct {
	ID            string     `schema:"id" filter:"array"`
	AccountID     string     `schema:"account_id" filter:"string"`
	Code          string     `schema:"code" filter:"string"`
	PartnerID     string     `schema:"partner_id" filter:"string"`
	BankAccount   string     `schema:"bank_account" filter:"string"`
	Amount        float64    `schema:"amount" filter:"string"`
	PaymentMethod string     `schema:"payment_method" filter:"string"`
	CreatedAt     *time.Time `schema:"created_at" filter:"string"`
	CreatedBy     string     `schema:"created_by" filter:"string"`
	UpdatedAt     *time.Time `schema:"updated_at" filter:"date"`
	UpdatedBy     string     `schema:"updated_by" filter:"date"`
	Status        string     `schema:"status" filter:"enum"`
}

//GetCashouts : GET list of Cashouts
func GetCashouts(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	qp := u.NewQueryParam(r)

	var decoder = schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)

	var f CashoutFilter
	if err := decoder.Decode(&f, r.Form); err != nil {
		res.SetError(JSONErrFatal.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	qp.SetFilterModel(f)

	Cashouts, next, err := model.GetCashouts(qp)
	if err != nil {
		res.SetError(JSONErrFatal.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	res.SetResponse(Cashouts)
	res.SetNewPagination(r, qp.Page, next, (*Cashouts)[0].Count)
	res.JSON(w, res, http.StatusOK)
}

//GetCashoutSummary : GET list of Cashout Summary
func GetCashoutSummary(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	qp := u.NewQueryParam(r)

	cashoutSummary, next, err := model.GetCashoutSummary(qp)
	if err != nil {
		res.SetError(JSONErrFatal.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	res.SetResponse(r)
	res.SetNewPagination(r, qp.Page, next, cashoutSummary[0].Count)
	res.JSON(w, res, http.StatusOK)
}

//GetCashoutUsedVoucher : GET list of Cashout Summary
func GetCashoutUsedVoucher(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	qp := u.NewQueryParam(r)
	program_id := bone.GetValue(r, "id")
	//add QueryParam Filter for used voucher
	//qp.

	usedVouchers, next, err := model.GetVouchersByProgramID(program_id, qp)
	if err != nil {
		res.SetError(JSONErrFatal.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	res.SetResponse(r)
	res.SetNewPagination(r, qp.Page, next, usedVouchers[0].Count)
	res.JSON(w, res, http.StatusOK)
}

//GetCashoutByID : GET
func GetCashoutByID(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	qp := u.NewQueryParam(r)
	id := bone.GetValue(r, "id")
	Cashout, _, err := model.GetCashoutByID(id, qp)
	if err != nil {
		res.SetError(JSONErrResourceNotFound)
		res.JSON(w, res, JSONErrResourceNotFound.Status)
		return
	}

	res.SetResponse(Cashout)
	res.JSON(w, res, http.StatusOK)
}

//DeleteCashout : remove Cashout
func DeleteCashout(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	id := bone.GetValue(r, "id")
	p := model.Cashout{ID: id}
	if err := p.Delete(); err != nil {
		res.SetError(JSONErrResourceNotFound)
		res.JSON(w, res, JSONErrResourceNotFound.Status)
		return
	}
	res.JSON(w, res, http.StatusOK)
}

//PostCashoutByPartner : POST Cashout by partner
func PostCashoutByPartner(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	var req model.Cashout
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		u.DEBUG(err)
		res.SetErrorWithDetail(JSONErrFatal, err)
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}
	// reqCashout.ID = bone.GetValue(r, "holder")
	response, err := req.Insert()
	if err != nil {
		u.DEBUG(err)
		res.SetErrorWithDetail(JSONErrFatal, err)
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	res.SetResponse(response)
	res.JSON(w, res, http.StatusCreated)
}
