package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gilkor/evoucher-v2/model"
	u "github.com/gilkor/evoucher-v2/util"
	"github.com/go-zoo/bone"
)

//PostProgram : post create program api
func PostProgram(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()
	qp := u.NewQueryParam(r)

	var program *model.Program
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&program); err != nil {
		res.SetError(JSONErrFatal.SetMessage(err.Error()))
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}
	// TO-DO
	//validate partners(?)
	for _, v := range program.Partners {
		_, _, err := model.GetPartnerByID(qp, v.ID)
		if err != nil {
			u.DEBUG(JSONErrBadRequest.Message, " OutletID:", v.ID)
			res.SetError(JSONErrBadRequest.SetMessage(err.Error()))
			res.JSON(w, res, JSONErrBadRequest.Status)
			return
		}
	}
	//insert program -> partners
	if err := program.Insert(); err != nil {
		fmt.Println(err)
		res.SetError(JSONErrFatal.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	// generate voucher

	res.SetResponse(program)
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
