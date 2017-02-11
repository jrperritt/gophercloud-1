package interfaces

type ProgressStatuser interface {
	BarID() string
}

type ProgressStartStatuser interface {
	ProgressStatuser
	BarSize() int
}

type ProgressUpdateStatuser interface {
	ProgressStatuser
}

type ProgressCompleteStatuser interface {
	ProgressStatuser
}

type ProgressErrorStatuser interface {
	ProgressStatuser
	Err() error
}
