package controller

import (
	"net/http"

	"github.com/gilkor/evoucher/internal/model"
	"github.com/ruizu/render"
)

func CheckToken(w http.ResponseWriter, r *http.Request) bool {
	res := NewResponse(nil)
	token := r.FormValue("token")
	if len(token) < 1 {
		res.AddError("401001", model.ErrCodeMissingToken, model.ErrMessageTokenNotFound, "token")
		render.JSON(w, res, http.StatusUnauthorized)
		return false
	}
	if _, _, _, valid, err := getValiditySession(r, token); err != nil {
		res.AddError("401002", model.ErrCodeMissingToken, model.ErrMessageTokenNotFound+"("+err.Error()+")", "token")
		render.JSON(w, res, http.StatusUnauthorized)
		return false
	} else if !valid {
		res.AddError("401003", model.ErrCodeMissingToken, model.ErrMessageTokenNotFound, "token")
		render.JSON(w, res, http.StatusUnauthorized)
		return false
	}
	return true
}

func check(f http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !CheckToken(w, r) {
			return
		}
		f.ServeHTTP(w, r)
	})
}

func CheckTokenAuth(f http.HandlerFunc) http.Handler {
	return check(http.HandlerFunc(f))
}
