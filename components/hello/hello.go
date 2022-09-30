package hello

import "github.com/Hongtao-Xu/langgo/core"

type Instance struct {
	Message string `yaml:"message"`
}

const name = "hello"

var instance *Instance

func (i *Instance) Load() error {
	instance = i
	err := core.GetComponentConfiguration(name, i)
	if err != nil {
		return err
	}
	return nil
}

func (i *Instance) GetName() string {
	return name
}

func GetInstance() *Instance {
	return instance
}
