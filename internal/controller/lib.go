package controller

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"regexp"
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
	return body
}

func getResponseData(param []byte) map[string]interface{} {
	var dat map[string]interface{}
	dat = make(map[string]interface{})
	if err := json.Unmarshal(param, &dat); err != nil {
		panic(err)
	}

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

func replaceSpecialCharacter(param string) string {
	reg, err := regexp.Compile("[^A-Za-z0-9]")
	if err != nil {
		log.Fatal(err)
	}

	safe := reg.ReplaceAllString(param, "x")
	safe = strings.Trim(safe, "-")
	return safe
}

func randomize(param int) int {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	return r1.Intn(param)
}

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
func bts(b bool) string {
	return strconv.FormatBool(b)
}

func stringInSlice(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}

func StrEncode(s string) string {
	base64.StdEncoding.DecodedLen(32)
	return base64.StdEncoding.EncodeToString([]byte(s))
}

func StrDecode(s string) string {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		log.Panic(err)
	}
	return string(data)
}

func stf(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}
