package controller

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gilkor/evoucher/internal/model"
)

func getUrlParam(url string) map[string]string {
	s := strings.Split(url, "?")
	param := strings.Split(s[1], "&")

	m := make(map[string]string)

	for _, v := range param {
		tempStr := strings.Split(v, "=")
		if tempStr[0] != "token" {
			m[tempStr[0]] = tempStr[1]
		}
	}

	return m
}

func sendPost(url string, param string) []byte {
	fmt.Println(param)
	var jsonStr = []byte(param)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
	return body
}

func getResponseData(param []byte) map[string]interface{} {
	var dat map[string]interface{}
	dat = make(map[string]interface{})
	if err := json.Unmarshal(param, &dat); err != nil {
		panic(err)
	}
	//fmt.Println(string(robots))

	if str, ok := dat["data"].(map[string]interface{}); ok {
		return str
	} else {
		return nil
	}
}

func hash(param string) string {
	password := []byte(param)
	hash := sha256.Sum256(password)
	return base64.StdEncoding.EncodeToString(hash[:])
}

// func basicAuth(w http.ResponseWriter, r *http.Request) (string, bool) {
// 	s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
// 	if len(s) != 2 {
// 		return "", false
// 	}
//
// 	b, err := base64.StdEncoding.DecodeString(s[1])
// 	if err != nil {
// 		return "", false
// 	}
//
// 	pair := strings.SplitN(string(b), ":", 2)
// 	if len(pair) != 2 {
// 		return "", false
// 	}
// 	fmt.Println(pair[0], hash(pair[1]))
// 	login, err := model.Login(pair[0], hash(pair[1]))
//
// 	if login == "" || err != nil {
// 		return "", false
// 	}
//
// 	return /*getAccountID(login)*/ "", true
// }

func randStr(ln int, fm string) string {
	CharsType := map[string]string{
		"Alphabet":     model.ALPHABET,
		"Numerals":     model.NUMERALS,
		"Alphanumeric": model.ALPHANUMERIC,
	}

	rand.Seed(time.Now().UTC().UnixNano())
	chars := CharsType[fm]
	result := make([]byte, ln)
	for i := 0; i < ln; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}
func its(i int) string {
	return strconv.Itoa(i)
}

func sti(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}
