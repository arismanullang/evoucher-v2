package controller

import (
	"encoding/json"
	"log"
	"net/http"
	//"time"
	"io/ioutil"

	"github.com/go-zoo/bone"
	"github.com/ruizu/render"

	"github.com/evoucher/voucher/internal/model"
)

type (
	User struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	RoleReq struct {
		Role string `json:"role"`
	}
)

func GetUserByRole(w http.ResponseWriter, r *http.Request) {
	var rd RoleReq
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		log.Panic(err)
	}

	user, err := model.FindAccountByRole(rd.Role)
	if err != nil && err != model.ErrResourceNotFound {
		log.Panic(err)
	}

	res := NewResponse(user)
	render.JSON(w, res)
}

func GetToken(w http.ResponseWriter, r *http.Request) {
	var rd User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		log.Panic(err)
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", "juno-staging.elys.id", nil)
	req.SetBasicAuth(rd.Username, rd.Password)
	resp, err := client.Do(req)
	if err != nil {
		log.Panic(err)
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	s := string(bodyText)

	res := NewResponse(s)
	render.JSON(w, res, http.StatusCreated)
}

func userpassing(w http.ResponseWriter, r *http.Request) {
	id := bone.GetValue(r, "id")
	plaintext := []byte(id)

	ec := model.Encrypt(plaintext)
	dc := model.Decrypt(ec)

	res := NewResponse(dc)
	render.JSON(w, res, http.StatusCreated)
}
