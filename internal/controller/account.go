package controller

import (
	"encoding/json"
	//"fmt"
	"log"
	"net/http"

	//"github.com/go-zoo/bone"
	"github.com/ruizu/render"

	"github.com/gilkor/evoucher/internal/model"
)

type (
	Account struct {
		AccountName string `json:"account_name"`
		Billing     string `json:"billing"`
		CreatedBy   string `json:"created_by"`
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

func GetAccountId(w http.ResponseWriter, r *http.Request) {
	account, err := model.FindAllAccounts()
	if err != nil && err != model.ErrResourceNotFound {
		log.Panic(err)
	}

	res := NewResponse(account)
	render.JSON(w, res)
}

func GetAllAccountRoles(w http.ResponseWriter, r *http.Request) {
	role, err := model.FindAllRole()
	if err != nil && err != model.ErrResourceNotFound {
		log.Panic(err)
	}

	res := NewResponse(role)
	render.JSON(w, res)
}
