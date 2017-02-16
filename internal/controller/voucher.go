package controller

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gilkor/evoucher/internal/model"
	"github.com/go-zoo/bone"
	"github.com/ruizu/render"
)

const (
	VoucherCodeLen int32 = 4 // VoucherCodeLen length of VoucherCode to run randStr
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
	//#####################################################//r

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
	DetailVoucherResponse struct {
		ID            string    `json:"id"`
		VoucherCode   string    `json:"voucher_code"`
		ReferenceNo   string    `json:"reference_no"`
		AccountID     string    `json:"account_id"`
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

	id := bone.GetValue(r, "id")

	voucher, err = model.FindVoucherByID(id)
	dv := voucher.VoucherData
	voucher.Status = model.ResponseStateOk

	if err == model.ErrResourceNotFound {
		vr.State = model.ErrCodeResourceNotFound
		vr.Description = model.ErrResourceNotFound.Error()
	} else if err != nil {
		vr.State = model.ResponseStateNok
		vr.Description = err.Error()
	} else {
		dvr := DetailVoucherResponse{
			ID:            dv.ID,
			VoucherCode:   dv.VoucherCode,
			ReferenceNo:   dv.ReferenceNo,
			AccountID:     dv.AccountID,
			VariantID:     dv.VariantID,
			ValidAt:       dv.ValidAt,
			ExpiredAt:     dv.ExpiredAt,
			DiscountValue: dv.DiscountValue,
			State:         dv.State,
			CreatedBy:     dv.CreatedBy,
			CreatedAt:     dv.CreatedAt,
			UpdatedBy:     dv.UpdatedBy.String,
			UpdatedAt:     dv.UpdatedAt.Time,
			DeletedBy:     dv.DeletedBy.String,
			DeletedAt:     dv.DeletedAt.Time,
			Status:        dv.Status,
		}
		vr = VoucherResponse{State: model.ResponseStateOk, Description: "", VoucherData: dvr}
	}

	res := NewResponse(vr)
	render.JSON(w, res)
}

//RedeemVoucher redeem
func RedeemVoucher(w http.ResponseWriter, r *http.Request) {
	var d model.UpdateDeleteRequest
	var vrr RedeemVoucherRequest
	var vr VoucherResponse

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&vrr); err != nil {
		log.Panic(err)
	}

	fv, err := model.FindVoucherByCode(vrr.VoucherCode)
	if err == model.ErrResourceNotFound {
		vr.State = model.ErrCodeResourceNotFound
		vr.Description = model.ErrResourceNotFound.Error()
	} else if err != nil {
		vr.State = model.ResponseStateNok
		vr.Description = err.Error()
	} else if fv.Status != model.ResponseStateOk {
		vr.State = fv.Status
		vr.Description = fv.Message
	} else {

		d.ID = fv.VoucherData.ID
		d.State = model.VoucherStateUsed
		d.User = vrr.AccountID

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
				AccountID:     uv.AccountID,
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

	res := NewResponse(vr)
	render.JSON(w, res)
}

// DeleteVoucher delete Voucher data from DB by ID
func DeleteVoucher(w http.ResponseWriter, r *http.Request) {
	var d model.UpdateDeleteRequest
	var rd DeleteVoucherRequest
	vr := VoucherResponse{State: model.ResponseStateOk, Description: ""}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		log.Panic(err)
	}

	fv, err := model.FindVoucherByCode(rd.VoucherCode)
	if err == model.ErrResourceNotFound {
		vr.State = model.ErrCodeResourceNotFound
		vr.Description = model.ErrResourceNotFound.Error()
	} else if err != nil {
		vr.State = model.ResponseStateNok
		vr.Description = err.Error()
	} else {
		d.ID = fv.VoucherData.ID
		d.State = model.VoucherStateDeleted
		d.User = fv.VoucherData.CreatedBy

		err := d.DeleteVc()
		if err != nil {
			vr.State = model.ResponseStateNok
			vr.Description = err.Error()
		} else {
			vr.State = model.ResponseStateOk
			vr.Description = ""
			vr.VoucherData = fv.VoucherData.ID
		}
	}

	res := NewResponse(vr)
	render.JSON(w, res)
}

func PayVoucher(w http.ResponseWriter, r *http.Request) {
	var d model.UpdateDeleteRequest
	var pvr PayVoucherRequest
	vr := VoucherResponse{State: model.ResponseStateOk, Description: ""}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&pvr); err != nil {
		log.Panic(err)
	}

	fv, err := model.FindVoucherByCode(pvr.VoucherCode)
	if err == model.ErrResourceNotFound {
		vr.State = model.ErrCodeResourceNotFound
		vr.Description = model.ErrResourceNotFound.Error()
	} else if err != nil {
		vr.State = model.ResponseStateNok
		vr.Description = err.Error()
	} else if fv.VoucherData.State == model.VoucherStatePaid {
		vr.State = model.ErrCodeVoucherAlreadyPaid
		vr.Description = model.ErrMessageVoucherAlreadyPaid
	} else {
		d.ID = fv.VoucherData.ID
		d.State = model.VoucherStatePaid
		d.User = fv.VoucherData.CreatedBy

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
				AccountID:     uv.AccountID,
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

	res := NewResponse(vr)
	render.JSON(w, res)
}

//GenerateVoucherOnDemand Generate singgle voucher request
func GenerateVoucherOnDemand(w http.ResponseWriter, r *http.Request) {
	var gvd GenerateVoucherRequest
	vr := VoucherResponse{State: model.ResponseStateOk, Description: ""}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&gvd); err != nil {
		log.Panic(err)
	}

	variant, err := model.FindVariantById(gvd.VariantID)
	if err == model.ErrResourceNotFound {
		vr.State = model.ErrCodeInvalidVoucher
		vr.Description = model.ErrMessageInvalidVoucher
	} else if err != nil {
		vr.State = model.ResponseStateNok
		vr.Description = err.Error()
	} else {

		if dt, ok := variant.Data.(model.Variant); ok {
			if (int(dt.MaxQuantityVoucher) - 1) <= 0 {
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

	res := NewResponse(vr)
	render.JSON(w, res)
}

//GenerateVoucher Generate bulk voucher request
func GenerateVoucher(w http.ResponseWriter, r *http.Request) {
	var gvd GenerateVoucherRequest
	vr := VoucherResponse{State: model.ResponseStateOk, Description: ""}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&gvd); err != nil {
		log.Panic(err)
	}
	variant, err := model.FindVariantById(gvd.VariantID)
	if err == model.ErrResourceNotFound {
		vr.State = model.ErrCodeResourceNotFound
		vr.Description = model.ErrMessageInvalidVoucher
	} else if err != nil {
		vr.State = model.ResponseStateNok
		vr.Description = err.Error()
	} else {

		if dt, ok := variant.Data.(model.Variant); ok {
			if (int(dt.MaxQuantityVoucher) - gvd.Quantity) <= 0 {
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

	res := NewResponse(vr)
	render.JSON(w, res)
}

// GenerateVoucher Genera te voucher and strore to DB
func (vr *GenerateVoucherRequest) generateVoucherBulk(v *model.Variant) ([]model.Voucher, error) {
	ret := make([]model.Voucher, vr.Quantity)
	var rt []string

	for i := 0; i <= vr.Quantity-1; i++ {
		rt = append(rt, randStr())

		tsd, err := time.Parse(time.RFC3339Nano, v.StartDate)
		if err != nil {
			log.Panic(err)
		}
		tea, err := time.Parse(time.RFC3339Nano, v.EndDate)
		if err != nil {
			log.Panic(err)
		}

		rd := model.Voucher{
			VoucherCode:   rt[i],
			ReferenceNo:   vr.ReferenceNo,
			AccountID:     vr.AccountID,
			VariantID:     v.Id,
			ValidAt:       tsd,
			ExpiredAt:     tea,
			DiscountValue: v.DiscountValue,
			State:         model.VoucherStateCreated,
			CreatedBy:     vr.AccountID,
			CreatedAt:     time.Now(),
		}

		if err := rd.InsertVc(); err != nil {
			log.Panic(err)
		}
		fmt.Println(i)
		ret[i] = rd
	}
	return ret, nil
}

func randStr() string {
	b := make([]byte, VoucherCodeLen)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	s := fmt.Sprintf("%X", b)

	return s
}
