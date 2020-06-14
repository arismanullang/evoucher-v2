package controller

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gilkor/evoucher-v2/model"
	u "github.com/gilkor/evoucher-v2/util"
	"github.com/go-zoo/bone"
	"github.com/gorilla/schema"
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

type (
	CashoutSummaryFilter struct {
		Date        *time.Time `schema:"date" filter:"date"`
		PartnerID   string     `schema:"partner_id" filter:"array"`
		PartnerName string     `schema:"partner_name" filter:"string"`
	}
	CashoutUnpaidFilter struct {
		Date        *time.Time `schema:"date" filter:"date"`
		PartnerID   string     `schema:"partner_id" filter:"array"`
		PartnerName string     `schema:"partner_name" filter:"string"`
	}
)

//GetCashoutsUnpaid : GET list of Unpaid Cashouts
func GetCashoutsUnpaid(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	qp := u.NewQueryParam(r)

	var decoder = schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)

	var f CashoutUnpaidFilter
	if err := decoder.Decode(&f, r.Form); err != nil {
		res.SetError(JSONErrFatal.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	qp.SetFilterModel(f)

	cashout, next, err := model.GetCashoutUnpaid(qp)
	if err != nil {
		res.SetError(JSONErrFatal.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	res.SetResponse(cashout)
	res.SetNewPagination(r, qp.Page, next, cashout[0].Count)
	res.JSON(w, res, http.StatusOK)
}

//GetUnpaidReimburse : GET list of unpaid reimburse
func GetUnpaidReimburse(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()
	qp := u.NewQueryParam(r)

	startDate := r.FormValue("start_date")
	endDate := r.FormValue("end_date")

	qp.SetCompanyID(bone.GetValue(r, "company"))

	unpaidReimburse, next, err := model.GetUnpaidReimburse(qp, startDate, endDate)
	if err != nil && err != model.ErrorResourceNotFound {
		res.SetError(JSONErrFatal.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	list := []model.UnpaidReimburse{}
	list = append(list, unpaidReimburse...)
	res.SetResponse(list)
	if len(unpaidReimburse) > 0 {
		res.SetNewPagination(r, qp.Page, next, (unpaidReimburse)[0].Count)
	}
	res.JSON(w, res, http.StatusOK)
}

//GetUnpaidVouchersByOutlet : GET list of unpaid vouchers by outlet
func GetUnpaidVouchersByOutlet(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()
	qp := u.NewQueryParam(r)

	partnerID := bone.GetValue(r, "partner_id")
	startDate := r.FormValue("start_date")
	endDate := r.FormValue("end_date")

	qp.SetCompanyID(bone.GetValue(r, "company"))

	unpaidVouchers, next, err := model.GetUnpaidVouchersByOutlet(qp, partnerID, startDate, endDate)
	if err != nil && err != model.ErrorResourceNotFound {
		res.SetError(JSONErrFatal.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	list := []model.UnpaidVouchers{}
	list = append(list, unpaidVouchers...)
	res.SetResponse(list)
	if len(unpaidVouchers) > 0 {
		res.SetNewPagination(r, qp.Page, next, (unpaidVouchers)[0].Count)
	}
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
	program_id := bone.GetValue(r, "program_id")
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

//GetCashoutUsedVoucherDate : GET list of Cashout Summary by date
func GetCashoutUsedVoucherDate(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	qp := u.NewQueryParam(r)
	program_id := bone.GetValue(r, "program_id")
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
