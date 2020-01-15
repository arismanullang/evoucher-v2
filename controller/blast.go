package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gilkor/evoucher-v2/model"
	u "github.com/gilkor/evoucher-v2/util"
)

// EmailBlast : POST email blast
func EmailBlast(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()
	qp := u.NewQueryParam(r)

	var reqBlast model.Blast
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&reqBlast); err != nil {
		u.DEBUG(err)
		res.SetError(JSONErrFatal)
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	program, err := model.GetProgramByID(reqBlast.Program.ID, qp)
	if err != nil {
		res.SetError(JSONErrBadRequest.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrBadRequest.Status)
		return
	}

	// Update program value needed on the blast
	reqBlast.Program = *program

	// get program detail
	// validate program channel -> should be blast
	// validate available voucher on program stock

	for _, recipient := range reqBlast.RecipientsData {
		// generate voucher for every recipient
		recipient.VoucherURL = "voucher-staging.elys.id"
	}

	if err := model.SendEmailBlast(reqBlast); err != nil {
		fmt.Println(err)
		res.SetError(JSONErrFatal.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	res.SetResponse(reqBlast)
	res.JSON(w, res, http.StatusOK)
}
