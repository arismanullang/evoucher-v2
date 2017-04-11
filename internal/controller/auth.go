package controller

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gilkor/evoucher/internal/model"
	"github.com/ruizu/render"
)

func basicAuth(w http.ResponseWriter, r *http.Request) (string, string, bool) {
	s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	if len(s) != 2 {
		return "", "", false
	}

	b, err := base64.StdEncoding.DecodeString(s[1])
	if err != nil {
		return "", "", false
	}

	pair := strings.SplitN(string(b), ":", 2)
	if len(pair) != 2 {
		return "", "", false
	}

	login, err := model.Login(pair[0], hash(pair[1]))
	if login == "" || err != nil {
		return "", "", false
	}
	fmt.Println(login)

	ac, err := model.GetAccountsByUser(login)
	if err != nil {
		return "", "", false
	}
	fmt.Println("a")

	return ac[0], login, true
}

func AuthToken(w http.ResponseWriter, r *http.Request) (string, string, time.Time, bool) {
	res := NewResponse(nil)
	token := r.FormValue("token")

	if len(token) < 1 {
		res.AddError(its(http.StatusUnauthorized), model.ErrCodeMissingToken, model.ErrMessageTokenNotFound, "token")
		render.JSON(w, res, http.StatusUnauthorized)
		return "", "", time.Now(), false
	}
	// Return : SessionData{ user_id, account_id, expired} , error
	s, err := model.GetSession(token)
	if err != nil {
		switch err {
		case model.ErrTokenNotFound:
			res.AddError(its(http.StatusUnauthorized), model.ErrCodeInvalidToken, model.ErrMessageTokenNotFound, "token")
			render.JSON(w, res, http.StatusUnauthorized)
		case model.ErrTokenExpired:
			res.AddError(its(http.StatusUnauthorized), model.ErrCodeInvalidToken, model.ErrMessageTokenExpired, "token")
			render.JSON(w, res, http.StatusUnauthorized)
		}
		return "", "", time.Now(), false
	}

	return s.AccountID, s.UserId, s.ExpiredAt, true
	//return "NNJs3Nfo", "IEC1cL77", time.Now().Add(time.Duration(model.TOKENLIFE) * time.Minute), true
}

func GetToken(w http.ResponseWriter, r *http.Request) {
	res := NewResponse(nil)
	status := http.StatusOK

	ac, ui, ok := basicAuth(w, r)
	if !ok {
		status = http.StatusUnauthorized
		res.AddError(its(http.StatusUnauthorized), model.ErrCodeInvalidUser, model.ErrMessageInvalidUser, "token")
		render.JSON(w, res, status)
		return
	}

	d := model.GenerateToken(ac, ui)
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
