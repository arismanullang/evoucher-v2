package model

import (
	"errors"
	"reflect"
	"time"
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

	//voucher config
	VOUCHER_URL string

	//Logger config
	LN_TRACE_ID int = 16

	//Token lifetime
	TOKENLIFE int
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
	ErrCodeInvalidProgramType       string = "invalid_program_type"
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
	ErrCodeJsonError                string = "json_error"

	ErrMessageAllowAccumulativeDisable string = "Accumulation is not allowed"
	ErrMessageResourceNotFound         string = "Resource not found"
	ErrMessageInternalError            string = "Internal error "
	ErrMessageVoucherNotActive         string = "Voucher is not active yet (before start date)"
	ErrMessageVoucherDisabled          string = "Voucher has been disabled (has already been used or paid)"
	ErrMessageVoucherExpired           string = "Voucher has already expired (after expiration date)"
	ErrMessageVoucherAlreadyUsed       string = "Voucher has already used "
	ErrMessageVoucherAlreadyPaid       string = "Voucher has already paid"
	ErrMessageInvalidVoucher           string = "Invalid voucher, voucher id not found"
	ErrMessageVoucherQtyExceeded       string = "Voucher's quantities limit has been exceeded"
	ErrMessageVoucherRulesViolated     string = "Order did not match validation rules"
	ErrMessageInvalidProgramType       string = "Invalid program type"
	ErrMessageMissingOrderItem         string = "Order items was not specified"
	ErrMessageTokenNotFound            string = "Token not found"
	ErrMessageTokenExpired             string = "Token has been expired"
	//ErrMessageInvalidProgram           string = "Invalid Program ID."
	ErrMessageInvalidProgram      string = "Program not found."
	ErrMessageInvalidHolder       string = "Invalid Holder."
	ErrMessageInvalidPaerner      string = "Invalid Patner."
	ErrMessageNilProgram          string = "Account doesn't have any program."
	ErrMessageNilPartner          string = "Program doesn't have any partner."
	ErrMessageOTPFailed           string = "Doesn't match OTP"
	ErrMessageInvalidQr           string = "Invalid partner QR"
	ErrMessageInvalidRedeemMethod string = "Invalid redemption method"
	ErrMessageInvalidUser         string = "Invalid username and password."
	ErrMessageRedeemNotValidDay   string = "Voucher cannot be used today."
	ErrMessageRedeemNotValidHour  string = "voucher cannot be used at current time."
	ErrMessageProgramHasBeenUsed  string = "Voucher has been redeemed."
	ErrMessageValidationError     string = "Validation error."
	ErrMessageParsingError        string = "Parsing error."

	EmailCreated string = "created"
	EmailSend    string = "send"
	EmailVoid    string = "void"

	StatusCreated string = "created"
	StatusDeleted string = "deleted"

	VoucherTypeCash     string = "cash"
	VoucherTypediscount string = "discount"
	VoucherTypePromo    string = "promo"

	ProgramTypePrivilege string = "privilege"
	ProgramTypeOnDemand  string = "on-demand"
	ProgramTypeGift      string = "gift"
	ProgramTypeBulk      string = "bulk"
	ProgramTypeStock     string = "stock"

	VoucherStatePrivilege string = "privilege"
	VoucherStateCreated   string = "created"
	VoucherStateActived   string = "actived"
	VoucherStateUsed      string = "used"
	VoucherStatePaid      string = "paid"
	VoucherStateDeleted   string = "deleted"
	VoucherStateRollback  string = "rollback"
	VoucherStateSend      string = "send"

	ALPHABET     = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	NUMERALS     = "1234567890"
	ALPHANUMERIC = ALPHABET + NUMERALS

	// defaut config Voucher format
	DEFAULT_CODE               string = "Numerals"
	DEFAULT_LENGTH             int    = 8
	DEFAULT_SEED_CODE          string = "Numerals"
	DEFAULT_SEED_LENGTH        int    = 4
	DEFAULT_TRANSACTION_SEED   string = "Numerals"
	DEFAULT_TRANSACTION_LENGTH int    = 3

	//default config tx code
	DEFAULT_TXCODE   string = "Numerals"
	DEFAULT_TXLENGTH int    = 5

	//Challenge code config
	CHALLENGE_FORMAT string = "Numerals"
	CHALLENGE_LENGTH int    = 4
	TIMEOUT_DURATION int    = 120 //in Second

	//Change Log
	ColumnChangeLogInsert string = "all"
	ColumnChangeLogSelect string = "custom"
	ColumnChangeLogDelete string = "all"
	ValueChangeLogNone    string = "none"
	ValueChangeLogAll     string = "all"
	ActionChangeLogInsert string = "insert"
	ActionChangeLogUpdate string = "update"
	ActionChangeLogDelete string = "delete"
	ActionChangeLogSelect string = "select"
	ActionChangeLogLogin  string = "login"
)

func getUpdate(paramUpdate, param2 reflect.Value) map[string]reflect.Value {
	updates := make(map[string]reflect.Value)

	dataParam2 := param2.Type()
	for i := 0; i < paramUpdate.NumField(); i++ {
		f := param2.Field(i)
		string1 := f.Interface()
		string2 := paramUpdate.FieldByName(dataParam2.Field(i).Name)

		if string1 != string2.Interface() && string2.Interface() != reflect.Zero(f.Type()).Interface() {
			col, _ := dataParam2.Field(i).Tag.Lookup("db")
			va := dataParam2.Field(i).Name

			updates[va+";"+col] = string2
		}
	}

	return updates
}

func StringToTimeJakarta(param string) time.Time {
	jakarta, _ := time.LoadLocation("Asia/Jakarta")
	layout := "2006-01-02T15:04:05.000Z"
	newTime, err := time.Parse(layout, param)
	if err != nil {
		return time.Date(1001, 1, 1, 1, 1, 1, 1, time.Local)
	}
	return newTime.In(jakarta)
}

func TimeToTimeJakarta(param time.Time) time.Time {
	jakarta, _ := time.LoadLocation("Asia/Jakarta")
	return param.In(jakarta)
}
