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

	UserResponse struct {
		Id    string
		Token string
	}
)

func DoLogin(w http.ResponseWriter, r *http.Request) {
	var rd UserLogin
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		log.Panic(err)
	}

	user, err := model.Login(rd.Username, hash(rd.Password), rd.AccountId)
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
		Path:   "/",
	}

	session.Values["expired"] = times.Format("2006-01-02 15:04:05")
	session.Save(r, w)
	resp := UserResponse{
		Id:    user,
		Token: encoded,
	}
	res := NewResponse(resp)
	render.JSON(w, res, http.StatusOK)
}

func FindUserByRole(w http.ResponseWriter, r *http.Request) {
	role := r.FormValue("role")

	var accountId string
	var user = model.Response{}
	var err error
	status := http.StatusUnauthorized
	token := r.FormValue("token")
	userId := r.FormValue("user")
	valid := false

	if token != "" && token != "null" {
		accountId, _, valid, _ = getValiditySession(r, userId, token)
	}

	if valid {
		status = http.StatusOK
		user, err = model.FindUserByRole(role, accountId)
		if err != nil {
			status = http.StatusInternalServerError
			if err != model.ErrResourceNotFound {
				//log.Panic(err)
				status = http.StatusNotFound
			}
		}
	}

	res := NewResponse(user)
	render.JSON(w, res, status)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	accountId := r.FormValue("account_id")

	var user = model.Response{}
	var err error
	status := http.StatusUnauthorized
	token := r.FormValue("token")
	userId := r.FormValue("user")
	valid := false

	if token != "" && token != "null" {
		_, _, valid, _ = getValiditySession(r, userId, token)
	}

	if valid {
		status = http.StatusOK
		user, err = model.FindAllUser(accountId)
		if err != nil {
			status = http.StatusInternalServerError
			if err != model.ErrResourceNotFound {
				//log.Panic(err)
				status = http.StatusNotFound
			}
		}
	}

	res := NewResponse(user)
	render.JSON(w, res, status)
}

func GetUserCustomParam(w http.ResponseWriter, r *http.Request) {
	param := getUrlParam(r.URL.String())

	var user = model.Response{}
	var err error
	status := http.StatusUnauthorized
	token := r.FormValue("token")
	userId := r.FormValue("user")
	valid := false

	if token != "" && token != "null" {
		_, _, valid, _ = getValiditySession(r, userId, token)
	}

	if valid {
		status = http.StatusOK
		user, err = model.FindUser(param)
		if err != nil {
			status = http.StatusInternalServerError
			if err != model.ErrResourceNotFound {
				//log.Panic(err)
				status = http.StatusNotFound
			}
		}
	}

	res := NewResponse(user)
	render.JSON(w, res, status)
}

func CheckSession(w http.ResponseWriter, r *http.Request) {
	token := r.FormValue("token")
	user := r.FormValue("user")
	valid := false

	if token != "" && token != "null" {
		_, _, valid, _ = getValiditySession(r, user, token)
	}

	res := NewResponse(valid)
	render.JSON(w, res, http.StatusOK)
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	token := r.FormValue("token")
	user := r.FormValue("user")
	valid := false
	var accountId string
	if token != "" && token != "null" {
		accountId, _, valid, _ = getValiditySession(r, user, token)
	}

	var rd User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		log.Panic(err)
	}

	status := http.StatusUnauthorized
	if valid {
		status = http.StatusOK
		param := model.User{
			AccountId: accountId,
			Username:  rd.Username,
			Password:  hash(rd.Password),
			Email:     rd.Email,
			Phone:     rd.Phone,
			RoleId:    rd.RoleId,
			CreatedBy: user,
		}

		if err := model.AddUser(param); err != nil {
			//log.Panic(err)
			status = http.StatusInternalServerError
		}
	}

	res := NewResponse(nil)
	render.JSON(w, res, status)
}

// local function
func getValiditySession(r *http.Request, user string, token string) (string, time.Time, bool, error) {
	valid := false
	decoded, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return "", time.Now(), false, err
	}

	fmt.Println(string(decoded))
	session := strings.Split(string(decoded), ";")

	if session[0] == user {
		sessionValue, err := store.Get(r, session[0])
		if err != nil {
			return "", time.Now(), false, err
			// log.Panic(err)
		}
		fmt.Println(session[0])
		ds := sessionValue.Values
		exp, err := time.Parse("2006-01-02 15:04:05", ds["expired"].(string))
		if err != nil {
			return "", time.Now(), false, err
			// log.Panic(err)
		}
		fmt.Println(session[0])

		if exp.After(time.Now()) {
			valid = true
		}
		return session[1], exp, valid, nil
	}

	return "", time.Now(), false, nil
}
