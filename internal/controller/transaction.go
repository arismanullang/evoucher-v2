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
	//"golang.org/x/tools/go/gcimporter15/testdata"
	"time"
)

type (
	TransactionRequest struct {
		ProgramID     string   `json:"program_id" valid:"required"`
		RedeemMethod  string   `json:"redeem_method" valid:"in(qr|token),required"`
		Partner       string   `json:"partner" valid:"required"`
		Challenge     string   `json:"challenge" valid:"numeric,optional"`
		Response      string   `json:"response" valid:"numeric,optional"`
		DiscountValue string   `json:"discount_value" valid:"required"`
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

	seedCode := randStr(model.DEFAULT_TRANSACTION_LENGTH, model.DEFAULT_TRANSACTION_SEED)
	txCode := seedCode + randStr(model.DEFAULT_TXLENGTH, model.DEFAULT_TXCODE)

	d := model.Transaction{
		AccountId:       a.User.Account.Id,
		PartnerId:       rd.Partner,
		TransactionCode: txCode,
		DiscountValue:   stf(rd.DiscountValue) * float64(len(rd.Vouchers)),
		Token:           rd.Response,
		User:            a.User.ID,
		VoucherIds:      rd.Vouchers,
	}
	//fmt.Println(d)
	tId, err := model.InsertTransaction(d)
	if err != nil {
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

	// get list email
	listEmail := []string{}
	emails, err := model.GetEmail(tId)

	if strings.Contains(emails.EmailAccount, ";") {
		tempEmailAccount := strings.Split(emails.EmailAccount, ";")
		for _, v := range tempEmailAccount {
			listEmail = append(listEmail, v)
		}
	} else {
		listEmail = append(listEmail, emails.EmailAccount)
	}

	if strings.Contains(emails.EmailPartner, ";") {
		tempEmailPartner := strings.Split(emails.EmailPartner, ";")
		for _, v := range tempEmailPartner {
			listEmail = append(listEmail, v)
		}
	} else {
		listEmail = append(listEmail, emails.EmailPartner)
	}

	if strings.Contains(emails.EmailMember, ";") {
		tempEmailMember := strings.Split(emails.EmailMember, ";")
		for _, v := range tempEmailMember {
			listEmail = append(listEmail, v)
		}
	} else {
		listEmail = append(listEmail, emails.EmailMember)
	}

	// voucher detail
	voucherDetail, err := model.FindVouchersById(rd.Vouchers)
	if err != nil {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, "voucher error "+model.ErrMessageInternalError+"("+err.Error()+")", logger.TraceID)
		logger.SetStatus(status).Log("param :", rd, "voucher error response :", res.Errors.ToString())
		render.JSON(w, res, status)
		return
	}

	listVoucher := []string{}
	for _, v := range voucherDetail.VoucherData {
		listVoucher = append(listVoucher, v.VoucherCode)
	}

	// partner detail
	partner, err := model.FindPartnerById(rd.Partner)
	if err != nil {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, "partner error "+model.ErrMessageInternalError+"("+err.Error()+")", logger.TraceID)
		logger.SetStatus(status).Log("param :", rd, "partner error response :", res.Errors.ToString())
		render.JSON(w, res, status)
		return
	}

	req := model.ConfirmationEmailRequest{
		Holder:          voucherDetail.VoucherData[0].HolderDescription.String,
		ProgramName:     voucherDetail.VoucherData[0].ProgramName,
		PartnerName:     partner.Name,
		TransactionCode: txCode,
		TransactionDate: time.Now().Format("2006-01-02 15:04:05"),
		ListEmail:       listEmail,
		ListVoucher:     listVoucher,
	}

	if err := model.SendConfirmationEmail(model.Domain, model.ApiKey, model.PublicApiKey, "Sedayu One Voucher Confirmation", req, a.User.Account.Id); err != nil {
		res := NewResponse(nil)
		status := http.StatusInternalServerError
		errTitle := model.ErrCodeInternalError
		if err == model.ErrResourceNotFound {
			status = http.StatusNotFound
			errTitle = model.ErrCodeResourceNotFound
		}

		res.AddError(its(status), errTitle, err.Error(), logger.TraceID)
		logger.SetStatus(status).Info("param :", listEmail, "response :", err.Error())
		render.JSON(w, res, status)
		return
	}

	res = NewResponse(TransactionResponse{TransactionCode: txCode, Vouchers: listVoucher})
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

	program, _ := model.FindProgramDetailsById(rd.ProgramID)

	seedCode := randStr(model.DEFAULT_TRANSACTION_LENGTH, model.DEFAULT_TRANSACTION_SEED)
	txCode := seedCode + randStr(model.DEFAULT_TXLENGTH, model.DEFAULT_TXCODE)

	d := model.Transaction{
		AccountId:       program.AccountId,
		PartnerId:       rd.Partner,
		TransactionCode: txCode,
		DiscountValue:   stf(rd.DiscountValue),
		Token:           rd.Response,
		User:            rd.CreatedBy,
		VoucherIds:      rd.Vouchers,
	}

	tId, err := model.InsertTransaction(d)
	if err != nil {
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

	listEmail := []string{}
	emails, err := model.GetEmail(tId)

	if strings.Contains(emails.EmailAccount, ";") {
		tempEmailAccount := strings.Split(emails.EmailAccount, ";")
		for _, v := range tempEmailAccount {
			listEmail = append(listEmail, v)
		}
	} else {
		listEmail = append(listEmail, emails.EmailAccount)
	}

	if strings.Contains(emails.EmailPartner, ";") {
		tempEmailPartner := strings.Split(emails.EmailPartner, ";")
		for _, v := range tempEmailPartner {
			listEmail = append(listEmail, v)
		}
	} else {
		listEmail = append(listEmail, emails.EmailPartner)
	}

	if strings.Contains(emails.EmailMember, ";") {
		tempEmailMember := strings.Split(emails.EmailMember, ";")
		for _, v := range tempEmailMember {
			listEmail = append(listEmail, v)
		}
	} else {
		listEmail = append(listEmail, emails.EmailMember)
	}

	// voucher detail
	voucherDetail, err := model.FindVouchersById(rd.Vouchers)
	if err != nil {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, "voucher error "+model.ErrMessageInternalError+"("+err.Error()+")", logger.TraceID)
		logger.SetStatus(status).Log("param :", rd, "voucher error response :", res.Errors.ToString())
		render.JSON(w, res, status)
		return
	}

	listVoucher := []string{}
	for _, v := range voucherDetail.VoucherData {
		listVoucher = append(listVoucher, v.VoucherCode)
	}

	// partner detail
	partner, err := model.FindPartnerById(rd.Partner)
	if err != nil {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, "partner error "+model.ErrMessageInternalError+"("+err.Error()+")", logger.TraceID)
		logger.SetStatus(status).Log("param :", rd, "partner error response :", res.Errors.ToString())
		render.JSON(w, res, status)
		return
	}

	req := model.ConfirmationEmailRequest{
		Holder:          voucherDetail.VoucherData[0].HolderDescription.String,
		ProgramName:     voucherDetail.VoucherData[0].ProgramName,
		PartnerName:     partner.Name,
		TransactionCode: txCode,
		TransactionDate: time.Now().Format("2006-01-02 15:04:05"),
		ListEmail:       listEmail,
		ListVoucher:     listVoucher,
	}
	fmt.Println(partner.AccountId)
	if err := model.SendConfirmationEmail(model.Domain, model.ApiKey, model.PublicApiKey, "Sedayu One Voucher Confirmation", req, partner.AccountId); err != nil {
		res := NewResponse(nil)
		status := http.StatusInternalServerError
		errTitle := model.ErrCodeInternalError
		if err == model.ErrResourceNotFound {
			status = http.StatusNotFound
			errTitle = model.ErrCodeResourceNotFound
		}

		res.AddError(its(status), errTitle, err.Error(), logger.TraceID)
		logger.SetStatus(status).Info("param :", listEmail, "response :", err.Error())
		render.JSON(w, res, status)
		return
	}

	res = NewResponse(TransactionResponse{TransactionCode: txCode})
	logger.SetStatus(status).Log("param :", rd, "response :", TransactionResponse{TransactionCode: txCode, Vouchers: rd.Vouchers})
	render.JSON(w, res, status)
}

func GetTransactionsByPartner(w http.ResponseWriter, r *http.Request) {
	apiName := "report_transaction"
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

	if CheckAPIRole(a, apiName) {
		logger.SetStatus(status).Info("param :", a.User.ID, "response :", "Invalid Role")

		status = http.StatusUnauthorized
		res.AddError(its(status), model.ErrCodeInvalidRole, model.ErrInvalidRole.Error(), logger.TraceID)
		render.JSON(w, res, status)
		return
	}

	transaction, err := model.FindTransactionsByPartner(a.User.Account.Id, partnerId)
	res = NewResponse(transaction)
	if err != nil {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, err.Error(), logger.TraceID)
		logger.SetStatus(status).Info("param :", a.User.Account.Id+" || "+partnerId, "response :", res.Errors)
	}

	render.JSON(w, res, status)
}

func GetVoucherTransactionDetails(w http.ResponseWriter, r *http.Request) {
	apiName := "report_transaction"
	voucherId := r.FormValue("id")

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

	transaction, err := model.FindVoucherCycle(a.User.Account.Id, voucherId)
	res = NewResponse(transaction)
	if err != nil {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, err.Error(), logger.TraceID)
		logger.SetStatus(status).Info("param :", a.User.Account.Id+" || "+voucherId, "response :", res.Errors)
	}

	render.JSON(w, res, status)
}

func CashoutTransactionDetails(w http.ResponseWriter, r *http.Request) {
	apiName := "transaction_get"

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

	if CheckAPIRole(a, apiName) {
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
