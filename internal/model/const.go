package model

import (
	"errors"
)

var (
	ErrValidationError = errors.New("validation error.")
	ErrInvalidPassword = errors.New("invalid password.")
	ErrNotModified     = errors.New("data not modified.")

	ErrResourceNotFound = errors.New("Resource Not Found.")
	ErrDuplicateEntry   = errors.New("Duplicate Entry.")
	ErrTokenExpired     = errors.New("Token Expired.")
	ErrTokenNotFound    = errors.New("Token Not Found.")
	ErrServerInternal   = errors.New("Server Internal Error.")
)

const (
	ResponseStateOk  string = "Ok"
	ResponseStateNok string = "Nok"

	RedeemtionMethodQr    string = "qr"
	RedeemtionMethodToken string = "token"

	ErrCodeAllowAccumulativeDisable string = "accumulation_is_not_allowed"
	ErrCodeInvalidRedeemMethod      string = "invalid_redeem_method"
	ErrCodeResourceNotFound         string = "resource_not_found"
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
	ErrCodeOTPFailed                string = "OTP_Failed"
	ErrCodeInvalidPartnerQr         string = "invalid_partner_qr"
	ErrCodeInvalidVariant           string = "invalid_variant"
	ErrCodeInvalidUser              string = "invalid_username_and_password"
	ErrCodeRedeemNotValidDay        string = "voucher_cannot_be_used_today"
	ErrCodeRedeemNotValidHour       string = "voucher_cannot_be_used_at_current_time"

	ErrMessageAllowAccumulativeDisable string = "accumulation is not allowed"
	ErrMessageResourceNotFound         string = "resource not found"
	ErrMessageInternalError            string = "internal error "
	ErrMessageVoucherNotActive         string = "voucher is not active yet (before start date)"
	ErrMessageVoucherDisabled          string = "voucher has been disabled (has already been used or paid)"
	ErrMessageVoucherExpired           string = "voucher has already expired (after expiration date)"
	ErrMessageVoucherAlreadyUsed       string = "voucher has already Used "
	ErrMessageVoucherAlreadyPaid       string = "voucher has already Paid"
	ErrMessageInvalidVoucher           string = "invalid voucher , VariantID not found"
	ErrMessageVoucherQtyExceeded       string = "voucher's quantities limit has been exceeded"
	ErrMessageVoucherRulesViolated     string = "order did not match validation rules"
	ErrMessageMissingOrderItem         string = "order items was not specified"
	ErrMessageTokenNotFound            string = "Token not found"
	ErrMessageTokenExpired             string = "Token has been expired"
	ErrMessageInvalidMerchant          string = "vouchers can not be used in this Partner."
	ErrMessageInvalidVariant           string = "Invalid Variant ID."
	ErrMessageInvalidHolder            string = "Invalid Holder."
	ErrMessageNilVariant               string = "Account doesn't have any Variant."
	ErrMessageNilPartner               string = " doesn't have any Partner."
	ErrMessageOTPFailed                string = "Doesn't match OTP"
	ErrMessageInvalidQr                string = "Invalid parner QR"
	ErrMessageInvalidRedeemMethod      string = "Invalid Redeemtion Method"
	ErrMessageInvalidUser              string = "Invalid Username and Password."
	ErrMessageRedeemNotValidDay        string = "Voucher cannot be used today."
	ErrMessageRedeemNotValidHour       string = "voucher cannot be used at current time."
	ErrMessageVariantHasBeenUsed       string = "Variant Has been Used"

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

	VariantTypeBulk     string = "bulk"
	VariantTypeOnDemand string = "on-demand"

	ALPHABET     = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	NUMERALS     = "1234567890"
	ALPHANUMERIC = ALPHABET + NUMERALS

	// defaut config Voucher format
	DEFAULT_CODE   string = "Numerals"
	DEFAULT_LENGTH int    = 8

	//default config tx code
	DEFAULT_TXCODE   string = "Numerals"
	DEFAULT_TXLENGTH int    = 10

	// Redis token life time
	TOKENLIFE int = 1440

	// Google Cloud Storage Config
	GCS_BUCKET     string = "e-voucher"
	GCS_PROJECT_ID string = "shared-project-159515"
	PublicURL      string = "https://storage.googleapis.com/%s/%s"

	//Voucber link
	VOUCHER_URL string = "https://voucher.elys.id/v1/redeem"
)
