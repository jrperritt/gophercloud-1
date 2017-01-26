package openstack

type ErrExit1 struct {
	Err error
}

func (e ErrExit1) Error() string {
	return e.Err.Error()
}

func (e ErrExit1) ExitCode() int {
	return 1
}
