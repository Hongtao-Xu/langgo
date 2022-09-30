package main

import (
	"github.com/Hongtao-Xu/langgo"
	"github.com/Hongtao-Xu/langgo/components/hello"
	"github.com/Hongtao-Xu/langgo/core/log"
	"time"
)

func main() {
	langgo.Run(&hello.Instance{Message: "hello component"})
	log.Logger("component", "hello").Info().Msg(hello.GetInstance().Message)
	loop()
}
func loop() {
	i := 0
	for {
		log.Logger("app", "sleep").Info().Int("index", i).Send()
		i++
		time.Sleep(500 * time.Millisecond)
	}
}
