package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-zoo/bone"
	"github.com/ruizu/render"

	"github.com/evoucher/voucher/internal/model"
)

type (
	TransactionRequest struct {
		CompanyID        string   `json:"companyId"`
		MerchantID       string   `json:"merchantId"`
		TransactionCode  string   `json:"transactionCode"`
		TotalTransaction float64  `json:"totalTransaction"`
		DiscountValue    float64  `json:"discountValue"`
		PaymentType      string   `json:"paymentType"`
		User             string   `json:"createdBy"`
		Vouchers         []string `json:"vouchers"`
	}
	DeleteRequest struct {
		User string `json:"createdBy"`
	}
)

func CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var rd TransactionRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		log.Panic(err)
	}

	d := &model.Transaction{
		CompanyID:        rd.CompanyID,
		PointNeeded:      rd.MerchantID,
		TransactionCode:  rd.TransactionCode,
		TotalTransaction: rd.TotalTransaction,
		DiscountValue:    rd.DiscountValue,
		PaymentType:      rd.PaymentType,
		User:             rd.User,
		Vouchers:         rd.Vouchers,
	}
	if err := d.Insert(); err != nil {
		log.Panic(err)
	}

	res := NewResponse(nil)
	render.JSON(w, res, http.StatusCreated)
}
