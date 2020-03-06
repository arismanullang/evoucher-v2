package model

import "errors"

const (
	//StatusCreated row satus "created"
	StatusCreated = "created"
	//StatusDeleted row satus "deleted"
	StatusDeleted = "deleted"
	//StatusSubmitted row satus "submitted"
	StatusSubmitted = "submitted"
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
	// HolderTypePartner tag
	HolderTypePartner = "partner"
	// HolderTypeCustomer tag
	HolderTypeCustomer = "customer"
)

var (
	//ErrorResourceNotFound :
	ErrorResourceNotFound = errors.New("Resource Not Found")
	// ErrorNoDataAffected :
	ErrorNoDataAffected = errors.New("No Data Affected")
	// ErrorInternalServer :
	ErrorInternalServer = errors.New("Internal Server Error")
	// ErrorForbidden :
	ErrorForbidden = errors.New("Forbidden")
)
