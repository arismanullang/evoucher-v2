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
		CompanyID          string   `json:"companyId"`
		VariantName        string   `json:"variantName"`
		VariantType        string   `json:"variantType"`
		VoucherType        string   `json:"voucherType"`
		PointNeeded        float64  `json:"pointNeeded"`
		MaxQuantityVoucher float64  `json:"maxQuantity"`
		MaxUsageVoucher    float64  `json:"maxUsage"`
		AllowAccumulative  bool     `json:"allowAccumulative"`
		RedeemtionMethod   string   `json:"redeem"`
		StartDate          string   `json:"startDate"`
		EndDate            string   `json:"endDate"`
		DiscountValue      float64  `json:"discountValue"`
		ImgUrl             string   `json:"imgUrl"`
		VariantTnc         string   `json:"variantTnc"`
		User               string   `json:"createdBy"`
		BlastUsers         []string `json:"blastUsers"`
		ValidTenants       []string `json:"validTenants"`
	}
	UserVariantRequest struct {
		User string `json:"user"`
	}
	DateVariantRequest struct {
		Start string `json:"start"`
		End   string `json:"end"`
	}
	MultiUserVariantRequest struct {
		CompanyID string   `json:"companyId"`
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
