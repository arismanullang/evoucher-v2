package model

import (
	"errors"
)

var (
	ErrValidationError = errors.New("validation error.")
	ErrInvalidPassword = errors.New("invalid password.")
	ErrNotModified     = errors.New("data not modified.")

	ErrAccountNotFound  = errors.New("account not found.")
	ErrResourceNotFound = errors.New("resource not found.")
	ErrRouteNotFound    = errors.New("route not found.")
	ErrDuplicateEntry   = errors.New("duplicate entry.")
	ErrTokenExpired     = errors.New("token expired.")
	ErrTokenNotFound    = errors.New("token not found.")
	ErrServerInternal   = errors.New("server internal error.")
	ErrInvalidRole      = errors.New("invalid role.")
	ErrBadRequest       = errors.New("bad request.")

	// custom
	ErrProgramNotNull = errors.New("Voucher Already Generated")

	// Google Cloud Storage Config
	GCS_BUCKET     string
	GCS_PROJECT_ID string
	PUBLIC_URL     string

	//OCRA config
	OCRA_URL               string
	OCRA_EVOUCHER_APPS_KEY string

	//voucher config
	VOUCHER_URL string

	//Logger config
	LN_TRACE_ID int = 16
)

const (
	ResponseStateOk  string = "Ok"
	ResponseStateNok string = "Nok"

	RedemptionMethodQr    string = "qr"
	RedemptionMethodToken string = "token"

	ErrCodeAllowAccumulativeDisable string = "accumulation_is_not_allowed"
	ErrCodeInvalidRedeemMethod      string = "invalid_redeem_method"
	ErrCodeResourceNotFound         string = "resource_not_found"
	ErrCodeRouteNotFound            string = "route_not_found"
	ErrCodeInternalError            string = "internal_error"
	ErrCodeVoucherNotActive         string = "voucher_not_active"
	ErrCodeVoucherDisabled          string = "voucher_disabled"
	ErrCodeVoucherExpired           string = "voucher_expired"
	ErrCodeVoucherAlreadyPaid       string = "voucher_already_paid"
	ErrCodeInvalidVoucher           string = "invalid_voucher"
	ErrCodeVoucherRulesViolated     string = "invalid_rules_violated"
	ErrCodeVoucherQtyExceeded       string = "voucher_quantity_exceeded"
	ErrCodeMissingOrderItem         string = "missing_order_items"
	ErrCodeMissingToken             string = "missing_token"
	ErrCodeInvalidToken             string = "invalid_token"
	ErrCodeOTPFailed                string = "OTP_failed"
	ErrCodeInvalidPartnerQr         string = "invalid_partner_qr"
	ErrCodeInvalidPartner           string = "invalid_partner"
	ErrCodeInvalidProgram           string = "invalid_program"
	ErrCodeInvalidUser              string = "invalid_username_and_password"
	ErrCodeInvalidRole              string = "invalid_role"
	ErrCodeRedeemNotValidDay        string = "voucher_cannot_be_used_today"
	ErrCodeRedeemNotValidHour       string = "voucher_cannot_be_used_at_current_time"
	ErrCodeValidationError          string = "validation_Error"

	ErrMessageAllowAccumulativeDisable string = "Accumulation is not allowed"
	ErrMessageResourceNotFound         string = "Resource not found"
	ErrMessageInternalError            string = "Internal error "
	ErrMessageVoucherNotActive         string = "Voucher is not active yet (before start date)"
	ErrMessageVoucherDisabled          string = "Voucher has been disabled (has already been used or paid)"
	ErrMessageVoucherExpired           string = "Voucher has already expired (after expiration date)"
	ErrMessageVoucherAlreadyUsed       string = "Voucher has already used "
	ErrMessageVoucherAlreadyPaid       string = "Voucher has already paid"
	ErrMessageInvalidVoucher           string = "Invalid voucher, program id not found"
	ErrMessageVoucherQtyExceeded       string = "Voucher's quantities limit has been exceeded"
	ErrMessageVoucherRulesViolated     string = "Order did not match validation rules"
	ErrMessageMissingOrderItem         string = "Order items was not specified"
	ErrMessageTokenNotFound            string = "Token not found"
	ErrMessageTokenExpired             string = "Token has been expired"
	ErrMessageInvalidProgram           string = "Invalid Program ID."
	ErrMessageInvalidHolder            string = "Invalid Holder."
	ErrMessageInvalidPaerner           string = "Invalid Patner."
	ErrMessageNilProgram               string = "Account doesn't have any program."
	ErrMessageNilPartner               string = "Program doesn't have any partner."
	ErrMessageOTPFailed                string = "Doesn't match OTP"
	ErrMessageInvalidQr                string = "Invalid partner QR"
	ErrMessageInvalidRedeemMethod      string = "Invalid redemption method"
	ErrMessageInvalidUser              string = "Invalid username and password."
	ErrMessageRedeemNotValidDay        string = "Voucher cannot be used today."
	ErrMessageRedeemNotValidHour       string = "voucher cannot be used at current time."
	ErrMessageProgramHasBeenUsed       string = "Program has been used"
	ErrMessageValidationError          string = "Validation error"

	StatusCreated string = "created"
	StatusDeleted string = "deleted"

	VoucherTypeCash     string = "cash"
	VoucherTypediscount string = "discount"
	VoucherTypePromo    string = "promo"

	VoucherStateCreated string = "created"
	VoucherStateActived string = "actived"
	VoucherStateUsed    string = "used"
	VoucherStatePaid    string = "paid"
	VoucherStateDeleted string = "deleted"

	ProgramTypeBulk     string = "bulk"
	ProgramTypeOnDemand string = "on-demand"

	ALPHABET     = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	NUMERALS     = "1234567890"
	ALPHANUMERIC = ALPHABET + NUMERALS

	// defaut config Voucher format
	DEFAULT_CODE        string = "Numerals"
	DEFAULT_LENGTH      int    = 8
	DEFAULT_SEED_CODE   string = "Numerals"
	DEFAULT_SEED_LENGTH int    = 4

	//default config tx code
	DEFAULT_TXCODE   string = "Numerals"
	DEFAULT_TXLENGTH int    = 6

	// Redis token life time
	TOKENLIFE int = 1440

	//Challenge code config
	CHALLENGE_FORMAT string = "Numerals"
	CHALLENGE_LENGTH int    = 4
	TIMEOUT_DURATION int    = 120 //in Second
)
