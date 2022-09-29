package core

import (
	"os"
	"os/signal"
	"syscall"
)

//============信号处理==============

type SignalHandler struct {
	Sig syscall.Signal
	F   func()
}

var handlers []*SignalHandler

func SignalHandle(handler *SignalHandler) {
	handlers = append(handlers, handler)
}

// SignalNotify 监听所有信号
func SignalNotify() {
	//1.在signals中记录所有的信号种类
	var signals []os.Signal
	for _, handler := range handlers {
		find := false
		for _, sig := range signals {
			if sig == handler.Sig {
				find = true
				break
			}
		}
		if !find {
			signals = append(signals, handler.Sig)
		}
	}

	//2.遍历handlers，触发信号
	c := make(chan os.Signal)
	signal.Notify(c, signals...) //监听signals种类的信号
	go func() {
		for {
			s := <-c
			for _, handler := range handlers {
				if handler.Sig == s {
					handler.F() //调用函数
				}
			}
		}
	}()
}
