package model

import (
	"time"

	"github.com/gilkor/evoucher-v2/util"
	"github.com/jmoiron/sqlx/types"
)

type (
	//Channel : represent of channels table model
	Channel struct {
		ID          string         `db:"id" json:"id,omitempty"`
		Name        string         `db:"name" json:"name,omitempty"`
		Description JSONExpr       `db:"description" json:"description,omitempty"`
		IsSuper     bool           `db:"is_super" json:"is_super,omitempty"`
		CreatedAt   *time.Time     `db:"created_at" json:"created_at,omitempty"`
		CreatedBy   string         `db:"created_by" json:"created_by,omitempty"`
		UpdatedAt   *time.Time     `db:"updated_at" json:"updated_at,omitempty"`
		UpdatedBy   string         `db:"updated_by" json:"updated_by,omitempty"`
		Status      string         `db:"status" json:"status,omitempty"`
		Tags        types.JSONText `db:"channel_tags" json:"tags,omitempty"`
	}
	//Channels :
	Channels []Channel

	// ChannelsWithTags struct {
	// 	ID          string         `db:"id" json:"id,omitempty"`
	// 	Name        string         `db:"name" json:"name,omitempty"`
	// 	Description JSONExpr       `db:"description" json:"description,omitempty"`
	// 	CompanyID   string         `db:"company_id" json:"company_id,omitempty"`
	// 	CreatedAt   *time.Time     `db:"created_at" json:"created_at,omitempty"`
	// 	CreatedBy   string         `db:"created_by" json:"created_by,omitempty"`
	// 	UpdatedAt   *time.Time     `db:"updated_at" json:"updated_at,omitempty"`
	// 	UpdatedBy   string         `db:"updated_by" json:"updated_by,omitempty"`
	// 	Status      string         `db:"status" json:"status,omitempty"`
	// 	Banks       types.JSONText `db:"channel_banks" json:"banks,omitempty"`
	// 	Tags        types.JSONText `db:"channel_tags" json:"tags,omitempty"`
	// }
)

//ChannelFields : default table field
var ChannelFields = []string{"id", "name", "description", "is_super", "created_at", "created_by", "updated_at", "updated_by", "status"}

//GetChannels : get list company by custom filter
func GetChannels(qp *util.QueryParam) (*Channels, bool, error) {
	return getChannels("1", "1", qp)
}

//GetChannelByID : get channel by specified ID
func GetChannelByID(qp *util.QueryParam, id string) (*Channels, bool, error) {
	return getChannels("id", id, qp)
}

func getChannels(k, v string, qp *util.QueryParam) (*Channels, bool, error) {

	q, err := qp.GetQueryByDefaultStruct(Channel{})
	if err != nil {
		return &Channels{}, false, err
	}
	// q := qp.GetQueryFields(ChannelFields)

	q += `
			FROM
				m_channels channel
			WHERE 
				status = ?
			AND ` + k + ` = ?`

	q += qp.GetQuerySort()
	q += qp.GetQueryLimit()
	var resd Channels
	err = db.Select(&resd, db.Rebind(q), StatusCreated, v)
	if err != nil {
		return &Channels{}, false, err
	}
	if len(resd) < 1 {
		return &Channels{}, false, ErrorResourceNotFound
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

//GetChannelsByTags : get channel by tag.id
func GetChannelsByTags(qp *util.QueryParam, v string) (*Channels, bool, error) {

	q, err := qp.GetQueryByDefaultStruct(Channel{})
	if err != nil {
		return &Channels{}, false, err
	}
	// q := qp.GetQueryFields(ChannelFields)

	q += `
			FROM
				channels Channel,
				tag_holders t
			WHERE 
				channel.status = ?
			AND t.status = ?
			AND channel.id = t.holder
			AND t.tag = ?`

	q += qp.GetQuerySort()
	q += qp.GetQueryLimit()
	// fmt.Println(q)
	util.DEBUG("query struct :", q)
	var resd Channels
	err = db.Select(&resd, db.Rebind(q), StatusCreated, StatusCreated, v)
	if err != nil {
		return &Channels{}, false, err
	}
	if len(resd) < 1 {
		return &Channels{}, false, ErrorResourceNotFound
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
func (p *Channel) Insert() error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	q := `INSERT INTO 
				channels ( name, description, is_super, created_by, updated_by, status)
			VALUES 
				( ?, ?, ?, ?, ?, ?)
			RETURNING
				id, name, description, is_super, created_at, created_by, updated_at, updated_by, status
	`
	// bank, err := json.Marshal(p.Bank)
	if err != nil {
		return err
	}
	var res []Channel
	// util.DEBUG(p.Bank)
	err = tx.Select(&res, tx.Rebind(q), p.Name, p.Description, p.IsSuper, p.CreatedBy, p.CreatedBy, StatusCreated)
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
func (p *Channel) Update() error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	q := `UPDATE
				channels 
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
	var res []Channel
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
func (p *Channel) Delete() error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	q := `UPDATE
				channels 
			SET
				updated_at = now(),
				updated_by = ?,
				status = ?			
			WHERE 
				id = ?	
			RETURNING
			id, name, description, is_super, created_at, created_by, updated_at, updated_by, status
	`
	var res []Channel
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
