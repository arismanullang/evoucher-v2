package controller

import (
	"net/http"

	"github.com/gilkor/evoucher/internal/model"
	"github.com/ruizu/render"
	"time"
)

type (
	ChallengeResponse struct {
		Challenge string `json:"challenge"`
		Timeout   string `json:"timeout"`
		Duration  int    `json:"duration"`
	}
)

func GetChallenge(w http.ResponseWriter, r *http.Request) {
	status := http.StatusOK

	c := randStr(model.CHALLENGE_LENGTH, model.CHALLENGE_FORMAT)
	d := model.TIMEOUT_DURATION
	t := time.Now().Add(time.Second * time.Duration(d))

	res := NewResponse(ChallengeResponse{Challenge: c, Timeout: t.Format("2006-01-02 15:04:05.000"), Duration: d})
	render.JSON(w, res, status)
}

func GetRedeemData(w http.ResponseWriter, r *http.Request) {
	status := http.StatusOK
	res := NewResponse(nil)
	VcID := r.FormValue("x")

	voucher, err := model.FindVoucher(map[string]string{"id": StrDecode(VcID)})
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

	variant, err := model.FindVariantDetailsById(voucher.VoucherData[0].VariantID)
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
	partner, err := model.FindVariantPartner(map[string]string{"variant_id": variant.Id})
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
	vcr := make([]VoucerResponse, 1)
	vcr[0].VoucherID = voucher.VoucherData[0].ID
	vcr[0].VoucherNo = voucher.VoucherData[0].VoucherCode
	vcr[0].State = voucher.VoucherData[0].State

	d := GetVoucherOfVariatListDetails{
		VariantID:          variant.Id,
		AccountId:          variant.AccountId,
		VariantName:        variant.VariantName,
		VoucherType:        variant.VoucherType,
		VoucherPrice:       variant.VoucherPrice,
		AllowAccumulative:  variant.AllowAccumulative,
		StartDate:          variant.StartDate,
		EndDate:            variant.EndDate,
		DiscountValue:      variant.DiscountValue,
		MaxQuantityVoucher: variant.MaxQuantityVoucher,
		MaxUsageVoucher:    variant.MaxUsageVoucher,
		RedeemtionMethod:   variant.RedeemtionMethod,
		ImgUrl:             variant.ImgUrl,
		VariantTnc:         variant.VariantTnc,
		VariantDescription: variant.VariantDescription,
		State:              voucher.VoucherData[0].State,
		Holder:             voucher.VoucherData[0].Holder.String,
		HolderDescription:  voucher.VoucherData[0].HolderDescription.String,
		Voucher:            vcr,
	}

	d.Partners = make([]Partner, len(partner))
	for i, pd := range partner {
		d.Partners[i].ID = pd.Id
		d.Partners[i].PartnerName = pd.PartnerName
		d.Partners[i].SerialNumber = pd.SerialNumber.String
	}

	d.Used = getCountVoucher(variant.Id)

	res = NewResponse(d)
	render.JSON(w, res, status)
}
