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

		fmt.Println(record)
		if len(record) > 0 {
			listTarget = append(listTarget, strings.Replace(record[1], "'", "", -1))
			listDescription = append(listDescription, strings.Replace(record[2], "'", "", -1))
		} else {
			for value := range record {
				fmt.Println(record[value])
				temp := strings.Split(record[value], ";")
				listTarget = append(listTarget, strings.Replace(temp[1], "'", "", -1))
				listDescription = append(listDescription, strings.Replace(temp[2], "'", "", -1))
				fmt.Println(temp)
			}
		}
	}

	res := NewResponse(nil)
	status := http.StatusUnauthorized
	err = model.ErrTokenNotFound
	errTitle := model.ErrCodeInvalidToken
	res.AddError(its(status), errTitle, err.Error(), "Insert Broadcast")
	a := AuthToken(w, r)
	if a.Valid {
		status = http.StatusCreated

		if err := model.InsertBroadcastUser(variantId, a.User.ID, listTarget, listDescription); err != nil {
			//log.Panic(err)
			status = http.StatusInternalServerError
			errTitle = model.ErrCodeInternalError
			res.AddError(its(status), errTitle, err.Error(), "Insert Broadcast")
		} else {
			res = NewResponse(nil)

		}

		if err := model.UpdateBulkVariant(variantId, len(listTarget)); err != nil {
			//log.Panic(err)
			status = http.StatusInternalServerError
			errTitle = model.ErrCodeInternalError
			res.AddError(its(status), errTitle, err.Error(), "Update Variant")
		}

	}else {
		res = a.res
		status = http.StatusUnauthorized
	}

	render.JSON(w, res, status)

}

func FindUserByRole(w http.ResponseWriter, r *http.Request) {
	role := r.FormValue("role")
	status := http.StatusUnauthorized
	err := model.ErrTokenNotFound
	res := NewResponse(nil)

	res.AddError(its(status), its(status), err.Error(), "user")
	a := AuthToken(w, r)
	if a.Valid {
		status = http.StatusOK
		user, err := model.FindUsersByRole(role, a.User.AccountID)
		if err != nil {
			status = http.StatusInternalServerError
			if err != model.ErrResourceNotFound {
				status = http.StatusNotFound
			}

			res.AddError(its(status), its(status), err.Error(), "user")
		} else {
			res = NewResponse(user)
		}
	}else {
		res = a.res
		status = http.StatusUnauthorized
	}

	render.JSON(w, res, status)
}

func GetUserDetails(w http.ResponseWriter, r *http.Request) {
	status := http.StatusUnauthorized
	err := model.ErrTokenNotFound
	res := NewResponse(nil)

	res.AddError(its(status), its(status), err.Error(), "user")
	a := AuthToken(w, r)
	if a.Valid {
		status = http.StatusOK
		user, err := model.FindUserDetail(a.User.AccountID)
		if err != nil {
			status = http.StatusInternalServerError
			if err != model.ErrResourceNotFound {
				status = http.StatusNotFound
			}

			res.AddError(its(status), its(status), err.Error(), "user")
		} else {
			res = NewResponse(user)
		}
	}else {
		res = a.res
		status = http.StatusUnauthorized
	}

	render.JSON(w, res, status)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	status := http.StatusUnauthorized
	err := model.ErrTokenNotFound
	res := NewResponse(nil)

	res.AddError(its(status), its(status), err.Error(), "user")

	a := AuthToken(w, r)
	if a.Valid {
		status = http.StatusOK
		user, err := model.FindAllUsers(a.User.AccountID)
		if err != nil {
			status = http.StatusInternalServerError
			if err != model.ErrResourceNotFound {
				status = http.StatusNotFound
			}

			res.AddError(its(status), its(status), err.Error(), "user")
		} else {
			res = NewResponse(user)
		}
	}else {
		res = a.res
		status = http.StatusUnauthorized
	}

	render.JSON(w, res, status)
}

func GetUserCustomParam(w http.ResponseWriter, r *http.Request) {
	param := getUrlParam(r.URL.String())

	status := http.StatusUnauthorized
	err := model.ErrTokenNotFound
	res := NewResponse(nil)

	res.AddError(its(status), its(status), err.Error(), "user")

	a := AuthToken(w, r)
	if a.Valid {
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


	var rd User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		log.Panic(err)
	}

	a := AuthToken(w, r)
	if a.Valid {
		fmt.Println("Valid")
		status = http.StatusOK
		param := model.User{
			AccountID: a.User.AccountID,
			Username:  rd.Username,
			Password:  hash(rd.Password),
			Email:     rd.Email,
			Phone:     rd.Phone,
			Role:      a.User.Role,
			CreatedBy: a.User.ID,
		}

		if err := model.AddUser(param); err != nil {
			fmt.Print(err.Error())
			status = http.StatusInternalServerError
			if err == model.ErrResourceNotFound {
				status = http.StatusNotFound
			}

			res.AddError(its(status), its(status), err.Error(), "user")
		} else {
			res = NewResponse(a.User.ID)
		}
	}else {
		res = a.res
		status = http.StatusUnauthorized
	}

	render.JSON(w, res, status)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	status := http.StatusUnauthorized
	err := model.ErrTokenNotFound
	res := NewResponse(nil)

	res.AddError(its(status), its(status), err.Error(), "user")


	var rd User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		log.Panic(err)
	}

	a := AuthToken(w, r)
	if a.Valid {
		fmt.Println("Valid")
		status = http.StatusOK
		param := model.User{
			Username:  rd.Username,
			Email:     rd.Email,
			Phone:     rd.Phone,
			CreatedBy: a.User.ID,
		}

		if err := model.UpdateUser(param); err != nil {
			fmt.Print(err.Error())
			status = http.StatusInternalServerError
			if err == model.ErrResourceNotFound {
				status = http.StatusNotFound
			}

			res.AddError(its(status), its(status), err.Error(), "user")
		} else {
			res = NewResponse(a.User.ID)
		}
	}else {
		res = a.res
		status = http.StatusUnauthorized
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

	err = model.UpdatePassword(user.User.ID, hash(rd.Password))
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


	var rd ChangePasswordReq
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		log.Panic(err)
	}

	a := AuthToken(w, r)
	if a.Valid {
		fmt.Println("Valid")
		status = http.StatusOK

		if err := model.ChangePassword(a.User.ID, hash(rd.OldPassword), hash(rd.NewPassword)); err != nil {
			fmt.Print(err.Error())
			status = http.StatusInternalServerError
			if err == model.ErrResourceNotFound {
				status = http.StatusNotFound
			}

			res.AddError(its(status), its(status), err.Error(), "user")
		} else {
			res = NewResponse(a.User.ID)
		}
	}else {
		res = a.res
		status = http.StatusUnauthorized
	}

	render.JSON(w, res, status)
}
