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
	// VoucherStateCreated state
	VoucherStateCreated = "created"
	// VoucherStateClaim state
	VoucherStateClaim = "claim"
	// VoucherStateUsed state
	VoucherStateUsed = "used"
	// VoucherStatePaid state
	VoucherStatePaid = "paid"
)

var (
	//ErrorResourceNotFound :
	ErrorResourceNotFound = errors.New("Resource Not Found")
	// ErrorNoDataAffected :
	ErrorNoDataAffected = errors.New("No Data Affected")
)
