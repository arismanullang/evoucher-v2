package model

import (
	"errors"
)

type (
	Response struct {
		Status  string
		Message string
		Data    interface{}
	}
)

var (
	ErrResourceNotFound = errors.New("resource not found.")
	ErrValidationError  = errors.New("validation error.")
	ErrInvalidPassword  = errors.New("invalid password.")
	ErrNotModified      = errors.New("data not modified.")
	ErrDuplicateEntry   = errors.New("duplicate entry.")
)

const (
	ResponseStateOk  string = "Ok"
	ResponseStateNok string = "Nok"

	ErrCodeResourceNotFound     string = "resource_not_found"
	ErrCodeInternalError        string = "internal_error"
	ErrCodeVoucherNotActive     string = "voucher_not_active"
	ErrCodeVoucherDisabled      string = "voucher_disabled"
	ErrCodeVoucherExpired       string = "voucher_expired"
	ErrCodeVoucherAlreadyPaid   string = "voucher_already_paid"
	ErrCodeInvalidVoucher       string = "invalid_voucher"
	ErrCodeVoucherRulesViolated string = "invalid_rules_violated"
	ErrCodeVoucherQtyExceeded   string = "voucher_quantity_exceeded"
	ErrCodeMissingOrderItem     string = "missing_order_items"

	ErrMessageResourceNotFound     string = "resource not found"
	ErrMessageInternalError        string = "internal error , failed when open the variant objects"
	ErrMessageVoucherNotActive     string = "voucher is not active yet (before start date)"
	ErrMessageVoucherDisabled      string = "voucher has been disabled (has already been used or paid)"
	ErrMessageVoucherExpired       string = "voucher has already expired (after expiration date)"
	ErrMessageVoucherAlreadyPaid   string = "voucher has already Paid"
	ErrMessageInvalidVoucher       string = "invalid voucher , VariantID not found"
	ErrMessageVoucherQtyExceeded   string = "voucher's quantities limit has been exceeded"
	ErrMessageVoucherRulesViolated string = "order did not match validation rules"
	ErrMessageMissingOrderItem     string = "order items was not specified"
	ErrMessageTokenNotFound        string = "Token not found"
	ErrMessageTokenExpired         string = "Token has been expired"

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

	ALPHABET     string = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	NUMERALS     string = "1234567890"
	ALPHANUMERIC string = ALPHABET + NUMERALS

	DEFAULT_CODE   string = NUMERALS
	DEFAULT_LENGTH int    = 8
)
