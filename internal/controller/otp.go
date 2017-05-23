package controller

import (
	"bufio"
	"bytes"
	"fmt"
	"encoding/json"
	"net/http"
	"github.com/dghubble/sling"
	"github.com/gilkor/evoucher/internal/model"
)

type (
	AuthResponse struct {
		Data struct {
			Code        string `json:"code"`
			Description string `json:"description"`
			Name        string `json:"name"`
			State       string `json:"state"`
		} `json:"data"`
	}

	ReqParams struct {
		Key     	string `url:"key,omitempty"`
		Challenge      	string `url:"challenge,omitempty"`
		Password 	string `url:"password,omitempty"`
		Response     	string `url:"response,omitempty"`
	}
)

func OTPAuth(key, challenge , response string) bool {
	req := ReqParams{Key:key,Challenge: challenge, Response: response}
	d, r, err := ocra(req)

	if r.StatusCode == 200 && d.Data.State == "success" {
		return true
	} else if err != nil {
		return false
	}

	return false
}

func ocra ( param ReqParams) (AuthResponse,*http.Response, error){
	return server("/v1/ocra/ocra",param)
}

func server(path string, param ReqParams) (AuthResponse,*http.Response, error){
	c := http.Client{}
	s := sling.New()
	req , err :=s.Get(model.OCRA_URL + path).QueryStruct(param).Request()
	resp , err := c.Do(req)
	scanner := bufio.NewScanner(resp.Body)
	scanner.Split(bufio.ScanRunes)
	if err != nil {
		return AuthResponse{}, resp , err
	}
	var buf bytes.Buffer
	for scanner.Scan() {
		var n int
		n, err = buf.WriteString(scanner.Text()); if err != nil{
			fmt.Println("return  :", n ," : " ,err )
			return AuthResponse{}, resp , err
		}
	}

	fmt.Println("request URL  : " , resp.Request.URL)
	fmt.Println("request Header  : " , resp.Header)
	fmt.Println("request Body  : " , buf.String())
	fmt.Println("response status  : " , resp.StatusCode)

	//decode response data
	var data AuthResponse
	err = json.Unmarshal([]byte(buf.String()) , &data);if err != nil {
		fmt.Println("unmarshall error :" , err)
		return AuthResponse{}, resp , err
	}
	fmt.Println("response data : " , data)
	return data, resp , nil
}
