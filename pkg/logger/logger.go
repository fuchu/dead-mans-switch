package logger

import (
	"os"

	log "github.com/sirupsen/logrus"
)

var Log = log.New()

func init() {
	// Log as JSON instead of the default ASCII formatter.
	Log.SetFormatter(&log.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05", // 自定义时间戳格式
	})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	Log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	Log.SetLevel(log.WarnLevel)
}
