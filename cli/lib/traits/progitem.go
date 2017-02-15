package traits

import "io"

type ProgressItem struct {
	endch, upch chan interface{}
	id          string
}

func (pi *ProgressItem) Init() {
	pi.endch = make(chan interface{})
	pi.upch = make(chan interface{})
}

func (pi *ProgressItem) SetEndCh(endch chan interface{}) {
	pi.endch = endch
}

func (pi *ProgressItem) EndCh() chan interface{} {
	return pi.endch
}

func (pi *ProgressItem) SetUpCh(upch chan interface{}) {
	pi.upch = upch
}

func (d *ProgressItem) UpCh() chan interface{} {
	return d.upch
}

func (pi *ProgressItem) SetID(id string) {
	pi.id = id
}

func (pi *ProgressItem) ID() string {
	return pi.id
}

type ProgressItemBytes struct {
	ProgressItem
	size int64
}

func (pi *ProgressItemBytes) SetSize(size int64) {
	pi.size = size
}

func (pi *ProgressItemBytes) Size() int64 {
	return pi.size
}

type ProgressItemBytesRead struct {
	ProgressItemBytes
	reader io.Reader
}

func (pi *ProgressItemBytesRead) SetReader(reader io.Reader) {
	pi.reader = reader
}

func (pi *ProgressItemBytesRead) Reader() io.Reader {
	return pi.reader
}

type ProgressItemBytesWrite struct {
	ProgressItemBytes
	writer io.Writer
}

func (pi *ProgressItemBytesWrite) SetWriter(writer io.Writer) {
	pi.writer = writer
}

func (pi *ProgressItemBytesWrite) Writer() io.Writer {
	return pi.writer
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
