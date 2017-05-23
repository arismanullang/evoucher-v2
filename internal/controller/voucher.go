package controller

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

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
		State 		   string	    `json:"state"`
		Holder		   string	    `json:"holder"`
		HolderDescription  string	    `json:"holder_description"`
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
	GetVoucherlinkdata struct {
		Url string `json:"url"`
		VoucherID  string `json:"voucher_id"`
		VoucherCode string `json:"voucher_code"`
	}
)

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

//GetVoucherOfVariant list voucher by holder
func GetVoucherOfVariant(w http.ResponseWriter, r *http.Request) {
	var voucher model.VoucherResponse
	var err error
	res := NewResponse(nil)
	var status int

	//Token Authentocation
	_, _, _, ok := AuthToken(w, r)
	if !ok {
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
	}
	// fmt.Println(voucher.VoucherData)

	distinctVariant := []string{}
	for _, v := range voucher.VoucherData {
		if !stringInSlice(v.VariantID, distinctVariant) {
			distinctVariant = append(distinctVariant, v.VariantID)
		}
	}
	// fmt.Println(distinctVariant)
	d := make(GetVoucherOfVariatList, len(distinctVariant))
	for k, v := range distinctVariant {
		dt, _ := model.FindVariantDetailsById(v)
		d[k].VariantID = dt.Id
		d[k].AccountId = dt.AccountId
		d[k].VariantName = dt.VariantName
		d[k].VoucherType = dt.VoucherType
		d[k].VoucherPrice = dt.VoucherPrice
		d[k].DiscountValue = dt.DiscountValue
		d[k].StartDate = dt.StartDate
		d[k].EndDate = dt.EndDate
		d[k].ImgUrl = dt.ImgUrl
		d[k].Used = getCountVoucher(v)
	}

	// d.Vouchers = make([]VoucerResponse, len(voucher.VoucherData))
	status = http.StatusOK
	res = NewResponse(d)
	render.JSON(w, res, status)
}

func GetVoucherOfVariantDetails(w http.ResponseWriter, r *http.Request) {
	variant := bone.GetValue(r, "id")
	res := NewResponse(nil)
	var status int

	//Token Authentocation
	_, _, _, ok := AuthToken(w, r)
	if !ok {
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

	d.Partners = make([]Partner, len(p))
	for i, pd := range p {
		d.Partners[i].ID = pd.Id
		d.Partners[i].PartnerName = pd.PartnerName
		d.Partners[i].SerialNumber = pd.SerialNumber.String
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
	accountID, userID, _, ok := AuthToken(w, r)
	if !ok {
		return
	}
	fmt.Println("auth result => ", accountID, userID)

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

	_, _, _, ok := AuthToken(w, r)
	if !ok {
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

	//Token Authentocation
	accountID, userID, _, ok := AuthToken(w, r)
	if !ok {
		return
	}
	fmt.Println("auth result => ", accountID, userID)

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

	if (int(dt.MaxQuantityVoucher) - getCountVoucher(gvd.VariantID) - 1) <= 0 {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeVoucherQtyExceeded, model.ErrMessageVoucherQtyExceeded, "voucher")
		render.JSON(w, res, status)
		return
	} else if dt.VariantType != model.VariantTypeOnDemand {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeVoucherRulesViolated, model.ErrMessageVoucherRulesViolated, "voucher")
		render.JSON(w, res, status)
		return
	}

	gvd.AccountID = dt.AccountId
	gvd.VariantID = dt.Id
	gvd.Quantity = 1
	gvd.CreatedBy = userID

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
	var gvd GenerateVoucherRequest
	var status int
	res := NewResponse(nil)
	vrID := r.FormValue("variant")
	fmt.Println("variant id = ", vrID)
	//Token Authentocation
	accountID, userID, _, ok := AuthToken(w, r)
	if !ok {
		return
	}
	fmt.Println("auth result => ", accountID, userID)

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

	gvd.AccountID = accountID
	gvd.VariantID = vrID
	gvd.Quantity = 1
	gvd.CreatedBy = userID

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

// GenerateVoucher Genera te voucher and strore to DB
func (vr *GenerateVoucherRequest) generateVoucher(v *model.Variant) ([]model.Voucher, error) {
	ret := make([]model.Voucher, vr.Quantity)
	var rt []string
	var vcf model.VoucherCodeFormat
	var code string
	var tsd , ted time.Time

	vcf, err := model.GetVoucherCodeFormat(v.VoucherFormat)
	if err != nil {
		return ret, err
	}
	if v.VoucherLifetime > 0 {
		tsd = time.Now()
		ted = time.Now().AddDate(0,0,v.VoucherLifetime)
	}else {
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
		switch {
		case v.VoucherFormat == 0:
			code = randStr(model.DEFAULT_LENGTH, model.DEFAULT_CODE)
		case vcf.Body.Valid == true && vcf.Body.String != "":
			code = vcf.Prefix.String + vcf.Body.String + vcf.Postfix.String
		default:
			code = vcf.Prefix.String + randStr(vcf.Length-(len(vcf.Prefix.String)+len(vcf.Postfix.String)), vcf.FormatType) + vcf.Postfix.String
		}
		rt = append(rt, code)

		// fmt.Println("generate data =>", vr.Holder)
		rd := model.Voucher{
			VoucherCode:   rt[i],
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


func GetVoucherlink(w http.ResponseWriter, r *http.Request) {
	status := http.StatusOK
	res := NewResponse(nil)
	varID := r.FormValue("variant")

	_, _, _, ok := AuthToken(w, r)
	if !ok {
		return
	}

	v,err := model.FindVoucher(map[string]string{"variant_id":varID})
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

	vl := make(GetVoucherlinkResponse, len(v.VoucherData))
	for k,v := range v.VoucherData{
		vl[k].Url = generateLink(v.ID)
		vl[k].VoucherID = v.ID
		vl[k].VoucherCode = v.VoucherCode
	}

	res = NewResponse(vl)
	render.JSON(w, res, status)
	return
}

func getCountVoucher(variantID string) int {
	fmt.Println(model.CountVoucher(variantID))
	return model.CountVoucher(variantID)
}

func rollback(vr string) {
	_ = model.HardDelete(vr)
}

func generateLink(id string) string{
	return model.VOUCHER_URL+"?x="+StrEncode(id)
}
