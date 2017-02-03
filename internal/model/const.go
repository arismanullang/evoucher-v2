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

	VoucherStateCreated string = "created"
	VoucherStateUsed    string = "used"
	VoucherStatePaid    string = "paid"
	VoucherStateDeleted string = "deleted"
)
