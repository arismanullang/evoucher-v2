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
	VariantRequest struct {
		CompanyID         string   `json:"companyId"`
		VariantName       string   `json:"variantName"`
		VariantType       string   `json:"variantType"`
		PointNeeded       float64  `json:"pointNeeded"`
		MaxVoucher        float64  `json:"maxVoucher"`
		AllowAccumulative bool     `json:"allowAccumulative"`
		StartDate         string   `json:"startDate"`
		FinishDate        string   `json:"finishDate"`
		DiscountValue     float64  `json:"discountValue"`
		ImgUrl            string   `json:"imgUrl"`
		VariantTnc        string   `json:"variantTnc"`
		User              string   `json:"createdBy"`
		ValidUsers        []string `json:"validUsers"`
	}
	DeleteVariantRequest struct {
		User string `json:"requestedBy"`
	}
	SearchVariantRequest struct {
		Field string `json:"field"`
		Value string `json:"value"`
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

	tf, err := time.Parse("01/02/2006", rd.FinishDate)
	if err != nil {
		log.Panic(err)
	}

	d := &model.Variant{
		CompanyID:         rd.CompanyID,
		PointNeeded:       rd.PointNeeded,
		VariantName:       rd.VariantName,
		VariantType:       rd.VariantType,
		MaxVoucher:        rd.MaxVoucher,
		AllowAccumulative: rd.AllowAccumulative,
		DiscountValue:     rd.DiscountValue,
		StartDate:         ts.Format("2006-01-02 15:04:05.000"),
		FinishDate:        tf.Format("2006-01-02 15:04:05.000"),
		ImgUrl:            rd.ImgUrl,
		VariantTnc:        rd.VariantTnc,
		User:              rd.User,
		ValidUsers:        rd.ValidUsers,
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

func SearchVariant(w http.ResponseWriter, r *http.Request) {
	var rd SearchVariantRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		log.Panic(err)
	}

	var variant model.VariantsResponse
	var err error

	switch rd.Field {
	case "variant_name":
		variant, err = model.FindVariantByName("%" + rd.Value + "%")
	case "company_id":
		variant, err = model.FindVariantByCompanyID(rd.Value)
	case "date":
		variant, err = model.FindVariantByDate(rd.Value)
	}

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
	//id := bone.GetValue(r, "id")
	id := "nZ9Xmo-2"
	variant, err := model.FindVariantByUser(id)
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
		ID:                id,
		CompanyID:         rd.CompanyID,
		PointNeeded:       rd.PointNeeded,
		VariantName:       rd.VariantName,
		VariantType:       rd.VariantType,
		MaxVoucher:        rd.MaxVoucher,
		AllowAccumulative: rd.AllowAccumulative,
		StartDate:         rd.StartDate,
		FinishDate:        rd.FinishDate,
		DiscountValue:     rd.DiscountValue,
		ImgUrl:            rd.ImgUrl,
		VariantTnc:        rd.VariantTnc,
		User:              rd.User,
		ValidUsers:        rd.ValidUsers,
	}
	if err := d.Update(); err != nil {
		log.Panic(err)
	}

	res := NewResponse(nil)
	render.JSON(w, res, http.StatusOK)
}

func DeleteVariant(w http.ResponseWriter, r *http.Request) {
	id := bone.GetValue(r, "id")
	var rd DeleteVariantRequest
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
