package controller

import (
	"encoding/json"
	"net/http"

	"github.com/gilkor/evoucher/internal/model"
	"github.com/ruizu/render"
)

type (
	CashoutRequest struct {
		PartnerId            string   `json:"partner_id"`
		BankAccount          string   `json:"bank_account"`
		BankAccountCompany   string   `json:"bank_account_company"`
		BankAccountNumber    string   `json:"bank_account_number"`
		BankAccountRefNumber string   `json:"bank_account_ref_number"`
		TotalCashout         float64  `json:"total_cashout"`
		PaymentMethod        string   `json:"payment_method"`
		Transactions         []string `json:"transactions"`
		Vouchers             []string `json:"vouchers"`
	}

	VoidRequest struct {
		CashoutID   string `json:"cashout_id"`
		Description string `json:"description"`
	}
)

func CashoutVoid(w http.ResponseWriter, r *http.Request) {
	apiName := "cashout_void"

	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag(apiName)

	res := NewResponse(nil)

	var vr VoidRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&vr); err != nil {
		logger.SetStatus(http.StatusBadRequest).Panic("param :", r.Body, "response :", err.Error())
		res.AddError(its(http.StatusBadRequest), model.ErrCodeJsonError, err.Error(), logger.TraceID)
	}

	status := http.StatusOK

	a := AuthTokenWithLogger(w, r, logger)
	if !a.Valid {
		res = a.res
		status = http.StatusUnauthorized
		render.JSON(w, res, status)
		return
	}

	if CheckAPIRole(a, apiName) {
		logger.SetStatus(status).Info("param :", a.User.ID, "response :", "Invalid Role")

		status = http.StatusUnauthorized
		res.AddError(its(status), model.ErrCodeInvalidRole, "Invalid Role : "+model.ErrInvalidRole.Error(), logger.TraceID)
		render.JSON(w, res, status)
		return
	}

	trxIDs, err := model.VoidCashout(vr.CashoutID, a.User.ID)
	if err != nil {
		status = http.StatusInternalServerError
		errTitle := model.ErrCodeInternalError
		res.AddError(its(status), errTitle, "Void Cashout : "+err.Error(), logger.TraceID)

		logger.SetStatus(status).Info("param :", vr, "response :", res.Errors)
	} else if len(trxIDs) > 0 {
		if err := model.UpdateCashoutTransactions(trxIDs, a.User.ID, model.VoucherStateUsed); err != nil {
			status = http.StatusInternalServerError
			errTitle := model.ErrCodeInternalError
			res.AddError(its(status), errTitle, "Update Voucher : "+err.Error(), logger.TraceID)

			logger.SetStatus(status).Info("param :", vr, "response :", res.Errors)
		} else {
			res = NewResponse("success")
		}
	} else if len(trxIDs) == 0 {
		status = http.StatusNotFound
		errTitle := model.ErrCodeCashoutNotFound
		res.AddError(its(status), errTitle, model.ErrMessageCashoutNotFound, logger.TraceID)

		logger.SetStatus(status).Info("param :", vr, "response :", res.Errors)
	}

	render.JSON(w, res, status)
}

func CashoutTransactions(w http.ResponseWriter, r *http.Request) {
	apiName := "transaction_cashout"

	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag(apiName)

	res := NewResponse(nil)

	var rd CashoutRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		logger.SetStatus(http.StatusBadRequest).Panic("param :", r.Body, "response :", err.Error())
		res.AddError(its(http.StatusBadRequest), model.ErrCodeJsonError, err.Error(), logger.TraceID)
	}

	status := http.StatusOK

	a := AuthTokenWithLogger(w, r, logger)
	if !a.Valid {
		res = a.res
		status = http.StatusUnauthorized
		render.JSON(w, res, status)
		return
	}

	if CheckAPIRole(a, apiName) {
		logger.SetStatus(status).Info("param :", a.User.ID, "response :", "Invalid Role")

		status = http.StatusUnauthorized
		res.AddError(its(status), model.ErrCodeInvalidRole, "Invalid Role : "+model.ErrInvalidRole.Error(), logger.TraceID)
		render.JSON(w, res, status)
		return
	}

	voucherState, err := model.FindVouchersState(rd.Vouchers)
	if err != nil {
		logger.SetStatus(status).Info("param :", a.User.ID, "response :", "Invalid Role")

		status = http.StatusUnauthorized
		res.AddError(its(status), model.ErrCodeInvalidRole, "Invalid Role : "+model.ErrInvalidRole.Error(), logger.TraceID)
		render.JSON(w, res, status)
		return
	}

	for _, v := range voucherState {
		if v.State != model.VoucherStateUsed {
			logger.SetStatus(status).Info("param :", a.User.ID, "response :", "Invalid Role")

			status = http.StatusBadRequest
			res.AddError(its(status), model.ErrServerInternal.Error(), model.ErrMessageVoucherAlreadyPaid, logger.TraceID)
			render.JSON(w, res, status)
			return
		}
	}

	transactions := []model.CashoutTransaction{}
	for i, v := range rd.Transactions {
		transactions = append(transactions, model.CashoutTransaction{TransactionId: v, VoucherId: rd.Vouchers[i]})
	}

	seedCode := randStr(model.DEFAULT_TRANSACTION_LENGTH, model.DEFAULT_TRANSACTION_SEED)
	csCode := seedCode + randStr(model.DEFAULT_TXLENGTH, model.DEFAULT_TXCODE)

	cashout := model.Cashout{
		AccountId:            a.User.Account.Id,
		CashoutCode:          csCode,
		PartnerId:            rd.PartnerId,
		BankAccount:          rd.BankAccount,
		BankAccountCompany:   rd.BankAccountCompany,
		BankAccountNumber:    rd.BankAccountNumber,
		BankAccountRefNumber: rd.BankAccountRefNumber,
		TotalCashout:         rd.TotalCashout,
		PaymentMethod:        rd.PaymentMethod,
		CreatedBy:            a.User.ID,
		Transactions:         transactions,
	}

	id, err := model.InsertCashout(cashout)
	if err != nil {
		status = http.StatusInternalServerError
		errTitle := model.ErrCodeInternalError
		res.AddError(its(status), errTitle, "Insert Cashout : "+err.Error(), logger.TraceID)

		logger.SetStatus(status).Info("param :", cashout, "response :", res.Errors)
	} else {
		res = NewResponse(id)
		if err := model.UpdateCashoutTransactions(rd.Transactions, a.User.ID, model.VoucherStatePaid); err != nil {
			status = http.StatusInternalServerError
			errTitle := model.ErrCodeInternalError
			res.AddError(its(status), errTitle, "Update Voucher : "+err.Error(), logger.TraceID)

			logger.SetStatus(status).Info("param :", cashout, "response :", res.Errors)
		}
	}

	render.JSON(w, res, status)
}

func PrintCashoutTransaction(w http.ResponseWriter, r *http.Request) {
	apiName := "transaction_cashout"

	cashoutCode := r.FormValue("transcation_code")

	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag(apiName)

	status := http.StatusOK
	res := NewResponse(nil)

	a := AuthTokenWithLogger(w, r, logger)
	if !a.Valid {
		res = a.res
		status = http.StatusUnauthorized
		render.JSON(w, res, status)
		return
	}

	if CheckAPIRole(a, apiName) {
		logger.SetStatus(status).Info("param :", a.User.ID, "response :", "Invalid Role")

		status = http.StatusUnauthorized
		res.AddError(its(status), model.ErrCodeInvalidRole, model.ErrInvalidRole.Error(), logger.TraceID)
		render.JSON(w, res, status)
		return
	}

	status = http.StatusOK
	transaction, err := model.PrintCashout(a.User.Account.Id, cashoutCode)
	res = NewResponse(transaction)
	if err != nil {
		status = http.StatusInternalServerError
		errTitle := model.ErrCodeInternalError
		res.AddError(its(status), errTitle, "Print : "+err.Error(), logger.TraceID)

		logger.SetStatus(status).Info("param :", a.User.Account.Id+" || "+cashoutCode, "response :", res.Errors)
	}

	render.JSON(w, res, status)
}

func GetReimburseHistory(w http.ResponseWriter, r *http.Request) {
	apiName := "report_cashout"
	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag(apiName)

	status := http.StatusOK
	res := NewResponse(nil)

	a := AuthTokenWithLogger(w, r, logger)
	if !a.Valid {
		res = a.res
		status = http.StatusUnauthorized
		render.JSON(w, res, status)
		return
	}

	if CheckAPIRole(a, apiName) {
		logger.SetStatus(status).Info("param :", a.User.ID, "response :", "Invalid Role")

		status = http.StatusUnauthorized
		res.AddError(its(status), model.ErrCodeInvalidRole, model.ErrInvalidRole.Error(), logger.TraceID)
		render.JSON(w, res, status)
		return
	}

	cashout, err := model.FindAllReimburse(a.User.Account.Id, a.User.ID)
	if err != nil {
		logger.SetStatus(status).Info("param :", a.User.ID, "response :", "Invalid Role")

		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrServerInternal.Error(), err.Error(), logger.TraceID)
		render.JSON(w, res, status)
		return
	}

	status = http.StatusOK
	res = NewResponse(cashout)
	logger.SetStatus(status).Log("response :", model.ProgramType)
	render.JSON(w, res, status)
}
