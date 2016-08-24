package lib

import "io"

type Outputter interface {
	OutputResult(interface{}) error
	LimitFields(interface{})
}

type ToTabler interface {
	ToTable() error
}

type ToJSONer interface {
	ToJSON() error
}

// PreJSONer is an interface that commands will satisfy if they have a `PreJSON` method.
type PreJSONer interface {
	PreJSON(interface{}) error
}

// PreTabler is an interface that commands will satisfy if they have a `PreTable` method.
type PreTabler interface {
	PreTable(interface{}) error
}

type CustomWriterer interface {
	CustomWriter() io.Writer
}
