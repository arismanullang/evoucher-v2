package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gilkor/evoucher-v2/model"
	u "github.com/gilkor/evoucher-v2/util"
)

type GetPrivilegesResponse struct {
	Data []string `json:"data"`
}

func CheckJWT(f http.Handler, access string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res := u.NewResponse()
		token := r.FormValue("token")

		if len(token) < 1 {
			res.SetError(JSONErrBadRequest.SetArgs("Token Not Found"))
			res.JSON(w, res, JSONErrBadRequest.Status)
			return
		}

		accData, err := model.GetSessionDataJWT(token)
		if err != nil {
			log.Panic(err)
			res.SetError(JSONErrUnauthorized)
			res.JSON(w, res, JSONErrUnauthorized.Status)
			return
		}

		urlPath := os.Getenv("JUNO_GET_ACCOUNT_PRIVILEGES_URL") + accData.AccountID + "?token=" + token
		httpRes, err := http.Get(urlPath)
		if err != nil {
			log.Panic(err)
			res.SetError(JSONErrUnauthorized)
			res.JSON(w, res, JSONErrUnauthorized.Status)
			return
		}

		if httpRes.StatusCode != 200 {
			log.Panic(err)
			res.SetError(JSONErrUnauthorized)
			res.JSON(w, res, JSONErrUnauthorized.Status)
			return
		}

		defer httpRes.Body.Close()
		var scope GetPrivilegesResponse
		json.NewDecoder(httpRes.Body).Decode(&scope)

		hasAccess := false
		for _, v := range scope.Data {
			fmt.Println(v)
			if v == access {
				hasAccess = true
				break
			}
		}

		if !hasAccess {
			res.SetError(JSONErrUnauthorized)
			res.JSON(w, res, JSONErrUnauthorized.Status)
			return
		}

		f.ServeHTTP(w, r)
	})

}

func CheckFuncJWT(f http.HandlerFunc, access string) http.Handler {
	return CheckJWT(http.HandlerFunc(f), access)
}
