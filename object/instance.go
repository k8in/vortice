package object

import (
	"fmt"
	"reflect"
)

type Instance interface {
	ID() string
	Definition() *Definition
	Value() any
	RefValue() reflect.Value
	Initialized() bool
	Initializable
	Destroyable
	fmt.Stringer
}

type Service interface {
	Instance
	Lifecycle
}
