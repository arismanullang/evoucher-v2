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

	"github.com/gilkor/evoucher/internal/model"
	"github.com/go-zoo/bone"
	"github.com/ruizu/render"
	"github.com/go-ozzo/ozzo-validation"
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
		AccountID string `json:"account_id"`
		VariantID string `json:"variant_id"`
		Quantity  int    `json:"quantity"`
		Holder    struct {
			Key         string `json:"id"`
			Phone       string `json:"phone"`
			Email       string `json:"email"`
			Description string `json:"description"`
		} `json:"holder"`
		ReferenceNo string `json:"reference_no"`
		CreatedBy   string `json:"user"`
	}

	GetVoucherOfVariatList []GetVoucherOfVariatdata
	GetVoucherOfVariatdata struct {
		VariantID     string  `json:"variant_id"`
		AccountId     string  `json:"account_id"`
		VariantName   string  `json:"variant_name"`
		VoucherType   string  `json:"voucher_type"`
		VoucherPrice  float64 `json:"voucher_price"`
		DiscountValue float64 `json:"discount_value"`
		MaxQty        float64 `json:"max_quantity_voucher"`
		ImgUrl        string  `json:"image_url"`
		StartDate     string  `json:"start_date"`
		EndDate       string  `json:"end_date"`
		Used          int     `json:"used"`
	}

	GetVoucherOfVariatListDetails struct {
		VariantID          string           `json:"variant_id"`
		AccountId          string           `json:"account_id"`
		VariantName        string           `json:"variant_name"`
		VariantType        string           `json:"variant_type"`
		VoucherFormat      int              `json:"voucher_format_id"`
		VoucherType        string           `json:"voucher_type"`
		VoucherPrice       float64          `json:"voucher_price"`
		AllowAccumulative  bool             `json:"allow_accumulative"`
		StartDate          string           `json:"start_date"`
		EndDate            string           `json:"end_date"`
		DiscountValue      float64          `json:"discount_value"`
		MaxQuantityVoucher float64          `json:"max_quantity_voucher"`
		MaxUsageVoucher    float64          `json:"max_usage_voucher"`
		RedeemtionMethod   string           `json:"redeemtion_method"`
		ImgUrl             string           `json:"image_url"`
		VariantTnc         string           `json:"variant_tnc"`
		VariantDescription string           `json:"variant_description"`
		CreatedBy          string           `json:"created_by"`
		CreatedAt          string           `json:"created_at"`
		Used               int              `json:"used"`
		State              string           `json:"state"`
		Holder             string           `json:"holder"`
		HolderDescription  string           `json:"holder_description"`
		Partners           []Partner        `json:"Partners"`
		Voucher            []VoucerResponse `json:"Vouchers"`
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
		VariantID         string    `json:"variant_id"`
		ValidAt           time.Time `json:"valid_at"`
		ExpiredAt         time.Time `json:"expired_at"`
		DiscountValue     float64   `json:"discount_value"`
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
		VariantID         string    `json:"variant_id"`
		ValidAt           time.Time `json:"valid_at"`
		ExpiredAt         time.Time `json:"expired_at"`
		DiscountValue     float64   `json:"discount_value"`
		State             string    `json:"state"`
		CreatedBy         string    `json:"created_by"`
		CreatedAt         time.Time `json:"created_at"`
		UpdatedBy         string    `json:"updated_by"`
		UpdatedAt         time.Time `json:"updated_at"`
		DeletedBy         string    `json:"deleted_by"`
		DeletedAt         time.Time `json:"deleted_at"`
		Status            string    `json:"status"`

		VariantName  string  `json:"Variant_name"`
		VoucherPrice float64 `json:"Voucher_price"`
		VoucherType  string  `json:"Voucher_type"`
	}

	GetVoucherlinkResponse []GetVoucherlinkdata
	GetVoucherlinkdata     struct {
		Url         string `json:"url"`
		VoucherID   string `json:"voucher_id"`
		VoucherCode string `json:"voucher_code"`
	}
)

// ## API ##//

//GetVoucherOfVariant list voucher by holder
func GetVoucherOfVariant(w http.ResponseWriter, r *http.Request) {
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

	param := getUrlParam(r.URL.String())
	delete(param, "token")

	if len(param) > 0 {
		voucher, err = model.FindAvailableVoucher(param)
	} else {
		status = http.StatusBadRequest
		res.AddError(its(status), model.ErrCodeMissingOrderItem, model.ErrMessageMissingOrderItem, "voucher")
		render.JSON(w, res, status)
		return
	}
	// fmt.Println(voucher, err)
	if err == model.ErrResourceNotFound {
		status = http.StatusNotFound
		res.AddError(its(status), model.ErrCodeResourceNotFound, model.ErrMessageResourceNotFound, "voucher")
		render.JSON(w, res, status)
		return
	} else if err != nil {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", "voucher")
		render.JSON(w, res, status)
		return
	}
	// fmt.Println(voucher.VoucherData)

	distinctVariant := []string{}
	for _, v := range voucher.VoucherData {
		if !stringInSlice(v.VariantID, distinctVariant) {
			distinctVariant = append(distinctVariant, v.VariantID)
		}
	}
	// fmt.Println(distinctVariant)
	d := []GetVoucherOfVariatListDetails{}
	for _, v := range distinctVariant {
		tempVoucherResponse := []VoucerResponse{}
		tempGetVoucherOfVariatListDetails := GetVoucherOfVariatListDetails{}
		tempPartners := []Partner{}

		dt, _ := model.FindVariantDetailsById(v)
		tempGetVoucherOfVariatListDetails.VariantID = dt.Id
		tempGetVoucherOfVariatListDetails.AccountId = dt.AccountId
		tempGetVoucherOfVariatListDetails.VariantName = dt.VariantName
		tempGetVoucherOfVariatListDetails.VoucherType = dt.VoucherType
		tempGetVoucherOfVariatListDetails.VoucherPrice = dt.VoucherPrice
		tempGetVoucherOfVariatListDetails.DiscountValue = dt.DiscountValue
		tempGetVoucherOfVariatListDetails.StartDate = dt.StartDate
		tempGetVoucherOfVariatListDetails.EndDate = dt.EndDate
		tempGetVoucherOfVariatListDetails.ImgUrl = dt.ImgUrl
		tempGetVoucherOfVariatListDetails.AllowAccumulative = dt.AllowAccumulative
		tempGetVoucherOfVariatListDetails.RedeemtionMethod = dt.RedeemtionMethod
		tempGetVoucherOfVariatListDetails.VariantTnc = dt.VariantTnc
		tempGetVoucherOfVariatListDetails.VariantDescription = dt.VariantDescription
		tempGetVoucherOfVariatListDetails.MaxQuantityVoucher = dt.MaxQuantityVoucher
		for _, vv := range voucher.VoucherData {
			if vv.VariantID == v {
				tempVoucher := VoucerResponse{
					VoucherID: vv.ID,
					VoucherNo: vv.VoucherCode,
					State:     vv.State,
				}
				tempVoucherResponse = append(tempVoucherResponse, tempVoucher)
			}
		}
		tempGetVoucherOfVariatListDetails.Voucher = tempVoucherResponse

		pv, _ := model.FindVariantPartners(v)
		for _, vv := range pv {
			tempPartner := Partner{
				ID:           vv.Id,
				PartnerName:  vv.PartnerName,
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
	render.JSON(w, res, status)
}

//GetVoucherOfVariantDetails voucher by holder
func GetVoucherOfVariantDetails(w http.ResponseWriter, r *http.Request) {
	variant := bone.GetValue(r, "id")
	res := NewResponse(nil)
	var status int

	//Token Authentocation
	a := AuthToken(w, r)
	if !a.Valid {
		render.JSON(w, a.res, http.StatusUnauthorized)
		return
	}

	param := getUrlParam(r.URL.String())
	param["state"] = model.VoucherStateCreated
	param["variant_id"] = variant
	delete(param, "token")

	if len(param) < 0 || variant == "" {
		status = http.StatusBadRequest
		res.AddError(its(status), model.ErrCodeMissingOrderItem, model.ErrMessageMissingOrderItem, "voucher")
		render.JSON(w, res, status)
		return
	}

	voucher, err := model.FindVoucher(param)
	if err == model.ErrResourceNotFound {
		status = http.StatusNotFound
		res.AddError(its(status), model.ErrCodeResourceNotFound, model.ErrMessageInvalidHolder, "voucher")
		render.JSON(w, res, status)
		return
	} else if err != nil {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", "voucher")
		render.JSON(w, res, status)
		return
	}
	dt, err := model.FindVariantDetailsById(variant)
	if err == model.ErrResourceNotFound {
		status = http.StatusNotFound
		res.AddError(its(status), model.ErrCodeResourceNotFound, model.ErrMessageInvalidVariant, "voucher")
		render.JSON(w, res, status)
		return
	} else if err != nil {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", "voucher")
		render.JSON(w, res, status)
		return
	}
	p, err := model.FindVariantPartner(map[string]string{"variant_id": variant})
	if err == model.ErrResourceNotFound {
		status = http.StatusNotFound
		res.AddError(its(status), model.ErrCodeResourceNotFound, model.ErrMessageInvalidVariant+"(Partner of Variant Not Found)", "voucher")
		render.JSON(w, res, status)
		return
	} else if err != nil {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", "voucher")
		render.JSON(w, res, status)
		return
	}

	d := GetVoucherOfVariatListDetails{}
	d.VariantID = dt.Id
	d.AccountId = dt.AccountId
	d.VariantName = dt.VariantName
	d.VariantType = dt.VariantType
	d.VariantTnc = dt.VariantTnc
	d.CreatedBy = dt.CreatedBy
	d.CreatedAt = dt.CreatedAt
	d.VoucherType = dt.VoucherType
	d.VoucherType = dt.VoucherType
	d.VoucherPrice = dt.VoucherPrice
	d.AllowAccumulative = dt.AllowAccumulative
	d.StartDate = dt.StartDate
	d.EndDate = dt.EndDate
	d.DiscountValue = dt.DiscountValue
	d.MaxQuantityVoucher = dt.MaxQuantityVoucher
	d.MaxUsageVoucher = dt.MaxUsageVoucher
	d.RedeemtionMethod = dt.RedeemtionMethod
	d.ImgUrl = dt.ImgUrl
	d.VariantTnc = dt.VariantTnc
	d.VariantDescription = dt.VariantDescription
	d.Used = getCountVoucher(dt.Id)

	d.Partners = make([]Partner, len(p))
	for i, pd := range p {
		d.Partners[i].ID = pd.Id
		d.Partners[i].PartnerName = pd.PartnerName
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

	// d.Vouchers = make([]VoucerResponse, len(voucher.VoucherData))
	status = http.StatusOK
	res = NewResponse(d)
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

	param := getUrlParam(r.URL.String())
	delete(param, "token")

	if len(param) > 0 {
		voucher, err = model.FindVoucher(param)
	} else {
		status = http.StatusBadRequest
		res.AddError(its(status), model.ErrCodeMissingOrderItem, model.ErrMessageMissingOrderItem, "voucher")
		render.JSON(w, res, status)
		return
	}
	// fmt.Println(voucher, err)
	if err == model.ErrResourceNotFound {
		status = http.StatusNotFound
		res.AddError(its(status), model.ErrCodeResourceNotFound, model.ErrMessageResourceNotFound, "voucher")
		render.JSON(w, res, status)
		return
	} else if err != nil {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", "voucher")
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
			dvr[i].VariantID = v.VariantID
			dvr[i].ValidAt = v.ValidAt
			dvr[i].ExpiredAt = v.ExpiredAt
			dvr[i].DiscountValue = v.DiscountValue
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

	d, err := model.FindVoucher(map[string]string{"id": vc})

	if err == model.ErrResourceNotFound {
		status = http.StatusNotFound
		res.AddError(its(status), model.ErrCodeResourceNotFound, model.ErrMessageResourceNotFound, "voucher")
		render.JSON(w, res, status)
		return
	} else if err != nil {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", "voucher")
		render.JSON(w, res, status)
		return
	}

	dt, err := model.FindVariantDetailsById(d.VoucherData[0].VariantID)
	if err == model.ErrResourceNotFound {
		status = http.StatusNotFound
		res.AddError(its(status), model.ErrCodeResourceNotFound, model.ErrMessageInvalidVariant, "voucher")
		render.JSON(w, res, status)
		return
	} else if err != nil {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", "voucher")
		render.JSON(w, res, status)
		return
	}
	// p, err := model.FindVariantPartner(map[string]string{"variant_id": d.VoucherData[0].VariantID})
	// if err == model.ErrResourceNotFound {
	// 	status = http.StatusNotFound
	// 	res.AddError(its(status), model.ErrCodeResourceNotFound, model.ErrMessageInvalidVariant+"(Partner of Variant Not Found)", "voucher")
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
		VariantID:         d.VoucherData[0].VariantID,
		ValidAt:           d.VoucherData[0].ValidAt,
		ExpiredAt:         d.VoucherData[0].ExpiredAt,
		DiscountValue:     d.VoucherData[0].DiscountValue,
		State:             d.VoucherData[0].State,
		CreatedBy:         d.VoucherData[0].CreatedBy,
		CreatedAt:         d.VoucherData[0].CreatedAt,
		UpdatedBy:         d.VoucherData[0].UpdatedBy.String,
		UpdatedAt:         d.VoucherData[0].UpdatedAt.Time,
		DeletedBy:         d.VoucherData[0].DeletedBy.String,
		DeletedAt:         d.VoucherData[0].DeletedAt.Time,
		Status:            d.VoucherData[0].Status,
		VariantName:       dt.VariantName,
		VoucherPrice:      dt.VoucherPrice,
		VoucherType:       dt.VoucherType,
	}

	status = http.StatusOK
	res = NewResponse(dvr)
	render.JSON(w, res, status)
	return
}

//GenerateVoucherOnDemand Generate singgle voucher request
func GenerateVoucherOnDemand(w http.ResponseWriter, r *http.Request) {
	var gvd GenerateVoucherRequest
	var status int
	res := NewResponse(nil)

	//err := gvd.Validate()
	//if err !=nil {
	//	status = http.StatusBadRequest
	//	res.AddError(its(status), model.ErrCodeValidationError, model.ErrMessageValidationError+"("+err.Error()+")", "transaction")
	//	render.JSON(w, res, status)
	//	return
	//}

	//Token Authentocation
	a := AuthToken(w, r)
	if !a.Valid {
		render.JSON(w, a.res, http.StatusUnauthorized)
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&gvd); err != nil {
		status = http.StatusBadRequest
		res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", "voucher")
		render.JSON(w, res, status)
		return
	}

	dt, err := model.FindVariantDetailsById(gvd.VariantID)
	if err == model.ErrResourceNotFound {
		status = http.StatusNotFound
		res.AddError(its(status), model.ErrCodeResourceNotFound, model.ErrMessageInvalidVariant, "voucher")
		render.JSON(w, res, status)
		return
	} else if err != nil {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInvalidVariant+"("+err.Error()+")", "voucher")
		render.JSON(w, res, status)
		return
	}
	sd, err := time.Parse(time.RFC3339Nano, dt.StartDate)
	ed, err := time.Parse(time.RFC3339Nano, dt.EndDate)
	if err != nil {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInvalidVariant+"("+err.Error()+")", "voucher")
		render.JSON(w, res, status)
		return
	}

	if (int(dt.MaxQuantityVoucher) - getCountVoucher(gvd.VariantID) - 1) < 0 {
		status = http.StatusBadRequest
		res.AddError(its(status), model.ErrCodeVoucherQtyExceeded, model.ErrMessageVoucherQtyExceeded, "voucher")
		render.JSON(w, res, status)
		return
	} else if dt.VariantType != model.VariantTypeOnDemand {
		status = http.StatusBadRequest
		res.AddError(its(status), model.ErrCodeVoucherRulesViolated, model.ErrMessageVoucherRulesViolated, "voucher")
		render.JSON(w, res, status)
		return
	} else if !sd.Before(time.Now()) {
		status = http.StatusBadRequest
		res.AddError(its(status), model.ErrCodeVoucherNotActive, model.ErrMessageVoucherNotActive, "voucher")
		render.JSON(w, res, status)
		return
	} else if !ed.After(time.Now()) {
		status = http.StatusBadRequest
		res.AddError(its(status), model.ErrCodeVoucherExpired, model.ErrMessageVoucherExpired, "voucher")
		render.JSON(w, res, status)
		return
	}

	gvd.AccountID = dt.AccountId
	gvd.VariantID = dt.Id
	gvd.Quantity = 1
	gvd.CreatedBy = a.User.ID

	// fmt.Println("request data =>", gvd.Holder)
	var voucher []model.Voucher
	voucher, err = gvd.generateVoucher(&dt)
	if err != nil {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"( failed Genarate Voucher :"+err.Error()+")", "voucher")
		render.JSON(w, res, status)
		return
	}

	gvr := VoucerResponse{
		VoucherID: voucher[0].ID,
		VoucherNo: voucher[0].VoucherCode,
	}

	status = http.StatusCreated
	res = NewResponse(gvr)
	render.JSON(w, res, status)
	return

}

//GenerateVoucher Generate bulk voucher request
func GenerateVoucherBulk(w http.ResponseWriter, r *http.Request) {
	apiName := "voucher_generate-bulk"
	valid := false
	var gvd GenerateVoucherRequest
	var status int
	res := NewResponse(nil)
	vrID := r.FormValue("variant")
	fmt.Println("variant id = ", vrID)
	//Token Authentocation
	a := AuthToken(w, r)
	if !a.Valid {

		render.JSON(w, a.res, http.StatusUnauthorized)
		return

	}

	for _, valueRole := range a.User.Role {
		features := model.ApiFeatures[valueRole.RoleDetail]
		for _, valueFeature := range features {
			if apiName == valueFeature {
				valid = true
			}
		}
	}

	if valid {
		render.JSON(w, model.ErrCodeInvalidRole, http.StatusUnauthorized)
		return
	}

	if getCountVoucher(vrID) > 0 {
		status = http.StatusBadRequest
		res.AddError(its(status), model.ErrCodeInvalidVariant, model.ErrMessageVariantHasBeenUsed, "voucher")
		render.JSON(w, res, status)
		return
	}

	variant, err := model.FindVariantDetailsById(vrID)
	if err == model.ErrResourceNotFound {
		status = http.StatusNotFound
		res.AddError(its(status), model.ErrCodeResourceNotFound, model.ErrMessageResourceNotFound, "variant")
		render.JSON(w, res, status)
		return
	} else if err != nil {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", "variant")
		render.JSON(w, res, status)
		return
	}

	var listBroadcast []model.BroadcastUser
	listBroadcast, err = model.FindBroadcastUser(map[string]string{"variant_id": vrID})
	if err != nil {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", "broadcast")
		render.JSON(w, res, status)
		return
	}

	gvd.AccountID = a.User.AccountID
	gvd.VariantID = vrID
	gvd.Quantity = 1
	gvd.CreatedBy = a.User.ID

	for _, v := range listBroadcast {

		gvd.ReferenceNo = its(v.ID)
		gvd.Holder.Key = v.Description
		gvd.Holder.Description = v.BroadcastTarget

		_, err = gvd.generateVoucher(&variant)
		if err != nil {
			fmt.Println(err)
			rollback(vrID)

			status = http.StatusInternalServerError
			res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", "voucher")
			render.JSON(w, res, status)
			return
		}
	}

	status = http.StatusCreated
	res = NewResponse("success")
	render.JSON(w, res, status)
	return

}

func GetVoucherlink(w http.ResponseWriter, r *http.Request) {
	apiName := "voucher_get-link"
	valid := false

	status := http.StatusOK
	res := NewResponse(nil)
	varID := r.FormValue("variant")

	a := AuthToken(w, r)
	if !a.Valid {
		render.JSON(w, a.res, http.StatusUnauthorized)
		return
	}

	for _, valueRole := range a.User.Role {
		features := model.ApiFeatures[valueRole.RoleDetail]
		for _, valueFeature := range features {
			if apiName == valueFeature {
				valid = true
			}
		}
	}

	if valid {
		render.JSON(w, model.ErrCodeInvalidRole, http.StatusUnauthorized)
		return
	}

	v, err := model.FindVoucher(map[string]string{"variant_id": varID})
	if err == model.ErrResourceNotFound {
		status = http.StatusNotFound
		res.AddError(its(status), model.ErrCodeResourceNotFound, model.ErrMessageInvalidHolder, "voucher")
		render.JSON(w, res, status)
		return
	} else if err != nil {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", "voucher")
		render.JSON(w, res, status)
		return
	}

	//vl := make(GetVoucherlinkResponse, len(v.VoucherData))
	//for k, v := range v.VoucherData {
	//	vl[k].Url = generateLink(v.ID)
	//	vl[k].VoucherID = v.ID
	//	vl[k].VoucherCode = v.VoucherCode
	//}

	//res = NewResponse(vl)
	//render.JSON(w, res, status)

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

	w.Header().Set("Content-Disposition", "attachment;filename=Report.csv")
	w.Write(b.Bytes())
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
		tempArray := []string{"index", "Email", "Name"}
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

func (gv GenerateVoucherRequest) Validate() error{
	return validation.ValidateStruct(&gv,
		validation.Field(&gv.VariantID, validation.Required),
		validation.Field(&gv.AccountID, validation.Skip),
		validation.Field(&gv.Quantity, validation.Skip),
		validation.Field(&gv.Holder, validation.Skip),
		validation.Field(&gv.ReferenceNo, validation.Required,validation.Length(1,64)),
		validation.Field(&gv.Holder.Key, validation.Required),
		validation.Field(&gv.Holder.Phone, validation.Skip,validation.Length(0,8)),
		validation.Field(&gv.Holder.Email, validation.Skip),
		validation.Field(&gv.Holder.Description, validation.Skip,validation.Length(0,64)),
		validation.Field(&gv.CreatedBy, validation.Skip),
	)
}

//CheckVoucherRedeemtion validation
func (r *TransactionRequest) CheckVoucherRedeemtion(voucherID string) (bool, error) {

	voucher, err := model.FindVoucher(map[string]string{"id": voucherID})

	if err != nil {
		return false, err
	} else if voucher.VoucherData[0].State == model.VoucherStateUsed {
		return false, errors.New(model.ErrMessageVoucherAlreadyUsed)
	} else if voucher.VoucherData[0].State == model.VoucherStatePaid {
		return false, errors.New(model.ErrMessageVoucherAlreadyPaid)
	} else if !voucher.VoucherData[0].ExpiredAt.After(time.Now()) {
		fmt.Println("expired date : ", voucher.VoucherData[0].ExpiredAt , voucher.VoucherData[0].ID)
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

func (vr *GenerateVoucherRequest) generateVoucher(v *model.Variant) ([]model.Voucher, error) {
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
			VoucherCode:   code[i],
			ReferenceNo:   vr.ReferenceNo,
			VariantID:     vr.VariantID,
			ValidAt:       tsd,
			ExpiredAt:     ted,
			DiscountValue: v.DiscountValue,
			State:         model.VoucherStateCreated,
			CreatedBy:     vr.CreatedBy, //note harus nya by user
			CreatedAt:     time.Now(),
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

func getCountVoucher(variantID string) int {
	return model.CountVoucher(variantID)
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
