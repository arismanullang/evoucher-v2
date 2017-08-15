package controller

import (
	//"fmt"
	"log"
	"net/http"

	//"github.com/go-zoo/bone"
	"github.com/ruizu/render"

	"encoding/json"
	"github.com/gilkor/evoucher/internal/model"
)

type (
	Role struct {
		Id       string   `json:"id"`
		Detail   string   `json:"detail"`
		Features []string `json:"features"`
	}
)

func AddRole(w http.ResponseWriter, r *http.Request) {
	apiName := "role_create"
	status := http.StatusCreated
	var rd Role
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		log.Panic(err)
	}

	res := NewResponse("")
	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag(apiName)

	a := AuthTokenWithLogger(w, r, logger)
	if !a.Valid {
		res = a.res
		status = http.StatusUnauthorized
		render.JSON(w, res, status)
		return
	}

	if CheckAPIRole(a, apiName) {
		logger.SetStatus(status).Info("param :", a.User.ID, "response :", "Invalid Role")

		status = http.StatusUnauthorized
		res.AddError(its(status), model.ErrCodeInvalidRole, model.ErrInvalidRole.Error(), logger.TraceID)
		render.JSON(w, res, status)
		return

	}

	param := model.Role{
		Detail:   rd.Detail,
		Features: rd.Features,
	}

	err := model.AddRole(param, a.User.ID)
	if err != nil {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, err.Error(), logger.TraceID)
		logger.SetStatus(status).Log("param :", param, "response :", err.Error())

	}

	res = NewResponse("Success")
	render.JSON(w, res, status)
}

func UpdateRole(w http.ResponseWriter, r *http.Request) {
	apiName := "role_update"
	status := http.StatusCreated
	var rd Role
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		log.Panic(err)
	}

	res := NewResponse("")
	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag(apiName)

	a := AuthTokenWithLogger(w, r, logger)
	if !a.Valid {
		res = a.res
		status = http.StatusUnauthorized
		render.JSON(w, res, status)
		return
	}

	if CheckAPIRole(a, apiName) {
		logger.SetStatus(status).Info("param :", a.User.ID, "response :", "Invalid Role")

		status = http.StatusUnauthorized
		res.AddError(its(status), model.ErrCodeInvalidRole, model.ErrInvalidRole.Error(), logger.TraceID)
		render.JSON(w, res, status)
		return

	}

	param := model.Role{
		Id:       rd.Id,
		Features: rd.Features,
	}

	err := model.UpdateRole(param, a.User.ID)
	if err != nil {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, err.Error(), logger.TraceID)
		logger.SetStatus(status).Log("param :", param, "response :", err.Error())

	}

	res = NewResponse("Success")
	render.JSON(w, res, status)
}

func GetAllAccountRoles(w http.ResponseWriter, r *http.Request) {
	role, err := model.FindAllRole()
	if err != nil && err != model.ErrResourceNotFound {
		log.Panic(err)
	}

	res := NewResponse(role)
	render.JSON(w, res)
}

func GetAllFeatures(w http.ResponseWriter, r *http.Request) {
	res := NewResponse("")
	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag("Get Feature Detail")

	a := AuthTokenWithLogger(w, r, logger)
	if !a.Valid {
		res = a.res
		render.JSON(w, res, http.StatusUnauthorized)
		return
	}

	features, err := model.GetAllFeatures()
	if err != nil && err != model.ErrResourceNotFound {
		log.Panic(err)
	}

	res = NewResponse(features)
	render.JSON(w, res)
}

func GetFeaturesDetail(w http.ResponseWriter, r *http.Request) {
	res := NewResponse("")
	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag("Get Feature Detail")

	a := AuthTokenWithLogger(w, r, logger)
	if !a.Valid {
		res = a.res
		render.JSON(w, res, http.StatusUnauthorized)
		return
	}

	id := r.FormValue("id")

	roles, err := model.GetRoleDetail(id)
	if err != nil && err != model.ErrResourceNotFound {
		logger.SetStatus(http.StatusInternalServerError).Log("param :", id, "response :", err.Error())
	}

	res = NewResponse(roles)
	render.JSON(w, res)
}
