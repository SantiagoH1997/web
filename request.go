package web

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"strings"

	"gopkg.in/go-playground/validator.v9"
)

// validate holds the settings and caches for validating request struct values.
var validate = validator.New()

func init() {
	// Use JSON tag names for errors instead of Go struct names.
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

// Decode decodes the request into some specific value.
// If that value is a struct, it gets validated.
func Decode(r *http.Request, value interface{}) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(value); err != nil {
		return NewError(err, http.StatusBadRequest)
	}

	if err := validate.Struct(value); err != nil {

		// Use a type assertion to get the real error value.
		verrors, ok := err.(validator.ValidationErrors)
		if !ok {
			return err
		}

		var fields []ErrorField
		for _, verror := range verrors {
			field := ErrorField{
				Field: verror.Field(),
				Error: verror.Translate(nil),
			}
			fields = append(fields, field)
		}

		return &Error{
			Err:        errors.New("field validation error"),
			StatusCode: http.StatusBadRequest,
			Fields:     fields,
		}
	}

	return nil
}
