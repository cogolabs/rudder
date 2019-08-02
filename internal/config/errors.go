package config

import "fmt"

// ErrMissingField is thrown when a required field is missing
type ErrMissingField struct {
	field string
}

// Error returns a string representation of the error
func (err *ErrMissingField) Error() string {
	return fmt.Sprintf("required field missing: %s", err.field)
}
