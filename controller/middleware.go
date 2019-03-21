package controller

import (
	"net/http"

	"github.com/gilkor/evoucher/model"
	u "github.com/gilkor/evoucher/util"

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
