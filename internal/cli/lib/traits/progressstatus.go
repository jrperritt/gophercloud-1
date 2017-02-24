package traits

type ProgressStatus struct {
	barid string
}

func (s *ProgressStatus) SetBarID(barid string) {
	s.barid = barid
}

func (s *ProgressStatus) BarID() string {
	return s.barid
}

type ProgressStatusStart struct {
	ProgressStatus
	size int64
}

func (s *ProgressStatusStart) SetBarSize(size int64) {
	s.size = size
}

func (s *ProgressStatusStart) BarSize() int64 {
	return s.size
}

type ProgressStatusError struct {
	ProgressStatus
	e error
}

func (s *ProgressStatusError) SetErr(e error) {
	s.e = e
}

func (s *ProgressStatusError) Err() error {
	return s.e
}

type ProgressStatusUpdate struct {
	ProgressStatus
	change interface{}
}

func (s *ProgressStatusUpdate) SetChange(change interface{}) {
	s.change = change
}

func (s *ProgressStatusUpdate) Change() interface{} {
	return s.change
}

type ProgressStatusComplete struct {
	ProgressStatus
	result interface{}
}

func (s *ProgressStatusComplete) SetResult(result interface{}) {
	s.result = result
}

func (s *ProgressStatusComplete) Result() interface{} {
	return s.result
}
