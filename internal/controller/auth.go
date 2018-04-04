package controller

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/gilkor/evoucher/internal/model"
	"github.com/ruizu/render"
)

type (
	Auth struct {
		User  model.User
		res   *Response
		Valid bool
	}
	LoginResponse struct {
		Token       model.Token  `json:"token"`
		Role        []model.Role `json:"role"`
		Ui          []string     `json:"ui"`
		Api         []string     `json:"api"`
		Destination string       `json:"destination"`
	}
	Check struct {
		User        model.User `json:"user"`
		Destination string     `json:"destination"`
	}
)

func basicAuth(r *http.Request) (model.User, bool) {
	logger := model.NewLog()
	logger.SetService("AUTH").SetMethod(r.Method).SetTag("Basic-Authentication")

	s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	if len(s) != 2 {
		logger.SetStatus(http.StatusUnauthorized).Log("User :" + s[1] + " # response : Authentication failed (missing username / password)")
		return model.User{}, false
	}

	b, err := base64.StdEncoding.DecodeString(s[1])
	if err != nil {
		logger.SetStatus(http.StatusUnauthorized).Log("User :" + s[1] + " # response : Authentication failed (error while decoding password)")
		return model.User{}, false
	}

	pair := strings.SplitN(string(b), ":", 2)
	if len(pair) != 2 {
		logger.SetStatus(http.StatusUnauthorized).Log("User :" + s[1] + " # response : Authentication failed (error while decoding password)")
		return model.User{}, false
	}

	login, err := model.Login(pair[0], hash(pair[1]))
	if login == "" || err != nil {
		logger.SetStatus(http.StatusUnauthorized).Log("User :" + s[1] + " # response : Authentication failed (user not found)")
		return model.User{}, false
	}

	user, err := model.FindUserDetail(login)
	if user.Username == "" || err != nil {
		logger.SetStatus(http.StatusUnauthorized).Log("User :" + s[1] + " # response : Authentication failed (user not found)")
		return model.User{}, false
	}

	ac, err := model.GetAccountDetailByUser(login)
	if err != nil {
		logger.SetStatus(http.StatusUnauthorized).Log("User :" + s[1] + " # response : Authentication failed (invalid user account)")
		return model.User{}, false
	}

	logger.SetStatus(http.StatusUnauthorized).Log("User :" + s[1] + " # response : success")
	user.Account = ac

	return user, true
}

func AuthToken(w http.ResponseWriter, r *http.Request) Auth {
	logger := model.NewLog()
	logger.SetService("AUTH").SetMethod(r.Method).SetTag("Client-Authentication")
	res := NewResponse(nil)
	token := r.FormValue("token")

	if len(token) < 1 {
		res.AddError(its(http.StatusUnauthorized), model.ErrCodeMissingToken, model.ErrMessageTokenNotFound, logger.TraceID)
		logger.SetStatus(http.StatusUnauthorized).Log("token :"+token+" # response :", res.Errors)
		return Auth{User: model.User{}, res: res, Valid: false}
	}
	// Return : SessionData{ user_id, account_id, expired} , error
	s, err := model.GetSession(token)
	if err != nil {
		switch err {
		case model.ErrTokenNotFound:
			res.AddError(its(http.StatusUnauthorized), model.ErrCodeInvalidToken, model.ErrMessageTokenNotFound, logger.TraceID)
		case model.ErrTokenExpired:
			res.AddError(its(http.StatusUnauthorized), model.ErrCodeInvalidToken, model.ErrMessageTokenExpired, logger.TraceID)
		}
		logger.SetStatus(http.StatusUnauthorized).Log("token :"+token+" # response :", res.Errors)
		return Auth{User: model.User{}, res: res, Valid: false}
	} else if !model.IsExistToken(token) {
		res.AddError(its(http.StatusUnauthorized), model.ErrCodeInvalidToken, model.ErrMessageTokenExpired, logger.TraceID)
		logger.SetStatus(http.StatusUnauthorized).Log("token :"+token+" # response :", res.Errors)
		//render.JSON(w, res, http.StatusUnauthorized)
		return Auth{User: model.User{}, res: res, Valid: false}
	}

	logger.SetStatus(http.StatusUnauthorized).Log("token :" + token + " # response : Authentication success")
	return Auth{User: s.User, res: res, Valid: true}
	//return "NNJs3Nfo", "IEC1cL77", time.Now().Add(time.Duration(model.TOKENLIFE) * time.Minute), true
}

func AuthTokenWithLogger(w http.ResponseWriter, r *http.Request, logger *model.LogField) Auth {
	res := NewResponse(nil)
	token := r.FormValue("token")

	if len(token) < 1 {
		res.AddError(its(http.StatusUnauthorized), model.ErrCodeMissingToken, model.ErrMessageTokenNotFound, logger.TraceID)
		logger.SetStatus(http.StatusUnauthorized).Log("token :"+token+" # response :", res.Errors)
		return Auth{User: model.User{}, res: res, Valid: false}
	}
	// Return : SessionData{ user_id, account_id, expired} , error
	s, err := model.GetSession(token)
	if err != nil {
		switch err {
		case model.ErrTokenNotFound:
			res.AddError(its(http.StatusUnauthorized), model.ErrCodeInvalidToken, model.ErrMessageTokenNotFound, logger.TraceID)
		case model.ErrTokenExpired:
			res.AddError(its(http.StatusUnauthorized), model.ErrCodeInvalidToken, model.ErrMessageTokenExpired, logger.TraceID)
		}
		logger.SetStatus(http.StatusUnauthorized).Log("token :"+token+" # response :", res.Errors)
		return Auth{User: model.User{}, res: res, Valid: false}
	} else if !model.IsExistToken(token) {
		res.AddError(its(http.StatusUnauthorized), model.ErrCodeInvalidToken, model.ErrMessageTokenExpired, logger.TraceID)
		logger.SetStatus(http.StatusUnauthorized).Log("token :"+token+" # response :", res.Errors)
		//render.JSON(w, res, http.StatusUnauthorized)
		return Auth{User: model.User{}, res: res, Valid: false}
	}

	logger.SetStatus(http.StatusUnauthorized).Log("token :" + token + " # response : Authentication success")
	return Auth{User: s.User, res: res, Valid: true}
	//return "NNJs3Nfo", "IEC1cL77", time.Now().Add(time.Duration(model.TOKENLIFE) * time.Minute), true
}

func GetToken(w http.ResponseWriter, r *http.Request) {
	res := NewResponse(nil)
	status := http.StatusOK

	u, ok := basicAuth(r)
	if !ok {
		status = http.StatusUnauthorized
		res.AddError(its(http.StatusUnauthorized), model.ErrCodeInvalidUser, model.ErrMessageInvalidUser, "token")
		render.JSON(w, res, status)
		return
	}

	d := model.GenerateToken(u)
	res = NewResponse(d)
	render.JSON(w, res, status)
}

func Login(w http.ResponseWriter, r *http.Request) {
	res := NewResponse(nil)
	status := http.StatusOK

	u, ok := basicAuth(r)
	if !ok {
		status = http.StatusUnauthorized
		res.AddError(its(http.StatusUnauthorized), model.ErrCodeInvalidUser, model.ErrMessageInvalidUser, "token")
		render.JSON(w, res, status)
		return
	}

	d := model.GenerateToken(u)

	uiFeatures := []string{}
	ui, err := model.GetUiFeatures(u.Role[0].Id)
	if err != nil {
	}

	for _, value := range ui {
		uiFeatures = append(uiFeatures, "/"+value.Category+"/"+value.Detail)
	}

	dest := "/program/index"
	if u.Role[0].Id == "Mn78I1wc" {
		dest = "/sa/search"
	}
	login := LoginResponse{
		Token:       d,
		Role:        u.Role,
		Ui:          uiFeatures,
		Destination: dest,
	}
	res = NewResponse(login)
	render.JSON(w, res, status)
}

func CheckToken(w http.ResponseWriter, r *http.Request) {
	res := NewResponse(nil)
	status := http.StatusOK

	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag("Check Token")

	a := AuthTokenWithLogger(w, r, logger)
	if !a.Valid {
		res = a.res
		status = http.StatusUnauthorized
		render.JSON(w, res, status)
		return
	}

	res = NewResponse(true)
	render.JSON(w, res, status)
}

func UICheckToken(w http.ResponseWriter, r *http.Request) {
	res := NewResponse(nil)
	url := r.FormValue("url")
	status := http.StatusOK

	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag("Check Token")

	valid := false
	a := AuthTokenWithLogger(w, r, logger)
	if !a.Valid {
		res = a.res
		status = http.StatusUnauthorized
		render.JSON(w, res, status)
		return
	}

	if url == "/" {
		valid = true
	}
	for _, valueRole := range a.User.Role {
		features, err := model.GetUiFeatures(valueRole.Id)
		if err != nil {
			valid = false
		}
		for _, valueFeature := range features {
			tempFeature := "/" + valueFeature.Category + "/" + valueFeature.Detail
			if url == tempFeature {
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

	dest := "/program/index"
	if a.User.Role[0].Detail == "sa" {
		dest = "/sa/search"
	}
	result := Check{
		User:        a.User,
		Destination: dest,
	}

	res = NewResponse(result)
	render.JSON(w, res, status)
}

func CheckAPIRole(a Auth, apiName string) bool {
	for _, valueRole := range a.User.Role {
		features, err := model.GetApiFeatures(valueRole.Id)
		if err != nil {
			return true
		}
		for _, valueFeature := range features {
			tempFeature := valueFeature.Category + "_" + valueFeature.Detail
			if apiName == tempFeature {
				return false
			}
		}
	}

	return true
}
