package controller

import (
	"fmt"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/serenize/snaker"
)

type errorData struct {
	Code   string `json:"code"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
	Name   string `json:"name"`
}

type ResponseData struct {
	State   string      `json:"state"`
	Error   string      `json:"Error"`
	Message string      `json:"messange"`
	Data    interface{} `json:"data"`
}

type paginationData struct {
	Next     string `json:"next,omitempty"`
	Previous string `json:"previous,omitempty"`
}

type Response struct {
	Pagination *paginationData `json:"pagination,omitempty"`
	Errors     []errorData     `json:"errors,omitempty"`
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

func (r *Response) AddError(code, title, detail, name string) {
	r.Errors = append(r.Errors, errorData{code, title, detail, name})
}

func (r *Response) AddGovalidatorErrors(errs govalidator.Errors) {
	for _, v := range errs {
		te := v.(govalidator.Error)
		tn := snaker.CamelToSnake(te.Name)
		r.AddError("000001", "Validation Error", fmt.Sprintf("%s: %s", tn, te.Err.Error()), tn)
	}
}
