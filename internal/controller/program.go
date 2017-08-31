package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-zoo/bone"
	"github.com/ruizu/render"

	"github.com/gilkor/evoucher/internal/model"
)

type (
	ProgramReq struct {
		ReqData Program `json:"program"`
		User    string  `json:"created_by"`
	}
	Program struct {
		Name               string    `json:"name"`
		Type               string    `json:"type"`
		VoucherFormat      FormatReq `json:"voucher_format"`
		VoucherType        string    `json:"voucher_type"`
		VoucherPrice       float64   `json:"voucher_price"`
		AllowAccumulative  bool      `json:"allow_accumulative"`
		StartDate          string    `json:"start_date"`
		EndDate            string    `json:"end_date"`
		StartHour          string    `json:"start_hour"`
		EndHour            string    `json:"end_hour"`
		ValidVoucherStart  string    `json:"valid_voucher_start"`
		ValidVoucherEnd    string    `json:"valid_voucher_end"`
		VoucherLifetime    int       `json:"voucher_lifetime"`
		ValidityDays       string    `json:"validity_days"`
		VoucherValue       float64   `json:"voucher_value"`
		MaxQuantityVoucher float64   `json:"max_quantity_voucher"`
		MaxGenerateVoucher float64   `json:"max_generate_voucher"`
		MaxRedeemVoucher   float64   `json:"max_redeem_voucher"`
		RedemptionMethod   string    `json:"redemption_method"`
		ImgUrl             string    `json:"image_url"`
		Tnc                string    `json:"tnc"`
		Description        string    `json:"description"`
		ValidPartners      []string  `json:"valid_partners"`
	}
	ProgramDetailResponse struct {
		Id                 string  `json:"id"`
		AccountId          string  `json:"account_id"`
		Name               string  `json:"name"`
		Type               string  `json:"type"`
		VoucherFormat      int     `json:"voucher_format"`
		VoucherType        string  `json:"voucher_type"`
		VoucherPrice       float64 `json:"voucher_price"`
		AllowAccumulative  bool    `json:"allow_accumulative"`
		StartDate          string  `json:"start_date"`
		EndDate            string  `json:"end_date"`
		VoucherValue       float64 `json:"voucher_value"`
		MaxQuantityVoucher float64 `json:"max_quantity_voucher"`
		MaxGenerateVoucher float64 `json:"max_generate_voucher"`
		MaxRedeemVoucher   float64 `json:"max_redeem_voucher"`
		RedemptionMethod   string  `json:"redemption_method"`
		ImgUrl             string  `json:"image_url"`
		Tnc                string  `json:"tnc"`
		Description        string  `json:"description"`
	}
	FormatReq struct {
		Prefix     string `json:"prefix"`
		Postfix    string `json:"postfix"`
		Body       string `json:"body"`
		FormatType string `json:"format_type"`
		Length     int    `json:"length"`
	}
	UserProgramRequest struct {
		User string `json:"user"`
	}
	MultiUserProgramRequest struct {
		User string   `json:"user"`
		Data []string `json:"data"`
	}
	QueryRequest struct {
		Query string `json:"query"`
	}
)

func CustomQuery(w http.ResponseWriter, r *http.Request) {
	var rd QueryRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		log.Panic(err)
	}

	result, err := model.CustomQuery(rd.Query)
	if err != nil {
		fmt.Println(err.Error())
	}

	res := NewResponse(result)
	render.JSON(w, res, http.StatusOK)
}

func ListPrograms(w http.ResponseWriter, r *http.Request) {
	res := NewResponse(nil)
	var status int

	a := AuthToken(w, r)
	if !a.Valid {
		render.JSON(w, a.res, http.StatusUnauthorized)
		return
	}

	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag("program_list")

	param := getUrlParam(r.URL.String())

	param["type"] = model.ProgramTypeOnDemand
	param["account_id"] = a.User.Account.Id
	delete(param, "token")

	program, err := model.FindAvailablePrograms(a.User.Account.Id)
	if err == model.ErrResourceNotFound {
		status = http.StatusNotFound
		res.AddError(its(status), model.ErrCodeResourceNotFound, model.ErrMessageNilProgram, logger.TraceID)
		logger.SetStatus(status).Log("param :", param, "response :", res.Errors)
		render.JSON(w, res, status)
		return
	} else if err != nil {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", logger.TraceID)
		logger.SetStatus(status).Log("param :", param, "response :", res.Errors)
		render.JSON(w, res, status)
		return
	}
	d := []GetVoucherOfVariatdata{}
	for _, dt := range program {
		if (int(dt.MaxVoucher) - sti(dt.Voucher)) > 0 {
			tempVoucher := GetVoucherOfVariatdata{}
			tempVoucher.ProgramID = dt.Id
			tempVoucher.AccountId = dt.AccountId
			tempVoucher.ProgramName = dt.Name
			tempVoucher.VoucherType = dt.VoucherType
			tempVoucher.VoucherPrice = dt.VoucherPrice
			tempVoucher.VoucherValue = dt.VoucherValue
			tempVoucher.StartDate = dt.StartDate
			tempVoucher.EndDate = dt.EndDate
			tempVoucher.ImgUrl = dt.ImgUrl
			tempVoucher.MaxQty = dt.MaxVoucher
			tempVoucher.Used = sti(dt.Voucher)

			d = append(d, tempVoucher)
		}
	}

	status = http.StatusOK
	res = NewResponse(d)
	logger.SetStatus(status).Log("param :", param, "response :", d)
	render.JSON(w, res, status)
}

func ListProgramsDetails(w http.ResponseWriter, r *http.Request) {
	program := bone.GetValue(r, "id")
	res := NewResponse(nil)
	var status int

	a := AuthToken(w, r)
	if !a.Valid {
		render.JSON(w, a.res, http.StatusUnauthorized)
		return
	}

	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag("program_details")

	dt, err := model.FindProgramDetailsById(program)
	if err == model.ErrResourceNotFound {
		status = http.StatusNotFound
		res.AddError(its(status), model.ErrCodeResourceNotFound, model.ErrMessageInvalidProgram, logger.TraceID)
		logger.SetStatus(status).Log("param :", program, "response :", res.Errors)
		render.JSON(w, res, status)
		return
	} else if err != nil {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", logger.TraceID)
		logger.SetStatus(status).Log("param :", program, "response :", res.Errors)
		render.JSON(w, res, status)
		return
	}
	p, err := model.FindProgramPartner(map[string]string{"program_id": program})
	if err == model.ErrResourceNotFound {
		status = http.StatusNotFound
		res.AddError(its(status), model.ErrCodeResourceNotFound, model.ErrMessageInvalidProgram+"(Partner of Program Not Found)", logger.TraceID)
		logger.SetStatus(status).Log("param :", program, "response :", res.Errors)
		render.JSON(w, res, status)
		return
	} else if err != nil {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", logger.TraceID)
		logger.SetStatus(status).Log("param :", program, "response :", res.Errors)
		render.JSON(w, res, status)
		return
	}

	d := GetVoucherOfVariatListDetails{}
	d.ProgramID = dt.Id
	d.AccountId = dt.AccountId
	d.ProgramName = dt.Name
	d.ProgramType = dt.Type
	d.VoucherType = dt.VoucherType
	d.VoucherPrice = dt.VoucherPrice
	d.AllowAccumulative = dt.AllowAccumulative
	d.StartDate = dt.StartDate
	d.EndDate = dt.EndDate
	d.VoucherValue = dt.VoucherValue
	d.MaxQuantityVoucher = dt.MaxQuantityVoucher
	d.MaxGenerateVoucher = dt.MaxGenerateVoucher
	d.MaxRedeemVoucher = dt.MaxRedeemVoucher
	d.RedemptionMethod = dt.RedemptionMethod
	d.ImgUrl = dt.ImgUrl
	d.ProgramTnc = dt.Tnc
	d.ProgramDescription = dt.Description

	d.Partners = make([]Partner, len(p))
	for i, pd := range p {
		d.Partners[i].ID = pd.Id
		d.Partners[i].Name = pd.Name
		d.Partners[i].SerialNumber = pd.SerialNumber.String
	}

	d.Used = getCountVoucher(dt.Id)

	status = http.StatusOK
	res = NewResponse(d)
	logger.SetStatus(status).Log("param :", program, "response :", d)
	render.JSON(w, res, status)
}

func GetAllPrograms(w http.ResponseWriter, r *http.Request) {
	status := http.StatusOK
	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag("program_all")

	res := NewResponse(nil)

	a := AuthTokenWithLogger(w, r, logger)
	if !a.Valid {
		res = a.res
		status = http.StatusUnauthorized
		render.JSON(w, res, status)
		return
	}

	program, err := model.FindAllPrograms(a.User.Account.Id)
	res = NewResponse(program)
	logger.SetStatus(status).Log("param account ID from token :", a.User.Account.Id, "response :", program)
	if err != nil {
		status = http.StatusInternalServerError
		errTitle := model.ErrCodeInternalError
		if err == model.ErrResourceNotFound {
			status = http.StatusNotFound
			errTitle = model.ErrCodeResourceNotFound
		}

		res.AddError(its(status), errTitle, err.Error(), logger.TraceID)
		logger.SetStatus(status).Log("param account ID from token :", a.User.Account.Id, "response :", res.Errors)
	}
	render.JSON(w, res, status)
}

func GetProgramDetailsCustom(w http.ResponseWriter, r *http.Request) {
	param := getUrlParam(r.URL.String())
	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag("program_detail")

	status := http.StatusOK
	res := NewResponse(nil)

	a := AuthTokenWithLogger(w, r, logger)
	if !a.Valid {
		res = a.res
		status = http.StatusUnauthorized
		render.JSON(w, res, status)
		return
	}

	program, err := model.FindProgramDetailsCustomParam(param)
	res = NewResponse(program)
	logger.SetStatus(status).Log("param :", param, "response :", program)
	if err != nil {
		status = http.StatusInternalServerError
		errTitle := model.ErrCodeInternalError
		if err == model.ErrResourceNotFound {
			status = http.StatusNotFound
			errTitle = model.ErrCodeResourceNotFound
		}

		res.AddError(its(status), errTitle, err.Error(), logger.TraceID)
		logger.SetStatus(status).Log("param :", param, "response :", res.Errors)
	}

	render.JSON(w, res, status)
}

func GetPrograms(w http.ResponseWriter, r *http.Request) {
	param := getUrlParam(r.URL.String())

	status := http.StatusOK
	res := NewResponse(nil)

	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag("program_get")

	a := AuthTokenWithLogger(w, r, logger)
	if !a.Valid {
		res = a.res
		status = http.StatusUnauthorized
		render.JSON(w, res, status)
		return
	}

	program, err := model.FindProgramsCustomParam(param)
	res = NewResponse(program)
	logger.SetStatus(status).Log("param :", param, "response :", program)
	if err != nil {
		status = http.StatusInternalServerError
		errTitle := model.ErrCodeInternalError
		if err == model.ErrResourceNotFound {
			status = http.StatusNotFound
			errTitle = model.ErrCodeResourceNotFound
		}

		res.AddError(its(status), errTitle, err.Error(), logger.TraceID)
		logger.SetStatus(status).Log("param :", param, "response :", res.Errors)
	}
	render.JSON(w, res, status)
}

func CreateProgram(w http.ResponseWriter, r *http.Request) {
	apiName := "program_create"
	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag(apiName)

	res := NewResponse(nil)
	status := http.StatusCreated

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

	var rd Program
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		logger.SetStatus(status).Panic("param :", rd, "response :", err.Error())
	}

	ts, err := time.Parse("01/02/2006", rd.StartDate)
	if err != nil {
		logger.SetStatus(status).Panic("param :", rd, "response :", err.Error())
	}
	te, err := time.Parse("01/02/2006", rd.EndDate)
	if err != nil {
		logger.SetStatus(status).Panic("param :", rd, "response :", err.Error())
	}
	tvs, err := time.Parse("01/02/2006", rd.ValidVoucherStart)
	if err != nil {
		logger.SetStatus(status).Panic("param :", rd, "response :", err.Error())
	}
	tve, err := time.Parse("01/02/2006", rd.ValidVoucherEnd)
	if err != nil {
		logger.SetStatus(status).Panic("param :", rd, "response :", err.Error())
	}

	vr := model.ProgramReq{
		AccountId:          a.User.Account.Id,
		Name:               rd.Name,
		Type:               rd.Type,
		VoucherType:        rd.VoucherType,
		VoucherPrice:       rd.VoucherPrice,
		MaxQuantityVoucher: rd.MaxQuantityVoucher,
		MaxRedeemVoucher:   rd.MaxRedeemVoucher,
		MaxGenerateVoucher: rd.MaxGenerateVoucher,
		AllowAccumulative:  rd.AllowAccumulative,
		RedemptionMethod:   rd.RedemptionMethod,
		VoucherValue:       rd.VoucherValue,
		StartDate:          ts.Format("2006-01-02 15:04:05.000"),
		EndDate:            te.Format("2006-01-02 15:04:05.000"),
		StartHour:          rd.StartHour,
		EndHour:            rd.EndHour,
		ValidVoucherStart:  tvs.Format("2006-01-02 15:04:05.000"),
		ValidVoucherEnd:    tve.Format("2006-01-02 15:04:05.000"),
		VoucherLifetime:    rd.VoucherLifetime,
		ValidityDays:       rd.ValidityDays,
		ImgUrl:             rd.ImgUrl,
		Tnc:                rd.Tnc,
		Description:        rd.Description,
		ValidPartners:      rd.ValidPartners,
	}

	accountDetail, err := model.GetAccountDetailByUser(a.User.ID)
	fr := model.FormatReq{
		Prefix:     rd.VoucherFormat.Prefix,
		Postfix:    accountDetail[0].Alias,
		Body:       rd.VoucherFormat.Body,
		FormatType: rd.VoucherFormat.FormatType,
		Length:     rd.VoucherFormat.Length,
	}
	id, err := model.InsertProgram(vr, fr, a.User.ID)
	res = NewResponse(id)
	logger.SetStatus(status).Info("param :", rd, "response :", id)
	if err != nil {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, err.Error(), logger.TraceID)
		logger.SetStatus(status).Info("param :", rd, "response :", res.Errors)
	}

	render.JSON(w, res, status)
}

func UpdateProgramRoute(w http.ResponseWriter, r *http.Request) {
	apiName := "program_update"

	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag(apiName)

	res := NewResponse(nil)
	status := http.StatusOK
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

	types := r.FormValue("type")
	if types == "detail" {
		UpdateProgram(w, r, logger, a)
	} else if types == "tenant" {
		UpdateProgramTenant(w, r, logger, a)
	} else if types == "broadcast" {
		UpdateProgramBroadcast(w, r, logger, a)
	} else {
		res.AddError(its(status), model.ErrCodeRouteNotFound, model.ErrRouteNotFound.Error(), logger.TraceID)
		logger.SetStatus(status).Log("param :", "response :", res.Errors)
		render.JSON(w, res, status)
	}

}

func UpdateProgram(w http.ResponseWriter, r *http.Request, logger *model.LogField, a Auth) {
	id := r.FormValue("id")

	res := NewResponse(nil)
	status := http.StatusOK

	var rd Program
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		logger.SetStatus(status).Panic("param :", rd, "response :", err.Error())
	}

	ts, err := time.Parse("01/02/2006", rd.StartDate)
	if err != nil {
		logger.SetStatus(status).Panic("param :", rd, "response :", err.Error())
	}
	te, err := time.Parse("01/02/2006", rd.EndDate)
	if err != nil {
		logger.SetStatus(status).Panic("param :", rd, "response :", err.Error())
	}

	tvs, err := time.Parse("2006-01-02T00:00:00Z", rd.ValidVoucherStart)
	if err != nil {
		logger.SetStatus(status).Panic("param :", rd, "response :", err.Error())
	}
	tve, err := time.Parse("2006-01-02T00:00:00Z", rd.ValidVoucherEnd)
	if err != nil {
		logger.SetStatus(status).Panic("param :", rd, "response :", err.Error())
	}

	vr := model.Program{
		Id:                 id,
		Name:               rd.Name,
		Type:               rd.Type,
		VoucherType:        rd.VoucherType,
		VoucherPrice:       rd.VoucherPrice,
		MaxQuantityVoucher: rd.MaxQuantityVoucher,
		MaxGenerateVoucher: rd.MaxGenerateVoucher,
		MaxRedeemVoucher:   rd.MaxRedeemVoucher,
		RedemptionMethod:   rd.RedemptionMethod,
		VoucherValue:       rd.VoucherValue,
		StartDate:          ts.Format("2006-01-02 15:04:05.000"),
		EndDate:            te.Format("2006-01-02 15:04:05.000"),
		StartHour:          rd.StartHour,
		EndHour:            rd.EndHour,
		AllowAccumulative:  rd.AllowAccumulative,
		ValidVoucherStart:  tvs.Format("2006-01-02 15:04:05.000"),
		ValidVoucherEnd:    tve.Format("2006-01-02 15:04:05.000"),
		VoucherLifetime:    rd.VoucherLifetime,
		ValidityDays:       rd.ValidityDays,
		ImgUrl:             rd.ImgUrl,
		Tnc:                rd.Tnc,
		Description:        rd.Description,
		CreatedBy:          a.User.ID,
	}
	err = model.UpdateProgram(vr)
	if err != nil {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, err.Error(), logger.TraceID)
		logger.SetStatus(status).Info("param :", rd, "response :", res.Errors)
	}

	render.JSON(w, res, status)
}

func UpdateProgramBroadcast(w http.ResponseWriter, r *http.Request, logger *model.LogField, a Auth) {
	id := r.FormValue("id")

	var rd MultiUserProgramRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		logger.SetStatus(http.StatusInternalServerError).Panic("param :", r.Body, "response :", err.Error())
	}

	res := NewResponse(nil)
	status := http.StatusOK

	d := model.UpdateProgramArrayRequest{
		ProgramId: id,
		User:      a.User.ID,
		Data:      rd.Data,
	}

	if err := model.UpdateProgramBroadcasts(d); err != nil {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, err.Error(), logger.TraceID)
		logger.SetStatus(status).Info("param :", d, "response :", res.Errors)
	}

	render.JSON(w, res, status)
}

func UpdateProgramTenant(w http.ResponseWriter, r *http.Request, logger *model.LogField, a Auth) {
	id := r.FormValue("id")

	var rd MultiUserProgramRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		logger.SetStatus(http.StatusInternalServerError).Panic("param :", r.Body, "response :", err.Error())
	}

	res := NewResponse(nil)
	status := http.StatusOK
	d := model.UpdateProgramArrayRequest{
		ProgramId: id,
		User:      a.User.ID,
		Data:      rd.Data,
	}

	if err := model.UpdateProgramPartners(d); err != nil {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, err.Error(), logger.TraceID)
		logger.SetStatus(status).Info("param :", d, "response :", res.Errors)
	}

	render.JSON(w, res, status)
}

func DeleteProgram(w http.ResponseWriter, r *http.Request) {
	apiName := "program_delete"

	res := NewResponse(nil)
	id := r.FormValue("id")
	status := http.StatusOK
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

	v := model.CountVoucher(id)
	if v > 0 {
		status = http.StatusBadRequest
		res.AddError(its(status), model.ErrMessageProgramHasBeenUsed, model.ErrBadRequest.Error(), logger.TraceID)
		logger.SetStatus(status).Log("param :", id, "response :", model.ErrBadRequest.Error())
		render.JSON(w, res, status)
		return
	}

	d := &model.DeleteProgramRequest{
		Id:   id,
		User: a.User.ID,
	}
	if err := d.Delete(); err != nil {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, err.Error(), logger.TraceID)
		logger.SetStatus(status).Log("param :", d, "response :", err.Error())
	}

	render.JSON(w, res, status)
}

func CheckProgram(rm, id string, qty int) (bool, error) {
	dt, err := model.FindProgramDetailsById(id)

	sd, err := time.Parse(time.RFC3339Nano, dt.StartDate)
	if err != nil {
		return false, err
	}
	ed, err := time.Parse(time.RFC3339Nano, dt.EndDate)
	if err != nil {
		return false, err
	}

	if !validdays(dt.ValidityDays) {
		return false, errors.New(model.ErrCodeRedeemNotValidDay)
	}

	if !validhours(dt.StartHour, dt.EndHour) {
		return false, errors.New(model.ErrCodeRedeemNotValidHour)
	}

	if !sd.Before(time.Now()) || !ed.After(time.Now()) {
		return false, errors.New(model.ErrCodeVoucherNotActive)
	}

	if dt.AllowAccumulative == false && qty > 1 {
		return false, errors.New(model.ErrCodeAllowAccumulativeDisable)
	}

	if dt.AllowAccumulative == false && qty > 1 {
		return false, errors.New(model.ErrCodeAllowAccumulativeDisable)
	}

	if dt.RedemptionMethod != rm {
		return false, errors.New(model.ErrCodeInvalidRedeemMethod)
	}

	return true, nil
}

func validdays(s string) bool {
	ret := false

	if s == "" || strings.ToUpper(s) == "ALL" {
		return true
	}

	vd := strings.Split(s, ";")

	if s == "all" {
		return true
	}

	for i := range vd {
		if strings.ToUpper(vd[i]) == strings.ToUpper(time.Now().Weekday().String()) {
			ret = true
			break
		}
	}
	return ret
}

func validhours(s, e string) bool {
	st := sti(strings.Replace(s, ":", "", 1))
	en := sti(strings.Replace(e, ":", "", 1))
	th, tm, _ := time.Now().Clock()
	tnow := sti(its(th) + its(tm))
	if tnow < st || tnow > en {
		return false
	}
	return true
}