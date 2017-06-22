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
	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type (
	TransactionRequest struct {
		VariantID     string   `json:"variant_id"`
		RedeemMethod  string   `json:"redeem_method"`
		Partner       string   `json:"partner"`
		Challenge     string   `json:"challenge"`
		Response      string   `json:"response"`
		DiscountValue string   `json:"discount_value"`
		Vouchers      []string `json:"vouchers"`
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

func (t TransactionRequest) validate() error {
	return validation.ValidateStruct(&t,
		validation.Field(&t.VariantID, validation.Required),
		validation.Field(&t.RedeemMethod, validation.Required, validation.In("qr", "token")),
		validation.Field(&t.Partner, validation.Required),
		validation.Field(&t.Challenge, is.Digit),
		validation.Field(&t.Response, is.Digit),
		validation.Field(&t.DiscountValue, validation.Required, is.Float),
		validation.Field(&t.Vouchers, validation.Required),
	)

}

func MobileCreateTransaction(w http.ResponseWriter, r *http.Request) {
	var rd TransactionRequest
	status := http.StatusCreated
	res := NewResponse(nil)

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		status = http.StatusBadRequest
		res.AddError(its(status), http.StatusText(status), http.StatusText(status)+"("+err.Error()+")", "transaction")
		render.JSON(w, res, status)
		return
	}

	//validate request param
	//err := rd.validate()
	//if err !=nil {
	//	status = http.StatusBadRequest
	//	res.AddError(its(status), model.ErrCodeValidationError, model.ErrMessageValidationError+"("+err.Error()+")", "transaction")
	//	render.JSON(w, res, status)
	//	return
	//}

	//Token Authentocation
	a := AuthToken(w, r)
	if !a.Valid {
		render.JSON(w, a.res, status)
		return
	}

	//check redeemtion method
	switch rd.RedeemMethod {
	case model.RedeemtionMethodQr:
		//to-do validate partner_id
		par := map[string]string{"variant_id": rd.VariantID, "id": rd.Partner}
		if _, err := model.FindVariantPartner(par); err == model.ErrResourceNotFound {
			status = http.StatusBadRequest
			res.AddError(its(status), model.ErrCodeResourceNotFound, model.ErrMessageInvalidQr, "partner")
			render.JSON(w, res, status)
			return
		} else if err != nil {
			status = http.StatusInternalServerError
			res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", "partner")
			render.JSON(w, res, status)
			return
		}
	case model.RedeemtionMethodToken:
		//to-do validate token
		par := map[string]string{"variant_id": rd.VariantID, "id": rd.Partner}
		if p, err := model.FindVariantPartner(par); err == model.ErrResourceNotFound {
			status = http.StatusBadRequest
			res.AddError(its(status), model.ErrCodeResourceNotFound, model.ErrMessageInvalidPaerner, "partner")
			render.JSON(w, res, status)
			return
		} else if err != nil {
			status = http.StatusInternalServerError
			res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", "partner")
			render.JSON(w, res, status)
			return
		} else {
			fmt.Println("panrner data : ", p[0].SerialNumber.String)

			if !OTPAuth(p[0].SerialNumber.String, rd.Challenge, rd.Response) {
				status = http.StatusBadRequest
				res.AddError(its(status), model.ErrCodeOTPFailed, model.ErrMessageOTPFailed, "transaction")
				render.JSON(w, res, status)
				return
			}
		}

	}

	if ok, err := CheckVariant(rd.RedeemMethod, rd.VariantID, len(rd.Vouchers)); !ok {
		switch err.Error() {
		case model.ErrCodeAllowAccumulativeDisable:
			status = http.StatusBadRequest
			res.AddError(its(status), err.Error(), model.ErrMessageAllowAccumulativeDisable, "variant")
			render.JSON(w, res, status)
		case model.ErrCodeInvalidRedeemMethod:
			status = http.StatusBadRequest
			res.AddError(its(status), err.Error(), model.ErrMessageInvalidRedeemMethod, "variant")
			render.JSON(w, res, status)
		case model.ErrCodeVoucherNotActive:
			status = http.StatusBadRequest
			res.AddError(its(status), err.Error(), model.ErrMessageVoucherNotActive, "variant")
			render.JSON(w, res, status)
		case model.ErrResourceNotFound.Error():
			status = http.StatusBadRequest
			res.AddError(its(status), model.ErrCodeResourceNotFound, model.ErrMessageInvalidVariant, "variant")
			render.JSON(w, res, status)
		case model.ErrCodeRedeemNotValidDay:
			status = http.StatusBadRequest
			res.AddError(its(status), err.Error(), model.ErrMessageRedeemNotValidDay, "variant")
			render.JSON(w, res, status)
		case model.ErrCodeRedeemNotValidHour:
			status = http.StatusBadRequest
			res.AddError(its(status), err.Error(), model.ErrMessageRedeemNotValidHour, "variant")
			render.JSON(w, res, status)

		default:
			status = http.StatusInternalServerError
			res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", "variant")
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
				res.AddError(its(status), err.Error(), model.ErrMessageVoucherNotActive, "transaction")
				render.JSON(w, res, status)
			case model.ErrResourceNotFound.Error():
				status = http.StatusBadRequest
				res.AddError(its(status), model.ErrCodeResourceNotFound, err.Error(), "transaction")
				render.JSON(w, res, status)
			case model.ErrMessageVoucherAlreadyUsed:
				status = http.StatusBadRequest
				res.AddError(its(status), model.ErrCodeVoucherDisabled, err.Error(), "transaction")
				render.JSON(w, res, status)
			case model.ErrMessageVoucherAlreadyPaid:
				status = http.StatusBadRequest
				res.AddError(its(status), model.ErrCodeVoucherDisabled, model.ErrMessageVoucherAlreadyUsed, "transaction")
				render.JSON(w, res, status)
			case model.ErrMessageVoucherExpired:
				status = http.StatusBadRequest
				res.AddError(its(status), model.ErrCodeVoucherExpired, err.Error(), "transaction")
				render.JSON(w, res, status)
			default:
				status = http.StatusInternalServerError
				res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", "transaction")
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
	fmt.Println(d)
	if err := model.InsertTransaction(d); err != nil {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", "transaction")
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
		res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", "transaction")
		render.JSON(w, res, status)
		return
	} else if err != nil {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", "transaction")
		render.JSON(w, res, status)
		return
	}

	res = NewResponse(TransactionResponse{TransactionCode: txCode})
	render.JSON(w, res, status)
}

func WebCreateTransaction(w http.ResponseWriter, r *http.Request) {
	var rd TransactionRequest
	status := http.StatusCreated
	res := NewResponse(nil)

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		status = http.StatusBadRequest
		res.AddError(its(status), http.StatusText(status), http.StatusText(status)+"("+err.Error()+")", "transaction")
		render.JSON(w, res, status)
		return
	}

	err := rd.validate()
	if err != nil {
		status = http.StatusBadRequest
		res.AddError(its(status), model.ErrCodeValidationError, model.ErrMessageValidationError+"("+err.Error()+")", "transaction")
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
			res.AddError(its(status), model.ErrCodeResourceNotFound, model.ErrMessageInvalidQr, "partner")
			render.JSON(w, res, status)
			return
		} else if err != nil {
			status = http.StatusInternalServerError
			res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", "partner")
			render.JSON(w, res, status)
			return
		}
	case model.RedeemtionMethodToken:
		//to-do validate token
		par := map[string]string{"variant_id": rd.VariantID, "id": rd.Partner}
		if p, err := model.FindVariantPartner(par); err == model.ErrResourceNotFound {
			status = http.StatusBadRequest
			res.AddError(its(status), model.ErrCodeResourceNotFound, model.ErrMessageInvalidQr, "partner")
			render.JSON(w, res, status)
			return
		} else if err != nil {
			status = http.StatusInternalServerError
			res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", "partner")
			render.JSON(w, res, status)
			return
		} else {
			fmt.Println("panrner data : ", p[0].SerialNumber.String)
			if !OTPAuth(p[0].SerialNumber.String, rd.Challenge, rd.Response) {
				status = http.StatusBadRequest
				res.AddError(its(status), model.ErrCodeOTPFailed, model.ErrMessageOTPFailed, "transaction")
				render.JSON(w, res, status)
				return
			}
		}

	}

	if ok, err := CheckVariant(rd.RedeemMethod, rd.VariantID, len(rd.Vouchers)); !ok {
		switch err.Error() {
		case model.ErrCodeAllowAccumulativeDisable:
			status = http.StatusBadRequest
			res.AddError(its(status), err.Error(), model.ErrMessageAllowAccumulativeDisable, "variant")
			render.JSON(w, res, status)
		case model.ErrCodeInvalidRedeemMethod:
			status = http.StatusBadRequest
			res.AddError(its(status), err.Error(), model.ErrMessageInvalidRedeemMethod, "variant")
			render.JSON(w, res, status)
		case model.ErrCodeVoucherNotActive:
			status = http.StatusBadRequest
			res.AddError(its(status), err.Error(), model.ErrMessageVoucherNotActive, "variant")
			render.JSON(w, res, status)
		case model.ErrResourceNotFound.Error():
			status = http.StatusBadRequest
			res.AddError(its(status), model.ErrCodeResourceNotFound, model.ErrMessageInvalidVariant, "variant")
			render.JSON(w, res, status)
		case model.ErrCodeRedeemNotValidDay:
			status = http.StatusBadRequest
			res.AddError(its(status), err.Error(), model.ErrMessageRedeemNotValidDay, "variant")
			render.JSON(w, res, status)
		case model.ErrCodeRedeemNotValidHour:
			status = http.StatusBadRequest
			res.AddError(its(status), err.Error(), model.ErrMessageRedeemNotValidHour, "variant")
			render.JSON(w, res, status)

		default:
			status = http.StatusInternalServerError
			res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", "variant")
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
				res.AddError(its(status), err.Error(), model.ErrMessageVoucherNotActive, "transaction")
				render.JSON(w, res, status)
			case model.ErrResourceNotFound.Error():
				status = http.StatusBadRequest
				res.AddError(its(status), model.ErrCodeResourceNotFound, err.Error(), "transaction")
				render.JSON(w, res, status)
			case model.ErrMessageVoucherAlreadyUsed:
				status = http.StatusBadRequest
				res.AddError(its(status), model.ErrCodeVoucherDisabled, err.Error(), "transaction")
				render.JSON(w, res, status)
			case model.ErrMessageVoucherAlreadyPaid:
				status = http.StatusBadRequest
				res.AddError(its(status), model.ErrCodeVoucherDisabled, model.ErrMessageVoucherAlreadyUsed, "transaction")
				render.JSON(w, res, status)
			case model.ErrMessageVoucherExpired:
				status = http.StatusBadRequest
				res.AddError(its(status), model.ErrCodeVoucherExpired, err.Error(), "transaction")
				render.JSON(w, res, status)
			default:
				status = http.StatusInternalServerError
				res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", "transaction")
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
		res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", "transaction")
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
		res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", "transaction")
		render.JSON(w, res, status)
		return
	} else if err != nil {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", "transaction")
		render.JSON(w, res, status)
		return
	}

	res = NewResponse(TransactionResponse{TransactionCode: txCode})
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

				res.AddError(its(status), errTitle, err.Error(), "Get Transaction")
			} else {
				res = NewResponse(variant)
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
				res.AddError(its(status), errTitle, err.Error(), "Cashout Transaction")
			} else {
				res = NewResponse("Transaction Success")
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
		}
	} else {
		res = a.res
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
