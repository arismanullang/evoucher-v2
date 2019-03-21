package controller

import (
	"net/http"

	"github.com/gilkor/athena/lib/x/jsonerr"
	u "github.com/gilkor/evoucher/util"
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
)
