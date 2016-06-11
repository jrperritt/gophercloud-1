package lib

type Outputter interface {
	OutputResult(Resulter) error
	GetFormatOptions() []string
	LimitFields(Resulter)
}

type Tabler interface {
	ToTable()
}

type JSONer interface {
	ToJSON()
}

// PreJSONer is an interface that commands will satisfy if they have a `PreJSON` method.
type PreJSONer interface {
	PreJSON(Resulter) error
}

// PreTabler is an interface that commands will satisfy if they have a `PreTable` method.
type PreTabler interface {
	PreTable(Resulter) error
}
