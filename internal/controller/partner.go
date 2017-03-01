package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	//"github.com/go-zoo/bone"
	"github.com/ruizu/render"

	"github.com/gilkor/evoucher/internal/model"
)

type (
	PartnerReq struct {
		Partner   string `json:"partner"`
		CreatedBy string `json:"created_by"`
	}
	Partner struct {
		PartnerName  string `json:"partner_name"`
		SerialNumber string `json:"serial_number"`
	}
)

func GetAllPartner(w http.ResponseWriter, r *http.Request) {
	var partner = model.Response{}
	var err error
	var status int
	if basicAuth(w, r) {
		partner, err = model.FindAllPartner()
		if err != nil && err != model.ErrResourceNotFound {
			log.Panic(err)
		}
		status = http.StatusOK
	} else {
		partner = model.Response{}
		status = http.StatusUnauthorized
	}

	res := NewResponse(partner)
	render.JSON(w, res, status)
}

func GetPartnerSerialName(w http.ResponseWriter, r *http.Request) {
	param := r.FormValue("param")

	var partner = model.Response{}
	var err error
	var status int
	if basicAuth(w, r) {
		partner, err = model.FindPartnerSerialNumber(param)
		if err != nil && err != model.ErrResourceNotFound {
			log.Panic(err)
		}
		status = http.StatusOK
	} else {
		partner = model.Response{}
		status = http.StatusUnauthorized
	}

	res := NewResponse(partner)
	render.JSON(w, res, status)
}

// dashboard
func AddPartner(w http.ResponseWriter, r *http.Request) {
	fmt.Print("Add")
	var rd Partner
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		log.Panic(err)
	}

	token := r.FormValue("token")
	valid := false
	var id string
	var exp time.Time
	if token != "" && token != "null" {
		id, _, exp = checkExpired(r, token)
		if exp.After(time.Now()) {
			valid = true
		}
	}

	var status int
	if valid {
		param := model.Partner{
			PartnerName:  rd.PartnerName,
			SerialNumber: rd.SerialNumber,
			CreatedBy:    id,
		}

		if err := model.InsertPartner(param); err != nil {
			log.Panic(err)
		}
		status = http.StatusCreated
	} else {
		status = http.StatusUnauthorized
	}

	res := NewResponse(nil)
	render.JSON(w, res, status)
}

func DashboardGetAllPartner(w http.ResponseWriter, r *http.Request) {
	token := r.FormValue("token")
	valid := false
	var exp time.Time
	if token != "" && token != "null" {
		_, _, exp = checkExpired(r, token)
		if exp.After(time.Now()) {
			valid = true
		}
	}

	var partner = model.Response{}
	var err error
	var status int
	if valid {
		partner, err = model.FindAllPartner()
		if err != nil && err != model.ErrResourceNotFound {
			log.Panic(err)
		}
		status = http.StatusOK
	} else {
		partner = model.Response{}
		status = http.StatusUnauthorized
	}

	res := NewResponse(partner)
	render.JSON(w, res, status)
}
