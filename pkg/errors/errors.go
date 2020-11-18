package errors

// Default errors

// Use this errors in your return values of your public methods.
// These errors can be transported via grpc

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// The entry to delete/read/update could not be found
var ErrNotFound = status.New(codes.NotFound, "Requested element not found").Err()

// ErrDuplicate is used when a record already exists.
var ErrDuplicate = status.New(codes.AlreadyExists, "Element already exists").Err()

func ErrConstraintViolation() *fieldError {
	return NewFieldError(codes.InvalidArgument, "Constraint violation")
}
