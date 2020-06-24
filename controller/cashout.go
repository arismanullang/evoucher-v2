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

type (
	CashoutFilter struct {
		ID            string `schema:"id" filter:"array"`
		AccountID     string `schema:"account_id" filter:"string"`
		Code          string `schema:"code" filter:"string"`
		OutletID      string `schema:"outlet_id" filter:"string"`
		BankAccount   string `schema:"bank_account" filter:"string"`
		Amount        string `schema:"amount" filter:"string"`
		PaymentMethod string `schema:"payment_method" filter:"string"`
		CreatedAt     string `schema:"created_at" filter:"date"`
		CreatedBy     string `schema:"created_by" filter:"string"`
		UpdatedAt     string `schema:"updated_at" filter:"date"`
		UpdatedBy     string `schema:"updated_by" filter:"string"`
		Status        string `schema:"status" filter:"enum"`
	}

	//CashoutRequest : cashout request struct
	CashoutRequest struct {
		OutletID    string `json:"outlet_id"`
		ReferenceNo string `json:"reference_no"`
		VoucherIDs  string `json:"voucher_ids"`
	}
)

//GetCashouts : GET list of Cashouts
func GetCashouts(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()
	qp := u.NewQueryParam(r)

	// startDate := r.FormValue("start_date")
	// endDate := r.FormValue("end_date")

	qp.SetCompanyID(bone.GetValue(r, "company"))

	var decoder = schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)

	var f CashoutFilter
	if err := decoder.Decode(&f, r.Form); err != nil {
		res.SetError(JSONErrFatal.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	qp.SetFilterModel(f)

	cashouts, next, err := model.GetCashouts(qp)
	if err != nil {
		res.SetError(JSONErrFatal.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	res.SetResponse(&cashouts)
	if len(*cashouts) > 0 {
		res.SetNewPagination(r, qp.Page, next, (*cashouts)[0].Count)
	}
	res.JSON(w, res, http.StatusOK)
}

type (
	CashoutSummaryFilter struct {
		Date       *time.Time `schema:"date" filter:"date"`
		OutletID   string     `schema:"outlet_id" filter:"array"`
		OutletName string     `schema:"outlet_name" filter:"string"`
	}
	CashoutUnpaidFilter struct {
		Date       *time.Time `schema:"date" filter:"date"`
		OutletID   string     `schema:"outlet_id" filter:"array"`
		OutletName string     `schema:"outlet_name" filter:"string"`
	}
)

//GetUnpaidCashout : GET list of unpaid cashout
func GetUnpaidCashout(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()
	qp := u.NewQueryParam(r)

	startDate := r.FormValue("start_date")
	endDate := r.FormValue("end_date")

	qp.SetCompanyID(bone.GetValue(r, "company"))

	unpaidReimburse, next, err := model.GetUnpaidCashout(qp, startDate, endDate)
	if err != nil && err != model.ErrorResourceNotFound {
		res.SetError(JSONErrFatal.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	list := []model.UnpaidCashout{}
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

	outletID := bone.GetValue(r, "outlet_id")
	startDate := r.FormValue("start_date")
	endDate := r.FormValue("end_date")

	qp.SetCompanyID(bone.GetValue(r, "company"))

	unpaidVouchers, next, err := model.GetUnpaidVouchersByOutlet(qp, outletID, startDate, endDate)
	if err != nil && err != model.ErrorResourceNotFound {
		res.SetError(JSONErrFatal.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	list := []model.VoucherTransaction{}
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
	cashout, _, err := model.GetCashoutByID(id, qp)
	if err != nil {
		res.SetError(JSONErrResourceNotFound)
		res.JSON(w, res, JSONErrResourceNotFound.Status)
		return
	}

	res.SetResponse(cashout)
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

//PostCashout : POST Cashout by outlet
func PostCashout(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()
	token := r.FormValue("token")

	companyID := bone.GetValue(r, "company")

	accData, err := model.GetSessionDataJWT(token, companyID)
	if err != nil {
		res.SetError(JSONErrUnauthorized)
		res.JSON(w, res, JSONErrUnauthorized.Status)
		return
	}

	var req CashoutRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		u.DEBUG(err)
		res.SetErrorWithDetail(JSONErrFatal, err)
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	outletBank, err := GetOutletBanks(r, req.OutletID)
	if err != nil {
		res.SetError(JSONErrBadRequest.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrBadRequest.Status)
		return
	}

	seedCode := u.RandomizeString(u.DEFAULT_LENGTH, u.NUMERALS)
	csCode := seedCode + u.RandomizeString(u.CASHOUT_CODE_LENGTH, u.NUMERALS)

	cashout := model.Cashout{
		CompanyID:       companyID,
		Code:            csCode,
		OutletID:        req.OutletID,
		BankName:        outletBank[0].BankName,
		BankCompanyName: outletBank[0].CompanyName,
		BankAccount:     outletBank[0].BankAccount,
		ReferenceNo:     req.ReferenceNo,
		PaymentMethod:   "bank_transfer",
		CreatedBy:       accData.AccountID,
		UpdatedBy:       accData.AccountID,
	}

	var f model.VoucherFilter
	f.ID = req.VoucherIDs
	voucherQP := u.NewQueryParam(r)
	voucherQP.Count = -1
	voucherQP.SetFilterModel(f)

	listVoucherByID, _, err := model.GetVouchers(voucherQP)
	if err != nil {
		res.SetError(JSONErrResourceNotFound)
		res.JSON(w, res, JSONErrResourceNotFound.Status)
		return
	}

	// transactionDetails, err := model.GetTransactionDetailByVoucherID(qp, req.VoucherIDs)
	// if err != nil {
	// 	res.SetError(JSONErrResourceNotFound)
	// 	res.JSON(w, res, JSONErrResourceNotFound.Status)
	// 	return
	// }

	// for _, td := range *transactionDetails{
	// 	transaction, err := model.GetTransactionByID(qp, td.TransactionId)
	// 	if err != nil {
	// Do we need to check if the vouchers used in the right outlet?
	// 	}
	// }

	totalAmount := float64(0)

	for _, voucher := range listVoucherByID {
		if voucher.State == model.VoucherStatePaid {
			res.SetError(JSONErrBadRequest.SetArgs("voucher has been paid"))
			res.JSON(w, res, JSONErrBadRequest.Status)
			return
		} else if voucher.State == model.VoucherStateCreated {
			res.SetError(JSONErrBadRequest.SetArgs("voucher has not been used yet"))
			res.JSON(w, res, JSONErrBadRequest.Status)
			return
		}

		totalAmount += voucher.ProgramMaxValue

		cashoutDetail := model.CashoutDetail{
			VoucherID: voucher.ID,
			CreatedBy: accData.AccountID,
			UpdatedBy: accData.AccountID,
		}
		cashout.CashoutDetails = append(cashout.CashoutDetails, cashoutDetail)
	}

	cashout.Amount = totalAmount

	response, err := cashout.Insert()
	if err != nil {
		u.DEBUG(err)
		res.SetErrorWithDetail(JSONErrFatal, err)
		res.JSON(w, res, JSONErrFatal.Status)
		return
	} else {
		for _, voucher := range listVoucherByID {
			voucher.State = model.VoucherStatePaid
			voucher.UpdatedBy = accData.AccountID
			if err = voucher.Update(); err != nil {
				res.SetError(JSONErrBadRequest.SetArgs(err.Error()))
				res.JSON(w, res, JSONErrBadRequest.Status)
				return
			}
		}
	}

	res.SetResponse(response)
	res.JSON(w, res, http.StatusCreated)
}

//GetCashoutVouchers : Get Reimburse Invoice / Detail
func GetCashoutVouchers(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	companyID := bone.GetValue(r, "company")
	id := bone.GetValue(r, "cashout_id")
	qp := u.NewQueryParam(r)

	qp.SetCompanyID(companyID)

	voucherTransactions, next, err := model.GetCashoutVouchers(qp, id)
	if err != nil {
		res.SetError(JSONErrBadRequest.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrBadRequest.Status)
		return
	}

	if len(voucherTransactions) > 0 {
		res.SetNewPagination(r, qp.Page, next, (voucherTransactions)[0].Count)
	}

	res.SetResponse(voucherTransactions)
	res.JSON(w, res, http.StatusCreated)
}
