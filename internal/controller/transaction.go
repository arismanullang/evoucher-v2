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
		CompanyID        string   `json:"companyId"`
		MerchantID       string   `json:"merchantId"`
		TransactionCode  string   `json:"transactionCode"`
		TotalTransaction float64  `json:"totalTransaction"`
		DiscountValue    float64  `json:"discountValue"`
		PaymentType      string   `json:"paymentType"`
		User             string   `json:"createdBy"`
		Vouchers         []string `json:"vouchers"`
	}
	DeleteTransactionRequest struct {
		User string `json:"requestedBy"`
	}
	DateTransactionRequest struct {
		Start string `json:"start"`
		End   string `json:"end"`
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
		MerchantID:       rd.MerchantID,
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
		ID:               id,
		CompanyID:        rd.CompanyID,
		MerchantID:       rd.MerchantID,
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
		ID:   id,
		User: rd.User,
	}
	if err := d.Delete(); err != nil {
		log.Panic(err)
	}

	res := NewResponse(nil)
	render.JSON(w, res, http.StatusOK)
}
