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

	valid := false
	for _, valueRole := range a.User.Role {
		features := model.ApiFeatures[valueRole.Detail]
		for _, valueFeature := range features {
			if apiName == valueFeature {
				valid = true
			}
		}
	}

	if !valid {
		logger.SetStatus(status).Info("param :", a.User.ID, "response :", "Invalid Role")

		status = http.StatusUnauthorized
		res.AddError(its(status), model.ErrCodeInvalidRole, model.ErrInvalidRole.Error(), logger.TraceID)
		render.JSON(w, res, status)
		return

	}

	param := model.Account{
		Name:    rd.Name,
		Billing: rd.Billing,
		Alias:   rd.Alias,
	}

	err := model.AddAccount(param, a.User.ID)
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
	features, err := model.GetAllFeatures()
	if err != nil && err != model.ErrResourceNotFound {
		log.Panic(err)
	}

	res := NewResponse(features)
	render.JSON(w, res)
}
