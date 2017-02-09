package controller

import (
	"bytes"
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
	UserReq struct {
		AccountRole string `json:"account_role"`
		AssignBy    string `json:"assign_by"`
		CreatedBy   string `json:"created_by"`
		UserValue   User   `json:"user"`
	}

	User struct {
		Username          string `json:"username"`
		Password          string `json:"password"`
		Gender            string `json:"gender"`
		MaritalStatus     string `json:"marital_status"`
		IssuedAt          string `json:"issued_at"`
		Name              string `json:"name"`
		BirthDate         string `json:"birthdate"`
		BirthPlace        string `json:"birthplace"`
		MobileCallingCode string `json:"mobile_calling_code"`
		MobileNo          string `json:"mobile_no"`
		Address           string `json:"address"`
		CountryCode       string `json:"country_code"`
		StateId           string `json:"state_id"`
		IdentityType      string `json:"identity_type"`
		IdentityNumber    string `json:"identity_no"`
		CityId            string `json:"city_id"`
	}

	UserLogin struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	RoleReq struct {
		Role string `json:"role"`
	}
)

func GetUser(w http.ResponseWriter, r *http.Request) {
	//param := getUrlParam(r.URL.String())

	user, err := model.FindAllAccount()
	if err != nil && err != model.ErrResourceNotFound {
		log.Panic(err)
	}

	res := NewResponse(user)
	render.JSON(w, res)
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	token := r.FormValue(`token`)

	var rd UserReq
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		log.Panic(err)
	}

	url := "http://juno-staging.elys.id/v1/api/accounts?token=" + token
	json, err := json.Marshal(rd.UserValue)
	if err != nil {
		log.Panic(err)
	}
	responseId := sendPost(url, string(json))
	id := getResponseData(responseId)
	fmt.Println(id)

	user := getProfile(token, id["id"].(string))

	d := model.AccountDetail{
		CompanyID:   user["company_id"].(string),
		UserID:      user["id"].(string),
		AccountRole: rd.AccountRole,
		AssignBy:    rd.AssignBy,
		CreatedBy:   rd.CreatedBy,
	}

	if err := model.AddAccount(d); err != nil {
		log.Panic(err)
	}

	res := NewResponse(user)
	render.JSON(w, res, http.StatusCreated)
}

func GetToken(w http.ResponseWriter, r *http.Request) {
	var rd UserLogin
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		log.Panic(err)
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", "http://juno-staging.elys.id/v1/api/token", nil)
	req.SetBasicAuth(rd.Username, rd.Password)
	resp, err := client.Do(req)
	if err != nil {
		log.Panic(err)
	}
	bodyText, err := ioutil.ReadAll(resp.Body)

	res := NewResponse(bodyText)
	render.JSON(w, res, http.StatusCreated)
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

func getProfile(token string, id string) map[string]interface{} {
	url := "http://juno-staging.elys.id/v1/api/accounts/" + id + "?token=" + token
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	robots, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	return getResponseData(robots)
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
