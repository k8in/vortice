package container

type Context interface {
	ObjectFactory
	Start() error
	Shutdown() error
}

//type Context struct {
//	context.Context
//	mux        *sync.RWMutex
//	singletons map[string]Object
//}
