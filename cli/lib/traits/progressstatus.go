package traits

type ProgressStatus struct {
	id string
}

func (p *ProgressStatus) BarID() string {
	return p.id
}

type ProgressStatusStart struct {
	ProgressStatus
	size int
}

func (s *ProgressStatusStart) BarSize() int {
	return s.size
}

type ProgressStatusError struct {
	ProgressStatus
	e error
}

func (s *ProgressStatusError) Err() error {
	return s.e
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
