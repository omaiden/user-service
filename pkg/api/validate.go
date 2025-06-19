package api

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/moonrhythm/validator"
)

func WrapError(e error) error {
	if err := (*validator.Error)(nil); errors.As(e, &err) {
		return &ValidateError{err}
	}
	if strings.HasPrefix(e.Error(), "error decoding string") {
		return &ValidateError{e}
	}
	{
		var jsonSyntaxError *json.SyntaxError
		if errors.As(e, &jsonSyntaxError) {
			return &ValidateError{e}
		}
	}
	return e
}

const ValidateErrorCode = "INVALID"

type ValidateError struct {
	err error
}

func (*ValidateError) OKError() {}

func (e *ValidateError) Error() string {
	return "validate error (" + e.err.Error() + ")"
}

func (e *ValidateError) Unwrap() error {
	return e.err
}

func (e *ValidateError) MarshalJSON() ([]byte, error) {
	xs := make([]string, 0)

	if err := (*validator.Error)(nil); errors.As(e.err, &err) {
		xs = append(xs, err.Strings()...)
	} else {
		xs = append(xs, e.err.Error())
	}

	return json.Marshal(struct {
		Code    string   `json:"code"`
		Message string   `json:"message"`
		Items   []string `json:"items"`
	}{ValidateErrorCode, "validate error", xs})
}
