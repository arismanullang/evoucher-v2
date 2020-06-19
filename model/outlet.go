package model

import (
	"time"

	"github.com/gilkor/evoucher-v2/util"
	"github.com/jmoiron/sqlx/types"
)

type (
	//Outlet : represent of outlets table model
	Outlet struct {
		ID          string         `db:"id" json:"id,omitempty"`
		Name        string         `db:"name" json:"name,omitempty"`
		Emails      *string        `db:"emails" json:"emails,omitempty"`
		Description types.JSONText `db:"description" json:"description,omitempty"`
		CompanyID   string         `db:"company_id" json:"company_id,omitempty"`
		CreatedAt   *time.Time     `db:"created_at" json:"created_at,omitempty"`
		CreatedBy   string         `db:"created_by" json:"created_by,omitempty"`
		UpdatedAt   *time.Time     `db:"updated_at" json:"updated_at,omitempty"`
		UpdatedBy   string         `db:"updated_by" json:"updated_by,omitempty"`
		Status      string         `db:"status" json:"status,omitempty"`
		Banks       types.JSONText `db:"outlet_banks" json:"banks,omitempty"`
		Tags        types.JSONText `db:"outlet_tags" json:"tags,omitempty"`
		Count       int            `db:"count" json:"-"`
	}
	//Outlets :
	Outlets []Outlet

	OutletsWithTags struct {
		ID          string         `db:"id" json:"id,omitempty"`
		Name        string         `db:"name" json:"name,omitempty"`
		Description types.JSONText `db:"description" json:"description,omitempty"`
		CompanyID   string         `db:"company_id" json:"company_id,omitempty"`
		CreatedAt   *time.Time     `db:"created_at" json:"created_at,omitempty"`
		CreatedBy   string         `db:"created_by" json:"created_by,omitempty"`
		UpdatedAt   *time.Time     `db:"updated_at" json:"updated_at,omitempty"`
		UpdatedBy   string         `db:"updated_by" json:"updated_by,omitempty"`
		Status      string         `db:"status" json:"status,omitempty"`
		Banks       types.JSONText `db:"outlet_banks" json:"banks,omitempty"`
		Tags        types.JSONText `db:"outlet_tags" json:"tags,omitempty"`
	}
)

//OutletFields : default table field
var OutletFields = []string{"id", "name", "description", "created_at", "created_by", "updated_at", "updated_by", "status"}

//MOutletFields : fields for 3rd party api
var MOutletFields = "id,name,description"

//GetOutlets : get list outlet by custom filter
func GetOutlets(qp *util.QueryParam) (*Outlets, bool, error) {
	return getOutlets("1", "1", qp)
}

//GetOutletByID : get outlet by specified ID
func GetOutletByID(qp *util.QueryParam, id string) (*Outlet, bool, error) {
	outlets, _, err := getOutlets("id", id, qp)
	if err != nil {
		return &Outlet{}, false, err
	}

	if len(*outlets) > 0 {
		outlet := (*outlets)[0]
		return &outlet, false, nil
	}

	return &Outlet{}, false, ErrorResourceNotFound
}

func getOutlets(k, v string, qp *util.QueryParam) (*Outlets, bool, error) {

	q, err := qp.GetQueryByDefaultStruct(Outlet{})
	if err != nil {
		return &Outlets{}, false, err
	}

	q += `
			FROM
				m_outlets outlet
			WHERE 
				status = ?
			AND ` + k + ` = ?`

	q = qp.GetQueryWhereClause(q, qp.Q)
	q = qp.GetQueryWithPagination(q, qp.GetQuerySort(), qp.GetQueryLimit())
	// fmt.Println(q)
	util.DEBUG("query struct :", q)
	// query := "select row_to_json(row) from (" + q + ") row"
	var resd Outlets
	err = db.Select(&resd, db.Rebind(q), StatusCreated, v)
	if err != nil {
		return &Outlets{}, false, err
	}
	if len(resd) < 1 {
		return &Outlets{}, false, ErrorResourceNotFound
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

//GetOutletsByTags : get outlet by tag.id
func GetOutletsByTags(qp *util.QueryParam, v string) (*Outlets, bool, error) {
	q, err := qp.GetQueryByDefaultStruct(Outlet{})
	if err != nil {
		return &Outlets{}, false, err
	}

	q += `
			FROM
				m_outlets outlet,
				object_tags object_tag
			WHERE 
				outlet.status = ?
			AND object_tag.status = ?
			AND outlet.id = object_tag.object_id
			AND object_tag.tag_id = ?`

	q = qp.GetQueryWithPagination(q, qp.GetQuerySort(), qp.GetQueryLimit())
	// fmt.Println(q)
	util.DEBUG("query struct :", q)
	var resd Outlets
	err = db.Select(&resd, db.Rebind(q), StatusCreated, StatusCreated, v)
	if err != nil {
		return &Outlets{}, false, err
	}
	if len(resd) < 1 {
		return &Outlets{}, false, ErrorResourceNotFound
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

//Insert : save data to database
func (p *Outlet) Insert() (*Outlets, error) {
	tx, err := db.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	q := `INSERT INTO 
				outlets ( name, description, emails, company_id, created_by, updated_by, status)
			VALUES 
				( ?, ?, ?, ?, ?, ?, ?)
			RETURNING
				id, name, description, emails, company_id, created_at, created_by, updated_at, updated_by, status
	`
	// bank, err := json.Marshal(p.Bank)
	// if err != nil {
	// 	return nil, err
	// }
	var res Outlets
	util.DEBUG(q)
	err = tx.Select(&res, tx.Rebind(q), p.Name, p.Description, p.Emails, p.CompanyID, p.CreatedBy, p.CreatedBy, StatusCreated)
	if err != nil {
		util.DEBUG(`la1-->`, err)
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		util.DEBUG(`la2-->`, err)
		return nil, err
	}
	*p = res[0]
	return &res, nil
}

//Update : modify data
func (p *Outlet) Update() error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	q := `UPDATE
				outlets 
			SET
				name = ?,
				emails = ?,
				description = ?,
				updated_at = now(),
				updated_by = ?				
			WHERE 
				id = ?	
			RETURNING
			id, name, emails, created_at, created_by, updated_at, updated_by, status
	`
	var res []Outlet
	err = tx.Select(&res, tx.Rebind(q), p.Name, p.Emails, p.Description, p.UpdatedBy, p.ID)
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
func (p *Outlet) Delete() error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	q := `UPDATE
				outlets 
			SET
				updated_at = now(),
				updated_by = ?
				status = ?			
			WHERE 
				id = ?	
			RETURNING
			id, name, description, created_at, created_by, updated_at, updated_by, status
	`
	var res []Outlet
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
// func (p *Outlets) DecodeDescription(i interface{}) *Outlets {
// 	data := make(Outlets, len(*p))
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
