package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-zoo/bone"
	"github.com/ruizu/render"

	"github.com/asaskevich/govalidator"
	"github.com/gilkor/evoucher/internal/model"
)

type (
	TransactionRequest struct {
		ProgramID     string   `json:"program_id" valid:"required"`
		RedeemMethod  string   `json:"redeem_method" valid:"in(qr|token),required"`
		Partner       string   `json:"partner" valid:"required"`
		Challenge     string   `json:"challenge" valid:"numeric,optional"`
		Response      string   `json:"response" valid:"numeric,optional"`
		DiscountValue string   `json:"discount_value" valid:"float,required"`
		Vouchers      []string `json:"vouchers" valid:"required"`
		CreatedBy     string   `json:"created_by" valid:"optional"`
	}
	DeleteTransactionRequest struct {
		User string `json:"requested_by"`
	}
	DateTransactionRequest struct {
		Start string `json:"start"`
		End   string `json:"end"`
	}
	TransactionResponse struct {
		TransactionCode string   `json:"transaction_code"`
		Vouchers        []string `json:"vouchers"`
	}
	TransactionCodeBulk struct {
		TransactionCode []string `json:"transaction_code"`
	}
)

func MobileCreateTransaction(w http.ResponseWriter, r *http.Request) {
	var rd TransactionRequest
	status := http.StatusCreated
	res := NewResponse(nil)

	//Token Authentocation
	a := AuthToken(w, r)
	if !a.Valid {
		render.JSON(w, a.res, status)
		return
	}

	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag("Mobile Redeem Transaction")

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		status = http.StatusBadRequest
		res.AddError(its(status), http.StatusText(status), http.StatusText(status)+"("+err.Error()+")", logger.TraceID)
		logger.SetStatus(status).Log("param :", r.Body, "response :", res.Errors.ToString())
		render.JSON(w, res, status)
		return
	}

	_, err := govalidator.ValidateStruct(rd)
	if err != nil {
		status = http.StatusBadRequest
		res.AddError(its(status), model.ErrCodeValidationError, model.ErrMessageValidationError+"("+err.Error()+")", logger.TraceID)
		logger.SetStatus(status).Log("param :", rd, "response :", res.Errors.ToString())
		render.JSON(w, res, status)
		return
	}

	//check redemption method
	switch rd.RedeemMethod {
	case model.RedemptionMethodQr:
		//to-do validate partner_id
		par := map[string]string{"program_id": rd.ProgramID, "id": rd.Partner}
		if _, err := model.FindProgramPartner(par); err == model.ErrResourceNotFound {
			status = http.StatusBadRequest
			res.AddError(its(status), model.ErrCodeResourceNotFound, model.ErrMessageInvalidQr, logger.TraceID)
			logger.SetStatus(status).Log("param :", rd, "response :", res.Errors.ToString())
			render.JSON(w, res, status)
			return
		} else if err != nil {
			status = http.StatusInternalServerError
			res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", logger.TraceID)
			logger.SetStatus(status).Log("param :", rd, "response :", res.Errors.ToString())
			render.JSON(w, res, status)
			return
		}
	case model.RedemptionMethodToken:
		//to-do validate token
		par := map[string]string{"program_id": rd.ProgramID, "id": rd.Partner}
		if p, err := model.FindProgramPartner(par); err == model.ErrResourceNotFound {
			status = http.StatusBadRequest
			res.AddError(its(status), model.ErrCodeResourceNotFound, model.ErrMessageInvalidPaerner, logger.TraceID)
			logger.SetStatus(status).Log("param :", rd, "response :", res.Errors.ToString())
			render.JSON(w, res, status)
			return
		} else if err != nil {
			status = http.StatusInternalServerError
			res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", logger.TraceID)
			logger.SetStatus(status).Log("param :", rd, "response :", res.Errors.ToString())
			render.JSON(w, res, status)
			return
		} else {
			fmt.Println("panrner data : ", p[0].SerialNumber.String)

			if !OTPAuth(p[0].SerialNumber.String, rd.Challenge, rd.Response) {
				status = http.StatusBadRequest
				res.AddError(its(status), model.ErrCodeOTPFailed, model.ErrMessageOTPFailed, logger.TraceID)
				logger.SetStatus(status).Log("param :", rd, "response :", res.Errors.ToString())
				render.JSON(w, res, status)
				return
			}
		}

	}

	if ok, err := CheckProgram(rd.RedeemMethod, rd.ProgramID, len(rd.Vouchers)); !ok {
		switch err.Error() {
		case model.ErrCodeAllowAccumulativeDisable:
			status = http.StatusBadRequest
			res.AddError(its(status), err.Error(), model.ErrMessageAllowAccumulativeDisable, logger.TraceID)
			logger.SetStatus(status).Log("param :", rd, "response :", res.Errors.ToString())
			render.JSON(w, res, status)
		case model.ErrCodeInvalidRedeemMethod:
			status = http.StatusBadRequest
			res.AddError(its(status), err.Error(), model.ErrMessageInvalidRedeemMethod, logger.TraceID)
			logger.SetStatus(status).Log("param :", rd, "response :", res.Errors.ToString())
			render.JSON(w, res, status)
		case model.ErrCodeVoucherNotActive:
			status = http.StatusBadRequest
			res.AddError(its(status), err.Error(), model.ErrMessageVoucherNotActive, logger.TraceID)
			logger.SetStatus(status).Log("param :", rd, "response :", res.Errors.ToString())
			render.JSON(w, res, status)
		case model.ErrResourceNotFound.Error():
			status = http.StatusBadRequest
			res.AddError(its(status), model.ErrCodeResourceNotFound, model.ErrMessageInvalidProgram, logger.TraceID)
			logger.SetStatus(status).Log("param :", rd, "response :", res.Errors.ToString())
			render.JSON(w, res, status)
		case model.ErrCodeRedeemNotValidDay:
			status = http.StatusBadRequest
			res.AddError(its(status), err.Error(), model.ErrMessageRedeemNotValidDay, logger.TraceID)
			logger.SetStatus(status).Log("param :", rd, "response :", res.Errors.ToString())
			render.JSON(w, res, status)
		case model.ErrCodeRedeemNotValidHour:
			status = http.StatusBadRequest
			res.AddError(its(status), err.Error(), model.ErrMessageRedeemNotValidHour, logger.TraceID)
			logger.SetStatus(status).Log("param :", rd, "response :", res.Errors.ToString())
			render.JSON(w, res, status)

		default:
			status = http.StatusInternalServerError
			res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", logger.TraceID)
			logger.SetStatus(status).Log("param :", rd, "response :", res.Errors.ToString())
			render.JSON(w, res, status)
		}
		return
	}

	// check validation all voucher & program
	for _, v := range rd.Vouchers {
		if ok, err := rd.CheckVoucherRedemption(v); !ok {
			switch err.Error() {
			case model.ErrCodeVoucherNotActive:
				status = http.StatusBadRequest
				res.AddError(its(status), err.Error(), model.ErrMessageVoucherNotActive, logger.TraceID)
				logger.SetStatus(status).Log("param :", rd, "response :", res.Errors.ToString())
				render.JSON(w, res, status)
			case model.ErrResourceNotFound.Error():
				status = http.StatusBadRequest
				res.AddError(its(status), model.ErrCodeResourceNotFound, err.Error(), logger.TraceID)
				logger.SetStatus(status).Log("param :", rd, "response :", res.Errors.ToString())
				render.JSON(w, res, status)
			case model.ErrMessageVoucherAlreadyUsed:
				status = http.StatusBadRequest
				res.AddError(its(status), model.ErrCodeVoucherDisabled, err.Error(), logger.TraceID)
				logger.SetStatus(status).Log("param :", rd, "response :", res.Errors.ToString())
				render.JSON(w, res, status)
			case model.ErrMessageVoucherAlreadyPaid:
				status = http.StatusBadRequest
				res.AddError(its(status), model.ErrCodeVoucherDisabled, model.ErrMessageVoucherAlreadyUsed, logger.TraceID)
				logger.SetStatus(status).Log("param :", rd, "response :", res.Errors.ToString())
				render.JSON(w, res, status)
			case model.ErrMessageVoucherExpired:
				status = http.StatusBadRequest
				res.AddError(its(status), model.ErrCodeVoucherExpired, err.Error(), logger.TraceID)
				logger.SetStatus(status).Log("param :", rd, "response :", res.Errors.ToString())
				render.JSON(w, res, status)
			default:
				status = http.StatusInternalServerError
				res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", logger.TraceID)
				render.JSON(w, res, status)
			}
			return
		}
	}

	txCode := randStr(model.DEFAULT_TXLENGTH, model.DEFAULT_TXCODE)
	d := model.Transaction{
		AccountId:       a.User.Account.Id,
		PartnerId:       rd.Partner,
		TransactionCode: txCode,
		DiscountValue:   stf(rd.DiscountValue) * float64(len(rd.Vouchers)),
		Token:           rd.Response,
		User:            a.User.ID,
		Vouchers:        rd.Vouchers,
	}
	//fmt.Println(d)
	if err := model.InsertTransaction(d); err != nil {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", logger.TraceID)
		render.JSON(w, res, status)
		return
	}

	rv := RedeemVoucherRequest{
		AccountID: a.User.Account.Id,
		User:      a.User.ID,
		State:     model.VoucherStateUsed,
		Vouchers:  rd.Vouchers,
	}
	fmt.Println("List valid voucher :", rv.Vouchers)
	// update voucher state "Used"
	if ok, err := rv.UpdateVoucher(); !ok {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", logger.TraceID)
		logger.SetStatus(status).Log("param :", rd, "response :", res.Errors.ToString())
		render.JSON(w, res, status)
		return
	} else if err != nil {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", logger.TraceID)
		logger.SetStatus(status).Log("param :", rd, "response :", res.Errors.ToString())
		render.JSON(w, res, status)
		return
	}

	res = NewResponse(TransactionResponse{TransactionCode: txCode, Vouchers: rd.Vouchers})
	logger.SetStatus(status).Log("param :", rd, "response :", res)
	render.JSON(w, res, status)
}

func WebCreateTransaction(w http.ResponseWriter, r *http.Request) {
	var rd TransactionRequest
	status := http.StatusCreated
	res := NewResponse(nil)

	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag("My Voucher")

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		status = http.StatusBadRequest
		res.AddError(its(status), http.StatusText(status), http.StatusText(status)+"("+err.Error()+")", logger.TraceID)
		logger.SetStatus(status).Log("param :", r.Body, "response :", res.Errors.ToString())
		render.JSON(w, res, status)
		return
	}

	_, err := govalidator.ValidateStruct(rd)
	if err != nil {
		status = http.StatusBadRequest
		res.AddError(its(status), model.ErrCodeValidationError, model.ErrMessageValidationError+"("+err.Error()+")", logger.TraceID)
		logger.SetStatus(status).Log("param :", rd, "response :", res.Errors.ToString())
		render.JSON(w, res, status)
		return
	}

	//check redemption method
	switch rd.RedeemMethod {
	case model.RedemptionMethodQr:
		//to-do validate partner_id
		par := map[string]string{"program_id": rd.ProgramID, "id": rd.Partner}
		if _, err := model.FindProgramPartner(par); err == model.ErrResourceNotFound {
			status = http.StatusBadRequest
			res.AddError(its(status), model.ErrCodeResourceNotFound, model.ErrMessageInvalidQr, logger.TraceID)
			logger.SetStatus(status).Log("param :", rd, "response :", res.Errors.ToString())
			render.JSON(w, res, status)
			return
		} else if err != nil {
			status = http.StatusInternalServerError
			res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", logger.TraceID)
			logger.SetStatus(status).Log("param :", rd, "response :", res.Errors.ToString())
			render.JSON(w, res, status)
			return
		}
	case model.RedemptionMethodToken:
		//to-do validate token
		par := map[string]string{"program_id": rd.ProgramID, "id": rd.Partner}
		if p, err := model.FindProgramPartner(par); err == model.ErrResourceNotFound {
			status = http.StatusBadRequest
			res.AddError(its(status), model.ErrCodeResourceNotFound, model.ErrMessageInvalidQr, "partner")
			logger.SetStatus(status).Log("param :", rd, "response :", res.Errors.ToString())
			render.JSON(w, res, status)
			return
		} else if err != nil {
			status = http.StatusInternalServerError
			res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", "partner")
			logger.SetStatus(status).Log("param :", rd, "response :", res.Errors.ToString())
			render.JSON(w, res, status)
			return
		} else {
			fmt.Println("panrner data : ", p[0].SerialNumber.String)
			if !OTPAuth(p[0].SerialNumber.String, rd.Challenge, rd.Response) {
				status = http.StatusBadRequest
				res.AddError(its(status), model.ErrCodeOTPFailed, model.ErrMessageOTPFailed, logger.TraceID)
				logger.SetStatus(status).Log("param :", rd, "response :", res.Errors.ToString())
				render.JSON(w, res, status)
				return
			}
		}

	}

	if ok, err := CheckProgram(rd.RedeemMethod, rd.ProgramID, len(rd.Vouchers)); !ok {
		switch err.Error() {
		case model.ErrCodeAllowAccumulativeDisable:
			status = http.StatusBadRequest
			res.AddError(its(status), err.Error(), model.ErrMessageAllowAccumulativeDisable, logger.TraceID)
			logger.SetStatus(status).Log("param :", rd, "response :", res.Errors.ToString())
			render.JSON(w, res, status)
		case model.ErrCodeInvalidRedeemMethod:
			status = http.StatusBadRequest
			res.AddError(its(status), err.Error(), model.ErrMessageInvalidRedeemMethod, logger.TraceID)
			logger.SetStatus(status).Log("param :", rd, "response :", res.Errors.ToString())
			render.JSON(w, res, status)
		case model.ErrCodeVoucherNotActive:
			status = http.StatusBadRequest
			res.AddError(its(status), err.Error(), model.ErrMessageVoucherNotActive, logger.TraceID)
			logger.SetStatus(status).Log("param :", rd, "response :", res.Errors.ToString())
			render.JSON(w, res, status)
		case model.ErrResourceNotFound.Error():
			status = http.StatusBadRequest
			res.AddError(its(status), model.ErrCodeResourceNotFound, model.ErrMessageInvalidProgram, logger.TraceID)
			logger.SetStatus(status).Log("param :", rd, "response :", res.Errors.ToString())
			render.JSON(w, res, status)
		case model.ErrCodeRedeemNotValidDay:
			status = http.StatusBadRequest
			res.AddError(its(status), err.Error(), model.ErrMessageRedeemNotValidDay, logger.TraceID)
			logger.SetStatus(status).Log("param :", rd, "response :", res.Errors.ToString())
			render.JSON(w, res, status)
		case model.ErrCodeRedeemNotValidHour:
			status = http.StatusBadRequest
			res.AddError(its(status), err.Error(), model.ErrMessageRedeemNotValidHour, logger.TraceID)
			logger.SetStatus(status).Log("param :", rd, "response :", res.Errors.ToString())
			render.JSON(w, res, status)

		default:
			status = http.StatusInternalServerError
			res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", logger.TraceID)
			logger.SetStatus(status).Log("param :", rd, "response :", res.Errors.ToString())
			render.JSON(w, res, status)
		}
		return
	}

	// check validation all voucher & program
	for _, v := range rd.Vouchers {
		if ok, err := rd.CheckVoucherRedemption(v); !ok {
			switch err.Error() {
			case model.ErrCodeVoucherNotActive:
				status = http.StatusBadRequest
				res.AddError(its(status), err.Error(), model.ErrMessageVoucherNotActive, logger.TraceID)
				logger.SetStatus(status).Log("param :", rd, "response :", res.Errors.ToString())
				render.JSON(w, res, status)
			case model.ErrResourceNotFound.Error():
				status = http.StatusBadRequest
				res.AddError(its(status), model.ErrCodeResourceNotFound, err.Error(), logger.TraceID)
				logger.SetStatus(status).Log("param :", rd, "response :", res.Errors.ToString())
				render.JSON(w, res, status)
			case model.ErrMessageVoucherAlreadyUsed:
				status = http.StatusBadRequest
				res.AddError(its(status), model.ErrCodeVoucherDisabled, err.Error(), logger.TraceID)
				logger.SetStatus(status).Log("param :", rd, "response :", res.Errors.ToString())
				render.JSON(w, res, status)
			case model.ErrMessageVoucherAlreadyPaid:
				status = http.StatusBadRequest
				res.AddError(its(status), model.ErrCodeVoucherDisabled, model.ErrMessageVoucherAlreadyUsed, logger.TraceID)
				logger.SetStatus(status).Log("param :", rd, "response :", res.Errors.ToString())
				render.JSON(w, res, status)
			case model.ErrMessageVoucherExpired:
				status = http.StatusBadRequest
				res.AddError(its(status), model.ErrCodeVoucherExpired, err.Error(), logger.TraceID)
				logger.SetStatus(status).Log("param :", rd, "response :", res.Errors.ToString())
				render.JSON(w, res, status)
			default:
				status = http.StatusInternalServerError
				res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", logger.TraceID)
				logger.SetStatus(status).Log("param :", rd, "response :", res.Errors.ToString())
				render.JSON(w, res, status)
			}
			return
		}
	}

	txCode := randStr(model.DEFAULT_TXLENGTH, model.DEFAULT_TXCODE)

	program, _ := model.FindProgramDetailsById(rd.ProgramID)

	d := model.Transaction{
		AccountId:       program.AccountId,
		PartnerId:       rd.Partner,
		TransactionCode: txCode,
		DiscountValue:   stf(rd.DiscountValue),
		Token:           rd.Response,
		User:            rd.CreatedBy,
		Vouchers:        rd.Vouchers,
	}
	fmt.Println(d)
	if err := model.InsertTransaction(d); err != nil {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", logger.TraceID)
		logger.SetStatus(status).Log("param :", rd, "response :", res.Errors.ToString())
		render.JSON(w, res, status)
		return
	}

	rv := RedeemVoucherRequest{
		AccountID: program.AccountId,
		User:      program.CreatedBy,
		State:     model.VoucherStateUsed,
		Vouchers:  rd.Vouchers,
	}
	fmt.Println("List valid voucher :", rv.Vouchers)
	// update voucher state "Used"
	if ok, err := rv.UpdateVoucher(); !ok {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", logger.TraceID)
		logger.SetStatus(status).Log("param :", rd, "response :", res.Errors.ToString())
		render.JSON(w, res, status)
		return
	} else if err != nil {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", logger.TraceID)
		logger.SetStatus(status).Log("param :", rd, "response :", res.Errors.ToString())
		render.JSON(w, res, status)
		return
	}

	res = NewResponse(TransactionResponse{TransactionCode: txCode})
	logger.SetStatus(status).Log("param :", rd, "response :", TransactionResponse{TransactionCode: txCode})
	render.JSON(w, res, status)
}

func GetAllTransactionsByPartner(w http.ResponseWriter, r *http.Request) {
	apiName := "report_transaction"
	valid := false
	partnerId := r.FormValue("partner")

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

	for _, valueRole := range a.User.Role {
		features := model.ApiFeatures[valueRole.Detail]
		for _, valueFeature := range features {
			if apiName == valueFeature {
				valid = true
			}
		}
	}

	if !valid {
		logger.SetStatus(status).Info("param :", a.User.ID, "response :", "Invalid Role")

		status = http.StatusUnauthorized
		res.AddError(its(status), model.ErrCodeInvalidRole, model.ErrInvalidRole.Error(), logger.TraceID)
		render.JSON(w, res, status)
		return
	}

	transaction, err := model.FindAllTransactionByPartner(a.User.Account.Id, partnerId)
	res = NewResponse(transaction)
	if err != nil {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, err.Error(), logger.TraceID)
		logger.SetStatus(status).Info("param :", a.User.Account.Id+" || "+partnerId, "response :", res.Errors)
	}

	render.JSON(w, res, status)
}

func CashoutTransactionDetails(w http.ResponseWriter, r *http.Request) {
	apiName := "transaction_get"
	valid := false

	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag(apiName)

	transactionCode := r.FormValue("id")
	status := http.StatusOK

	res := NewResponse(nil)

	a := AuthTokenWithLogger(w, r, logger)
	if !a.Valid {
		res = a.res
		status = http.StatusUnauthorized
		render.JSON(w, res, status)
		return
	}

	for _, valueRole := range a.User.Role {
		features := model.ApiFeatures[valueRole.Detail]
		for _, valueFeature := range features {
			if apiName == valueFeature {
				valid = true
			}
		}
	}

	if !valid {
		logger.SetStatus(status).Info("param :", a.User.ID, "response :", "Invalid Role")

		status = http.StatusUnauthorized
		res.AddError(its(status), model.ErrCodeInvalidRole, model.ErrInvalidRole.Error(), logger.TraceID)
		render.JSON(w, res, status)
		return
	}

	program, err := model.FindCashoutTransactionDetails(transactionCode)
	res = NewResponse(program)
	if err != nil {
		status = http.StatusInternalServerError
		errTitle := model.ErrCodeInternalError
		if err == model.ErrResourceNotFound {
			status = http.StatusNotFound
			errTitle = model.ErrCodeResourceNotFound
		}

		res.AddError(its(status), errTitle, err.Error(), logger.TraceID)
		logger.SetStatus(status).Info("param :", transactionCode, "response :", res.Errors)
	}

	render.JSON(w, res, status)
}

func PublicCashoutTransactionDetails(w http.ResponseWriter, r *http.Request) {
	transactionCode := bone.GetValue(r, "id")

	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag("public_transaction_details")

	status := http.StatusOK
	program, err := model.FindCashoutTransactionDetails(transactionCode)
	res := NewResponse(program)
	if err != nil {
		status = http.StatusInternalServerError
		errTitle := model.ErrCodeInternalError
		if err == model.ErrResourceNotFound {
			status = http.StatusNotFound
			errTitle = model.ErrCodeResourceNotFound
		}

		res.AddError(its(status), errTitle, err.Error(), logger.TraceID)
	}

	render.JSON(w, res, status)
}

func CashoutTransactions(w http.ResponseWriter, r *http.Request) {
	apiName := "transaction_cashout"
	valid := false

	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag(apiName)

	var rd TransactionCodeBulk
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		logger.SetStatus(http.StatusInternalServerError).Panic("param :", r.Body, "response :", err.Error())
	}

	res := NewResponse(nil)
	status := http.StatusOK

	a := AuthTokenWithLogger(w, r, logger)
	if !a.Valid {
		res = a.res
		status = http.StatusUnauthorized
		render.JSON(w, res, status)
		return
	}

	for _, valueRole := range a.User.Role {
		features := model.ApiFeatures[valueRole.Detail]
		for _, valueFeature := range features {
			if apiName == valueFeature {
				valid = true
			}
		}
	}

	if !valid {
		logger.SetStatus(status).Info("param :", a.User.ID, "response :", "Invalid Role")

		status = http.StatusUnauthorized
		res.AddError(its(status), model.ErrCodeInvalidRole, model.ErrInvalidRole.Error(), logger.TraceID)
		render.JSON(w, res, status)
		return
	}

	if err := model.UpdateCashoutTransactions(rd.TransactionCode, a.User.ID); err != nil {
		status = http.StatusInternalServerError
		errTitle := model.ErrCodeInternalError
		res.AddError(its(status), errTitle, err.Error(), logger.TraceID)

		logger.SetStatus(status).Info("param :", a.User.ID+" || "+strings.Join(rd.TransactionCode, ";"), "response :", res.Errors)
	}

	render.JSON(w, res, status)
}

func PrintCashoutTransaction(w http.ResponseWriter, r *http.Request) {
	apiName := "transaction_cashout"
	valid := false

	transactionCode := r.FormValue("transcation_code")
	transactionCodeArr := strings.Split(transactionCode, ";")

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

	for _, valueRole := range a.User.Role {
		features := model.ApiFeatures[valueRole.Detail]
		for _, valueFeature := range features {
			if apiName == valueFeature {
				valid = true
			}
		}
	}

	if !valid {
		logger.SetStatus(status).Info("param :", a.User.ID, "response :", "Invalid Role")

		status = http.StatusUnauthorized
		res.AddError(its(status), model.ErrCodeInvalidRole, model.ErrInvalidRole.Error(), logger.TraceID)
		render.JSON(w, res, status)
		return
	}

	status = http.StatusOK
	transaction, err := model.PrintCashout(a.User.Account.Id, transactionCodeArr)
	res = NewResponse(transaction)
	if err != nil {
		status = http.StatusInternalServerError
		errTitle := model.ErrCodeInternalError
		res.AddError(its(status), errTitle, err.Error(), logger.TraceID)

		logger.SetStatus(status).Info("param :", a.User.Account.Id+" || "+transactionCode, "response :", res.Errors)
	}

	render.JSON(w, res, status)
}
