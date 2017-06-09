package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	//"fmt"
	"log"
	"net/http"

	//"github.com/go-zoo/bone"
	"github.com/ruizu/render"

	"github.com/gilkor/evoucher/internal/model"
)

type (
	Account struct {
		AccountName string         `json:"account_name"`
		Billing     sql.NullString `json:"billing"`
		CreatedBy   string         `json:"created_by"`
	}
)

func RegisterAccount(w http.ResponseWriter, r *http.Request) {
	status := http.StatusCreated
	var rd Account
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		log.Panic(err)
	}

	param := model.Account{
		AccountName: rd.AccountName,
		Billing:     rd.Billing,
		CreatedBy:   rd.CreatedBy,
	}

	if err := model.AddAccount(param); err != nil {
		//log.Panic(err)
		status = http.StatusInternalServerError
	}

	res := NewResponse(nil)
	render.JSON(w, res, status)
}

func GetAccount(w http.ResponseWriter, r *http.Request) {
	account, err := model.FindAllAccounts()
	if err != nil && err != model.ErrResourceNotFound {
		log.Panic(err)
	}

	res := NewResponse(account)
	render.JSON(w, res)
}

func GetAccountDetailByUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get Account Details")
	status := http.StatusUnauthorized
	err := model.ErrTokenNotFound
	errTitle := model.ErrCodeInvalidToken
	res := NewResponse(nil)

	res.AddError(its(status), errTitle, err.Error(), "Get Account")

	a := AuthToken(w, r)
	if a.Valid{
		status = http.StatusOK
		account, err := model.GetAccountDetailByUser(a.User.ID)
		if err != nil {
			status = http.StatusInternalServerError
			errTitle = model.ErrCodeInternalError
			if err != model.ErrResourceNotFound {
				status = http.StatusNotFound
				errTitle = model.ErrCodeResourceNotFound
			}

			res.AddError(its(status), errTitle, err.Error(), "Get Acccount")
		} else {
			res = NewResponse(account)
		}
	}else {
		res = a.res
		status = http.StatusUnauthorized
	}
	render.JSON(w, res, status)
}

func GetAccountsByUser(w http.ResponseWriter, r *http.Request) {
	status := http.StatusUnauthorized
	err := model.ErrTokenNotFound
	errTitle := model.ErrCodeInvalidToken
	res := NewResponse(nil)

	res.AddError(its(status), errTitle, err.Error(), "Get ccount")

	a := AuthToken(w, r)
	if a.Valid {
		status = http.StatusOK
		account, err := model.GetAccountsByUser(a.User.ID)
		if err != nil {
			status = http.StatusInternalServerError
			errTitle = model.ErrCodeInternalError
			if err != model.ErrResourceNotFound {
				status = http.StatusNotFound
				errTitle = model.ErrCodeResourceNotFound
			}

			res.AddError(its(status), errTitle, err.Error(), "Get ccount")
		} else {
			res = NewResponse(account)
		}
	}else {
		res = a.res
		status = http.StatusUnauthorized
	}
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
