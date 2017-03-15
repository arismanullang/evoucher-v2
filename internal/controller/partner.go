package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	//"time"

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
	status := http.StatusOK
	partner, err := model.FindAllPartner()
	if err != nil && err != model.ErrResourceNotFound {
		//log.Panic(err)
		status = http.StatusInternalServerError
	}

	res := NewResponse(partner)
	render.JSON(w, res, status)
}

func GetPartnerSerialName(w http.ResponseWriter, r *http.Request) {
	param := r.FormValue("param")

	var partner = model.Response{}
	var err error
	status := http.StatusUnauthorized
	token := r.FormValue("token")
	valid := false

	if token != "" && token != "null" {
		_, _, _, valid, _ = getValiditySession(r, token)
	}

	if valid {
		status = http.StatusOK
		partner, err = model.FindPartnerSerialNumber(param)
		if err != nil && err != model.ErrResourceNotFound {
			log.Panic(err)
		}
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
	user := r.FormValue("user")
	valid := false
	status := http.StatusUnauthorized
	if token != "" && token != "null" {
		_, _, _, valid, _ = getValiditySession(r, token)
	}

	if valid {
		status = http.StatusCreated
		param := model.Partner{
			PartnerName:  rd.PartnerName,
			SerialNumber: rd.SerialNumber,
			CreatedBy:    user,
		}

		if err := model.InsertPartner(param); err != nil {
			//log.Panic(err)
			status = http.StatusInternalServerError
		}
	}

	res := NewResponse(nil)
	render.JSON(w, res, status)
}
