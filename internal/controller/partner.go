package controller

import (
	"encoding/json"
	"net/http"

	"github.com/gilkor/evoucher/internal/model"
	u "github.com/gilkor/evoucher/internal/util"
	"github.com/go-zoo/bone"
	"github.com/ruizu/render"
)

//PostPartner : POST partner data
func PostPartner(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	var reqPartner model.Partner
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&reqPartner); err != nil {
		res.SetError(JSONErrFatal)
		render.JSON(w, res, JSONErrFatal.Status)
		return
	}
	if err := reqPartner.Insert(); err != nil {
		res.SetError(JSONErrFatal)
		render.JSON(w, res, JSONErrFatal.Status)
		return
	}

	render.JSON(w, res, http.StatusCreated)
}

//GetPartner : GET list of partners
func GetPartner(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	qp := u.NewQueryParam(r)
	partners, next, err := model.GetPartners(qp)
	if err != nil {
		res.SetError(JSONErrFatal.SetArgs(err.Error()))
		render.JSON(w, res, JSONErrFatal.Status)
		return
	}

	res.SetResponse(partners)
	res.SetPagination(r, qp.Page, next)
	render.JSON(w, res, http.StatusOK)
}

//GetPartnerByID : GET
func GetPartnerByID(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	qp := u.NewQueryParam(r)
	id := bone.GetValue(r, "id")
	partner, _, err := model.GetPartnerByID(qp, id)
	if err != nil {
		res.SetError(JSONErrResourceNotFound)
		render.JSON(w, res, JSONErrResourceNotFound.Status)
		return
	}

	res.SetResponse(partner)
	render.JSON(w, res, http.StatusOK)
}

// UpdatePartner :
func UpdatePartner(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	var reqPartner model.Partner
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&reqPartner); err != nil {
		res.SetError(JSONErrFatal)
		render.JSON(w, res, JSONErrFatal.Status)
		return
	}
	if err := reqPartner.Update(); err != nil {
		res.SetError(JSONErrFatal)
		render.JSON(w, res, JSONErrFatal.Status)
		return
	}
	render.JSON(w, res, http.StatusCreated)
}

//DeletePartner : remove partner
func DeletePartner(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	id := bone.GetValue(r, "id")
	p := model.Partner{ID: id}
	if err := p.Delete(); err != nil {
		res.SetError(JSONErrResourceNotFound)
		render.JSON(w, res, JSONErrResourceNotFound.Status)
		return
	}
	render.JSON(w, res, http.StatusCreated)
}
