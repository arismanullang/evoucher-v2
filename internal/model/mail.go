package model

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"strings"
	"time"

	"gopkg.in/mailgun/mailgun-go.v1"
)

var (
	Config       map[string]map[string]string
	Domain       string
	ApiKey       string
	PublicApiKey string
	RootTemplate string
	RootURL      string
	Email        string
)

type (
	CampaignData struct {
		HolderEmail  string  `json:"email"`
		HolderName   string  `json:"name"`
		VoucherUrl   string  `json:"voucher_url"`
		VoucherObj   Voucher `json:"voucher"`
		ProgramName  string  `json:"program_name"`
		ImageHeader  string  `json:"image_header"`
		ImageVoucher string  `json:"image_voucher"`
		ImageFooter  string  `json:"image_footer"`
		EmailContent string  `json:"email_content"`
		EmailSubject string  `json:"email_subject"`
	}

	TargetEmail struct {
		EmailAddress string      `json:"email_address"`
		Data         interface{} `json:"data"`
	}

	BaseMail struct {
		From     string        `json:"from"`
		To       []TargetEmail `json:"to"`
		Subject  string        `json:"subject,omitempty"`
		Message  string        `json:"message,omitempty"`
		Template string        `json:"template,omitempty"`
	}

	ProgramCampaign struct {
		Id           string    `db:"id" json:"id"`
		ProgramID    string    `db:"program_id" json:"program_id"`
		ProgramName  string    `db:"program_name" json:"program_name"`
		AccountID    string    `db:"account_id" json:"account_id"`
		ImageHeader  string    `db:"header_image" json:"header_image"`
		ImageVoucher string    `db:"voucher_image" json:"voucher_image"`
		ImageFooter  string    `db:"footer_image" json:"footer_image"`
		CreatedBy    string    `db:"created_by" json:"created_by"`
		CreatedAt    time.Time `db:"created_at" json:"created_at"`
	}

	ProgramCampaignV2 struct {
		ID           string    `db:"id" json:"id"`
		ProgramID    string    `db:"program_id" json:"program_id"`
		ProgramName  string    `db:"program_name" json:"program_name"`
		AccountID    string    `db:"account_id" json:"account_id"`
		Template     string    `db:"email_template" json:"email_template"`
		EmailSubject string    `db:"email_subject" json:"email_subject"`
		EmailSender  string    `db:"email_sender" json:"email_sender"`
		EmailContent string    `db:"email_content" json:"email_content"`
		ImageHeader  string    `db:"header_image" json:"header_image"`
		ImageVoucher string    `db:"voucher_image" json:"voucher_image"`
		ImageFooter  string    `db:"footer_image" json:"footer_image"`
		CreatedBy    string    `db:"created_by" json:"created_by"`
		CreatedAt    time.Time `db:"created_at" json:"created_at"`
	}

	ConfirmationEmail struct {
		EmailAccount string `db:"account"`
		EmailPartner string `db:"partner"`
		EmailMember  string `db:"member"`
	}

	ConfirmationEmailRequest struct {
		Holder          string    `json:"name"`
		ProgramName     string    `json:"program_name"`
		TransactionCode string    `json:"transaction_code"`
		TransactionDate time.Time `json:"transaction_date"`
		PartnerName     string    `json:"partner_name"`
		ListEmail       []string  `json:"list_email"`
		ListVoucher     []string  `json:"list_voucher"`
		EmailSubject    string    `json:"email_subject"`
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

	url := "https://" + RootURL + "/user/recover?key=" + tok.Token
	//element := "<a href='"+url+"'>"+url+"</a>"
	result := string(str) + url
	return result
}

func SendMailSedayuOne(domain, apiKey, publicApiKey, subject string, target []CampaignData, program ProgramCampaign) error {
	mg := mailgun.NewMailgun(domain, apiKey, publicApiKey)

	for _, v := range target {
		message := mailgun.NewMessage(
			Email,
			subject,
			subject,
			v.HolderEmail)
		message.SetHtml(makeMessageEmailSedayuOne(program, v))
		resp, id, err := mg.Send(message)
		if err != nil {
			return err
		}
		UpdateBroadcastUserState(v.HolderEmail, program.ProgramID, program.CreatedBy)
		fmt.Printf("ID: %s Resp: %s\n", id, resp)
	}

	return nil
}
func makeMessageEmailSedayuOne(program ProgramCampaign, target CampaignData) string {
	// %%full-name%%
	// %%link-voucher%%
	templateCampaign := Config[program.AccountID]["email_campaign"]
	str, err := ioutil.ReadFile(RootTemplate + templateCampaign)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}

	imageHeader := "http://mailer.gilkor.com/admin/temp/newsletters/137/header_oct2017.jpg"
	imageVoucher := "http://coma.greenparksolo.com/gilkor/images/testvoucher_image2.jpg"
	imageFooter := "http://mailer.gilkor.com/admin/temp/newsletters/137/footer_oct-2017.jpg"

	if program.ImageHeader != "" {
		imageHeader = program.ImageHeader
	}
	if program.ImageVoucher != "" {
		imageVoucher = program.ImageVoucher
	}
	if program.ImageFooter != "" {
		imageFooter = program.ImageFooter
	}

	result := string(str)
	result = strings.Replace(result, "%%full-name%%", target.HolderName, 1)
	result = strings.Replace(result, "%%link-voucher%%", target.VoucherUrl, 1)
	result = strings.Replace(result, "%%program-name%%", program.ProgramName, 1)
	result = strings.Replace(result, "%%image-header%%", imageHeader, 1)
	result = strings.Replace(result, "%%image-voucher%%", imageVoucher, 1)
	result = strings.Replace(result, "%%image-footer%%", imageFooter, 1)
	return result
}

func SendVoucherMail(domain, apiKey, publicApiKey, subject string, target []CampaignData, program ProgramCampaign) error {
	mg := mailgun.NewMailgun(domain, apiKey, publicApiKey)

	for _, v := range target {
		message := mailgun.NewMessage(
			Email,
			subject,
			subject,
			v.HolderEmail)
		message.SetHtml(makeMessageVoucherEmail(program, v))
		resp, id, err := mg.Send(message)
		if err != nil {
			return err
		}
		fmt.Printf("ID: %s Resp: %s\n", id, resp)
	}

	return nil
}
func makeMessageVoucherEmail(program ProgramCampaign, target CampaignData) string {
	templateCampaign := Config[program.AccountID]["email_campaign"]
	str, err := ioutil.ReadFile(RootTemplate + templateCampaign)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}

	imageHeader := "https://voucher.elys.id/assets/img/template_demo_email_01.jpg"
	imageVoucher := "https://voucher.elys.id/assets/img/template_demo_email_02.jpg"
	imageFooter := "https://voucher.elys.id/assets/img/template_demo_email_03.jpg"

	if program.ImageHeader != "" {
		imageHeader = program.ImageHeader
	}
	if program.ImageVoucher != "" {
		imageVoucher = program.ImageVoucher
	}
	if program.ImageFooter != "" {
		imageFooter = program.ImageFooter
	}
	fmt.Println(program)
	result := string(str)
	result = strings.Replace(result, "%%full-name%%", target.HolderName, 1)
	result = strings.Replace(result, "%%link-voucher%%", target.VoucherUrl, 1)
	result = strings.Replace(result, "%%program-name%%", program.ProgramName, 1)
	result = strings.Replace(result, "%%image-header%%", imageHeader, 1)
	result = strings.Replace(result, "%%image-voucher%%", imageVoucher, 1)
	result = strings.Replace(result, "%%image-footer%%", imageFooter, 1)
	return result
}

func SendVoucherMailV3(program ProgramCampaignV2, CampaignData []TargetEmail) error {
	url := "/v3/email/messages?key="
	param := BaseMail{
		From:     program.EmailSender,
		To:       CampaignData,
		Template: "basic-template1",
	}

	jsonParam, _ := json.Marshal(param)

	err := mailService("POST", url, jsonParam)
	if err != nil {
		return err
	}

	return nil
}

func SendVoucherMailV2(domain, apiKey, publicApiKey string, program ProgramCampaignV2, CampaignData []CampaignData) (int, error) {

	newTarget := []TargetEmail{}

	imageHeader := "https://voucher.elys.id/assets/img/template_demo_email_01.jpg"
	imageVoucher := "https://voucher.elys.id/assets/img/template_demo_email_02.jpg"
	imageFooter := "https://voucher.elys.id/assets/img/template_demo_email_03.jpg"

	if program.ImageHeader != "" {
		imageHeader = program.ImageHeader
	}
	if program.ImageVoucher != "" {
		imageVoucher = program.ImageVoucher
	}
	if program.ImageFooter != "" {
		imageFooter = program.ImageFooter
	}

	for _, v := range CampaignData {
		v.ProgramName = program.ProgramName
		v.ImageHeader = imageHeader
		v.ImageVoucher = imageVoucher
		v.ImageFooter = imageFooter
		v.EmailContent = program.EmailContent
		v.EmailSubject = program.EmailSubject
		target := TargetEmail{
			EmailAddress: v.HolderEmail,
			Data:         v,
		}

		newTarget = append(newTarget, target)
	}

	err := SendVoucherMailV3(program, newTarget)
	if err != nil {
		return 0, err
	}

	return len(CampaignData), nil
}

func makeMessageVoucherEmailV2(program ProgramCampaignV2, target CampaignData) string {
	str, err := ioutil.ReadFile(RootTemplate + program.Template)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}

	imageHeader := "https://voucher.elys.id/assets/img/template_demo_email_01.jpg"
	imageVoucher := "https://voucher.elys.id/assets/img/template_demo_email_02.jpg"
	imageFooter := "https://voucher.elys.id/assets/img/template_demo_email_03.jpg"

	if program.ImageHeader != "" {
		imageHeader = program.ImageHeader
	}
	if program.ImageVoucher != "" {
		imageVoucher = program.ImageVoucher
	}
	if program.ImageFooter != "" {
		imageFooter = program.ImageFooter
	}
	fmt.Println(program)
	result := string(str)
	result = strings.Replace(result, "%%full-name%%", target.HolderName, 1)
	result = strings.Replace(result, "%%link-voucher%%", target.VoucherUrl, 1)
	result = strings.Replace(result, "%%program-name%%", program.ProgramName, 1)
	result = strings.Replace(result, "%%content-email%%", program.EmailContent, 1)
	result = strings.Replace(result, "%%image-header%%", imageHeader, 1)
	result = strings.Replace(result, "%%image-voucher%%", imageVoucher, 1)
	result = strings.Replace(result, "%%image-footer%%", imageFooter, 1)
	return result
}

func SendConfirmationEmail(emailSender, subject string, target ConfirmationEmailRequest, accountId string) error {
	// mg := mailgun.NewMailgun(domain, apiKey, publicApiKey)

	// for _, v := range target.ListEmail {
	// 	fmt.Println(v)
	// 	if v != "" {
	// 		message := mailgun.NewMessage(
	// 			Email,
	// 			subject,
	// 			subject,
	// 			v)
	// 		message.SetHtml(makeMessageConfirmationEmail(accountId, target))
	// 		resp, id, err := mg.Send(message)
	// 		if err != nil {
	// 			fmt.Println(message)
	// 			fmt.Println(err.Error())
	// 			return err
	// 		}
	// 		fmt.Printf("ID: %s Resp: %s\n", id, resp)
	// 	}
	// }

	targetMail := []TargetEmail{}
	target.EmailSubject = subject

	for _, v := range target.ListEmail {
		tempTarget := TargetEmail{
			EmailAddress: v,
			Data:         target,
		}

		targetMail = append(targetMail, tempTarget)
	}

	url := "/v3/email/messages?key="
	param := BaseMail{
		From:     emailSender,
		To:       targetMail,
		Template: "s1-confirmation-email",
	}

	jsonParam, _ := json.Marshal(param)

	err := mailService("POST", url, jsonParam)
	if err != nil {
		return err
	}

	return nil
}

func makeMessageConfirmationEmail(accountId string, target ConfirmationEmailRequest) string {
	// %%full-name%%
	// %%link-voucher%%
	fmt.Println(Config[accountId]["email_transaction_confirmation"])
	templateCampaign := Config[accountId]["email_transaction_confirmation"]
	str, err := ioutil.ReadFile(RootTemplate + templateCampaign)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}

	voucher := ""
	for _, v := range target.ListVoucher {
		voucher += "<tr><td style='color:#ffffff; padding:10px 0px; background-color: #69cdcd;'>"
		voucher += v
		voucher += "</td><td style='color:#ffffff; padding:10px 0px; background-color: #69cdcd;'>"
		voucher += target.PartnerName
		voucher += "</td></tr>"
	}

	result := string(str)
	result = strings.Replace(result, "%%full-name%%", target.Holder, 1)
	result = strings.Replace(result, "%%transaction-code%%", target.TransactionCode, 1)
	result = strings.Replace(result, "%%transaction-date%%", target.TransactionDate.Format("2006-01-02 15:04:05"), 1)
	result = strings.Replace(result, "%%program-name%%", target.ProgramName, 1)
	result = strings.Replace(result, "%%voucher-code%%", voucher, 1)

	return result
}

// Query Database
func InsertCampaign(request ProgramCampaign, user string) (string, error) {
	tx, err := db.Beginx()
	if err != nil {
		fmt.Println(err.Error())
		return "", ErrServerInternal
	}
	defer tx.Rollback()

	var res []string
	campaign, err := GetCampaign(request.ProgramID)
	if campaign.ProgramID == "" {
		q := `
			INSERT INTO program_campaigns(
				program_id
				, header_image
				, voucher_image
				, footer_image
				, created_by
				, created_at
				, status
			)
			VALUES (?, ?, ?, ?, ?, ?, ?)
			RETURNING
				id
		`

		err = tx.Select(&res, tx.Rebind(q), request.ProgramID, request.ImageHeader, request.ImageVoucher, request.ImageFooter, user, time.Now(), StatusCreated)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println(q)
			return "", ErrServerInternal
		}
	} else {
		q := `
			UPDATE program_campaigns
			SET
				header_image = ?
				, voucher_image = ?
				, footer_image = ?
				, updated_by = ?
				, updated_at = ?
			WHERE
				program_id = ?
				AND status = ?
		`

		_, err = tx.Exec(tx.Rebind(q), request.ImageHeader, request.ImageVoucher, request.ImageFooter, user, time.Now(), request.ProgramID, StatusCreated)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println(q)
			return "", ErrServerInternal
		}
		res = append(res, request.ProgramID)
	}

	if err := tx.Commit(); err != nil {
		fmt.Println(err.Error())
		return "", ErrServerInternal
	}

	return res[0], nil
}

func InsertCampaignV2(request ProgramCampaignV2, user string) (string, error) {
	tx, err := db.Beginx()
	if err != nil {
		fmt.Println(err.Error())
		return "", ErrServerInternal
	}
	defer tx.Rollback()

	var res []string
	campaign, err := GetCampaign(request.ProgramID)
	if campaign.ProgramID == "" {
		q := `
			INSERT INTO program_campaigns(
				program_id
				, email_template
				, email_subject
				, email_sender
				, email_content
				, header_image
				, voucher_image
				, footer_image
				, created_by
				, created_at
				, status
			)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ? ,? ,?)
			RETURNING
				id
		`

		err = tx.Select(&res, tx.Rebind(q), request.ProgramID, request.Template, request.EmailSubject, request.EmailSender, request.EmailContent, request.ImageHeader, request.ImageVoucher, request.ImageFooter, user, time.Now(), StatusCreated)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println(q)
			return "", ErrServerInternal
		}
	} else {
		q := `
			UPDATE program_campaigns
			SET
				email_template = ?
				, email_subject = ?
				, email_sender = ?
				, email_content = ?
				, header_image = ?
				, voucher_image = ?
				, footer_image = ?
				, updated_by = ?
				, updated_at = ?
			WHERE
				program_id = ?
				AND status = ?
		`

		_, err = tx.Exec(tx.Rebind(q), request.Template, request.EmailSubject, request.EmailSender, request.EmailContent, request.ImageHeader, request.ImageVoucher, request.ImageFooter, user, time.Now(), request.ProgramID, StatusCreated)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println(q)
			return "", ErrServerInternal
		}
		res = append(res, request.ProgramID)
	}

	if err := tx.Commit(); err != nil {
		fmt.Println(err.Error())
		return "", ErrServerInternal
	}

	return res[0], nil
}

func GetCampaign(programId string) (ProgramCampaign, error) {
	q := `
		SELECT
			id
			, program_id
			, header_image
			, voucher_image
			, footer_image
			, created_by
		FROM
			program_campaigns
		WHERE
			program_id = ?
			AND status = ?
	`

	var resv []ProgramCampaign
	if err := db.Select(&resv, db.Rebind(q), programId, StatusCreated); err != nil {
		fmt.Println(err.Error())
		return ProgramCampaign{}, ErrServerInternal
	}
	if len(resv) < 1 {
		return ProgramCampaign{}, nil
	}

	return resv[0], nil
}

func GetCampaignV2(programId string) (ProgramCampaignV2, error) {
	q := `
		SELECT
			pc.id
			, pc.program_id
			, p.name as program_name
			, p.account_id
			, pc.email_template
			, pc.email_subject
			, pc.email_sender
			, pc.email_content
			, pc.header_image
			, pc.voucher_image
			, pc.footer_image
			, pc.created_by
			, pc.created_at
		FROM
			program_campaigns as pc
		JOIN
			programs as p
		ON
			pc.program_id = p.id
		WHERE
			pc.program_id = ?
			AND pc.status = ?
	`

	var resv []ProgramCampaignV2
	if err := db.Select(&resv, db.Rebind(q), programId, StatusCreated); err != nil {
		fmt.Println(err.Error())
		return ProgramCampaignV2{}, ErrServerInternal
	}
	if len(resv) < 1 {
		return ProgramCampaignV2{}, nil
	}

	return resv[0], nil
}

func GetEmail(transaction string) (ConfirmationEmail, error) {
	q := `
		SELECT DISTINCT
			v.holder_email AS member, a.email AS account, p.email AS partner
		FROM
			vouchers AS v
		JOIN
			transaction_details AS td
		ON
			td.voucher_id = v.id
		JOIN
			transactions AS t
		ON
			t.id = td.transaction_id
		JOIN
			partners as p
		ON
			p.id = t.partner_id
		JOIN
			accounts AS a
		ON
			a.id = t.account_id
		WHERE
			( t.id = ?
			OR t.transaction_code = ? )
			AND t.status = ?
	`

	var resv []ConfirmationEmail
	if err := db.Select(&resv, db.Rebind(q), transaction, transaction, StatusCreated); err != nil {
		fmt.Println(err.Error())
		return ConfirmationEmail{}, ErrServerInternal
	}
	if len(resv) < 1 {
		return ConfirmationEmail{}, nil
	}

	return resv[0], nil
}

func mailService(method, url string, param []byte) error {
	domain := os.Getenv("MAIL_DOMAIN")
	key := os.Getenv("MAIL_KEY")

	fmt.Printf("url = " + domain + url + key)
	fmt.Printf("%s", param)

	req, err := http.NewRequest(method, domain+url+key, bytes.NewBuffer(param))
	if err != nil {
		panic(err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New(resp.Status)
	}

	return nil
}
