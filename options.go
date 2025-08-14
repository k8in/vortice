package vortice

import "vortice/object"

// Option is a function type for configuring properties of an object,
// allowing modification of its attributes.
type Option func(prop *object.Property)

// WithDesc sets the description of a property, providing a brief explanation
// or additional context.
func WithDesc(desc string) Option {
	return func(prop *object.Property) {
		prop.Desc = desc
	}
}

// WithSingleton sets the scope of a property to Singleton,
// ensuring it is instantiated once and shared.
func WithSingleton() Option {
	return func(prop *object.Property) {
		prop.Scope = object.Singleton
	}
}

// WithPrototype sets the scope of a property to Prototype, indicating it should be
// instantiated each time it is requested.
func WithPrototype() Option {
	return func(prop *object.Property) {
		prop.Scope = object.Prototype
	}
}

// WithLazyInit returns an Option that sets the LazyInit flag of a Property to true,
// enabling lazy initialization.
func WithLazyInit() Option {
	return func(prop *object.Property) {
		prop.LazyInit = true
	}
}

// WithAutoStartup sets the AutoStartup property to true, enabling automatic startup for the object.
func WithAutoStartup() Option {
	return func(prop *object.Property) {
		prop.AutoStartup = true
	}
}

func WithExtension() Option {
	return func(prop *object.Property) {
		prop.AutoStartup = true
	}
}
