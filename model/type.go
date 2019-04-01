package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

// VoucherFormat model
type VoucherFormat struct {
	Type       string `json:"type"`
	Properties struct {
		Code    string `json:"code,omitempty"`
		Random  string `json:"random,omitempty"`
		Prefix  string `json:"prefix,omitempty"`
		Postfix string `json:"postfix,omitempty"`
		Length  int    `json:"length,omitempty"`
	} `json:"properties"`
}

//ToString :
func (vf VoucherFormat) ToString() string {
	j, _ := json.Marshal(vf)
	return string(j)
}

func castInterface(dist interface{}, src interface{}) error {
	err := json.Unmarshal([]byte(src.(string)), &dist)
	if err != nil {
		fmt.Println("error:", err)
		return err
	}
	return nil
}

// JSONExpr :
type JSONExpr json.RawMessage

var emptyJSON = JSONExpr("{}")

// MarshalJSON returns the *j as the JSON encoding of j.
func (j JSONExpr) MarshalJSON() ([]byte, error) {
	if len(j) == 0 {
		return emptyJSON, nil
	}
	return j, nil
}

// UnmarshalJSON sets *j to a copy of data
func (j *JSONExpr) UnmarshalJSON(data []byte) error {
	if j == nil {
		return errors.New("JSONExpr: UnmarshalJSON on nil pointer")
	}
	*j = append((*j)[0:0], data...)
	return nil
}

// Value returns j as a value.  This does a validating unmarshal into another
// RawMessage.  If j is invalid json, it returns an error.
func (j JSONExpr) Value() (driver.Value, error) {
	var m json.RawMessage
	var err = j.Unmarshal(&m)
	if err != nil {
		return []byte{}, err
	}
	return []byte(j), nil
}

// Scan stores the src in *j.  No validation is done.
func (j *JSONExpr) Scan(src interface{}) error {
	var source []byte
	switch t := src.(type) {
	case string:
		source = []byte(t)
	case []byte:
		if len(t) == 0 {
			source = emptyJSON
		} else {
			source = t
		}
	case nil:
		*j = emptyJSON
	default:
		return errors.New("Incompatible type for JSONExpr")
	}
	*j = JSONExpr(append((*j)[0:0], source...))
	return nil
}

// Unmarshal unmarshal's the json in j to v, as in json.Unmarshal.
func (j *JSONExpr) Unmarshal(v interface{}) error {
	if len(*j) == 0 {
		*j = emptyJSON
	}
	return json.Unmarshal([]byte(*j), v)
}

// String supports pretty printing for JSONExpr types.
func (j JSONExpr) String() string {
	return string(j)
}

// NullJSONExpr represents a JSONExpr that may be null.
// NullJSONExpr implements the scanner interface so
// it can be used as a scan destination, similar to NullString.
type NullJSONExpr struct {
	JSONExpr
	Valid bool // Valid is true if JSONExpr is not NULL
}

// Scan implements the Scanner interface.
func (n *NullJSONExpr) Scan(value interface{}) error {
	if value == nil {
		n.JSONExpr, n.Valid = emptyJSON, false
		return nil
	}
	n.Valid = true
	return n.JSONExpr.Scan(value)
}

// Value implements the driver Valuer interface.
func (n NullJSONExpr) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.JSONExpr.Value()
}
