package model

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/mailgun/mailgun-go.v1"
	"strings"
	"time"
)

var (
	Config       map[string]map[string]string
	Domain       string
	ApiKey       string
	PublicApiKey string
	RootTemplate string
	RootUrl      string
	Email        string
)

type (
	TargetEmail struct {
		HolderEmail string
		HolderName  string
		VoucherUrl  string
	}

	ProgramCampaign struct {
		Id           string `db:"id"`
		ProgramId    string `db:"program_id"`
		ProgramName  string `db:"program_name"`
		AccountId    string `db:"account_id"`
		ImageHeader  string `db:"header_image"`
		ImageVoucher string `db:"voucher_image"`
		ImageFooter  string `db:"footer_image"`
		CreatedBy    string `db:"created_by"`
		CreatedAt    string `db:"created_at"`
	}

	ConfirmationEmail struct {
		EmailAccount string `db:"account"`
		EmailPartner string `db:"partner"`
		EmailMember  string `db:"member"`
	}

	ConfirmationEmailRequest struct {
		Holder          string
		ProgramName     string
		TransactionCode string
		TransactionDate string
		PartnerName     string
		ListEmail       []string
		ListVoucher     []string
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

func SendMailSedayuOne(domain, apiKey, publicApiKey, subject string, target []TargetEmail, program ProgramCampaign) error {
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
		UpdateBroadcastUserState(v.HolderEmail, program.ProgramId)
		fmt.Printf("ID: %s Resp: %s\n", id, resp)
	}

	return nil
}
func makeMessageEmailSedayuOne(program ProgramCampaign, target TargetEmail) string {
	// %%full-name%%
	// %%link-voucher%%
	fmt.Println(program.AccountId)
	templateCampaign := Config[program.AccountId]["email_campaign"]
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

func SendConfirmationEmail(domain, apiKey, publicApiKey, subject string, target ConfirmationEmailRequest, accountId string) error {
	mg := mailgun.NewMailgun(domain, apiKey, publicApiKey)

	for _, v := range target.ListEmail {
		message := mailgun.NewMessage(
			Email,
			subject,
			subject,
			v)
		message.SetHtml(makeMessageConfirmationEmail(accountId, target))
		resp, id, err := mg.Send(message)
		if err != nil {
			return err
		}
		fmt.Printf("ID: %s Resp: %s\n", id, resp)
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
	result = strings.Replace(result, "%%transaction-date%%", target.TransactionDate, 1)
	result = strings.Replace(result, "%%program-name%%", target.ProgramName, 1)
	result = strings.Replace(result, "%%voucher-code%%", voucher, 1)
	return result
}

// Query Database
func InsertCampaign(request ProgramCampaign, user string) (string, error) {
	fmt.Println("Add")
	tx, err := db.Beginx()
	if err != nil {
		fmt.Println(err.Error())
		return "", ErrServerInternal
	}
	defer tx.Rollback()

	var res []string

	campaign, err := GetCampaign(request.ProgramId)
	if campaign.ProgramId == "" {
		q := `
			INSERT INTO program_campaigns(
				program_id
				, header_image
				, voucher_image
				, footer_image
				, created_by
				, status
			)
			VALUES (?, ?, ?, ?, ?, ?)
			RETURNING
				id
		`

		err = tx.Select(&res, tx.Rebind(q), request.ProgramId, request.ImageHeader, request.ImageVoucher, request.ImageFooter, user, StatusCreated)
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

		_, err = tx.Exec(tx.Rebind(q), request.ImageHeader, request.ImageVoucher, request.ImageFooter, user, time.Now(), request.ProgramId, StatusCreated)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println(q)
			return "", ErrServerInternal
		}
		res = append(res, request.ProgramId)
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
