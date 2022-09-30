package langgo

import (
	"github.com/Hongtao-Xu/langgo/core"
	"github.com/Hongtao-Xu/langgo/core/log"
	"github.com/Hongtao-Xu/langgo/heplers/io"
	"github.com/joho/godotenv"
	"os"
	"path"
)

func Init() {
	core.EnvName = os.Getenv("langgo_env")
	//1.配置处理
	if core.EnvName == "" { //默认开发环境
		core.EnvName = core.Development
	}
	if core.WorkDir == "" {
		core.WorkDir = os.Getenv("langgo_workdir_dir")
	}
	if core.WorkDir == "" {
		core.WorkDir, _ = os.Getwd() //当前工作目录
		os.Setenv("langgo_workdir_dir", core.WorkDir)
	}
	envPath := path.Join(core.WorkDir, ".env."+core.EnvName+".yml")
	confName := "app" //默认配置文件名称
	//加载配置文件
	if io.FileExists(envPath) {
		err := godotenv.Load(envPath)
		if err != nil {
			log.Logger("langgo", "run").Warn().Err(err).Msg("load env file")
		}
		confName = os.Getenv("langgo_configuration_name")
	} else {
		log.Logger("langgo", "run").Warn().Msg("env file not find")
	}
	l := log.Instance{}
	confPath := path.Join(core.WorkDir, confName+".yml")
	err := core.LoadConfigurationFile(confPath)
	if err != nil {
		log.Logger("langgo", "run").Warn().Str("path", confPath).Msg("load app config failed.")
	}
	//2.加载日志
	l.Load()
}

func Run(instances ...core.Component) {
	Init()
	core.AddComponents(instances...)
	core.LoadComponents()
	core.SignalNotify()
}
