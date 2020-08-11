package controller

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/gilkor/evoucher-v2/model"
	u "github.com/gilkor/evoucher-v2/util"
	"github.com/go-zoo/bone"
)

type GetPrivilegesResponse struct {
	Data []string `json:"data"`
}

func CheckJWT(f http.Handler, access string) http.Handler {
	return f
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res := u.NewResponse()
		companyId := bone.GetValue(r, "company")
		token := r.FormValue("token")

		if len(token) < 1 {
			res.SetError(JSONErrBadRequest.SetArgs("Token Not Found"))
			res.JSON(w, res, JSONErrBadRequest.Status)
			return
		}

		accData, err := model.GetSessionDataJWT(token, companyId)
		if err != nil {
			res.SetError(JSONErrUnauthorized)
			res.JSON(w, res, JSONErrUnauthorized.Status)
			return
		}

		configs, err := model.GetConfigs(companyId, "juno")
		if len(configs) > 0 && accData.ClientKey != configs["client_key"].(string) {
			token, err = GetVoucherToken(token, accData.AccountID, configs["client_key"].(string))
			if err != nil {
				res.SetError(JSONErrUnauthorized)
				res.JSON(w, res, JSONErrUnauthorized.Status)
				return
			}
		}

		urlPath := os.Getenv("JUNO_GET_ACCOUNT_PRIVILEGES_URL") + accData.AccountID + "?token=" + token
		httpRes, err := http.Get(urlPath)
		if err != nil {
			res.SetError(JSONErrUnauthorized)
			res.JSON(w, res, JSONErrUnauthorized.Status)
			return
		}

		if httpRes.StatusCode != 200 {
			res.SetError(JSONErrUnauthorized)
			res.JSON(w, res, JSONErrUnauthorized.Status)
			return
		}

		defer httpRes.Body.Close()
		var scope GetPrivilegesResponse
		json.NewDecoder(httpRes.Body).Decode(&scope)

		hasAccess := false
		for _, v := range scope.Data {
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

type JunoToken struct {
	Token string `json:"token"`
}

type JunoRes struct {
	Data JunoToken `json:"data"`
}

func GetVoucherToken(token, account_id, clientKey string) (string, error) {
	urlPath := os.Getenv("JUNO_REGISTER_CLIENT_URL") + clientKey + "?token=" + token
	httpRes, err := http.Get(urlPath)
	if err != nil {
		return "", err
	}

	if httpRes.StatusCode != 200 {
		urlPath = os.Getenv("JUNO_GET_TOKEN_URL") + clientKey + "?token=" + token
		httpRes, err = http.Get(urlPath)
		if err != nil {
			return "", err
		}
	}

	defer httpRes.Body.Close()
	var res JunoRes
	json.NewDecoder(httpRes.Body).Decode(&token)
	return res.Data.Token, nil
}

type JunoRoleResponse struct {
	ClientSecret string   `json:"client_secret"`
	Scope        []string `json:"scope"`
}

func GetJunoaBasicRole(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	companyID := bone.GetValue(r, "company")

	configs, err := model.GetConfigs(companyID, "juno")
	if err != nil {
		res.SetError(JSONErrFatal.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	scope := strings.Split(configs["basic_roles"].(string), ",")

	roleRes := JunoRoleResponse{
		ClientSecret: configs["client_secret"].(string),
		Scope:        scope,
	}
	// res.SetPagination(r, qp.Page, next)
	res.JSON(w, roleRes, http.StatusOK)
}
