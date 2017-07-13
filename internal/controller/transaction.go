package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-zoo/bone"
	"github.com/ruizu/render"

	"github.com/gilkor/evoucher/internal/model"
	"github.com/asaskevich/govalidator"
)

type (
	TransactionRequest struct {
		VariantID     string   `json:"variant_id" valid:"alphanum,required"`
		RedeemMethod  string   `json:"redeem_method" valid:"in(qr|token),required"`
		Partner       string   `json:"partner" valid:"alphanum,required"`
		Challenge     string   `json:"challenge" valid:"numeric,optional"`
		Response      string   `json:"response" valid:"numeric,optional"`
		DiscountValue string   `json:"discount_value" valid:"float,required"`
		Vouchers      []string `json:"vouchers" valid:"alphanum,required"`
	}
	DeleteTransactionRequest struct {
		User string `json:"requested_by"`
	}
	DateTransactionRequest struct {
		Start string `json:"start"`
		End   string `json:"end"`
	}
	TransactionResponse struct {
		TransactionCode string `json:"transaction_code"`
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
		logger.SetStatus(status).Log("param :", r.Body , "response :" , res.Errors.ToString())
		render.JSON(w, res, status)
		return
	}


	_, err := govalidator.ValidateStruct(rd)
	if err != nil {
		status = http.StatusBadRequest
		res.AddError(its(status), model.ErrCodeValidationError, model.ErrMessageValidationError+"("+err.Error()+")", logger.TraceID)
		logger.SetStatus(status).Log("param :", rd , "response :" , res.Errors.ToString())
		render.JSON(w, res, status)
		return
	}


	//check redeemtion method
	switch rd.RedeemMethod {
	case model.RedeemtionMethodQr:
		//to-do validate partner_id
		par := map[string]string{"variant_id": rd.VariantID, "id": rd.Partner}
		if _, err := model.FindVariantPartner(par); err == model.ErrResourceNotFound {
			status = http.StatusBadRequest
			res.AddError(its(status), model.ErrCodeResourceNotFound, model.ErrMessageInvalidQr, logger.TraceID)
			logger.SetStatus(status).Log("param :", rd , "response :" , res.Errors.ToString())
			render.JSON(w, res, status)
			return
		} else if err != nil {
			status = http.StatusInternalServerError
			res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", logger.TraceID)
			logger.SetStatus(status).Log("param :", rd , "response :" , res.Errors.ToString())
			render.JSON(w, res, status)
			return
		}
	case model.RedeemtionMethodToken:
		//to-do validate token
		par := map[string]string{"variant_id": rd.VariantID, "id": rd.Partner}
		if p, err := model.FindVariantPartner(par); err == model.ErrResourceNotFound {
			status = http.StatusBadRequest
			res.AddError(its(status), model.ErrCodeResourceNotFound, model.ErrMessageInvalidPaerner, logger.TraceID)
			logger.SetStatus(status).Log("param :", rd , "response :" , res.Errors.ToString())
			render.JSON(w, res, status)
			return
		} else if err != nil {
			status = http.StatusInternalServerError
			res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", logger.TraceID)
			logger.SetStatus(status).Log("param :", rd , "response :" , res.Errors.ToString())
			render.JSON(w, res, status)
			return
		} else {
			fmt.Println("panrner data : ", p[0].SerialNumber.String)

			if !OTPAuth(p[0].SerialNumber.String, rd.Challenge, rd.Response) {
				status = http.StatusBadRequest
				res.AddError(its(status), model.ErrCodeOTPFailed, model.ErrMessageOTPFailed, logger.TraceID)
				logger.SetStatus(status).Log("param :", rd , "response :" , res.Errors.ToString())
				render.JSON(w, res, status)
				return
			}
		}

	}

	if ok, err := CheckVariant(rd.RedeemMethod, rd.VariantID, len(rd.Vouchers)); !ok {
		switch err.Error() {
		case model.ErrCodeAllowAccumulativeDisable:
			status = http.StatusBadRequest
			res.AddError(its(status), err.Error(), model.ErrMessageAllowAccumulativeDisable, logger.TraceID)
			logger.SetStatus(status).Log("param :", rd , "response :" , res.Errors.ToString())
			render.JSON(w, res, status)
		case model.ErrCodeInvalidRedeemMethod:
			status = http.StatusBadRequest
			res.AddError(its(status), err.Error(), model.ErrMessageInvalidRedeemMethod, logger.TraceID)
			logger.SetStatus(status).Log("param :", rd , "response :" , res.Errors.ToString())
			render.JSON(w, res, status)
		case model.ErrCodeVoucherNotActive:
			status = http.StatusBadRequest
			res.AddError(its(status), err.Error(), model.ErrMessageVoucherNotActive, logger.TraceID)
			logger.SetStatus(status).Log("param :", rd , "response :" , res.Errors.ToString())
			render.JSON(w, res, status)
		case model.ErrResourceNotFound.Error():
			status = http.StatusBadRequest
			res.AddError(its(status), model.ErrCodeResourceNotFound, model.ErrMessageInvalidVariant, logger.TraceID)
			logger.SetStatus(status).Log("param :", rd , "response :" , res.Errors.ToString())
			render.JSON(w, res, status)
		case model.ErrCodeRedeemNotValidDay:
			status = http.StatusBadRequest
			res.AddError(its(status), err.Error(), model.ErrMessageRedeemNotValidDay, logger.TraceID)
			logger.SetStatus(status).Log("param :", rd , "response :" , res.Errors.ToString())
			render.JSON(w, res, status)
		case model.ErrCodeRedeemNotValidHour:
			status = http.StatusBadRequest
			res.AddError(its(status), err.Error(), model.ErrMessageRedeemNotValidHour, logger.TraceID)
			logger.SetStatus(status).Log("param :", rd , "response :" , res.Errors.ToString())
			render.JSON(w, res, status)

		default:
			status = http.StatusInternalServerError
			res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", logger.TraceID)
			logger.SetStatus(status).Log("param :", rd , "response :" , res.Errors.ToString())
			render.JSON(w, res, status)
		}
		return
	}

	// check validation all voucher & variant
	for _, v := range rd.Vouchers {
		if ok, err := rd.CheckVoucherRedeemtion(v); !ok {
			switch err.Error() {
			case model.ErrCodeVoucherNotActive:
				status = http.StatusBadRequest
				res.AddError(its(status), err.Error(), model.ErrMessageVoucherNotActive, logger.TraceID)
				logger.SetStatus(status).Log("param :", rd , "response :" , res.Errors.ToString())
				render.JSON(w, res, status)
			case model.ErrResourceNotFound.Error():
				status = http.StatusBadRequest
				res.AddError(its(status), model.ErrCodeResourceNotFound, err.Error(), logger.TraceID)
				logger.SetStatus(status).Log("param :", rd , "response :" , res.Errors.ToString())
				render.JSON(w, res, status)
			case model.ErrMessageVoucherAlreadyUsed:
				status = http.StatusBadRequest
				res.AddError(its(status), model.ErrCodeVoucherDisabled, err.Error(), logger.TraceID)
				logger.SetStatus(status).Log("param :", rd , "response :" , res.Errors.ToString())
				render.JSON(w, res, status)
			case model.ErrMessageVoucherAlreadyPaid:
				status = http.StatusBadRequest
				res.AddError(its(status), model.ErrCodeVoucherDisabled, model.ErrMessageVoucherAlreadyUsed, logger.TraceID)
				logger.SetStatus(status).Log("param :", rd , "response :" , res.Errors.ToString())
				render.JSON(w, res, status)
			case model.ErrMessageVoucherExpired:
				status = http.StatusBadRequest
				res.AddError(its(status), model.ErrCodeVoucherExpired, err.Error(), logger.TraceID)
				logger.SetStatus(status).Log("param :", rd , "response :" , res.Errors.ToString())
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
		DiscountValue:   stf(rd.DiscountValue),
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
		logger.SetStatus(status).Log("param :", rd , "response :" , res.Errors.ToString())
		render.JSON(w, res, status)
		return
	} else if err != nil {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", logger.TraceID)
		logger.SetStatus(status).Log("param :", rd , "response :" , res.Errors.ToString())
		render.JSON(w, res, status)
		return
	}

	res = NewResponse(TransactionResponse{TransactionCode: txCode})
	logger.SetStatus(status).Log("param :", rd , "response :" , TransactionResponse{TransactionCode: txCode})
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
		logger.SetStatus(status).Log("param :", r.Body , "response :" , res.Errors.ToString())
		render.JSON(w, res, status)
		return
	}

	_, err := govalidator.ValidateStruct(rd)
	if err != nil {
		status = http.StatusBadRequest
		res.AddError(its(status), model.ErrCodeValidationError, model.ErrMessageValidationError+"("+err.Error()+")", logger.TraceID)
		logger.SetStatus(status).Log("param :", rd , "response :" , res.Errors.ToString())
		render.JSON(w, res, status)
		return
	}

	//check redeemtion method
	switch rd.RedeemMethod {
	case model.RedeemtionMethodQr:
		//to-do validate partner_id
		par := map[string]string{"variant_id": rd.VariantID, "id": rd.Partner}
		if _, err := model.FindVariantPartner(par); err == model.ErrResourceNotFound {
			status = http.StatusBadRequest
			res.AddError(its(status), model.ErrCodeResourceNotFound, model.ErrMessageInvalidQr, logger.TraceID)
			logger.SetStatus(status).Log("param :", rd , "response :" , res.Errors.ToString())
			render.JSON(w, res, status)
			return
		} else if err != nil {
			status = http.StatusInternalServerError
			res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", logger.TraceID)
			logger.SetStatus(status).Log("param :", rd , "response :" , res.Errors.ToString())
			render.JSON(w, res, status)
			return
		}
	case model.RedeemtionMethodToken:
		//to-do validate token
		par := map[string]string{"variant_id": rd.VariantID, "id": rd.Partner}
		if p, err := model.FindVariantPartner(par); err == model.ErrResourceNotFound {
			status = http.StatusBadRequest
			res.AddError(its(status), model.ErrCodeResourceNotFound, model.ErrMessageInvalidQr, "partner")
			logger.SetStatus(status).Log("param :", rd , "response :" , res.Errors.ToString())
			render.JSON(w, res, status)
			return
		} else if err != nil {
			status = http.StatusInternalServerError
			res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", "partner")
			logger.SetStatus(status).Log("param :", rd , "response :" , res.Errors.ToString())
			render.JSON(w, res, status)
			return
		} else {
			fmt.Println("panrner data : ", p[0].SerialNumber.String)
			if !OTPAuth(p[0].SerialNumber.String, rd.Challenge, rd.Response) {
				status = http.StatusBadRequest
				res.AddError(its(status), model.ErrCodeOTPFailed, model.ErrMessageOTPFailed, logger.TraceID)
				logger.SetStatus(status).Log("param :", rd , "response :" , res.Errors.ToString())
				render.JSON(w, res, status)
				return
			}
		}

	}

	if ok, err := CheckVariant(rd.RedeemMethod, rd.VariantID, len(rd.Vouchers)); !ok {
		switch err.Error() {
		case model.ErrCodeAllowAccumulativeDisable:
			status = http.StatusBadRequest
			res.AddError(its(status), err.Error(), model.ErrMessageAllowAccumulativeDisable, logger.TraceID)
			logger.SetStatus(status).Log("param :", rd , "response :" , res.Errors.ToString())
			render.JSON(w, res, status)
		case model.ErrCodeInvalidRedeemMethod:
			status = http.StatusBadRequest
			res.AddError(its(status), err.Error(), model.ErrMessageInvalidRedeemMethod, logger.TraceID)
			logger.SetStatus(status).Log("param :", rd , "response :" , res.Errors.ToString())
			render.JSON(w, res, status)
		case model.ErrCodeVoucherNotActive:
			status = http.StatusBadRequest
			res.AddError(its(status), err.Error(), model.ErrMessageVoucherNotActive, logger.TraceID)
			logger.SetStatus(status).Log("param :", rd , "response :" , res.Errors.ToString())
			render.JSON(w, res, status)
		case model.ErrResourceNotFound.Error():
			status = http.StatusBadRequest
			res.AddError(its(status), model.ErrCodeResourceNotFound, model.ErrMessageInvalidVariant, logger.TraceID)
			logger.SetStatus(status).Log("param :", rd , "response :" , res.Errors.ToString())
			render.JSON(w, res, status)
		case model.ErrCodeRedeemNotValidDay:
			status = http.StatusBadRequest
			res.AddError(its(status), err.Error(), model.ErrMessageRedeemNotValidDay, logger.TraceID)
			logger.SetStatus(status).Log("param :", rd , "response :" , res.Errors.ToString())
			render.JSON(w, res, status)
		case model.ErrCodeRedeemNotValidHour:
			status = http.StatusBadRequest
			res.AddError(its(status), err.Error(), model.ErrMessageRedeemNotValidHour, logger.TraceID)
			logger.SetStatus(status).Log("param :", rd , "response :" , res.Errors.ToString())
			render.JSON(w, res, status)

		default:
			status = http.StatusInternalServerError
			res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", logger.TraceID)
			logger.SetStatus(status).Log("param :", rd , "response :" , res.Errors.ToString())
			render.JSON(w, res, status)
		}
		return
	}

	// check validation all voucher & variant
	for _, v := range rd.Vouchers {
		if ok, err := rd.CheckVoucherRedeemtion(v); !ok {
			switch err.Error() {
			case model.ErrCodeVoucherNotActive:
				status = http.StatusBadRequest
				res.AddError(its(status), err.Error(), model.ErrMessageVoucherNotActive, logger.TraceID)
				logger.SetStatus(status).Log("param :", rd , "response :" , res.Errors.ToString())
				render.JSON(w, res, status)
			case model.ErrResourceNotFound.Error():
				status = http.StatusBadRequest
				res.AddError(its(status), model.ErrCodeResourceNotFound, err.Error(), logger.TraceID)
				logger.SetStatus(status).Log("param :", rd , "response :" , res.Errors.ToString())
				render.JSON(w, res, status)
			case model.ErrMessageVoucherAlreadyUsed:
				status = http.StatusBadRequest
				res.AddError(its(status), model.ErrCodeVoucherDisabled, err.Error(), logger.TraceID)
				logger.SetStatus(status).Log("param :", rd , "response :" , res.Errors.ToString())
				render.JSON(w, res, status)
			case model.ErrMessageVoucherAlreadyPaid:
				status = http.StatusBadRequest
				res.AddError(its(status), model.ErrCodeVoucherDisabled, model.ErrMessageVoucherAlreadyUsed, logger.TraceID)
				logger.SetStatus(status).Log("param :", rd , "response :" , res.Errors.ToString())
				render.JSON(w, res, status)
			case model.ErrMessageVoucherExpired:
				status = http.StatusBadRequest
				res.AddError(its(status), model.ErrCodeVoucherExpired, err.Error(), logger.TraceID)
				logger.SetStatus(status).Log("param :", rd , "response :" , res.Errors.ToString())
				render.JSON(w, res, status)
			default:
				status = http.StatusInternalServerError
				res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", logger.TraceID)
				logger.SetStatus(status).Log("param :", rd , "response :" , res.Errors.ToString())
				render.JSON(w, res, status)
			}
			return
		}
	}

	txCode := randStr(model.DEFAULT_TXLENGTH, model.DEFAULT_TXCODE)

	variant, _ := model.FindVariantDetailsById(rd.VariantID)

	d := model.Transaction{
		AccountId:       variant.AccountId,
		PartnerId:       rd.Partner,
		TransactionCode: txCode,
		DiscountValue:   stf(rd.DiscountValue),
		Token:           rd.Response,
		User:            variant.CreatedBy,
		Vouchers:        rd.Vouchers,
	}
	fmt.Println(d)
	if err := model.InsertTransaction(d); err != nil {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", logger.TraceID)
		logger.SetStatus(status).Log("param :", rd , "response :" , res.Errors.ToString())
		render.JSON(w, res, status)
		return
	}

	rv := RedeemVoucherRequest{
		AccountID: variant.AccountId,
		User:      variant.CreatedBy,
		State:     model.VoucherStateUsed,
		Vouchers:  rd.Vouchers,
	}
	fmt.Println("List valid voucher :", rv.Vouchers)
	// update voucher state "Used"
	if ok, err := rv.UpdateVoucher(); !ok {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", logger.TraceID)
		logger.SetStatus(status).Log("param :", rd , "response :" , res.Errors.ToString())
		render.JSON(w, res, status)
		return
	} else if err != nil {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", logger.TraceID)
		logger.SetStatus(status).Log("param :", rd , "response :" , res.Errors.ToString())
		render.JSON(w, res, status)
		return
	}

	res = NewResponse(TransactionResponse{TransactionCode: txCode})
	logger.SetStatus(status).Log("param :", rd , "response :" , TransactionResponse{TransactionCode: txCode})
	render.JSON(w, res, status)
}

func GetAllTransactionsByPartner(w http.ResponseWriter, r *http.Request) {
	apiName := "report_transaction"
	valid := false
	partnerId := r.FormValue("partner")

	status := http.StatusUnauthorized
	err := model.ErrInvalidRole
	errTitle := model.ErrCodeInvalidRole
	res := NewResponse(nil)
	res.AddError(its(status), errTitle, err.Error(), "Get Transaction")

	a := AuthToken(w, r)
	if a.Valid {
		for _, valueRole := range a.User.Role {
			features := model.ApiFeatures[valueRole.RoleDetail]
			for _, valueFeature := range features {
				if apiName == valueFeature {
					valid = true
				}
			}
		}

		if valid {
			status = http.StatusOK
			transaction, _ := model.FindAllTransactionByPartner(a.User.Account.Id, partnerId)
			res = NewResponse(transaction)
		}
	} else {
		res = a.res
		status = http.StatusUnauthorized
	}
	render.JSON(w, res, status)
}

func CashoutTransactionDetails(w http.ResponseWriter, r *http.Request) {
	apiName := "transaction_get"
	valid := false

	transactionCode := r.FormValue("id")
	status := http.StatusUnauthorized
	err := model.ErrInvalidRole
	errTitle := model.ErrCodeInvalidRole
	res := NewResponse(nil)
	res.AddError(its(status), errTitle, err.Error(), "Get Transaction")

	a := AuthToken(w, r)

	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag("Get Transaction")

	if a.Valid {
		for _, valueRole := range a.User.Role {
			features := model.ApiFeatures[valueRole.RoleDetail]
			for _, valueFeature := range features {
				if apiName == valueFeature {
					valid = true
				}
			}
		}

		if valid {
			status = http.StatusOK
			variant, err := model.FindCashoutTransactionDetails(transactionCode)
			fmt.Println(err)
			if err != nil {
				status = http.StatusInternalServerError
				errTitle = model.ErrCodeInternalError
				if err == model.ErrResourceNotFound {
					status = http.StatusNotFound
					errTitle = model.ErrCodeResourceNotFound
				}

				res.AddError(its(status), errTitle, err.Error(), logger.TraceID)
				logger.SetStatus(status).Log("param :", transactionCode , "response :" , res.Errors.ToString())
			} else {
				res = NewResponse(variant)
				logger.SetStatus(status).Log("param :", transactionCode , "response :" , variant)
			}
		}
	} else {
		res = a.res
		status = http.StatusUnauthorized
	}

	render.JSON(w, res, status)
}

func PublicCashoutTransactionDetails(w http.ResponseWriter, r *http.Request) {
	transactionCode := bone.GetValue(r, "id")
	status := http.StatusUnauthorized
	err := model.ErrTokenNotFound
	errTitle := model.ErrCodeInvalidToken
	res := NewResponse(nil)


	res.AddError(its(status), errTitle, err.Error(), "Get Transaction")

	status = http.StatusOK
	variant, err := model.FindCashoutTransactionDetails(transactionCode)
	fmt.Println(err)
	if err != nil {
		status = http.StatusInternalServerError
		errTitle = model.ErrCodeInternalError
		if err == model.ErrResourceNotFound {
			status = http.StatusNotFound
			errTitle = model.ErrCodeResourceNotFound
		}

		res.AddError(its(status), errTitle, err.Error(), "Get Transaction")
	} else {
		res = NewResponse(variant)
	}

	render.JSON(w, res, status)
}

func CashoutTransactions(w http.ResponseWriter, r *http.Request) {
	apiName := "transaction_cashout"
	valid := false

	var rd TransactionCodeBulk
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		log.Panic(err)
	}

	res := NewResponse(nil)
	status := http.StatusUnauthorized
	err := model.ErrInvalidRole
	errTitle := model.ErrCodeInvalidRole
	res.AddError(its(status), errTitle, err.Error(), "Cashout Transaction")

	a := AuthToken(w, r)

	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag("Cash-Out Transaction")

	if a.Valid {
		for _, valueRole := range a.User.Role {
			features := model.ApiFeatures[valueRole.RoleDetail]
			for _, valueFeature := range features {
				if apiName == valueFeature {
					valid = true
				}
			}
		}

		if valid {
			status = http.StatusOK
			if err := model.UpdateCashoutTransactions(rd.TransactionCode, a.User.ID); err != nil {
				status = http.StatusInternalServerError
				errTitle = model.ErrCodeInternalError
				res.AddError(its(status), errTitle, err.Error(), logger.TraceID)
				logger.SetStatus(status).Log("param :", rd , "response :" , res.Errors.ToString())
			} else {
				res = NewResponse("Transaction Success")
				logger.SetStatus(status).Log("param :", rd , "response :" , res.Errors.ToString())
			}
		}
	} else {
		res = a.res
	}

	render.JSON(w, res, status)
}

func PrintCashoutTransaction(w http.ResponseWriter, r *http.Request) {
	apiName := "transaction_cashout"
	valid := false

	transactionCode := r.FormValue("transcation_code")
	transactionCodeArr := strings.Split(transactionCode, ";")

	status := http.StatusUnauthorized
	err := model.ErrInvalidRole
	errTitle := model.ErrCodeInvalidRole
	res := NewResponse(nil)
	res.AddError(its(status), errTitle, err.Error(), "Get Transaction")

	a := AuthToken(w, r)

	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag("Print Cash-Out Transaction")

	if a.Valid {
		for _, valueRole := range a.User.Role {
			features := model.ApiFeatures[valueRole.RoleDetail]
			for _, valueFeature := range features {
				if apiName == valueFeature {
					valid = true
				}
			}
		}

		if valid {
			status = http.StatusOK
			transaction, _ := model.PrintCashout(a.User.Account.Id, transactionCodeArr)
			res = NewResponse(transaction)
			logger.SetStatus(status).Log("param :", transactionCode , "response :" , transaction)
		}
	} else {
		res = a.res
		logger.SetStatus(status).Log("param :", transactionCode , "response :" , res.Errors.ToString())
	}

	render.JSON(w, res, status)
}

func DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	id := bone.GetValue(r, "id")
	var rd DeleteTransactionRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		log.Panic(err)
	}

	d := &model.DeleteTransactionRequest{
		Id:   id,
		User: rd.User,
	}
	if err := d.Delete(); err != nil {
		log.Panic(err)
	}

	res := NewResponse(nil)
	render.JSON(w, res, http.StatusOK)
}
