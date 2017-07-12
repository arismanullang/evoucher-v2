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
	"strconv"
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
	PartnerPerformance struct {
		TransactionCode  int     `json:"transaction_code"`
		TransactionValue float32 `json:"transaction_value"`
		Program          int     `json:"program"`
		VoucherGenerated int     `json:"voucher_generated"`
		VoucherUsed      int     `json:"voucher_used"`
		Customer         int     `json:"customer"`
	}
)

func GetVariantPartners(w http.ResponseWriter, r *http.Request) {
	param := getUrlParam(r.URL.String())
	status := http.StatusOK
	res := NewResponse(nil)
	partner, err := model.FindVariantPartner(param)
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

func GetPartners(w http.ResponseWriter, r *http.Request) {
	param := getUrlParam(r.URL.String())
	status := http.StatusOK
	res := NewResponse(nil)
	partner, err := model.FindPartners(param)
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
	a := AuthToken(w, r)
	if !a.Valid {
		render.JSON(w, a.res, status)
		return
	}

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

func UpdatePartner(w http.ResponseWriter, r *http.Request) {
	apiName := "partner_update"
	valid := false

	id := r.FormValue("id")
	var rd Partner
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		log.Panic(err)
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

				res.AddError(its(status), errorTitle, err.Error(), "Get Partner")
			} else {
				res = NewResponse("Success")
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

				res.AddError(its(status), errorTitle, err.Error(), "Get Partner")
			} else {
				res = NewResponse("Success")
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

				res.AddError(its(status), errorTitle, err.Error(), "Add Partner")
			} else {
				res = NewResponse("Success")
			}
		}
	} else {
		res = a.res
		status = http.StatusUnauthorized
	}
	render.JSON(w, res, status)
}

func GetPerformancePartner(w http.ResponseWriter, r *http.Request) {
	apiName := "partner_performance"
	valid := false

	id := r.FormValue("id")

	status := http.StatusUnauthorized
	err := model.ErrInvalidRole
	errorTitle := model.ErrCodeInvalidRole
	res := NewResponse(nil)
	res.AddError(its(status), errorTitle, err.Error(), "Get Partner")

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
			var result PartnerPerformance

			transactions, err := model.FindAllTransactionByPartner(a.User.Account.Id, id)
			if err != nil {
				status = http.StatusInternalServerError
				errorTitle = model.ErrCodeInternalError
				if err == model.ErrResourceNotFound {
					status = http.StatusNotFound
					errorTitle = model.ErrCodeResourceNotFound
				}

				res.AddError(its(status), errorTitle, err.Error(), "Get Transaction")
			}
			result.TransactionCode = len(transactions)

			var trValue float32
			for _, v := range transactions {
				trValue += v.VoucherValue
			}
			result.TransactionValue = trValue

			programs, err := model.FindVariantsPartner(id, a.User.Account.Id)
			if err != nil {
				status = http.StatusInternalServerError
				errorTitle = model.ErrCodeInternalError
				if err == model.ErrResourceNotFound {
					status = http.StatusNotFound
					errorTitle = model.ErrCodeResourceNotFound
				}

				res.AddError(its(status), errorTitle, err.Error(), "Get Program")
			}
			result.Program = len(programs)

			var tVoucher int
			for _, v := range programs {
				tempVoucher, _ := strconv.Atoi(v.Voucher)
				tVoucher += tempVoucher
			}
			result.VoucherGenerated = tVoucher
			result.VoucherUsed = len(transactions)

			var nameList []string
			for _, v := range transactions {
				for _, vv := range v.Voucher {
					exist := false
					for _, name := range nameList {
						if vv.HolderDescription.String == name {
							exist = true
						}
					}

					if !exist {
						nameList = append(nameList, vv.HolderDescription.String)
					}
				}
			}
			result.Customer = len(nameList)

			res = NewResponse(result)
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
	tag, err := model.FindAllTags()
	if err != nil {
		fmt.Println(err.Error())
		status = http.StatusInternalServerError
		errorTitle := model.ErrCodeInternalError
		if err == model.ErrResourceNotFound {
			status = http.StatusNotFound
			errorTitle = model.ErrCodeResourceNotFound
		}

		res.AddError(its(status), errorTitle, err.Error(), "Get Tags")
	} else {
		res = NewResponse(tag)
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

				res.AddError(its(status), errorTitle, err.Error(), "Add Tag")
			}

		} else {
			res = NewResponse("Success")
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

				res.AddError(its(status), errorTitle, err.Error(), "Get tag")
			} else {
				res = NewResponse("Success")
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
			} else {
				res = NewResponse("Success")
			}
		}
	} else {
		res = a.res
		status = http.StatusUnauthorized
	}
	render.JSON(w, res, status)
}
