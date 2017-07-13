package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-zoo/bone"
	"github.com/ruizu/render"

	"github.com/gilkor/evoucher/internal/model"
)

type (
	PartnerReq struct {
		Partner   string `json:logger.TraceID`
		CreatedBy string `json:"created_by"`
	}
	Partner struct {
		ID           string `json:"id"`
		PartnerName  string `json:"partner_name"`
		SerialNumber string `json:"serial_number"`
		Tag          string `json:"tag"`
		Description  string `json:"description"`
	}
	PartnerResponseDetails []PartnerResponse
	PartnerResponse        struct {
		PartnerName  string `json:"partner_name"`
		SerialNumber string `json:"serial_number"`
		VariantID    string `json:"variant_id"`
		CreatedBy    string `json:"created_by"`
	}
	Tag struct {
		Value string `json:"tag"`
	}
	Tags struct {
		Value []string `json:"tags"`
	}
)

func GetVariantPartners(w http.ResponseWriter, r *http.Request) {
	param := getUrlParam(r.URL.String())
	status := http.StatusOK
	res := NewResponse(nil)

	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag("Get Partner of Variant")

	partner, err := model.FindVariantPartner(param)
	if err != nil {
		fmt.Println(err.Error())
		status = http.StatusInternalServerError
		errorTitle := model.ErrCodeInternalError
		if err == model.ErrResourceNotFound {
			status = http.StatusNotFound
			errorTitle = model.ErrCodeResourceNotFound
		}

		res.AddError(its(status), errorTitle, err.Error(), logger.TraceID)
		logger.SetStatus(status).Log("param :", param , "response :" , res.Errors.ToString())
	} else {
		res = NewResponse(partner)
		logger.SetStatus(status).Log("param :", param , "response :" , partner)
	}

	render.JSON(w, res, status)
}

func GetPartners(w http.ResponseWriter, r *http.Request) {
	param := getUrlParam(r.URL.String())
	status := http.StatusOK
	res := NewResponse(nil)
	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag("Get Partner")

	partner, err := model.FindPartners(param)
	if err != nil {
		fmt.Println(err.Error())
		status = http.StatusInternalServerError
		errorTitle := model.ErrCodeInternalError
		if err == model.ErrResourceNotFound {
			status = http.StatusNotFound
			errorTitle = model.ErrCodeResourceNotFound
		}

		res.AddError(its(status), errorTitle, err.Error(), logger.TraceID)
		logger.SetStatus(status).Log("param :", param , "response :" , res.Errors.ToString())
	} else {
		res = NewResponse(partner)
		logger.SetStatus(status).Log("param :", param , "response :" , partner)
	}

	render.JSON(w, res, status)
}

func GetAllPartners(w http.ResponseWriter, r *http.Request) {
	status := http.StatusOK
	res := NewResponse(nil)

	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag("Get All Partner")

	partner, err := model.FindAllPartners()
	if err != nil {
		fmt.Println(err.Error())
		status = http.StatusInternalServerError
		errorTitle := model.ErrCodeInternalError
		if err == model.ErrResourceNotFound {
			status = http.StatusNotFound
			errorTitle = model.ErrCodeResourceNotFound
		}

		res.AddError(its(status), errorTitle, err.Error(), logger.TraceID)
		logger.SetStatus(status).Log("param :" , "response :" , res.Errors.ToString())
	} else {
		res = NewResponse(partner)
		logger.SetStatus(status).Log("param :" , "response :" , partner)
	}

	render.JSON(w, res, status)
}

func GetAllPartnersCustomParam(w http.ResponseWriter, r *http.Request) {

	res := NewResponse(nil)
	var status int
	var err error

	//Token Authentocation
	a := AuthToken(w, r)
	if !a.Valid {
		render.JSON(w, a.res, status)
		return
	}

	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag("Get All Partner Custom Param")

	param := getUrlParam(r.URL.String())
	delete(param, "token")

	p := []model.Partner{}
	if len(param) > 0 {
		p, err = model.FindVariantPartner(param)
	} else {
		status = http.StatusBadRequest
		res.AddError(its(status), model.ErrCodeMissingOrderItem, model.ErrMessageMissingOrderItem, logger.TraceID)
		logger.SetStatus(status).Log("param :", param , "response :" , res.Errors.ToString())
		render.JSON(w, res, status)
		return
	}

	if err == model.ErrResourceNotFound {
		status = http.StatusNotFound
		res.AddError(its(status), model.ErrCodeResourceNotFound, model.ErrMessageResourceNotFound, logger.TraceID)
		logger.SetStatus(status).Log("param :", param , "response :" , res.Errors.ToString())
		render.JSON(w, res, status)
		return
	} else if err != nil {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", logger.TraceID)
		logger.SetStatus(status).Log("param :", param , "response :" , res.Errors.ToString())
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
	logger.SetStatus(status).Log("param :", param , "response :" , d)
	render.JSON(w, res, status)
}

func UpdatePartner(w http.ResponseWriter, r *http.Request) {
	apiName := "partner_update"
	valid := false

	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag("Update Partner")

	id := r.FormValue("id")
	var rd Partner
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		//log.Panic(err)
		logger.SetStatus(http.StatusBadRequest).Log("param :", id, r.Body , "response :" , err.Error())
	}

	status := http.StatusUnauthorized
	err := model.ErrInvalidRole
	errorTitle := model.ErrCodeInvalidRole
	res := NewResponse(nil)
	res.AddError(its(status), errorTitle, err.Error(), "Update Partner")

	a := AuthToken(w, r)
	if a.Valid {
		for _, valueRole := range a.User.Role {
			features := model.ApiFeatures[valueRole.RoleDetail]
			for _, valueFeature := range features {
				if apiName == valueFeature {
					valid = true
				}
			}
		}

		if valid {
			status = http.StatusOK
			err := model.UpdatePartner(id, rd.SerialNumber, a.User.ID)
			if err != nil {
				status = http.StatusInternalServerError
				errorTitle = model.ErrCodeInternalError
				if err == model.ErrResourceNotFound {
					status = http.StatusNotFound
					errorTitle = model.ErrCodeResourceNotFound
				}

				res.AddError(its(status), errorTitle, err.Error(), logger.TraceID)
				logger.SetStatus(http.StatusBadRequest).Log("param :", id, rd , "response :" , res.Errors.ToString())
			} else {
				res = NewResponse("Success")
				logger.SetStatus(http.StatusBadRequest).Log("param :", id, rd , "response : success")
			}

		}
	} else {
		res = a.res
		status = http.StatusUnauthorized
	}

	render.JSON(w, res, status)
}

func DeletePartner(w http.ResponseWriter, r *http.Request) {
	apiName := "partner_delete"
	valid := false

	id := r.FormValue("id")

	status := http.StatusUnauthorized
	err := model.ErrInvalidRole
	errorTitle := model.ErrCodeInvalidRole
	res := NewResponse(nil)
	res.AddError(its(status), errorTitle, err.Error(), "Delete Partner")

	a := AuthToken(w, r)

	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag("Delete Partner")

	if a.Valid {
		for _, valueRole := range a.User.Role {
			features := model.ApiFeatures[valueRole.RoleDetail]
			for _, valueFeature := range features {
				if apiName == valueFeature {
					valid = true
				}
			}
		}

		if valid {
			status = http.StatusOK
			err := model.DeletePartner(id, a.User.ID)
			if err != nil {
				status = http.StatusInternalServerError
				errorTitle = model.ErrCodeInternalError
				if err == model.ErrResourceNotFound {
					status = http.StatusNotFound
					errorTitle = model.ErrCodeResourceNotFound
				}

				res.AddError(its(status), errorTitle, err.Error(), logger.TraceID)
				logger.SetStatus(http.StatusBadRequest).Log("param :", id , "response :" , res.Errors.ToString())
			} else {
				res = NewResponse("Success")
				logger.SetStatus(http.StatusBadRequest).Log("param :", id , "response : success " )
			}
		}
	} else {
		res = a.res
		status = http.StatusUnauthorized
	}
	render.JSON(w, res, status)
}

// dashboard
func AddPartner(w http.ResponseWriter, r *http.Request) {
	apiName := "partner_create"
	valid := false

	var rd Partner
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		log.Panic(err)
	}

	status := http.StatusUnauthorized
	err := model.ErrInvalidRole
	errorTitle := model.ErrCodeInvalidRole
	res := NewResponse(nil)
	res.AddError(its(status), errorTitle, err.Error(), "Create Partner")

	a := AuthToken(w, r)

	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag("Add Partner")

	if a.Valid {
		for _, valueRole := range a.User.Role {
			features := model.ApiFeatures[valueRole.RoleDetail]
			for _, valueFeature := range features {
				if apiName == valueFeature {
					valid = true
				}
			}
		}

		if valid {

			status = http.StatusCreated
			param := model.Partner{
				PartnerName: rd.PartnerName,
				SerialNumber: sql.NullString{
					String: rd.SerialNumber,
					Valid:  true,
				},
				CreatedBy: sql.NullString{
					String: a.User.ID,
					Valid:  true,
				},
				Tag: sql.NullString{
					String: rd.Tag,
					Valid:  true,
				},
				Description: sql.NullString{
					String: rd.Description,
					Valid:  true,
				},
			}
			err = model.InsertPartner(param)
			if err != nil {
				status = http.StatusInternalServerError
				errorTitle = model.ErrCodeInternalError
				if err == model.ErrResourceNotFound {
					status = http.StatusNotFound
					errorTitle = model.ErrCodeResourceNotFound
				}

				res.AddError(its(status), errorTitle, err.Error(), logger.TraceID)
				logger.SetStatus(http.StatusBadRequest).Log("param :", rd , "response :" , res.Errors.ToString())
			} else {
				res = NewResponse("Success")
				logger.SetStatus(http.StatusBadRequest).Log("param :", rd , "response : success")
			}
		}
	} else {
		res = a.res
		status = http.StatusUnauthorized
	}
	render.JSON(w, res, status)
}

// ------------------------------------------------------------------------------
// Tag

func GetAllTags(w http.ResponseWriter, r *http.Request) {
	status := http.StatusOK
	res := NewResponse(nil)

	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag("Get All Tag")

	tag, err := model.FindAllTags()
	if err != nil {
		fmt.Println(err.Error())
		status = http.StatusInternalServerError
		errorTitle := model.ErrCodeInternalError
		if err == model.ErrResourceNotFound {
			status = http.StatusNotFound
			errorTitle = model.ErrCodeResourceNotFound
		}

		res.AddError(its(status), errorTitle, err.Error(), logger.TraceID)
		logger.SetStatus(status).Log("param :" , "response :" , res.Errors.ToString())
	} else {
		res = NewResponse(tag)
		logger.SetStatus(status).Log("param :" , "response :" , tag)
	}

	render.JSON(w, res, status)
}

func AddTag(w http.ResponseWriter, r *http.Request) {
	apiName := "tag_create"
	valid := false

	var rd Tag
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		log.Panic(err)
	}

	status := http.StatusUnauthorized
	err := model.ErrInvalidRole
	errorTitle := model.ErrCodeInvalidRole
	res := NewResponse(nil)
	res.AddError(its(status), errorTitle, err.Error(), "Add Tag")

	a := AuthToken(w, r)

	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag("Add Tag")
	if a.Valid {
		for _, valueRole := range a.User.Role {
			features := model.ApiFeatures[valueRole.RoleDetail]
			for _, valueFeature := range features {
				if apiName == valueFeature {
					valid = true
				}
			}
		}

		if valid {
			status = http.StatusCreated
			err := model.InsertTag(rd.Value, a.User.ID)
			if err != nil {
				status = http.StatusInternalServerError
				errorTitle = model.ErrCodeInternalError
				if err == model.ErrResourceNotFound {
					status = http.StatusNotFound
					errorTitle = model.ErrCodeResourceNotFound
				}

				res.AddError(its(status), errorTitle, err.Error(), logger.TraceID)
				logger.SetStatus(http.StatusBadRequest).Log("param :", rd , "response :" , res.Errors.ToString())
			}

		} else {
			res = NewResponse("Success")
			logger.SetStatus(http.StatusBadRequest).Log("param :", rd , "response : success")
		}
	} else {
		res = a.res
		status = http.StatusUnauthorized
	}
	render.JSON(w, res, status)
}

func DeleteTag(w http.ResponseWriter, r *http.Request) {
	apiName := "tag_delete"
	valid := false

	id := bone.GetValue(r, "id")

	status := http.StatusUnauthorized
	err := model.ErrInvalidRole
	errorTitle := model.ErrCodeInvalidRole
	res := NewResponse(nil)
	res.AddError(its(status), errorTitle, err.Error(), "Delete Tag")

	a := AuthToken(w, r)

	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag("Delete Tag")

	if a.Valid {
		for _, valueRole := range a.User.Role {
			features := model.ApiFeatures[valueRole.RoleDetail]
			for _, valueFeature := range features {
				if apiName == valueFeature {
					valid = true
				}
			}
		}

		if valid {
			status = http.StatusOK
			err := model.DeleteTag(id, a.User.ID)
			if err != nil {
				status = http.StatusInternalServerError
				errorTitle = model.ErrCodeInternalError
				if err == model.ErrResourceNotFound {
					status = http.StatusNotFound
					errorTitle = model.ErrCodeResourceNotFound
				}

				res.AddError(its(status), errorTitle, err.Error(), logger.TraceID)
				logger.SetStatus(http.StatusBadRequest).Log("param :", id , "response :" , res.Errors.ToString())
			} else {
				res = NewResponse("Success")
				logger.SetStatus(http.StatusBadRequest).Log("param :", id , "response : success")
			}
		}
	} else {
		res = a.res
		status = http.StatusUnauthorized
	}
	render.JSON(w, res, status)
}

func DeleteTagBulk(w http.ResponseWriter, r *http.Request) {
	apiName := "tag_delete"
	valid := false

	var rd Tags
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		log.Panic(err)
	}

	status := http.StatusUnauthorized
	err := model.ErrInvalidRole
	errorTitle := model.ErrCodeInvalidRole
	res := NewResponse(nil)
	res.AddError(its(status), errorTitle, err.Error(), "Delete Tag")

	a := AuthToken(w, r)

	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag("Delete Tag bulk")

	if a.Valid {
		for _, valueRole := range a.User.Role {
			features := model.ApiFeatures[valueRole.RoleDetail]
			for _, valueFeature := range features {
				if apiName == valueFeature {
					valid = true
				}
			}
		}

		if valid {
			status = http.StatusOK
			err := model.DeleteTagBulk(rd.Value, a.User.ID)
			if err != nil {
				status = http.StatusInternalServerError
				errorTitle = model.ErrCodeInternalError
				if err == model.ErrResourceNotFound {
					status = http.StatusNotFound
					errorTitle = model.ErrCodeResourceNotFound
				}

				res.AddError(its(status), errorTitle, err.Error(), "Get tag")
				logger.SetStatus(http.StatusBadRequest).Log("param :", rd , "response :" , res.Errors.ToString())
			} else {
				res = NewResponse("Success")
				logger.SetStatus(http.StatusBadRequest).Log("param :", rd , "response : success")
			}
		}
	} else {
		res = a.res
		status = http.StatusUnauthorized
	}
	render.JSON(w, res, status)
}
