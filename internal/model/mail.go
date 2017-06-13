package model

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/mailgun/mailgun-go.v1"
)

var (
	Domain       string
	ApiKey       string
	PublicApiKey string
	RootTemplate string
)

func SendMail(domain, apiKey, publicApiKey, username string) error {
	id, err := CheckUsername(username)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	user, err := FindUserDetail(id)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Println(user)
	mg := mailgun.NewMailgun(domain, apiKey, publicApiKey)
	message := mailgun.NewMessage(
		"evoucher@gilkor.com",
		"Forgot Password E-Voucher",
		makeMessage(id),
		user.Email)
	resp, id, err := mg.Send(message)
	if err != nil {
		return err
	}
	fmt.Printf("ID: %s Resp: %s\n", id, resp)
	return nil
}

func makeMessage(id string) string {
	u, err := FindUserDetail(id)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	tok := GenerateToken(u)
	str, err := ioutil.ReadFile(RootTemplate + "template")
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}

	url := "http://voucher.apps.id:8889/user/recover?key=" + tok.Token
	result := string(str) + url
	return result
}
