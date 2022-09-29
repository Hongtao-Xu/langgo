package main

import (
	"fmt"
	"langgo"
	"langgo/components/mysql"
	"langgo/examples/custom/my"
)

func main() {
	langgo.Run(&my.Instance{}, &mysql.Instance{})
	fmt.Printf("component name is `%s`, message is `%s`\n", my.GetInstance().GetName(), my.GetInstance().Message)
}
