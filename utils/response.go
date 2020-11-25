package utils

import (
	"fmt"
	"net/http"
)

// RespError interface
type RespError interface {
	Message() string
	Status() int
	Error() string
	Causes() []interface{}
}

type respError struct {
	ErrMessage string        `json:"message"`
	ErrStatus  int           `json:"status"`
	ErrError   string        `json:"error"`
	ErrCauses  []interface{} `json:"causes"`
}

func (e respError) Error() string {
	return fmt.Sprintf("message: %s - status: %d - error: %s - causes: %v",
		e.ErrMessage, e.ErrStatus, e.ErrError, e.ErrCauses)
}

func (e respError) Message() string {
	return e.ErrMessage
}

func (e respError) Status() int {
	return e.ErrStatus
}

func (e respError) Causes() []interface{} {
	return e.ErrCauses
}

// NewNotFoundError return response not found
func NewNotFoundError(msg string) RespError {
	return respError{
		ErrMessage: msg,
		ErrStatus:  http.StatusNotFound,
		ErrError:   "not found",
	}
}

// NewBadRequestError return response not found
func NewBadRequestError(msg string) RespError {
	return respError{
		ErrMessage: msg,
		ErrStatus:  http.StatusBadRequest,
		ErrError:   "bad request",
	}
}

// NewInternalServerError function
func NewInternalServerError(msg string, err error) RespError {
	result := respError{
		ErrMessage: msg,
		ErrStatus:  http.StatusInternalServerError,
		ErrError:   "internal_server_error",
	}
	if err != nil {
		result.ErrCauses = append(result.ErrCauses, err.Error())
	}
	return result
}
