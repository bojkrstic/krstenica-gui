package errorx

import (
	"errors"
	"fmt"
)

var (
	ErrTampleNotFound = errors.New("tample not found")
	ErrPriestNotFound = errors.New("priest not found")
)

type ValidationError error

func GetValidationError(resource, method, message string) error {
	return fmt.Errorf("%s %s failed with message %s", resource, method, message)
}
