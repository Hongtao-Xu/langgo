package core

var EnvName string
var WorkDir string

const (
	Development = "development"
	Production  = "production"
)

type DeferHandler func()

var deferHandles []DeferHandler

// DeferRun 运行任务
func DeferRun() {
	for _, foo := range deferHandles {
		foo()
	}
}

// DeferAdd 添加任务
func DeferAdd(handle DeferHandler) {
	deferHandles = append(deferHandles, handle)
}
