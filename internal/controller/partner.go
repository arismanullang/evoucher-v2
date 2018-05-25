package controller

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-zoo/bone"
	"github.com/ruizu/render"

	"strconv"

	"github.com/gilkor/evoucher/internal/model"
)

type (
	Partner struct {
		ID                string `json:"id"`
		Name              string `json:"name"`
		SerialNumber      string `json:"serial_number"`
		Email             string `json:"email"`
		Tag               string `json:"tag"`
		Description       string `json:"description"`
		Address           string `json:"address"`
		Building          string `json:"building"`
		City              string `json:"city"`
		Province          string `json:"province"`
		ZipCode           string `json:"zip_code"`
		CompanyName       string `json:"company_name"`
		CompanyPic        string `json:"company_pic"`
		CompanyTelp       string `json:"company_telp"`
		CompanyEmail      string `json:"company_email"`
		BankName          string `json:"bank_name"`
		BankBranch        string `json:"bank_branch"`
		BankAccountNumber string `json:"bank_account_number"`
		BankAccountHolder string `json:"bank_account_holder"`
		CreatedBy         string `json:"created_by"`
		CreatedAt         string `json:"created_at"`
	}
	PartnerResponseDetails []PartnerResponse
	PartnerResponse        struct {
		Name         string `json:"name"`
		SerialNumber string `json:"serial_number"`
		ProgramID    string `json:"program_id"`
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

	MobilePartnerObj struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
)

func GetProgramPartners(w http.ResponseWriter, r *http.Request) {
	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag("partner_get_program")

	param := getUrlParam(r.URL.String())
	status := http.StatusOK
	res := NewResponse(nil)

	partner, err := model.FindProgramPartner(param)
	res = NewResponse(partner)
	if err != nil {
		status = http.StatusInternalServerError
		errorTitle := model.ErrCodeInternalError
		if err == model.ErrResourceNotFound {
			status = http.StatusNotFound
			errorTitle = model.ErrCodeResourceNotFound
		}

		res.AddError(its(status), errorTitle, err.Error(), logger.TraceID)
		logger.SetStatus(status).Info("param :", param, "response :", err.Error())
	}

	render.JSON(w, res, status)
}

func GetPartners(w http.ResponseWriter, r *http.Request) {
	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag("partner_get")

	param := getUrlParam(r.URL.String())
	status := http.StatusOK
	res := NewResponse(nil)
	partner, err := model.FindPartners(param)
	res = NewResponse(partner)
	if err != nil {
		status = http.StatusInternalServerError
		errorTitle := model.ErrCodeInternalError
		if err == model.ErrResourceNotFound {
			status = http.StatusNotFound
			errorTitle = model.ErrCodeResourceNotFound
		}

		res.AddError(its(status), errorTitle, err.Error(), "Get Partner")
		logger.SetStatus(status).Info("param :", param, "response :", err.Error())
	}

	render.JSON(w, res, status)
}

func GetAllPartners(w http.ResponseWriter, r *http.Request) {
	status := http.StatusOK
	res := NewResponse(nil)

	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag("partner_all")

	a := AuthTokenWithLogger(w, r, logger)
	if !a.Valid {
		res = a.res
		status = http.StatusUnauthorized
		render.JSON(w, res, status)
		return
	}

	partner, err := model.FindAllPartners(a.User.Account.Id)
	res = NewResponse(partner)
	if err != nil {
		status = http.StatusInternalServerError
		errorTitle := model.ErrCodeInternalError
		if err == model.ErrResourceNotFound {
			status = http.StatusNotFound
			errorTitle = model.ErrCodeResourceNotFound
		}

		res.AddError(its(status), errorTitle, err.Error(), logger.TraceID)
		logger.SetStatus(status).Info("param :", a.User.Account.Id, "response :", err.Error())
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
		p, err = model.FindProgramPartner(param)
	} else {
		status = http.StatusBadRequest
		res.AddError(its(status), model.ErrCodeMissingOrderItem, model.ErrMessageMissingOrderItem, logger.TraceID)
		logger.SetStatus(status).Log("param :", param, "response :", res.Errors.ToString())
		render.JSON(w, res, status)
		return
	}

	if err == model.ErrResourceNotFound {
		status = http.StatusNotFound
		res.AddError(its(status), model.ErrCodeResourceNotFound, model.ErrMessageResourceNotFound, logger.TraceID)
		logger.SetStatus(status).Log("param :", param, "response :", res.Errors.ToString())
		render.JSON(w, res, status)
		return
	} else if err != nil {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", logger.TraceID)
		logger.SetStatus(status).Log("param :", param, "response :", res.Errors.ToString())
		render.JSON(w, res, status)
		return
	}

	d := make(PartnerResponseDetails, len(p))
	for i, v := range p {
		d[i].Name = v.Name
		d[i].SerialNumber = v.SerialNumber.String
		d[i].ProgramID = v.ProgramID
		d[i].CreatedBy = v.CreatedBy.String
	}

	status = http.StatusOK
	res = NewResponse(d)
	logger.SetStatus(status).Log("param :", param, "response :", d)
	render.JSON(w, res, status)
}

func UpdatePartner(w http.ResponseWriter, r *http.Request) {
	apiName := "partner_update"

	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag(apiName)

	id := r.FormValue("id")
	var rd Partner
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		logger.SetStatus(http.StatusInternalServerError).Panic("param :", r.Body, "response :", err.Error())
	}

	status := http.StatusOK
	res := NewResponse(nil)

	a := AuthTokenWithLogger(w, r, logger)
	if !a.Valid {
		res = a.res
		status = http.StatusUnauthorized
		render.JSON(w, res, status)
		return
	}

	if CheckAPIRole(a, apiName) {
		logger.SetStatus(status).Info("param :", a.User.ID, "response :", "Invalid Role")

		status = http.StatusUnauthorized
		res.AddError(its(status), model.ErrCodeInvalidRole, model.ErrInvalidRole.Error(), logger.TraceID)
		render.JSON(w, res, status)
		return
	}

	var serial, desc, tag sql.NullString
	var email, name string

	if rd.SerialNumber != "" {
		serial = sql.NullString{String: rd.SerialNumber, Valid: true}
	}

	if rd.Description != "" {
		desc = sql.NullString{String: rd.Description, Valid: true}
	}

	if rd.Email != "" {
		email = strings.Trim(rd.Email, ";")
	}

	if rd.Name != "" {
		name = rd.Name
	}

	if rd.Tag != "" {
		tag = sql.NullString{String: rd.Tag, Valid: true}
	}

	partner := model.PartnerUpdateRequest{
		Id:                id,
		Name:              name,
		SerialNumber:      serial.String,
		Email:             email,
		Description:       desc.String,
		Building:          rd.Building,
		Address:           rd.Address,
		City:              rd.City,
		Province:          rd.Province,
		ZipCode:           rd.ZipCode,
		Tag:               tag.String,
		BankName:          rd.BankName,
		BankBranch:        rd.BankBranch,
		BankAccountNumber: rd.BankAccountNumber,
		BankAccountHolder: rd.BankAccountHolder,
		CompanyName:       rd.CompanyName,
		CompanyEmail:      rd.CompanyEmail,
		CompanyTelp:       rd.CompanyTelp,
		CompanyPic:        rd.CompanyPic,
	}

	err := model.UpdatePartner(partner, a.User.ID, a.User.Account.Id)
	if err != nil {
		status = http.StatusInternalServerError
		errorTitle := model.ErrCodeInternalError
		if err == model.ErrResourceNotFound {
			status = http.StatusNotFound
			errorTitle = model.ErrCodeResourceNotFound
		}

		res.AddError(its(status), errorTitle, err.Error(), logger.TraceID)
		logger.SetStatus(status).Info("param :", id+" || "+rd.SerialNumber, "response :", err.Error())
	}

	render.JSON(w, res, status)
}

func DeletePartner(w http.ResponseWriter, r *http.Request) {
	apiName := "partner_delete"

	id := r.FormValue("id")
	status := http.StatusOK
	res := NewResponse(nil)

	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag(apiName)

	a := AuthTokenWithLogger(w, r, logger)
	if !a.Valid {
		res = a.res
		status = http.StatusUnauthorized
		render.JSON(w, res, status)
		return
	}

	if CheckAPIRole(a, apiName) {
		logger.SetStatus(status).Info("param :", a.User.ID, "response :", "Invalid Role")

		status = http.StatusUnauthorized
		res.AddError(its(status), model.ErrCodeInvalidRole, model.ErrInvalidRole.Error(), logger.TraceID)
		render.JSON(w, res, status)
		return
	}

	status = http.StatusOK
	err := model.DeletePartner(id, a.User.ID)
	if err != nil {
		status = http.StatusInternalServerError
		errorTitle := model.ErrCodeInternalError
		if err == model.ErrResourceNotFound {
			status = http.StatusNotFound
			errorTitle = model.ErrCodeResourceNotFound
		}

		res.AddError(its(status), errorTitle, err.Error(), logger.TraceID)
		logger.SetStatus(status).Info("param :", id, "response :", err.Error())
	}

	render.JSON(w, res, status)
}

// dashboard
func AddPartner(w http.ResponseWriter, r *http.Request) {
	apiName := "partner_create"

	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag(apiName)

	var rd Partner
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		logger.SetStatus(http.StatusInternalServerError).Panic("param :", r.Body, "response :", err.Error())
	}

	status := http.StatusCreated
	res := NewResponse(nil)

	a := AuthTokenWithLogger(w, r, logger)
	if !a.Valid {
		res = a.res
		status = http.StatusUnauthorized
		render.JSON(w, res, status)
		return
	}

	if CheckAPIRole(a, apiName) {
		logger.SetStatus(status).Info("param :", a.User.ID, "response :", "Invalid Role")

		status = http.StatusUnauthorized
		res.AddError(its(status), model.ErrCodeInvalidRole, model.ErrInvalidRole.Error(), logger.TraceID)
		render.JSON(w, res, status)
		return
	}

	param := model.Partner{
		Name:      rd.Name,
		AccountId: a.User.Account.Id,
		SerialNumber: sql.NullString{
			String: rd.SerialNumber,
			Valid:  true,
		},
		Email: strings.Trim(rd.Email, ";"),
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
		Building:          rd.Building,
		Address:           rd.Address,
		City:              rd.City,
		Province:          rd.Province,
		ZipCode:           rd.ZipCode,
		CompanyName:       rd.CompanyName,
		CompanyPic:        rd.CompanyPic,
		CompanyTelp:       rd.CompanyTelp,
		CompanyEmail:      rd.CompanyEmail,
		BankName:          rd.BankName,
		BankBranch:        rd.BankBranch,
		BankAccountHolder: rd.BankAccountHolder,
		BankAccountNumber: rd.BankAccountNumber,
	}
	err := model.InsertPartner(param)
	if err != nil {
		status = http.StatusInternalServerError
		errorTitle := model.ErrCodeInternalError
		if err == model.ErrResourceNotFound {
			status = http.StatusNotFound
			errorTitle = model.ErrCodeResourceNotFound
		}

		res.AddError(its(status), errorTitle, err.Error(), logger.TraceID)
		logger.SetStatus(status).Info("param :", param, "response :", err.Error())
	}

	render.JSON(w, res, status)
}

func GetProgramsPartner(w http.ResponseWriter, r *http.Request) {
	apiName := "partner_performance"

	id := r.FormValue("id")

	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag(apiName)

	status := http.StatusOK
	res := NewResponse(nil)

	a := AuthToken(w, r)
	if !a.Valid {
		res = a.res
		status = http.StatusUnauthorized
		render.JSON(w, res, status)
		return
	}

	if CheckAPIRole(a, apiName) {
		logger.SetStatus(status).Info("param :", a.User.ID, "response :", "Invalid Role")

		status = http.StatusUnauthorized
		res.AddError(its(status), model.ErrCodeInvalidRole, model.ErrInvalidRole.Error(), logger.TraceID)
		render.JSON(w, res, status)
		return
	}

	programs, err := model.FindProgramsPartner(id, a.User.Account.Id)
	if err != nil {
		if err != model.ErrResourceNotFound {
			status = http.StatusInternalServerError
			errorTitle := model.ErrCodeInternalError

			res.AddError(its(status), errorTitle, err.Error(), logger.TraceID)
			logger.SetStatus(status).Info("param :", id+" || "+a.User.Account.Id, "response :", err.Error())
			render.JSON(w, res, status)
			return
		}
	}

	res = NewResponse(programs)
	render.JSON(w, res, status)
}

func GetPerformancePartner(w http.ResponseWriter, r *http.Request) {
	apiName := "partner_performance"

	id := r.FormValue("id")

	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag(apiName)

	status := http.StatusOK
	res := NewResponse(nil)

	a := AuthToken(w, r)
	if !a.Valid {
		res = a.res
		status = http.StatusUnauthorized
		render.JSON(w, res, status)
		return
	}

	if CheckAPIRole(a, apiName) {
		logger.SetStatus(status).Info("param :", a.User.ID, "response :", "Invalid Role")

		status = http.StatusUnauthorized
		res.AddError(its(status), model.ErrCodeInvalidRole, model.ErrInvalidRole.Error(), logger.TraceID)
		render.JSON(w, res, status)
		return
	}

	var result PartnerPerformance

	transactions, err := model.FindTransactionsByPartner(a.User.Account.Id, id)
	if err != nil {
		if err != model.ErrResourceNotFound {
			status = http.StatusInternalServerError
			errorTitle := model.ErrCodeInternalError

			res.AddError(its(status), errorTitle, err.Error(), logger.TraceID)
			logger.SetStatus(status).Info("param :", id+" || "+a.User.Account.Id, "response :", err.Error())
			render.JSON(w, res, status)
			return
		}

		result.TransactionCode = 0
		result.TransactionValue = 0
		result.Customer = 0
	} else {
		var trValue float32
		for _, v := range transactions {
			trValue += v.VoucherValue
		}

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
		result.TransactionCode = len(transactions)
		result.TransactionValue = trValue
	}

	programs, err := model.FindProgramsPartner(id, a.User.Account.Id)
	if err != nil {
		if err != model.ErrResourceNotFound {
			status = http.StatusInternalServerError
			errorTitle := model.ErrCodeInternalError

			res.AddError(its(status), errorTitle, err.Error(), logger.TraceID)
			logger.SetStatus(status).Info("param :", id+" || "+a.User.Account.Id, "response :", err.Error())
			render.JSON(w, res, status)
			return
		}

		result.Program = 0
		result.VoucherGenerated = 0
		result.VoucherUsed = 0
	} else {
		result.Program = len(programs)

		var tVoucher int
		for _, v := range programs {
			tempVoucher, _ := strconv.Atoi(v.Voucher)
			tVoucher += tempVoucher
		}
		result.VoucherGenerated = tVoucher
		result.VoucherUsed = len(transactions)
	}

	res = NewResponse(result)
	render.JSON(w, res, status)
}

func GetDailyPerformancePartner(w http.ResponseWriter, r *http.Request) {
	apiName := "partner_performance"

	id := r.FormValue("id")

	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag(apiName)

	status := http.StatusOK
	res := NewResponse(nil)

	a := AuthToken(w, r)
	if !a.Valid {
		res = a.res
		status = http.StatusUnauthorized
		render.JSON(w, res, status)
		return
	}

	if CheckAPIRole(a, apiName) {
		logger.SetStatus(status).Info("param :", a.User.ID, "response :", "Invalid Role")

		status = http.StatusUnauthorized
		res.AddError(its(status), model.ErrCodeInvalidRole, model.ErrInvalidRole.Error(), logger.TraceID)
		render.JSON(w, res, status)
		return
	}

	var result PartnerPerformance

	transactions, err := model.FindTodayTransactionByPartner(a.User.Account.Id, id)
	if err != nil {
		if err != model.ErrResourceNotFound {
			status = http.StatusInternalServerError
			errorTitle := model.ErrCodeInternalError

			res.AddError(its(status), errorTitle, "Find transaction | "+err.Error(), logger.TraceID)
			logger.SetStatus(status).Info("param :", id+" || "+a.User.Account.Id, "response :", err.Error())
			render.JSON(w, res, status)
			return
		}

		result.TransactionCode = 0
		result.TransactionValue = 0
		result.Customer = 0
	} else {
		result.TransactionCode = len(transactions)

		var trValue float32
		for _, v := range transactions {
			trValue += v.VoucherValue
		}
		result.TransactionValue = trValue

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
	}

	programs, err := model.FindTodayProgramsPartner(id, a.User.Account.Id)
	if err != nil {
		if err != model.ErrResourceNotFound {
			status = http.StatusInternalServerError
			errorTitle := model.ErrCodeInternalError

			res.AddError(its(status), errorTitle, "Find program | "+err.Error(), logger.TraceID)
			logger.SetStatus(status).Info("param :", id+" || "+a.User.Account.Id, "response :", err.Error())
			render.JSON(w, res, status)
			return
		}

		result.Program = 0
		result.VoucherGenerated = 0
		result.VoucherUsed = 0
	} else {
		result.Program = len(programs)

		var tVoucher int
		for _, v := range programs {
			tempVoucher, _ := strconv.Atoi(v.Voucher)
			tVoucher += tempVoucher
		}
		result.VoucherGenerated = tVoucher
		result.VoucherUsed = len(transactions)
	}

	res = NewResponse(result)
	render.JSON(w, res, status)
}

func GetProgramPartnerSummary(w http.ResponseWriter, r *http.Request) {
	apiName := "partner_performance"
	programId := r.FormValue("program_id")

	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag(apiName)

	status := http.StatusOK
	res := NewResponse(nil)

	a := AuthTokenWithLogger(w, r, logger)
	if !a.Valid {
		res = a.res
		status = http.StatusUnauthorized
		render.JSON(w, res, status)
		return
	}

	if CheckAPIRole(a, apiName) {
		logger.SetStatus(status).Info("param :", a.User.ID, "response :", "Invalid Role")

		status = http.StatusUnauthorized
		res.AddError(its(status), model.ErrCodeInvalidRole, model.ErrInvalidRole.Error(), logger.TraceID)
		render.JSON(w, res, status)
		return
	}

	transaction, err := model.FindProgramPartnerSummary(a.User.Account.Id, programId)
	res = NewResponse(transaction)
	if err != nil {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, err.Error(), logger.TraceID)
		logger.SetStatus(status).Info("param :", programId, "response :", res.Errors)
	}

	render.JSON(w, res, status)
}

// ------------------------------------------------------------------------------
// Tag

func GetAllTags(w http.ResponseWriter, r *http.Request) {
	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag("tag_all")

	status := http.StatusOK
	res := NewResponse(nil)

	a := AuthTokenWithLogger(w, r, logger)
	if !a.Valid {
		res = a.res
		status = http.StatusUnauthorized
		render.JSON(w, res, status)
		return
	}

	tag, err := model.FindAllTags(a.User.Account.Id)
	res = NewResponse(tag)
	if err != nil {
		status = http.StatusInternalServerError
		errorTitle := model.ErrCodeInternalError
		if err == model.ErrResourceNotFound {
			status = http.StatusNotFound
			errorTitle = model.ErrCodeResourceNotFound
		}

		res.AddError(its(status), errorTitle, err.Error(), "Get Tags")
		logger.SetStatus(status).Info("param :", "", "response :", err.Error())
	}

	render.JSON(w, res, status)
}

func AddTag(w http.ResponseWriter, r *http.Request) {
	apiName := "tag_create"

	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag(apiName)

	var rd Tag
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		logger.SetStatus(http.StatusInternalServerError).Panic("param :", r.Body, "response :", err.Error())
	}

	status := http.StatusCreated
	res := NewResponse(nil)

	a := AuthToken(w, r)
	if !a.Valid {
		res = a.res
		status = http.StatusUnauthorized
		render.JSON(w, res, status)
		return
	}

	if CheckAPIRole(a, apiName) {
		logger.SetStatus(status).Info("param :", a.User.ID, "response :", "Invalid Role")

		status = http.StatusUnauthorized
		res.AddError(its(status), model.ErrCodeInvalidRole, model.ErrInvalidRole.Error(), logger.TraceID)
		render.JSON(w, res, status)
		return
	}

	err := model.InsertTag(rd.Value, a.User.ID, a.User.Account.Id)
	if err != nil {
		status = http.StatusInternalServerError
		errorTitle := model.ErrCodeInternalError
		if err == model.ErrResourceNotFound {
			status = http.StatusNotFound
			errorTitle = model.ErrCodeResourceNotFound
		}

		res.AddError(its(status), errorTitle, err.Error(), logger.TraceID)
		logger.SetStatus(status).Info("param :", rd.Value, "response :", logger.TraceID)
	}

	render.JSON(w, res, status)
}

func DeleteTag(w http.ResponseWriter, r *http.Request) {
	apiName := "tag_delete"

	id := bone.GetValue(r, "id")

	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag(apiName)

	status := http.StatusOK
	res := NewResponse(nil)

	a := AuthToken(w, r)
	if !a.Valid {
		res = a.res
		status = http.StatusUnauthorized
		render.JSON(w, res, status)
		return
	}

	if CheckAPIRole(a, apiName) {
		logger.SetStatus(status).Info("param :", a.User.ID, "response :", "Invalid Role")

		status = http.StatusUnauthorized
		res.AddError(its(status), model.ErrCodeInvalidRole, model.ErrInvalidRole.Error(), logger.TraceID)
		render.JSON(w, res, status)
		return
	}

	err := model.DeleteTag(id, a.User.ID)
	if err != nil {
		status = http.StatusInternalServerError
		errorTitle := model.ErrCodeInternalError
		if err == model.ErrResourceNotFound {
			status = http.StatusNotFound
			errorTitle = model.ErrCodeResourceNotFound
		}

		res.AddError(its(status), errorTitle, err.Error(), logger.TraceID)
		logger.SetStatus(status).Info("param :", id, "response :", err.Error())
	}

	render.JSON(w, res, status)
}

func DeleteTagBulk(w http.ResponseWriter, r *http.Request) {
	apiName := "tag_delete"

	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag(apiName)

	var rd Tags
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		logger.SetStatus(http.StatusInternalServerError).Panic("param :", r.Body, "response :", err.Error())
	}

	status := http.StatusOK
	res := NewResponse(nil)

	a := AuthTokenWithLogger(w, r, logger)
	if !a.Valid {
		res = a.res
		status = http.StatusUnauthorized
		render.JSON(w, res, status)
		return
	}

	if CheckAPIRole(a, apiName) {
		logger.SetStatus(status).Info("param :", a.User.ID, "response :", "Invalid Role")

		status = http.StatusUnauthorized
		res.AddError(its(status), model.ErrCodeInvalidRole, model.ErrInvalidRole.Error(), logger.TraceID)
		render.JSON(w, res, status)
		return
	}

	status = http.StatusOK
	err := model.DeleteTagBulk(rd.Value, a.User.ID)
	if err != nil {
		status = http.StatusInternalServerError
		errorTitle := model.ErrCodeInternalError
		if err == model.ErrResourceNotFound {
			status = http.StatusNotFound
			errorTitle = model.ErrCodeResourceNotFound
		}

		res.AddError(its(status), errorTitle, err.Error(), logger.TraceID)
		logger.SetStatus(status).Info("param :", rd.Value, "response :", err.Error())
	}

	render.JSON(w, res, status)
}
