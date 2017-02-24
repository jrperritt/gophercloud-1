package openstack

// ErrExit1 represents an error that causes the program to exit
type ErrExit1 struct {
	Err error
}

// Error returns the error's message as a string.
func (e ErrExit1) Error() string {
	return e.Err.Error()
}

// ExitCode returns the error's exit code
func (e ErrExit1) ExitCode() int {
	return 1
}
