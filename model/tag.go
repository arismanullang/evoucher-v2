package model

import (
	"bytes"
	"time"

	"github.com/gilkor/evoucher-v2/util"
)

type (
	//Tag : represent of tags table model
	Tag struct {
		ID          string     `db:"id" json:"id,omitempty"`
		Name        string     `db:"name" json:"name,omitempty"`
		Key         string     `db:"key" json:"key,omitempty"`
		CompanyID   string     `db:"company_id" json:"company_id,omitempty"`
		AccessLevel string     `db:"access_level" json:"access_level,omitempty"`
		CreatedAt   *time.Time `db:"created_at" json:"created_at,omitempty"`
		CreatedBy   string     `db:"created_by" json:"created_by,omitempty"`
		UpdatedAt   *time.Time `db:"updated_at" json:"updated_at,omitempty"`
		UpdatedBy   string     `db:"updated_by" json:"updated_by,omitempty"`
		Status      string     `db:"status" json:"status,omitempty"`
		Action      string     `json:"action,omitempty"`
		Count       int        `db:"count" json:"-"`
	}
	//Tags :
	Tags []Tag
	//ObjectTag tags to n
	ObjectTag struct {
		ID             string     `db:"id" json:"id,omitempty"`
		TagID          string     `db:"tag_id" json:"tag_id,omitempty"`
		ObjectID       string     `db:"object_id" json:"object_id,omitempty" validate:"required"`
		ObjectCategory string     `db:"object_category" json:"object_category,omitempty" validate:"required"`
		CreatedAt      *time.Time `db:"created_at" json:"created_at,omitempty"`
		CreatedBy      string     `db:"created_by" json:"created_by,omitempty" validate:"required"`
		Status         string     `db:"status" json:"status,omitempty"`
		Tags           Tags       `json:"tags,omitempty"  validate:"required"`
	}
	//ObjectTags : list of ObjectTag
	ObjectTags []ObjectTag
)

//TagFields : default table field
var TagFields = []string{"id", "name", "key", "company_id", "access_level", "created_at", "created_by", "updated_at", "updated_by", "status"}

//ObjectTagCategory : type const object tag category
type ObjectTagCategory string

//OTGPrograms : programs
//OTGChannels : channels
//OTGOutlets : outlets
const (
	OTGPrograms ObjectTagCategory = "programs"
	OTGChannels ObjectTagCategory = "channels"
	OTGOutlets  ObjectTagCategory = "outlets"
)

//GetTags : get list company by custom filter
func GetTags(qp *util.QueryParam) (*Tags, bool, error) {
	return getTags("1", "1", qp)
}

//GetTagByCompanyID : get tag by specified ID
func GetTagByCompanyID(qp *util.QueryParam, id string) (*Tags, bool, error) {
	return getTags("company_id", id, qp)
}

//GetTagByID : get tag by specified ID
func GetTagByID(qp *util.QueryParam, id string) (*Tags, bool, error) {
	return getTags("id", id, qp)
}

//GetTagByKey : get tag by specified key
func GetTagByKey(qp *util.QueryParam, key string) (*Tags, bool, error) {
	return getSearchTags("name", util.SimplifyKeyString("%"+key+"%"), qp)
}

//GetTagByCategory : get tag by specified category
func GetTagByCategory(qp *util.QueryParam, val string) (*Tags, bool, error) {
	return getSearchTags("object_category", val, qp)
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

	q = qp.GetQueryWhereClause(q, qp.Q)
	q = qp.GetQueryWithPagination(q, qp.GetQuerySort(), qp.GetQueryLimit())
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
		resd = resd[:qp.Count]
	}
	if len(resd) < qp.Count {
		qp.Count = len(resd)
	}

	return &resd, next, nil
}

func getSearchTags(k, v string, qp *util.QueryParam) (*Tags, bool, error) {
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
			AND ` + k + ` ILIKE ?`

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
func (t *Tag) Insert() (*Tags, error) {
	tx, err := db.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	q := `SELECT id, name, key, company_id, access_level, created_at, created_by, updated_at, updated_by, status
				FROM 
					tags
				WHERE 
					key = ?
					AND company_id = ?
					AND status != ?`
	var res Tags
	err = tx.Select(&res, tx.Rebind(q), util.SimplifyKeyString(t.Name), t.CompanyID, StatusDeleted)
	if err != nil {
		return nil, err
	}
	//if no row return
	if len(res) <= 0 {
		q = `INSERT INTO 
					tags ( name, key, company_id, access_level, created_by, updated_by, status)
				VALUES 
					( ?, ?, ?, ?, ?, ?, ?)
				RETURNING
					id, name, key, company_id, access_level, created_at, created_by, updated_at, updated_by, status
			`
		err = tx.Select(&res, tx.Rebind(q), util.StandardizeSpaces(t.Name), util.SimplifyKeyString(t.Name), t.CompanyID, t.AccessLevel, t.CreatedBy, t.CreatedBy, StatusCreated)
		// util.DEBUG("Lamhoot:", pqerr.Code.Name())
		if err != nil {
			return nil, err
		}
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	*t = res[0]
	return &res, nil
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
				key = ?,
				company_id = ?,
				updated_at = now(),
				updated_by = ?	
			WHERE 
				id = ?	
			RETURNING
			id, name, key, company_id ,created_at, created_by, updated_at, updated_by, status
	`
	var res Tags
	err = tx.Select(&res, tx.Rebind(q), util.StandardizeSpaces(t.Name), util.SimplifyKeyString(t.Name), t.CompanyID, t.UpdatedBy, t.ID)
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
				updated_by = ?,
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

// ObjectTags ###
func getObjectTags(k, v string, qp *util.QueryParam) (*ObjectTags, bool, error) {

	q := qp.GetQueryFields(OutletFields)
	q += `
			FROM
				object_tags object_tags
			WHERE 
				status = ?
			AND ` + k + ` = ?`

	q += qp.GetQuerySort()
	q += qp.GetQueryLimit()
	// fmt.Println(q)
	var resd ObjectTags
	err := db.Select(&resd, db.Rebind(q), StatusCreated, v)
	if err != nil {
		return &ObjectTags{}, false, err
	}
	if len(resd) < 1 {
		return &ObjectTags{}, false, ErrorResourceNotFound
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
func (t *ObjectTag) Insert() (*ObjectTags, error) {

	valuesInsert := new(bytes.Buffer)
	valuesDelete := new(bytes.Buffer)
	var argsInsert []interface{}
	var argsDelete []interface{}
	isInsert := false
	isDelete := false
	// var status string
	//verify existing and deleted data from front guys
	for _, v := range t.Tags {
		if v.Action == "add" {
			isInsert = true
			valuesInsert.WriteString("(?, ?, ?, ?, ?),")
			argsInsert = append(argsInsert, t.ObjectCategory, t.ObjectID, v.ID, t.CreatedBy, StatusCreated)
		} else if v.Action == "remove" {
			isDelete = true
			valuesDelete.WriteString("?,")
			argsDelete = append(argsDelete, v.ID)
		}
	}

	tx, err := db.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var res ObjectTags

	if isDelete {
		q := `DELETE FROM object_tags WHERE tag_id in (`

		valuestr := valuesDelete.String()
		q += valuestr[:len(valuestr)-1]
		q += ")"
		util.DEBUG(q, argsDelete)
		err = tx.Select(&res, tx.Rebind(q), argsDelete...)
		if err != nil {
			return nil, err
		}
	}

	if isInsert {
		q := `INSERT INTO 
				object_tags ( object_category, object_id, tag_id, created_by, status)
				VALUES 
		`
		valuestr := valuesInsert.String()
		q += valuestr[:len(valuestr)-1]
		// ON CONFLICT (object_category, object_id, tag_id, status)
		// 		DO UPDATE SET status = EXCLUDED.status,
		// 					object_category = EXCLUDED.object_category
		q += `  ON CONFLICT DO NOTHING
				RETURNING
					id, tag_id, object_id, object_category, created_by, created_at, status
		`
		util.DEBUG(q, argsInsert)
		err = tx.Select(&res, tx.Rebind(q), argsInsert...)
		if err != nil {
			return nil, err
		}
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	// *t = res[0]
	return &res, nil
}

//Assign :
//nothing wrong here, just do as front guy requested
func (t *ObjectTag) Assign() error {
	values := new(bytes.Buffer)
	var args []interface{}
	// var status string
	//verify existing and deleted data from front guys

	//create tag : NEW

	//delete object tag :

	//assign tag object

	for _, v := range t.Tags {
		values.WriteString("(?, ?, ?, ?, ?),")
		args = append(args, t.ObjectCategory, t.ObjectID, v.ID, t.CreatedBy, StatusCreated)
	}

	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	q := `INSERT INTO 
			object_tags ( object_category, object_id, tag_id, created_by, status)
			VALUES 
	`
	valuestr := values.String()
	q += valuestr[:len(valuestr)-1]
	// ON CONFLICT (object_category, object_id, tag_id, status)
	// 		DO UPDATE SET status = EXCLUDED.status,
	// 					object_category = EXCLUDED.object_category
	q += `  
			RETURNING
				id, tag_id, object_id, object_category, created_by, created_at, status
	`
	var res []ObjectTag
	err = tx.Select(&res, tx.Rebind(q), args...)
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
// func (t *ObjectTag) InsertByHolderID(holder, holderType string, tags []string) error {
// 	tx, err := db.Beginx()
// 	if err != nil {
// 		return err
// 	}
// 	defer tx.Rollback()

// 	q := `INSERT INTO
// 				ObjectTags ( holder_type, holder, tag, created_by, status)
// 			VALUES
// 				( ?, ?, ?, ?, ?)
// 			RETURNING
// 				id, name, company_id, created_at, created_by, updated_at, updated_by, status
// 	`
// 	var res ObjectTags
// 	err = tx.Select(&res, tx.Rebind(q), t.HolderType, holder, t.Tag, t.CreatedBy, StatusCreated)
// 	if err != nil {
// 		return err
// 	}
// 	err = tx.Commit()
// 	if err != nil {
// 		return err
// 	}
// 	*t = res[0]
// 	return nil
// }

//Delete :
func (t *ObjectTag) Delete() error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	q := `DELETE FROM tag_holders where id = ?`
	var res ObjectTags
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
