package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gilkor/evoucher/internal/model"
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
		AccountID   string `json:"account_id"`
		VariantID   string `json:"variant_id"`
		Quantity    int    `json:"quantity"`
		Holder      string `json:"holder"`
		ReferenceNo string `json:"reference_no"`
		CreatedBy   string `json:"user"`
	}

	GetVoucherOfVariatList []GetVoucherOfVariatdata
	GetVoucherOfVariatdata struct {
		AccountId          string           `json:"account_id"`
		VariantName        string           `json:"variant_name"`
		VariantType        string           `json:"variant_type"`
		VoucherType        string           `json:"voucher_type"`
		RedeemMethod       string           `json:"redeemtion_method"`
		VoucherPrice       float64          `json:"voucher_price"`
		DiscountValue      float64          `json:"discount_value"`
		AllowAccumulative  bool             `json:"allow_accumulative"`
		StartDate          string           `json:"start_date"`
		EndDate            string           `json:"end_date"`
		ImgUrl             string           `json:"image_url"`
		VariantTnc         string           `json:"variant_tnc"`
		VariantDescription string           `json:"variant_description"`
		Vouchers           []VoucerResponse `json:"vouchers,omitempty"`
	}

	// VoucerResponse represent list of voucher data
	VoucerResponse struct {
		VoucherID string `json:"voucher_id"`
		VoucherNo string `json:"voucher_code"`
		State     string `json:"state"`
	}
	// DetailListResponseData represent list of voucher data
	DetailListResponseData []DetailResponseData
	// DetailResponseData represent list of voucher data
	DetailResponseData struct {
		ID            string    `json:"id"`
		VoucherCode   string    `json:"voucher_code"`
		ReferenceNo   string    `json:"reference_no"`
		Holder        string    `json:"holder"`
		VariantID     string    `json:"variant_id"`
		ValidAt       time.Time `json:"valid_at"`
		ExpiredAt     time.Time `json:"expired_at"`
		DiscountValue float64   `json:"discount_value"`
		State         string    `json:"state"`
		CreatedBy     string    `json:"created_by"`
		CreatedAt     time.Time `json:"created_at"`
		UpdatedBy     string    `json:"updated_by"`
		UpdatedAt     time.Time `json:"updated_at"`
		DeletedBy     string    `json:"deleted_by"`
		DeletedAt     time.Time `json:"deleted_at"`
		Status        string    `json:"status"`
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
	}

	dt, err := model.FindVariantDetailsById(voucher.VoucherData[0].VariantID)

	sd, err := time.Parse(time.RFC3339Nano, dt.StartDate)
	if err != nil {
		return false, err
	}
	ed, err := time.Parse(time.RFC3339Nano, dt.EndDate)
	if err != nil {
		return false, err
	}

	// fmt.Println(dt.RedeemtionMethod, " vs ", r.RedeemMethod)
	if sd.After(time.Now()) && ed.After(time.Now()) {
		return false, errors.New(model.ErrCodeVoucherNotActive)
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

//MyVoucher list voucher by holder
func GetVoucherOfVariant(w http.ResponseWriter, r *http.Request) {
	var voucher model.VoucherResponse
	var err error
	res := NewResponse(nil)
	var status int

	//Token Authentocation
	_, _, _, ok := CheckToken(w, r)
	if !ok {
		return
	}

	param := getUrlParam(r.URL.String())
	viewresponse := param["list"]
	delete(param, "list")
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
	fmt.Println(voucher.VoucherData)

	distinctVariant := []string{}
	for _, v := range voucher.VoucherData {
		if !stringInSlice(v.VariantID, distinctVariant) {
			distinctVariant = append(distinctVariant, v.VariantID)
		}
	}
	// fmt.Println(distinctVariant)
	d := make(GetVoucherOfVariatList, len(distinctVariant))

	switch viewresponse {
	case "variant":
		for k, v := range distinctVariant {
			dt, _ := model.FindVariantDetailsById(v)
			d[k].AccountId = dt.AccountId
			d[k].VariantName = dt.VariantName
			d[k].VariantType = dt.VariantType
			d[k].VoucherType = dt.VoucherType
			d[k].RedeemMethod = dt.RedeemtionMethod
			d[k].VoucherPrice = dt.VoucherPrice
			d[k].DiscountValue = dt.DiscountValue
			d[k].AllowAccumulative = dt.AllowAccumulative
			d[k].StartDate = dt.StartDate
			d[k].EndDate = dt.EndDate
			d[k].ImgUrl = dt.ImgUrl
			d[k].VariantTnc = dt.VariantTnc
			d[k].VariantDescription = dt.VariantDescription
		}
	case "voucher":
		dd := make([]VoucerResponse, len(voucher.VoucherData))
		for i, val := range voucher.VoucherData {
			dd[i].VoucherID = val.ID
			dd[i].VoucherNo = val.VoucherCode
			dd[i].State = val.State
		}
		status = http.StatusOK
		res = NewResponse(dd)
		render.JSON(w, res, status)
		return
	default:
		for k, v := range distinctVariant {
			dt, _ := model.FindVariantDetailsById(v)
			d[k].AccountId = dt.AccountId
			d[k].VariantName = dt.VariantName
			d[k].VariantType = dt.VariantType
			d[k].VoucherType = dt.VoucherType
			d[k].RedeemMethod = dt.RedeemtionMethod
			d[k].VoucherPrice = dt.VoucherPrice
			d[k].DiscountValue = dt.DiscountValue
			d[k].AllowAccumulative = dt.AllowAccumulative
			d[k].StartDate = dt.StartDate
			d[k].EndDate = dt.EndDate
			d[k].ImgUrl = dt.ImgUrl
			d[k].VariantTnc = dt.VariantTnc
			d[k].VariantDescription = dt.VariantDescription

			d[k].Vouchers = make([]VoucerResponse, len(voucher.VoucherData))
			for i, val := range voucher.VoucherData {
				d[k].Vouchers[i].VoucherID = val.ID
				d[k].Vouchers[i].VoucherNo = val.VoucherCode
				d[k].Vouchers[i].State = val.State
			}
		}
	}

	// d.Vouchers = make([]VoucerResponse, len(voucher.VoucherData))
	status = http.StatusOK
	res = NewResponse(d)
	render.JSON(w, res, status)
}

// GetVoucherDetail get Voucher detail from DB
func GetVoucherDetail(w http.ResponseWriter, r *http.Request) {
	var voucher model.VoucherResponse
	var err error
	res := NewResponse(nil)
	var status int

	//Token Authentocation
	accountID, userID, _, ok := CheckToken(w, r)
	if !ok {
		return
	}
	fmt.Println(accountID, userID)

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
			dvr[i].Holder = v.Holder
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

//GenerateVoucherOnDemand Generate singgle voucher request
func GenerateVoucherOnDemand(w http.ResponseWriter, r *http.Request) {
	var gvd GenerateVoucherRequest
	var status int
	res := NewResponse(nil)

	//Token Authentocation
	accountID, userID, _, ok := CheckToken(w, r)
	if !ok {
		return
	}
	fmt.Println(accountID, userID)

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
		res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", "voucher")
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

	d := GenerateVoucherRequest{
		AccountID:   dt.AccountId,
		VariantID:   dt.Id,
		Quantity:    1,
		Holder:      gvd.Holder,
		ReferenceNo: gvd.ReferenceNo,
		CreatedBy:   userID,
	}

	var voucher []model.Voucher
	voucher, err = d.generateVoucherBulk(&dt)
	if err != nil {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", "voucher")
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
func GenerateVoucher(w http.ResponseWriter, r *http.Request) {
	var gvd GenerateVoucherRequest
	var status int
	res := NewResponse(nil)

	//Token Authentocation
	accountID, userID, _, ok := CheckToken(w, r)
	if !ok {
		return
	}
	fmt.Println(accountID, userID)

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
		res.AddError(its(status), model.ErrCodeResourceNotFound, model.ErrMessageResourceNotFound, "voucher")
		render.JSON(w, res, status)
		return
	} else if err != nil {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", "voucher")
		render.JSON(w, res, status)
		return
	}

	if (int(dt.MaxQuantityVoucher) - getCountVoucher(gvd.VariantID) - gvd.Quantity) <= 0 {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeVoucherQtyExceeded, model.ErrMessageVoucherQtyExceeded, "voucher")
		render.JSON(w, res, status)
		return
	} else if dt.VariantType != model.VariantTypeBulk {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeVoucherRulesViolated, model.ErrMessageVoucherRulesViolated, "voucher")
		render.JSON(w, res, status)
		return
	}

	d := GenerateVoucherRequest{
		AccountID:   dt.AccountId,
		VariantID:   dt.Id,
		Quantity:    gvd.Quantity,
		Holder:      gvd.Holder,
		ReferenceNo: gvd.ReferenceNo,
		CreatedBy:   userID,
	}

	var voucher []model.Voucher
	voucher, err = d.generateVoucherBulk(&dt)
	if err != nil {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeVoucherRulesViolated, model.ErrMessageVoucherRulesViolated, "voucher")
		render.JSON(w, res, status)
		return
	}

	gvr := make([]VoucerResponse, len(voucher))
	for i, v := range voucher {
		gvr[i].VoucherID = v.ID
		gvr[i].VoucherNo = v.VoucherCode
	}

	status = http.StatusCreated
	res = NewResponse(gvr)
	render.JSON(w, res, status)
	return

}

// GenerateVoucher Genera te voucher and strore to DB
func (vr *GenerateVoucherRequest) generateVoucherBulk(v *model.Variant) ([]model.Voucher, error) {
	ret := make([]model.Voucher, vr.Quantity)
	var rt []string
	var vcf model.VoucherCodeFormat
	var code string

	vcf, err := model.GetVoucherCodeFormat(v.VoucherFormat)
	if err != nil {
		return ret, err
	}

	for i := 0; i <= vr.Quantity-1; i++ {

		switch {
		case v.VoucherFormat == 0:
			code = randStr(model.DEFAULT_LENGTH, model.DEFAULT_CODE)
		case vcf.Body.Valid == true:
			code = vcf.Prefix.String + vcf.Body.String + vcf.Postfix.String
		default:
			code = vcf.Prefix.String + randStr(vcf.Length-(len(vcf.Prefix.String)+len(vcf.Postfix.String)), vcf.FormatType) + vcf.Postfix.String
			// fmt.Println("3 :", code)
		}
		rt = append(rt, code)

		tsd, err := time.Parse(time.RFC3339Nano, v.StartDate)
		if err != nil {
			log.Panic(err)
		}
		ted, err := time.Parse(time.RFC3339Nano, v.EndDate)
		if err != nil {
			log.Panic(err)
		}
		rd := model.Voucher{
			VoucherCode:   rt[i],
			ReferenceNo:   vr.ReferenceNo,
			Holder:        vr.Holder,
			VariantID:     vr.VariantID,
			ValidAt:       tsd,
			ExpiredAt:     ted,
			DiscountValue: v.DiscountValue,
			State:         model.VoucherStateCreated,
			CreatedBy:     vr.CreatedBy, //note harus nya by user
			CreatedAt:     time.Now(),
		}

		if err := rd.InsertVc(); err != nil {
			log.Panic(err)
		}
		// fmt.Println(i)
		ret[i] = rd
	}
	return ret, nil
}

func getCountVoucher(variantID string) int {
	fmt.Println(model.CountVoucher(variantID))
	return model.CountVoucher(variantID)
}
