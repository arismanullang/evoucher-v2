package controller

import (
	"net/http"

	"github.com/gilkor/evoucher-v2/model"
	u "github.com/gilkor/evoucher-v2/util"
)

//GetVoucherByHolder : GET list of program and vouchers by holder
func GetVoucherByHolder(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()
	qp := u.NewQueryParam(r)
	qp.Count = -1
	accountToken := r.FormValue("token")

	// var mobileVoucherData map[string]interface{}

	claims, err := model.VerifyAccountToken(accountToken)
	if err != nil {
		u.DEBUG(err)
		res.SetError(JSONErrUnauthorized)
		res.JSON(w, res, JSONErrUnauthorized.Status)
		return
	}

	vouchers, err := model.GetVouchersByHolder(claims.AccountID, qp)
	if err != nil {
		u.DEBUG(err)
		res.SetError(JSONErrBadRequest)
		res.JSON(w, res, JSONErrBadRequest.Status)
		return
	}

	distinctProgram := []string{}
	for _, v := range *vouchers {
		if !u.StringInSlice(v.ProgramID, distinctProgram) {
			distinctProgram = append(distinctProgram, v.ProgramID)
		}
	}

	listPrograms := model.Programs{}
	for _, programID := range distinctProgram {

		program := model.Program{}
		vouchersByProgram := model.Vouchers{}
		partnersByProgram := model.Partners{}

		detailProgram, err := model.GetProgramByID(programID, qp)
		if err != nil {
			u.DEBUG(err)
			res.SetError(JSONErrBadRequest)
			res.JSON(w, res, JSONErrBadRequest.Status)
			return
		}

		program.ID = detailProgram.ID
		program.Name = detailProgram.Name
		program.Type = detailProgram.Type
		program.Value = detailProgram.Value
		program.MaxValue = detailProgram.MaxValue
		program.StartDate = detailProgram.StartDate
		program.EndDate = detailProgram.EndDate
		program.Description = detailProgram.Description
		program.ImageURL = detailProgram.ImageURL
		program.Price = detailProgram.Price
		program.ProgramChannels = detailProgram.ProgramChannels
		program.State = detailProgram.State
		program.Status = detailProgram.Status

		for _, voucher := range *vouchers {
			if voucher.ProgramID == programID {
				tempVoucher := model.Voucher{
					ID:        voucher.ID,
					Code:      voucher.Code,
					ExpiredAt: voucher.ExpiredAt,
					ValidAt:   voucher.ValidAt,
					State:     voucher.State,
				}
				vouchersByProgram = append(vouchersByProgram, tempVoucher)
			}
		}
		program.Vouchers = vouchersByProgram

		for _, outlet := range detailProgram.Partners {
			tempOutlet := model.Partner{
				ID:          outlet.ID,
				Name:        outlet.Name,
				Description: outlet.Description,
				Status:      outlet.Status,
			}
			partnersByProgram = append(partnersByProgram, tempOutlet)
		}

		program.Partners = partnersByProgram

		listPrograms = append(listPrograms, program)
	}

	res.SetResponse(listPrograms)
	res.JSON(w, res, http.StatusOK)
}
