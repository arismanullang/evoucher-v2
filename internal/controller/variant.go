package controller

import (
	"encoding/json"
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

func GetAllVariants(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get Variant")
	accountId := r.FormValue("account_id")
	var variant model.Response
	var err error
	status := http.StatusOK

	variant, err = model.FindAllVariants(accountId)

	if err != nil && err != model.ErrResourceNotFound {
		//log.Panic(err)
		variant.Message = err.Error()
	}

	variant.Status = its(status)
	res := NewResponse(variant)
	render.JSON(w, res, status)
}

func GetVariants(w http.ResponseWriter, r *http.Request) {
	param := getUrlParam(r.URL.String())
	var variant model.Response
	var err error
	var status int
	if _, ok := basicAuth(w, r); ok {
		variant, err = model.FindVariantMultipleParam(param)
		if err != nil && err != model.ErrResourceNotFound {
			log.Panic(err)
		}
		status = http.StatusOK
		variant.Message = http.StatusText(status)
	} else {
		status = http.StatusUnauthorized
		variant.Message = http.StatusText(status)
	}

	variant.Status = its(status)
	res := NewResponse(variant)
	render.JSON(w, res, status)
}

func GetVariantDetailsById(w http.ResponseWriter, r *http.Request) {
	id := bone.GetValue(r, "id")
	var variant model.Response
	var err error
	var status int
	if _, ok := basicAuth(w, r); ok {
		variant, err = model.FindVariantById(id)
		if err != nil && err != model.ErrResourceNotFound {
			log.Panic(err)
		}
		status = http.StatusOK
		variant.Message = http.StatusText(status)
	} else {
		status = http.StatusUnauthorized
		variant.Message = http.StatusText(status)
	}

	variant.Status = its(status)
	res := NewResponse(variant)
	render.JSON(w, res, status)
}

func GetVariantDetailsByDate(w http.ResponseWriter, r *http.Request) {
	start := r.FormValue("start")
	end := r.FormValue("end")
	var variant model.Response
	var err error
	var status int
	if _, ok := basicAuth(w, r); ok {
		variant, err = model.FindVariantByDate(start, end)
		if err != nil && err != model.ErrResourceNotFound {
			log.Panic(err)
		}
		status = http.StatusOK
		variant.Message = http.StatusText(status)
	} else {
		status = http.StatusUnauthorized
		variant.Message = http.StatusText(status)
	}

	variant.Status = its(status)
	res := NewResponse(variant)
	render.JSON(w, res, status)
}

// dashboard
func CreateVariant(w http.ResponseWriter, r *http.Request) {
	token := r.FormValue("token")
	user := r.FormValue("user")
	valid := false
	res := NewResponse(nil)
	var account string
	if token != "" && token != "null" {
<<<<<<< HEAD
		user, account, exp, _ = checkExpired(r, token)
		if exp.After(time.Now()) {
			valid = true
		}
=======
		account, _, valid = getValiditySession(r, user, token)
>>>>>>> ed60a86f315f4521641200631c93d01b0c9c855e
	}

	if valid {
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
			AccountId:          account,
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
			log.Panic(err)
		}

		render.JSON(w, res, http.StatusCreated)
	} else {
		render.JSON(w, res, http.StatusUnauthorized)
	}
}

func UpdateVariant(w http.ResponseWriter, r *http.Request) {
	id := bone.GetValue(r, "id")
	token := r.FormValue("token")
	user := r.FormValue("user")
	valid := false
	res := NewResponse(nil)
	if token != "" && token != "null" {
<<<<<<< HEAD
		user, _, exp, _ = checkExpired(r, token)
		if exp.After(time.Now()) {
			valid = true
		}
=======
		_, _, valid = getValiditySession(r, user, token)
>>>>>>> ed60a86f315f4521641200631c93d01b0c9c855e
	}

	if valid {
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
			log.Panic(err)
		}
		render.JSON(w, res, http.StatusOK)
	} else {
		render.JSON(w, res, http.StatusUnauthorized)
	}
}

func UpdateVariantBroadcast(w http.ResponseWriter, r *http.Request) {
	id := bone.GetValue(r, "id")
	var rd MultiUserVariantRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		log.Panic(err)
	}

	d := model.UpdateVariantUsersRequest{
		VariantId: id,
		User:      rd.User,
		Data:      rd.Data,
	}

	if err := model.UpdateBroadcast(d); err != nil {
		log.Panic(err)
	}

	res := NewResponse(nil)
	render.JSON(w, res)
}

func UpdateVariantTenant(w http.ResponseWriter, r *http.Request) {
	id := bone.GetValue(r, "id")
	var rd MultiUserVariantRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		log.Panic(err)
	}

	d := model.UpdateVariantUsersRequest{
		VariantId: id,
		User:      rd.User,
		Data:      rd.Data,
	}

	if err := model.UpdatePartner(d); err != nil {
		log.Panic(err)
	}

	res := NewResponse(nil)
	render.JSON(w, res)
}

func DeleteVariant(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get Variant")
	token := r.FormValue("token")
	user := r.FormValue("user")
	id := bone.GetValue(r, "id")
	valid := false
	var status int
	if token != "" && token != "null" {
		_, _, valid = getValiditySession(r, user, token)
	}

	if valid {
		d := &model.DeleteVariantRequest{
			Id:   id,
			User: user,
		}
		if err := d.Delete(); err != nil {
			status = http.StatusInternalServerError
		} else {
			status = http.StatusOK
		}
	} else {
		status = http.StatusUnauthorized
	}
	res := NewResponse(nil)
	render.JSON(w, res, status)
}

func DashboardGetAllVariants(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get Variant")
	token := r.FormValue("token")
	user := r.FormValue("user")
	valid := false
	var variant model.Response
	var accountId string
	var err error
	var status int
	if token != "" && token != "null" {
		accountId, _, valid = getValiditySession(r, user, token)
	}

	if valid {
		variant, err = model.FindAllVariants(accountId)
		if err != nil && err != model.ErrResourceNotFound {
			log.Panic(err)
		}
		status = http.StatusOK
	} else {
		status = http.StatusUnauthorized
	}

	res := NewResponse(variant)
	render.JSON(w, res, status)
}

func DashboardGetVariantDetailsById(w http.ResponseWriter, r *http.Request) {
	token := r.FormValue("token")
	user := r.FormValue("user")
	id := bone.GetValue(r, "id")
	valid := false
	var variant model.Response
	var err error
	var status int
	if token != "" && token != "null" {
		_, _, valid = getValiditySession(r, user, token)
	}

	if valid {
		variant, err = model.FindVariantById(id)
		if err != nil && err != model.ErrResourceNotFound {
			log.Panic(err)
		}
		status = http.StatusOK
	} else {
		status = http.StatusUnauthorized
	}

	res := NewResponse(variant)
	render.JSON(w, res, status)
}
