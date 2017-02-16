package controller

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	//"github.com/go-zoo/bone"
	"github.com/ruizu/render"

	"github.com/gilkor/evoucher/internal/model"
)

type (
	User struct {
		AccountId string   `json:"account_id"`
		Username  string   `json:"username"`
		Password  string   `json:"password"`
		Email     string   `json:"email"`
		Phone     string   `json:"phone"`
		RoleId    []string `json:"role_id"`
		CreatedBy string   `json:"created_by"`
	}

	UserLogin struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	RoleReq struct {
		AccountId string `json:"account_id"`
		Role      string `json:"role"`
	}

	AccountReq struct {
		AccountId string `json:"account_id"`
	}
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {
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

func FindUserByRole(w http.ResponseWriter, r *http.Request) {
	var rd RoleReq
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		log.Panic(err)
	}

	var user = model.UserResponse{}
	var err error
	if basicAuth(w, r) {
		user, err = model.FindUserByRole(rd.Role, rd.AccountId)
		if err != nil && err != model.ErrResourceNotFound {
			log.Panic(err)
		}
	} else {
		user = model.UserResponse{}
	}

	res := NewResponse(user)
	render.JSON(w, res, http.StatusCreated)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	//param := getUrlParam(r.URL.String())
	var rd AccountReq
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		log.Panic(err)
	}

	var user = model.UserResponse{}
	var err error
	if basicAuth(w, r) {
		user, err = model.FindAllUser(rd.AccountId)
		if err != nil && err != model.ErrResourceNotFound {
			log.Panic(err)
		}
	} else {
		user = model.UserResponse{}
	}

	res := NewResponse(user)
	render.JSON(w, res)
}

func GetUserCustomParam(w http.ResponseWriter, r *http.Request) {
	param := getUrlParam(r.URL.String())

	var user = model.UserResponse{}
	var err error
	if basicAuth(w, r) {
		user, err = model.FindUser(param)
		if err != nil && err != model.ErrResourceNotFound {
			log.Panic(err)
		}
	} else {
		user = model.UserResponse{}
	}

	res := NewResponse(user)
	render.JSON(w, res)
}

func getUrlParam(url string) map[string]string {
	s := strings.Split(url, "?")
	param := strings.Split(s[1], "&")

	m := make(map[string]string)

	for _, v := range param {
		tempStr := strings.Split(v, "=")
		m[tempStr[0]] = tempStr[1]
	}

	return m
}

func sendPost(url string, param string) []byte {
	fmt.Println(param)
	var jsonStr = []byte(param)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
	return body
}

func getResponseData(param []byte) map[string]interface{} {
	var dat map[string]interface{}
	dat = make(map[string]interface{})
	if err := json.Unmarshal(param, &dat); err != nil {
		panic(err)
	}
	//fmt.Println(string(robots))

	if str, ok := dat["data"].(map[string]interface{}); ok {
		return str
	} else {
		return nil
	}
}

func hash(param string) string {
	password := []byte(param)
	hash := sha256.Sum256(password)
	return base64.StdEncoding.EncodeToString(hash[:])
}

func basicAuth(w http.ResponseWriter, r *http.Request) bool {
	s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	if len(s) != 2 {
		return false
	}

	b, err := base64.StdEncoding.DecodeString(s[1])
	if err != nil {
		return false
	}

	pair := strings.SplitN(string(b), ":", 2)
	if len(pair) != 2 {
		return false
	}

	login, err := model.Login(pair[0], hash(pair[1]))

	if login == 0 || err != nil {
		return false
	}

	return true
}
