package controller

import (
	"encoding/json"
	"net/http"

	"github.com/gilkor/evoucher-v2/model"
	u "github.com/gilkor/evoucher-v2/util"
	"github.com/go-zoo/bone"
)

//PostTag : POST Tag data
func PostTag(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	var reqTag model.Tag
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&reqTag); err != nil {
		res.SetError(JSONErrFatal)
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}
	tags, err := reqTag.Insert()
	if err != nil {
		res.SetErrorWithDetail(JSONErrFatal, err)
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	res.SetResponse(tags)
	res.JSON(w, res, http.StatusCreated)
}

//GetTags : GET list of Tags
func GetTags(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	qp := u.NewQueryParam(r)
	tags, next, err := model.GetTags(qp)
	if err != nil {
		u.DEBUG(err)
		res.SetError(JSONErrFatal.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	res.SetResponse(tags)
	res.SetPagination(r, qp.Page, next)
	res.JSON(w, res, http.StatusOK)
}

//GetTagByID : GET
func GetTagByID(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	qp := u.NewQueryParam(r)
	id := bone.GetValue(r, "id")
	Tag, _, err := model.GetTagByID(qp, id)
	if err != nil {
		res.SetError(JSONErrResourceNotFound)
		res.JSON(w, res, JSONErrResourceNotFound.Status)
		return
	}

	res.SetResponse(Tag)
	res.JSON(w, res, http.StatusOK)
}

//GetTagByKey : GET
func GetTagByKey(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	qp := u.NewQueryParam(r)
	key := bone.GetValue(r, "key")
	Tag, _, err := model.GetTagByKey(qp, key)
	if err != nil {
		res.SetError(JSONErrResourceNotFound)
		res.JSON(w, res, JSONErrResourceNotFound.Status)
		return
	}

	res.SetResponse(Tag)
	res.JSON(w, res, http.StatusOK)
}

// UpdateTag :
func UpdateTag(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	id := bone.GetValue(r, "id")
	var reqTag model.Tag
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&reqTag); err != nil {
		res.SetError(JSONErrFatal)
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}
	reqTag.ID = id
	if err := reqTag.Update(); err != nil {
		res.SetError(JSONErrFatal)
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}
	res.SetResponse(reqTag)
	res.JSON(w, res, http.StatusOK)
}

//DeleteTag : remove Tag
func DeleteTag(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	id := bone.GetValue(r, "id")
	p := model.Tag{ID: id}
	if err := p.Delete(); err != nil {
		u.DEBUG(err)
		res.SetError(JSONErrResourceNotFound)
		res.JSON(w, res, JSONErrResourceNotFound.Status)
		return
	}
	res.JSON(w, res, http.StatusOK)
}

//PostObjectTags : submit holder to tags
func PostObjectTags(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	var reqTag model.ObjectTag
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&reqTag); err != nil {
		res.SetError(JSONErrFatal)
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}
	// reqTag.ObjectCategory = OTG
	if err := reqTag.Insert(); err != nil {
		res.SetErrorWithDetail(JSONErrFatal, err)
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	res.SetResponse(reqTag)
	res.JSON(w, res, http.StatusCreated)
}

func PostAssignObjectTags(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	var reqTag model.ObjectTag
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&reqTag); err != nil {
		res.SetError(JSONErrBadRequest)
		res.JSON(w, res, JSONErrBadRequest.Status)
	}

	// switch reqTag.Action {
	// case "add":
	// 	status = StatusCreated
	// 	break
	// case "exist":
	// 	continue
	// case "remove":
	// 	status = StatusDeleted
	// 	break
	// default:
	// 	break
	// }

}
