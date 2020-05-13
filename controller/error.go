package controller

import (
	"net/http"

	"github.com/gilkor/athena/lib/x/jsonerr"
	u "github.com/gilkor/evoucher-v2/util"
)

//ErrorResponse embedding of type jsonerr.ErrorResponse
var (
	//JSONErrFatal :
	JSONErrFatal = u.NewError(*jsonerr.ErrFatal)
	//JsonErrUnauthorized :
	JSONErrUnauthorized = u.NewError(*jsonerr.ErrUnauthorized)
	//JSONErrForbidden :
	JSONErrForbidden = u.NewError(*jsonerr.ErrForbidden)
	//JSONErrResourceNotFound :
	JSONErrResourceNotFound = u.NewError(
		jsonerr.Error{
			Status:  http.StatusNotFound,
			Code:    "ERR_RESOURCE_NOT_FOUND",
			Message: "Can not find requested resource.",
		})
	//JSONErrBadRequest :
	JSONErrBadRequest = u.NewError(
		jsonerr.Error{
			Status:  http.StatusBadRequest,
			Code:    "ERR_BAD_REQUEST",
			Message: "Can not find requested resource.",
		})
	//JSONErrInvalidRule :
	JSONErrInvalidRule = u.NewError(
		jsonerr.Error{
			Status:  http.StatusForbidden,
			Code:    "ERR_INVALID_RULE",
			Message: "Rule Checking Return Invalid.",
		})
	//JSONErrExceedAmount :
	JSONErrExceedAmount = u.NewError(
		jsonerr.Error{
			Status:  http.StatusOK,
			Code:    "AMOUNT_EXCEED_MAXIMUM",
			Message: "Amount Exceed Maximum.",
		})
)
