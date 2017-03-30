package controller

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	//"github.com/go-zoo/bone"
	"github.com/gorilla/sessions"
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
		Username  string `json:"username"`
		Password  string `json:"password"`
		AccountId string `json:"account_id"`
	}

	UserResponse struct {
		Id    string
		Token string
	}

	Session struct {
		AccountId string
		UserId    string
		Expired   time.Time
	}
)

var store = sessions.NewCookieStore([]byte("lalala"))

func DoLogin(w http.ResponseWriter, r *http.Request) {
	var rd UserLogin
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		log.Panic(err)
	}

	status := http.StatusOK
	res := NewResponse(nil)

	user, err := model.Login(rd.Username, hash(rd.Password), rd.AccountId)
	if err != nil {
		status = http.StatusUnauthorized
		if err == model.ErrResourceNotFound {
			status = http.StatusNotFound
		}
		res.AddError(its(status), its(status), err.Error(), "variant")
	} else {
		times := time.Now()
		times = times.AddDate(0, 0, 1)
		encoded := base64.StdEncoding.EncodeToString([]byte(user + ";" + rd.AccountId + ";" + times.String()))
		fmt.Println(user)
		encoded = replaceSpecialCharacter(encoded)
		session, err := store.Get(r, encoded)

		if err != nil {
			status = http.StatusInternalServerError
			res.AddError(its(status), its(status), err.Error(), "user")
		} else {
			session.Options = &sessions.Options{
				MaxAge: 86400,
				Path:   "/",
			}
			fmt.Println(user)
			session.Values["user"] = user
			session.Values["account"] = rd.AccountId
			session.Values["expired"] = times.Format("2006-01-02 15:04:05")
			session.Save(r, w)
			resp := UserResponse{
				Id:    user,
				Token: encoded,
			}

			res = NewResponse(resp)
		}

	}

	render.JSON(w, res, status)
}

func FindUserByRole(w http.ResponseWriter, r *http.Request) {
	role := r.FormValue("role")
	accountId := ""
	token := r.FormValue("token")
	status := http.StatusUnauthorized
	err := model.ErrTokenNotFound
	res := NewResponse(nil)

	res.AddError(its(status), its(status), err.Error(), "user")

	valid := false
	if token != "" && token != "null" {
		_, accountId, _, valid, _ = getValiditySession(r, token)
	}

	if valid {
		status = http.StatusOK
		user, err := model.FindUsersByRole(role, accountId)
		if err != nil {
			status = http.StatusInternalServerError
			if err != model.ErrResourceNotFound {
				status = http.StatusNotFound
			}

			res.AddError(its(status), its(status), err.Error(), "user")
		} else {
			res = NewResponse(user)
		}
	}

	render.JSON(w, res, status)
}

func GetUserDetails(w http.ResponseWriter, r *http.Request) {
	user := ""
	token := r.FormValue("token")
	status := http.StatusUnauthorized
	err := model.ErrTokenNotFound
	res := NewResponse(nil)

	res.AddError(its(status), its(status), err.Error(), "user")

	valid := false
	if token != "" && token != "null" {
		user, _, _, valid, _ = getValiditySession(r, token)
	}

	if valid {
		status = http.StatusOK
		user, err := model.FindUserDetail(user)
		if err != nil {
			status = http.StatusInternalServerError
			if err != model.ErrResourceNotFound {
				status = http.StatusNotFound
			}

			res.AddError(its(status), its(status), err.Error(), "user")
		} else {
			res = NewResponse(user)
		}
	}

	render.JSON(w, res, status)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	accountId := ""
	token := r.FormValue("token")
	status := http.StatusUnauthorized
	err := model.ErrTokenNotFound
	res := NewResponse(nil)

	res.AddError(its(status), its(status), err.Error(), "user")

	valid := false
	if token != "" && token != "null" {
		_, accountId, _, valid, _ = getValiditySession(r, token)
	}

	if valid {
		status = http.StatusOK
		user, err := model.FindAllUsers(accountId)
		if err != nil {
			status = http.StatusInternalServerError
			if err != model.ErrResourceNotFound {
				status = http.StatusNotFound
			}

			res.AddError(its(status), its(status), err.Error(), "user")
		} else {
			res = NewResponse(user)
		}
	}

	render.JSON(w, res, status)
}

func GetUserCustomParam(w http.ResponseWriter, r *http.Request) {
	param := getUrlParam(r.URL.String())

	token := r.FormValue("token")
	status := http.StatusUnauthorized
	err := model.ErrTokenNotFound
	res := NewResponse(nil)

	res.AddError(its(status), its(status), err.Error(), "user")

	valid := false
	if token != "" && token != "null" {
		_, _, _, valid, _ = getValiditySession(r, token)
	}

	if valid {
		status = http.StatusOK
		user, err := model.FindUsersCustomParam(param)
		if err != nil {
			status = http.StatusInternalServerError
			if err != model.ErrResourceNotFound {
				status = http.StatusNotFound
			}

			res.AddError(its(status), its(status), err.Error(), "user")
		} else {
			res = NewResponse(user)
		}
	}

	render.JSON(w, res, status)
}

func CheckSession(w http.ResponseWriter, r *http.Request) {
	token := r.FormValue("token")
	valid := false

	if token != "" && token != "null" {
		_, _, _, valid, _ = getValiditySession(r, token)
	}

	res := NewResponse(valid)
	render.JSON(w, res, http.StatusOK)
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	user := ""
	accountId := ""
	token := r.FormValue("token")
	status := http.StatusUnauthorized
	err := model.ErrTokenNotFound
	res := NewResponse(nil)

	res.AddError(its(status), its(status), err.Error(), "user")

	valid := false
	if token != "" && token != "null" {
		user, accountId, _, valid, _ = getValiditySession(r, token)
	}

	var rd User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		log.Panic(err)
	}

	if valid {
		fmt.Println("Valid")
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
			fmt.Print(err.Error())
			status = http.StatusInternalServerError
			if err == model.ErrResourceNotFound {
				status = http.StatusNotFound
			}

			res.AddError(its(status), its(status), err.Error(), "user")
		} else {
			res = NewResponse(user)
		}
	}

	render.JSON(w, res, status)
}

// Return : user_id, account_id, expired, boolean, error
func getValiditySession(r *http.Request, token string) (string, string, time.Time, bool, error) {
	fmt.Println("Check Session")
	fmt.Println(r)
	sessionValue, err := store.Get(r, token)
	if err != nil {
		return "", "", time.Now(), false, model.ErrTokenNotFound
	}
	ds := sessionValue.Values
	if len(ds) == 0 {
		return "", "", time.Now(), false, model.ErrTokenNotFound
	}

	exp, err := time.Parse("2006-01-02 15:04:05", ds["expired"].(string))
	if err != nil {
		//log.Panic(err)
		return "", "", time.Now(), false, model.ErrTokenExpired
	}

	if exp.Before(time.Now()) {
		return "", "", time.Now(), false, model.ErrTokenExpired
	}

	return ds["user"].(string), ds["account"].(string), exp, true, nil
	// return "g6mrRguA", "So-GAf-G", time.Now(), true, nil
}
