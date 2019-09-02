package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gilkor/evoucher/model"
	u "github.com/gilkor/evoucher/util"
	"github.com/go-zoo/bone"
)

//PostProgram : POST partner data
func PostProgram(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()
	// qp := u.NewQueryParam(r)

	var program *model.Program
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&program); err != nil {
		res.SetError(JSONErrFatal.SetMessage(err.Error()))
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	// TO-DO
	//validate partners(?)
	// for _, v := range program.Partners {
	// 	_, _, err := model.GetPartnerByID(qp, v.ID)
	// 	if err != nil {
	// 		u.DEBUG(err)
	// 		res.SetError(JSONErrBadRequest.SetMessage(err.Error()))
	// 		res.JSON(w, res, JSONErrBadRequest.Status)
	// 		return
	// 	}
	// }
	//insert program -> partners
	// fmt.Println(program)
	if err := program.Insert(); err != nil {
		fmt.Println(err)
		res.SetError(JSONErrFatal.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	// generate voucher
	var vouchers model.Vouchers
	var vf model.VoucherFormat
	program.VoucherFormat.Unmarshal(&vf)

	// generateVoucher := func(str string, i []interface{}) {
	// 	var code string
	// 	// specific code for coupon
	// 	// else generate random unique code
	// 	if vf.IsSpecifiedCode() {
	// 		code = vf.Properties.Code
	// 	} else {
	// 		code = vf.Properties.Prefix + str + vf.Properties.Postfix
	// 	}
	// 	vc := model.Voucher{
	// 		Code:      code,
	// 		ProgramID: program.ID,
	// 		CreatedBy: program.CreatedBy,
	// 		UpdatedBy: &program.CreatedBy,
	// 	}
	// 	// u.DEBUG("LAMHOT", vc.ProgramID)

	// 	i = append(i, vc)
	// }

	// generate stock
	gco := NewGenerateCode()
	codes := gco.GetUniqueStrings(vf.Properties.Length, int(program.Stock))
	for _, code := range codes {
		vouchers = append(vouchers, model.Voucher{
			Code:      code,
			ProgramID: program.ID,
			CreatedBy: program.CreatedBy,
			UpdatedBy: &program.CreatedBy,
		})
	}

	if int64(len(vouchers)) != program.Stock {
		u.DEBUG("[", len(vouchers), "|", program.Stock, "]Stock Generate went wrong, terminate request. Please Contact us if this error frequently happen.")
		res.SetError(JSONErrFatal.SetArgs("Stock Generate went wrong, terminate request. Please Contact us if this error frequently happen."))
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	if err := vouchers.Insert(); err != nil {
		u.DEBUG("zz_generate_voucher", err)
		fmt.Println(err)
		res.SetError(JSONErrFatal.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	// res.SetResponse(program)
	res.JSON(w, res, http.StatusCreated)
}

// GetProgram :
func GetProgram(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	qp := u.NewQueryParam(r)
	programs, next, err := model.GetPrograms(qp)
	if err != nil {
		res.SetError(JSONErrBadRequest.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrBadRequest.Status)
		return
	}

	res.SetResponse(programs)
	res.SetPagination(r, qp.Page, next)
	res.JSON(w, res, http.StatusOK)
}

//GetProgramByID :
func GetProgramByID(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	qp := u.NewQueryParam(r)
	id := bone.GetValue(r, "id")
	program, err := model.GetProgramByID(id, qp)
	if err != nil {
		res.SetError(JSONErrBadRequest.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrBadRequest.Status)
		return
	}

	res.SetResponse(program)
	res.JSON(w, res, http.StatusOK)
}

// DeleteProgram :
func DeleteProgram(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	qp := u.NewQueryParam(r)
	id := bone.GetValue(r, "id")
	fmt.Println("program id ", id)
	program, err := model.GetProgramByID(id, qp)
	if err != nil {
		res.SetError(JSONErrBadRequest.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrBadRequest.Status)
		return
	}
	// Delete program
	fmt.Println("delete prog ", program)
	if err := program.Delete(); err != nil {
		res.SetError(JSONErrFatal.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}
	res.JSON(w, res, http.StatusOK)
}
