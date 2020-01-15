package model

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"
)

type (
	// Blast : represent of blast table model
	Blast struct {
		ID             string          `db:"id" json:"id,omitempty"`
		Subject        string          `db:"subject" json:"subject,omitempty"`
		Program        Program         `db:"program" json:"program,omitempty"`
		CompanyID      string          `db:"company_id" json:"company_id,omitempty"`
		ImageHeader    string          `db:"image_header" json:"image_header,omitempty"`
		ImageFooter    string          `db:"image_footer" json:"image_footer,omitempty"`
		EmailContent   string          `db:"email_content" json:"email_content,omitempty"`
		Template       string          `db:"template" json:"template,omitempty"`
		RecipientsData []RecipientData `db:"recipients" json:"recipients,omitempty"`
		CreatedAt      *time.Time      `db:"created_at" json:"created_at,omitempty"`
		CreatedBy      string          `db:"created_by" json:"created_by,omitempty"`
		UpdatedAt      *time.Time      `db:"updated_at" json:"updated_at,omitempty"`
		UpdatedBy      string          `db:"updated_by" json:"updated_by,omitempty"`
		Status         string          `db:"status" json:"status,omitempty"`
	}

	// Blasts : List of Blast
	Blasts []Blast

	// Recipient : recipient email and data
	Recipient struct {
		EmailAddress string        `json:"email_address"`
		Data         RecipientData `json:"data"`
	}

	// BlastRequest : body data of post blast
	BlastRequest struct {
		From     string      `json:"from"`
		To       []Recipient `json:"to"`
		Subject  string      `json:"subject,omitempty"`
		Message  string      `json:"message,omitempty"`
		Template string      `json:"template,omitempty"`
	}

	// RecipientData : detail data for each recipient per blast
	RecipientData struct {
		HolderEmail  string  `json:"email"`
		HolderName   string  `json:"name"`
		HolderPhone  string  `json:"phone"`
		VoucherURL   string  `json:"voucher_url"`
		VoucherObj   Voucher `json:"voucher"`
		ProgramName  string  `json:"program_name"`
		ImageHeader  string  `json:"image_header"`
		ImageVoucher string  `json:"image_voucher"`
		ImageFooter  string  `json:"image_footer"`
		EmailContent string  `json:"email_content"`
		EmailSubject string  `json:"email_subject"`
	}
)

// SendEmailBlast send email blast
func SendEmailBlast(reqBlast Blast) error {
	recipients := []Recipient{}

	imageHeader := reqBlast.ImageHeader
	imageVoucher := reqBlast.Program.ImageURL
	imageFooter := reqBlast.ImageFooter

	for _, recipientData := range reqBlast.RecipientsData {
		recipientData.ProgramName = reqBlast.Program.Name
		recipientData.ImageHeader = imageHeader
		recipientData.ImageVoucher = imageVoucher
		recipientData.ImageFooter = imageFooter
		recipientData.EmailContent = reqBlast.EmailContent
		recipientData.EmailSubject = reqBlast.Subject

		recipient := Recipient{
			EmailAddress: recipientData.HolderEmail,
			Data:         recipientData,
		}

		recipients = append(recipients, recipient)
	}

	url := "/v3/email/messages?key="
	param := BlastRequest{
		From:     "voucher@elys.id",
		To:       recipients,
		Template: "blast-template",
	}

	jsonParam, _ := json.Marshal(param)

	err := mailService("POST", url, jsonParam)
	if err != nil {
		return err
	}

	return nil
}

func mailService(method, url string, param []byte) error {
	domain := os.Getenv("MAIL_DOMAIN")
	mailKey := os.Getenv("MAIL_KEY")

	fmt.Printf("url = " + domain + url + mailKey)
	fmt.Printf("%s", param)

	req, err := http.NewRequest(method, domain+url+mailKey, bytes.NewBuffer(param))
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
