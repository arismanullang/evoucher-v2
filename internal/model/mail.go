package model

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/mailgun/mailgun-go.v1"
	"strings"
)

var (
	Domain       string
	ApiKey       string
	PublicApiKey string
	RootTemplate string
	RootUrl      string
	Email        string
)

type (
	SedayuOneEmail struct {
		Name       string
		VoucherUrl string
	}
)

func SendMailForgotPassword(domain, apiKey, publicApiKey, username, accountId string) error {
	id, err := CheckUsername(username, accountId)
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
		Email,
		"Forgot Password E-Voucher",
		makeMessageForgotPassword(id),
		user.Email)
	resp, id, err := mg.Send(message)
	if err != nil {
		return err
	}
	fmt.Printf("ID: %s Resp: %s\n", id, resp)
	return nil
}

func makeMessageForgotPassword(id string) string {
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

	url := "https://" + RootUrl + "/user/recover?key=" + tok.Token
	//element := "<a href='"+url+"'>"+url+"</a>"
	result := string(str) + url
	return result
}

func SendMailSedayuOne(domain, apiKey, publicApiKey, subject string, emailTarget []string, param []SedayuOneEmail) error {
	mg := mailgun.NewMailgun(domain, apiKey, publicApiKey)

	for i, v := range emailTarget {
		message := mailgun.NewMessage(
			Email,
			subject,
			subject,
			v)
		message.SetHtml(makeMessageEmailSedayuOne(param[i]))
		resp, id, err := mg.Send(message)
		if err != nil {
			return err
		}
		fmt.Printf("ID: %s Resp: %s\n", id, resp)
	}

	return nil
}

func makeMessageEmailSedayuOne(param SedayuOneEmail) string {
	// %%full-name%%
	// %%link-voucher%%
	str, err := ioutil.ReadFile(RootTemplate + "sedayu_one")
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}

	result := string(str)
	result = strings.Replace(result, "%%full-name%%", param.Name, 1)
	result = strings.Replace(result, "%%link-voucher%%", param.VoucherUrl, 1)
	return result
}
