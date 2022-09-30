package main

import (
	"fmt"
	"github.com/Hongtao-Xu/langgo"
	"github.com/Hongtao-Xu/langgo/components/mysql"
	"github.com/Hongtao-Xu/langgo/examples/custom/my"
)

func main() {
	langgo.Run(&my.Instance{}, &mysql.Instance{})
	fmt.Printf("component name is `%s`, message is `%s`\n", my.GetInstance().GetName(), my.GetInstance().Message)
}
