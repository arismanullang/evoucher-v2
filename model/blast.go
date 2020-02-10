package model

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gilkor/evoucher-v2/util"
	"github.com/jmoiron/sqlx/types"
)

type (

	// Blast : represent of blast table model
	Blast struct {
		ID             string          `db:"id" json:"id,omitempty"`
		Subject        string          `db:"subject" json:"subject,omitempty"`
		ProgramID      string          `db:"program_id" json:"program_id,omitempty"`
		Program        *Program        `json:"program,omitempty"`
		BlastProgram   types.JSONText  `db:"program" json:"-"`
		CompanyID      string          `db:"company_id" json:"company_id,omitempty"`
		ImageHeader    string          `db:"image_header" json:"image_header,omitempty"`
		ImageFooter    string          `db:"image_footer" json:"image_footer,omitempty"`
		EmailContent   string          `db:"email_content" json:"email_content,omitempty"`
		Template       string          `db:"template" json:"template,omitempty"`
		BlastRecipient BlastRecipients `json:"recipients,omitempty"`
		Recipient      types.JSONText  `db:"recipients" json:"-"`
		CreatedAt      *time.Time      `db:"created_at" json:"created_at,omitempty"`
		CreatedBy      string          `db:"created_by" json:"created_by,omitempty"`
		UpdatedAt      *time.Time      `db:"updated_at" json:"updated_at,omitempty"`
		UpdatedBy      string          `db:"updated_by" json:"updated_by,omitempty"`
		Status         string          `db:"status" json:"status,omitempty"`
	}

	// Blasts : List of Blast
	Blasts []Blast

	// BlastRecipient : detail data for each recipient per blast
	BlastRecipient struct {
		ID          int        `db:"id" json:"id,omitempty"`
		BlastID     string     `db:"blast_id" json:"blast_id,omitempty"`
		HolderEmail string     `db:"email" json:"email,omitempty"`
		HolderName  string     `db:"name" json:"name,omitempty"`
		VoucherID   string     `db:"voucher_id" json:"voucher_id,omitempty"`
		VoucherObj  Voucher    `json:"voucher,omitempty"`
		CreatedAt   *time.Time `db:"created_at" json:"created_at,omitempty"`
		CreatedBy   string     `db:"created_by" json:"created_by,omitempty"`
		UpdatedAt   *time.Time `db:"updated_at" json:"updated_at,omitempty"`
		UpdatedBy   string     `db:"updated_by" json:"updated_by,omitempty"`
		Status      string     `db:"status" json:"status,omitempty"`
	}

	// BlastRecipients : List of BlastRecipient
	BlastRecipients []BlastRecipient

	// For NUDGE

	// BlastRequest : body data of post blast
	BlastRequest struct {
		From     string      `json:"from"`
		To       []Recipient `json:"to"`
		Subject  string      `json:"subject,omitempty"`
		Message  string      `json:"message,omitempty"`
		Template string      `json:"template,omitempty"`
	}

	// Recipient : recipient email and data for nudge
	Recipient struct {
		EmailAddress string               `json:"email_address"`
		Data         RecipientRequestData `json:"data"`
	}

	// RecipientRequestData : recipient data request
	RecipientRequestData struct {
		ProgramName  string `json:"program_name,omitempty"`
		ImageHeader  string `json:"image_header,omitempty"`
		ImageVoucher string `json:"image_voucher,omitempty"`
		ImageFooter  string `json:"image_footer,omitempty"`
		EmailContent string `json:"email_content,omitempty"`
		EmailSubject string `json:"email_subject,omitempty"`
	}
)

//GetBlasts : get list blast by custom filter
func GetBlasts(qp *util.QueryParam) (*Blasts, bool, error) {
	return getBlasts("1", "1", qp)
}

//GetBlastByID : get blast by specified ID
func GetBlastByID(qp *util.QueryParam, id string) (*Blast, error) {
	blasts, _, err := getBlasts("id", id, qp)
	if err != nil {
		return &Blast{}, errors.New("Failed when select on blast ," + err.Error())
	}
	blast := &(*blasts)[0]
	//get program detail
	program, err := GetProgramByID(blast.ProgramID, qp)
	if err != nil {
		return &Blast{}, errors.New("Failed when select on blast recipient ," + err.Error())
	}
	blast.Program = program

	//get blast recipient data
	recipients, _, err := getBlastRecipient(blast.ID, qp)
	if err != nil {
		return &Blast{}, errors.New("Failed when select on blast recipient ," + err.Error())
	}
	blast.BlastRecipient = *recipients

	return blast, nil
}

func getBlastRecipient(blastID string, qp *util.QueryParam) (*BlastRecipients, bool, error) {
	q, err := qp.GetQueryByDefaultStruct(BlastRecipient{})
	if err != nil {
		return &BlastRecipients{}, false, err
	}
	q += `
			FROM
				blast_recipients as BlastRecipient
			WHERE 
				status = ?
			AND
				blast_id = ?
`
	q += qp.GetQuerySort()
	// q += qp.GetQueryLimit()
	// fmt.Println(q)
	fmt.Println("query struct :", q)
	var resd BlastRecipients
	err = db.Select(&resd, db.Rebind(q), StatusCreated, blastID)
	if err != nil {
		fmt.Println("blast_id = ", blastID)
		fmt.Println("err = ", err)
		return &BlastRecipients{}, false, err
	}
	if len(resd) < 1 {
		fmt.Println("blast_id res not found = ", blastID)
		return &BlastRecipients{}, false, ErrorResourceNotFound
	}
	next := false
	if len(resd) > qp.Count {
		next = true
	}
	if len(resd) < qp.Count {
		qp.Count = len(resd)
	}
	return &resd, next, nil
}

func getBlasts(k, v string, qp *util.QueryParam) (*Blasts, bool, error) {

	q, err := qp.GetQueryByDefaultStruct(Blast{})
	if err != nil {
		return &Blasts{}, false, err
	}

	q += `
			FROM
				m_blasts blast
			WHERE 
				status IN (?, ?)
			AND ` + k + ` = ?`

	q += qp.GetQuerySort()
	q += qp.GetQueryLimit()
	util.DEBUG(q)
	var resd Blasts
	err = db.Select(&resd, db.Rebind(q), StatusCreated, StatusSubmitted, v)
	if err != nil {
		return &Blasts{}, false, err
	}
	if len(resd) < 1 {
		return &Blasts{}, false, ErrorResourceNotFound
	}
	next := false
	if len(resd) > qp.Count {
		next = true
	}
	if len(resd) < qp.Count {
		qp.Count = len(resd)
	}

	err = json.Unmarshal([]byte(resd[0].Recipient), &resd[0].BlastRecipient)
	if err != nil {
		return &Blasts{}, false, err
	}

	var programs Programs
	err = json.Unmarshal([]byte(resd[0].BlastProgram), &programs)
	if err != nil {
		return &Blasts{}, false, err
	}

	resd[0].Program = &programs[0]

	return &resd, next, nil
}

//Insert : single row inset into table
func (b *Blast) Insert() (*Blasts, error) {
	tx, err := db.Beginx()
	if err != nil {
		return nil, errors.New("Failed when insert new blast ," + err.Error())
	}
	defer tx.Rollback()

	q := `INSERT INTO 
				blasts
				( 
					company_id
					, subject
					, program_id
					, image_header
					, image_footer
					, email_content
					, template
					, created_by
					, updated_by					
					, status
				)
			VALUES 
				( ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			RETURNING 
			id
			, company_id
			, subject
			, program_id
			, image_header
			, image_footer
			, email_content
			, template
			, created_at
			, created_by
			, updated_at
			, updated_by					
			, status
	`

	util.DEBUG(q)
	var res Blasts

	err = tx.Select(&res, tx.Rebind(q), b.CompanyID, b.Subject, b.Program.ID, b.ImageHeader, b.ImageFooter, b.EmailContent, b.Template, b.CreatedBy, b.UpdatedBy, StatusCreated)
	if err != nil {
		return nil, err
	}

	//insert blast detail
	for _, r := range b.BlastRecipient {
		q := `
			INSERT INTO blast_recipients(
				blast_id
				, name
				, email
				, voucher_id
				, created_by
				, updated_by
				, status					
			)
			VALUES (?, ?, ?, ?, ?, ?, ?)
		`

		_, err := tx.Exec(tx.Rebind(q), res[0].ID, r.HolderName, r.HolderEmail, r.VoucherObj.ID, res[0].CreatedBy, res[0].UpdatedBy, StatusCreated)
		if err != nil {
			return nil, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	res[0].BlastRecipient = b.BlastRecipient
	*b = res[0]
	return &res, nil
}

//Update : modify data
func (b *Blast) Update() error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	q := `UPDATE
				blasts 
			SET
				subject = ?
				, program_id = ?
				, image_header = ?
				, image_footer = ?
				, email_content = ?
				, template = ?
				, updated_at = now()
				, updated_by = ?				
			WHERE 
				id = ?	
			RETURNING
				id
				, company_id
				, subject
				, program_id
				, image_header
				, image_footer
				, email_content
				, template
				, created_at
				, created_by
				, updated_at
				, updated_by					
				, status
	`
	var res []Blast
	err = tx.Select(&res, tx.Rebind(q),
		b.Subject, b.ProgramID, b.ImageHeader, b.ImageFooter, b.EmailContent, b.Template, b.UpdatedBy, b.ID)
	if err != nil {
		return err
	}

	// update blast detail
	// for _, r := range b.RecipientsData {
	// 	q := `
	// 		INSERT INTO blast_recipients(
	// 			blast_id
	// 			, name
	// 			, email
	// 			, voucher_id
	// 			, created_at
	// 			, status
	// 		)
	// 		VALUES (?, ?, ?, ?, ?, ?)
	// 	`

	// 	_, err := tx.Exec(tx.Rebind(q), res[0].ID, r.HolderName, r.HolderEmail, r.VoucherObj.ID, res[0].CreatedAt, StatusCreated)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }

	err = tx.Commit()
	if err != nil {
		return err
	}

	*b = res[0]
	return nil
}

// SendEmailBlast send email blast
func (blast *Blast) SendEmailBlast() error {
	recipients := []Recipient{}

	imageHeader := blast.ImageHeader
	imageVoucher := blast.Program.ImageURL
	imageFooter := blast.ImageFooter

	for _, recipientData := range blast.BlastRecipient {

		data := RecipientRequestData{
			ProgramName:  blast.Program.Name,
			ImageHeader:  imageHeader,
			ImageVoucher: imageVoucher,
			ImageFooter:  imageFooter,
			EmailContent: blast.EmailContent,
			EmailSubject: blast.Subject,
		}

		recipient := Recipient{
			EmailAddress: recipientData.HolderEmail,
			Data:         data,
		}

		recipients = append(recipients, recipient)
	}

	url := "/v3/email/messages?key="
	param := BlastRequest{
		From:     "voucher@elys.id",
		To:       recipients,
		Template: blast.Template,
	}

	jsonParam, _ := json.Marshal(param)

	success, err := mailService("POST", url, jsonParam)
	if err != nil {
		return err
	}

	if success {
		// Update blast status
		err = blast.UpdateBlastStatus()
		if err != nil {
			return errors.New("Failed when update blast status ," + err.Error())
		}
	}

	return nil
}

func (b *Blast) UpdateBlastStatus() error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	q := `UPDATE
				blasts 
			SET
				status = ?		
			WHERE 
				id = ?	
			RETURNING
			id, subject, program_id, image_header, image_footer, email_content, template, created_by, updated_at, updated_by, status
	`
	var res Blasts
	err = tx.Select(&res, tx.Rebind(q), StatusSubmitted, b.ID)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	*b = res[0]
	return nil
}

func mailService(method, url string, param []byte) (bool, error) {
	domain := os.Getenv("MAIL_DOMAIN")
	mailKey := os.Getenv("MAIL_KEY")

	fmt.Printf("url = " + domain + url + mailKey)
	fmt.Printf("%s", param)

	req, err := http.NewRequest(method, domain+url+mailKey, bytes.NewBuffer(param))
	if err != nil {
		panic(err)
		return false, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return false, errors.New(resp.Status)
	}

	return true, nil
}