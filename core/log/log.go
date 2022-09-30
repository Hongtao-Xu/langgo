package log

import (
	"fmt"
	"github.com/rs/zerolog"
	sysio "io"

	"github.com/Hongtao-Xu/langgo/core"
	"github.com/Hongtao-Xu/langgo/heplers/io"
	"github.com/Hongtao-Xu/langgo/heplers/reopen"

	"log"
	"os"
	"path"
	"syscall"
	"time"
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

// Logger 日志处理
func Logger(name string, tag string) *zerolog.Logger {
	rp := path.Join(core.WorkDir, "logs")
	io.CreateFolder(rp, true)
	if _, ok := loggers[name]; !ok {
		p := path.Join(rp, fmt.Sprintf("%s.log", name))
		rf, err := reopen.NewFileWriter(p)
		if err != nil {
			log.Fatalf("%s create %s log file %s : %v", "langgo", name, p, err)
		}
		//开发环境
		if core.EnvName == core.Development {
			mf := sysio.MultiWriter(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.Kitchen, NoColor: false}, //输出到控制台
				zerolog.ConsoleWriter{Out: rf, TimeFormat: time.RFC3339, NoColor: true}) //输出到文件
			l := zerolog.New(mf).With().Str("tag", tag).Timestamp().Logger()
			loggers[name] = item{
				logger: l,
				writer: rf,
			}
		} else { //生成环境
			zc := zerolog.ConsoleWriter{Out: rf, TimeFormat: time.RFC3339, NoColor: true}
			l := zerolog.New(zc).With().Str("tag", tag).Timestamp().Logger()
			loggers[name] = item{
				logger: l,
				writer: rf,
			}
		}
	}
	it := loggers[name]
	return &it.logger
}
