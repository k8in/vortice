package business

type Kind string

const (
	KindExtension Kind = "extension"
	KindAbility   Kind = "ability"

	KindTag = "kind"
)

func (kind Kind) String() string {
	return string(kind)
}

type Extension interface{}
