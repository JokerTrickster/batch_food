package utils

import (
	"context"

	"github.com/go-playground/validator/v10"
)

// validator instance
var val = validator.New()

// ValidateStruct : validate struct
func ValidateStruct(class interface{}) error {
	if err := val.Struct(class); err != nil {
		return ErrorMsg(context.TODO(), ErrBadParameter, Trace(), err.Error(), ErrFromClient)
	}
	return nil
}
