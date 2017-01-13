package model

import (
	"errors"
)

var (
	ErrResourceNotFound = errors.New("resource not found.")
	ErrValidationError  = errors.New("validation error.")
	ErrInvalidPassword  = errors.New("invalid password.")
	ErrNotModified      = errors.New("data not modified.")
)

const (
	StatusCreated string = "created"
	StatusDeleted string = "deleted"

	TransactionStateOutstanding string = "outstanding"
	TransactionStatePaid        string = "paid"
	TransactionStateVoided      string = "voided"
)
