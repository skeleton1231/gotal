package code

import (
	"net/http"

	"github.com/skeleton1231/gotal/internal/pkg/errors"
)

type ErrCode struct {
	// C refers to the code of the ErrCode.
	C int

	// HTTP status that should be used for the associated error code.
	HTTP int

	// External (user) facing error text.
	Ext string
}

var _ errors.Coder = &ErrCode{}

// Code returns the integer code of ErrCode.
func (coder ErrCode) Code() int {
	return coder.C
}

// String implements stringer. String returns the external error message,
// if any.
func (coder ErrCode) String() string {
	return coder.Ext
}

// HTTPStatus returns the associated HTTP status code, if any. Otherwise,
// returns 200.
func (coder ErrCode) HTTPStatus() int {
	if coder.HTTP == 0 {
		return http.StatusInternalServerError
	}

	return coder.HTTP
}

func isValidHTTPStatus(code int) bool {
	validStatusCodes := []int{http.StatusOK, http.StatusBadRequest, http.StatusUnauthorized, http.StatusForbidden, http.StatusNotFound, http.StatusInternalServerError}
	for _, v := range validStatusCodes {
		if code == v {
			return true
		}
	}
	return false
}

func register(code int, httpStatus int, message string, refs ...string) {
	if !isValidHTTPStatus(httpStatus) {
		panic("http code not in `200, 400, 401, 403, 404, 500`")
	}

	coder := &ErrCode{
		C:    code,
		HTTP: httpStatus,
		Ext:  message,
	}

	errors.MustRegister(coder)
}
