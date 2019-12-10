package controller

import (
	"encoding/json"
	"net/http"

	"github.com/gilkor/evoucher-v2/model"
	u "github.com/gilkor/evoucher-v2/util"
	"github.com/go-zoo/bone"
)

//PostPartner : POST partner data
func PostPartner(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	var reqPartner model.Partner
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&reqPartner); err != nil {
		res.SetErrorWithDetail(JSONErrFatal, err)
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}
	response, err := reqPartner.Insert()
	if err != nil {
		res.SetErrorWithDetail(JSONErrFatal, err)
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	res.SetResponse(response)
	res.JSON(w, res, http.StatusCreated)
}

//GetPartners : GET list of partners
func GetPartners(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	qp := u.NewQueryParam(r)

	partners, next, err := model.GetPartners(qp)
	if err != nil {
		res.SetError(JSONErrFatal.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	res.SetResponse(partners)
	res.SetPagination(r, qp.Page, next)
	res.JSON(w, res, http.StatusOK)
}

//GetPartnerByID : GET
func GetPartnerByID(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	qp := u.NewQueryParam(r)
	id := bone.GetValue(r, "id")
	partner, _, err := model.GetPartnerByID(qp, id)
	if err != nil {
		res.SetError(JSONErrResourceNotFound)
		res.JSON(w, res, JSONErrResourceNotFound.Status)
		return
	}

	res.SetResponse(partner)
	res.JSON(w, res, http.StatusOK)
}

// UpdatePartner :
func UpdatePartner(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	id := bone.GetValue(r, "id")
	var reqPartner model.Partner
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&reqPartner); err != nil {
		res.SetErrorWithDetail(JSONErrFatal, err)
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}
	reqPartner.ID = id
	err := reqPartner.Update()
	if err != nil {
		res.SetErrorWithDetail(JSONErrFatal, err)
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}
	res.SetResponse(reqPartner)
	res.JSON(w, res, http.StatusOK)
}

//DeletePartner : remove partner
func DeletePartner(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	id := bone.GetValue(r, "id")
	p := model.Partner{ID: id}
	if err := p.Delete(); err != nil {
		res.SetError(JSONErrResourceNotFound)
		res.JSON(w, res, JSONErrResourceNotFound.Status)
		return
	}
	res.JSON(w, res, http.StatusOK)
}

//PostPartnerTags : POST tags of partner
func PostPartnerTags(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	var req model.ObjectTag
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		u.DEBUG(err)
		res.SetErrorWithDetail(JSONErrFatal, err)
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}
	// reqPartner.ID = bone.GetValue(r, "holder")
	response, err := req.Insert()
	if err != nil {
		u.DEBUG(err)
		res.SetErrorWithDetail(JSONErrFatal, err)
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	res.SetResponse(response)
	res.JSON(w, res, http.StatusCreated)
}

//GetPartnerByTags : GET
func GetPartnerByTags(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	qp := u.NewQueryParam(r)
	id := bone.GetValue(r, "tag_id")
	partner, _, err := model.GetPartnersByTags(qp, id)
	if err != nil {
		u.DEBUG(err)
		res.SetError(JSONErrResourceNotFound)
		res.JSON(w, res, JSONErrResourceNotFound.Status)
		return
	}

	res.SetResponse(partner)
	res.JSON(w, res, http.StatusOK)
}

//PostBank : POST bank data
func PostBank(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	var reqBank model.Bank
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&reqBank); err != nil {
		u.DEBUG(err)
		res.SetErrorWithDetail(JSONErrFatal, err)
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}
	reqBank.PartnerID = bone.GetValue(r, "pid")
	response, err := reqBank.Insert()
	if err != nil {
		u.DEBUG(err)
		res.SetErrorWithDetail(JSONErrFatal, err)
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	res.SetResponse(response)
	res.JSON(w, res, http.StatusCreated)
}

//GetBanks : GET list of banks
func GetBanks(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	qp := u.NewQueryParam(r)

	banks, next, err := model.GetBanks(qp)
	if err != nil {
		res.SetError(JSONErrFatal.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	res.SetResponse(banks)
	res.SetPagination(r, qp.Page, next)
	res.JSON(w, res, http.StatusOK)
}

//GetBankByPartnerID : GET bank by partner id
func GetBankByPartnerID(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	qp := u.NewQueryParam(r)
	partnerID := bone.GetValue(r, "pid")
	bank, _, err := model.GetBankByPartnerID(qp, partnerID)
	if err != nil {
		res.SetError(JSONErrResourceNotFound)
		res.JSON(w, res, JSONErrResourceNotFound.Status)
		return
	}

	res.SetResponse(bank)
	res.JSON(w, res, http.StatusOK)
}

// UpdateBank :
func UpdateBank(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	id := bone.GetValue(r, "pid")
	var reqBank model.Bank
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&reqBank); err != nil {
		res.SetErrorWithDetail(JSONErrFatal, err)
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}
	reqBank.PartnerID = id
	err := reqBank.Update()
	if err != nil {
		res.SetErrorWithDetail(JSONErrFatal, err)
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}
	res.SetResponse(reqBank)
	res.JSON(w, res, http.StatusOK)
}

//DeleteBank : remove bank
func DeleteBank(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	id := bone.GetValue(r, "id")
	p := model.Bank{ID: id}
	if err := p.Delete(); err != nil {
		res.SetError(JSONErrResourceNotFound)
		res.JSON(w, res, JSONErrResourceNotFound.Status)
		return
	}
	res.JSON(w, res, http.StatusOK)
}
