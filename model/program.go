package model

import (
	"errors"
	"time"

	"github.com/gilkor/evoucher-v2/util"
	"github.com/jmoiron/sqlx/types"
)

type (
	//Program : base model
	Program struct {
		ID              string         `db:"id" json:"id,omitempty"`
		CompanyID       string         `db:"company_id" json:"company_id,omitempty"`
		Name            string         `db:"name" json:"name,omitempty"`
		Type            string         `db:"type" json:"type,omitempty"`
		Value           float64        `db:"value" json:"value,omitempty"`
		MaxValue        float64        `db:"max_value" json:"max_value,omitempty"`
		StartDate       *time.Time     `db:"start_date" json:"start_date,omitempty"`
		EndDate         *time.Time     `db:"end_date" json:"end_date,omitempty"`
		Description     types.JSONText `db:"description" json:"description,omitempty"`
		ImageURL        string         `db:"image_url" json:"image_url,omitempty"`
		Template        string         `db:"template" json:"template,omitempty"`
		Rule            types.JSONText `db:"rule" json:"rule,omitempty"`
		State           string         `db:"state" json:"state,omitempty"`
		Stock           int64          `db:"stock" json:"stock,omitempty"`
		CreatedAt       *time.Time     `db:"created_at" json:"created_at,omitempty"`
		CreatedBy       string         `db:"created_by" json:"created_by,omitempty"`
		UpdatedAt       *time.Time     `db:"updated_at" json:"updated_at,omitempty"`
		UpdatedBy       string         `db:"updated_by" json:"updated_by,omitempty"`
		Status          string         `db:"status" json:"status,omitempty"`
		Outlets         Outlets        `json:"outlets,omitempty"`
		Vouchers        Vouchers       `json:"vouchers,omitempty"`
		VoucherFormat   types.JSONText `db:"voucher_format" json:"voucher_format,omitempty"`
		IsReimburse     bool           `db:"is_reimburse" json:"is_reimburse,omitempty"`
		Price           float64        `db:"price" json:"price,omitempty"`
		ChannelID       string         `db:"channel_id" json:"channel_id,omitempty"`
		Channels        Channels       `json:"channels,omitempty"`
		ProgramChannels types.JSONText `db:"program_channels" json:"program_channels,omitempty"`
		Count           int            `db:"count" json:"-"`
		ClaimedVoucher  int            `db:"claimed" json:"claimed"`
		UsedVoucher     int            `db:"used" json:"used"`
		PaidVoucher     int            `db:"paid" json:"paid,omitempty"`
		// WithTransactionCount bool       `json:"with_transaction_count,omitempty"`
	}
	// Programs : base model
	Programs []Program
)

//MProgramFields : fields for 3rd party api
var MProgramFields = "id, name, type, value, max_value, rule, start_date, end_date, description, image_url, stock, is_reimburse, claimed, used"

// GetProgramByHolder :
func GetProgramByHolder(id string, qp *util.QueryParam) (*Programs, bool, error) {
	return getPrograms("holder", id, qp)
}

// GetProgramByID :  program details
func GetProgramByID(id string, qp *util.QueryParam) (*Program, error) {

	programs, _, err := getPrograms("id", id, qp)
	if err != nil {
		return &Program{}, errors.New("Failed when select on program ," + err.Error())
	}
	if len(*programs) < 1 {
		return &Program{}, nil
	}
	program := &(*programs)[0]

	return program, nil
}

// GetPrograms : get program list
func GetPrograms(qp *util.QueryParam) (*Programs, bool, error) {
	return getPrograms("1", "1", qp)
}

func getPrograms(key, value string, qp *util.QueryParam) (*Programs, bool, error) {
	q, err := qp.GetQueryByDefaultStruct(Program{})
	if err != nil {
		return &Programs{}, false, errors.New("Failed when select on program ," + err.Error())
	}
	q += `
			FROM
				m_programs program
			WHERE 
				status = ?
			AND ` + key + ` = ?`

	q = qp.GetQueryWhereClause(q, qp.Q)
	q = qp.GetQueryWithPagination(q, qp.GetQuerySort(), qp.GetQueryLimit())
	util.DEBUG(q)
	var resd Programs
	err = db.Select(&resd, db.Rebind(q), StatusCreated, value)
	if err != nil {
		return &Programs{}, false, errors.New("Failed when select on program ," + err.Error())
	}
	if len(resd) < 1 {
		return &Programs{}, false, nil
	}

	next := false
	if len(resd) > qp.Count && qp.Count > 0 {
		next = true
		resd = resd[:qp.Count]
	}
	if len(resd) < qp.Count {
		qp.Count = len(resd)
	}

	return &resd, next, nil
}

//Insert : single row inset into table
func (p *Program) Insert() (*Programs, error) {
	tx, err := db.Beginx()
	if err != nil {
		return nil, errors.New("Failed when insert new program ," + err.Error())
	}
	defer tx.Rollback()

	q := `INSERT INTO 
				programs 
				( 
					company_id
					, name
					, type
					, value
					, price
					, max_value
					, start_date
					, end_date
					, description
					, image_url
					, template
					, rule
					, state
					, stock
					, is_reimburse
					, channel_id
					, voucher_format
					, created_by
					, updated_by					
					, status
				)
			VALUES 
				( ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			RETURNING 
			id
			, company_id
			, name
			, type
			, value
			, price
			, max_value
			, start_date
			, end_date
			, description
			, image_url
			, template
			, rule
			, state
			, stock
			, is_reimburse
			, channel_id
			, voucher_format
			, created_at
			, created_by
			, updated_at
			, updated_by
			, status
			, (SELECT array_to_json(array_agg(row_to_json(c.*))) AS array_to_json
			FROM channels c
			WHERE (((c.id)::text = (programs.channel_id)::text) AND (c.status = 'created'::status))) AS program_channels
	`

	util.DEBUG(q)
	var res Programs
	channelID := p.ChannelID
	if channelID == "" && (len(p.Channels) > 0) {
		channelID = p.Channels[0].ID

		util.DEBUG(channelID)
		err = tx.Select(&res, tx.Rebind(q), p.CompanyID, p.Name, p.Type, p.Value, p.Price, p.MaxValue, p.StartDate, p.EndDate,
			p.Description, p.ImageURL, p.Template, util.StandardizeSpaces(p.Rule.String()), p.State, p.Stock, p.IsReimburse, channelID,
			util.StandardizeSpaces(p.VoucherFormat.String()), p.CreatedBy, p.CreatedBy, StatusCreated)
		if err != nil {
			return nil, err
		}
	}

	//update returning ID into obj program
	p.ID = res[0].ID
	//insert program outlets
	if err := NewProgramOutlets(p.ID, p.Outlets).Upsert(tx); err != nil {
		return nil, errors.New("Failed when insert new outlets ," + err.Error())
	}

	err = tx.Commit()
	if err != nil {
		return nil, errors.New("Failed when insert new program ," + err.Error())
	}

	return &res, nil
}

//Update : update program
func (p *Program) Update() error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	q := `UPDATE
				programs 
			SET
					company_id = ?
					, name = ?
					, type = ?
					, value = ?
					, price = ?
					, max_value = ?
					, start_date = ?
					, end_date = ?
					, description = ?
					, image_url = ?
					, template = ?
					, rule = ?
					, state = ?
					, stock = ?
					, is_reimburse = ?
					, updated_at = now()
					, updated_by = ?		
			WHERE 
				id = ?
			RETURNING
				id
				, company_id
				, name
				, type
				, value
				, price
				, max_value
				, start_date
				, end_date
				, description
				, image_url
				, template
				, rule
				, state
				, stock
				, is_reimburse
				, channel_id
				, voucher_format
				, created_at
				, created_by
				, updated_at
				, updated_by
				, status
				, (SELECT array_to_json(array_agg(row_to_json(c.*))) AS array_to_json
				FROM channels c
				WHERE (((c.id)::text = (programs.channel_id)::text) AND (c.status = 'created'::status))) AS program_channels
	`
	util.DEBUG(q)
	var res []Program
	err = tx.Select(&res, tx.Rebind(q), p.CompanyID, p.Name, p.Type, p.Value, p.Price, p.MaxValue, p.StartDate, p.EndDate,
		p.Description, p.ImageURL, p.Template, p.Rule, p.State, p.Stock, p.IsReimburse, p.UpdatedBy, p.ID)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

//Delete : soft delated data by updateting row status to "deleted"
func (p *Program) Delete() error {

	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	//delete outlets
	if err := NewProgramOutlets(p.ID, p.Outlets).Delete(tx); err != nil {
		return errors.New("Failed when delete outlets ," + err.Error())
	}

	q := `UPDATE
				customers 
			SET
				updated_at = now(),
				updated_by = ?,
				status = ?			
			WHERE 
				id = ?
			RETURNING 
				id			
	`
	var res []string
	err = tx.Select(&res, tx.Rebind(q), p.UpdatedBy, StatusDeleted, p.ID)
	if err != nil {
		return errors.New("Failed when delete program ," + ErrorNoDataAffected.Error() + " , " + err.Error())
	}
	err = tx.Commit()
	if err != nil {
		return errors.New("Failed when delete program ," + err.Error())
	}

	return nil
}
