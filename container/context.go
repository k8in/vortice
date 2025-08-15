package container

type Context interface {
	ObjectFactory
	Start() error
	Shutdown() error
}
