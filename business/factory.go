package business

type AbilityFactoryFunc[T, O any, E Extension] interface {
	~func(O, E) T
}
