package controller

import (
	"encoding/json"
	"fmt"
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
	var rd User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		log.Panic(err)
	}

	fmt.Println(len(hash(rd.Password)))

	param := model.User{
		AccountId: rd.AccountId,
		Username:  rd.Username,
		Password:  hash(rd.Password),
		Email:     rd.Email,
		Phone:     rd.Phone,
		RoleId:    rd.RoleId,
		CreatedBy: rd.CreatedBy,
	}

	if err := model.AddUser(param); err != nil {
		log.Panic(err)
	}

	res := NewResponse(nil)
	render.JSON(w, res, http.StatusCreated)
}

func GetAccountId(w http.ResponseWriter, r *http.Request) {
	account, err := model.FindAllAccount()
	if err != nil && err != model.ErrResourceNotFound {
		log.Panic(err)
	}

	res := NewResponse(account)
	render.JSON(w, res)
}
