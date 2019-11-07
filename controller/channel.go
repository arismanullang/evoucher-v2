package controller

import (
	"encoding/json"
	"net/http"

	"github.com/gilkor/evoucher-v2/model"
	u "github.com/gilkor/evoucher-v2/util"
	"github.com/go-zoo/bone"
)

//PostChannel : POST channel data
func PostChannel(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	var reqChannel model.Channel
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&reqChannel); err != nil {
		u.DEBUG(err)
		res.SetError(JSONErrFatal)
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}
	if err := reqChannel.Insert(); err != nil {
		u.DEBUG(err)
		res.SetError(JSONErrFatal)
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	res.JSON(w, res, http.StatusCreated)
}

//GetChannels : GET list of channels
func GetChannels(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	qp := u.NewQueryParam(r)

	channels, next, err := model.GetChannels(qp)
	if err != nil {
		res.SetError(JSONErrFatal.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	res.SetResponse(channels)
	res.SetPagination(r, qp.Page, next)
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

	id := bone.GetValue(r, "id")
	var reqChannel model.Channel
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&reqChannel); err != nil {
		res.SetError(JSONErrFatal)
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}
	reqChannel.ID = id
	if err := reqChannel.Update(); err != nil {
		res.SetError(JSONErrFatal)
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}
	res.JSON(w, res, http.StatusOK)
}

//DeleteChannel : remove channel
func DeleteChannel(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	id := bone.GetValue(r, "id")
	p := model.Channel{ID: id}
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

	var req model.TagHolder
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		u.DEBUG(err)
		res.SetError(JSONErrFatal)
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}
	// reqChannel.ID = bone.GetValue(r, "holder")
	if err := req.Insert(); err != nil {
		u.DEBUG(err)
		res.SetError(JSONErrFatal)
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

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
