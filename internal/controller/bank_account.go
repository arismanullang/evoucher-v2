package controller

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	//"github.com/go-zoo/bone"
	"github.com/ruizu/render"

	"github.com/gilkor/evoucher/internal/model"
)

type (
	BankAccount struct {
		Id                string         `db:"id" json:"id"`
		CompanyName       string         `db:"company_name" json:"company_name"`
		CompanyPic        string         `db:"company_pic" json:"company_pic"`
		CompanyTelp       string         `db:"company_telp" json:"company_telp"`
		CompanyEmail      string         `db:"company_email" json:"company_email"`
		BankName          string         `db:"bank_name" json:"bank_name"`
		BankBranch        string         `db:"bank_branch" json:"bank_branch"`
		BankAccountNumber string         `db:"bank_account_number" json:"bank_account_number"`
		BankAccountHolder string         `db:"bank_account_holder" json:"bank_account_holder"`
		AccountId         string         `db:"account_id" json:"account_id"`
		CreatedAt         time.Time      `db:"created_at" json:"created_at"`
		CreatedBy         string         `db:"created_by" json:"created_by"`
		UpdatedAt         sql.NullString `db:"updated_at" json:"updated_at"`
		UpdatedBy         sql.NullString `db:"updated_by" json:"updated_by"`
		Status            string         `db:"status" json:"status"`
	}
)

func RegisterBankAccount(w http.ResponseWriter, r *http.Request) {
	apiName := "partner_create"
	status := http.StatusCreated
	var rd BankAccount
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

	param := model.BankAccount{
		CompanyName:       rd.CompanyName,
		CompanyPic:        rd.CompanyPic,
		CompanyTelp:       rd.CompanyTelp,
		CompanyEmail:      rd.CompanyEmail,
		BankName:          rd.BankName,
		BankBranch:        rd.BankBranch,
		BankAccountHolder: rd.BankAccountHolder,
		BankAccountNumber: rd.BankAccountNumber,
	}

	err := model.AddBankAccount(param, a.User)
	if err != nil {
		status = http.StatusInternalServerError
		res = NewResponse("")
		res.AddError(its(status), model.ErrCodeInternalError, err.Error(), logger.TraceID)
		logger.SetStatus(status).Log("param :", param, "response :", err.Error())

	}

	render.JSON(w, res, status)
}

func UpdateBankAccount(w http.ResponseWriter, r *http.Request) {
	apiName := "partner_update"
	status := http.StatusCreated
	var rd BankAccount
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

	param := model.BankAccount{
		CompanyName:       rd.CompanyName,
		CompanyPic:        rd.CompanyPic,
		CompanyTelp:       rd.CompanyTelp,
		CompanyEmail:      rd.CompanyEmail,
		BankName:          rd.BankName,
		BankBranch:        rd.BankBranch,
		BankAccountHolder: rd.BankAccountHolder,
		BankAccountNumber: rd.BankAccountNumber,
	}

	err := model.UpdateBankAccount(param, a.User.ID)
	if err != nil {
		status = http.StatusInternalServerError
		res = NewResponse("")
		res.AddError(its(status), model.ErrCodeInternalError, err.Error(), logger.TraceID)
		logger.SetStatus(status).Log("param :", param, "response :", err.Error())

	}

	render.JSON(w, res, status)
}

func GetAllBankAccounts(w http.ResponseWriter, r *http.Request) {
	apiName := "partner_select"
	status := http.StatusCreated

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

	account, err := model.FindAllBankAccounts(a.User.Account.Id)
	if err != nil && err != model.ErrResourceNotFound {
		log.Panic(err)
	}

	res = NewResponse(account)
	render.JSON(w, res)
}

func GetBankAccountDetailByBankAccountNumber(w http.ResponseWriter, r *http.Request) {
	status := http.StatusOK
	res := NewResponse(nil)
	accountNumber := r.FormValue("account_number")

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

	account, err := model.FindBankAccount(a.User.Account.Id, accountNumber)
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

func GetBankAccountDetailByPartner(w http.ResponseWriter, r *http.Request) {
	status := http.StatusOK
	res := NewResponse(nil)
	partner := r.FormValue("partner")

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

	account, err := model.FindBankAccountByPartner(a.User.Account.Id, partner)
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
