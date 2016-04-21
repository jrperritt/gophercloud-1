package lib

import "reflect"

type Resulter interface {
	Types() []reflect.Type
}
