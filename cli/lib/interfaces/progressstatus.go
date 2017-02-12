package interfaces

type ProgressStatuser interface {
	SetBarID(barid string)
	BarID() string
}

type ProgressStartStatuser interface {
	ProgressStatuser
	BarSize() int64
	SetBarSize(size int64)
}

type ProgressUpdateStatuser interface {
	ProgressStatuser
	Change() interface{}
	SetChange(change interface{})
}

type ProgressCompleteStatuser interface {
	ProgressStatuser
	Result() interface{}
	SetResult(result interface{})
}

type ProgressErrorStatuser interface {
	ProgressStatuser
	Err() error
	SetErr(err error)
}
