package model

import (
	"time"

	"github.com/gilkor/evoucher/util"
)

type (
	//Tag : represent of tags table model
	Tag struct {
		ID        string     `db:"id" json:"id,omitempty"`
		Name      string     `db:"name" json:"name,omitempty"`
		CompanyID string     `db:"company_id" json:"company_id"`
		CreatedAt *time.Time `db:"created_at" json:"created_at,omitempty"`
		CreatedBy string     `db:"created_by" json:"created_by,omitempty"`
		UpdatedAt *time.Time `db:"updated_at" json:"updated_at,omitempty"`
		UpdatedBy string     `db:"updated_by" json:"updated_by,omitempty"`
		Status    string     `db:"status" json:"status,omitempty"`
	}
	//Tags :
	Tags []Tag
	//TagHolder tags to n
	TagHolder struct {
		ID         int        `db:"id" json:"id,omitempty"`
		HolderType string     `db:"holder_type" json:"holder_type,omitempty"`
		Holder     string     `db:"holder" json:"holder,omitempty"`
		Tag        string     `db:"tag" json:"tag,omitempty"`
		CreatedAt  *time.Time `db:"created_at" json:"created_at,omitempty"`
		CreatedBy  string     `db:"created_by" json:"created_by,omitempty"`
		Status     string     `db:"status" json:"status,omitempty"`
	}
	//TagHolders : list of TagHolder
	TagHolders []TagHolder
)

//TagFields : default table field
var TagFields = []string{"id", "name", "company_id", "created_at", "created_by", "updated_at", "updated_by", "status"}

//GetTags : get list company by custom filter
func GetTags(qp *util.QueryParam) (*Tags, bool, error) {
	return getTags("1", "1", qp)
}

//GetTagByCompanyID : get partner by specified ID
func GetTagByCompanyID(qp *util.QueryParam, id string) (*Tags, bool, error) {
	return getTags("company_id", id, qp)
}

//GetTagByID : get partner by specified ID
func GetTagByID(qp *util.QueryParam, id string) (*Tags, bool, error) {
	return getTags("id", id, qp)
}

func getTags(k, v string, qp *util.QueryParam) (*Tags, bool, error) {
	q, err := qp.GetQueryByDefaultStruct(Tag{})
	if err != nil {
		return &Tags{}, false, err
	}
	// q := qp.GetQueryFields(TagFields)
	q += `
			FROM
				tags tag
			WHERE 
				status = ?
			AND ` + k + ` = ?`

	q += qp.GetQuerySort()
	q += qp.GetQueryLimit()
	// fmt.Println(q)
	var resd Tags
	err = db.Select(&resd, db.Rebind(q), StatusCreated, v)
	if err != nil {
		return &Tags{}, false, err
	}
	if len(resd) < 1 {
		return &Tags{}, false, ErrorResourceNotFound
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
func (t *Tag) Insert() error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	q := `INSERT INTO 
				tags ( name, company_id, created_by, updated_by, status)
			VALUES 
				( ?, ?, ?, ?, ?)
			RETURNING
				id, name, company_id, created_at, created_by, updated_at, updated_by, status
	`
	var res Tags
	err = tx.Select(&res, tx.Rebind(q), t.Name, t.CompanyID, t.CreatedBy, t.CreatedBy, StatusCreated)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	*t = res[0]
	return nil
}

//Update : modify data
func (t *Tag) Update() error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	q := `UPDATE
				tags 
			SET
				name = ?,
				company_id = ?,
				updated_at = now(),
				updated_by = ?				
			WHERE 
				id = ?	
			RETURNING
			id, name, company_id ,created_at, created_by, updated_at, updated_by, status
	`
	var res Tags
	err = tx.Select(&res, tx.Rebind(q), t.Name, t.CompanyID, t.UpdatedBy, t.ID)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	*t = res[0]
	return nil
}

//Delete : soft delated data by updateting row status to "deleted"
func (t *Tag) Delete() error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	q := `UPDATE
				tags 
			SET
				updated_at = now(),
				updated_by = ?
				status = ?			
			WHERE 
				id = ?	
			RETURNING
			id, name, company_id, created_at, created_by, updated_at, updated_by, status
	`
	var res Tags
	err = tx.Select(&res, tx.Rebind(q), t.UpdatedBy, StatusDeleted, t.ID)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// TagHolders ###
func getTagHolders(k, v string, qp *util.QueryParam) (*TagHolders, bool, error) {

	// q := qp.GetQueryByDefaultStruct(Partner{})
	q := qp.GetQueryFields(PartnerFields)
	q += `
			FROM
				tag_holders tag_holders
			WHERE 
				status = ?
			AND ` + k + ` = ?`

	q += qp.GetQuerySort()
	q += qp.GetQueryLimit()
	// fmt.Println(q)
	var resd TagHolders
	err := db.Select(&resd, db.Rebind(q), StatusCreated, v)
	if err != nil {
		return &TagHolders{}, false, err
	}
	if len(resd) < 1 {
		return &TagHolders{}, false, ErrorResourceNotFound
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
func (t *TagHolder) Insert() error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	q := `INSERT INTO 
				TagHolders ( holder_type, holder, tag, created_by, status)
			VALUES 
				( ?, ?, ?, ?, ?)
			RETURNING
				id, name, company_id, created_at, created_by, updated_at, updated_by, status
	`
	var res TagHolders
	err = tx.Select(&res, tx.Rebind(q), t.HolderType, t.Holder, t.Tag, t.CreatedBy, StatusCreated)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	*t = res[0]
	return nil
}

//InsertByHolderID : save data to database
func (t *TagHolder) InsertByHolderID(holder, holderType string, tags []string) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	q := `INSERT INTO 
				TagHolders ( holder_type, holder, tag, created_by, status)
			VALUES 
				( ?, ?, ?, ?, ?)
			RETURNING
				id, name, company_id, created_at, created_by, updated_at, updated_by, status
	`
	var res TagHolders
	err = tx.Select(&res, tx.Rebind(q), t.HolderType, holder, t.Tag, t.CreatedBy, StatusCreated)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	*t = res[0]
	return nil
}

//Delete :
func (t *TagHolder) Delete() error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	q := `DELETE FROM tag_holders where id = ?`
	var res TagHolders
	err = tx.Select(&res, tx.Rebind(q), t.ID)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
