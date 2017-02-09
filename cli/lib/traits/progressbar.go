package traits

import (
	"fmt"

	"github.com/gosuri/uiprogress"
)

type BarType uint8

const (
	BarPercentage BarType = iota
	BarBytes
	BarText
)

type ProgressBar struct {
	*uiprogress.Bar
	id string
}

func (b *ProgressBar) ID() string {
	return b.id
}

func (b *ProgressBar) Complete() {
	b.Set(b.Total)
}

func (b *ProgressBar) Error(err error) string {
	return fmt.Sprintf("[ERROR: %s] %s", err, b.id)
}

func (b *ProgressBar) Reset() {
	b.Set(0)
}

type ProgressBarPercentage struct {
	*ProgressBar
}

func (b *ProgressBarPercentage) Update(up interface{}) {
	b.Incr()
	b.Set(up.(int) - 1)
}

type ProgressBarBytes struct {
	*ProgressBar
}

func NewProgressBarBytes(id string) *ProgressBarBytes {
	b := new(ProgressBarBytes)
	b.id = id
	return b
}

func (b *ProgressBarBytes) Update(up interface{}) {
	incr := up.(*ProgressStatusUpdate)
	b.Set(b.Current() + incr.Increment)
	//lib.Log.Debugf("b: %+v\n", b)
}

type ProgressBarText struct {
	*ProgressBar
}

func (b *ProgressBarText) Update(up interface{}) {
	b.Incr()
	b.Set(0)
}

func (b *ProgressBarText) setBarToText() {
	b.LeftEnd = ' '
	b.RightEnd = ' '
	b.Head = ' '
	b.Fill = ' '
	b.Empty = ' '
}
