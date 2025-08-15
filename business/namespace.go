package business

import "vortice/object"

type Namespace struct {
	ns object.Namespace
}

func NewNamespace(ns string) *Namespace {
	return &Namespace{}
}

// RegisterExtension ...
func (ns *Namespace) RegisterExtension() {

}

// RegisterAbility ...
func (ns *Namespace) RegisterAbility() {

}
