package business

type Ability[O any] interface {
	Extension
	Support() bool
	Enabled() bool
	Priority() int
}

// BaseAbility 默认实现
type BaseAbility[O any] struct {
	obj O
}

func NewAbility[O any](obj O) Ability[O] {
	return BaseAbility[O]{obj: obj}
}

func (BaseAbility[O]) Support() bool {
	return true
}

func (BaseAbility[O]) Enabled() bool {
	return true
}

func (BaseAbility[O]) Priority() int {
	return 0
}
