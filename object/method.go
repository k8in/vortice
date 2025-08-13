package object

import (
	"reflect"
)

type (
	// Initializable 可初始化对象
	Initializable interface {
		// Init 初始化
		Init()
	}

	// Destroyable 销毁钩子
	Destroyable interface {
		// Destroy 销毁初始化
		Destroy()
	}

	// Lifecycle 生命周期管理
	Lifecycle interface {
		// Start 启动组件
		Start()
		// Stop 关闭组件
		Stop()
		// Running 组件是否在运行
		Running() bool
	}
)

var (
	initMethodType    = reflect.TypeOf((*Initializable)(nil)).Elem()
	destroyMethodType = reflect.TypeOf((*Destroyable)(nil)).Elem()
	lifecycleType     = reflect.TypeOf((*Lifecycle)(nil)).Elem()

	initMethodName    = "Init"
	destroyMethodName = "Destroy"
	startMethodName   = "Start"
	stopMethodName    = "Stop"
	runningMethod     = "Running"
)

type Methods struct {
	initMethod    *reflect.Method
	destroyMethod *reflect.Method
	startMethod   *reflect.Method
	stopMethod    *reflect.Method
	runningMethod *reflect.Method
}

func (m *Methods) init(rv reflect.Value) {
	if m.initMethod == nil {
		return
	}
	rv.MethodByName(m.initMethod.Name).Call([]reflect.Value{})
}
