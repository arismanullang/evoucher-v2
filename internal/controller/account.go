package controller

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	//"github.com/go-zoo/bone"
	"github.com/ruizu/render"

	"github.com/gilkor/evoucher/internal/model"
)

type (
	AccountId struct {
		Id string `json:"id"`
	}
	Account struct {
		Id        string         `json:"id"`
		Name      string         `json:"name"`
		Alias     string         `json:"alias"`
		Email     string         `json:"email"`
		Billing   sql.NullString `json:"billing"`
		Address   string         `json:"address"`
		City      string         `json:"city"`
		Province  string         `json:"province"`
		CreatedBy string         `json:"created_by"`
	}
)

func RegisterAccount(w http.ResponseWriter, r *http.Request) {
	apiName := "sa_a-create"
	status := http.StatusCreated
	var rd Account
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

	param := model.Account{
		Name:    rd.Name,
		Billing: rd.Billing,
		Email:   rd.Email,
		Alias:   rd.Alias,
	}

	err := model.AddAccount(param, a.User.ID)
	if err != nil {
		status = http.StatusInternalServerError
		res = NewResponse("")
		res.AddError(its(status), model.ErrCodeInternalError, err.Error(), logger.TraceID)
		logger.SetStatus(status).Log("param :", param, "response :", err.Error())

	}

	render.JSON(w, res, status)
}

func UpdateAccount(w http.ResponseWriter, r *http.Request) {
	apiName := "sa_a-update"
	status := http.StatusCreated
	var rd Account
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

	param := model.Account{
		Id:      rd.Id,
		Name:    rd.Name,
		Billing: rd.Billing,
		Email:   rd.Email,
		Alias:   rd.Alias,
	}

	err := model.UpdateAccount(param, a.User.ID)
	if err != nil {
		status = http.StatusInternalServerError
		res = NewResponse("")
		res.AddError(its(status), model.ErrCodeInternalError, err.Error(), logger.TraceID)
		logger.SetStatus(status).Log("param :", param, "response :", err.Error())

	}

	render.JSON(w, res, status)
}

func GetAllAccountsDetail(w http.ResponseWriter, r *http.Request) {
	account, err := model.FindAllAccountsDetail()
	if err != nil && err != model.ErrResourceNotFound {
		log.Panic(err)
	}

	res := NewResponse(account)
	render.JSON(w, res)
}

func GetAllAccount(w http.ResponseWriter, r *http.Request) {
	account, err := model.FindAllAccounts()
	if err != nil && err != model.ErrResourceNotFound {
		log.Panic(err)
	}

	res := NewResponse(account)
	render.JSON(w, res)
}

func GetAccountDetailByUser(w http.ResponseWriter, r *http.Request) {
	status := http.StatusOK
	res := NewResponse(nil)

	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag("Get Account By User")

	a := AuthTokenWithLogger(w, r, logger)
	if !a.Valid {
		res = a.res
		status = http.StatusUnauthorized
		render.JSON(w, res, status)
		return
	}

	account, err := model.GetAccountDetailByUser(a.User.ID)
	if err != nil {
		status = http.StatusInternalServerError
		errTitle := model.ErrCodeInternalError
		if err != model.ErrResourceNotFound {
			status = http.StatusNotFound
			errTitle = model.ErrCodeResourceNotFound
		}

		res.AddError(its(status), errTitle, err.Error(), logger.TraceID)
		logger.SetStatus(status).Log("param :", a.User.ID, "response :", err.Error())
	} else {
		res = NewResponse(account)
	}

	render.JSON(w, res, status)
}

func GetAccountDetailByOtherUser(w http.ResponseWriter, r *http.Request) {
	status := http.StatusOK
	res := NewResponse(nil)
	id := r.FormValue("id")

	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag("Get Account By User")

	a := AuthTokenWithLogger(w, r, logger)
	if !a.Valid {
		res = a.res
		status = http.StatusUnauthorized
		render.JSON(w, res, status)
		return
	}

	account, err := model.GetAccountDetailByAccountId(id)
	if err != nil {
		status = http.StatusInternalServerError
		errTitle := model.ErrCodeInternalError
		if err != model.ErrResourceNotFound {
			status = http.StatusNotFound
			errTitle = model.ErrCodeResourceNotFound
		}

		res.AddError(its(status), errTitle, err.Error(), logger.TraceID)
		logger.SetStatus(status).Log("param :", a.User.ID, "response :", err.Error())
	} else {
		res = NewResponse(account)
	}

	render.JSON(w, res, status)
}

func BlockAccount(w http.ResponseWriter, r *http.Request) {
	apiName := "sa_a-delete"
	status := http.StatusCreated
	var rd AccountId
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

	err := model.BlockAccount(rd.Id, a.User.ID)
	if err != nil {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, err.Error(), logger.TraceID)
		logger.SetStatus(status).Log("param :", rd.Id, "response :", err.Error())

	}

	res = NewResponse("Success")
	render.JSON(w, res, status)
}

func ActivateAccount(w http.ResponseWriter, r *http.Request) {
	apiName := "sa_a-delete"
	status := http.StatusCreated
	var rd AccountId
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

	err := model.ActivateAccount(rd.Id, a.User.ID)
	if err != nil {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, err.Error(), logger.TraceID)
		logger.SetStatus(status).Log("param :", rd.Id, "response :", err.Error())

	}

	res = NewResponse("Success")
	render.JSON(w, res, status)
}
