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
		VariantID        string   `json:"variant_id"`
		PartnerID        string   `json:"partner_id"`
		RedeemMethod     string   `json:"Redeem_method"`
		RedeemCode       string   `json:"Redeem_code"`
		TotalTransaction float64  `json:"total_transaction"`
		PaymentType      string   `json:"payment_type"`
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
	var tr ResponseData

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		log.Panic(err)
	}

	status := http.StatusCreated

	d := model.Transaction{
		AccountId:        "",
		PartnerId:        rd.PartnerID,
		TransactionCode:  randStr(12, model.NUMERALS),
		TotalTransaction: rd.TotalTransaction,
		DiscountValue:    0,
		PaymentType:      rd.PaymentType,
		User:             "",
	}
	for i, _ := range rd.Vouchers {
		rd := RedeemVoucherRequest{
			VoucherCode: rd.Vouchers[i],
			AccountID:   "",
		}
		if vr := rd.RedeemVoucherValidation(); vr.State != model.ResponseStateOk {
			break
		}
		//append(d.Vouchers,vr.)
	}

	if err := model.InsertTransaction(d); err != nil {
		log.Panic(err)
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
		AccountId:        "",
		PartnerId:        rd.PartnerID,
		TransactionCode:  "",
		TotalTransaction: rd.TotalTransaction,
		DiscountValue:    0,
		PaymentType:      rd.PaymentType,
		User:             "",
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
