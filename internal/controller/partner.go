package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gilkor/evoucher/internal/model"
	u "github.com/gilkor/evoucher/internal/util"
	"github.com/go-zoo/bone"
	"github.com/ruizu/render"
)

//PostPartner : POST partner data
func PostPartner(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse(nil)

	var reqPartner model.Partner
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&reqPartner); err != nil {
		res.SetError(ErrFatal)
	}
	if err := reqPartner.Insert(); err != nil {
		fmt.Println(err)
	}

	render.JSON(w, res, http.StatusCreated)
}

//GetPartner : GET list of partners
func GetPartner(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse(nil)

	f := u.NewFilter(r)
	partners, next, err := model.GetPartners(f)
	if err != nil {
		res.SetError(ErrFatal.SetArgs(err.Error()))
	}

	res.SetResponse(partners)
	res.SetPagination(r, f.Page, next)
	render.JSON(w, res, http.StatusOK)
}

//GetPartnerByID : GET
func GetPartnerByID(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse(nil)

	f := u.NewFilter(r)
	id := bone.GetValue(r, "id")
	partner, _, err := model.GetPartnerByID(f, id)
	if err != nil {
		fmt.Println(err)
		res.SetError(ErrResourceNotFound)
	}

	res.SetResponse(partner)
	render.JSON(w, res, http.StatusOK)
}

// UpdatePartner :
func UpdatePartner(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse(nil)

	var reqPartner model.Partner
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&reqPartner); err != nil {
		res.SetError(ErrFatal)
	}
	if err := reqPartner.Update(); err != nil {
		fmt.Println(err)
	}
	render.JSON(w, res, http.StatusCreated)
}

//DeleltePartner : remove partner
func DeleltePartner(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse(nil)

	id := bone.GetValue(r, "id")
	p := model.Partner{ID: id}
	if err := p.Delete(); err != nil {
		res.SetError(ErrResourceNotFound)
	}
	render.JSON(w, res, http.StatusCreated)
}
