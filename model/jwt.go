package model

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gilkor/evoucher-v2/util"
)

// "account_id": "S51uSQ6Z",
//   "username": "6281230335880",
//   "name": "Aris Manullang",
//   "mobile_calling_code": "62",
//   "mobile_no": "81230335880",
//   "email": "lamhot@gilkor.com",
//   "gender": "male",
//   "identity_type": "ktp",
//   "identity_no": "5555333322221111",
//   "client_key": "ESrAB66amXanqiVgkXXX5yy7wKInnGWC5",
//   "account_verifications": [
//     {
//   ....
//     }
//   ],
//   "aud": "system",
//   "exp": 1578554513,
//   "iat": 1578468113,
//   "iss": "http://juno-staging.elys.id"
// }
type (
	//JWTJunoClaims Juno format of claims
	JWTJunoClaims struct {
		Algorithm string `json:"alg"`
		KeyID     string `json:"kid"`
		Type      string `json:"typ"`
		// SessionData SessionData `json:"session_data"`
		CompanyID        string `json:"aud"`
		AccountID        string `json:"account_id"`
		Username         string `json:"username"`
		Name             string `json:"name"`
		MobileCalingCode string `json:"mobile_caling_code"`
		MobileNo         string `json:"mobile_no"`
		Email            string `json:"email"`
		Gender           string `json:"gender"`
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
			util.DEBUG(err)
			return nil, fmt.Errorf("[API] Something went wrong. Error:%v", "Public Key")
		}

		if res.StatusCode != 200 {
			return nil, fmt.Errorf("[API] Unexpected status code:%v", res.StatusCode)
		}

		// Read the token out of the response body
		buf := new(bytes.Buffer)
		io.Copy(buf, res.Body)
		res.Body.Close()

		key, err := jwt.ParseRSAPublicKeyFromPEM(buf.Bytes())
		if err != nil {
			return nil, fmt.Errorf("failed to load key:%s", err)
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
