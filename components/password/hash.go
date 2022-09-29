package password

import (
	"crypto/sha1"
	"encoding/hex"
	"langgo/core"
)

type Instance struct {
	Salt string `yaml:"salt"`
}

var instance *Instance

const name = "password"

func (i *Instance) Load() error {
	instance = i
	core.GetComponentConfiguration(name, i)
	return nil
}

func (i *Instance) GetName() string {
	return name
}

//Hash 对原始字符串加盐hash加密
func Hash(orig string) string {
	hn := sha1.New()
	//sha1加密
	hn.Write([]byte(orig))
	//加盐
	hn.Write([]byte(instance.Salt))
	data := hn.Sum([]byte(""))
	return hex.EncodeToString(data)
}
