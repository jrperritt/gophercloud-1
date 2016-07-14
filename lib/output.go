package lib

type Outputter interface {
	OutputResult(interface{}) error
	LimitFields(interface{})
}

type Tabler interface {
	ToTable()
}

type JSONer interface {
	ToJSON()
}

// PreJSONer is an interface that commands will satisfy if they have a `PreJSON` method.
type PreJSONer interface {
	PreJSON(interface{}) error
}

// PreTabler is an interface that commands will satisfy if they have a `PreTable` method.
type PreTabler interface {
	PreTable(interface{}) error
}
