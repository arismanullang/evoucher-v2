package controller

import (
	"encoding/json"
	"net/http"

	"github.com/gilkor/evoucher-v2/model"
	u "github.com/gilkor/evoucher-v2/util"
	// "github.com/go-zoo/bone"
)

//PostVoucherAssignHolder :
func PostVoucherAssignHolder(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	var req model.Voucher
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	req.State = model.VoucherStateUsed
	if err != nil {
		u.DEBUG(err)
		res.SetError(JSONErrBadRequest)
		res.JSON(w, res, JSONErrBadRequest.Status)
		return
	}
	if err := req.Update(); err != nil {
		u.DEBUG(err)
		res.SetErrorWithDetail(JSONErrFatal, err)
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}
	res.JSON(w, res, http.StatusCreated)
}

//PostVoucherUse :
func PostVoucherUse(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	var req model.Voucher
	var rule model.Rules

	decoder := json.NewDecoder(r.Body)
	qp := u.NewQueryParam(r)
	err := decoder.Decode(&req)
	req.State = model.VoucherStateUsed

	//Validate Rule Program
	program, err := model.GetProgramByID(req.ProgramID, qp)
	if err != nil {

	}
	err = rule.Unmarshal(program.Rule)
	if err != nil {
		u.DEBUG(err)
		res.SetError(JSONErrBadRequest)
		res.JSON(w, res, JSONErrBadRequest.Status)
		return
	}

	result, err := rule.Validate()
	if result {
		//Unexpected error
		if err != nil {
			u.DEBUG(err)
			res.SetError(JSONErrBadRequest)
			res.JSON(w, res, JSONErrBadRequest.Status)
			return
		}
		//expected error
		res.SetError(JSONErrBadRequest)
		res.JSON(w, res, JSONErrBadRequest.Status)
		return
	}

	//
	//
	//
	if err != nil {
		u.DEBUG(err)
		res.SetError(JSONErrBadRequest)
		res.JSON(w, res, JSONErrBadRequest.Status)
		return
	}
	if err = req.Update(); err != nil {
		u.DEBUG(err)
		res.SetErrorWithDetail(JSONErrFatal, err)
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}
	res.JSON(w, res, http.StatusCreated)
}

type (
	VoucherUseRequest struct {
		Reference    string                 `json:'reference'`
		Transactions VoucherUseTransactions `json:"transactions"`
		Vouchers     []string               `json:"vouchers"`
		OutletID     string                 `json:"outlet_id"`
	}
	VoucherUseTransactions struct {
		TotalAmount float64 `json:"total_amount"`
		Details     string  `json:"details"`
	}
)

//PostVoucherUse :
func PostVoucherUset(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	var req model.Voucher
	var rule model.Rules

	decoder := json.NewDecoder(r.Body)
	qp := u.NewQueryParam(r)
	err := decoder.Decode(&req)
	req.State = model.VoucherStateUsed

	//Validate Rule Program
	program, err := model.GetProgramByID(req.ProgramID, qp)
	if err != nil {

	}
	err = rule.Unmarshal(program.Rule)
	if err != nil {
		u.DEBUG(err)
		res.SetError(JSONErrBadRequest)
		res.JSON(w, res, JSONErrBadRequest.Status)
		return
	}

	result, err := rule.Validate()
	if result {
		//Unexpected error
		if err != nil {
			u.DEBUG(err)
			res.SetError(JSONErrBadRequest)
			res.JSON(w, res, JSONErrBadRequest.Status)
			return
		}
		//expected error
		res.SetError(JSONErrBadRequest)
		res.JSON(w, res, JSONErrBadRequest.Status)
		return
	}
	//
	//
	//
	if err != nil {
		u.DEBUG(err)
		res.SetError(JSONErrBadRequest)
		res.JSON(w, res, JSONErrBadRequest.Status)
		return
	}
	if err = req.Update(); err != nil {
		u.DEBUG(err)
		res.SetErrorWithDetail(JSONErrFatal, err)
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}
	res.JSON(w, res, http.StatusCreated)
}

type (
	VoucherClaimRequest struct {
	}
)

//PostVoucherClaim :
func PostVoucherClaim(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	var req model.Voucher
	var rule model.Rules

	decoder := json.NewDecoder(r.Body)
	qp := u.NewQueryParam(r)
	err := decoder.Decode(&req)
	req.State = model.VoucherStateUsed

	//Validate Rule Program
	program, err := model.GetProgramByID(req.ProgramID, qp)
	if err != nil {

	}
	err = rule.Unmarshal(program.Rule)
	if err != nil {
		u.DEBUG(err)
		res.SetError(JSONErrBadRequest)
		res.JSON(w, res, JSONErrBadRequest.Status)
		return
	}

	result, err := rule.Validate()
	if result {
		//Unexpected error
		if err != nil {
			u.DEBUG(err)
			res.SetError(JSONErrBadRequest)
			res.JSON(w, res, JSONErrBadRequest.Status)
			return
		}
		//expected error
		res.SetError(JSONErrBadRequest)
		res.JSON(w, res, JSONErrBadRequest.Status)
		return
	}
	//
	//
	//
	if err != nil {
		u.DEBUG(err)
		res.SetError(JSONErrBadRequest)
		res.JSON(w, res, JSONErrBadRequest.Status)
		return
	}
	if err = req.Update(); err != nil {
		u.DEBUG(err)
		res.SetErrorWithDetail(JSONErrFatal, err)
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}
	res.JSON(w, res, http.StatusCreated)
}
