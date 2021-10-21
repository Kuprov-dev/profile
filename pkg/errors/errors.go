package errors

import "fmt"

type RequestError struct {
	StatusCode int
	Errors     uint16
	Err        error
}

func (r *RequestError) Error() string {
	return fmt.Sprintf("status %d: err %v", r.StatusCode, r.Err)
}

func NewRequestError(statusCode int, errors uint16, err error) *RequestError {
	return &RequestError{
		StatusCode: statusCode,
		Errors:     errors,
		Err:        err,
	}
}

const (
	CredsMarshalingError uint16 = 1 << iota
	ClientRequestError
	BadRequestError
	UnauthorisedError
	ForbiddenError
	AuthServiceBadGatewayError
	AuthServiceUnavailableError
)
