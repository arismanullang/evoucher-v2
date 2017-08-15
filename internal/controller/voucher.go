package controller

import (
	"bytes"
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gilkor/evoucher/internal/model"
	"github.com/go-zoo/bone"
	"github.com/ruizu/render"
)

type (
	// RedeemVoucherRequest represent a Request of GenerateVoucher
	RedeemVoucherRequest struct {
		AccountID string
		User      string
		State     string
		Vouchers  []string
	}
	// PayVoucherRequest represent a Request of GenerateVoucher
	PayVoucherRequest struct {
		VoucherCode string `json:"voucher_code"`
		AccountID   string `json:"account_id"`
	}
	// DeleteVoucherRequest represent a Request of GenerateVoucher
	DeleteVoucherRequest struct {
		VoucherCode string `json:"voucher_code"`
		AccountID   string `json:"account_id"`
	}
	// GenerateVoucherRequest represent a Request of GenerateVoucher
	GenerateVoucherRequest struct {
		AccountID string `json:"account_id" valid:"-"`
		ProgramID string `json:"program_id" valid:"required"`
		Quantity  int    `json:"quantity" valid:"numeric,optional"`
		Holder    struct {
			Key         string `json:"id" valid:"required"`
			Phone       string `json:"phone" valid:"numeric,optional"`
			Email       string `json:"email" valid:"email,optional"`
			Description string `json:"description" valid:"-"`
		} `json:"holder"`
		ReferenceNo string `json:"reference_no" valid:"required"`
		CreatedBy   string `json:"user" valid:"-"`
	}

	GetVoucherOfVariatList []GetVoucherOfVariatdata
	GetVoucherOfVariatdata struct {
		ProgramID    string  `json:"program_id"`
		AccountId    string  `json:"account_id"`
		ProgramName  string  `json:"program_name"`
		VoucherType  string  `json:"voucher_type"`
		VoucherPrice float64 `json:"voucher_price"`
		VoucherValue float64 `json:"voucher_value"`
		MaxQty       float64 `json:"max_quantity_voucher"`
		ImgUrl       string  `json:"image_url"`
		StartDate    string  `json:"start_date"`
		EndDate      string  `json:"end_date"`
		Used         int     `json:"used"`
	}

	GetVoucherOfVariatListDetails struct {
		ProgramID          string           `json:"program_id"`
		AccountId          string           `json:"account_id"`
		ProgramName        string           `json:"program_name"`
		ProgramType        string           `json:"program_type"`
		VoucherFormat      int              `json:"voucher_format_id"`
		VoucherType        string           `json:"voucher_type"`
		VoucherPrice       float64          `json:"voucher_price"`
		AllowAccumulative  bool             `json:"allow_accumulative"`
		StartDate          string           `json:"start_date"`
		EndDate            string           `json:"end_date"`
		VoucherValue       float64          `json:"voucher_value"`
		MaxQuantityVoucher float64          `json:"max_quantity_voucher"`
		MaxGenerateVoucher float64          `json:"max_generate_voucher"`
		MaxRedeemVoucher   float64          `json:"max_redeem_voucher"`
		RedemptionMethod   string           `json:"redemption_method"`
		ImgUrl             string           `json:"image_url"`
		ProgramTnc         string           `json:"program_tnc"`
		ProgramDescription string           `json:"program_description"`
		CreatedBy          string           `json:"created_by"`
		CreatedAt          string           `json:"created_at"`
		Used               int              `json:"used"`
		State              string           `json:"state"`
		Holder             string           `json:"holder"`
		HolderDescription  string           `json:"holder_description"`
		Partners           []Partner        `json:"partners"`
		Voucher            []VoucerResponse `json:"vouchers"`
	}

	// VoucerResponse represent list of voucher data
	VoucerResponse struct {
		VoucherID string `json:"voucher_id"`
		VoucherNo string `json:"voucher_code"`
		State     string `json:"state,omitempty"`
	}
	// DetailListResponseData represent list of voucher data
	DetailListResponseData []RespomseData
	// DetailResponseData represent list of voucher data
	RespomseData struct {
		ID                string    `json:"id"`
		VoucherCode       string    `json:"voucher_code"`
		ReferenceNo       string    `json:"reference_no"`
		Holder            string    `json:"holder"`
		HolderPhone       string    `json:"holder_phone"`
		HolderEmail       string    `json:"holder_email"`
		HolderDescription string    `json:"holder_description"`
		ProgramID         string    `json:"program_id"`
		ValidAt           time.Time `json:"valid_at"`
		ExpiredAt         time.Time `json:"expired_at"`
		VoucherValue      float64   `json:"voucher_value"`
		State             string    `json:"state"`
		CreatedBy         string    `json:"created_by"`
		CreatedAt         time.Time `json:"created_at"`
		UpdatedBy         string    `json:"updated_by"`
		UpdatedAt         time.Time `json:"updated_at"`
		DeletedBy         string    `json:"deleted_by"`
		DeletedAt         time.Time `json:"deleted_at"`
		Status            string    `json:"status"`
	}
	DetailResponseData struct {
		ID                string    `json:"id"`
		VoucherCode       string    `json:"voucher_code"`
		ReferenceNo       string    `json:"reference_no"`
		Holder            string    `json:"holder"`
		HolderPhone       string    `json:"holder_phone"`
		HolderEmail       string    `json:"holder_email"`
		HolderDescription string    `json:"holder_description"`
		ProgramID         string    `json:"program_id"`
		ValidAt           time.Time `json:"valid_at"`
		ExpiredAt         time.Time `json:"expired_at"`
		VoucherValue      float64   `json:"voucher_value"`
		State             string    `json:"state"`
		CreatedBy         string    `json:"created_by"`
		CreatedAt         time.Time `json:"created_at"`
		UpdatedBy         string    `json:"updated_by"`
		UpdatedAt         time.Time `json:"updated_at"`
		DeletedBy         string    `json:"deleted_by"`
		DeletedAt         time.Time `json:"deleted_at"`
		Status            string    `json:"status"`

		ProgramName  string  `json:"program_name"`
		VoucherPrice float64 `json:"voucher_price"`
		VoucherType  string  `json:"voucher_type"`
	}

	GetVoucherlinkResponse []GetVoucherlinkdata
	GetVoucherlinkdata     struct {
		Url         string `json:"url"`
		VoucherID   string `json:"voucher_id"`
		VoucherCode string `json:"voucher_code"`
	}
)

// ## API ##//

//GetVoucherOfProgram list voucher by holder
func GetVoucherOfProgram(w http.ResponseWriter, r *http.Request) {
	var voucher model.VoucherResponse
	var err error
	res := NewResponse(nil)
	var status int

	//Token Authentocation
	a := AuthToken(w, r)
	if !a.Valid {
		render.JSON(w, a.res, http.StatusUnauthorized)
		return
	}

	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag("My Voucher")

	param := getUrlParam(r.URL.String())
	delete(param, "token")

	if len(param) > 0 {
		voucher, err = model.FindAvailableVoucher(param)
	} else {
		status = http.StatusBadRequest
		res.AddError(its(status), model.ErrCodeMissingOrderItem, model.ErrMessageMissingOrderItem, logger.TraceID)
		logger.SetStatus(status).Log("param :", param, "response :", res.Errors.ToString())
		render.JSON(w, res, status)
		return
	}
	// fmt.Println(voucher, err)
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
	// fmt.Println(voucher.VoucherData)

	distinctProgram := []string{}
	for _, v := range voucher.VoucherData {
		if !stringInSlice(v.ProgramID, distinctProgram) {
			distinctProgram = append(distinctProgram, v.ProgramID)
		}
	}
	// fmt.Println(distinctProgram)
	d := []GetVoucherOfVariatListDetails{}
	for _, v := range distinctProgram {
		tempVoucherResponse := []VoucerResponse{}
		tempGetVoucherOfVariatListDetails := GetVoucherOfVariatListDetails{}
		tempPartners := []Partner{}

		dt, _ := model.FindProgramDetailsById(v)
		tempGetVoucherOfVariatListDetails.ProgramID = dt.Id
		tempGetVoucherOfVariatListDetails.AccountId = dt.AccountId
		tempGetVoucherOfVariatListDetails.ProgramName = dt.Name
		tempGetVoucherOfVariatListDetails.VoucherType = dt.VoucherType
		tempGetVoucherOfVariatListDetails.VoucherPrice = dt.VoucherPrice
		tempGetVoucherOfVariatListDetails.VoucherValue = dt.VoucherValue
		tempGetVoucherOfVariatListDetails.StartDate = dt.StartDate
		tempGetVoucherOfVariatListDetails.EndDate = dt.EndDate
		tempGetVoucherOfVariatListDetails.ImgUrl = dt.ImgUrl
		tempGetVoucherOfVariatListDetails.AllowAccumulative = dt.AllowAccumulative
		tempGetVoucherOfVariatListDetails.RedemptionMethod = dt.RedemptionMethod
		tempGetVoucherOfVariatListDetails.ProgramTnc = dt.Tnc
		tempGetVoucherOfVariatListDetails.ProgramDescription = dt.Description
		tempGetVoucherOfVariatListDetails.MaxQuantityVoucher = dt.MaxQuantityVoucher
		tempGetVoucherOfVariatListDetails.MaxRedeemVoucher = dt.MaxRedeemVoucher
		tempGetVoucherOfVariatListDetails.MaxGenerateVoucher = dt.MaxGenerateVoucher
		for _, vv := range voucher.VoucherData {
			if vv.ProgramID == v {
				tempVoucher := VoucerResponse{
					VoucherID: vv.ID,
					VoucherNo: vv.VoucherCode,
					State:     vv.State,
				}
				tempVoucherResponse = append(tempVoucherResponse, tempVoucher)
			}
		}
		tempGetVoucherOfVariatListDetails.Voucher = tempVoucherResponse

		pv, _ := model.FindProgramPartners(v)
		for _, vv := range pv {
			tempPartner := Partner{
				ID:           vv.Id,
				Name:         vv.Name,
				SerialNumber: vv.SerialNumber.String,
				Tag:          vv.Tag.String,
			}
			tempPartners = append(tempPartners, tempPartner)
		}
		tempGetVoucherOfVariatListDetails.Partners = tempPartners

		d = append(d, tempGetVoucherOfVariatListDetails)
	}

	// d.Vouchers = make([]VoucerResponse, len(voucher.VoucherData))
	status = http.StatusOK
	res = NewResponse(d)
	logger.SetStatus(status).Log("param :", param, "response :", d)
	render.JSON(w, res, status)
}

//GetVoucherOfProgramDetails voucher by holder
func GetVoucherOfProgramDetails(w http.ResponseWriter, r *http.Request) {
	program := bone.GetValue(r, "id")
	res := NewResponse(nil)
	var status int

	//Token Authentocation
	a := AuthToken(w, r)
	if !a.Valid {
		render.JSON(w, a.res, http.StatusUnauthorized)
		return
	}

	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag("My-Voucher-Details")

	param := getUrlParam(r.URL.String())
	param["state"] = model.VoucherStateCreated
	param["program_id"] = program
	delete(param, "token")

	if len(param) < 0 || program == "" {
		status = http.StatusBadRequest
		res.AddError(its(status), model.ErrCodeMissingOrderItem, model.ErrMessageMissingOrderItem, logger.TraceID)
		logger.SetStatus(status).Log("param :", param, "response :", res.Errors.ToString())
		render.JSON(w, res, status)
		return
	}

	voucher, err := model.FindVoucher(param)
	if err == model.ErrResourceNotFound {
		status = http.StatusNotFound
		res.AddError(its(status), model.ErrCodeResourceNotFound, model.ErrMessageInvalidHolder, logger.TraceID)
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
	dt, err := model.FindProgramDetailsById(program)
	if err == model.ErrResourceNotFound {
		status = http.StatusNotFound
		res.AddError(its(status), model.ErrCodeResourceNotFound, model.ErrMessageInvalidProgram, logger.TraceID)
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
	p, err := model.FindProgramPartner(map[string]string{"program_id": program})
	if err == model.ErrResourceNotFound {
		status = http.StatusNotFound
		res.AddError(its(status), model.ErrCodeResourceNotFound, model.ErrMessageInvalidProgram+"(Partner of Variant Not Found)", logger.TraceID)
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

	fmt.Println(dt)
	d := GetVoucherOfVariatListDetails{}
	d.ProgramID = dt.Id
	d.AccountId = dt.AccountId
	d.ProgramName = dt.Name
	d.ProgramType = dt.Type
	d.CreatedBy = dt.CreatedBy
	d.CreatedAt = dt.CreatedAt
	d.VoucherType = dt.VoucherType
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
	d.Used = getCountVoucher(dt.Id)

	d.Partners = make([]Partner, len(p))
	for i, pd := range p {
		d.Partners[i].ID = pd.Id
		d.Partners[i].Name = pd.Name
		d.Partners[i].SerialNumber = pd.SerialNumber.String
		d.Partners[i].Tag = pd.Tag.String
		d.Partners[i].Description = pd.Description.String
	}

	d.Voucher = make([]VoucerResponse, len(voucher.VoucherData))
	for j, vd := range voucher.VoucherData {
		d.Voucher[j].VoucherID = vd.ID
		d.Voucher[j].VoucherNo = vd.VoucherCode
		d.Voucher[j].State = vd.State
	}

	fmt.Println(d)
	// d.Vouchers = make([]VoucerResponse, len(voucher.VoucherData))
	status = http.StatusOK
	res = NewResponse(d)
	logger.SetStatus(status).Log("param :", param, "response :", d)
	render.JSON(w, res, status)
}

// GetVoucherDetail get Voucher detail from DB
func GetVoucherList(w http.ResponseWriter, r *http.Request) {
	var voucher model.VoucherResponse
	var err error
	res := NewResponse(nil)
	var status int

	//Token Authentocation
	a := AuthToken(w, r)
	if !a.Valid {
		render.JSON(w, a.res, http.StatusUnauthorized)
		return
	}

	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag("Voucher-List")

	param := getUrlParam(r.URL.String())
	delete(param, "token")

	if len(param) > 0 {
		voucher, err = model.FindVoucher(param)
	} else {
		status = http.StatusBadRequest
		res.AddError(its(status), model.ErrCodeMissingOrderItem, model.ErrMessageMissingOrderItem, logger.TraceID)
		logger.SetStatus(status).Log("param :", param, "response :", res.Errors.ToString())
		render.JSON(w, res, status)
		return
	}
	// fmt.Println(voucher, err)
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
	} else if voucher.Message != "" {

		dvr := make(DetailListResponseData, len(voucher.VoucherData))
		for i, v := range voucher.VoucherData {
			dvr[i].ID = v.ID
			dvr[i].VoucherCode = v.VoucherCode
			dvr[i].ReferenceNo = v.ReferenceNo
			dvr[i].Holder = v.Holder.String
			dvr[i].HolderPhone = v.HolderPhone.String
			dvr[i].HolderEmail = v.HolderEmail.String
			dvr[i].HolderDescription = v.HolderDescription.String
			dvr[i].ProgramID = v.ProgramID
			dvr[i].ValidAt = v.ValidAt
			dvr[i].ExpiredAt = v.ExpiredAt
			dvr[i].VoucherValue = v.VoucherValue
			dvr[i].State = v.State
			dvr[i].CreatedBy = v.CreatedBy
			dvr[i].CreatedAt = v.CreatedAt
			dvr[i].UpdatedBy = v.UpdatedBy.String
			dvr[i].UpdatedAt = v.UpdatedAt.Time
			dvr[i].DeletedBy = v.DeletedBy.String
			dvr[i].DeletedAt = v.DeletedAt.Time
			dvr[i].Status = v.Status
		}
		status = http.StatusOK
		res = NewResponse(dvr)
		logger.SetStatus(status).Log("param :", param, "response :", dvr)
		render.JSON(w, res, status)
		return
	}
}

func GetVoucherDetails(w http.ResponseWriter, r *http.Request) {
	vc := bone.GetValue(r, "id")
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
		SetTag("Voucher-Details")

	d, err := model.FindVoucher(map[string]string{"id": vc})

	if err == model.ErrResourceNotFound {
		status = http.StatusNotFound
		res.AddError(its(status), model.ErrCodeResourceNotFound, model.ErrMessageResourceNotFound, logger.TraceID)
		logger.SetStatus(status).Log("param :", vc, "response :", res.Errors.ToString())
		render.JSON(w, res, status)
		return
	} else if err != nil {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", logger.TraceID)
		logger.SetStatus(status).Log("param :", vc, "response :", res.Errors.ToString())
		render.JSON(w, res, status)

		return
	}

	dt, err := model.FindProgramDetailsById(d.VoucherData[0].ProgramID)
	if err == model.ErrResourceNotFound {
		status = http.StatusNotFound
		res.AddError(its(status), model.ErrCodeResourceNotFound, model.ErrMessageInvalidProgram, logger.TraceID)
		logger.SetStatus(status).Log("param :", vc, "response :", res.Errors.ToString())
		render.JSON(w, res, status)
		return
	} else if err != nil {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", logger.TraceID)
		logger.SetStatus(status).Log("param :", vc, "response :", res.Errors.ToString())
		render.JSON(w, res, status)
		return
	}
	// p, err := model.FindProgramPartner(map[string]string{"program_id": d.VoucherData[0].ProgramID})
	// if err == model.ErrResourceNotFound {
	// 	status = http.StatusNotFound
	// 	res.AddError(its(status), model.ErrCodeResourceNotFound, model.ErrMessageInvalidProgram+"(Partner of Program Not Found)", "voucher")
	// 	render.JSON(w, res, status)
	// 	return
	// } else if err != nil {
	// 	status = http.StatusInternalServerError
	// 	res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", "voucher")
	// 	render.JSON(w, res, status)
	// 	return
	// }

	dvr := DetailResponseData{
		ID:                d.VoucherData[0].ID,
		VoucherCode:       d.VoucherData[0].VoucherCode,
		ReferenceNo:       d.VoucherData[0].ReferenceNo,
		Holder:            d.VoucherData[0].Holder.String,
		HolderPhone:       d.VoucherData[0].HolderPhone.String,
		HolderEmail:       d.VoucherData[0].HolderEmail.String,
		HolderDescription: d.VoucherData[0].HolderDescription.String,
		ProgramID:         d.VoucherData[0].ProgramID,
		ValidAt:           d.VoucherData[0].ValidAt,
		ExpiredAt:         d.VoucherData[0].ExpiredAt,
		VoucherValue:      d.VoucherData[0].VoucherValue,
		State:             d.VoucherData[0].State,
		CreatedBy:         d.VoucherData[0].CreatedBy,
		CreatedAt:         d.VoucherData[0].CreatedAt,
		UpdatedBy:         d.VoucherData[0].UpdatedBy.String,
		UpdatedAt:         d.VoucherData[0].UpdatedAt.Time,
		DeletedBy:         d.VoucherData[0].DeletedBy.String,
		DeletedAt:         d.VoucherData[0].DeletedAt.Time,
		Status:            d.VoucherData[0].Status,
		ProgramName:       dt.Name,
		VoucherPrice:      dt.VoucherPrice,
		VoucherType:       dt.VoucherType,
	}

	status = http.StatusOK
	res = NewResponse(dvr)
	logger.SetStatus(status).Log("param :", vc, "response :", dvr)
	render.JSON(w, res, status)
	return
}

//GenerateVoucherOnDemand Generate singgle voucher request
func GenerateVoucherOnDemand(w http.ResponseWriter, r *http.Request) {
	var gvd GenerateVoucherRequest
	var status int
	res := NewResponse(nil)

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&gvd); err != nil {
		status = http.StatusBadRequest
		res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", "voucher")
		render.JSON(w, res, status)
		return
	}

	_, err := govalidator.ValidateStruct(gvd)
	if err != nil {
		status = http.StatusBadRequest
		res.AddError(its(status), model.ErrCodeValidationError, model.ErrMessageValidationError+"("+err.Error()+")", "transaction")
		render.JSON(w, res, status)
		return
	}

	//Token Authentocation
	a := AuthToken(w, r)
	if !a.Valid {
		render.JSON(w, a.res, http.StatusUnauthorized)
		return
	}

	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag("Generate-Voucher-Single")

	_, err = govalidator.ValidateStruct(gvd)
	if err != nil {
		status = http.StatusBadRequest
		res.AddError(its(status), model.ErrCodeValidationError, model.ErrMessageValidationError+"("+err.Error()+")", logger.TraceID)
		logger.SetStatus(status).Log("param :", gvd, "response :", res.Errors.ToString())
		render.JSON(w, res, status)
		return
	}

	dt, err := model.FindProgramDetailsById(gvd.ProgramID)
	if err == model.ErrResourceNotFound {
		status = http.StatusNotFound
		res.AddError(its(status), model.ErrCodeResourceNotFound, model.ErrMessageInvalidProgram, logger.TraceID)
		logger.SetStatus(status).Log("param :", gvd, "response :", res.Errors.ToString())
		render.JSON(w, res, status)
		return
	} else if err != nil {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInvalidProgram+"("+err.Error()+")", logger.TraceID)
		logger.SetStatus(status).Log("param :", gvd, "response :", res.Errors.ToString())
		render.JSON(w, res, status)
		return
	}
	sd, err := time.Parse(time.RFC3339Nano, dt.StartDate)
	ed, err := time.Parse(time.RFC3339Nano, dt.EndDate)
	if err != nil {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInvalidProgram+"("+err.Error()+")", logger.TraceID)
		logger.SetStatus(status).Log("param :", gvd, "response :", res.Errors.ToString())
		render.JSON(w, res, status)
		return
	}

	if (int(dt.MaxGenerateVoucher) - model.CountHolderVoucher(gvd.ProgramID, gvd.Holder.Key)) < 0 {
		status = http.StatusBadRequest
		res.AddError(its(status), model.ErrCodeVoucherQtyExceeded, model.ErrMessageVoucherQtyExceeded, logger.TraceID)
		logger.SetStatus(status).Log("param :", gvd, "response :", res.Errors.ToString())
		render.JSON(w, res, status)
		return
	} else if dt.Type != model.ProgramTypeOnDemand {
		status = http.StatusBadRequest
		res.AddError(its(status), model.ErrCodeVoucherRulesViolated, model.ErrMessageVoucherRulesViolated, logger.TraceID)
		logger.SetStatus(status).Log("param :", gvd, "response :", res.Errors.ToString())
		render.JSON(w, res, status)
		return
	} else if !sd.Before(time.Now()) {
		status = http.StatusBadRequest
		res.AddError(its(status), model.ErrCodeVoucherNotActive, model.ErrMessageVoucherNotActive, logger.TraceID)
		logger.SetStatus(status).Log("param :", gvd, "response :", res.Errors.ToString())
		render.JSON(w, res, status)
		return
	} else if !ed.After(time.Now()) {
		status = http.StatusBadRequest
		res.AddError(its(status), model.ErrCodeVoucherExpired, model.ErrMessageVoucherExpired, logger.TraceID)
		logger.SetStatus(status).Log("param :", gvd, "response :", res.Errors.ToString())
		render.JSON(w, res, status)
		return
	}

	gvd.AccountID = dt.AccountId
	gvd.ProgramID = dt.Id
	gvd.Quantity = 1
	gvd.CreatedBy = a.User.ID

	// fmt.Println("request data =>", gvd.Holder)
	var voucher []model.Voucher
	voucher, err = gvd.generateVoucher(&dt)
	if err != nil {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"( failed Genarate Voucher :"+err.Error()+")", logger.TraceID)
		logger.SetStatus(status).Log("param :", gvd, "response :", res.Errors.ToString())
		render.JSON(w, res, status)
		return
	}

	gvr := VoucerResponse{
		VoucherID: voucher[0].ID,
		VoucherNo: voucher[0].VoucherCode,
	}

	status = http.StatusCreated
	res = NewResponse(gvr)
	logger.SetStatus(status).Log("param :", gvd, "response :", gvr)
	render.JSON(w, res, status)
	return

}

//GenerateVoucher Generate bulk voucher request
func GenerateVoucherBulk(w http.ResponseWriter, r *http.Request) {
	apiName := "voucher_generate-bulk"
	var gvd GenerateVoucherRequest
	var status int
	res := NewResponse(nil)
	vrID := r.FormValue("program")
	fmt.Println("program id = ", vrID)

	//Token Authentocation
	a := AuthToken(w, r)
	if !a.Valid {
		render.JSON(w, a.res, http.StatusUnauthorized)
		return
	}

	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag("Generate-Voucher-Bulk")

	if CheckAPIRole(a, apiName) {
		render.JSON(w, model.ErrCodeInvalidRole, http.StatusUnauthorized)
		return
	}

	if getCountVoucher(vrID) > 0 {
		status = http.StatusBadRequest
		res.AddError(its(status), model.ErrMessageInvalidProgram, model.ErrMessageProgramHasBeenUsed, logger.TraceID)
		logger.SetStatus(status).Log("param :", vrID, "response :", res.Errors.ToString())
		render.JSON(w, res, status)
		return
	}

	program, err := model.FindProgramDetailsById(vrID)
	if err == model.ErrResourceNotFound {
		status = http.StatusNotFound
		res.AddError(its(status), model.ErrCodeResourceNotFound, model.ErrMessageResourceNotFound, logger.TraceID)
		logger.SetStatus(status).Log("param :", vrID, "response :", res.Errors.ToString())
		render.JSON(w, res, status)
		return
	} else if err != nil {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", logger.TraceID)
		logger.SetStatus(status).Log("param :", vrID, "response :", res.Errors.ToString())
		render.JSON(w, res, status)
		return
	}

	var listBroadcast []model.BroadcastUser
	listBroadcast, err = model.FindBroadcastUser(map[string]string{"program_id": vrID})
	if err != nil {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", logger.TraceID)
		logger.SetStatus(status).Log("param :", vrID, "response :", res.Errors.ToString())
		render.JSON(w, res, status)
		return
	}

	gvd.AccountID = a.User.Account.Id
	gvd.ProgramID = vrID
	gvd.Quantity = 1
	gvd.CreatedBy = a.User.ID

	for _, v := range listBroadcast {
		gvd.ReferenceNo = its(v.ID)
		gvd.Holder.Key = v.Description
		gvd.Holder.Description = v.Target

		_, err = gvd.generateVoucher(&program)
		if err != nil {
			fmt.Println(err)
			rollback(vrID)

			status = http.StatusInternalServerError
			res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", logger.TraceID)
			logger.SetStatus(status).Log("param :", vrID, "response :", res.Errors.ToString())
			render.JSON(w, res, status)
			return
		}
	}

	status = http.StatusCreated
	res = NewResponse("success")
	logger.SetStatus(status).Log("param :", vrID, "response : success")
	render.JSON(w, res, status)
	return

}

func GetVoucherlink(w http.ResponseWriter, r *http.Request) {
	apiName := "voucher_get-link"

	status := http.StatusOK
	res := NewResponse(nil)
	varID := r.FormValue("program")

	a := AuthToken(w, r)
	if !a.Valid {
		render.JSON(w, a.res, http.StatusUnauthorized)
		return
	}

	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag("Generate-Voucher-Link")

	if CheckAPIRole(a, apiName) {
		render.JSON(w, model.ErrCodeInvalidRole, http.StatusUnauthorized)
		return
	}

	v, err := model.FindVoucher(map[string]string{"program_id": varID})
	if err == model.ErrResourceNotFound {
		status = http.StatusNotFound
		res.AddError(its(status), model.ErrCodeResourceNotFound, model.ErrMessageInvalidHolder, logger.TraceID)
		logger.SetStatus(status).Log("param :", varID, "response :", res.Errors.ToString())
		render.JSON(w, res, status)
		return
	} else if err != nil {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", logger.TraceID)
		logger.SetStatus(status).Log("param :", varID, "response :", res.Errors.ToString())
		render.JSON(w, res, status)
		return
	}

	vl := [][]string{}
	for _, v := range v.VoucherData {
		tempName := v.Holder.String
		tempEmail := v.HolderEmail.String
		if v.HolderEmail.String == "" {
			tempEmail = v.HolderDescription.String
		}
		tempArray := []string{generateLink(v.ID), tempEmail, tempName}
		vl = append(vl, tempArray)
	}

	b := &bytes.Buffer{}   // creates IO Writer
	wr := csv.NewWriter(b) // creates a csv writer that uses the io buffer.
	for _, value := range vl {
		err := wr.Write(value) // converts array of string to comma seperated values for 1 row.
		if err != nil {
			log.Fatal("", err)
		}
	}
	wr.Flush() // writes the csv writer data to  the buffered data io writer(b(bytes.buffer))

	w.Header().Set("Content-Type", "text/csv") // setting the content type header to text/csv

	w.Header().Set("Content-Disposition", "attachment;filename=voucher.csv")
	w.Write(b.Bytes())

	logger.SetStatus(status).Log("param :", varID, "response :", vl)

	return
}

func GetCsvSample(w http.ResponseWriter, r *http.Request) {
	a := AuthToken(w, r)
	if !a.Valid {
		render.JSON(w, a.res, http.StatusUnauthorized)
		return
	}

	vl := [][]string{}
	for i := 0; i < 3; i++ {
		tempArray := []string{"index", "Name", "Email"}
		vl = append(vl, tempArray)
	}

	b := &bytes.Buffer{}   // creates IO Writer
	wr := csv.NewWriter(b) // creates a csv writer that uses the io buffer.
	for _, value := range vl {
		err := wr.Write(value) // converts array of string to comma seperated values for 1 row.
		if err != nil {
			log.Fatal("", err)
		}
	}
	wr.Flush() // writes the csv writer data to  the buffered data io writer(b(bytes.buffer))

	w.Header().Set("Content-Type", "text/csv") // setting the content type header to text/csv

	w.Header().Set("Content-Disposition", "attachment;filename=Sample.csv")
	w.Write(b.Bytes())
	return
}

// ## ### ##//

//CheckVoucherRedemption validation
func (r *TransactionRequest) CheckVoucherRedemption(voucherID string) (bool, error) {

	voucher, err := model.FindVoucher(map[string]string{"id": voucherID})

	if err != nil {
		return false, err
	} else if voucher.VoucherData[0].State == model.VoucherStateUsed {
		return false, errors.New(model.ErrMessageVoucherAlreadyUsed)
	} else if voucher.VoucherData[0].State == model.VoucherStatePaid {
		return false, errors.New(model.ErrMessageVoucherAlreadyPaid)
	} else if !voucher.VoucherData[0].ExpiredAt.After(time.Now()) {
		fmt.Println("expired date : ", voucher.VoucherData[0].ExpiredAt, voucher.VoucherData[0].ID)
		return false, errors.New(model.ErrMessageVoucherExpired)
	}

	return true, nil
}

//UpdateVoucher redeem
func (r *RedeemVoucherRequest) UpdateVoucher() (bool, error) {
	var d model.UpdateDeleteRequest

	d.State = model.VoucherStateUsed
	d.User = r.User
	d.State = r.State

	for _, v := range r.Vouchers {
		d.ID = v
		fmt.Println("update voucher :", d.ID)
		if _, err := d.UpdateVc(); err != nil {
			return false, err
		}
	}

	return true, nil
}

func (vr *GenerateVoucherRequest) generateVoucher(v *model.Program) ([]model.Voucher, error) {
	ret := make([]model.Voucher, vr.Quantity)
	var code []string
	var vcf model.VoucherCodeFormat
	var tsd, ted time.Time

	vcf, err := model.GetVoucherCodeFormat(v.VoucherFormat)
	if err != nil {
		return ret, err
	}
	if v.VoucherLifetime > 0 {
		tsd = time.Now()
		ted = time.Now().AddDate(0, 0, v.VoucherLifetime)
	} else {
		tsd, err = time.Parse(time.RFC3339Nano, v.ValidVoucherStart)
		if err != nil {
			log.Panic(err)
		}

		ted, err = time.Parse(time.RFC3339Nano, v.ValidVoucherEnd)
		if err != nil {
			log.Panic(err)
		}
	}

	for i := 0; i <= vr.Quantity-1; i++ {

		code = append(code, voucherCode(vcf, v.VoucherFormat))

		// fmt.Println("generate data =>", vr.Holder)
		rd := model.Voucher{
			VoucherCode:  code[i],
			ReferenceNo:  vr.ReferenceNo,
			ProgramID:    vr.ProgramID,
			ValidAt:      tsd,
			ExpiredAt:    ted,
			VoucherValue: v.VoucherValue,
			State:        model.VoucherStateCreated,
			CreatedBy:    vr.CreatedBy, //note harus nya by user
			CreatedAt:    time.Now(),
		}
		rd.Holder = sql.NullString{String: vr.Holder.Key, Valid: true}
		rd.HolderPhone = sql.NullString{String: vr.Holder.Phone, Valid: true}
		rd.HolderEmail = sql.NullString{String: vr.Holder.Email, Valid: true}
		rd.HolderDescription = sql.NullString{String: vr.Holder.Description, Valid: true}

		if err := rd.InsertVc(); err != nil {
			log.Panic(err)
		}
		// fmt.Println(i)
		ret[i] = rd
	}
	return ret, nil
}

func getCountVoucher(programID string) int {
	return model.CountVoucher(programID)
}

func rollback(vr string) {
	_ = model.HardDelete(vr)
}

func generateLink(id string) string {
	return model.VOUCHER_URL + "?x=" + StrEncode(id)
}

func voucherCode(vcf model.VoucherCodeFormat, flag int) string {
	var code string

	seedCode := func() string {
		return randStr(model.DEFAULT_SEED_LENGTH, model.DEFAULT_SEED_CODE)
	}

	if vcf.Prefix.Valid {
		code += vcf.Prefix.String + "-"
	}

	switch {
	case flag == 0:
		code += seedCode() + "-" + randStr(model.DEFAULT_LENGTH, model.DEFAULT_CODE)
	case vcf.Body.Valid == true && vcf.Body.String != "":
		code += seedCode() + "-" + vcf.Body.String
	default:
		code += seedCode() + "-" + randStr(vcf.Length, vcf.FormatType)
	}

	if vcf.Postfix.Valid {
		code += "-" + vcf.Postfix.String
	}

	return code
}

func SendVoucherBulk(w http.ResponseWriter, r *http.Request) {
	apiName := "voucher_generate-bulk"
	var gvd GenerateVoucherRequest
	var status int
	res := NewResponse(nil)
	vrID := r.FormValue("program")

	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag(apiName)

	a := AuthTokenWithLogger(w, r, logger)
	if !a.Valid {
		render.JSON(w, a.res, http.StatusUnauthorized)
		return
	}

	if CheckAPIRole(a, apiName) {
		logger.SetStatus(status).Info("param :", a.User.ID, "response :", "Invalid Role")

		status = http.StatusUnauthorized
		res.AddError(its(status), model.ErrCodeInvalidRole, model.ErrInvalidRole.Error(), logger.TraceID)
		render.JSON(w, res, status)
		return
	}

	if getCountVoucher(vrID) > 0 {
		status = http.StatusBadRequest
		res.AddError(its(status), model.ErrCodeInvalidProgram, model.ErrMessageProgramHasBeenUsed, logger.TraceID)
		logger.SetStatus(status).Log("param :", vrID, "response :", res.Errors)
		render.JSON(w, res, status)
		return
	}

	program, err := model.FindProgramDetailsById(vrID)
	if err == model.ErrResourceNotFound {
		status = http.StatusNotFound
		res.AddError(its(status), model.ErrCodeResourceNotFound, model.ErrMessageResourceNotFound, logger.TraceID)
		logger.SetStatus(status).Log("param :", vrID, "response :", res.Errors)
		render.JSON(w, res, status)
		return
	} else if err != nil {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", logger.TraceID)
		logger.SetStatus(status).Log("param :", vrID, "response :", res.Errors)
		render.JSON(w, res, status)
		return
	}

	var listBroadcast []model.BroadcastUser
	listBroadcast, err = model.FindBroadcastUser(map[string]string{"program_id": vrID})
	if err != nil {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", logger.TraceID)
		logger.SetStatus(status).Log("param :", vrID, "response :", res.Errors)
		render.JSON(w, res, status)
		return
	}

	gvd.AccountID = a.User.Account.Id
	gvd.ProgramID = vrID
	gvd.Quantity = 1
	gvd.CreatedBy = a.User.ID

	totalVoucher := []model.Voucher{}
	tempListVoucher := []model.Voucher{}

	for _, v := range listBroadcast {
		gvd.ReferenceNo = its(v.ID)
		gvd.Holder.Key = v.Target
		gvd.Holder.Description = v.Description

		tempListVoucher, err = gvd.generateVoucher(&program)
		if err != nil {
			fmt.Println(err)
			rollback(vrID)

			status = http.StatusInternalServerError
			res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", logger.TraceID)
			logger.SetStatus(status).Log("param :", vrID, "response :", err.Error())
			render.JSON(w, res, status)
			return
		}

		for _, vv := range tempListVoucher {
			totalVoucher = append(totalVoucher, vv)
		}
	}

	listEmail := []string{}
	listParam := []model.SedayuOneEmail{}

	for _, v := range totalVoucher {
		listEmail = append(listEmail, v.Holder.String)
		listParam = append(listParam, model.SedayuOneEmail{Name: v.HolderDescription.String, VoucherUrl: generateLink(v.ID)})
	}

	fmt.Println(listEmail)
	fmt.Println(listParam)

	if err := model.SendMailSedayuOne(model.Domain, model.ApiKey, model.PublicApiKey, "Sedayu One Voucher Test", listEmail, listParam); err != nil {
		res := NewResponse(nil)
		status := http.StatusInternalServerError
		errTitle := model.ErrCodeInternalError
		if err == model.ErrResourceNotFound {
			status = http.StatusNotFound
			errTitle = model.ErrCodeResourceNotFound
		}

		res.AddError(its(status), errTitle, err.Error(), logger.TraceID)
		logger.SetStatus(status).Info("param :", listParam, "response :", err.Error())
		render.JSON(w, res, status)
		return
	}

	status = http.StatusCreated
	res = NewResponse("success")
	logger.SetStatus(status).Log("param :", vrID, "response :", res.Data)
	render.JSON(w, res, status)
	return

}
