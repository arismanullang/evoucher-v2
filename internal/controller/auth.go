package controller

import (
	"net/http"
	"time"

	"github.com/gilkor/evoucher/internal/model"
	"github.com/ruizu/render"
)

func CheckToken(w http.ResponseWriter, r *http.Request) (ResponseData, bool) {
	var rs ResponseData
	token := r.FormValue("token")
	if len(token) < 1 {
		rs.State = its(http.StatusUnauthorized)
		rs.Error = http.StatusText(http.StatusUnauthorized)
		rs.Message = model.ErrMessageTokenNotFound
		return rs, false
	}
	if _, _, exp, err := checkExpired(r, token); err != nil {
		rs.State = its(http.StatusUnauthorized)
		rs.Error = http.StatusText(http.StatusUnauthorized)
		rs.Message = model.ErrMessageTokenNotFound + "(" + err.Error() + ")"
		return rs, false
	} else if exp.Before(time.Now()) {
		rs.State = its(http.StatusUnauthorized)
		rs.Error = http.StatusText(http.StatusUnauthorized)
		rs.Message = model.ErrMessageTokenExpired
		return rs, false
	}
	return rs, true
}

func check(f http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if rs, ok := CheckToken(w, r); !ok {
			res := NewResponse(rs)
			render.JSON(w, res, sti(rs.State))
			return
		}
		f.ServeHTTP(w, r)
	})
}

func CheckTokenAuth(f http.HandlerFunc) http.Handler {
	return check(http.HandlerFunc(f))
}
