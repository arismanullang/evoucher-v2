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

//PostTag : POST Tag data
func PostTag(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse(nil)

	var reqTag model.Tag
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&reqTag); err != nil {
		res.SetError(ErrFatal)
	}
	if err := reqTag.Insert(); err != nil {
		fmt.Println(err)
	}

	render.JSON(w, res, http.StatusCreated)
}

//GetTag : GET list of Tags
func GetTag(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse(nil)

	f := u.NewFilter(r)
	Tags, next, err := model.GetTags(f)
	if err != nil {
		res.SetError(ErrFatal.SetArgs(err.Error()))
	}

	res.SetResponse(Tags)
	res.SetPagination(r, f.Page, next)
	render.JSON(w, res, http.StatusOK)
}

//GetTagByID : GET
func GetTagByID(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse(nil)

	f := u.NewFilter(r)
	id := bone.GetValue(r, "id")
	Tag, _, err := model.GetTagByID(f, id)
	if err != nil {
		fmt.Println(err)
		res.SetError(ErrResourceNotFound)
	}

	res.SetResponse(Tag)
	render.JSON(w, res, http.StatusOK)
}

// UpdateTag :
func UpdateTag(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse(nil)

	var reqTag model.Tag
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&reqTag); err != nil {
		res.SetError(ErrFatal)
	}
	if err := reqTag.Update(); err != nil {
		fmt.Println(err)
	}
	render.JSON(w, res, http.StatusCreated)
}

//DelelteTag : remove Tag
func DelelteTag(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse(nil)

	id := bone.GetValue(r, "id")
	p := model.Tag{ID: id}
	if err := p.Delete(); err != nil {
		res.SetError(ErrResourceNotFound)
	}
	render.JSON(w, res, http.StatusCreated)
}
