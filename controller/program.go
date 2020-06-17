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

//PostProgram : post create program api
func PostProgram(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()
	qp := u.NewQueryParam(r)

	token := r.FormValue("token")
	accData, err := model.GetSessionDataJWT(token)
	if err != nil {
		res.SetError(JSONErrUnauthorized)
		res.JSON(w, res, JSONErrUnauthorized.Status)
		return
	}

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

	//var rules model.RulesExpression
	//program.Rule.Unmarshal(&rules)
	//
	//for k, val := range rules.And {
	//	u.DEBUG(k, "----> ", val, ":s:", len(val.In))
	//}
	//ras := (rules.And["some"])
	//u.DEBUG("isempty", ras.IsEmpty())

	//insert program -> partners
	program.CreatedBy = accData.AccountID
	response, err := program.Insert()
	if err != nil {
		res.SetError(JSONErrFatal.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	// generate voucher

	res.SetResponse(response)
	res.JSON(w, res, http.StatusCreated)
}

type ProgramFilter struct {
	ID              string `schema:"id" filter:"array"`
	CompanyID       string `schema:"company_id" filter:"string"`
	Name            string `schema:"name" filter:"string"`
	Type            string `schema:"type" filter:"enum"`
	Value           string `schema:"value" filter:"number"`
	MaxValue        string `schema:"max_value" filter:"number"`
	StartDate       string `schema:"start_date" filter:"date"`
	EndDate         string `schema:"end_date" filter:"date"`
	Description     string `schema:"description" filter:"record"`
	State           string `schema:"state" filter:"enum"`
	Stock           string `schema:"stock" filter:"number"`
	CreatedAt       string `schema:"created_at" filter:"date"`
	CreatedBy       string `schema:"created_by" filter:"string"`
	UpdatedAt       string `schema:"updated_at" filter:"date"`
	UpdatedBy       string `schema:"updated_by" filter:"string"`
	Status          string `schema:"status" filter:"enum"`
	VoucherFormat   string `schema:"voucher_format" filter:"record"`
	IsReimburse     string `schema:"is_reimburse" filter:"bool"`
	Price           string `schema:"price" filter:"number"`
	ChannelID       string `schema:"channel_id" filter:"string"`
	ProgramChannels string `schema:"program_channels" filter:"json_array"`
}

// GetProgram :
func GetProgram(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()
	qp := u.NewQueryParam(r)

	qp.SetCompanyID(bone.GetValue(r, "company"))

	var decoder = schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)

	var f ProgramFilter
	if err := decoder.Decode(&f, r.Form); err != nil {
		res.SetError(JSONErrFatal.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	qp.SetFilterModel(f)

	programs, next, err := model.GetPrograms(qp)
	if err != nil {
		res.SetError(JSONErrBadRequest.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrBadRequest.Status)
		return
	}

	res.SetResponse(programs)
	if len(*programs) > 0 {
		res.SetNewPagination(r, qp.Page, next, (*programs)[0].Count)
	}
	res.JSON(w, res, http.StatusOK)
}

//GetProgramByID :
func GetProgramByID(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	qp := u.NewQueryParam(r)
	companyID := bone.GetValue(r, "company")
	id := bone.GetValue(r, "id")

	qp.SetCompanyID(companyID)
	program, err := model.GetProgramByID(id, qp)
	if err != nil {
		res.SetError(JSONErrBadRequest.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrBadRequest.Status)
		return
	}

	res.SetResponse(program)
	res.JSON(w, res, http.StatusOK)
}

//UpdateProgram :
func UpdateProgram(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	id := bone.GetValue(r, "id")
	token := r.FormValue("token")

	accData, err := model.GetSessionDataJWT(token)
	if err != nil {
		res.SetError(JSONErrUnauthorized)
		res.JSON(w, res, JSONErrUnauthorized.Status)
		return
	}

	var req model.Program
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		res.SetErrorWithDetail(JSONErrFatal, err)
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}
	req.ID = id
	req.UpdatedBy = accData.AccountID
	err = req.Update()
	if err != nil {
		res.SetErrorWithDetail(JSONErrFatal, err)
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}
	res.SetResponse(model.Programs{req})
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

//UploadProgramImage : post upload and update program image
func UploadProgramImage(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()
	qp := u.NewQueryParam(r)

	companyID := bone.GetValue(r, "company")
	id := bone.GetValue(r, "id")

	program, err := model.GetProgramByID(id, qp)
	if err != nil {
		res.SetError(JSONErrBadRequest.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrBadRequest.Status)
		return
	}

	err = r.ParseMultipartForm(2 << 20)
	if err != nil {
		res.SetError(JSONErrBadRequest.SetArgs(err.Error(), "parse error"))
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	sourceURL, err := UploadFileFromForm(r, id, companyID+"/programs/"+id+"/")
	if err != nil {
		res.SetError(JSONErrBadRequest.SetArgs(err.Error(), "upload fail"))
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	program.ImageURL = sourceURL

	err = program.Update()
	if err != nil {
		res.SetError(JSONErrFatal.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}
	res.SetResponse(model.Programs{*program})
	res.JSON(w, res, http.StatusOK)
}

//GetProgramsByChannel : get program by channel id
func GetProgramsByChannel(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()
	r.Form.Set("fields", model.MProgramFields)
	qp := u.NewQueryParam(r)

	qp.SetCompanyID(bone.GetValue(r, "company"))
	channelID := bone.GetValue(r, "channel_id")

	var decoder = schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)

	var f ProgramFilter
	if err := decoder.Decode(&f, r.Form); err != nil {
		res.SetError(JSONErrFatal.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	f.ChannelID = channelID

	qp.SetFilterModel(f)

	programs, next, err := model.GetPrograms(qp)
	if err != nil {
		res.SetError(JSONErrBadRequest.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrBadRequest.Status)
		return
	}

	res.SetResponse(programs)
	if len(*programs) > 0 {
		res.SetNewPagination(r, qp.Page, next, (*programs)[0].Count)
	}
	res.JSON(w, res, http.StatusOK)
}

//GetMProgramByID : get program by id with fields for 3rd party
func GetMProgramByID(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()
	r.Form.Set("fields", model.MProgramFields)
	qp := u.NewQueryParam(r)
	companyID := bone.GetValue(r, "company")
	id := bone.GetValue(r, "id")

	qp.SetCompanyID(companyID)
	program, err := model.GetProgramByID(id, qp)
	if err != nil {
		res.SetError(JSONErrBadRequest.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrBadRequest.Status)
		return
	}

	qp2 := u.NewQueryParam(r)
	r.Form.Set("fields", model.MOutletFields)
	qp2.SetCompanyID(companyID)
	partners, _, err := model.GetPartnerByProgramID(id, qp2)
	if err != nil {
		res.SetError(JSONErrBadRequest.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrBadRequest.Status)
		return
	}

	program.Partners = *partners

	res.SetResponse(program)
	res.JSON(w, res, http.StatusOK)
}
