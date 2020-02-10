package model

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/gilkor/evoucher-v2/util"
)

type (

	// Config : represent of config table model
	Config struct {
		CompanyID string     `db:"company_id" json:"company_id,omitempty"`
		Category  string     `db:"category" json:"category,omitempty"`
		Key       string     `db:"key" json:"key,omitempty"`
		Value     string     `db:"value" json:"value,omitempty"`
		CreatedAt *time.Time `db:"created_at" json:"created_at,omitempty"`
		CreatedBy string     `db:"created_by" json:"created_by,omitempty"`
		UpdatedAt *time.Time `db:"updated_at" json:"updated_at,omitempty"`
		UpdatedBy string     `db:"updated_by" json:"updated_by,omitempty"`
		Status    string     `db:"status" json:"status,omitempty"`
	}

	// Configs : List of Config
	Configs []Config
)

// GetConfigs : get list config by custom filter
func GetConfigs(companyID, category string) (map[string]interface{}, error) {

	q := `
			SELECT 
				key,
				value
			FROM
				config
			WHERE 
				status = ?
			AND
				company_id = ?
			AND
				category = ?
			ORDER BY created_at`

	util.DEBUG(q)

	// err := db.Select(&res, db.Rebind(q), StatusCreated, v)
	rows, err := db.Query(db.Rebind(q), StatusCreated, companyID, category)

	cols, err := rows.Columns()
	if err != nil {
		fmt.Println(err.Error())
		return map[string]interface{}{}, err
	}
	defer rows.Close()
	index := 0
	// var result []map[string]interface{}
	m := make(map[string]interface{})

	for rows.Next() {
		pointer := make([]interface{}, len(cols))
		pointers := make([]interface{}, len(cols))
		for i, _ := range cols {
			pointers[i] = &pointer[i]
		}
		err := rows.Scan(pointers...)
		if err != nil {
			log.Panic(err)
		}

		key := fmt.Sprintf("%v", pointer[0])
		value := pointer[1]
		m[key] = value

		index++
	}

	if err != nil {
		return map[string]interface{}{}, err
	}
	if len(m) < 1 {
		return map[string]interface{}{}, ErrorResourceNotFound
	}

	return m, nil
}

//Insert : single row inset into table
func (c *Config) Insert() (*Configs, error) {
	tx, err := db.Beginx()
	if err != nil {
		return nil, errors.New("Failed when insert new config ," + err.Error())
	}
	defer tx.Rollback()

	q := `INSERT INTO 
				config
				( 
					company_id
					, category
					, key
					, value
					, created_by
					, updated_by					
					, status
				)
			VALUES 
				( ?, ?, ?, ?, ?, ?, ?)
			RETURNING 
			company_id
			, category
			, key
			, value
			, created_at
			, created_by
			, updated_at
			, updated_by					
			, status
	`

	util.DEBUG(q)
	var res Configs

	err = tx.Select(&res, tx.Rebind(q), c.CompanyID, c.Category, c.Key, c.Value, c.CreatedBy, c.UpdatedBy, StatusCreated)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	*c = res[0]
	return &res, nil
}

//Update : modify data
func (c *Config) Update() error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	q := `UPDATE
				config 
			SET
				value = ?
				, updated_at = now()
				, updated_by = ?				
			WHERE 
				company_id = ?
				AND key = ?
			RETURNING
				, company_id
				, category
				, key
				, value
				, created_at
				, created_by
				, updated_at
				, updated_by					
				, status
	`
	var res Configs
	err = tx.Select(&res, tx.Rebind(q),
		c.Value, c.UpdatedBy, c.CompanyID, c.Key)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	*c = res[0]
	return nil
}

//Delete : soft delated data by updating row status to "deleted"
func (c *Config) Delete() error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	q := `UPDATE
				config 
			SET
				updated_at = now(),
				updated_by = ?,
				status = ?			
			WHERE 
				company_id = ?	
				AND key = ?
			RETURNING
			company_id, category, key, value, created_at, created_by, updated_at, updated_by, status
	`
	var res Configs
	err = tx.Select(&res, tx.Rebind(q), c.UpdatedBy, StatusDeleted, c.CompanyID, c.Key)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
