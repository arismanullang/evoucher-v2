package model

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type (

	// ConfirmationEmailRecipient : list of emails from separated tables, who will received an email confirmation
	ConfirmationEmailRecipient struct {
		OutletEmails  string `db:"outlet"`
		CompanyEmails string `db:"company"`
		HolderEmail   string `db:"holder"`
	}

	ConfirmationEmailData struct {
		Holder          string    `json:"name"`
		ProgramName     string    `json:"program_name"`
		TransactionCode string    `json:"transaction_code"`
		TransactionDate time.Time `json:"transaction_date"`
		OutletName      string    `json:"outlet_name"`
		ListVoucher     []string  `json:"list_voucher"`
		EmailSubject    string    `json:"email_subject"`
	}
)

// GetEmail : get emails of email confirmation recipient
func GetEmail(transactionID string, companyID string) (ConfirmationEmailRecipient, error) {
	q := `
		SELECT DISTINCT
			v.holder_detail::json->>'holder_email' AS holder, c.value AS company, p.emails AS outlet
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
			outlets as p
		ON
			p.id = t.outlet_id
		JOIN
			config AS c
		ON
			c.company_id = t.company_id
		WHERE t.id = ?
			AND t.company_id = ?
			AND c.key = 'finance_emails'
			AND c.category = 'company'
			AND t.status = ?
			
	`

	var resv []ConfirmationEmailRecipient
	if err := db.Select(&resv, db.Rebind(q), transactionID, companyID, StatusCreated); err != nil {
		fmt.Println(err.Error())
		return ConfirmationEmailRecipient{}, err
	}
	if len(resv) < 1 {
		return ConfirmationEmailRecipient{}, nil
	}

	return resv[0], nil
}

// SendEmailConfirmation : send email confirmation after use voucher
func (t Transaction) SendEmailConfirmation() error {
	// get list email
	listEmail := []string{}
	emails, err := GetEmail(t.ID, t.CompanyID)
	if err != nil {
		return err
	}

	if strings.Contains(emails.CompanyEmails, ",") {
		tmpCompanyEmails := strings.Split(emails.CompanyEmails, ",")
		for _, v := range tmpCompanyEmails {
			listEmail = append(listEmail, strings.Replace(v, " ", "", -1))
		}
	} else {
		listEmail = append(listEmail, emails.CompanyEmails)
	}

	if strings.Contains(emails.OutletEmails, ",") {
		tmpOutletEmails := strings.Split(emails.OutletEmails, ",")
		for _, v := range tmpOutletEmails {
			listEmail = append(listEmail, strings.Replace(v, " ", "", -1))
		}
	} else {
		listEmail = append(listEmail, emails.OutletEmails)
	}

	if strings.Contains(emails.HolderEmail, ",") {
		tmpHolderEmail := strings.Split(emails.HolderEmail, ",")
		for _, v := range tmpHolderEmail {
			listEmail = append(listEmail, strings.Replace(v, " ", "", -1))
		}
	} else {
		listEmail = append(listEmail, emails.HolderEmail)
	}

	// voucher detail
	listVoucher := []string{}
	for _, voucher := range t.Vouchers {
		listVoucher = append(listVoucher, voucher.Code)
	}

	//HolderDetail :type struct Voucher.types.JSONText.Unmarshal(&HolderDetail)
	holderDetail := HolderDetail{}
	t.Vouchers[0].HolderDetail.Unmarshal(&holderDetail)

	recipients := []Recipient{}

	for _, emailRecipient := range listEmail {

		data := ConfirmationEmailData{
			Holder:          holderDetail.Name,
			ProgramName:     t.Vouchers[0].ProgramName,
			OutletName:      t.OutletName,
			TransactionCode: t.TransactionCode,
			TransactionDate: *t.CreatedAt,
			ListVoucher:     listVoucher,
			EmailSubject:    "Elys Voucher Confirmation",
		}

		recipient := Recipient{
			EmailAddress: emailRecipient,
			Data:         data,
		}

		recipients = append(recipients, recipient)
	}

	configs, err := GetConfigs(t.CompanyID, "company")
	if err != nil {
		return ErrorInternalServer
	}

	sender := ""
	template := ""

	tmpSender := configs[CompanyEmailSender]
	if tmpSender != nil {
		sender = tmpSender.(string)
	}

	tmpTemplate := configs[CompanyEmailTemplate]
	if tmpTemplate != nil {
		template = tmpTemplate.(string)
	}

	url := "/v3/email/messages?key="
	param := BlastRequest{
		From:     sender,
		To:       recipients,
		Template: template,
	}

	jsonParam, _ := json.Marshal(param)

	_, err = mailService("POST", url, jsonParam)
	if err != nil {
		return err
	}

	return nil

}
