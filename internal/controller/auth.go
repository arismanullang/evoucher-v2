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
	fmt.Println(pair[0], hash(pair[1]))
	login, err := model.Login(pair[0], hash(pair[1]), "So-GAf-G")

	if login == "" || err != nil {
		return "", "", false
	}

	return "So-GAf-G", login, true
}

func CheckToken(w http.ResponseWriter, r *http.Request) (string, string, time.Time, bool) {
	res := NewResponse(nil)
	token := r.FormValue("token")
	var exp time.Time
	var valid bool
	var err error
	var a, u string

	if len(token) < 1 {
		res.AddError(its(http.StatusUnauthorized), model.ErrCodeMissingToken, model.ErrMessageTokenNotFound, "token")
		render.JSON(w, res, http.StatusUnauthorized)
		return "", "", time.Now(), false
	}
	// Return : user_id, account_id, expired, boolean, error
	if u, a, exp, valid, err = getValiditySession(r, token); err != nil {
		switch err {
		case model.ErrTokenNotFound:
			res.AddError(its(http.StatusUnauthorized), model.ErrCodeInvalidToken, model.ErrMessageTokenNotFound, "token")
			render.JSON(w, res, http.StatusUnauthorized)
		case model.ErrTokenExpired:
			res.AddError(its(http.StatusUnauthorized), model.ErrCodeInvalidToken, model.ErrMessageTokenExpired, "token")
			render.JSON(w, res, http.StatusUnauthorized)
		}
		return "", "", time.Now(), false
	} else if !valid {
		res.AddError(its(http.StatusUnauthorized), model.ErrCodeInvalidToken, model.ErrMessageTokenNotFound, "token")
		render.JSON(w, res, http.StatusUnauthorized)
		return "", "", time.Now(), false
	}
	return a, u, exp, true
}

func GetToken(w http.ResponseWriter, r *http.Request) {
	res := NewResponse(nil)

	ac, ui, ok := basicAuth(w, r)
	if !ok {
		res.AddError(its(http.StatusUnauthorized), model.ErrCodeMissingToken, model.ErrMessageTokenNotFound, "token")
		render.JSON(w, res, http.StatusUnauthorized)
		return
	}

	d := model.GenerateToken(ac, ui)

	res = NewResponse(d)
	render.JSON(w, res, http.StatusUnauthorized)
}
