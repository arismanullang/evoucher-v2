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
		StartHour          string    `json:"start_hour"`
		EndHour            string    `json:"end_hour"`
		ValidVoucherStart  string    `json:"valid_voucher_start"`
		ValidVoucherEnd    string    `json:"valid_voucher_end"`
		VoucherLifetime    int       `json:"voucher_lifetime"`
		ValidityDays       string    `json:"validity_days"`
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
		Id                 string  `json:"id"`
		AccountId          string  `json:"account_id"`
		VariantName        string  `json:"variant_name"`
		VariantType        string  `json:"variant_type"`
		VoucherFormat      int     `json:"voucher_format"`
		VoucherType        string  `json:"voucher_type"`
		VoucherPrice       float64 `json:"voucher_price"`
		AllowAccumulative  bool    `json:"allow_accumulative"`
		StartDate          string  `json:"start_date"`
		EndDate            string  `json:"end_date"`
		DiscountValue      float64 `json:"discount_value"`
		MaxQuantityVoucher float64 `json:"max_quantity_voucher"`
		MaxUsageVoucher    float64 `json:"max_usage_voucher"`
		RedeemtionMethod   string  `json:"redeemtion_method"`
		ImgUrl             string  `json:"image_url"`
		VariantTnc         string  `json:"variant_tnc"`
		VariantDescription string  `json:"variant_description"`
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

func ListVariants(w http.ResponseWriter, r *http.Request) {
	res := NewResponse(nil)
	var status int

	a := AuthToken(w, r)
	if !a.Valid {
		render.JSON(w, a.res, http.StatusUnauthorized)
		return
	}
	param := getUrlParam(r.URL.String())

	param["variant_type"] = model.VariantTypeOnDemand
	param["account_id"] = a.User.Account.Id
	delete(param, "token")

	variant, err := model.FindAvailableVariants()
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
	d := []GetVoucherOfVariatdata{}
	for _, dt := range variant {
		if (int(dt.MaxVoucher) - sti(dt.Voucher)) > 0 {
			tempVoucher := GetVoucherOfVariatdata{}
			tempVoucher.VariantID = dt.Id
			tempVoucher.AccountId = dt.AccountId
			tempVoucher.VariantName = dt.VariantName
			tempVoucher.VoucherType = dt.VoucherType
			tempVoucher.VoucherPrice = dt.VoucherPrice
			tempVoucher.DiscountValue = dt.DiscountValue
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
	render.JSON(w, res, status)
}

func ListVariantsDetails(w http.ResponseWriter, r *http.Request) {
	variant := bone.GetValue(r, "id")
	res := NewResponse(nil)
	var status int

	a := AuthToken(w, r)
	if !a.Valid {
		render.JSON(w, a.res, http.StatusUnauthorized)
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

	d.Used = getCountVoucher(dt.Id)

	status = http.StatusOK
	res = NewResponse(d)
	render.JSON(w, res, status)
}

func GetAllVariants(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get Variant")
	status := http.StatusUnauthorized
	err := model.ErrTokenNotFound
	errTitle := model.ErrCodeInvalidToken
	res := NewResponse(nil)
	res.AddError(its(status), errTitle, err.Error(), "Get Variant")

	fmt.Println("Check Session")
	a := AuthToken(w, r)
	if a.Valid {
		status = http.StatusOK
		variant, err := model.FindAllVariants(a.User.Account.Id)
		fmt.Println(err)
		if err != nil {
			status = http.StatusInternalServerError
			errTitle = model.ErrCodeInternalError
			if err == model.ErrResourceNotFound {
				status = http.StatusNotFound
				errTitle = model.ErrCodeResourceNotFound
			}

			res.AddError(its(status), errTitle, err.Error(), "Get Variant")
		} else {
			res = NewResponse(variant)
		}
	}
	render.JSON(w, res, status)
}

func GetTotalVariant(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get Variant")
	status := http.StatusUnauthorized
	err := model.ErrTokenNotFound
	errTitle := model.ErrCodeInvalidToken
	res := NewResponse(nil)
	res.AddError(its(status), errTitle, err.Error(), "Get Variant")

	fmt.Println("Check Session")
	a := AuthToken(w, r)

	if a.Valid {
		status = http.StatusOK
		variant, err := model.FindAllVariants(a.User.Account.Id)
		if err != nil {
			status = http.StatusInternalServerError
			errTitle = model.ErrCodeInternalError
			if err == model.ErrResourceNotFound {
				status = http.StatusNotFound
				errTitle = model.ErrCodeResourceNotFound
			}

			res.AddError(its(status), errTitle, err.Error(), "Get Variant")
		} else {
			res = NewResponse(len(variant))
		}
	}
	render.JSON(w, res, status)
}

func GetVariantDetailsCustom(w http.ResponseWriter, r *http.Request) {
	param := getUrlParam(r.URL.String())

	status := http.StatusUnauthorized
	err := model.ErrTokenNotFound
	errTitle := model.ErrCodeInvalidToken
	res := NewResponse(nil)
	res.AddError(its(status), errTitle, err.Error(), "Get Variant")

	a := AuthToken(w, r)
	if a.Valid {
		status = http.StatusOK
		variant, err := model.FindVariantDetailsCustomParam(param)
		if err != nil {
			status = http.StatusInternalServerError
			errTitle = model.ErrCodeInternalError
			if err == model.ErrResourceNotFound {
				status = http.StatusNotFound
				errTitle = model.ErrCodeResourceNotFound
			}

			res.AddError(its(status), errTitle, err.Error(), "Get Variant")
		} else {
			res = NewResponse(variant)
		}
	} else {
		res = a.res
		status = http.StatusUnauthorized
	}

	render.JSON(w, res, status)
}

func GetVariants(w http.ResponseWriter, r *http.Request) {
	param := getUrlParam(r.URL.String())

	status := http.StatusUnauthorized
	err := model.ErrTokenNotFound
	errTitle := model.ErrCodeInvalidToken
	res := NewResponse(nil)
	res.AddError(its(status), errTitle, err.Error(), "Get Variant")

	a := AuthToken(w, r)
	if a.Valid {
		status = http.StatusOK
		variant, err := model.FindVariantsCustomParam(param)
		if err != nil {
			status = http.StatusInternalServerError
			errTitle = model.ErrCodeInternalError
			if err == model.ErrResourceNotFound {
				status = http.StatusNotFound
				errTitle = model.ErrCodeResourceNotFound
			}

			res.AddError(its(status), errTitle, err.Error(), "Get Variant")
		} else {
			res = NewResponse(variant)
		}
	} else {
		res = a.res
		status = http.StatusUnauthorized
	}

	render.JSON(w, res, status)
}

func GetVariantDetailsById(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	status := http.StatusUnauthorized
	err := model.ErrTokenNotFound
	errTitle := model.ErrCodeInvalidToken
	res := NewResponse(nil)
	res.AddError(its(status), errTitle, err.Error(), "Get Variant")

	a := AuthToken(w, r)
	if a.Valid {
		status = http.StatusOK
		variant, err := model.FindVariantDetailsById(id)
		if err != nil {
			status = http.StatusInternalServerError
			errTitle = model.ErrCodeInternalError
			if err == model.ErrResourceNotFound {
				status = http.StatusNotFound
				errTitle = model.ErrCodeResourceNotFound
			}

			res.AddError(its(status), errTitle, err.Error(), "Get Variant")
		} else {
			res = NewResponse(variant)

		}
	} else {
		res = a.res
		status = http.StatusUnauthorized
	}

	render.JSON(w, res, status)
}

func GetVariantDetailsByDate(w http.ResponseWriter, r *http.Request) {
	start := r.FormValue("start")
	end := r.FormValue("end")
	status := http.StatusUnauthorized
	err := model.ErrTokenNotFound
	errTitle := model.ErrCodeInvalidToken
	res := NewResponse(nil)

	res.AddError(its(status), errTitle, err.Error(), "Get Variant")

	a := AuthToken(w, r)
	if a.Valid {
		status = http.StatusOK
		variant, err := model.FindVariantsByDate(start, end, a.User.Account.Id)
		if err != nil {
			status = http.StatusInternalServerError
			errTitle = model.ErrCodeInternalError
			if err == model.ErrResourceNotFound {
				status = http.StatusNotFound
				errTitle = model.ErrCodeResourceNotFound
			} else {
				res = NewResponse(variant)
			}

			res.AddError(its(status), errTitle, err.Error(), "Get Variant")
		} else {
			res = NewResponse(variant)
		}
	} else {
		res = a.res
		status = http.StatusUnauthorized
	}

	render.JSON(w, res, status)
}

func CreateVariant(w http.ResponseWriter, r *http.Request) {
	apiName := "variant_create"
	valid := false

	res := NewResponse(nil)
	status := http.StatusUnauthorized
	err := model.ErrInvalidRole
	errTitle := model.ErrCodeInvalidRole
	res.AddError(its(status), errTitle, err.Error(), "Create Variant")

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
			fmt.Println(rd.ValidVoucherStart)
			tvs, err := time.Parse("01/02/2006", rd.ValidVoucherStart)
			if err != nil {
				log.Panic(err)
			}
			tve, err := time.Parse("01/02/2006", rd.ValidVoucherEnd)
			if err != nil {
				log.Panic(err)
			}

			vr := model.VariantReq{
				AccountId:          a.User.Account.Id,
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
				StartHour:          rd.StartHour,
				EndHour:            rd.EndHour,
				ValidVoucherStart:  tvs.Format("2006-01-02 15:04:05.000"),
				ValidVoucherEnd:    tve.Format("2006-01-02 15:04:05.000"),
				VoucherLifetime:    rd.VoucherLifetime,
				ValidityDays:       rd.ValidityDays,
				ImgUrl:             rd.ImgUrl,
				VariantTnc:         rd.VariantTnc,
				VariantDescription: rd.VariantDescription,
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
			if id, err := model.InsertVariant(vr, fr, a.User.ID); err != nil {
				//log.Panic(err)
				status = http.StatusInternalServerError
				errTitle = model.ErrCodeInternalError
				res.AddError(its(status), errTitle, err.Error(), "Create Variant")
			} else {
				res = NewResponse(id)
			}
		}
	} else {
		res = a.res
		status = http.StatusUnauthorized
	}

	render.JSON(w, res, status)
}

func UpdateVariantRoute(w http.ResponseWriter, r *http.Request) {
	types := r.FormValue("type")
	apiName := "variant_update"
	valid := false

	res := NewResponse(nil)
	status := http.StatusUnauthorized
	err := model.ErrInvalidRole
	errTitle := model.ErrCodeInvalidRole
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
	}

	if types == "" {
		res.AddError(its(status), errTitle, err.Error(), "Update type not found")
		render.JSON(w, res, status)
	} else {
		if valid {
			if types == "detail" {
				UpdateVariant(w, r)
			} else if types == "tenant" {
				UpdateVariantTenant(w, r)
			} else if types == "broadcast" {
				UpdateVariantBroadcast(w, r)
			} else {
				res.AddError(its(status), errTitle, err.Error(), "Update type not allowed")
				render.JSON(w, res, status)
			}
		} else {
			res.AddError(its(status), errTitle, err.Error(), "Update Variant")
			render.JSON(w, res, status)
		}
	}
}

func UpdateVariant(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")

	res := NewResponse(nil)
	status := http.StatusUnauthorized
	err := model.ErrTokenNotFound
	errTitle := model.ErrCodeInvalidToken
	res.AddError(its(status), errTitle, err.Error(), "Update Variant")
	a := AuthToken(w, r)
	if a.Valid {
		status = http.StatusOK
		var rd Variant
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&rd); err != nil {
			log.Panic(err)
		}

		fmt.Println(rd)
		ts, err := time.Parse("01/02/2006", rd.StartDate)
		if err != nil {
			log.Panic(err)
		}
		te, err := time.Parse("01/02/2006", rd.EndDate)
		if err != nil {
			log.Panic(err)
		}

		tvs, err := time.Parse("2006-01-02T00:00:00Z", rd.ValidVoucherStart)
		if err != nil {
			log.Panic(err)
		}
		tve, err := time.Parse("2006-01-02T00:00:00Z", rd.ValidVoucherEnd)
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
			StartHour:          rd.StartHour,
			EndHour:            rd.EndHour,
			AllowAccumulative:  rd.AllowAccumulative,
			ValidVoucherStart:  tvs.Format("2006-01-02 15:04:05.000"),
			ValidVoucherEnd:    tve.Format("2006-01-02 15:04:05.000"),
			VoucherLifetime:    rd.VoucherLifetime,
			ValidityDays:       rd.ValidityDays,
			ImgUrl:             rd.ImgUrl,
			VariantTnc:         rd.VariantTnc,
			VariantDescription: rd.VariantDescription,
			CreatedBy:          a.User.ID,
		}
		if err := model.UpdateVariant(vr); err != nil {
			status = http.StatusInternalServerError
			errTitle = model.ErrCodeInternalError
			res.AddError(its(status), errTitle, err.Error(), "Update Variant")
		}
		res = NewResponse("")
	} else {
		res = a.res
		status = http.StatusUnauthorized
	}
	render.JSON(w, res, status)
}

func UpdateVariantBroadcast(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")

	var rd MultiUserVariantRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		log.Panic(err)
	}

	res := NewResponse(nil)
	status := http.StatusUnauthorized
	err := model.ErrTokenNotFound
	errTitle := model.ErrCodeInvalidToken
	res.AddError(its(status), errTitle, err.Error(), "Update Variant")

	a := AuthToken(w, r)
	if a.Valid {
		status = http.StatusOK
		d := model.UpdateVariantArrayRequest{
			VariantId: id,
			User:      a.User.ID,
			Data:      rd.Data,
		}

		if err := model.UpdateVariantBroadcasts(d); err != nil {
			//log.Panic(err)
			status = http.StatusInternalServerError
			errTitle = model.ErrCodeInternalError
			res.AddError(its(status), errTitle, err.Error(), "Update Variant")
		}
		res = NewResponse("")
	} else {
		res = a.res
		status = http.StatusUnauthorized
	}

	render.JSON(w, res, status)
}

func UpdateVariantTenant(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")

	var rd MultiUserVariantRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		log.Panic(err)
	}

	res := NewResponse(nil)
	status := http.StatusUnauthorized
	err := model.ErrTokenNotFound
	errTitle := model.ErrCodeInvalidToken
	res.AddError(its(status), errTitle, err.Error(), "Update Variant")
	a := AuthToken(w, r)

	if a.Valid {
		status = http.StatusOK
		d := model.UpdateVariantArrayRequest{
			VariantId: id,
			User:      a.User.ID,
			Data:      rd.Data,
		}

		if err := model.UpdateVariantPartners(d); err != nil {
			status = http.StatusInternalServerError
			errTitle = model.ErrCodeInternalError
			res.AddError(its(status), errTitle, err.Error(), "Update Variant")
		}
		res = NewResponse("")
	} else {
		res = a.res
		status = http.StatusUnauthorized
	}

	render.JSON(w, res, status)
}

func DeleteVariant(w http.ResponseWriter, r *http.Request) {
	apiName := "variant_delete"
	valid := false

	res := NewResponse(nil)
	id := r.FormValue("id")

	status := http.StatusUnauthorized
	err := model.ErrInvalidRole
	errTitle := model.ErrCodeInvalidRole
	res.AddError(its(status), errTitle, err.Error(), "Delete Variant")
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
			d := &model.DeleteVariantRequest{
				Id:   id,
				User: a.User.ID,
			}
			if err := d.Delete(); err != nil {
				status = http.StatusInternalServerError
				errTitle = model.ErrCodeInternalError
				res.AddError(its(status), errTitle, err.Error(), "Delete Variant")

			}
			objName := strings.Split(d.Img_url, "/")
			if deleteFile(w, r, objName[4]) {
				return
			}
			res = NewResponse("")
		}
	} else {
		res = a.res
		status = http.StatusUnauthorized
	}
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

	if dt.RedeemtionMethod != rm {
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
