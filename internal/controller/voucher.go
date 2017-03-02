package controller

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gilkor/evoucher/internal/model"
	"github.com/ruizu/render"
)

type (
	//###################REQUEST FORMAT###################//
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
		ReferenceNo string `json:"reference_no"`
	}
	//#####################################################//

	//###################RESPONSE FORMAT###################//
	// VoucherResponse represent a Response of GenerateVouche
	VoucherResponse struct {
		State       string      `json:"state"`
		Description string      `json:"messange"`
		VoucherData interface{} `json:"data"`
	}
	// GenerateVoucerResponse represent list of voucher data
	GenerateVoucerResponse struct {
		VoucherID string `json:"voucher_id"`
		VoucherNo string `json:"voucher_code"`
	}
	// DetailVoucherResponse represent list of voucher data
	DetailListVoucherResponse []DetailVoucherResponse
	DetailVoucherResponse     struct {
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
	//#####################################################//r
)

// GetVoucherDetail get Voucher detail from DB
func GetVoucherDetail(w http.ResponseWriter, r *http.Request) {
	var voucher model.VoucherResponse
	var err error
	var vr VoucherResponse
	status := http.StatusOK

	id := r.FormValue("id")
	code := r.FormValue("code")
	variantid := r.FormValue("variant_id")

	if id != "" {
		voucher, err = model.FindVoucherByID(id)
	} else if code != "" {
		voucher, err = model.FindVoucherByCode(code)
	} else if variantid != "" {
		voucher, err = model.FindVoucherByVariant(variantid)
	} else {
		vr.State = model.ErrCodeMissingOrderItem
		vr.Description = model.ErrMessageMissingOrderItem
		status = http.StatusBadRequest
	}

	if err == model.ErrResourceNotFound {
		vr.State = model.ErrCodeResourceNotFound
		vr.Description = model.ErrResourceNotFound.Error()
	} else if err != nil {
		vr.State = model.ResponseStateNok
		vr.Description = err.Error()
		status = http.StatusInternalServerError
	} else if voucher.Message != "" {

		dvr := make(DetailListVoucherResponse, len(voucher.VoucherData))
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
		vr = VoucherResponse{State: model.ResponseStateOk, Description: "", VoucherData: dvr}
	}

	res := NewResponse(vr)
	render.JSON(w, res, status)
}

//RedeemVoucher redeem
func RedeemVoucher(w http.ResponseWriter, r *http.Request) {
	var d model.UpdateDeleteRequest
	var vrr RedeemVoucherRequest
	var vr VoucherResponse
	status := http.StatusOK

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&vrr); err != nil {
		status = http.StatusInternalServerError
		log.Panic(err)
	}

	if basicAuth(w, r) {
		fv, err := model.FindVoucherByCode(vrr.VoucherCode)
		if err == model.ErrResourceNotFound {
			vr.State = model.ErrCodeResourceNotFound
			vr.Description = model.ErrResourceNotFound.Error()
		} else if err != nil {
			vr.State = model.ResponseStateNok
			vr.Description = err.Error()
			status = http.StatusInternalServerError
		} else if fv.Status != model.ResponseStateOk {
			vr.State = fv.Status
			vr.Description = fv.Message
		} else {

			d.ID = fv.VoucherData[0].ID
			d.State = model.VoucherStateUsed
			d.User = vrr.AccountID

			uv, err := d.UpdateVc()
			if err != nil {
				vr.State = model.ResponseStateNok
				vr.Description = err.Error()
				status = http.StatusInternalServerError
			} else {
				vr.State = model.ResponseStateOk
				vr.Description = ""

				dvr := DetailVoucherResponse{
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
				vr.VoucherData = dvr
			}
		}
	} else {
		vr = VoucherResponse{}
		status = http.StatusUnauthorized
	}

	res := NewResponse(vr)
	render.JSON(w, res, status)
}

// DeleteVoucher delete Voucher data from DB by ID
func DeleteVoucher(w http.ResponseWriter, r *http.Request) {
	var d model.UpdateDeleteRequest
	var rd DeleteVoucherRequest
	status := http.StatusOK
	vr := VoucherResponse{State: model.ResponseStateOk, Description: ""}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		log.Panic(err)
		status = http.StatusInternalServerError
	}

	if basicAuth(w, r) {

		fv, err := model.FindVoucherByCode(rd.VoucherCode)
		if err == model.ErrResourceNotFound {
			vr.State = model.ErrCodeResourceNotFound
			vr.Description = model.ErrResourceNotFound.Error()
		} else if err != nil {
			vr.State = model.ResponseStateNok
			vr.Description = err.Error()
			status = http.StatusInternalServerError
		} else {
			d.ID = fv.VoucherData[0].ID
			d.State = model.VoucherStateDeleted
			d.User = fv.VoucherData[0].CreatedBy

			err := d.DeleteVc()
			if err != nil {
				vr.State = model.ResponseStateNok
				vr.Description = err.Error()
			} else {
				vr.State = model.ResponseStateOk
				vr.Description = ""
				vr.VoucherData = fv.VoucherData[0].ID
			}
		}

	} else {
		vr = VoucherResponse{}
		status = http.StatusUnauthorized
	}

	res := NewResponse(vr)
	render.JSON(w, res, status)
}

func PayVoucher(w http.ResponseWriter, r *http.Request) {
	var d model.UpdateDeleteRequest
	var pvr PayVoucherRequest
	status := http.StatusOK
	vr := VoucherResponse{State: model.ResponseStateOk, Description: ""}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&pvr); err != nil {
		vr.State = model.ResponseStateNok
		vr.Description = err.Error()
		status = http.StatusInternalServerError
		// log.Panic(err)
	}

	if basicAuth(w, r) {

		fv, err := model.FindVoucherByCode(pvr.VoucherCode)
		if err == model.ErrResourceNotFound {
			vr.State = model.ErrCodeResourceNotFound
			vr.Description = model.ErrResourceNotFound.Error()
		} else if err != nil {
			vr.State = model.ResponseStateNok
			vr.Description = err.Error()
			status = http.StatusInternalServerError
		} else if fv.VoucherData[0].State == model.VoucherStatePaid {
			vr.State = model.ErrCodeVoucherAlreadyPaid
			vr.Description = model.ErrMessageVoucherAlreadyPaid
		} else {
			d.ID = fv.VoucherData[0].ID
			d.State = model.VoucherStatePaid
			d.User = fv.VoucherData[0].CreatedBy

			uv, err := d.UpdateVc()
			if err != nil {
				vr.State = model.ResponseStateNok
				vr.Description = err.Error()
			} else {
				vr.State = model.ResponseStateOk
				vr.Description = ""
				dvr := DetailVoucherResponse{
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
				vr.VoucherData = dvr
			}
		}

	} else {
		vr = VoucherResponse{}
		status = http.StatusUnauthorized
	}

	res := NewResponse(vr)
	render.JSON(w, res, status)
}

//GenerateVoucherOnDemand Generate singgle voucher request
func GenerateVoucherOnDemand(w http.ResponseWriter, r *http.Request) {
	var gvd GenerateVoucherRequest
	status := http.StatusOK
	vr := VoucherResponse{State: model.ResponseStateOk, Description: ""}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&gvd); err != nil {
		status = http.StatusInternalServerError
		log.Panic(err)
	}

	if basicAuth(w, r) {

		variant, err := model.FindVariantById(gvd.VariantID)
		if err == model.ErrResourceNotFound {
			vr.State = model.ErrCodeInvalidVoucher
			vr.Description = model.ErrMessageInvalidVoucher
		} else if err != nil {
			vr.State = model.ResponseStateNok
			vr.Description = err.Error()
			status = http.StatusInternalServerError
		} else {

			if dt, ok := variant.Data.(model.Variant); ok {
				if (int(dt.MaxQuantityVoucher) - getCountVoucher(gvd.VariantID) - 1) <= 0 {
					vr.State = model.ErrCodeVoucherQtyExceeded
					vr.Description = model.ErrMessageVoucherQtyExceeded
				} else if dt.VariantType != model.VariantTypeOnDemand {
					vr.State = model.ErrCodeVoucherRulesViolated
					vr.Description = model.ErrMessageVoucherRulesViolated
				} else {
					d := GenerateVoucherRequest{
						VariantID:   dt.Id,
						Quantity:    1,
						ReferenceNo: gvd.ReferenceNo,
						AccountID:   gvd.AccountID,
					}

					var voucher []model.Voucher
					voucher, err = d.generateVoucherBulk(&dt)
					if err != nil {
						vr.State = model.ResponseStateNok
						vr.Description = err.Error()
						status = http.StatusInternalServerError
					} else {
						gvr := make([]GenerateVoucerResponse, len(voucher))
						for i, v := range voucher {
							gvr[i].VoucherID = v.ID
							gvr[i].VoucherNo = v.VoucherCode
						}

						vr.State = model.ResponseStateOk
						vr.Description = ""
						vr.VoucherData = gvr
					}
				}
			} else {
				vr.State = model.ErrCodeInternalError
				vr.Description = model.ErrMessageInternalError
			}
		}

	} else {
		vr = VoucherResponse{}
		status = http.StatusUnauthorized
	}

	res := NewResponse(vr)
	render.JSON(w, res, status)
}

//GenerateVoucher Generate bulk voucher request
func GenerateVoucher(w http.ResponseWriter, r *http.Request) {
	var gvd GenerateVoucherRequest
	status := http.StatusOK
	vr := VoucherResponse{State: model.ResponseStateOk, Description: ""}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&gvd); err != nil {
		status = http.StatusInternalServerError
		log.Panic(err)
	}

	if basicAuth(w, r) {

		variant, err := model.FindVariantById(gvd.VariantID)
		if err == model.ErrResourceNotFound {
			vr.State = model.ErrCodeResourceNotFound
			vr.Description = model.ErrMessageInvalidVoucher
		} else if err != nil {
			vr.State = model.ResponseStateNok
			vr.Description = err.Error()
		} else {

			if dt, ok := variant.Data.(model.Variant); ok {
				if (int(dt.MaxQuantityVoucher) - getCountVoucher(gvd.VariantID) - gvd.Quantity) <= 0 {
					vr.State = model.ErrCodeVoucherQtyExceeded
					vr.Description = model.ErrMessageVoucherQtyExceeded
				} else if dt.VariantType != model.VariantTypeBulk {
					vr.State = model.ErrCodeVoucherRulesViolated
					vr.Description = model.ErrMessageVoucherRulesViolated
				} else {
					d := GenerateVoucherRequest{
						VariantID:   dt.Id,
						Quantity:    gvd.Quantity,
						ReferenceNo: gvd.ReferenceNo,
						AccountID:   gvd.AccountID,
					}

					var voucher []model.Voucher
					voucher, err = d.generateVoucherBulk(&dt)
					if err != nil {
						vr.State = model.ResponseStateNok
						vr.Description = err.Error()
					} else {

						gvr := make([]GenerateVoucerResponse, len(voucher))
						for i, v := range voucher {
							gvr[i].VoucherID = v.ID
							gvr[i].VoucherNo = v.VoucherCode
						}

						vr.State = model.ResponseStateOk
						vr.Description = ""
						vr.VoucherData = gvr
					}

				}
			} else {
				vr.State = model.ErrCodeInternalError
				vr.Description = model.ErrMessageInternalError
			}
		}

	} else {
		vr = VoucherResponse{}
		status = http.StatusUnauthorized
	}

	res := NewResponse(vr)
	render.JSON(w, res, status)
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
			// fmt.Println("1 :", code)
		case vcf.Body.Valid == true:
			code = vcf.Prefix.String + vcf.Body.String + vcf.Postfix.String
			// fmt.Println("2 :", code)
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
			CreatedBy:     vr.AccountID,
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

func randStr(ln int, fm string) string {
	CharsType := map[string]string{
		"Alphabet":     model.ALPHABET,
		"Numerals":     model.NUMERALS,
		"Alphanumeric": model.ALPHANUMERIC,
	}

	rand.Seed(time.Now().UTC().UnixNano())
	chars := CharsType[fm]
	result := make([]byte, ln)
	for i := 0; i < ln; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

func getCountVoucher(variantID string) int {
	return model.CountVoucher(variantID)
}
