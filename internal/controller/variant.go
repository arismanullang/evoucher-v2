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
		CompanyID          string   `json:"company_id"`
		VariantName        string   `json:"variant_name"`
		VariantType        string   `json:"variant_type"`
		VoucherType        string   `json:"voucher_type"`
		PointNeeded        float64  `json:"point_needed"`
		MaxQuantityVoucher float64  `json:"max_quantity"`
		MaxUsageVoucher    float64  `json:"max_usage"`
		AllowAccumulative  bool     `json:"allow_accumulative"`
		RedeemtionMethod   string   `json:"redeem"`
		StartDate          string   `json:"start_date"`
		EndDate            string   `json:"end_date"`
		DiscountValue      float64  `json:"discount_value"`
		ImgUrl             string   `json:"img_url"`
		VariantTnc         string   `json:"variant_tnc"`
		User               string   `json:"created_by"`
		BlastUsers         []string `json:"blast_users"`
		ValidTenants       []string `json:"valid_tenants"`
	}
	UserVariantRequest struct {
		User string `json:"user"`
	}
	DateVariantRequest struct {
		Start string `json:"start"`
		End   string `json:"end"`
	}
	MultiUserVariantRequest struct {
		CompanyID string   `json:"company_id"`
		User      string   `json:"user"`
		Data      []string `json:"data"`
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

	tf, err := time.Parse("01/02/2006", rd.EndDate)
	if err != nil {
		log.Panic(err)
	}

	d := &model.Variant{
		CompanyID:          rd.CompanyID,
		VariantName:        rd.VariantName,
		VariantType:        rd.VariantType,
		VoucherType:        rd.VoucherType,
		PointNeeded:        rd.PointNeeded,
		MaxQuantityVoucher: rd.MaxQuantityVoucher,
		MaxUsageVoucher:    rd.MaxUsageVoucher,
		AllowAccumulative:  rd.AllowAccumulative,
		RedeemtionMethod:   rd.RedeemtionMethod,
		DiscountValue:      rd.DiscountValue,
		StartDate:          ts.Format("2006-01-02 15:04:05.000"),
		EndDate:            tf.Format("2006-01-02 15:04:05.000"),
		ImgUrl:             rd.ImgUrl,
		VariantTnc:         rd.VariantTnc,
		User:               rd.User,
		BlastUsers:         rd.BlastUsers,
		ValidTenants:       rd.ValidTenants,
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

func GetVariantDetailsByID(w http.ResponseWriter, r *http.Request) {
	id := bone.GetValue(r, "id")
	variant, err := model.FindVariantByID(id)
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
		ID:                 id,
		CompanyID:          rd.CompanyID,
		VariantName:        rd.VariantName,
		VariantType:        rd.VariantType,
		VoucherType:        rd.VoucherType,
		PointNeeded:        rd.PointNeeded,
		MaxQuantityVoucher: rd.MaxQuantityVoucher,
		MaxUsageVoucher:    rd.MaxUsageVoucher,
		AllowAccumulative:  rd.AllowAccumulative,
		RedeemtionMethod:   rd.RedeemtionMethod,
		DiscountValue:      rd.DiscountValue,
		StartDate:          rd.StartDate,
		EndDate:            rd.EndDate,
		ImgUrl:             rd.ImgUrl,
		VariantTnc:         rd.VariantTnc,
		User:               rd.User,
		BlastUsers:         rd.BlastUsers,
		ValidTenants:       rd.ValidTenants,
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
		ID:        id,
		CompanyID: rd.CompanyID,
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
		ID:        id,
		CompanyID: rd.CompanyID,
		User:      rd.User,
		Data:      rd.Data,
	}

	if err := model.UpdateTenant(d); err != nil {
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
		ID:   id,
		User: rd.User,
	}
	if err := d.Delete(); err != nil {
		log.Panic(err)
	}

	res := NewResponse(nil)
	render.JSON(w, res, http.StatusOK)
}
