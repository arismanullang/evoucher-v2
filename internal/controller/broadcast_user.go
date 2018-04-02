package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	//"github.com/go-zoo/bone"
	"github.com/ruizu/render"

	"github.com/gilkor/evoucher/internal/model"
)

type (
	EmailUser struct {
		Id    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	ListEmailUser struct {
		Id         string   `json:"id"`
		Name       string   `json:"name"`
		EmailUsers []string `json:"email_users"`
	}
	RequestAddEmailUser struct {
		Name   string `json:"name"`
		Email  string `json:"email"`
		ListId string `json:"list_id"`
	}
	RequestEmailUser struct {
		EmailUserId string `json:"email_user_id"`
		ListId      string `json:"list_id"`
	}
	RequestEmailUsers struct {
		EmailUserId []string `json:"email_user_id"`
		ListId      string   `json:"list_id"`
	}
)

func InsertEmailUser(w http.ResponseWriter, r *http.Request) {
	apiName := "email_create"
	status := http.StatusCreated
	var rd EmailUser
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		log.Panic(err)
	}

	res := NewResponse("Success")
	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag(apiName)

	a := AuthTokenWithLogger(w, r, logger)
	if !a.Valid {
		res = a.res
		status = http.StatusUnauthorized
		res = NewResponse("")
		render.JSON(w, res, status)
		return
	}

	if CheckAPIRole(a, apiName) {
		logger.SetStatus(status).Info("param :", a.User.ID, "response :", "Invalid Role")

		status = http.StatusUnauthorized
		res = NewResponse("")
		res.AddError(its(status), model.ErrCodeInvalidRole, model.ErrInvalidRole.Error(), logger.TraceID)
		render.JSON(w, res, status)
		return
	}

	param := model.EmailUser{
		Name:      strings.ToLower(rd.Name),
		Email:     strings.ToLower(rd.Email),
		AccountID: a.User.Account.Id,
		CreatedBy: a.User.ID,
	}

	_, err := model.InsertEmailUser(param)
	if err != nil {
		status = http.StatusInternalServerError
		res = NewResponse("")
		res.AddError(its(status), model.ErrCodeInternalError, err.Error(), logger.TraceID)
		logger.SetStatus(status).Log("param :", param, "response :", err.Error())
	}

	render.JSON(w, res, status)
}

func GetAllEmailUser(w http.ResponseWriter, r *http.Request) {
	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag("select-email")

	res := NewResponse("")
	a := AuthTokenWithLogger(w, r, logger)
	if !a.Valid {
		res = a.res
		status := http.StatusUnauthorized
		render.JSON(w, res, status)
		return
	}

	user, err := model.GetAllEmailUser(a.User.Account.Id)
	if err != nil && err != model.ErrResourceNotFound {
		log.Panic(err)
	}

	res = NewResponse(user)
	render.JSON(w, res)
}

func GetListEmailUserByID(w http.ResponseWriter, r *http.Request) {
	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag("select-email")
	id := r.FormValue("id")

	res := NewResponse("")
	a := AuthTokenWithLogger(w, r, logger)
	if !a.Valid {
		res = a.res
		status := http.StatusUnauthorized
		render.JSON(w, res, status)
		return
	}

	user, err := model.GetListEmailUserById(id, a.User.Account.Id)
	if err != nil && err != model.ErrResourceNotFound {
		log.Panic(err)
	}

	res = NewResponse(user)
	render.JSON(w, res)
}

func GetListEmailUserByIDs(w http.ResponseWriter, r *http.Request) {
	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag("select-email")
	id := r.FormValue("id")

	res := NewResponse("")
	a := AuthTokenWithLogger(w, r, logger)
	if !a.Valid {
		res = a.res
		status := http.StatusUnauthorized
		render.JSON(w, res, status)
		return
	}

	idArr := strings.Split(id, "`")
	user, err := model.GetListEmailUserByIds(idArr, a.User.Account.Id)
	if err != nil && err != model.ErrResourceNotFound {
		log.Panic(err)
	}

	res = NewResponse(user)
	render.JSON(w, res)
}

func SearchEmailUser(w http.ResponseWriter, r *http.Request) {
	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag("select-email")

	res := NewResponse("")
	a := AuthTokenWithLogger(w, r, logger)
	if !a.Valid {
		res = a.res
		status := http.StatusUnauthorized
		render.JSON(w, res, status)
		return
	}

	param := getUrlParam(r.URL.String())

	param["account_id"] = a.User.Account.Id

	user, err := model.GetEmailUser(param)
	if err != nil && err != model.ErrResourceNotFound {
		log.Panic(err)
	}

	res = NewResponse(user)
	render.JSON(w, res)
}

func GetEmailUserByIDs(w http.ResponseWriter, r *http.Request) {
	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag("select-email")

	res := NewResponse("")
	a := AuthTokenWithLogger(w, r, logger)
	if !a.Valid {
		res = a.res
		status := http.StatusUnauthorized
		render.JSON(w, res, status)
		return
	}

	id := r.FormValue("id")
	idArr := strings.Split(id, "`")

	user, err := model.GetEmailUserByIDs(idArr)
	if err != nil && err != model.ErrResourceNotFound {
		log.Panic(err)
	}

	res = NewResponse(user)
	render.JSON(w, res)
}

func DeleteEmailUser(w http.ResponseWriter, r *http.Request) {
	apiName := "email_delete"
	status := http.StatusCreated
	var rd ChangeUserStatusReq
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		log.Panic(err)
	}

	res := NewResponse("Success")
	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag(apiName)

	a := AuthTokenWithLogger(w, r, logger)
	if !a.Valid {
		res = a.res
		status = http.StatusUnauthorized
		res = NewResponse("")
		render.JSON(w, res, status)
		return
	}

	if CheckAPIRole(a, apiName) {
		logger.SetStatus(status).Info("param :", a.User.ID, "response :", "Invalid Role")

		status = http.StatusUnauthorized
		res = NewResponse("")
		res.AddError(its(status), model.ErrCodeInvalidRole, model.ErrInvalidRole.Error(), logger.TraceID)
		render.JSON(w, res, status)
		return
	}

	err := model.DeleteEmailUser(rd.Id, a.User.ID)
	if err != nil {
		status = http.StatusInternalServerError
		res = NewResponse("")
		res.AddError(its(status), model.ErrCodeInternalError, err.Error(), logger.TraceID)
		logger.SetStatus(status).Log("param :", rd.Id, "response :", err.Error())
	}

	render.JSON(w, res, status)
}

// List email
func InsertListEmailUser(w http.ResponseWriter, r *http.Request) {
	apiName := "email_create"
	status := http.StatusCreated
	var rd ListEmailUser
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		log.Panic(err)
	}

	res := NewResponse("Success")
	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag(apiName)

	a := AuthTokenWithLogger(w, r, logger)
	if !a.Valid {
		res = a.res
		status = http.StatusUnauthorized
		res = NewResponse("")
		render.JSON(w, res, status)
		return
	}

	if CheckAPIRole(a, apiName) {
		logger.SetStatus(status).Info("param :", a.User.ID, "response :", "Invalid Role")

		status = http.StatusUnauthorized
		res = NewResponse("")
		res.AddError(its(status), model.ErrCodeInvalidRole, model.ErrInvalidRole.Error(), logger.TraceID)
		render.JSON(w, res, status)
		return
	}

	emailUser := []model.EmailUser{}
	for _, v := range rd.EmailUsers {
		tempEmailUser := model.EmailUser{
			ID: v,
		}
		emailUser = append(emailUser, tempEmailUser)
	}

	param := model.ListEmailUser{
		Name:       rd.Name,
		EmailUsers: emailUser,
		AccountID:  a.User.Account.Id,
		CreatedBy:  a.User.ID,
	}

	err := model.InsertListEmail(param)
	if err != nil {
		status = http.StatusInternalServerError
		res = NewResponse("")
		res.AddError(its(status), model.ErrCodeInternalError, err.Error(), logger.TraceID)
		logger.SetStatus(status).Log("param :", param, "response :", err.Error())
	}

	render.JSON(w, res, status)
}

func AddEmailUser(w http.ResponseWriter, r *http.Request) {
	apiName := "email_create"
	status := http.StatusCreated
	var rd RequestEmailUsers
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		log.Panic(err)
	}

	res := NewResponse("Success")
	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag(apiName)

	a := AuthTokenWithLogger(w, r, logger)
	if !a.Valid {
		res = a.res
		status = http.StatusUnauthorized
		res = NewResponse("")
		render.JSON(w, res, status)
		return
	}

	if CheckAPIRole(a, apiName) {
		logger.SetStatus(status).Info("param :", a.User.ID, "response :", "Invalid Role")

		status = http.StatusUnauthorized
		res = NewResponse("")
		res.AddError(its(status), model.ErrCodeInvalidRole, model.ErrInvalidRole.Error(), logger.TraceID)
		render.JSON(w, res, status)
		return
	}

	err := model.AddEmailUserToList(rd.EmailUserId, rd.ListId, a.User.Account.Id, a.User.ID)
	if err != nil {
		status = http.StatusInternalServerError
		res = NewResponse("")
		res.AddError(its(status), model.ErrCodeInternalError, err.Error(), logger.TraceID)
		logger.SetStatus(status).Log("param :", rd, "response :", err.Error())
	}

	render.JSON(w, res, status)
}

func AddNewEmailUser(w http.ResponseWriter, r *http.Request) {
	apiName := "email_create"
	status := http.StatusCreated
	var rd RequestAddEmailUser
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		log.Panic(err)
	}

	res := NewResponse("Success")
	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag(apiName)

	a := AuthTokenWithLogger(w, r, logger)
	if !a.Valid {
		res = a.res
		status = http.StatusUnauthorized
		res = NewResponse("")
		render.JSON(w, res, status)
		return
	}

	if CheckAPIRole(a, apiName) {
		logger.SetStatus(status).Info("param :", a.User.ID, "response :", "Invalid Role")

		status = http.StatusUnauthorized
		res = NewResponse("")
		res.AddError(its(status), model.ErrCodeInvalidRole, model.ErrInvalidRole.Error(), logger.TraceID)
		render.JSON(w, res, status)
		return
	}

	param := model.EmailUser{
		Name:      strings.ToLower(rd.Name),
		Email:     strings.ToLower(rd.Email),
		AccountID: a.User.Account.Id,
		CreatedBy: a.User.ID,
	}

	err := model.AddNewEmailUserToList(param, rd.ListId)
	if err != nil {
		status = http.StatusInternalServerError
		res = NewResponse("")
		res.AddError(its(status), model.ErrCodeInternalError, err.Error(), logger.TraceID)
		logger.SetStatus(status).Log("param :", param, "response :", err.Error())
	}

	render.JSON(w, res, status)
}

func RemoveEmailUser(w http.ResponseWriter, r *http.Request) {
	apiName := "email_delete"
	status := http.StatusOK
	var rd RequestEmailUser
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		log.Panic(err)
	}

	res := NewResponse("Success")
	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag(apiName)

	a := AuthTokenWithLogger(w, r, logger)
	if !a.Valid {
		res = a.res
		status = http.StatusUnauthorized
		res = NewResponse("")
		render.JSON(w, res, status)
		return
	}

	if CheckAPIRole(a, apiName) {
		logger.SetStatus(status).Info("param :", a.User.ID, "response :", "Invalid Role")

		status = http.StatusUnauthorized
		res = NewResponse("")
		res.AddError(its(status), model.ErrCodeInvalidRole, model.ErrInvalidRole.Error(), logger.TraceID)
		render.JSON(w, res, status)
		return
	}

	err := model.RemoveEmailUserFromList(rd.EmailUserId, rd.ListId, a.User.ID)
	if err != nil {
		status = http.StatusInternalServerError
		res = NewResponse("")
		res.AddError(its(status), model.ErrCodeInternalError, err.Error(), logger.TraceID)
		logger.SetStatus(status).Log("param :", rd, "response :", err.Error())
	}

	render.JSON(w, res, status)
}

func GetAllListEmailUser(w http.ResponseWriter, r *http.Request) {
	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag("select-email")

	res := NewResponse("")
	a := AuthTokenWithLogger(w, r, logger)
	if !a.Valid {
		res = a.res
		status := http.StatusUnauthorized
		render.JSON(w, res, status)
		return
	}

	user, err := model.GetAllListEmailUser(a.User.Account.Id)
	if err != nil && err != model.ErrResourceNotFound {
		log.Panic(err)
	}

	res = NewResponse(user)
	render.JSON(w, res)
}

func DeleteListUser(w http.ResponseWriter, r *http.Request) {
	apiName := "email_delete"
	status := http.StatusCreated
	var rd ChangeUserStatusReq
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		log.Panic(err)
	}

	res := NewResponse("Success")
	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag(apiName)

	a := AuthTokenWithLogger(w, r, logger)
	if !a.Valid {
		res = a.res
		status = http.StatusUnauthorized
		res = NewResponse("")
		render.JSON(w, res, status)
		return
	}

	if CheckAPIRole(a, apiName) {
		logger.SetStatus(status).Info("param :", a.User.ID, "response :", "Invalid Role")

		status = http.StatusUnauthorized
		res = NewResponse("")
		res.AddError(its(status), model.ErrCodeInvalidRole, model.ErrInvalidRole.Error(), logger.TraceID)
		render.JSON(w, res, status)
		return
	}

	err := model.DeleteListUser(rd.Id, a.User.ID)
	if err != nil {
		status = http.StatusInternalServerError
		res = NewResponse("")
		res.AddError(its(status), model.ErrCodeInternalError, err.Error(), logger.TraceID)
		logger.SetStatus(status).Log("param :", rd.Id, "response :", err.Error())
	}

	render.JSON(w, res, status)
}
