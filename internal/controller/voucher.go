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
	ErrorStateOk   string = "success" // ErrorStateOk Error state Success
	ErrorStateNok  string = "failed"  // ErrorStateNok Error state failed
	VoucherCodeLen int32  = 4         // VoucherCodeLen length of VoucherCode to run randStr
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
		ReferenceNo string `json:"reference_no"`
	}
	// VoucherResponse represent a Response of GenerateVoucher
	VoucherResponse struct {
		State       string        `json:"state"`
		Description string        `json:"messange"`
		VoucherList []VoucherData `json:"data"`
	}
	// VoucherData represent list of voucher data
	VoucherData struct {
		VoucherNo string `json:"voucher"`
	}
)

// GetVoucherDetail get Voucher detail from DB
func GetVoucherDetail(w http.ResponseWriter, r *http.Request) {
	var voucher model.VoucherResponse
	var err error
	id := bone.GetValue(r, "id")
	// code := bone.GetValue(r, "code")

	// if id != "" {
	voucher, err = model.FindVoucherByID(id)
	voucher.Status = ErrorStateOk
	// }
	// if code != "" {
	// voucher, err = model.FindVoucherByCode(id)
	// }

	if err != nil {
		voucher.Status = ErrorStateNok
		voucher.Message = err.Error()
	}

	res := NewResponse(voucher)
	render.JSON(w, res)
}

//RedeemVoucher redeem
func RedeemVoucher(w http.ResponseWriter, r *http.Request) {
	var d model.UpdateDeleteRequest
	rv := VoucherResponse{State: ErrorStateOk, Description: ""}
	// id := bone.GetValue(r, "id")
	var rd RedeemVoucherRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		log.Panic(err)
	}

	fv, err := model.FindVoucherByCode(rd.VoucherCode)
	if err != nil {
		rv.State = ErrorStateNok
		rv.Description = err.Error()
	} else if fv.Status == "Nok" {
		rv.State = ErrorStateNok
		rv.Description = fv.Message
	} else {
		d.ID = fv.VoucherData.ID

		d.State = model.VoucherStateUsed
		d.User = rd.AccountID

		err := d.UpdateVc()
		if err != nil {
			rv.State = ErrorStateNok
			rv.Description = err.Error()
		}

		rv.VoucherList = make([]VoucherData, 1)
		rv.VoucherList[0].VoucherNo = fv.VoucherData.VoucherCode
	}

	res := NewResponse(rv)
	render.JSON(w, res)
}

// DeleteVoucher delete Voucher data from DB by ID
func DeleteVoucher(w http.ResponseWriter, r *http.Request) {
	var d model.UpdateDeleteRequest
	rv := VoucherResponse{State: ErrorStateOk, Description: ""}
	// id := bone.GetValue(r, "id")

	var rd DeleteVoucherRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		log.Panic(err)
	}

	fv, err := model.FindVoucherByCode(rd.VoucherCode)
	if err != nil {
		rv.State = ErrorStateNok
		rv.Description = err.Error()
	} else {
		d.ID = fv.VoucherData.ID
		d.State = model.VoucherStateDeleted
		// note user belum ambil dari value API
		d.User = fv.VoucherData.CreatedBy

		err := d.DeleteVc()
		if err != nil {
			rv.State = ErrorStateNok
			rv.Description = err.Error()
		}

		rv.VoucherList = make([]VoucherData, 1)
		rv.VoucherList[0].VoucherNo = fv.VoucherData.VoucherCode
	}

	res := NewResponse(rv)
	render.JSON(w, res)
}

func PayVoucher(w http.ResponseWriter, r *http.Request) {
	var d model.UpdateDeleteRequest
	rv := VoucherResponse{State: ErrorStateOk, Description: ""}
	// id := bone.GetValue(r, "id")

	var rd PayVoucherRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		log.Panic(err)
	}

	fv, err := model.FindVoucherByCode(rd.VoucherCode)

	if err != nil {
		rv.State = ErrorStateNok
		rv.Description = err.Error()
	} else if fv.VoucherData.State == model.VoucherStatePaid {
		rv.State = ErrorStateNok
		rv.Description = "voucher has already Paid"
	} else {
		d.ID = fv.VoucherData.ID
		d.State = model.VoucherStatePaid
		d.User = fv.VoucherData.CreatedBy

		err := d.UpdateVc()
		if err != nil {
			rv.State = ErrorStateNok
			rv.Description = err.Error()
		}

		rv.VoucherList = make([]VoucherData, 1)
		rv.VoucherList[0].VoucherNo = fv.VoucherData.VoucherCode
	}

	res := NewResponse(rv)
	render.JSON(w, res)
}

//GenerateVoucherOnDemand Generate singgle voucher request
func GenerateVoucherOnDemand(w http.ResponseWriter, r *http.Request) {
	var rv GenerateVoucherRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rv); err != nil {
		log.Panic(err)
	}

	rd := VoucherResponse{State: ErrorStateOk, Description: ""}
	v, err := model.FindVariantByID(rv.VariantID)
	if err == model.ErrResourceNotFound {
		rd.State = ErrorStateNok
		rd.Description = "invalid variant"
	} else if err != nil {
		rd.State = ErrorStateNok
		rd.Description = err.Error()
	} else {
		d := GenerateVoucherRequest{VariantID: rv.VariantID, Quantity: 1, ReferenceNo: rv.ReferenceNo, AccountID: v.VariantValue.AccountID}
		var voucher []string
		voucher, err = d.GenerateVoucherbulk(&v.VariantValue)
		if err != nil {
			rd.State = ErrorStateNok
			rd.Description = err.Error()
		}

		rd.VoucherList = make([]VoucherData, len(voucher))
		for i, vd := range voucher {
			rd.VoucherList[i].VoucherNo = vd
		}

	}

	res := NewResponse(rd)
	render.JSON(w, res)
}

//GenerateVoucher Generate bulk voucher request
func GenerateVoucher(w http.ResponseWriter, r *http.Request) {
	rd := VoucherResponse{State: ErrorStateOk, Description: ""}
	var voucher []string
	var variant model.VariantResponse
	var rv GenerateVoucherRequest

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rv); err != nil {
		log.Panic(err)
	}
	//check stock voucher
	variant, err := model.FindVariantByID(rv.VariantID)
	if err != nil {
		rd.State = ErrorStateNok
		rd.Description = err.Error()
	} else if int(variant.VariantValue.MaxVoucher) < rv.Quantity {
		rd.State = ErrorStateNok
		rd.Description = "Out Off stock "
	} else {

		voucher, err = rv.GenerateVoucherbulk(&variant.VariantValue)
		if err != nil {
			rd.State = ErrorStateNok
			rd.Description = err.Error()
		}

		rd.VoucherList = make([]VoucherData, len(voucher))
		for i, vd := range voucher {
			rd.VoucherList[i].VoucherNo = vd
		}
	}

	res := NewResponse(rd)
	render.JSON(w, res)
}

// GenerateVoucherbulk Genera te voucher and strore to DB
func (vr *GenerateVoucherRequest) GenerateVoucherbulk(v *model.Variant) ([]string, error) {
	var rt []string

	for i := 0; i <= vr.Quantity-1; i++ {
		rt = append(rt, randStr())

		rd := model.Voucher{
			VoucherCode:   rt[i],
			ReferenceNo:   vr.ReferenceNo,
			AccountID:     vr.AccountID,
			VariantID:     vr.VariantID,
			ValidAt:       time.Now(),
			ExpiredAt:     v.FinishDate,
			VoucherType:   model.VoucherTypeCash, // v.VariantType,
			DiscountValue: v.DiscountValue,
			State:         model.VoucherStateCreated,
			PaymentType:   model.VoucherTypeCash, // default VALUES
			CreatedBy:     vr.AccountID,
			CreatedAt:     time.Now(),
		}

		if err := rd.InsertVc(); err != nil {
			log.Panic(err)
		}

	}
	return rt, nil
}

func randStr() string {
	b := make([]byte, VoucherCodeLen)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	s := fmt.Sprintf("%X", b)

	return s
}
