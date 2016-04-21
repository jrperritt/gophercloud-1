package lib

type Outputter interface {
	Options() []string
}

type Tabler interface {
	Outputter
	ToTable()
}

type JSONer interface {
	Outputter
	ToJSON()
}

// PreJSONer is an interface that commands will satisfy if they have a `PreJSON` method.
type PreJSONer interface {
	PreJSON(*Resource) error
}

// PreTabler is an interface that commands will satisfy if they have a `PreTable` method.
type PreTabler interface {
	PreTable(*Resource) error
}

type Output struct{}

func (o Output) Options() []string {
	return []string{
		"json",
		"table",
	}
}

func (o Output) ToTable() {

}

func (o Output) ToJSON() {

}
