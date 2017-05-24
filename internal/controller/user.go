package controller

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
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

	ChangePasswordReq struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}

	PasswordReq struct {
		Password string `json:"password"`
	}
)

func InsertBroadcastUser(w http.ResponseWriter, r *http.Request) {
	var listTarget []string
	var listDescription []string
	variantId := r.FormValue("variant-id")

	r.ParseMultipartForm(32 << 20)
	f, _, err := r.FormFile("list-target")
	if err == http.ErrMissingFile {
		err = model.ErrResourceNotFound
	}
	if err != nil {
		err = model.ErrServerInternal
	}
	//	fmt.Println(r)
	fmt.Println("f : ")
	fmt.Println(f)

	read := csv.NewReader(bufio.NewReader(f))
	for {
		record, err := read.Read()
		// Stop at EOF.
		if err == io.EOF {
			break
		}

		for value := range record {
			temp := strings.Split(record[value], ",")
			listTarget = append(listTarget, temp[1])
			listDescription = append(listDescription, temp[2])
			fmt.Println(temp)
		}
	}

	res := NewResponse(nil)
	status := http.StatusUnauthorized
	err = model.ErrTokenNotFound
	errTitle := model.ErrCodeInvalidToken
	res.AddError(its(status), errTitle, err.Error(), "Insert Broadcast")
	_, user, _, valid := AuthToken(w, r)
	if valid {
		status = http.StatusCreated

		if err := model.InsertBroadcastUser(variantId, user, listTarget, listDescription); err != nil {
			//log.Panic(err)
			status = http.StatusInternalServerError
			errTitle = model.ErrCodeInternalError
			res.AddError(its(status), errTitle, err.Error(), "Insert Broadcast")
		} else {
			res = NewResponse(nil)

		}

	}

	render.JSON(w, res, status)

}

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

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	status := http.StatusUnauthorized
	err := model.ErrTokenNotFound
	res := NewResponse(nil)

	res.AddError(its(status), its(status), err.Error(), "user")

	_, user, _, valid := AuthToken(w, r)

	var rd User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		log.Panic(err)
	}

	if valid {
		fmt.Println("Valid")
		status = http.StatusOK
		param := model.User{
			Username:  rd.Username,
			Email:     rd.Email,
			Phone:     rd.Phone,
			CreatedBy: user,
		}

		if err := model.UpdateUser(param); err != nil {
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

func ChangePassword(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Change Password")
	status := http.StatusUnauthorized
	err := model.ErrTokenNotFound
	res := NewResponse(nil)

	res.AddError(its(status), its(status), err.Error(), "user")

	_, user, _, valid := AuthToken(w, r)

	var rd ChangePasswordReq
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		log.Panic(err)
	}

	if valid {
		fmt.Println("Valid")
		status = http.StatusOK

		if err := model.ChangePassword(user, hash(rd.OldPassword), hash(rd.NewPassword)); err != nil {
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
