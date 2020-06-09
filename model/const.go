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
	// ErrorInvalidToken :
	ErrorInvalidToken = errors.New("Invalid Token")
	// ErrorExpiredToken :
	ErrorExpiredToken = errors.New("Token has expired")

	// ErrorMaxAssignByDay :
	ErrorMaxAssignByDay = errors.New("You have reach the maximum limit of voucher today, try again tomorrow")
	// ErrorMaxAssignByProgram :
	ErrorMaxAssignByProgram = errors.New("You have reach the maximum limit of voucher in this program")
	// ErrorStockEmpty :
	ErrorStockEmpty = errors.New("voucher stock is empty")
	//ErrorInvalidDate :
	ErrorInvalidDate = errors.New("voucher can't be used at current date")
	//ErrorInvalidDay :
	ErrorInvalidDay = errors.New("voucher can't be used at current day")
	//ErrorInvalidTime :
	ErrorInvalidTime = errors.New("voucher can't be used at current time")
	//ErrorInvalidOutlet :
	ErrorInvalidOutlet = errors.New("voucher can't be used at current outlet")
)
