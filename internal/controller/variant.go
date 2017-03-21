package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-zoo/bone"
	"github.com/ruizu/render"

	"github.com/gilkor/evoucher/internal/model"
)

type (
	VariantReq struct {
		ReqData Variant `json:"variant"`
		User    string  `json:"created_by"`
	}
	Variant struct {
		VariantName        string    `json:"variant_name"`
		VariantType        string    `json:"variant_type"`
		VoucherFormat      FormatReq `json:"voucher_format"`
		VoucherType        string    `json:"voucher_type"`
		VoucherPrice       float64   `json:"voucher_price"`
		AllowAccumulative  bool      `json:"allow_accumulative"`
		StartDate          string    `json:"start_date"`
		EndDate            string    `json:"end_date"`
		DiscountValue      float64   `json:"discount_value"`
		MaxQuantityVoucher float64   `json:"max_quantity_voucher"`
		MaxUsageVoucher    float64   `json:"max_usage_voucher"`
		RedeemtionMethod   string    `json:"redeemtion_method"`
		ImgUrl             string    `json:"image_url"`
		VariantTnc         string    `json:"variant_tnc"`
		VariantDescription string    `json:"variant_description"`
		ValidPartners      []string  `json:"valid_partners"`
	}
	VariantDetailResponse struct {
		VariantName        string    `json:"variant_name"`
		VariantType        string    `json:"variant_type"`
		VoucherFormat      FormatReq `json:"voucher_format"`
		VoucherType        string    `json:"voucher_type"`
		VoucherPrice       float64   `json:"voucher_price"`
		AllowAccumulative  bool      `json:"allow_accumulative"`
		StartDate          string    `json:"start_date"`
		EndDate            string    `json:"end_date"`
		DiscountValue      float64   `json:"discount_value"`
		MaxQuantityVoucher float64   `json:"max_quantity_voucher"`
		MaxUsageVoucher    float64   `json:"max_usage_voucher"`
		RedeemtionMethod   string    `json:"redeemtion_method"`
		ImgUrl             string    `json:"image_url"`
		VariantTnc         string    `json:"variant_tnc"`
		VariantDescription string    `json:"variant_description"`
		ValidPartners      []string  `json:"valid_partners"`
		Voucher            int       `json:"-"`
	}
	FormatReq struct {
		Prefix     string `json:"prefix"`
		Postfix    string `json:"postfix"`
		Body       string `json:"body"`
		FormatType string `json:"format_type"`
		Length     int    `json:"length"`
	}
	UserVariantRequest struct {
		User string `json:"user"`
	}
	MultiUserVariantRequest struct {
		User string   `json:"user"`
		Data []string `json:"data"`
	}
)

func ListVariants(w http.ResponseWriter, r *http.Request) {
	res := NewResponse(nil)
	var status int

	accountID, _, _, ok := CheckToken(w, r)
	if !ok {
		return
	}
	param := getUrlParam(r.URL.String())

	param["variant_type"] = model.VariantTypeOnDemand
	param["account_id"] = accountID
	delete(param, "token")

	variant, err := model.FindVariantsCustomParam(param)
	if err == model.ErrResourceNotFound {
		status = http.StatusNotFound
		res.AddError(its(status), model.ErrCodeResourceNotFound, model.ErrMessageNilVariant, "voucher")
		render.JSON(w, res, status)
		return
	} else if err != nil {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", "voucher")
		render.JSON(w, res, status)
		return
	}
	d := make(GetVoucherOfVariatList, len(variant))
	for k, dt := range variant {
		d[k].VariantID = dt.Id
		d[k].AccountId = dt.AccountId
		d[k].VariantName = dt.VariantName
		d[k].VoucherType = dt.VoucherType
		d[k].VoucherPrice = dt.VoucherPrice
		d[k].DiscountValue = dt.DiscountValue
		d[k].StartDate = dt.StartDate
		d[k].EndDate = dt.EndDate
		d[k].ImgUrl = dt.ImgUrl
		d[k].Used = its(getCountVoucher(dt.Id))
	}

	status = http.StatusOK
	res = NewResponse(d)
	render.JSON(w, res, status)
}

func ListVariantsDetails(w http.ResponseWriter, r *http.Request) {
	variant := bone.GetValue(r, "id")
	res := NewResponse(nil)
	var status int

	_, _, _, ok := CheckToken(w, r)
	if !ok {
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

	status = http.StatusOK
	res = NewResponse(d)
	render.JSON(w, res, status)
}

func GetAllVariants(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get Variant")
	accountId := ""

	token := r.FormValue("token")
	status := http.StatusUnauthorized
	err := model.ErrTokenNotFound
	res := NewResponse(nil)

	res.AddError(its(status), its(status), err.Error(), "variant")

	valid := false
	if token != "" && token != "null" {
		fmt.Println("Check Session")
		_, accountId, _, valid, err = getValiditySession(r, token)
	}

	if valid {
		status = http.StatusOK
		variant, err := model.FindAllVariants(accountId)
		if err != nil {
			status = http.StatusInternalServerError
			if err != model.ErrResourceNotFound {
				status = http.StatusNotFound
			}

			res.AddError(its(status), its(status), err.Error(), "variant")
		} else {
			res = NewResponse(variant)
		}
	}
	render.JSON(w, res, status)
}

func GetTotalVariant(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get Variant")
	accountId := ""

	token := r.FormValue("token")
	status := http.StatusUnauthorized
	err := model.ErrTokenNotFound
	res := NewResponse(nil)

	res.AddError(its(status), its(status), err.Error(), "variant")

	valid := false
	if token != "" && token != "null" {
		fmt.Println("Check Session")
		_, accountId, _, valid, err = getValiditySession(r, token)
	}

	if valid {
		status = http.StatusOK
		variant, err := model.FindAllVariants(accountId)
		if err != nil {
			status = http.StatusInternalServerError
			if err != model.ErrResourceNotFound {
				status = http.StatusNotFound
			}

			res.AddError(its(status), its(status), err.Error(), "variant")
		} else {
			res = NewResponse(len(variant))
		}
	}
	render.JSON(w, res, status)
}

func GetVariantDetailsCustom(w http.ResponseWriter, r *http.Request) {
	param := getUrlParam(r.URL.String())
	token := r.FormValue("token")

	status := http.StatusUnauthorized
	err := model.ErrTokenNotFound
	res := NewResponse(nil)

	res.AddError(its(status), its(status), err.Error(), "variant")

	valid := false
	if token != "" && token != "null" {
		_, _, _, valid, _ = getValiditySession(r, token)
	}

	if valid {
		status = http.StatusOK
		variant, err := model.FindVariantDetailsCustomParam(param)
		if err != nil {
			status = http.StatusInternalServerError
			if err != model.ErrResourceNotFound {
				status = http.StatusNotFound
			}

			res.AddError(its(status), its(status), err.Error(), "variant")
		} else {
			res = NewResponse(variant)
		}
	}

	render.JSON(w, res, status)
}

func GetVariants(w http.ResponseWriter, r *http.Request) {
	param := getUrlParam(r.URL.String())
	token := r.FormValue("token")

	status := http.StatusUnauthorized
	err := model.ErrTokenNotFound
	res := NewResponse(nil)

	res.AddError(its(status), its(status), err.Error(), "variant")

	valid := false

	if token != "" && token != "null" {
		_, _, _, valid, _ = getValiditySession(r, token)
	}

	if valid {
		status = http.StatusOK
		variant, err := model.FindVariantsCustomParam(param)
		if err != nil {
			status = http.StatusInternalServerError
			if err != model.ErrResourceNotFound {
				status = http.StatusNotFound
			}

			res.AddError(its(status), its(status), err.Error(), "variant")
		} else {
			res = NewResponse(variant)
		}
	}

	render.JSON(w, res, status)
}

func GetVariantDetailsById(w http.ResponseWriter, r *http.Request) {
	id := bone.GetValue(r, "id")
	token := r.FormValue("token")

	status := http.StatusUnauthorized
	err := model.ErrTokenNotFound
	res := NewResponse(nil)

	res.AddError(its(status), its(status), err.Error(), "variant")

	valid := false
	if token != "" && token != "null" {
		_, _, _, valid, _ = getValiditySession(r, token)
	}

	if valid {
		status = http.StatusOK
		variant, err := model.FindVariantDetailsById(id)
		if err != nil {
			status = http.StatusInternalServerError
			if err != model.ErrResourceNotFound {
				status = http.StatusNotFound
			}

			res.AddError(its(status), its(status), err.Error(), "variant")
		} else {
			res = NewResponse(variant)

		}
	}

	render.JSON(w, res, status)
}

func GetVariantDetailsByDate(w http.ResponseWriter, r *http.Request) {
	start := r.FormValue("start")
	end := r.FormValue("end")
	token := r.FormValue("token")

	accountId := ""
	status := http.StatusUnauthorized
	err := model.ErrTokenNotFound
	res := NewResponse(nil)

	res.AddError(its(status), its(status), err.Error(), "variant")

	valid := false
	if token != "" && token != "null" {
		_, accountId, _, valid, _ = getValiditySession(r, token)
	}

	if valid {
		status = http.StatusOK
		variant, err := model.FindVariantsByDate(start, end, accountId)
		if err != nil {
			status = http.StatusInternalServerError
			if err != model.ErrResourceNotFound {
				status = http.StatusNotFound
			} else {
				res = NewResponse(variant)
			}

			res.AddError(its(status), its(status), err.Error(), "variant")
		} else {
			res = NewResponse(variant)
		}
	}

	render.JSON(w, res, status)
}

func CreateVariant(w http.ResponseWriter, r *http.Request) {
	token := r.FormValue("token")

	valid := false
	res := NewResponse(nil)
	accountId := ""
	user := ""
	status := http.StatusUnauthorized
	if token != "" && token != "null" {
		user, accountId, _, valid, _ = getValiditySession(r, token)
	}

	if valid {
		status = http.StatusCreated

		var rd Variant
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&rd); err != nil {
			log.Panic(err)
		}

		ts, err := time.Parse("01/02/2006", rd.StartDate)
		if err != nil {
			log.Panic(err)
		}
		te, err := time.Parse("01/02/2006", rd.EndDate)
		if err != nil {
			log.Panic(err)
		}

		vr := model.VariantReq{
			AccountId:          accountId,
			VariantName:        rd.VariantName,
			VariantType:        rd.VariantType,
			VoucherType:        rd.VoucherType,
			VoucherPrice:       rd.VoucherPrice,
			MaxQuantityVoucher: rd.MaxQuantityVoucher,
			MaxUsageVoucher:    rd.MaxUsageVoucher,
			AllowAccumulative:  rd.AllowAccumulative,
			RedeemtionMethod:   rd.RedeemtionMethod,
			DiscountValue:      rd.DiscountValue,
			StartDate:          ts.Format("2006-01-02 15:04:05.000"),
			EndDate:            te.Format("2006-01-02 15:04:05.000"),
			ImgUrl:             rd.ImgUrl,
			VariantTnc:         rd.VariantTnc,
			VariantDescription: rd.VariantDescription,
			ValidPartners:      rd.ValidPartners,
		}
		fr := model.FormatReq{
			Prefix:     rd.VoucherFormat.Prefix,
			Postfix:    rd.VoucherFormat.Postfix,
			Body:       rd.VoucherFormat.Body,
			FormatType: rd.VoucherFormat.FormatType,
			Length:     rd.VoucherFormat.Length,
		}

		if err := model.InsertVariant(vr, fr, user); err != nil {
			//log.Panic(err)
			status = http.StatusInternalServerError
		}

	}

	render.JSON(w, res, status)
}

func UpdateVariant(w http.ResponseWriter, r *http.Request) {
	id := bone.GetValue(r, "id")
	token := r.FormValue("token")
	user := ""
	valid := false
	status := http.StatusUnauthorized
	if token != "" && token != "null" {
		user, _, _, valid, _ = getValiditySession(r, token)
	}

	if valid {
		status = http.StatusOK
		var rd Variant
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&rd); err != nil {
			log.Panic(err)
		}

		ts, err := time.Parse("01/02/2006", rd.StartDate)
		if err != nil {
			log.Panic(err)
		}
		te, err := time.Parse("01/02/2006", rd.EndDate)
		if err != nil {
			log.Panic(err)
		}

		vr := model.Variant{
			Id:                 id,
			VariantName:        rd.VariantName,
			VariantType:        rd.VariantType,
			VoucherType:        rd.VoucherType,
			VoucherPrice:       rd.VoucherPrice,
			MaxQuantityVoucher: rd.MaxQuantityVoucher,
			MaxUsageVoucher:    rd.MaxUsageVoucher,
			RedeemtionMethod:   rd.RedeemtionMethod,
			DiscountValue:      rd.DiscountValue,
			StartDate:          ts.Format("2006-01-02 15:04:05.000"),
			EndDate:            te.Format("2006-01-02 15:04:05.000"),
			ImgUrl:             rd.ImgUrl,
			VariantTnc:         rd.VariantTnc,
			VariantDescription: rd.VariantDescription,
			CreatedBy:          user,
		}
		if err := model.UpdateVariant(vr); err != nil {
			//log.Panic(err)
			status = http.StatusInternalServerError
		}
	}

	res := NewResponse(nil)
	render.JSON(w, res, status)
}

func UpdateVariantBroadcast(w http.ResponseWriter, r *http.Request) {
	id := bone.GetValue(r, "id")
	token := r.FormValue("token")
	user := ""

	var rd MultiUserVariantRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		log.Panic(err)
	}

	valid := false
	res := NewResponse(nil)
	status := http.StatusUnauthorized
	if token != "" && token != "null" {
		user, _, _, valid, _ = getValiditySession(r, token)
	}

	if valid {
		status = http.StatusOK
		d := model.UpdateVariantUsersRequest{
			VariantId: id,
			User:      user,
			Data:      rd.Data,
		}

		if err := model.UpdateVariantBroadcasts(d); err != nil {
			//log.Panic(err)
			status = http.StatusInternalServerError
		}

	}

	res = NewResponse(nil)
	render.JSON(w, res, status)
}

func UpdateVariantTenant(w http.ResponseWriter, r *http.Request) {
	id := bone.GetValue(r, "id")
	token := r.FormValue("token")
	user := ""

	var rd MultiUserVariantRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		log.Panic(err)
	}

	valid := false
	res := NewResponse(nil)
	status := http.StatusUnauthorized
	if token != "" && token != "null" {
		user, _, _, valid, _ = getValiditySession(r, token)
	}

	if valid {
		status = http.StatusOK
		d := model.UpdateVariantUsersRequest{
			VariantId: id,
			User:      user,
			Data:      rd.Data,
		}

		if err := model.UpdateVariantPartners(d); err != nil {
			//log.Panic(err)
			status = http.StatusInternalServerError
		}
	}

	res = NewResponse(nil)
	render.JSON(w, res, status)
}

func DeleteVariant(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get Variant")
	token := r.FormValue("token")
	user := ""
	id := bone.GetValue(r, "id")

	valid := false
	status := http.StatusUnauthorized
	if token != "" && token != "null" {
		user, _, _, valid, _ = getValiditySession(r, token)
	}

	if valid {
		status = http.StatusOK
		d := &model.DeleteVariantRequest{
			Id:   id,
			User: user,
		}
		if err := d.Delete(); err != nil {
			status = http.StatusInternalServerError
		}
	}

	res := NewResponse(nil)
	render.JSON(w, res, status)
}

func CheckVariant(rm, id string, qty int) (bool, error) {
	dt, err := model.FindVariantDetailsById(id)
	sd, err := time.Parse(time.RFC3339Nano, dt.StartDate)
	if err != nil {
		return false, err
	}
	ed, err := time.Parse(time.RFC3339Nano, dt.EndDate)
	if err != nil {
		return false, err
	}

	if !sd.Before(time.Now()) || !ed.After(time.Now()) {
		return false, errors.New(model.ErrCodeVoucherNotActive)
	}

	if dt.AllowAccumulative == false && qty > 1 {
		return false, errors.New(model.ErrCodeAllowAccumulativeDisable)
	}

	if dt.RedeemtionMethod != rm {
		return false, errors.New(model.ErrCodeInvalidRedeemMethod)
	}

	return true, nil
}
