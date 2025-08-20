package business

import "vortice/container"

type (
	// Target is a type that can represent any value, used for generic or flexible type handling.
	Target any
	// Ability defines an interface for capabilities with support, enablement, priority, and a default extension implementation.
	Ability[O Target, E Extension] interface {
		// Support returns true, indicating that the ability is supported.
		Support() bool
		// Enabled returns true, indicating that the ability is currently enabled.
		Enabled() bool
		// Priority returns the priority of the ability, with 100 being the default value.
		Priority() int
		// DefaultImpl returns the default implementation of the extension for this ability.
		DefaultImpl() E
	}
	// AbilityFactory is a function type that creates an instance of T given an input of type O.
	AbilityFactory[O Target, E Extension, T Ability[O, E]] interface {
		~func(O, E) T
	}
	// AbilityObject combines the Ability interface with a container.Object, allowing for managed, extensible, and prioritizable objects.
	AbilityObject[O Target, E Extension] struct {
		Ability[O, E]
		container.Object
	}
)

// BaseAbility is a generic type that serves as a base for abilities, encapsulating an object of any type.
type BaseAbility[O Target, E Extension] struct {
	obj O
	ext E
}

// NewAbility creates a new Ability instance wrapping the provided object.
func NewAbility[O Target, E Extension](obj O, ext E) Ability[O, E] {
	return BaseAbility[O, E]{obj: obj, ext: ext}
}

// Support indicates whether the ability is supported, always returning true.
func (ba BaseAbility[O, E]) Support() bool {
	return true
}

// Enabled checks if the ability is currently enabled, always returning true.
func (ba BaseAbility[O, E]) Enabled() bool {
	return true
}

// Priority returns the priority level of the ability, with a default value of 0.
func (ba BaseAbility[O, E]) Priority() int {
	return 0
}

// DefaultImpl returns the extension associated with the BaseAbility, implementing a default behavior.
func (ba BaseAbility[O, E]) DefaultImpl() E {
	return ba.ext
}
