package controller

import (
	"fmt"
	"net/http"
)

type errorData struct {
	Code   	string `json:"code"`
	Title  	string `json:"title"`
	Detail 	string `json:"detail"`
	TraceID string `json:"traceID"`
}

type paginationData struct {
	Next     string `json:"next,omitempty"`
	Previous string `json:"previous,omitempty"`
}

type Response struct {
	Pagination *paginationData `json:"pagination,omitempty"`
	Errors     *errorData      `json:"errors,omitempty"`
	Data       interface{}     `json:"data,omitempty"`
}

func NewResponse(data interface{}) *Response {
	return &Response{Data: data}
}

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

func (r *Response) AddError(code, title, detail, traceID string) {
	r.Errors = &errorData{code, title, detail, traceID}
}

//func (r *Response) AddGovalidatorErrors(errs govalidator.Errors) {
//	for _, v := range errs {
//		te := v.(govalidator.Error)
//		tn := snaker.CamelToSnake(te.Name)
//		r.AddError("000001", "Validation Error", fmt.Sprintf("%s: %s", tn, te.Err.Error()), tn)
//	}
//}
