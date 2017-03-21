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
	user := ""
	token := r.FormValue("token")
	status := http.StatusUnauthorized
	err := model.ErrTokenNotFound
	res := NewResponse(nil)

	res.AddError(its(status), its(status), err.Error(), "account")

	valid := false
	if token != "" && token != "null" {
		user, _, _, valid, _ = getValiditySession(r, token)
	}
	fmt.Println("user " + user)
	if valid {
		status = http.StatusOK
		account, err := model.GetAccountDetailByUser(user)
		if err != nil {
			status = http.StatusInternalServerError
			if err != model.ErrResourceNotFound {
				status = http.StatusNotFound
			}

			res.AddError(its(status), its(status), err.Error(), "account")
		} else {
			res = NewResponse(account)
		}
	}
	render.JSON(w, res, status)
}

func GetAccountsByUser(w http.ResponseWriter, r *http.Request) {
	user := ""
	token := r.FormValue("token")
	status := http.StatusUnauthorized
	err := model.ErrTokenNotFound
	res := NewResponse(nil)

	res.AddError(its(status), its(status), err.Error(), "account")

	valid := false
	if token != "" && token != "null" {
		user, _, _, valid, _ = getValiditySession(r, token)
	}
	if valid {
		status = http.StatusOK
		account, err := model.GetAccountsByUser(user)
		if err != nil {
			status = http.StatusInternalServerError
			if err != model.ErrResourceNotFound {
				status = http.StatusNotFound
			}

			res.AddError(its(status), its(status), err.Error(), "account")
		} else {
			res = NewResponse(account)
		}
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
