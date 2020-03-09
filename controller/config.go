package controller

import (
	"encoding/json"
	"net/http"

	"github.com/gilkor/evoucher-v2/model"
	u "github.com/gilkor/evoucher-v2/util"
	"github.com/go-zoo/bone"
)

// PostConfig : Create new config
func PostConfig(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	var configs model.Configs
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&configs); err != nil {
		u.DEBUG(err)
		res.SetError(JSONErrFatal)
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	// insert config
	response, err := configs.Upsert()
	if err != nil {
		res.SetError(JSONErrFatal.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	res.SetResponse(response)
	res.JSON(w, res, http.StatusOK)
}

// UpdateConfig :
func UpdateConfig(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	// companyID := bone.GetValue(r, "company_id")
	// key := bone.GetValue(r, "key")

	var config model.Config
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&config); err != nil {
		res.SetErrorWithDetail(JSONErrFatal, err)
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	err := config.Update()
	if err != nil {
		res.SetErrorWithDetail(JSONErrFatal, err)
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}
	res.SetResponse(model.Configs{config})
	res.JSON(w, res, http.StatusOK)
}

//GetConfigs : GET list of configs
func GetConfigs(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	companyID := bone.GetValue(r, "company")
	category := r.FormValue("category")

	configs, err := model.GetConfigs(companyID, category)
	if err != nil {
		res.SetError(JSONErrFatal.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	res.SetResponse(configs)
	// res.SetPagination(r, qp.Page, next)
	res.JSON(w, res, http.StatusOK)
}
