package log

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/Hongtao-Xu/langgo/core"
	"github.com/Hongtao-Xu/langgo/heplers/io"
	helperString "github.com/Hongtao-Xu/langgo/heplers/string"

	"github.com/rs/zerolog"
)

//BenchmarkLoggerSystemLog 系统Log的性能
func BenchmarkLoggerSystemLog(b *testing.B) {
	io.CreateFolder("logs", true)
	logfile := "logs/syslog.log"
	os.Remove(logfile)
	f, _ := os.OpenFile(logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	defer f.Close()
	log.SetOutput(f)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		log.Println("INF ok tag=test")
	}
}

//BenchmarkZerologFile Zerolog的性能
func BenchmarkZerologFile(b *testing.B) {
	io.CreateFolder("logs", true)
	logfile := "logs/zerolog.log"
	os.Remove(logfile)
	f, _ := os.OpenFile(logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	defer f.Close()
	l := zerolog.New(f).With().Str("tag", "test").Timestamp().Logger()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Info().Msg("ok")
	}
}

//BenchmarkZerologConsole Zerolog的性能
func BenchmarkZerologConsole(b *testing.B) {
	io.CreateFolder("logs", true)
	logfile := "logs/zerolog-console.log"
	os.Remove(logfile)
	f, _ := os.OpenFile(logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	//序列化
	zc := zerolog.ConsoleWriter{Out: f, TimeFormat: time.RFC3339, NoColor: true}
	defer f.Close()
	l := zerolog.New(zc).With().Str("tag", "test").Timestamp().Logger()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Info().Msg("ok")
	}
}

//BenchmarkLoggerLanggo Langgo中Log的性能
func BenchmarkLoggerLanggo(b *testing.B) {
	core.EnvName = core.Production
	io.CreateFolder("logs", true)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Logger("langgo", "test").Info().Msg("ok")
	}
}

//BenchmarkLoggerLanggoMulti Langgo写成多个log文件的性能
func BenchmarkLoggerLanggoMulti(b *testing.B) {
	core.EnvName = core.Production
	io.CreateFolder("logs", true)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s, _ := helperString.RandString(2, helperString.LettersNumberNoZero)
		Logger(s, "test").Info().Msg("ok")
	}
}
