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
		CompanyID         string    `json:"companyId"`
		VariantName       string    `json:"variantName"`
		VariantType       string    `json:"variantType"`
		PointNeeded       float64   `json:"pointNeeded"`
		MaxVoucher        float64   `json:"maxVoucher"`
		AllowAccumulative bool      `json:"allowAccumulative"`
		StartDate         time.Time `json:"startDate"`
		FinishDate        time.Time `json:"finishDate"`
		ImgUrl            string    `json:"imgUrl"`
		VariantTnc        string    `json:"variantTnc"`
		User              string    `json:"createdBy"`
		ValidUsers        []string  `json:"validUsers"`
	}
	DeleteRequest struct {
		User string `json:"createdBy"`
	}
)

func CreateVariant(w http.ResponseWriter, r *http.Request) {
	var rd VariantRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		log.Panic(err)
	}

	d := &model.Variant{
		CompanyID:         rd.CompanyID,
		PointNeeded:       rd.PointNeeded,
		VariantName:       rd.VariantName,
		VariantType:       rd.VariantType,
		MaxVoucher:        rd.MaxVoucher,
		AllowAccumulative: rd.AllowAccumulative,
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

func GetVariantDetails(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	name := r.FormValue("name")

	var variant model.VariantResponse
	var err error
	if id != "" && name != "" {
		w.Write([]byte("Please choose id or name"))
		return
	}

	if id != "" {
		variant, err = model.FindVariantByID(id)
	}
	if name != "" {
		variant, err = model.FindVariantByName(name)
	}

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
	var rd DeleteRequest
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
