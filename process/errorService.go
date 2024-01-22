package process

import "fmt"

type ResponseError struct {
	StatusCode int
	Err        error
}

func (r *ResponseError) Error() string {
	return fmt.Sprintf("statusCode: %d and errMsg: %v", r.StatusCode, r.Err)
}
