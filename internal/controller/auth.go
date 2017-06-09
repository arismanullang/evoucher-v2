package controller

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/gilkor/evoucher/internal/model"
	"github.com/ruizu/render"
)

func basicAuth(r *http.Request) (model.User, bool) {
	s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	if len(s) != 2 {
		return model.User{},false
	}

	b, err := base64.StdEncoding.DecodeString(s[1])
	if err != nil {
		return model.User{},false
	}

	pair := strings.SplitN(string(b), ":", 2)
	if len(pair) != 2 {
		return model.User{},false
	}

	login, err := model.Login(pair[0], hash(pair[1]))
	if login == "" || err != nil {
		return model.User{},false
	}

	user, err := model.FindUserDetail(login)
	if user.Username == "" || err != nil {
		return model.User{},false
	}

	ac, err := model.GetAccountsByUser(login)
	if err != nil {
		return model.User{},false
	}
	user.AccountID = ac[0]



	return model.User{}, true
}

type (
	Auth struct {
		User 	model.User
		res  	*Response
		Valid 	bool
	}
)

func AuthToken(w http.ResponseWriter, r *http.Request) (Auth) {
	res := NewResponse(nil)
	token := r.FormValue("token")

	if len(token) < 1 {
		res.AddError(its(http.StatusUnauthorized), model.ErrCodeMissingToken, model.ErrMessageTokenNotFound, "token")
		return Auth{User: model.User{}, res: res, Valid: false,}
	}
	// Return : SessionData{ user_id, account_id, expired} , error
	s, err := model.GetSession(token)
	if err != nil {
		switch err {
		case model.ErrTokenNotFound:
			res.AddError(its(http.StatusUnauthorized), model.ErrCodeInvalidToken, model.ErrMessageTokenNotFound, "token")
		case model.ErrTokenExpired:
			res.AddError(its(http.StatusUnauthorized), model.ErrCodeInvalidToken, model.ErrMessageTokenExpired, "token")
		}
		return Auth{User: model.User{}, res: res, Valid: false,}
	} else if !model.IsExistToken(token) {
		res.AddError(its(http.StatusUnauthorized), model.ErrCodeInvalidToken, model.ErrMessageTokenExpired, "token")
		return Auth{User: model.User{}, res: res, Valid: false,}
	}

	return Auth{User: s.User, res: res, Valid: true,}
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

func CheckToken(w http.ResponseWriter, r *http.Request) {
	res := NewResponse(nil)
	token := r.FormValue("token")
	status := http.StatusOK

	Exists := true
	if !model.IsExistToken(token) {
		Exists = false
		status = http.StatusUnauthorized
	}

	res = NewResponse(Exists)
	render.JSON(w, res, status)
}
