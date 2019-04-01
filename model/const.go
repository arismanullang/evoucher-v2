package model

import "errors"

const (
	//StatusCreated row satus "created"
	StatusCreated = "created"
	//StatusDeleted row satus "deleted"
	StatusDeleted = "deleted"
	// VoucherFormatTypeFix type
	VoucherFormatTypeFix = "fix"
	// VoucherFormatTypeRandom type
	VoucherFormatTypeRandom = "random"
)

var (
	//ErrorResourceNotFound :
	ErrorResourceNotFound = errors.New("Resource Not Found")
	// ErrorNoDataAffected :
	ErrorNoDataAffected = errors.New("No Data Affected")
)
