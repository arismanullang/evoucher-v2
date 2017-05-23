package controller

import (
	"net/http"

	"github.com/gilkor/evoucher/internal/model"
	"github.com/ruizu/render"
)

func RedeemPage(w http.ResponseWriter, r *http.Request) {
	render.FileInLayout(w, "layout.html", "redeem.html", nil)
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
