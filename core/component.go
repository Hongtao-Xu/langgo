package core

//===============组件处理===============

type Component interface {
	Load() error
	GetName() string
}

var components = make(map[string]Component)

//LoadComponents 加载组件
func LoadComponents() {

	for _, c := range components {
		c.Load()
	}
}

//AddComponents 添加组件
func AddComponents(instances ...Component) {
	for _, c := range instances {
		components[c.GetName()] = c //添加到map
	}
}
