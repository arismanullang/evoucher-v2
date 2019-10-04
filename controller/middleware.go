package controller

import (
	"context"
	"net/http"

	"github.com/gilkor/evoucher-v2/model"
	u "github.com/gilkor/evoucher-v2/util"

	"github.com/go-zoo/bone"

	"github.com/urfave/negroni"
)

// CompanyID : global companyID
var CompanyID *string

//CompanyParamMiddleware : global middleware to get company id from alias inside URL param
func CompanyParamMiddleware() negroni.Handler {
	return negroni.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
			alias := bone.GetValue(r, "company")
			c, _, err := model.GetCompanyByAlias(alias, u.NewQueryParam(r))
			if err != nil {
				if err == model.ErrorResourceNotFound {
					u.NewResponse().SetError(JSONErrResourceNotFound)
					return
				}
				u.NewResponse().SetError(JSONErrFatal.SetArgs(err.Error()))
				return
			}
			CompanyID = &c[0].ID

			//next Handler
			next(w, r)

		})
}

//VerifyJunoJWTAuthMiddleware : set session user data from auth
func VerifyJunoJWTAuthMiddleware() negroni.Handler {
	return negroni.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
			tokenString := bone.GetValue(r, "token")
			token, err := VerifyJWT(tokenString)
			if err != nil {
				if err == model.ErrorForbidden {
					u.NewResponse().SetError(JSONErrForbidden)
					return
				} else if err == model.ErrorInternalServer {
					u.NewResponse().SetError(JSONErrFatal)
					return
				}
			} else {
				claims, ok := token.Claims.(*JWTJunoClaims)
				if ok && token.Valid {
					// fmt.Printf("Key:%v", token.Header)
					ctx := context.WithValue(r.Context(), KeyContextAuth, claims)
					r.WithContext(ctx)
				} else {
					u.NewResponse().SetError(JSONErrForbidden)
					return
				}
			}
			next(w, r)
		})
}

type (
	contextKey string
)

const (
	//KeyContextAuth :
	KeyContextAuth = contextKey("auth")
	//KeyContextCompanyID :
	KeyContextCompanyID = contextKey("company_id")
	//KeyContextUserID :
	KeyContextUserID = contextKey("user_id")
	//KeyContextAccountID :
	KeyContextAccountID = contextKey("account_id")
)
