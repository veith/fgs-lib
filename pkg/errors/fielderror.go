package errors

// Default errors for task

// Use this errors in your return values of your public methods.
// These errors can be transported via grpc

import (
	"fmt"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type fieldError struct {
	status     *status.Status
	badRequest *errdetails.BadRequest
}

func NewFieldError(statusCode codes.Code, statusMessage string) *fieldError {
	fieldError := &fieldError{}
	fieldError.status = status.New(statusCode, statusMessage)
	fieldError.badRequest = &errdetails.BadRequest{FieldViolations: []*errdetails.BadRequest_FieldViolation{}}
	return fieldError
}

func (fe *fieldError) AddFieldViolation(fieldName string, fieldLevelErrorMessage string) {

	fe.badRequest.FieldViolations = append(fe.badRequest.FieldViolations, &errdetails.BadRequest_FieldViolation{
		Field:       fieldName,
		Description: fieldLevelErrorMessage,
	})
}

func (fe *fieldError) GetErr() error {
	st, err := fe.status.WithDetails(fe.badRequest)
	if err != nil {
		// If this errored, it will always error
		// here, so better panic so we can figure
		// out why than have this silently passing.
		panic(fmt.Sprintf("Unexpected error attaching metadata: %v", err))
	}
	return st.Err()
}

func (fe *fieldError) HasFieldViolation() bool {
	return len(fe.badRequest.FieldViolations) > 0
}
