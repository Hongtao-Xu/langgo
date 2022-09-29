package main

import (
	"fmt"
	"langgo"
	"langgo/components/hello"
	"langgo/core/log"
)

func main() {
	langgo.Run(&hello.Instance{})
	log.Logger("component", "hello").Info().Msg(hello.GetInstance().Message)
	fmt.Printf("component name is `%s`, message is `%s`\n", hello.GetInstance().GetName(), hello.GetInstance().Message)
}
