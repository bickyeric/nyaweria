package entity

import liberr "github.com/bickyeric/nyaweria/errors"

type ResponseBody struct {
	Message string               `json:"message,omitempty"`
	Errors  []liberr.ErrorDetail `json:"errors,omitempty"`
	Data    any                  `json:"data,omitempty"`
}
