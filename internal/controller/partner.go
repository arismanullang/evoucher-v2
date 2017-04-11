package controller

import (
	"database/sql"
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
		ID           string `json:"id"`
		PartnerName  string `json:"partner_name"`
		SerialNumber string `json:"serial_number"`
	}
	PartnerResponseDetails []PartnerResponse
	PartnerResponse        struct {
		PartnerName  string `json:"partner_name"`
		SerialNumber string `json:"serial_number"`
		VariantID    string `json:"variant_id"`
		CreatedBy    string `json:"created_by"`
	}
)

func GetAllPartners(w http.ResponseWriter, r *http.Request) {
	status := http.StatusOK
	res := NewResponse(nil)
	partner, err := model.FindAllPartners()
	if err != nil {
		fmt.Println(err.Error())
		status = http.StatusInternalServerError
		errorTitle := model.ErrCodeInternalError
		if err == model.ErrResourceNotFound {
			status = http.StatusNotFound
			errorTitle = model.ErrCodeResourceNotFound
		}

		res.AddError(its(status), errorTitle, err.Error(), "Get Partner")
	} else {
		res = NewResponse(partner)
	}

	render.JSON(w, res, status)
}

func GetAllPartnersCustomParam(w http.ResponseWriter, r *http.Request) {

	res := NewResponse(nil)
	var status int
	var err error

	//Token Authentocation
	accountID, userID, _, ok := AuthToken(w, r)
	if !ok {
		return
	}
	fmt.Println(accountID, userID)

	param := getUrlParam(r.URL.String())
	delete(param, "token")

	p := []model.Partner{}
	if len(param) > 0 {
		p, err = model.FindVariantPartner(param)
	} else {
		status = http.StatusBadRequest
		res.AddError(its(status), model.ErrCodeMissingOrderItem, model.ErrMessageMissingOrderItem, "partner")
		render.JSON(w, res, status)
		return
	}

	if err == model.ErrResourceNotFound {
		status = http.StatusNotFound
		res.AddError(its(status), model.ErrCodeResourceNotFound, model.ErrMessageResourceNotFound, "partner")
		render.JSON(w, res, status)
		return
	} else if err != nil {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", "partner")
		render.JSON(w, res, status)
		return
	}

	d := make(PartnerResponseDetails, len(p))
	for i, v := range p {
		d[i].PartnerName = v.PartnerName
		d[i].SerialNumber = v.SerialNumber.String
		d[i].VariantID = v.VariantID
		d[i].CreatedBy = v.CreatedBy.String
	}

	status = http.StatusOK
	res = NewResponse(d)
	render.JSON(w, res, status)
}

func GetPartnerSerialName(w http.ResponseWriter, r *http.Request) {
	param := r.FormValue("param")

	status := http.StatusUnauthorized
	err := model.ErrTokenNotFound
	errorTitle := model.ErrCodeInvalidToken
	res := NewResponse(nil)
	res.AddError(its(status), errorTitle, err.Error(), "Get Partner")

	_, _, _, valid := AuthToken(w, r)
	if valid {
		status = http.StatusOK
		partner, err := model.FindPartnerSerialNumber(param)
		if err != nil {
			status = http.StatusInternalServerError
			errorTitle = model.ErrCodeInternalError
			if err != model.ErrResourceNotFound {
				status = http.StatusNotFound
				errorTitle = model.ErrCodeResourceNotFound
			}

			res.AddError(its(status), errorTitle, err.Error(), "Get Partner")
		} else {
			res = NewResponse(partner)
		}
	}
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

	status := http.StatusUnauthorized
	err := model.ErrTokenNotFound
	errorTitle := model.ErrCodeInvalidToken
	res := NewResponse(nil)
	res.AddError(its(status), errorTitle, err.Error(), "Add Partner")

	fmt.Println("Check Session")
	_, user, _, valid := AuthToken(w, r)
	if valid {
		status = http.StatusCreated
		param := model.Partner{
			PartnerName: rd.PartnerName,
			SerialNumber: sql.NullString{
				String: rd.SerialNumber,
				Valid:  true,
			},
			CreatedBy: sql.NullString{
				String: user,
				Valid:  true,
			},
		}
		err := model.InsertPartner(param)
		if err != nil {
			status = http.StatusInternalServerError
			errorTitle = model.ErrCodeInternalError
			if err != model.ErrResourceNotFound {
				status = http.StatusNotFound
				errorTitle = model.ErrCodeResourceNotFound
			}

			res.AddError(its(status), errorTitle, err.Error(), "Add Partner")
		}
	}
	render.JSON(w, res, status)
}
