package log

import (
	"fmt"
	"github.com/Hongtao-Xu/langgo/heplers/io"
	sysio "io"
	"log"
	"os"
	"path"
	"sync"
	"syscall"
	"time"

	"github.com/Hongtao-Xu/langgo/core"
	"github.com/Hongtao-Xu/langgo/heplers/reopen"

	"github.com/rs/zerolog"
)

type Instance struct {
	ReopenSignal syscall.Signal `yaml:"reopen_signal"`
}

const name = "log"

var instance *Instance

type item struct {
	logger zerolog.Logger
	writer reopen.Writer
}

var loggers = make(map[string]item)
var lock sync.Mutex //解决log的并发问题

// Load 加载日志实例
func (i *Instance) Load() {
	instance = i
	core.GetComponentConfiguration(name, i)
	//日志级别
	if core.EnvName == core.Development {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
	if i.ReopenSignal > 0 {
		core.SignalHandle(&core.SignalHandler{
			Sig: i.ReopenSignal,
			F: func() {
				for _, it := range loggers {
					it.writer.Reopen()
				}
			},
		})
	}
}

func (i *Instance) GetName() string {
	return name
}

// Logger 日志处理
func Logger(name string, tag string) *zerolog.Logger {
	rp := path.Join(core.WorkDir, "logs")

	if _, ok := loggers[name]; !ok {
		func() {
			//非sync.Map的线程安全解决：锁+双重检查
			lock.Lock()
			if _, ok := loggers[name]; ok { //第一个线程完成后，其他相同name的线程，直接返回
				fmt.Println("name", "exists")
				return
			}
			defer lock.Unlock()
			io.CreateFolder(rp, true)
			p := path.Join(rp, fmt.Sprintf("%s.log", name))
			rf, err := reopen.NewFileWriter(p)
			if err != nil {
				log.Fatalf("%s create %s log file %s : %v", "langgo", name, p, err)
			}
			//开发环境
			if core.EnvName == core.Development {
				//关闭序列化，使用默认输出:
				//writer := zerolog.ConsoleWriter{Out: rf, TimeFormat: time.RFC3339, NoColor: true}
				mf := sysio.MultiWriter(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.Kitchen, NoColor: false}, //输出到控制台
					rf) //输出到文件
				l := zerolog.New(mf).With().Str("tag", tag).Timestamp().Logger()
				loggers[name] = item{
					logger: l,
					writer: rf,
				}
			} else { //生产环境
				//zc := zerolog.ConsoleWriter{Out: rf, TimeFormat: time.RFC3339, NoColor: true}
				l := zerolog.New(rf).With().Str("tag", tag).Timestamp().Logger()
				loggers[name] = item{
					logger: l,
					writer: rf,
				}
			}
		}()
	}
	it := loggers[name]
	return &it.logger
}
