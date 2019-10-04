package model

import (
	"time"

	"github.com/gilkor/evoucher-v2/util"
)

type (
	//Partner : represent of partners table model
	Partner struct {
		ID          string     `db:"id" json:"id,omitempty"`
		Name        string     `db:"name" json:"name,omitempty"`
		Description JSONExpr   `db:"description" json:"description,omitempty"`
		CompanyID   string     `db:"company_id" json:"company_id,omitempty"`
		CreatedAt   *time.Time `db:"created_at" json:"created_at,omitempty"`
		CreatedBy   string     `db:"created_by" json:"created_by,omitempty"`
		UpdatedAt   *time.Time `db:"updated_at" json:"updated_at,omitempty"`
		UpdatedBy   string     `db:"updated_by" json:"updated_by,omitempty"`
		Status      string     `db:"status" json:"status,omitempty"`
		Bank        *Bank      `json:"bank,omitempty"`
		Tags        *Tags      `json:"tags,omitempty"`
	}
	//Partners :
	Partners []Partner

	PartnersWithTags struct {
		ID          string     `db:"id" json:"id,omitempty"`
		Name        string     `db:"name" json:"name,omitempty"`
		Description JSONExpr   `db:"description" json:"description,omitempty"`
		CompanyID   string     `db:"company_id" json:"company_id,omitempty"`
		CreatedAt   *time.Time `db:"created_at" json:"created_at,omitempty"`
		CreatedBy   string     `db:"created_by" json:"created_by,omitempty"`
		UpdatedAt   *time.Time `db:"updated_at" json:"updated_at,omitempty"`
		UpdatedBy   string     `db:"updated_by" json:"updated_by,omitempty"`
		Status      string     `db:"status" json:"status,omitempty"`
		Bank        *Bank      `json:"bank,omitempty"`
		Tags        *Tags      `json:"tags,omitempty"`
	}

	//Bank :
	/**
	Company Name
	Person In Charge
	Contact Number
	Company Email
	Bank Name
	Bank Branch
	Bank Account Number
	Bank Account Holder
	**/
	Bank struct {
		Bank            string `json:"bank,omitempty"`
		BankAccount     string `json:"bank_account,omitempty"`
		BankAccountName string `json:"bank_acount_name,omitempty"`
	}
)

//PartnerFields : default table field
var PartnerFields = []string{"id", "name", "description", "created_at", "created_by", "updated_at", "updated_by", "status"}

//GetPartners : get list company by custom filter
func GetPartners(qp *util.QueryParam) (*Partners, bool, error) {
	return getPartners("1", "1", qp)
}

//GetPartnerByID : get partner by specified ID
func GetPartnerByID(qp *util.QueryParam, id string) (*Partners, bool, error) {
	return getPartners("id", id, qp)
}

func getPartners(k, v string, qp *util.QueryParam) (*Partners, bool, error) {

	q, err := qp.GetQueryByDefaultStruct(Partner{})
	if err != nil {
		return &Partners{}, false, err
	}
	// q := qp.GetQueryFields(PartnerFields)

	q += `
			FROM
				partners partner
			WHERE 
				status = ?
			AND ` + k + ` = ?`

	q += qp.GetQuerySort()
	q += qp.GetQueryLimit()
	// fmt.Println(q)
	util.DEBUG("query struct :", q)
	var resd Partners
	err = db.Select(&resd, db.Rebind(q), StatusCreated, v)
	if err != nil {
		return &Partners{}, false, err
	}
	if len(resd) < 1 {
		return &Partners{}, false, ErrorResourceNotFound
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

//GetPartnersByTags : get partner by tag.id
func GetPartnersByTags(qp *util.QueryParam, v string) (*Partners, bool, error) {

	q, err := qp.GetQueryByDefaultStruct(Partner{})
	if err != nil {
		return &Partners{}, false, err
	}
	// q := qp.GetQueryFields(PartnerFields)

	q += `
			FROM
				partners Partner,
				tag_holders t
			WHERE 
				partner.status = ?
			AND t.status = ?
			AND partner.id = t.holder
			AND t.tag = ?`

	q += qp.GetQuerySort()
	q += qp.GetQueryLimit()
	// fmt.Println(q)
	util.DEBUG("query struct :", q)
	var resd Partners
	err = db.Select(&resd, db.Rebind(q), StatusCreated, StatusCreated, v)
	if err != nil {
		return &Partners{}, false, err
	}
	if len(resd) < 1 {
		return &Partners{}, false, ErrorResourceNotFound
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

//Insert : save data to database
func (p *Partner) Insert() error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	q := `INSERT INTO 
				partners ( name, description, company_id, created_by, updated_by, status)
			VALUES 
				( ?, ?, ?, ?, ?, ?)
			RETURNING
				id, name, description, company_id, created_at, created_by, updated_at, updated_by, status
	`
	var res []Partner
	err = tx.Select(&res, tx.Rebind(q), p.Name, p.Description, p.CompanyID, p.CreatedBy, p.CreatedBy, StatusCreated)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	*p = res[0]
	return nil
}

//Update : modify data
func (p *Partner) Update() error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	q := `UPDATE
				partners 
			SET
				name = ?,
				description = ?,
				updated_at = now(),
				updated_by = ?				
			WHERE 
				id = ?	
			RETURNING
			id, name, created_at, created_by, updated_at, updated_by, status
	`
	var res []Partner
	err = tx.Select(&res, tx.Rebind(q), p.Name, p.Description, p.UpdatedBy, p.ID)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	*p = res[0]
	return nil
}

//Delete : soft delated data by updateting row status to "deleted"
func (p *Partner) Delete() error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	q := `UPDATE
				partners 
			SET
				updated_at = now(),
				updated_by = ?
				status = ?			
			WHERE 
				id = ?	
			RETURNING
			id, name, description, created_at, created_by, updated_at, updated_by, status
	`
	var res []Partner
	err = tx.Select(&res, tx.Rebind(q), p.UpdatedBy, StatusDeleted, p.ID)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// //DecodeDescription :
// func (p *Partners) DecodeDescription(i interface{}) *Partners {
// 	data := make(Partners, len(*p))
// 	for k, v := range *p {
// 		data[k].ID = v.ID
// 		data[k].Name = v.Name
// 		data[k].Description = v.Description
// 		data[k].CreatedAt = v.CreatedAt
// 		data[k].CreatedBy = v.CreatedBy
// 		data[k].UpdatedAt = v.UpdatedAt
// 		data[k].UpdatedAt = v.UpdatedAt
// 		data[k].Status = v.Status

// 		if v.Description != nil {
// 			desc := v.Description.Unmarshal(&i)
// 			data[k].Bank = i
// 		}
// 	}
// 	// copy(data, *p)
// 	return &data
// }
