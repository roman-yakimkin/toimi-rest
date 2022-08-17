package errors

import "fmt"

type ValidationError struct {
	Err error
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error: %v", e.Err)
}
