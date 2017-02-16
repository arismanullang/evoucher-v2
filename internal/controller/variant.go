package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
	//"fmt"

	"github.com/go-zoo/bone"
	"github.com/ruizu/render"

	"github.com/gilkor/evoucher/internal/model"
)

type (
	VariantRequest struct {
		AccountId          string   `json:"account_id"`
		VariantName        string   `json:"variant_name"`
		VariantType        string   `json:"variant_type"`
		VoucherFormatId    string   `json:"voucher_format_id"`
		VoucherType        string   `json:"voucher_type"`
		VoucherPrice       float64  `json:"voucher_price"`
		AllowAccumulative  bool     `json:"allow_accumulative"`
		StartDate          string   `json:"start_date"`
		EndDate            string   `json:"end_date"`
		DiscountValue      float64  `json:"discount_value"`
		MaxQuantityVoucher float64  `json:"max_quantity_voucher"`
		MaxUsageVoucher    float64  `json:"max_usage_voucher"`
		RedeemtionMethod   string   `json:"redeem_method"`
		ImgUrl             string   `json:"img_url"`
		VariantTnc         string   `json:"variant_tnc"`
		VariantDescription string   `json:"variant_description"`
		User               string   `json:"created_by"`
		ValidPartners      []string `json:"valid_partners"`
	}
	UserVariantRequest struct {
		User string `json:"user"`
	}
	DateVariantRequest struct {
		Start string `json:"start"`
		End   string `json:"end"`
	}
	MultiUserVariantRequest struct {
		User string   `json:"user"`
		Data []string `json:"data"`
	}
	SearchVariantRequests struct {
		Fields []string `json:"fields"`
		Values []string `json:"values"`
	}
)

func CreateVariant(w http.ResponseWriter, r *http.Request) {
	var rd VariantRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		log.Panic(err)
	}

	ts, err := time.Parse("01/02/2006", rd.StartDate)
	if err != nil {
		log.Panic(err)
	}

	d := &model.Variant{
		AccountId:          rd.AccountId,
		VariantName:        rd.VariantName,
		VariantType:        rd.VariantType,
		VoucherFormatId:    rd.VoucherFormatId,
		VoucherType:        rd.VoucherType,
		VoucherPrice:       rd.VoucherPrice,
		MaxQuantityVoucher: rd.MaxQuantityVoucher,
		MaxUsageVoucher:    rd.MaxUsageVoucher,
		AllowAccumulative:  rd.AllowAccumulative,
		RedeemtionMethod:   rd.RedeemtionMethod,
		DiscountValue:      rd.DiscountValue,
		StartDate:          ts.Format("2006-01-02 15:04:05.000"),
		EndDate:            ts.Format("2006-01-02 15:04:05.000"),
		ImgUrl:             rd.ImgUrl,
		VariantTnc:         rd.VariantTnc,
		VariantDescription: rd.VariantDescription,
		CreatedBy:          rd.User,
		ValidPartners:      rd.ValidPartners,
	}
	if err := d.Insert(); err != nil {
		log.Panic(err)
	}

	res := NewResponse(nil)
	render.JSON(w, res, http.StatusCreated)
}

func GetAllVariant(w http.ResponseWriter, r *http.Request) {
	variant, err := model.FindAllVariants()
	if err != nil && err != model.ErrResourceNotFound {
		log.Panic(err)
	}

	res := NewResponse(variant)
	render.JSON(w, res)
}

func GetVariantDetails(w http.ResponseWriter, r *http.Request) {
	var rd SearchVariantRequests
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		log.Panic(err)
	}

	variant, err := model.FindVariantMultipleParam(rd.Fields, rd.Values)
	if err != nil && err != model.ErrResourceNotFound {
		log.Panic(err)
	}

	res := NewResponse(variant)
	render.JSON(w, res)
}

func GetVariantDetailsById(w http.ResponseWriter, r *http.Request) {
	id := bone.GetValue(r, "id")
	variant, err := model.FindVariantById(id)
	if err != nil && err != model.ErrResourceNotFound {
		log.Panic(err)
	}

	res := NewResponse(variant)
	render.JSON(w, res)
}

func GetVariantDetailsByUser(w http.ResponseWriter, r *http.Request) {
	//userId := "nZ9Xmo-2"
	var rd UserVariantRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		log.Panic(err)
	}

	field := []string{"created_by"}
	value := []string{rd.User}

	param := SearchVariantRequests{
		Fields: field,
		Values: value,
	}

	variant, err := model.FindVariantMultipleParam(param.Fields, param.Values)
	if err != nil && err != model.ErrResourceNotFound {
		log.Panic(err)
	}

	res := NewResponse(variant)
	render.JSON(w, res)
}

func GetVariantDetailsByDate(w http.ResponseWriter, r *http.Request) {
	var rd DateVariantRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		log.Panic(err)
	}

	variant, err := model.FindVariantByDate(rd.Start, rd.End)
	if err != nil && err != model.ErrResourceNotFound {
		log.Panic(err)
	}

	res := NewResponse(variant)
	render.JSON(w, res)
}

func UpdateVariant(w http.ResponseWriter, r *http.Request) {
	id := bone.GetValue(r, "id")
	var rd VariantRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		log.Panic(err)
	}

	d := &model.Variant{
		Id:                 id,
		AccountId:          rd.AccountId,
		VariantName:        rd.VariantName,
		VariantType:        rd.VariantType,
		VoucherFormatId:    rd.VoucherFormatId,
		VoucherType:        rd.VoucherType,
		VoucherPrice:       rd.VoucherPrice,
		MaxQuantityVoucher: rd.MaxQuantityVoucher,
		MaxUsageVoucher:    rd.MaxUsageVoucher,
		AllowAccumulative:  rd.AllowAccumulative,
		RedeemtionMethod:   rd.RedeemtionMethod,
		DiscountValue:      rd.DiscountValue,
		StartDate:          rd.StartDate,
		EndDate:            rd.EndDate,
		ImgUrl:             rd.ImgUrl,
		VariantTnc:         rd.VariantTnc,
		VariantDescription: rd.VariantDescription,
		ValidPartners:      rd.ValidPartners,
	}
	if err := d.Update(); err != nil {
		log.Panic(err)
	}

	res := NewResponse(nil)
	render.JSON(w, res, http.StatusOK)
}

func UpdateVariantBroadcast(w http.ResponseWriter, r *http.Request) {
	id := bone.GetValue(r, "id")
	var rd MultiUserVariantRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		log.Panic(err)
	}

	d := model.UpdateVariantUsersRequest{
		VariantId: id,
		User:      rd.User,
		Data:      rd.Data,
	}

	if err := model.UpdateBroadcast(d); err != nil {
		log.Panic(err)
	}

	res := NewResponse(nil)
	render.JSON(w, res)
}

func UpdateVariantTenant(w http.ResponseWriter, r *http.Request) {
	id := bone.GetValue(r, "id")
	var rd MultiUserVariantRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		log.Panic(err)
	}

	d := model.UpdateVariantUsersRequest{
		VariantId: id,
		User:      rd.User,
		Data:      rd.Data,
	}

	if err := model.UpdatePartner(d); err != nil {
		log.Panic(err)
	}

	res := NewResponse(nil)
	render.JSON(w, res)
}

func DeleteVariant(w http.ResponseWriter, r *http.Request) {
	id := bone.GetValue(r, "id")
	var rd UserVariantRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		log.Panic(err)
	}

	d := &model.DeleteVariantRequest{
		Id:   id,
		User: rd.User,
	}
	if err := d.Delete(); err != nil {
		log.Panic(err)
	}

	res := NewResponse(nil)
	render.JSON(w, res, http.StatusOK)
}
