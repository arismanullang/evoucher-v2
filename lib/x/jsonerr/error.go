package jsonerr

import "fmt"

type ErrorResponse struct {
	Status  int      `json:"status"`
	Code    string   `json:"code"`
	Message string   `json:"message"`
	Args    []string `json:"args,omitempty"`
}

func (e *ErrorResponse) WithArgs(args ...string) *ErrorResponse {
	args1 := make([]interface{}, len(args))
	for k, v := range args {
		args1[k] = v
	}
	return &ErrorResponse{
		Status:  e.Status,
		Code:    e.Code,
		Message: fmt.Sprintf(e.Message, args1...),
		Args:    args,
	}
}
