package model

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type (
	//Voucher model
	// Voucher struct {
	// 	ID           string         `json:"id,omitempty" db:"id"`
	// 	Code         string         `json:"code,omitempty" db:"code"`
	// 	ReferenceNo  string         `json:"reference_no,omitempty" db:"reference_no"`
	// 	Holder       *string        `json:"holder,omitempty" db:"holder"`
	// 	HolderDetail types.JSONText `json:"holder_detail,omitempty" db:"holder_detail"`
	// 	ProgramID    string         `json:"program_id,omitempty" db:"program_id"`
	// 	ValidAt      *time.Time     `json:"valid_at,omitempty" db:"valid_at"`
	// 	ExpiredAt    *time.Time     `json:"expired_at,omitempty" db:"expired_at"`
	// 	State        string         `json:"state,omitempty" db:"state"`
	// 	CreatedBy    string         `json:"created_by,omitempty" db:"created_by"`
	// 	CreatedAt    *time.Time     `json:"created_at,omitempty" db:"created_at"`
	// 	UpdatedBy    string         `json:"updated_by,omitempty" db:"updated_by"`
	// 	UpdatedAt    *time.Time     `json:"updated_at,omitempty" db:"updated_at"`
	// 	Status       string         `json:"status,omitempty" db:"status"`
	// 	Count        int            `db:"count" json:"-"`
	// }

	// ConfirmationEmailRecipient : list of emails from separated tables, who will received an email confirmation
	ConfirmationEmailRecipient struct {
		PartnerEmails string `db:"partner"`
		CompanyEmails string `db:"company"`
		HolderEmail   string `db:"holder"`
	}

	ConfirmationEmailData struct {
		Holder          string    `json:"name"`
		ProgramName     string    `json:"program_name"`
		TransactionCode string    `json:"transaction_code"`
		TransactionDate time.Time `json:"transaction_date"`
		PartnerName     string    `json:"partner_name"`
		ListVoucher     []string  `json:"list_voucher"`
		EmailSubject    string    `json:"email_subject"`
	}
)

// GetEmail : get emails of email confirmation recipient
func GetEmail(transactionID string, companyID string) (ConfirmationEmailRecipient, error) {
	q := `
		SELECT DISTINCT
			v.holder_detail::json->>'holder_email' AS holder, c.value AS company, p.emails AS partner
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
	emails, err := GetEmail(t.ID, t.CompanyId)
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

	if strings.Contains(emails.PartnerEmails, ",") {
		tmpPartnerEmails := strings.Split(emails.PartnerEmails, ",")
		for _, v := range tmpPartnerEmails {
			listEmail = append(listEmail, strings.Replace(v, " ", "", -1))
		}
	} else {
		listEmail = append(listEmail, emails.PartnerEmails)
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

	// partner, err := GetPartnerByID(n)

	recipients := []Recipient{}

	for _, emailRecipient := range listEmail {

		data := ConfirmationEmailData{
			Holder:          *t.Vouchers[0].Holder,
			ProgramName:     t.Programs[0].Name,
			PartnerName:     t.Partner.Name,
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

	configs, err := GetConfigs(t.CompanyId, "company")
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

// // partner detail
// partner, err := model.FindPartnerById(rd.Partner)
// if err != nil {
// 	status = http.StatusInternalServerError
// 	res.AddError(its(status), model.ErrCodeInternalError, "partner error "+model.ErrMessageInternalError+"("+err.Error()+")", logger.TraceID)
// 	logger.SetStatus(status).Log("param :", rd, "partner error response :", res.Errors.ToString())
// 	render.JSON(w, res, status)
// 	return
// }

// req := model.ConfirmationEmailRequest{
// 	Holder:          voucherDetail.VoucherData[0].HolderDescription.String,
// 	ProgramName:     voucherDetail.VoucherData[0].ProgramName,
// 	PartnerName:     partner.Name,
// 	TransactionCode: transaction.TransactionCode,
// 	TransactionDate: transaction.CreatedAt,
// 	ListEmail:       listEmail,
// 	ListVoucher:     listVoucher,
// }

// senderMail := a.User.Account.SenderEmail
// mailKey := a.User.Account.MailKey.String
// title := "Elys Voucher Confirmation"

// if err := model.SendConfirmationEmail(senderMail, title, req, a.User.Account.Id, mailKey); err != nil {
// 	logger.SetStatus(status).Info("param :", listEmail, "response :", err.Error())
// }

// res = NewResponse(TransactionResponse{
// 	TransactionID:   transaction.Id,
// 	TransactionCode: transaction.TransactionCode,
// 	DiscountValue:   transaction.DiscountValue,
// 	Created_at:      transaction.CreatedAt,
// 	Vouchers:        listVoucher,
// 	Voucher:         voucher,
// 	Partner:         MobilePartnerObj{partner.Id, partner.Name}})

// SendEmailConfirmation : send email confirmation after use voucher
// func (t *Transaction) SendEmailConfirmation() error {
// 	recipients := []Recipient{}

// 	data := RecipientRequestData{
// 		ProgramName:  b.Program.Name,
// 		ImageHeader:  imageHeader,
// 		ImageVoucher: imageVoucher,
// 		ImageFooter:  imageFooter,
// 		EmailContent: b.EmailContent,
// 		EmailSubject: b.Subject,
// 		VoucherURL:   generateLink(b.CompanyID, recipientData.VoucherID),
// 	}

// 	recipient := Recipient{
// 		EmailAddress: recipientData.HolderEmail,
// 		Data:         data,
// 	}

// 	recipients = append(recipients, recipient)

// 	url := "/v3/email/messages?key="
// 	param := BlastRequest{
// 		From:     b.Sender,
// 		To:       recipients,
// 		Template: b.Template,
// 	}

// 	jsonParam, _ := json.Marshal(param)

// 	success, err := mailService("POST", url, jsonParam)
// 	if err != nil {
// 		return err
// 	}

// 	if success {
// 		// Update blast status
// 		err = b.UpdateBlastStatus()
// 		if err != nil {
// 			return errors.New("Failed when update blast status ," + err.Error())
// 		}
// 	}

// 	return nil
// }
