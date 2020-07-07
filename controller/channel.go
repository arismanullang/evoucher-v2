package controller

import (
	"encoding/json"
	"net/http"

	"github.com/gilkor/evoucher-v2/model"
	u "github.com/gilkor/evoucher-v2/util"
	"github.com/go-zoo/bone"
	"github.com/gorilla/schema"
)

//PostChannel : POST channel data
func PostChannel(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()
	token := r.FormValue("token")

	companyID := bone.GetValue(r, "company")

	accData, err := model.GetSessionDataJWT(token, companyID)
	if err != nil {
		res.SetError(JSONErrUnauthorized)
		res.JSON(w, res, JSONErrUnauthorized.Status)
		return
	}

	var reqChannel model.Channel
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&reqChannel); err != nil {
		u.DEBUG(err)
		res.SetErrorWithDetail(JSONErrFatal, err)
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}
	reqChannel.CreatedBy = accData.AccountID
	reqChannel.UpdatedBy = accData.AccountID
	response, err := reqChannel.Insert()
	if err != nil {
		u.DEBUG(err)
		res.SetErrorWithDetail(JSONErrFatal, err)
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	res.SetResponse(response)
	res.JSON(w, res, http.StatusCreated)
}

type ChannelFilter struct {
	ID          string `schema:"id" filter:"array"`
	Name        string `schema:"name" filter:"string"`
	Description string `schema:"description" filter:"record"`
	IsSuper     string `schema:"is_super" filter:"bool"`
	CreatedAt   string `schema:"created_at" filter:"date"`
	CreatedBy   string `schema:"created_by" filter:"string"`
	UpdatedAt   string `schema:"updated_at" filter:"date"`
	UpdatedBy   string `schema:"updated_by" filter:"string"`
	Status      string `schema:"status" filter:"enum"`
	Tags        string `schema:"channel_tags" filter:"json"`
}

//GetChannels : GET list of channels
func GetChannels(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()
	qp := u.NewQueryParam(r)

	qp.SetCompanyID(bone.GetValue(r, "company"))

	var decoder = schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)

	var f ChannelFilter
	if err := decoder.Decode(&f, r.Form); err != nil {
		res.SetError(JSONErrFatal.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	qp.SetFilterModel(f)

	channels, next, err := model.GetChannels(qp)
	if err != nil {
		res.SetError(JSONErrFatal.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	res.SetResponse(channels)
	if len(*channels) > 0 {
		res.SetNewPagination(r, qp.Page, next, (*channels)[0].Count)
	}
	res.JSON(w, res, http.StatusOK)
}

//GetChannelByID : GET
func GetChannelByID(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	qp := u.NewQueryParam(r)
	id := bone.GetValue(r, "id")
	channel, _, err := model.GetChannelByID(qp, id)
	if err != nil {
		res.SetError(JSONErrResourceNotFound)
		res.JSON(w, res, JSONErrResourceNotFound.Status)
		return
	}

	res.SetResponse(channel)
	res.JSON(w, res, http.StatusOK)
}

// UpdateChannel :
func UpdateChannel(w http.ResponseWriter, r *http.Request) {
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
	var reqChannel model.Channel
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&reqChannel); err != nil {
		res.SetErrorWithDetail(JSONErrFatal, err)
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}
	reqChannel.ID = id
	reqChannel.UpdatedBy = accData.AccountID
	err = reqChannel.Update()
	if err != nil {
		res.SetErrorWithDetail(JSONErrFatal, err)
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}
	res.SetResponse(model.Channels{reqChannel})
	res.JSON(w, res, http.StatusOK)
}

//DeleteChannel : remove channel
func DeleteChannel(w http.ResponseWriter, r *http.Request) {
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
	p := model.Channel{ID: id}
	p.UpdatedBy = accData.AccountID
	qp := u.NewQueryParam(r)
	datas, _, err := model.GetChannelByID(qp, id)
	if err != nil {
		res.SetErrorWithDetail(JSONErrFatal, err)
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}
	if len(*datas) <= 0 {
		res.SetError(JSONErrResourceNotFound)
		res.JSON(w, res, JSONErrResourceNotFound.Status)
		return
	}
	if (*datas)[0].IsSuper {
		res.SetError(JSONErrUnauthorized)
		res.JSON(w, res, JSONErrUnauthorized.Status)
		return
	}
	if err := p.Delete(); err != nil {
		res.SetError(JSONErrResourceNotFound)
		res.JSON(w, res, JSONErrResourceNotFound.Status)
		return
	}
	res.JSON(w, res, http.StatusOK)
}

//PostChannelTags : POST tags of channel
func PostChannelTags(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	var req model.ObjectTag
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		u.DEBUG(err)
		res.SetErrorWithDetail(JSONErrFatal, err)
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}
	// reqChannel.ID = bone.GetValue(r, "holder")
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

//GetChannelByTags : GET
func GetChannelByTags(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	qp := u.NewQueryParam(r)
	id := bone.GetValue(r, "tag_id")
	channel, _, err := model.GetChannelsByTags(qp, id)
	if err != nil {
		u.DEBUG(err)
		res.SetError(JSONErrResourceNotFound)
		res.JSON(w, res, JSONErrResourceNotFound.Status)
		return
	}

	res.SetResponse(channel)
	res.JSON(w, res, http.StatusOK)
}
