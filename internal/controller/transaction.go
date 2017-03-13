package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-zoo/bone"
	"github.com/ruizu/render"

	"github.com/gilkor/evoucher/internal/model"
)

type (
	TransactionRequest struct {
		AccountId        string   `json:"account_id"`
		PartnerId        string   `json:"partner_id"`
		TransactionCode  string   `json:"transaction_code"`
		TotalTransaction float64  `json:"total_transaction"`
		DiscountValue    float64  `json:"discount_value"`
		PaymentType      string   `json:"payment_type"`
		User             string   `json:"created_by"`
		Vouchers         []string `json:"vouchers"`
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

func CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var rd TransactionRequest
	var tr VoucherResponse

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		log.Panic(err)
	}

	status := http.StatusCreated
	if _, ok := basicAuth(w, r); ok {

		d := model.Transaction{
			AccountId:        rd.AccountId,
			PartnerId:        rd.PartnerId,
			TransactionCode:  randStr(12, model.NUMERALS),
			TotalTransaction: rd.TotalTransaction,
			DiscountValue:    rd.DiscountValue,
			PaymentType:      rd.PaymentType,
			User:             rd.User,
		}
		for i, _ := range rd.Vouchers {
			rd := RedeemVoucherRequest{
				VoucherCode: rd.Vouchers[i],
				AccountID:   rd.AccountId,
			}
			if vr := rd.RedeemVoucherValidation(); vr.State != model.ResponseStateOk {
				break
			}
			//append(d.Vouchers,vr.)
		}

		if err := model.InsertTransaction(d); err != nil {
			log.Panic(err)
		}

	} else {
		status = http.StatusUnauthorized
	}

	res := NewResponse(tr)
	render.JSON(w, res, status)
}

func GetTransactionDetails(w http.ResponseWriter, r *http.Request) {
	id := bone.GetValue(r, "id")
	variant, err := model.FindTransactionByID(id)
	if err != nil && err != model.ErrResourceNotFound {
		log.Panic(err)
	}

	res := NewResponse(variant)
	render.JSON(w, res)
}

func GetTransactionByDate(w http.ResponseWriter, r *http.Request) {
	var rd DateTransactionRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		log.Panic(err)
	}

	variant, err := model.FindTransactionByDate(rd.Start, rd.End)
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

	d := &model.Transaction{
		Id:               id,
		AccountId:        rd.AccountId,
		PartnerId:        rd.PartnerId,
		TransactionCode:  rd.TransactionCode,
		TotalTransaction: rd.TotalTransaction,
		DiscountValue:    rd.DiscountValue,
		PaymentType:      rd.PaymentType,
		User:             rd.User,
		Vouchers:         rd.Vouchers,
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
