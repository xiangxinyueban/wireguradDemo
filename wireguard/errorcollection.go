package kernelspace

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// ErrorCollection knows how to combine multiple error strings
type ErrorCollection []error

// Add puts given error to collection
func (ec *ErrorCollection) Add(errors ...error) {
	for _, err := range errors {
		if err != nil {
			*ec = append(*ec, err)
		}
	}
}

// String concatenates collection to single string
func (ec *ErrorCollection) String() string {
	return ec.Stringf("ErrorCollection: %s", ", ")
}

// Stringf returns a string representation of the underlying errors with the given format
func (ec *ErrorCollection) Stringf(format, errorDelimiter string) string {
	errorStrings := make([]string, 0)
	for _, err := range *ec {
		errorStrings = append(errorStrings, err.Error())
	}

	return fmt.Sprintf(format, strings.Join(errorStrings, errorDelimiter))
}

// Error converts collection to single error
func (ec *ErrorCollection) Error() error {
	if len(*ec) == 0 {
		return nil
	}
	return errors.New(ec.String())
}

// Errorf converts collection to single error by wanted format
func (ec *ErrorCollection) Errorf(format, errorDelimiter string) error {
	if len(*ec) == 0 {
		return nil
	}
	return errors.New(ec.Stringf(format, errorDelimiter))
}
