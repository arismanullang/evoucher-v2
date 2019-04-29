package controller

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gilkor/evoucher/model"
)

type (
	//JWTJunoClaims Juno format of claims
	JWTJunoClaims struct {
		Algorithm string `json:"alg"`
		KeyID     string `json:"kid"`
		Type      string `json:"typ"`
		// SessionData SessionData `json:"session_data"`
		// CompanyID        string `json:"aud"`
		AccountID        string `json:"account_id"`
		Username         string `json:"username"`
		Name             string `json:"name"`
		MobileCalingCode string `json:"mobile_caling_code"`
		MobileNo         string `json:"mobile_no"`
		Email            string `json:"email"`
		Gender           string `json:"gender"`
		IdentityType     string `json:"identity_type"`
		IdentityNo       string `json:"identity_no"`
		ClientKey        string `json:"client_key"`
		jwt.StandardClaims
	}
	//SessionData token auth
	SessionData struct {
		AccountID        string `json:"account_id"`
		Username         string `json:"username"`
		Name             string `json:"name"`
		MobileCalingCode string `json:"mobile_caling_code"`
		MobileNo         string `json:"mobile_no"`
		Email            string `json:"email"`
		Gender           string `json:"gender"`
	}
)

//VerifyJWT Auth JWT JUNO
func VerifyJWT(tokenString string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenString, &JWTJunoClaims{}, func(token *jwt.Token) (interface{}, error) {
		urlFormat := os.Getenv("JUNO_PUBLIC_KEY_URL")
		urlPath := fmt.Sprintf(urlFormat, token.Header["kid"])
		res, err := http.Get(urlPath)
		if err != nil {
			// e(err)
			return nil, model.ErrorForbidden
		}

		if res.StatusCode != 200 {
			return nil, model.ErrorForbidden
		}

		// Read the token out of the response body
		buf := new(bytes.Buffer)
		io.Copy(buf, res.Body)
		res.Body.Close()

		key, err := jwt.ParseRSAPublicKeyFromPEM(buf.Bytes())
		if err != nil {
			return nil, err
		}

		parts := strings.Split(tokenString, ".")
		keyMethod := (token.Header["alg"]).(string)
		method := jwt.GetSigningMethod(keyMethod)
		err = method.Verify(strings.Join(parts[0:2], "."), parts[2], key)
		if err != nil {
			return nil, err
		}

		return key, nil
	})
}
