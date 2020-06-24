package controller

import (
	"encoding/json"
	"net/http"

	"github.com/gilkor/evoucher-v2/model"
	u "github.com/gilkor/evoucher-v2/util"
	"github.com/go-zoo/bone"
	"github.com/gorilla/schema"
)

//PostTag : POST Tag data
func PostTag(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	token := r.FormValue("token")
	companyID := bone.GetValue(r, "company")

	accData, err := model.GetSessionDataJWT(token, companyID)
	if err != nil {
		res.SetError(JSONErrUnauthorized)
		res.JSON(w, res, JSONErrUnauthorized.Status)
		return
	}

	var reqTag model.Tag
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&reqTag); err != nil {
		res.SetErrorWithDetail(JSONErrFatal, err)
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	reqTag.CreatedBy = accData.AccountID
	reqTag.UpdatedBy = accData.AccountID
	response, err := reqTag.Insert()
	if err != nil {
		res.SetErrorWithDetail(JSONErrFatal, err)
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	res.SetResponse(response)
	res.JSON(w, res, http.StatusCreated)
}

type TagFilter struct {
	ID   string `schema:"id" filter:"string"`
	Name string `schema:"name" filter:"string"`
}

//GetTags : GET list of Tags
func GetTags(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()
	qp := u.NewQueryParam(r)

	qp.SetCompanyID(bone.GetValue(r, "company"))

	var decoder = schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)

	var f TagFilter
	if err := decoder.Decode(&f, r.Form); err != nil {
		res.SetError(JSONErrFatal.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	qp.SetFilterModel(f)

	tags, next, err := model.GetTags(qp)
	if err != nil {
		u.DEBUG(err)
		res.SetError(JSONErrFatal.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	res.SetResponse(tags)
	res.SetNewPagination(r, qp.Page, next, (*tags)[0].Count)
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

//GetTagByCategory : GET
func GetTagByCategory(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	qp := u.NewQueryParam(r)
	val := bone.GetValue(r, "category")
	Tag, _, err := model.GetTagByKey(qp, val)
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
	token := r.FormValue("token")

	companyID := bone.GetValue(r, "company")

	accData, err := model.GetSessionDataJWT(token, companyID)
	if err != nil {
		res.SetError(JSONErrUnauthorized)
		res.JSON(w, res, JSONErrUnauthorized.Status)
		return
	}

	id := bone.GetValue(r, "id")
	var reqTag model.Tag
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&reqTag); err != nil {
		res.SetErrorWithDetail(JSONErrFatal, err)
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}
	reqTag.ID = id
	reqTag.UpdatedBy = accData.AccountID
	err = reqTag.Update()
	if err != nil {
		res.SetErrorWithDetail(JSONErrFatal, err)
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}
	res.SetResponse(model.Tags{reqTag})
	res.JSON(w, res, http.StatusOK)
}

//DeleteTag : remove Tag
func DeleteTag(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()
	token := r.FormValue("token")

	companyID := bone.GetValue(r, "company")

	accData, err := model.GetSessionDataJWT(token, companyID)
	if err != nil {
		res.SetError(JSONErrUnauthorized)
		res.JSON(w, res, JSONErrUnauthorized.Status)
		return
	}

	id := bone.GetValue(r, "id")
	p := model.Tag{ID: id}
	p.UpdatedBy = accData.AccountID
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
		res.SetErrorWithDetail(JSONErrFatal, err)
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}
	// reqTag.ObjectCategory = OTG
	response, err := reqTag.Insert()
	if err != nil {
		res.SetErrorWithDetail(JSONErrFatal, err)
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	res.SetResponse(response)
	res.JSON(w, res, http.StatusCreated)
}

//PostAssignObjectTags :
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
