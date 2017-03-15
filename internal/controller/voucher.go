package controller

import (
	"encoding/json"
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
		VoucherCode string `json:"voucher_code"`
		AccountID   string `json:"account_id"`
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
	}

	// GenerateVoucerResponse represent list of voucher data
	GenerateVoucerResponse struct {
		VoucherID string `json:"voucher_id"`
		VoucherNo string `json:"voucher_code"`
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

// GetVoucherDetail get Voucher detail from DB
func GetVoucherDetail(w http.ResponseWriter, r *http.Request) {
	var voucher model.VoucherResponse
	var err error
	// var vr ResponseData
	res := NewResponse(nil)
	var status int

	param := getUrlParam(r.URL.String())
	delete(param, "token")
	delete(param, "user")

	if len(param) > 0 {
		voucher, err = model.FindVoucher(param)
	} else {
		status = http.StatusBadRequest
		res.AddError("400001", model.ErrCodeMissingOrderItem, model.ErrMessageMissingOrderItem, "voucher")
		render.JSON(w, res, status)
		return
	}
	// fmt.Println(voucher, err)
	if err == model.ErrResourceNotFound {
		status = http.StatusBadRequest
		res.AddError("400002", model.ErrCodeResourceNotFound, model.ErrMessageResourceNotFound, "voucher")
		render.JSON(w, res, status)
		return
	} else if err != nil {
		status = http.StatusInternalServerError
		res.AddError("500001", model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", "voucher")
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

//RedeemVoucherValidation redeem
func (r *RedeemVoucherRequest) RedeemVoucherValidation() ResponseData {
	var d model.UpdateDeleteRequest
	var vr ResponseData

	fv, err := model.FindVoucher(map[string]string{"voucher_code": r.VoucherCode})
	if err == model.ErrResourceNotFound {
		vr.Error = model.ErrCodeResourceNotFound
		vr.Message = model.ErrResourceNotFound.Error()
	} else if err != nil {
		vr.Error = model.ResponseStateNok
		vr.Message = err.Error()
	} else if fv.Status != model.ResponseStateOk {
		vr.Error = fv.Status
		vr.Message = fv.Message
	} else {

		d.ID = fv.VoucherData[0].ID
		d.State = model.VoucherStateUsed
		d.User = r.AccountID

		uv, err := d.UpdateVc()
		if err != nil {
			vr.Error = model.ResponseStateNok
			vr.Message = err.Error()
		} else {
			vr.Error = model.ResponseStateOk
			vr.Message = ""

			dvr := DetailResponseData{
				ID:            uv.ID,
				VoucherCode:   uv.VoucherCode,
				ReferenceNo:   uv.ReferenceNo,
				Holder:        uv.Holder,
				VariantID:     uv.VariantID,
				ValidAt:       uv.ValidAt,
				ExpiredAt:     uv.ExpiredAt,
				DiscountValue: uv.DiscountValue,
				State:         uv.State,
				CreatedBy:     uv.CreatedBy,
				CreatedAt:     uv.CreatedAt,
				UpdatedBy:     uv.UpdatedBy.String,
				UpdatedAt:     uv.UpdatedAt.Time,
				DeletedBy:     uv.DeletedBy.String,
				DeletedAt:     uv.DeletedAt.Time,
				Status:        uv.Status,
			}
			vr.Data = dvr
		}
	}
	return vr
}

// DeleteVoucher delete Voucher data from DB by ID
func DeleteVoucher(w http.ResponseWriter, r *http.Request) {
	var d model.UpdateDeleteRequest
	var rd DeleteVoucherRequest
	status := http.StatusOK
	vr := ResponseData{Error: model.ResponseStateOk, Message: ""}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		log.Panic(err)
		status = http.StatusInternalServerError
	}

	fv, err := model.FindVoucher(map[string]string{"voucher_code": rd.VoucherCode})
	if err == model.ErrResourceNotFound {
		vr.Error = model.ErrCodeResourceNotFound
		vr.Message = model.ErrResourceNotFound.Error()
	} else if err != nil {
		status = http.StatusInternalServerError
		vr.Error = http.StatusText(status)
		vr.Message = err.Error()
	} else {
		d.ID = fv.VoucherData[0].ID
		d.State = model.VoucherStateDeleted
		d.User = fv.VoucherData[0].CreatedBy

		err := d.DeleteVc()
		if err != nil {
			vr.Error = model.ResponseStateNok
			vr.Message = err.Error()
		} else {
			vr.Error = model.ResponseStateOk
			vr.Message = ""
			vr.Data = fv.VoucherData[0].ID
		}
	}

	vr.State = its(status)

	res := NewResponse(vr)
	render.JSON(w, res, status)
}

// PayVoucher
func PayVoucher(w http.ResponseWriter, r *http.Request) {
	var d model.UpdateDeleteRequest
	var pvr PayVoucherRequest
	status := http.StatusOK
	vr := ResponseData{State: model.ResponseStateOk, Message: ""}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&pvr); err != nil {
		vr.Error = model.ResponseStateNok
		vr.Message = err.Error()
		status = http.StatusInternalServerError
		// log.Panic(err)
	}

	fv, err := model.FindVoucher(map[string]string{"voucher_code": pvr.VoucherCode})
	if err == model.ErrResourceNotFound {
		vr.Error = model.ErrCodeResourceNotFound
		vr.Message = model.ErrResourceNotFound.Error()
	} else if err != nil {
		vr.Error = model.ResponseStateNok
		vr.Message = err.Error()
		status = http.StatusInternalServerError
	} else if fv.VoucherData[0].State == model.VoucherStatePaid {
		vr.Error = model.ErrCodeVoucherAlreadyPaid
		vr.Message = model.ErrMessageVoucherAlreadyPaid
	} else {
		d.ID = fv.VoucherData[0].ID
		d.State = model.VoucherStatePaid
		d.User = fv.VoucherData[0].CreatedBy

		uv, err := d.UpdateVc()
		if err != nil {
			vr.Error = model.ResponseStateNok
			vr.Message = err.Error()
		} else {
			vr.Error = model.ResponseStateOk
			vr.Message = ""
			dvr := DetailResponseData{
				ID:            uv.ID,
				VoucherCode:   uv.VoucherCode,
				ReferenceNo:   uv.ReferenceNo,
				Holder:        uv.Holder,
				VariantID:     uv.VariantID,
				ValidAt:       uv.ValidAt,
				ExpiredAt:     uv.ExpiredAt,
				DiscountValue: uv.DiscountValue,
				State:         uv.State,
				CreatedBy:     uv.CreatedBy,
				CreatedAt:     uv.CreatedAt,
				UpdatedBy:     uv.UpdatedBy.String,
				UpdatedAt:     uv.UpdatedAt.Time,
				DeletedBy:     uv.DeletedBy.String,
				DeletedAt:     uv.DeletedAt.Time,
				Status:        uv.Status,
			}
			vr.Data = dvr
		}
	}

	res := NewResponse(vr)
	render.JSON(w, res, status)
}

//GenerateVoucherOnDemand Generate singgle voucher request
func GenerateVoucherOnDemand(w http.ResponseWriter, r *http.Request) {
	var gvd GenerateVoucherRequest
	var status int
	res := NewResponse(nil)

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&gvd); err != nil {
		status = http.StatusBadRequest
		res.AddError("400002", its(status), http.StatusText(status)+"("+err.Error()+")", "voucher")
		render.JSON(w, res, status)
		return
	}

	variant, err := model.FindVariantById(gvd.VariantID)
	if err == model.ErrResourceNotFound {
		status = http.StatusBadRequest
		res.AddError("400002", model.ErrCodeResourceNotFound, model.ErrMessageResourceNotFound, "voucher")
		render.JSON(w, res, status)
		return
	} else if err != nil {
		status = http.StatusInternalServerError
		res.AddError("500002", its(status), http.StatusText(status)+"("+err.Error()+")", "voucher")
		render.JSON(w, res, status)
		return
	}

	if dt, ok := variant.Data.(model.Variant); !ok {
		status = http.StatusInternalServerError
		res.AddError("500001", model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", "voucher")
		render.JSON(w, res, status)
		return
	} else {
		if (int(dt.MaxQuantityVoucher) - getCountVoucher(gvd.VariantID) - 1) <= 0 {
			status = http.StatusInternalServerError
			res.AddError("500003", model.ErrCodeVoucherQtyExceeded, model.ErrMessageVoucherQtyExceeded, "voucher")
			render.JSON(w, res, status)
			return
		} else if dt.VariantType != model.VariantTypeOnDemand {
			status = http.StatusInternalServerError
			res.AddError("500004", model.ErrCodeVoucherRulesViolated, model.ErrMessageVoucherRulesViolated, "voucher")
			render.JSON(w, res, status)
			return
		}

		d := GenerateVoucherRequest{
			AccountID:   dt.AccountId,
			VariantID:   dt.Id,
			Quantity:    1,
			Holder:      gvd.Holder,
			ReferenceNo: gvd.ReferenceNo,
		}

		var voucher []model.Voucher
		voucher, err = d.generateVoucherBulk(&dt)
		if err != nil {
			status = http.StatusInternalServerError
			res.AddError("500005", its(status), http.StatusText(status)+"("+err.Error()+")", "voucher")
			render.JSON(w, res, status)
			return
		}

		gvr := make([]GenerateVoucerResponse, len(voucher))
		for i, v := range voucher {
			gvr[i].VoucherID = v.ID
			gvr[i].VoucherNo = v.VoucherCode
		}

		status = http.StatusCreated
		res = NewResponse(gvr)
		render.JSON(w, res, status)
		return
	}
}

//GenerateVoucher Generate bulk voucher request
func GenerateVoucher(w http.ResponseWriter, r *http.Request) {
	var gvd GenerateVoucherRequest
	var status int
	res := NewResponse(nil)

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&gvd); err != nil {
		status = http.StatusBadRequest
		res.AddError("400002", its(status), http.StatusText(status)+"("+err.Error()+")", "voucher")
		render.JSON(w, res, status)
		return
	}

	variant, err := model.FindVariantById(gvd.VariantID)
	if err == model.ErrResourceNotFound {
		status = http.StatusBadRequest
		res.AddError("400002", model.ErrCodeResourceNotFound, model.ErrMessageResourceNotFound, "voucher")
		render.JSON(w, res, status)
		return
	} else if err != nil {
		status = http.StatusInternalServerError
		res.AddError("500002", its(status), http.StatusText(status)+"("+err.Error()+")", "voucher")
		render.JSON(w, res, status)
		return
	}

	if dt, ok := variant.Data.(model.Variant); !ok {
		status = http.StatusInternalServerError
		res.AddError("500001", model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", "voucher")
		render.JSON(w, res, status)
		return
	} else {
		if (int(dt.MaxQuantityVoucher) - getCountVoucher(gvd.VariantID) - 1) <= 0 {
			status = http.StatusInternalServerError
			res.AddError("500003", model.ErrCodeVoucherQtyExceeded, model.ErrMessageVoucherQtyExceeded, "voucher")
			render.JSON(w, res, status)
			return
		} else if dt.VariantType != model.VariantTypeBulk {
			status = http.StatusInternalServerError
			res.AddError("500004", model.ErrCodeVoucherRulesViolated, model.ErrMessageVoucherRulesViolated, "voucher")
			render.JSON(w, res, status)
			return
		}

		d := GenerateVoucherRequest{
			AccountID:   dt.AccountId,
			VariantID:   dt.Id,
			Quantity:    1,
			Holder:      gvd.Holder,
			ReferenceNo: gvd.ReferenceNo,
		}

		var voucher []model.Voucher
		voucher, err = d.generateVoucherBulk(&dt)
		if err != nil {
			status = http.StatusInternalServerError
			res.AddError("500005", its(status), http.StatusText(status)+"("+err.Error()+")", "voucher")
			render.JSON(w, res, status)
			return
		}

		gvr := make([]GenerateVoucerResponse, len(voucher))
		for i, v := range voucher {
			gvr[i].VoucherID = v.ID
			gvr[i].VoucherNo = v.VoucherCode
		}

		status = http.StatusCreated
		res = NewResponse(gvr)
		render.JSON(w, res, status)
		return
	}
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
			Holder:        vr.AccountID,
			VariantID:     vr.VariantID,
			ValidAt:       tsd,
			ExpiredAt:     ted,
			DiscountValue: v.DiscountValue,
			State:         model.VoucherStateCreated,
			CreatedBy:     vr.AccountID, //note harus nya by user
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
