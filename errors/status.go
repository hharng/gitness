// Copyright 2023 Harness, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package errors

import (
	"errors"
	"fmt"
)

type Status string

const (
	StatusConflict           Status = "conflict"
	StatusInternal           Status = "internal"
	StatusInvalidArgument    Status = "invalid"
	StatusNotFound           Status = "not_found"
	StatusNotImplemented     Status = "not_implemented"
	StatusUnauthorized       Status = "unauthorized"
	StatusFailed             Status = "failed"
	StatusPreconditionFailed Status = "precondition_failed"
	StatusAborted            Status = "aborted"
)

type Error struct {
	// Source error
	Err error

	// Machine-readable status code.
	Status Status

	// Human-readable error message.
	Message string

	// Details
	Details map[string]any
}

func (e *Error) Unwrap() error {
	return e.Err
}

// Error implements the error interface.
func (e *Error) Error() string {
	return e.Message
}

// AsStatus unwraps an error and returns its code.
// Non-application errors always return StatusInternal.
func AsStatus(err error) Status {
	if err == nil {
		return ""
	}
	e := AsError(err)
	if e != nil {
		return e.Status
	}
	return StatusInternal
}

// Message unwraps an error and returns its message.
func Message(err error) string {
	if err == nil {
		return ""
	}
	e := AsError(err)
	if e != nil {
		return e.Message
	}
	return err.Error()
}

// Details unwraps an error and returns its details.
func Details(err error) map[string]any {
	if err == nil {
		return nil
	}
	e := AsError(err)
	if e != nil {
		return e.Details
	}
	return nil
}

// AsError return err as Error.
func AsError(err error) (e *Error) {
	if err == nil {
		return nil
	}
	if errors.As(err, &e) {
		return
	}
	return
}

type Arg struct {
	Key   string
	Value any
}

// Format is a helper function to return an Error with a given status and formatted message.
func Format(code Status, format string, args ...interface{}) *Error {
	var (
		err     error
		details []Arg
		newArgs []any
	)

	for _, arg := range args {
		if arg == nil {
			continue
		}
		switch obj := arg.(type) {
		case error:
			err = obj
		case Arg:
			details = append(details, obj)
		case []Arg:
			details = append(details, obj...)
		default:
			newArgs = append(newArgs, obj)
		}
	}

	msg := fmt.Sprintf(format, newArgs...)
	newErr := &Error{
		Status:  code,
		Message: msg,
		Err:     err,
	}
	if len(details) > 0 {
		newErr.Details = make(map[string]any)
		for _, arg := range details {
			newErr.Details[arg.Key] = arg.Value
		}
	}
	return newErr
}

// NotFound is a helper function to return an not found Error.
func NotFound(format string, args ...interface{}) *Error {
	return Format(StatusNotFound, format, args...)
}

// InvalidArgument is a helper function to return an invalid argument Error.
func InvalidArgument(format string, args ...interface{}) *Error {
	return Format(StatusInvalidArgument, format, args...)
}

// Internal is a helper function to return an internal Error.
func Internal(format string, args ...interface{}) *Error {
	return Format(StatusInternal, format, args...)
}

// Conflict is a helper function to return an conflict Error.
func Conflict(format string, args ...interface{}) *Error {
	return Format(StatusConflict, format, args...)
}

// PreconditionFailed is a helper function to return an precondition
// failed error.
func PreconditionFailed(format string, args ...interface{}) *Error {
	return Format(StatusPreconditionFailed, format, args...)
}

// Failed is a helper function to return failed error status.
func Failed(format string, args ...interface{}) *Error {
	return Format(StatusFailed, format, args...)
}

// Aborted is a helper function to return aborted error status.
func Aborted(format string, args ...interface{}) *Error {
	return Format(StatusAborted, format, args...)
}

// IsNotFound checks if err is not found error.
func IsNotFound(err error) bool {
	return AsStatus(err) == StatusNotFound
}

// IsConflict checks if err is conflict error.
func IsConflict(err error) bool {
	return AsStatus(err) == StatusConflict
}

// IsInvalidArgument checks if err is invalid argument error.
func IsInvalidArgument(err error) bool {
	return AsStatus(err) == StatusInvalidArgument
}

// IsInternal checks if err is internal error.
func IsInternal(err error) bool {
	return AsStatus(err) == StatusInternal
}

// IsPreconditionFailed checks if err is precondition failed error.
func IsPreconditionFailed(err error) bool {
	return AsStatus(err) == StatusPreconditionFailed
}

// IsAborted checks if err is aborted error.
func IsAborted(err error) bool {
	return AsStatus(err) == StatusAborted
}
