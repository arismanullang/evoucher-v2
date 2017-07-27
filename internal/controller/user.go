package controller

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	//"github.com/go-zoo/bone"

	"github.com/ruizu/render"

	"github.com/gilkor/evoucher/internal/model"
)

type (
	User struct {
		Id        string   `json:"id"`
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

	ChangeUserStatusReq struct {
		Id string `json:"id"`
	}
)

func InsertBroadcastUser(w http.ResponseWriter, r *http.Request) {
	var listTarget []string
	var listDescription []string
	programId := r.FormValue("program-id")
	apiName := "broadcast_create"

	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag(apiName)

	r.ParseMultipartForm(32 << 20)
	f, _, err := r.FormFile("list-target")
	if err == http.ErrMissingFile {
		err = model.ErrResourceNotFound
	}
	if err != nil {
		err = model.ErrServerInternal
	}

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
			}
		}
	}

	res := NewResponse(nil)
	status := http.StatusCreated

	a := AuthToken(w, r)
	if !a.Valid {
		res = a.res
		status = http.StatusUnauthorized
		render.JSON(w, res, status)
		return
	}

	if err := model.InsertBroadcastUser(programId, a.User.ID, listTarget, listDescription); err != nil {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, err.Error(), logger.TraceID)
		logger.SetStatus(status).Info("param :", programId+" || "+strings.Join(listTarget, ";"), "response :", err.Error())
	}

	if err := model.UpdateBulkProgram(programId, len(listTarget)); err != nil {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, err.Error(), logger.TraceID)
		logger.SetStatus(status).Info("param :", programId+" || "+its(len(listTarget)), "response :", err.Error())
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
		user, err := model.FindUsersByRole(role, a.User.Account.Id)
		if err != nil {
			status = http.StatusInternalServerError
			if err != model.ErrResourceNotFound {
				status = http.StatusNotFound
			}

			res.AddError(its(status), its(status), err.Error(), "user")
		} else {
			res = NewResponse(user)
		}
	} else {
		res = a.res
		status = http.StatusUnauthorized
	}

	render.JSON(w, res, status)
}

func GetOtherUserDetails(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	res := NewResponse(nil)
	status := http.StatusOK

	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag("user_other_detail")

	a := AuthTokenWithLogger(w, r, logger)
	if !a.Valid {
		res = a.res
		status = http.StatusUnauthorized
		render.JSON(w, res, status)
		return
	}

	user, err := model.FindUserDetail(id)
	res = NewResponse(user)
	if err != nil {
		status = http.StatusInternalServerError
		errTitle := model.ErrCodeInternalError
		if err != model.ErrResourceNotFound {
			status = http.StatusNotFound
			errTitle = model.ErrCodeResourceNotFound
		}

		res.AddError(its(status), errTitle, err.Error(), logger.TraceID)
		logger.SetStatus(status).Info("param :", id, "response :", err.Error())
	}

	render.JSON(w, res, status)
}

func GetUserDetails(w http.ResponseWriter, r *http.Request) {
	status := http.StatusOK
	res := NewResponse(nil)

	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag("user_detail")

	a := AuthTokenWithLogger(w, r, logger)
	if !a.Valid {
		res = a.res
		status = http.StatusUnauthorized
		render.JSON(w, res, status)
		return
	}

	user, err := model.FindUserDetail(a.User.ID)
	res = NewResponse(user)
	if err != nil {
		status = http.StatusInternalServerError
		errTitle := model.ErrCodeInternalError
		if err != model.ErrResourceNotFound {
			status = http.StatusNotFound
			errTitle = model.ErrCodeResourceNotFound
		}

		res.AddError(its(status), errTitle, err.Error(), logger.TraceID)
		logger.SetStatus(status).Info("param :", a.User.ID, "response :", err.Error())
	}

	render.JSON(w, res, status)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	status := http.StatusOK
	res := NewResponse(nil)

	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag("user_get")

	a := AuthTokenWithLogger(w, r, logger)
	if !a.Valid {
		res = a.res
		status = http.StatusUnauthorized
	}

	user, err := model.FindAllUsers(a.User.Account.Id)
	res = NewResponse(user)
	if err != nil {
		status = http.StatusInternalServerError
		errTitle := model.ErrCodeInternalError
		if err != model.ErrResourceNotFound {
			status = http.StatusNotFound
			errTitle = model.ErrCodeResourceNotFound
		}

		res.AddError(its(status), errTitle, err.Error(), logger.TraceID)
		logger.SetStatus(status).Info("param :", a.User.Account.Id, "response :", err.Error())
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
	apiName := "user_create"
	valid := false

	status := http.StatusCreated
	res := NewResponse(nil)

	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag(apiName)

	var rd User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		logger.SetStatus(http.StatusInternalServerError).Panic("param :", r.Body, "response :", err.Error())
	}

	a := AuthTokenWithLogger(w, r, logger)
	if !a.Valid {
		res = a.res
		status = http.StatusUnauthorized
		render.JSON(w, res, status)
		return
	}

	for _, valueRole := range a.User.Role {
		features := model.ApiFeatures[valueRole.Detail]
		for _, valueFeature := range features {
			if apiName == valueFeature {
				valid = true
			}
		}
	}

	if !valid {
		logger.SetStatus(status).Info("param :", a.User.ID, "response :", "Invalid Role")

		status = http.StatusUnauthorized
		res.AddError(its(status), model.ErrCodeInvalidRole, model.ErrInvalidRole.Error(), logger.TraceID)
		render.JSON(w, res, status)
		return
	}

	param := model.RegisterUser{
		Username: rd.Username,
		Password: hash(rd.Password),
		Email:    rd.Email,
		Phone:    rd.Phone,
		Role:     rd.RoleId,
	}

	if err := model.AddUser(param, a.User.ID, a.User.Account.Id); err != nil {
		status = http.StatusInternalServerError
		errTitle := model.ErrCodeInternalError
		if err == model.ErrResourceNotFound {
			status = http.StatusNotFound
			errTitle = model.ErrCodeResourceNotFound
		}

		res.AddError(its(status), errTitle, err.Error(), logger.TraceID)
		logger.SetStatus(status).Info("param :", param, "response :", err.Error())
	}

	render.JSON(w, res, status)
}

func UpdateUserRoute(w http.ResponseWriter, r *http.Request) {
	types := r.FormValue("type")
	apiName := "user_update"

	res := NewResponse(nil)
	status := http.StatusInternalServerError
	errTitle := model.ErrCodeInternalError
	err := model.ErrServerInternal

	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag(apiName)

	a := AuthTokenWithLogger(w, r, logger)
	if !a.Valid {
		res = a.res
		status = http.StatusUnauthorized
		render.JSON(w, res, status)
		return
	}

	if types == "" {
		res.AddError(its(status), errTitle, err.Error(), logger.TraceID)
		logger.SetStatus(status).Info("param :", "empty types", "response :", err.Error())
		render.JSON(w, res, status)
	} else if types == "detail" {
		UpdateUser(w, r, logger, a)
	} else if types == "other" {
		UpdateOtherUser(w, r, logger, a)
	} else if types == "password" {
		ChangePassword(w, r, logger, a)
	} else if types == "reset" {
		ResetOtherPassword(w, r, logger, a)
	} else {
		res.AddError(its(status), errTitle, err.Error(), logger.TraceID)
		logger.SetStatus(status).Info("param :", "types not allowed", "response :", err.Error())
		render.JSON(w, res, status)
	}
}

func UpdateOtherUser(w http.ResponseWriter, r *http.Request, logger *model.LogField, a Auth) {
	apiName := "user_update"
	valid := false

	status := http.StatusOK
	res := NewResponse(nil)

	var rd User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		logger.SetStatus(http.StatusInternalServerError).Panic("param :", r.Body, "response :", err.Error())
	}

	for _, valueRole := range a.User.Role {
		features := model.ApiFeatures[valueRole.Detail]
		for _, valueFeature := range features {
			if apiName == valueFeature {
				valid = true
			}
		}
	}

	if !valid {
		logger.SetStatus(status).Info("param :", a.User.ID, "response :", "Invalid Role")

		status = http.StatusUnauthorized
		res.AddError(its(status), model.ErrCodeInvalidRole, model.ErrInvalidRole.Error(), logger.TraceID)
		render.JSON(w, res, status)
		return
	}

	roles := []model.Role{}
	for _, v := range rd.RoleId {
		tempRole := model.Role{
			Id:     v,
			Detail: "",
		}
		roles = append(roles, tempRole)
	}

	param := model.User{
		ID:        rd.Id,
		Username:  rd.Username,
		Email:     rd.Email,
		Phone:     rd.Phone,
		Role:      roles,
		CreatedBy: a.User.ID,
	}

	if err := model.UpdateOtherUser(param); err != nil {
		status = http.StatusInternalServerError
		errTitle := model.ErrCodeInternalError
		if err == model.ErrResourceNotFound {
			status = http.StatusNotFound
			errTitle = model.ErrCodeResourceNotFound
		}

		res.AddError(its(status), errTitle, err.Error(), logger.TraceID)
		logger.SetStatus(status).Info("param :", param, "response :", err.Error())
	}

	render.JSON(w, res, status)
}

func UpdateUser(w http.ResponseWriter, r *http.Request, logger *model.LogField, a Auth) {
	status := http.StatusOK
	res := NewResponse(nil)

	var rd User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		logger.SetStatus(http.StatusInternalServerError).Panic("param :", r.Body, "response :", err.Error())
	}

	param := model.User{
		Username:  rd.Username,
		Email:     rd.Email,
		Phone:     rd.Phone,
		CreatedBy: a.User.ID,
	}

	if err := model.UpdateUser(param); err != nil {
		status = http.StatusInternalServerError
		errTitle := model.ErrCodeInternalError
		if err == model.ErrResourceNotFound {
			status = http.StatusNotFound
			errTitle = model.ErrCodeResourceNotFound
		}

		res.AddError(its(status), errTitle, err.Error(), logger.TraceID)
		logger.SetStatus(status).Info("param :", param, "response :", err.Error())
	}

	render.JSON(w, res, status)
}

func ChangePassword(w http.ResponseWriter, r *http.Request, logger *model.LogField, a Auth) {
	status := http.StatusOK
	res := NewResponse(nil)

	var rd ChangePasswordReq
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		logger.SetStatus(http.StatusInternalServerError).Panic("param :", r.Body, "response :", err.Error())
	}

	if err := model.ChangePassword(a.User.ID, hash(rd.OldPassword), hash(rd.NewPassword)); err != nil {
		status = http.StatusInternalServerError
		errTitle := model.ErrCodeInternalError
		if err == model.ErrResourceNotFound {
			status = http.StatusNotFound
			errTitle = model.ErrCodeResourceNotFound
		}

		res.AddError(its(status), errTitle, err.Error(), logger.TraceID)
		logger.SetStatus(status).Info("param :", a.User.ID+" || "+rd.OldPassword, "response :", err.Error())
	}

	render.JSON(w, res, status)
}

func ResetOtherPassword(w http.ResponseWriter, r *http.Request, logger *model.LogField, a Auth) {
	status := http.StatusOK
	res := NewResponse(nil)

	var rd ChangeUserStatusReq
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		logger.SetStatus(http.StatusInternalServerError).Panic("param :", r.Body, "response :", err.Error())
	}

	newPass := strings.ToLower(randStr(8, "Alphanumeric"))
	if err := model.ResetPassword(rd.Id, hash(newPass)); err != nil {
		status = http.StatusInternalServerError
		errTitle := model.ErrCodeInternalError
		if err == model.ErrResourceNotFound {
			status = http.StatusNotFound
			errTitle = model.ErrCodeResourceNotFound
		}

		res.AddError(its(status), errTitle, err.Error(), logger.TraceID)
		logger.SetStatus(status).Info("param :", rd, "response :", err.Error())
		render.JSON(w, res, status)
		return
	}

	res = NewResponse(newPass)
	render.JSON(w, res, status)
}

func BlockUser(w http.ResponseWriter, r *http.Request) {
	apiName := "user_block"
	valid := false

	status := http.StatusOK
	res := NewResponse(nil)

	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag(apiName)

	var rd ChangeUserStatusReq
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		logger.SetStatus(http.StatusInternalServerError).Panic("param :", r.Body, "response :", err.Error())
	}

	a := AuthTokenWithLogger(w, r, logger)
	if !a.Valid {
		res = a.res
		status = http.StatusUnauthorized
		render.JSON(w, res, status)
		return
	}

	for _, valueRole := range a.User.Role {
		features := model.ApiFeatures[valueRole.Detail]
		for _, valueFeature := range features {
			if apiName == valueFeature {
				valid = true
			}
		}
	}

	if !valid {
		logger.SetStatus(status).Info("param :", a.User.ID, "response :", "Invalid Role")

		status = http.StatusUnauthorized
		res.AddError(its(status), model.ErrCodeInvalidRole, model.ErrInvalidRole.Error(), logger.TraceID)
		render.JSON(w, res, status)
		return

	}

	if err := model.BlockUser(rd.Id); err != nil {
		status = http.StatusInternalServerError
		errTitle := model.ErrCodeInternalError
		if err == model.ErrResourceNotFound {
			status = http.StatusNotFound
			errTitle = model.ErrCodeResourceNotFound
		}

		res.AddError(its(status), errTitle, err.Error(), logger.TraceID)
		logger.SetStatus(status).Info("param :", rd.Id, "response :", err.Error())
	}

	render.JSON(w, res, status)
}

func ActivateUser(w http.ResponseWriter, r *http.Request) {
	apiName := "user_release"
	valid := false

	status := http.StatusOK
	res := NewResponse(nil)

	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag(apiName)

	var rd ChangeUserStatusReq
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		logger.SetStatus(http.StatusInternalServerError).Panic("param :", r.Body, "response :", err.Error())
	}

	a := AuthTokenWithLogger(w, r, logger)
	if !a.Valid {
		res = a.res
		status = http.StatusUnauthorized
		render.JSON(w, res, status)
		return
	}

	for _, valueRole := range a.User.Role {
		features := model.ApiFeatures[valueRole.Detail]
		for _, valueFeature := range features {
			if apiName == valueFeature {
				valid = true
			}
		}
	}

	if !valid {
		logger.SetStatus(status).Info("param :", a.User.ID, "response :", "Invalid Role")

		status = http.StatusUnauthorized
		res.AddError(its(status), model.ErrCodeInvalidRole, model.ErrInvalidRole.Error(), logger.TraceID)
		render.JSON(w, res, status)
		return
	}

	if err := model.ReleaseUser(rd.Id); err != nil {
		status = http.StatusInternalServerError
		errTitle := model.ErrCodeInternalError
		if err == model.ErrResourceNotFound {
			status = http.StatusNotFound
			errTitle = model.ErrCodeResourceNotFound
		}

		res.AddError(its(status), errTitle, err.Error(), logger.TraceID)
		logger.SetStatus(status).Info("param :", rd.Id, "response :", err.Error())
	}

	render.JSON(w, res, status)
}

func ForgotPassword(w http.ResponseWriter, r *http.Request) {
	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag("user_forgot")

	var rd PasswordReq
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		logger.SetStatus(http.StatusInternalServerError).Panic("param :", r.Body, "response :", err.Error())
	}

	password := r.FormValue("key")
	user, err := model.GetSession(password)
	if err != nil {
		logger.SetStatus(http.StatusInternalServerError).Info("param :", password, "response :", err.Error())
	}

	err = model.UpdatePassword(user.User.ID, hash(rd.Password))
	if err != nil {
		logger.SetStatus(http.StatusInternalServerError).Info("param :", rd.Password, "response :", err.Error())
	}
	render.JSON(w, http.StatusOK)
}
