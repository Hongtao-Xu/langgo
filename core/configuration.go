package core

import (
	"errors"
	"gopkg.in/yaml.v3"
	"os"
)

//========配置文件解析============

var componentConfiguration = make(map[string]interface{})

//LoadConfigurationFile 加载配置
func LoadConfigurationFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(data, &componentConfiguration)
	if err != nil {
		return err
	}
	return nil
}

//GetComponentConfiguration 获取组件配置
func GetComponentConfiguration(name string, conf interface{}) error {
	if obj, ok := componentConfiguration[name]; ok {
		marshal, err := yaml.Marshal(obj)
		if err != nil {
			return err
		}
		err = yaml.Unmarshal(marshal, conf)
		if err != nil {
			return err
		}
		return nil
	} else {
		return errors.New("Component configuration not find")
	}
}
