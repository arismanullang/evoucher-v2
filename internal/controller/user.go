package controller

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	//"github.com/go-zoo/bone"
	"github.com/gorilla/sessions"
	"github.com/ruizu/render"

	"github.com/gilkor/evoucher/internal/model"
)

var store = sessions.NewCookieStore([]byte("lalala"))

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
		Username  string `json:"username"`
		Password  string `json:"password"`
		AccountId string `json:"account_id"`
	}
)

func DoLogin(w http.ResponseWriter, r *http.Request) {
	var rd UserLogin
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		log.Panic(err)
	}

	user, err := model.Login(rd.Username, hash(rd.Password))
	if err != nil && err != model.ErrResourceNotFound {
		log.Panic(err)
	}

	times := time.Now()
	times = times.AddDate(0, 0, 1)
	encoded := base64.StdEncoding.EncodeToString([]byte(user + ";" + rd.AccountId + ";" + times.String()))

	session, err := store.Get(r, user)
	if err != nil {
		log.Panic(err)
	}

	session.Options = &sessions.Options{
		MaxAge: 86400,
	}

	session.Values["expired"] = times.Format("2006-01-02 15:04:05")
	session.Save(r, w)

	res := NewResponse(encoded)
	render.JSON(w, res, http.StatusOK)
}

func FindUserByRole(w http.ResponseWriter, r *http.Request) {
	accountId := r.FormValue("account_id")
	role := r.FormValue("role")

	var user = model.Response{}
	var err error
	var status int
	if _, ok := basicAuth(w, r); ok {
		user, err = model.FindUserByRole(role, accountId)
		if err != nil && err != model.ErrResourceNotFound {
			log.Panic(err)
		}
		status = http.StatusOK
	} else {
		status = http.StatusUnauthorized
	}

	res := NewResponse(user)
	render.JSON(w, res, status)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	accountId := r.FormValue("account_id")

	var user = model.Response{}
	var err error
	var status int
	if _, ok := basicAuth(w, r); ok {
		user, err = model.FindAllUser(accountId)
		if err != nil && err != model.ErrResourceNotFound {
			log.Panic(err)
		}
		status = http.StatusOK
	} else {
		status = http.StatusUnauthorized
	}

	res := NewResponse(user)
	render.JSON(w, res, status)
}

func GetUserCustomParam(w http.ResponseWriter, r *http.Request) {
	param := getUrlParam(r.URL.String())

	var user = model.Response{}
	var err error
	var status int
	if _, ok := basicAuth(w, r); ok {
		user, err = model.FindUser(param)
		if err != nil && err != model.ErrResourceNotFound {
			log.Panic(err)
		}
		status = http.StatusOK
	} else {
		status = http.StatusUnauthorized
	}

	res := NewResponse(user)
	render.JSON(w, res, status)
}

// only dashboard api
func CheckSession(w http.ResponseWriter, r *http.Request) {
	token := r.FormValue("token")
	valid := false
	if token != "" && token != "null" {

		_, _, exp, _ := checkExpired(r, token)

		if exp.After(time.Now()) {
			valid = true
		}
	}

	res := NewResponse(valid)
	render.JSON(w, res, http.StatusOK)
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	token := r.FormValue("token")
	valid := false
	if token != "" && token != "null" {

		_, _, exp, _ := checkExpired(r, token)

		if exp.After(time.Now()) {
			valid = true
		}
	}

	var rd User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		log.Panic(err)
	}

	fmt.Println(len(hash(rd.Password)))

	var status int
	if valid {
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
		status = http.StatusOK
	} else {
		status = http.StatusUnauthorized
	}

	res := NewResponse(nil)
	render.JSON(w, res, status)
}

// local function
func checkExpired(r *http.Request, token string) (string, string, time.Time, error) {
	decoded, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return "", "", time.Now(), err
		// log.Panic(err)
	}

	fmt.Println(string(decoded))
	session := strings.Split(string(decoded), ";")
	sessionValue, err := store.Get(r, session[0])
	if err != nil {
		return "", "", time.Now(), err
		// log.Panic(err)
	}

	user := sessionValue.Values
	exp, err := time.Parse("2006-01-02 15:04:05", user["expired"].(string))
	if err != nil {
		return "", "", time.Now(), err
		// log.Panic(err)
	}

	return session[0], session[1], exp, nil
}

func getAccountID(userID string) string {
	return model.GetAccountByUser(userID)
}
