package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gilkor/evoucher-v2/model"
	u "github.com/gilkor/evoucher-v2/util"
	"github.com/go-zoo/bone"
	"github.com/gorilla/schema"
)

// CreateEmailBlast : Create email blast
func CreateEmailBlast(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()
	qp := u.NewQueryParam(r)

	companyID := bone.GetValue(r, "company")
	accountToken := r.FormValue("token")

	auth, err := model.VerifyAccountToken(accountToken)
	if err != nil {
		u.DEBUG(err)
		res.SetError(JSONErrUnauthorized)
		res.JSON(w, res, JSONErrUnauthorized.Status)
		return
	}

	var blast model.Blast
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&blast); err != nil {
		u.DEBUG(err)
		res.SetError(JSONErrFatal)
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	program, err := model.GetProgramByID(blast.Program.ID, qp)
	if err != nil {
		res.SetError(JSONErrBadRequest.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrBadRequest.Status)
		return
	}

	configs, err := model.GetConfigs(companyID, "")
	if err != nil {
		res.SetError(JSONErrFatal.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	sender := configs["sender"]
	if sender != nil {
		blast.Sender = sender.(string)
	}

	templateName := configs["template_name"]
	if templateName != nil {
		blast.Template = templateName.(string)
	}

	blast.Program = program

	if blast.Template == "" || blast.Sender == "" {
		res.SetError(JSONErrBadRequest.SetMessage("Please setup the blast config"))
		res.JSON(w, res, JSONErrBadRequest.Status)
		return
	}

	// validate program channel -> should be blast
	// if program.ChannelID != "blast" {
	// 	res.SetError(JSONErrBadRequest.SetMessage("Please choose channel blast"))
	// 	res.JSON(w, res, JSONErrBadRequest.Status)
	// 	return
	// }

	voucherCreated, err := model.GetVoucherCreatedAmountByProgramID(program.ID)
	if err != nil {
		res.SetError(JSONErrFatal.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrFatal.Status)
	}

	// validate available voucher on program stock
	var availableVoucher = int(program.Stock) - voucherCreated
	if availableVoucher >= len(blast.BlastRecipient) {

		for idx, recipient := range blast.BlastRecipient {

			tmp := model.HolderDetail{
				Name:  recipient.HolderName,
				Phone: recipient.HolderPhone,
				Email: recipient.HolderEmail,
			}

			holderDetail, err := json.Marshal(tmp)
			if err != nil {
				fmt.Println(err)
				return
			}

			// need further discussion on Holder value for web voucher
			gvr := GenerateVoucherRequest{
				ProgramID:    program.ID,
				Quantity:     1,
				ReferenceNo:  "email-blast",
				HolderID:     "",
				HolderDetail: holderDetail,
				UpdatedBy:    auth.AccountID,
			}

			vouchers, err := gvr.GenerateVoucher(fmt.Sprint(configs["timezone"]), *program)
			if err != nil {
				res.SetError(JSONErrFatal.SetArgs(err.Error()))
				res.JSON(w, res, JSONErrFatal.Status)
			}

			response, err := vouchers.Insert()
			if err != nil {
				fmt.Println(err)
				res.SetError(JSONErrFatal.SetArgs(err.Error()))
				res.JSON(w, res, JSONErrFatal.Status)
				return
			}

			recipient.VoucherID = (&(*response)[0]).ID
			recipient.UpdatedBy = auth.AccountID

			// update the BlastRecipient Data
			blast.BlastRecipient[idx] = recipient

		}

		blast.UpdatedBy = auth.AccountID

		// insert blast
		response, err := blast.Insert()
		if err != nil {
			res.SetError(JSONErrFatal.SetArgs(err.Error()))
			res.JSON(w, res, JSONErrFatal.Status)
			return
		}

		res.SetResponse(response)
		res.JSON(w, res, http.StatusOK)

	} else {
		res.SetError(JSONErrFatal.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrFatal.Status)
	}

}

// UpdateBlast :
func UpdateBlast(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	id := bone.GetValue(r, "id")
	var reqBlast model.Blast
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&reqBlast); err != nil {
		res.SetErrorWithDetail(JSONErrFatal, err)
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}
	reqBlast.ID = id
	err := reqBlast.Update()
	if err != nil {
		res.SetErrorWithDetail(JSONErrFatal, err)
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}
	res.SetResponse(model.Blasts{reqBlast})
	res.JSON(w, res, http.StatusOK)
}

type BlastFilter struct {
	ID           string `schema:"id" filter:"array"`
	Sender       string `schema:"sender" filter:"array"`
	Subject      string `schema:"subject" filter:"string"`
	ProgramID    string `schema:"program_id" filter:"array"`
	BlastProgram string `schema:"program" filter:"json_array"`
	CompanyID    string `schema:"company_id" filter:"string"`
	Recipient    string `schema:"recipients" filter:"json_array"`
	CreatedAt    string `schema:"created_at" filter:"date"`
	CreatedBy    string `schema:"created_by" filter:"string"`
	UpdatedAt    string `schema:"updated_at" filter:"date"`
	UpdatedBy    string `schema:"updated_by" filter:"string"`
	Status       string `schema:"status" filter:"enum"`
}

//GetBlasts : GET list of blasts
func GetBlasts(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()
	qp := u.NewQueryParam(r)

	var decoder = schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)

	var f BlastFilter
	if err := decoder.Decode(&f, r.Form); err != nil {
		res.SetError(JSONErrFatal.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	qp.SetFilterModel(f)

	blasts, next, err := model.GetBlasts(qp)
	if err != nil {
		res.SetError(JSONErrFatal.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	res.SetResponse(blasts)
	res.SetNewPagination(r, qp.Page, next, (*blasts)[0].Count)
	res.JSON(w, res, http.StatusOK)
}

//GetBlastByID : GET
func GetBlastByID(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	qp := u.NewQueryParam(r)
	id := bone.GetValue(r, "id")

	blast, err := model.GetBlastByID(qp, id)
	if err != nil {
		res.SetError(JSONErrResourceNotFound)
		res.JSON(w, res, JSONErrResourceNotFound.Status)
		return
	}

	res.SetResponse(model.Blasts{*blast})
	res.JSON(w, res, http.StatusOK)
}

// GetTemplateByName : get template of blast by using nudge blast name
func GetTemplateByName(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()
	name := bone.GetValue(r, "name")

	template, err := model.GetBlastsTemplate(name)
	if err != nil {
		res.SetError(JSONErrResourceNotFound)
		res.JSON(w, res, JSONErrResourceNotFound.Status)
		return
	}

	res.SetResponse(template.Data)
	res.JSON(w, res, http.StatusOK)
}

// GetBlastsTemplate : get template of blast
func GetBlastsTemplate(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	qp := u.NewQueryParam(r)
	id := bone.GetValue(r, "id")

	blast, err := model.GetBlastByID(qp, id)
	if err != nil {
		res.SetError(JSONErrResourceNotFound)
		res.JSON(w, res, JSONErrResourceNotFound.Status)
		return
	}

	template, err := model.GetBlastsTemplate(blast.Template)
	if err != nil {
		res.SetError(JSONErrResourceNotFound)
		res.JSON(w, res, JSONErrResourceNotFound.Status)
		return
	}

	res.SetResponse(template.Data)
	res.JSON(w, res, http.StatusOK)
}

// SendEmailBlast : Send email blast
func SendEmailBlast(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()
	qp := u.NewQueryParam(r)
	id := bone.GetValue(r, "id")

	blast, err := model.GetBlastByID(qp, id)
	if err != nil {
		res.SetError(JSONErrResourceNotFound)
		res.JSON(w, res, JSONErrResourceNotFound.Status)
		return
	}

	if blast.Status == model.StatusCreated {
		err := blast.SendEmailBlast()
		if err != nil {
			fmt.Println(err)
			res.SetError(JSONErrFatal.SetArgs(err.Error()))
			res.JSON(w, res, JSONErrFatal.Status)
			return
		}
		res.SetResponse(model.Blasts{*blast})
		res.JSON(w, res, http.StatusOK)
	} else {
		res.SetError(JSONErrBadRequest.SetMessage("Blast already submitted"))
		res.JSON(w, res, JSONErrBadRequest.Status)
		return
	}
}
