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
		ProgamName string
		VoucherImg string
		HolderName string
		VoucherUrl string
	}
	TargetEmail struct {
		HolderEmail string
		HolderName  string
		VoucherUrl  string
	}

	ProgramCampaign struct {
		Id           string `db:"id"`
		ProgramId    string `db:"program_id"`
		ProgramName  string `db:"program_name"`
		ImageHeader  string `db:"header_image"`
		ImageVoucher string `db:"voucher_image"`
		ImageFooter  string `db:"footer_image"`
		CreatedBy    string `db:"created_by"`
		CreatedAt    string `db:"created_at"`
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
		fmt.Printf("ID: %s Resp: %s\n", id, resp)
	}

	return nil
}
func makeMessageEmailSedayuOne(program ProgramCampaign, target TargetEmail) string {
	// %%full-name%%
	// %%link-voucher%%
	str, err := ioutil.ReadFile(RootTemplate + "sedayu_one_campaign")
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}

	imageHeader := "http://mailer.gilkor.com/admin/temp/newsletters/33/header_apr_fix.jpg"
	imageVoucher := "http://mailer.gilkor.com/admin/temp/newsletters/116/pik_marvelous_prizes.jpg"
	imageFooter := "http://mailer.gilkor.com/admin/temp/newsletters/116/bannerlogo_footer.jpg"

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

// Query Database
func InsertCampaign(request ProgramCampaign, user string) (string, error) {
	fmt.Println("Add")
	tx, err := db.Beginx()
	if err != nil {
		fmt.Println(err.Error())
		return "", ErrServerInternal
	}
	defer tx.Rollback()

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

	var res []string
	err = tx.Select(&res, tx.Rebind(q), request.ProgramId, request.ImageHeader, request.ImageVoucher, request.ImageFooter, user, StatusCreated)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println(q)
		return "", ErrServerInternal
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
