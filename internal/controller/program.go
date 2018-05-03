package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-zoo/bone"
	"github.com/ruizu/render"

	"io"
	"os"

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
		Visibility         string    `json:"visibility"`
	}
	UpdateProgramRequest struct {
		Name               string   `json:"name"`
		Type               string   `json:"type"`
		VoucherFormat      string   `json:"voucher_format"`
		VoucherType        string   `json:"voucher_type"`
		VoucherPrice       float64  `json:"voucher_price"`
		AllowAccumulative  bool     `json:"allow_accumulative"`
		StartDate          string   `json:"start_date"`
		EndDate            string   `json:"end_date"`
		StartHour          string   `json:"start_hour"`
		EndHour            string   `json:"end_hour"`
		ValidVoucherStart  string   `json:"valid_voucher_start"`
		ValidVoucherEnd    string   `json:"valid_voucher_end"`
		VoucherLifetime    int      `json:"voucher_lifetime"`
		ValidityDays       string   `json:"validity_days"`
		VoucherValue       float64  `json:"voucher_value"`
		MaxQuantityVoucher float64  `json:"max_quantity_voucher"`
		MaxGenerateVoucher float64  `json:"max_generate_voucher"`
		MaxRedeemVoucher   float64  `json:"max_redeem_voucher"`
		RedemptionMethod   string   `json:"redemption_method"`
		ImgUrl             string   `json:"image_url"`
		Tnc                string   `json:"tnc"`
		Description        string   `json:"description"`
		ValidPartners      []string `json:"valid_partners"`
		Visibility         string   `json:"visibility"`
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

	SearchProgram struct {
		Id                string                 `json:"id"`
		AccountId         string                 `json:"account_id"`
		Name              string                 `json:"name"`
		Type              string                 `json:"type"`
		VoucherType       string                 `json:"voucher_type"`
		VoucherPrice      float64                `json:"voucher_price"`
		VoucherValue      float64                `json:"voucher_value"`
		AllowAccumulative bool                   `json:"allow_accumulative"`
		MaxVoucher        float64                `json:"max_quantity_voucher"`
		ImgUrl            string                 `json:"image_url"`
		StartDate         string                 `json:"start_date"`
		EndDate           string                 `json:"end_date"`
		Partners          []SearchProgramPartner `json:"partners"`
		Status            string                 `json:"status"`
		CreatedAt         string                 `json:"created_at"`
		UpdatedAt         string                 `json:"updated_at"`
	}
	SearchProgramPartner struct {
		Id      string                 `json:"id"`
		Name    string                 `json:"name"`
		Voucher []SearchProgramVoucher `json:"vouchers"`
	}
	SearchProgramVoucher struct {
		Voucher string `json:"voucher"`
		State   string `json:"state"`
	}
	SpinResponse struct {
		ColorArray              []string       `json:"colorArray"`
		SegmentValueArray       []SegmentValue `json:"segmentValuesArray"`
		SvgWidth                int            `json:"svgWidth"`
		SvgHeight               int            `json:"svgHeight"`
		WheelStrokeColor        string         `json:"wheelStrokeColor"`
		WheelStrokeWidth        int            `json:"wheelStrokeWidth"`
		WheelSize               int            `json:"wheelSize"`
		WheelTextOffsetY        int            `json:"wheelTextOffsetY"`
		WheelTextColor          string         `json:"wheelTextColor"`
		WheelTextSize           string         `json:"wheelTextSize"`
		WheelImageOffsetY       int            `json:"wheelImageOffsetY"`
		WheelImageSize          int            `json:"wheelImageSize"`
		CenterCircleSize        int            `json:"centerCircleSize"`
		CenterCircleStrokeColor string         `json:"centerCircleStrokeColor"`
		CenterCircleStrokeWidth int            `json:"centerCircleStrokeWidth"`
		CenterCircleFillColor   string         `json:"centerCircleFillColor"`
		SegmentStrokeColor      string         `json:"segmentStrokeColor"`
		SegmentStrokeWidth      int            `json:"segmentStrokeWidth"`
		CenterX                 int            `json:"centerX"`
		CenterY                 int            `json:"centerY"`
		HasShadows              bool           `json:"hasShadows"`
		NumSpins                int            `json:"numSpins"`
		SpinDestinationArray    []string       `json:"spinDestinationArray"`
		MinSpinDuration         int            `json:"minSpinDuration"`
		GameOverText            string         `json:"gameOverText"`
		InvalidSpinText         string         `json:"invalidSpinText"`
		IntroText               string         `json:"introText"`
		HasSound                bool           `json:"hasSound"`
		GameId                  string         `json:"gameId"`
		ClickToSpin             bool           `json:"clickToSpin"`
		SpinDirection           string         `json:"spinDirection"`
	}
	SegmentValue struct {
		Probability int          `json:"probability"`
		Type        string       `json:"type"`
		Value       string       `json:"value"`
		Win         bool         `json:"win"`
		ResultText  string       `json:"resultText"`
		UserData    UserDataSpin `json:"userData"`
	}
	UserDataSpin struct {
		Score int `json:"score"`
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

func ListMobilePrograms(w http.ResponseWriter, r *http.Request) {
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

	param["end_date"] = " > now()"
	param["start_date"] = " < now()"
	param["type"] = model.ProgramTypeOnDemand
	param["account_id"] = a.User.Account.Id
	delete(param, "token")

	program, err := model.FindAvailablePrograms(param)
	//if err == model.ErrResourceNotFound {
	//	status = http.StatusNotFound
	//	res.AddError(its(status), model.ErrCodeResourceNotFound, model.ErrMessageNilProgram, logger.TraceID)
	//	logger.SetStatus(status).Log("param :", param, "response :", res.Errors)
	//	render.JSON(w, res, status)
	//	return
	//} else
	if err != nil && err != model.ErrResourceNotFound {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", logger.TraceID)
		logger.SetStatus(status).Log("param :", param, "response :", res.Errors)
		render.JSON(w, res, status)
		return
	}
	d := GetVoucherOfVariatList{}
	for _, dt := range program {
		if (int(dt.MaxVoucher) - sti(dt.Voucher)) > 0 {
			tempVoucher := GetVoucherOfVariatdata{}
			tempVoucher.ProgramID = dt.Id
			tempVoucher.AccountId = dt.AccountId
			tempVoucher.ProgramName = dt.Name
			tempVoucher.VoucherType = dt.VoucherType
			tempVoucher.VoucherPrice = dt.VoucherPrice
			tempVoucher.VoucherValue = dt.VoucherValue
			tempVoucher.AllowAccumulative = dt.AllowAccumulative
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

func ListMallPrograms(w http.ResponseWriter, r *http.Request) {
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

	param["end_date"] = " > now()"
	param["start_date"] = " < now()"
	param["account_id"] = a.User.Account.Id
	delete(param, "token")

	program, err := model.FindAvailablePrograms(param)
	//if err == model.ErrResourceNotFound {
	//	status = http.StatusNotFound
	//	res.AddError(its(status), model.ErrCodeResourceNotFound, model.ErrMessageNilProgram, logger.TraceID)
	//	logger.SetStatus(status).Log("param :", param, "response :", res.Errors)
	//	render.JSON(w, res, status)
	//	return
	//} else
	if err != nil && err != model.ErrResourceNotFound {
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
			tempVoucher.ProgramType = dt.Type
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
	if err != nil {
		status = http.StatusInternalServerError
		errTitle := model.ErrCodeInternalError
		if err == model.ErrResourceNotFound {
			status = http.StatusNotFound
			errTitle = model.ErrCodeResourceNotFound
		}

		res.AddError(its(status), errTitle, err.Error(), logger.TraceID)
	}
	render.JSON(w, res, status)
}

func GetOnGoingPrograms(w http.ResponseWriter, r *http.Request) {
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

	param := make(map[string]string)
	param["end_date"] = " > now()"
	param["start_date"] = " < now()"
	param["account_id"] = a.User.Account.Id

	program, err := model.FindAllProgramsCustom(param)
	if err != nil {
		status = http.StatusInternalServerError
		errTitle := model.ErrCodeInternalError
		if err == model.ErrResourceNotFound {
			status = http.StatusNotFound
			errTitle = model.ErrCodeResourceNotFound
		}

		res.AddError(its(status), errTitle, err.Error(), logger.TraceID)
	}

	resProgram := []SearchProgram{}
	for _, v := range program {
		tempProgram := SearchProgram{
			Id:                v.Id,
			AccountId:         v.AccountId,
			Name:              v.Name,
			Type:              v.Type,
			VoucherType:       v.VoucherType,
			VoucherPrice:      v.VoucherPrice,
			VoucherValue:      v.VoucherValue,
			AllowAccumulative: v.AllowAccumulative,
			MaxVoucher:        v.MaxVoucher,
			ImgUrl:            v.ImgUrl,
			StartDate:         v.StartDate,
			EndDate:           v.EndDate,
			CreatedAt:         v.CreatedAt,
			UpdatedAt:         v.UpdatedAt.String,
		}

		tempVoucher := make(map[string][]SearchProgramVoucher)
		tempPartners := make(map[string]string)
		for _, vv := range v.Vouchers {
			tempPartner := strings.ToLower(vv.PartnerId.String)
			if tempPartner == "" {
				tempPartner = "0"
			}
			tempVoucher[tempPartner] = append(tempVoucher[tempPartner], SearchProgramVoucher{vv.Voucher, vv.State})
			tempPartners[tempPartner] = vv.PartnerName.String
		}
		result := []SearchProgramPartner{}
		for kk, vv := range tempVoucher {
			tempPartners := SearchProgramPartner{
				Id:      kk,
				Name:    tempPartners[kk],
				Voucher: vv,
			}
			result = append(result, tempPartners)
		}

		tempProgram.Partners = result
		resProgram = append(resProgram, tempProgram)
	}
	res = NewResponse(resProgram)
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
		logger.SetStatus(http.StatusBadRequest).Panic("param :", rd, "response :", err.Error())
	}

	ts, err := time.Parse("01/02/2006", rd.StartDate)
	if err != nil {
		logger.SetStatus(status).Panic("param :", rd, "response :", err.Error())
	}
	te, err := time.Parse("01/02/2006", rd.EndDate)
	if err != nil {
		logger.SetStatus(status).Panic("param :", rd, "response :", err.Error())
	}
	te = time.Date(te.Year(), te.Month(), te.Day(), 23, 59, 59, 0, time.Local)
	tvs, err := time.Parse("01/02/2006", rd.ValidVoucherStart)
	if err != nil {
		logger.SetStatus(status).Panic("param :", rd, "response :", err.Error())
	}
	tve, err := time.Parse("01/02/2006", rd.ValidVoucherEnd)
	if err != nil {
		logger.SetStatus(status).Panic("param :", rd, "response :", err.Error())
	}
	tve = time.Date(tve.Year(), tve.Month(), tve.Day(), 23, 59, 59, 0, time.Local)

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
		Postfix:    accountDetail.Alias,
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

	var rd UpdateProgramRequest
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
	te = time.Date(te.Year(), te.Month(), te.Day(), 23, 59, 59, 0, time.Local)

	tvs, err := time.Parse("01/02/2006", rd.ValidVoucherStart)
	if err != nil {
		logger.SetStatus(status).Panic("param :", rd, "response :", err.Error())
	}
	tve, err := time.Parse("01/02/2006", rd.ValidVoucherEnd)
	if err != nil {
		logger.SetStatus(status).Panic("param :", rd, "response :", err.Error())
	}
	tve = time.Date(tve.Year(), tve.Month(), tve.Day(), 23, 59, 59, 0, time.Local)

	bo, err := strconv.ParseBool(rd.Visibility)
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
		VoucherFormat:      sti(rd.VoucherFormat),
		Visibility:         bo,
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
		res.AddError(its(status), model.ErrBadRequest.Error(), model.ErrMessageProgramHasBeenUsed, logger.TraceID)
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

func VisibilityProgram(w http.ResponseWriter, r *http.Request) {
	apiName := "program_delete"

	res := NewResponse(nil)
	id := r.FormValue("id")
	visible := r.FormValue("visible")
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

	d := model.DeleteProgramRequest{
		Id:   id,
		User: a.User.ID,
	}
	visibility := true
	if visible == "true" {
		visibility = false
	}

	if err := model.VisibilityProgram(d, visibility); err != nil {
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

	if !validdays(dt.ValidityDays) {
		return false, errors.New(model.ErrCodeRedeemNotValidDay)
	}

	if !validhours(dt.StartHour, dt.EndHour) {
		return false, errors.New(model.ErrCodeRedeemNotValidHour)
	}

	if !sd.Before(time.Now()) {
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

func CreateTemplateCampaign(w http.ResponseWriter, r *http.Request) {
	var listTarget []string
	var listDescription []string
	programId := r.FormValue("program-id")
	apiName := "broadcast_create"

	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag(apiName)

	r.ParseMultipartForm(32 << 20)
	f, handler, err := r.FormFile("template")
	if err == http.ErrMissingFile {
		err = model.ErrResourceNotFound
	}
	if err != nil {
		err = model.ErrServerInternal
	}

	defer f.Close()
	file, err := os.OpenFile("./test/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	io.Copy(file, f)

	res := NewResponse(nil)
	status := http.StatusCreated

	a := AuthToken(w, r)
	if !a.Valid {
		res = a.res
		status = http.StatusUnauthorized
		render.JSON(w, res, status)
		return
	}

	if err := model.InsertBroadcastUser(programId, a.User.ID, listTarget, listDescription); err != nil {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, err.Error(), logger.TraceID)
		logger.SetStatus(status).Info("param :", programId+" || "+strings.Join(listTarget, ";"), "response :", err.Error())
	}

	if err := model.UpdateBulkProgram(programId, a.User.ID, len(listTarget)); err != nil {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, err.Error(), logger.TraceID)
		logger.SetStatus(status).Info("param :", programId+" || "+its(len(listTarget)), "response :", err.Error())
	}

	render.JSON(w, res, status)
}

func GetProgramType(w http.ResponseWriter, r *http.Request) {
	res := NewResponse(nil)
	var status int

	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag("program_list")

	status = http.StatusOK
	res = NewResponse(model.ProgramType)
	logger.SetStatus(status).Log("response :", model.ProgramType)
	render.JSON(w, res, status)
}

func GetListSpinPrograms(w http.ResponseWriter, r *http.Request) {
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

	param["end_date"] = " > now()"
	param["start_date"] = " < now()"
	param["type"] = model.ProgramTypeOnDemand
	param["account_id"] = a.User.Account.Id
	delete(param, "token")

	program, err := model.FindAvailablePrograms(param)
	if err != nil && err != model.ErrResourceNotFound {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", logger.TraceID)
		logger.SetStatus(status).Log("param :", param, "response :", res.Errors)
		render.JSON(w, res, status)
		return
	}

	segments := []SegmentValue{}
	for _, dt := range program {
		if (int(dt.MaxVoucher) - sti(dt.Voucher)) > 0 {
			tempUserData := UserDataSpin{
				Score: 0,
			}
			tempVoucher := SegmentValue{}
			tempVoucher.Probability = int(dt.MaxVoucher) - sti(dt.Voucher)
			tempVoucher.Type = "string"
			tempVoucher.Value = dt.Name
			tempVoucher.Win = true
			tempVoucher.ResultText = "<img width='567' height='283' src='" + dt.ImgUrl + "' /><p>" + dt.Name + "</p> `" + dt.Name + "`" + dt.Id
			tempVoucher.UserData = tempUserData

			segments = append(segments, tempVoucher)
		}
	}

	d := SpinResponse{}
	d.ColorArray = []string{"#ed1c24", "#f15a29", "#f7941e", "#8dc63f", "#39b54a", "#00a79d", "#27aae1", "#1c75bc", "#5c099c", "#9c099a"}
	d.SegmentValueArray = segments
	d.SvgWidth = 1024
	d.SvgHeight = 768
	d.WheelStrokeColor = "#ffffff" // warna lingkar
	d.WheelStrokeWidth = 18
	d.WheelSize = 700
	d.WheelTextOffsetY = 80
	d.WheelTextColor = "#EDEDED"
	d.WheelTextSize = "1.4em"
	d.WheelImageOffsetY = 40
	d.WheelImageSize = 50
	d.CenterCircleSize = 400
	d.CenterCircleStrokeColor = "#27AAE1"
	d.CenterCircleStrokeWidth = 12
	d.CenterCircleFillColor = "#EDEDED"
	d.SegmentStrokeColor = "#E2E2E2"
	d.SegmentStrokeWidth = 4
	d.CenterX = 512 //buletan tengah
	d.CenterY = 384
	d.HasShadows = false
	d.NumSpins = 1
	d.SpinDestinationArray = []string{}
	d.MinSpinDuration = 6
	d.GameOverText = ""
	d.InvalidSpinText = "INVALID SPIN. PLEASE SPIN AGAIN."
	d.IntroText = ""
	d.HasSound = false
	d.GameId = "9a0232ec06bc431114e2a7f3aea03bbe2164f1aa"
	d.ClickToSpin = true
	d.SpinDirection = "ccw"

	status = http.StatusOK
	res = NewResponse(d)
	logger.SetStatus(status).Log("param :", param, "response :", d)
	render.JSON(w, res, status)
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
	tnow := model.TimeToTimeJakarta(time.Now())
	dateNow := tnow.Format("2006-01-02")

	st, err := time.Parse(time.RFC3339, dateNow+"T"+s+":00+07:00")
	if err != nil {
		return false
	}
	en, err := time.Parse(time.RFC3339, dateNow+"T"+e+":00+07:00")
	if err != nil {
		return false
	}

	return tnow.Before(en) && tnow.After(st)
}
