package util

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gilkor/athena/lib/x/jsonerr"
)

const (
	//APIIsDebugError :
	APIIsDebugError = true
)

type paginationData struct {
	Next     string `json:"next,omitempty"`
	Previous string `json:"previous,omitempty"`
	// Count     interface{} `json:"count,omitempty"`
	// TotalData interface{} `json:"total_data,omitempty"`
}

//Response object JSON
type Response struct {
	Pagination  *paginationData `json:"pagination,omitempty"`
	Error       *ErrResponse    `json:"error,omitempty"`
	ErrorDetail interface{}     `json:"error_detail,omitempty"`
	Data        interface{}     `json:"data,omitempty"`
}

//NewResponse : new response
func NewResponse() *Response {
	return &Response{}
}

//SetResponse :
func (r *Response) SetResponse(data interface{}) *Response {
	r.Data = data
	return r
}

//SetErrorDetail :
func (r *Response) SetErrorDetail(errorDetail interface{}) *Response {
	r.ErrorDetail = errorDetail
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

//SetTotalData : Pagination
// func (r *Response) SetTotalData(td interface{}) {
// 	r.Pagination.TotalData = td
// }

//SetError : set error type of jsonerr.ErrorResponse
func (r *Response) SetError(e ErrResponse) {
	r.SetErrorWithDetail(e, nil)
}

//SetErrorWithDetail : set error type of jsonerr.ErrorResponse
func (r *Response) SetErrorWithDetail(e ErrResponse, err error) {
	r.Error = &e
	if APIIsDebugError && err != nil {
		r.Error.SetArgs(err.Error())
		r.ErrorDetail = err
	}
}

//JSON : render response with JSON format
func (r *Response) JSON(w http.ResponseWriter, data interface{}, code ...int) {
	b, err := json.Marshal(data)
	if err != nil {
		log.Panic(err)
		return
	}

	// set HTTP response code
	c := http.StatusOK
	if len(code) > 0 {
		c = code[0]
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(c)
	w.Write(b)
}

// ErrResponse :
type ErrResponse struct {
	*jsonerr.Error
}

// NewError :
func NewError(je jsonerr.Error) ErrResponse {
	return ErrResponse{&jsonerr.Error{Status: je.Status, Code: je.Code, Message: je.Message, Args: je.Args}}
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
