package traits

type ProgressStatus struct {
	Name string
}

type ProgressStatusStart struct {
	ProgressStatus
	TotalSize int
}

type ProgressStatusError struct {
	ProgressStatus
	Err error
}

type ProgressStatusUpdate struct {
	ProgressStatus
	Increment int
	Msg       string
}

type ProgressStatusComplete struct {
	ProgressStatus
	Result interface{}
}
