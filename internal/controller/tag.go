package controller

import (
	"encoding/json"
	"net/http"

	"github.com/gilkor/evoucher/internal/model"
	u "github.com/gilkor/evoucher/internal/util"
	"github.com/go-zoo/bone"
	"github.com/ruizu/render"
)

//PostTag : POST Tag data
func PostTag(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	var reqTag model.Tag
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&reqTag); err != nil {
		res.SetError(JSONErrFatal)
		render.JSON(w, res, JSONErrFatal.Status)
		return
	}
	if err := reqTag.Insert(); err != nil {
		res.SetError(JSONErrFatal)
		render.JSON(w, res, JSONErrFatal.Status)
		return
	}

	render.JSON(w, res, http.StatusCreated)
}

//GetTag : GET list of Tags
func GetTag(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	qp := u.NewQueryParam(r)
	Tags, next, err := model.GetTags(qp)
	if err != nil {
		res.SetError(JSONErrFatal.SetArgs(err.Error()))
		render.JSON(w, res, JSONErrFatal.Status)
		return
	}

	res.SetResponse(Tags)
	res.SetPagination(r, qp.Page, next)
	render.JSON(w, res, http.StatusOK)
}

//GetTagByID : GET
func GetTagByID(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	qp := u.NewQueryParam(r)
	id := bone.GetValue(r, "id")
	Tag, _, err := model.GetTagByID(qp, id)
	if err != nil {
		res.SetError(JSONErrResourceNotFound)
		render.JSON(w, res, JSONErrResourceNotFound.Status)
		return
	}

	res.SetResponse(Tag)
	render.JSON(w, res, http.StatusOK)
}

// UpdateTag :
func UpdateTag(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	var reqTag model.Tag
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&reqTag); err != nil {
		res.SetError(JSONErrFatal)
		render.JSON(w, res, JSONErrFatal.Status)
		return
	}
	if err := reqTag.Update(); err != nil {
		res.SetError(JSONErrFatal)
		render.JSON(w, res, JSONErrFatal.Status)
		return
	}
	render.JSON(w, res, http.StatusCreated)
}

//DeleteTag : remove Tag
func DeleteTag(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	id := bone.GetValue(r, "id")
	p := model.Tag{ID: id}
	if err := p.Delete(); err != nil {
		res.SetError(JSONErrResourceNotFound)
		render.JSON(w, res, JSONErrResourceNotFound.Status)
		return
	}
	render.JSON(w, res, http.StatusCreated)
}
