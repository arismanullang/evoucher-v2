package model

import "errors"

const (
	//StatusCreated row satus "created"
	StatusCreated = "created"
	//StatusDeleted row satus "deleted"
	StatusDeleted = "deleted"
	//StatusSubmitted row satus "submitted"
	StatusSubmitted = "submitted"
	//StatusApproved row satus "approved"
	StatusApproved = "approved"
	//StatusPaid row satus "paid"
	StatusPaid = "paid"
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
	// HolderTypeOutlet tag
	HolderTypeOutlet = "outlet"
	// HolderTypeCustomer tag
	HolderTypeCustomer = "customer"

	CompanyEmailSender   = "email_confirmation_sender"
	CompanyEmailTemplate = "email_confirmation_template"
	CompanyFinanceEmails = "finance_emails"
	CompanyTimezone      = "timezone"

	BlastSender      = "sender"
	BlastTemplate    = "template_name"
	BlastImageHeader = "image_header"
	BlastImageFooter = "image_footer"
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
	ErrorExpiredToken  = errors.New("Token has expired")
	ErrorTokenNotFound = errors.New("Token not found")

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

	ErrorVoucherUsed        = errors.New("Voucher has been used")
	ErrorVoucherPaid        = errors.New("Voucher has been paid")
	ErrorVoucherExpired     = errors.New("Voucher has expired")
	ErrorVoucherInvalidTime = errors.New("Voucher can't be used at current time, please check the terms & conditions")

	ErrorBankNotFound = errors.New("Please complete the outlet bank details")
)
