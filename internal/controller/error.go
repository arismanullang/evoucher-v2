package controller

import (
	"net/http"

	"github.com/gilkor/evoucher/internal/util"

	"github.com/gilkor/athena/lib/x/jsonerr"
)

//ErrorResponse embedding of type jsonerr.ErrorResponse
var (
	//ErrFatal :
	ErrFatal = util.NewError(*jsonerr.ErrFatal)
	//ErrUnauthorized :
	ErrUnauthorized = util.NewError(*jsonerr.ErrUnauthorized)
	//ErrForbidden :
	ErrForbidden = util.NewError(*jsonerr.ErrForbidden)
	//ErrResourceNotFound :
	ErrResourceNotFound = util.NewError(
		jsonerr.ErrorResponse{
			Status:  http.StatusNotFound,
			Code:    "ERR_RESOURCE_NOT_FOUND",
			Message: "Can not find requested resource.",
		})
)
