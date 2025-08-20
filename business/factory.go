package business

// AbilityFactoryFunc 定义 Ability 工厂函数类型
type AbilityFactoryFunc[T, O any, E Extension] interface {
	~func(O, E) T
}
