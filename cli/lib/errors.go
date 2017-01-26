package lib

import (
	"fmt"
	"strings"
)

// MultiError contains one or more errors encountered while trying to
// run a command
type MultiError []error

func (e MultiError) Error() string {
	errs := make([]string, len(e))
	for i, err := range e {
		errs[i] = err.Error()
	}
	return strings.Join(errs, "\n")
}

// ErrMissingFlagPrefix is the prefix for when a required flag is missing.
var ErrMissingFlagPrefix = "Missing flag:"

// ErrMissingFlag is used when a user doesn't provide a required flag.
type ErrMissingFlag struct {
	Msg string
}

func (e ErrMissingFlag) Error() string {
	return fmt.Sprintf("%s %s\n", ErrMissingFlagPrefix, e.Msg)
}

// ErrFlagFormatting is the prefix for when a flag's format is invalid.
var ErrFlagFormattingPrefix = "Invalid flag formatting:"

// ErrFlagFormatting is used when a flag's format is invalid.
type ErrFlagFormatting struct {
	Msg string
}

func (e ErrFlagFormatting) Error() string {
	return fmt.Sprintf("%s %s\n", ErrFlagFormattingPrefix, e.Msg)
}

// ErrArgsFlagPrefix is the prefix for when a flag's argument is invalid.
var ErrArgsFlagPrefix = "Argument error:"

// ErrArgs is used when a flag's arguments are invalid.
type ErrArgs struct {
	Msg string
}

func (e ErrArgs) Error() string {
	return fmt.Sprintf("%s %s\n", ErrArgsFlagPrefix, e.Msg)
}
