package traits

type ProgressItem struct {
	endch, upch chan interface{}
}

func (pi *ProgressItem) SetEndCh(endch chan interface{}) {
	pi.endch = endch
}

func (pi *ProgressItem) EndCh() chan interface{} {
	return pi.endch
}

type ProgressItemBytes struct {
	ProgressItem
}

type ProgressItemText struct {
	ProgressItem
}

func (pip *ProgressItemText) Size() int64 {
	return 2
}

type ProgressItemPct struct {
	ProgressItem
}

func (pip *ProgressItemPct) Size() int64 {
	return 100
}
