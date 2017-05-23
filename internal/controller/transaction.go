package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-zoo/bone"
	"github.com/ruizu/render"

	"github.com/gilkor/evoucher/internal/model"
)

type (
	TransactionRequest struct {
		VariantID     string   `json:"variant_id"`
		RedeemMethod  string   `json:"redeem_method"`
		Partner       string   `json:"partner"`
		Challenge     string   `json:"challenge"`
		Response      string   `json:"response"`
		DiscountValue float64  `json:"discount_value"`
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
)

func MobileCreateTransaction(w http.ResponseWriter, r *http.Request) {
	var rd TransactionRequest
	status := http.StatusCreated
	res := NewResponse(nil)

	//Token Authentocation
	accountID, userID, _, ok := AuthToken(w, r)
	if !ok {
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		status = http.StatusBadRequest
		res.AddError(its(status), http.StatusText(status), http.StatusText(status)+"("+err.Error()+")", "voucher")
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
			if !OTPAuth(p[0].Id, rd.Challenge, rd.Response) {
				status = http.StatusBadRequest
				res.AddError(its(status), model.ErrCodeOTPFailed, model.ErrMessageOTPFailed, "voucher")
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
				res.AddError(its(status), err.Error(), model.ErrMessageVoucherNotActive, "voucher")
				render.JSON(w, res, status)
			case model.ErrResourceNotFound.Error():
				status = http.StatusBadRequest
				res.AddError(its(status), model.ErrCodeResourceNotFound, err.Error(), "voucher")
				render.JSON(w, res, status)
			case model.ErrMessageVoucherAlreadyUsed:
				status = http.StatusBadRequest
				res.AddError(its(status), model.ErrCodeVoucherDisabled, err.Error(), "voucher")
				render.JSON(w, res, status)
			case model.ErrMessageVoucherAlreadyPaid:
				status = http.StatusBadRequest
				res.AddError(its(status), model.ErrCodeVoucherDisabled, model.ErrMessageVoucherAlreadyUsed, "voucher")
				render.JSON(w, res, status)
			case model.ErrMessageVoucherExpired:
				status = http.StatusBadRequest
				res.AddError(its(status), model.ErrCodeVoucherExpired, err.Error(), "voucher")
				render.JSON(w, res, status)
			default:
				status = http.StatusInternalServerError
				res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", "voucher")
				render.JSON(w, res, status)
			}
			return
		}
	}

	txCode := randStr(model.DEFAULT_TXLENGTH, model.DEFAULT_TXCODE)
	d := model.Transaction{
		AccountId:       accountID,
		PartnerId:       rd.Partner,
		TransactionCode: txCode,
		DiscountValue:   rd.DiscountValue,
		Token:           rd.Response,
		User:            userID,
		Vouchers:        rd.Vouchers,
	}
	fmt.Println(d)
	if err := model.InsertTransaction(d); err != nil {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", "voucher")
		render.JSON(w, res, status)
		return
	}

	rv := RedeemVoucherRequest{
		AccountID: accountID,
		User:      userID,
		State:     model.VoucherStateUsed,
		Vouchers:  rd.Vouchers,
	}
	fmt.Println("List valid voucher :", rv.Vouchers)
	// update voucher state "Used"
	if ok, err := rv.UpdateVoucher(); !ok {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", "voucher")
		render.JSON(w, res, status)
		return
	} else if err != nil {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", "voucher")
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
		res.AddError(its(status), http.StatusText(status), http.StatusText(status)+"("+err.Error()+")", "voucher")
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
			if !OTPAuth(p[0].Id, rd.Challenge, rd.Response) {
				status = http.StatusBadRequest
				res.AddError(its(status), model.ErrCodeOTPFailed, model.ErrMessageOTPFailed, "voucher")
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
				res.AddError(its(status), err.Error(), model.ErrMessageVoucherNotActive, "voucher")
				render.JSON(w, res, status)
			case model.ErrResourceNotFound.Error():
				status = http.StatusBadRequest
				res.AddError(its(status), model.ErrCodeResourceNotFound, err.Error(), "voucher")
				render.JSON(w, res, status)
			case model.ErrMessageVoucherAlreadyUsed:
				status = http.StatusBadRequest
				res.AddError(its(status), model.ErrCodeVoucherDisabled, err.Error(), "voucher")
				render.JSON(w, res, status)
			case model.ErrMessageVoucherAlreadyPaid:
				status = http.StatusBadRequest
				res.AddError(its(status), model.ErrCodeVoucherDisabled, model.ErrMessageVoucherAlreadyUsed, "voucher")
				render.JSON(w, res, status)
			case model.ErrMessageVoucherExpired:
				status = http.StatusBadRequest
				res.AddError(its(status), model.ErrCodeVoucherExpired, err.Error(), "voucher")
				render.JSON(w, res, status)
			default:
				status = http.StatusInternalServerError
				res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", "voucher")
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
		DiscountValue:   rd.DiscountValue,
		Token:           rd.Response,
		User:            variant.CreatedBy,
		Vouchers:        rd.Vouchers,
	}
	fmt.Println(d)
	if err := model.InsertTransaction(d); err != nil {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", "voucher")
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
		res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", "voucher")
		render.JSON(w, res, status)
		return
	} else if err != nil {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", "voucher")
		render.JSON(w, res, status)
		return
	}

	res = NewResponse(TransactionResponse{TransactionCode: txCode})
	render.JSON(w, res, status)
}

func GetAllTransactions(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get Transaction")
	status := http.StatusUnauthorized
	err := model.ErrTokenNotFound
	errTitle := model.ErrCodeInvalidToken
	res := NewResponse(nil)
	res.AddError(its(status), errTitle, err.Error(), "Get Transaction")

	fmt.Println("Check Session")
	accountId, _, _, valid := AuthToken(w, r)
	if valid {
		status = http.StatusOK
		transaction, err := model.FindAllTransaction(accountId)
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
			res = NewResponse(transaction)
		}
	}
	render.JSON(w, res, status)
}

func GetAllTransactionsByPartner(w http.ResponseWriter, r *http.Request) {
	partnerId := r.FormValue("partner")

	fmt.Println("Get Transaction")
	status := http.StatusUnauthorized
	err := model.ErrTokenNotFound
	errTitle := model.ErrCodeInvalidToken
	res := NewResponse(nil)
	res.AddError(its(status), errTitle, err.Error(), "Get Transaction")

	fmt.Println("Check Session")
	accountId, _, _, valid := AuthToken(w, r)
	if valid {
		status = http.StatusOK
		transaction, _ := model.FindAllTransactionByPartner(accountId, partnerId)
		// transaction, err := model.FindAllTransactionByPartner(accountId, partnerId)
		fmt.Println(err)
		// if err != nil {
		// 	status = http.StatusInternalServerError
		// 	errTitle = model.ErrCodeInternalError
		// 	if err == model.ErrResourceNotFound {
		// 		status = http.StatusNotFound
		// 		errTitle = model.ErrCodeResourceNotFound
		// 	}
		//
		// 	res.AddError(its(status), errTitle, err.Error(), "Get Transaction")
		// } else {
		res = NewResponse(transaction)
		// }
	}
	render.JSON(w, res, status)
}

func GetTransaction(w http.ResponseWriter, r *http.Request) {
	status := http.StatusOK

	accountID, _, _, ok := AuthToken(w, r)
	if !ok {
		return
	}

	tx, err := model.FindTransactionDetailsById(accountID)
	if err != nil && err != model.ErrResourceNotFound {
		log.Panic(err)
	}

	res := NewResponse(tx)
	render.JSON(w, res, status)
}

func GetTransactionByDate(w http.ResponseWriter, r *http.Request) {
	var rd DateTransactionRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		log.Panic(err)
	}

	variant, err := model.FindTransactionDetailsByDate(rd.Start, rd.End)
	if err != nil && err != model.ErrResourceNotFound {
		log.Panic(err)
	}

	res := NewResponse(variant)
	render.JSON(w, res)
}

func UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	id := bone.GetValue(r, "id")
	var rd TransactionRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		log.Panic(err)
	}

	accountID, userID, _, ok := AuthToken(w, r)
	if !ok {
		return
	}

	d := &model.Transaction{
		Id:              id,
		AccountId:       accountID,
		PartnerId:       rd.Partner,
		TransactionCode: "",
		DiscountValue:   0,
		User:            userID,
		Vouchers:        rd.Vouchers,
	}
	if err := d.Update(); err != nil {
		log.Panic(err)
	}

	res := NewResponse(nil)
	render.JSON(w, res, http.StatusOK)
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
