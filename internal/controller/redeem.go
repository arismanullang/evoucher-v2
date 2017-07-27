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

	program, err := model.FindProgramDetailsById(voucher.VoucherData[0].ProgramID)
	if err == model.ErrResourceNotFound {
		status = http.StatusNotFound
		res.AddError(its(status), model.ErrCodeResourceNotFound, model.ErrMessageInvalidProgram, "voucher")
		render.JSON(w, res, status)
		return
	} else if err != nil {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", "voucher")
		render.JSON(w, res, status)
		return
	}
	partner, err := model.FindProgramPartner(map[string]string{"program_id": program.Id})
	if err == model.ErrResourceNotFound {
		status = http.StatusNotFound
		res.AddError(its(status), model.ErrCodeResourceNotFound, model.ErrMessageInvalidProgram+"(Partner of Program Not Found)", "voucher")
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
		ProgramID:          program.Id,
		AccountId:          program.AccountId,
		ProgramName:        program.Name,
		VoucherType:        program.VoucherType,
		VoucherPrice:       program.VoucherPrice,
		AllowAccumulative:  program.AllowAccumulative,
		StartDate:          program.StartDate,
		EndDate:            program.EndDate,
		VoucherValue:       program.VoucherValue,
		MaxQuantityVoucher: program.MaxQuantityVoucher,
		MaxGenerateVoucher: program.MaxGenerateVoucher,
		RedeemtionMethod:   program.RedeemtionMethod,
		ImgUrl:             program.ImgUrl,
		ProgramTnc:         program.Tnc,
		ProgramDescription: program.Description,
		State:              voucher.VoucherData[0].State,
		Holder:             voucher.VoucherData[0].Holder.String,
		HolderDescription:  voucher.VoucherData[0].HolderDescription.String,
		Voucher:            vcr,
	}

	d.Partners = make([]Partner, len(partner))
	for i, pd := range partner {
		d.Partners[i].ID = pd.Id
		d.Partners[i].Name = pd.Name
		d.Partners[i].SerialNumber = pd.SerialNumber.String
	}

	d.Used = getCountVoucher(program.Id)

	res = NewResponse(d)
	render.JSON(w, res, status)
}
