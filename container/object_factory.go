package container

type ObjectFactory interface {
	Init() error
	Destroy() error
	GetObject() (Object, error)
}
