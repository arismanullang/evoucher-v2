package util

import (
	"fmt"
	"net/http"

	"github.com/gilkor/athena/lib/x/jsonerr"
)

type (
	paginationData struct {
		Next     string `json:"next,omitempty"`
		Previous string `json:"previous,omitempty"`
	}
)

type (
	// ErrResponse : embeding type of jsonerr.ErrorResponse
	ErrResponse struct {
		*jsonerr.ErrorResponse
	}
	//Response object JSON
	Response struct {
		Pagination *paginationData `json:"pagination,omitempty"`
		Error      *ErrResponse    `json:"error,omitempty"`
		Data       interface{}     `json:"data,omitempty"`
	}
)

//NewResponse : new response
func NewResponse() *Response {
	return &Response{}
}

//SetResponse :
func (r *Response) SetResponse(data interface{}) *Response {
	r.Data = data
	return r
}

//SetPagination :
func (r *Response) SetPagination(req *http.Request, page int, next bool) {
	u := *req.URL
	u.User = nil

	nextp := ""
	if next {
		q := u.Query()
		q.Set("page", fmt.Sprintf("%d", page+1))
		u.RawQuery = q.Encode()
		nextp = u.String()
	}

	prev := ""
	if page > 1 {
		q := u.Query()
		q.Set("page", fmt.Sprintf("%d", page-1))
		u.RawQuery = q.Encode()
		prev = u.String()
	}

	if nextp != "" || prev != "" {
		r.Pagination = &paginationData{nextp, prev}
	}
}

//SetError : set error type of jsonerr.ErrorResponse
func (r *Response) SetError(e ErrResponse) {
	r.Error = &e
}

// NewError :
func NewError(je jsonerr.ErrorResponse) ErrResponse {
	return ErrResponse{&jsonerr.ErrorResponse{Status: je.Status, Code: je.Code, Message: je.Message, Args: je.Args}}
}

// SetStatus :
func (e ErrResponse) SetStatus(i int) ErrResponse {
	e.Status = i
	return e
}

//SetMessage :
func (e ErrResponse) SetMessage(m string) ErrResponse {
	e.Message = m
	return e
}

//SetArgs :
func (e ErrResponse) SetArgs(args ...string) ErrResponse {
	e.Args = args
	return e
}
