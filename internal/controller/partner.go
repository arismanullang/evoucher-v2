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
		res.SetError(ErrFatal)
		render.JSON(w, res, ErrFatal.Status)
		return
	}
	if err := reqPartner.Insert(); err != nil {
		res.SetError(ErrFatal)
		render.JSON(w, res, ErrFatal.Status)
		return
	}

	render.JSON(w, res, http.StatusCreated)
}

//GetPartner : GET list of partners
func GetPartner(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	f := u.NewFilter(r)
	partners, next, err := model.GetPartners(f)
	if err != nil {
		res.SetError(ErrFatal.SetArgs(err.Error()))
		render.JSON(w, res, ErrFatal.Status)
		return
	}

	res.SetResponse(partners)
	res.SetPagination(r, f.Page, next)
	render.JSON(w, res, http.StatusOK)
}

//GetPartnerByID : GET
func GetPartnerByID(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	f := u.NewFilter(r)
	id := bone.GetValue(r, "id")
	partner, _, err := model.GetPartnerByID(f, id)
	if err != nil {
		res.SetError(ErrResourceNotFound)
		render.JSON(w, res, ErrResourceNotFound.Status)
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
		res.SetError(ErrFatal)
		render.JSON(w, res, ErrFatal.Status)
		return
	}
	if err := reqPartner.Update(); err != nil {
		res.SetError(ErrFatal)
		render.JSON(w, res, ErrFatal.Status)
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
		res.SetError(ErrResourceNotFound)
		render.JSON(w, res, ErrResourceNotFound.Status)
		return
	}
	render.JSON(w, res, http.StatusCreated)
}
