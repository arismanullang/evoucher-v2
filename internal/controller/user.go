package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

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

	PasswordReq struct {
		Password string `json:"password"`
	}
)

func FindUserByRole(w http.ResponseWriter, r *http.Request) {
	role := r.FormValue("role")
	status := http.StatusUnauthorized
	err := model.ErrTokenNotFound
	res := NewResponse(nil)

	res.AddError(its(status), its(status), err.Error(), "user")
	accountId, _, _, valid := AuthToken(w, r)
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
	status := http.StatusUnauthorized
	err := model.ErrTokenNotFound
	res := NewResponse(nil)

	res.AddError(its(status), its(status), err.Error(), "user")
	_, user, _, valid := AuthToken(w, r)
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
	status := http.StatusUnauthorized
	err := model.ErrTokenNotFound
	res := NewResponse(nil)

	res.AddError(its(status), its(status), err.Error(), "user")

	accountId, _, _, valid := AuthToken(w, r)
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

	status := http.StatusUnauthorized
	err := model.ErrTokenNotFound
	res := NewResponse(nil)

	res.AddError(its(status), its(status), err.Error(), "user")

	_, _, _, valid := AuthToken(w, r)
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

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	status := http.StatusUnauthorized
	err := model.ErrTokenNotFound
	res := NewResponse(nil)

	res.AddError(its(status), its(status), err.Error(), "user")

	accountId, user, _, valid := AuthToken(w, r)

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

func ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var username = r.FormValue("username")
	fmt.Println(username)
	if err := model.SendMail(model.Domain, model.ApiKey, model.PublicApiKey, username); err != nil {
		log.Fatal(err)
	}
}

func UpdatePassword(w http.ResponseWriter, r *http.Request) {
	var rd PasswordReq
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		log.Panic(err)
	}

	password := r.FormValue("key")
	user, err := model.GetSession(password)
	if err != nil {
		log.Panic(err)
	}

	err = model.UpdatePassword(user.UserId, hash(rd.Password))
	if err != nil {
		log.Panic(err)
	}
	render.JSON(w, http.StatusOK)
}
