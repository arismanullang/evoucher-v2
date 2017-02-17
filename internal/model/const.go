package model

import (
	"errors"
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

	ErrMessageInternalError        string = "internal error , failed when open the variant objects"
	ErrMessageVoucherNotActive     string = "voucher is not active yet (before start date)"
	ErrMessageVoucherDisabled      string = "voucher has been disabled (has already been used or paid)"
	ErrMessageVoucherExpired       string = "voucher has already expired (after expiration date)"
	ErrMessageVoucherAlreadyPaid   string = "voucher has already Paid"
	ErrMessageInvalidVoucher       string = "invalid voucher , VariantID not found"
	ErrMessageVoucherQtyExceeded   string = "voucher's quantities limit has been exceeded"
	ErrMessageVoucherRulesViolated string = "order did not match validation rules"

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
)
