package main

import "github.com/gophercloud/cli/lib"

// Context satisfies the Provider interface
type Context struct {
	lib.Context
}

// Name satisfies the Provider.Name method
func (c Context) Name() string {
	return "stack"
}

// NewGlobalOptions satisfies the Provider.NewGlobalOptions method
func (c Context) NewGlobalOptions() GlobalOptionser {
	return GlobalOptions{}
}

// NewResult satisfies the Provider.NewResult method
func (c Context) NewResult() Resulter {
	return Result{}
}
