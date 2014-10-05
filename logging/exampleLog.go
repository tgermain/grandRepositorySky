package logging

import (
	"github.com/op/go-logging"
	"os"
	"time"
)

var Log = logging.MustGetLogger("example")

// Example format string. Everything except the message has a custom color
// which is dependent on the log level. Many fields have a custom output
// formatting too, eg. the time returns the hour down to the milli second.
var format = "%{color}%{time:15:04:05.000000} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}"

// Password is just an example type implementing the Redactor interface. Any
// time this is logged, the Redacted() function will be called.
type Password string

func (p Password) Redacted() interface{} {
	return logging.Redact(string(p))
}

func main() {
	// Setup one stderr and one file backend and combine them both into one
	// logging backend. By default stderr is used with the standard log flag.

	//stdErr backend
	logBackend := logging.NewLogBackend(os.Stderr, "", 0)

	//file creation and opening
	logFileBaseName := "sampleLog-" + time.Now().Format(time.RFC3339)
	logFileName := "./" + logFileBaseName
	logFile, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
	if err != nil {
		panic("Could not open log file: " + err.Error())
	}
	//file backend
	logFileBackend := logging.NewLogBackend(logFile, "", 0)

	logging.SetBackend(logBackend, logFileBackend)
	logging.SetFormatter(logging.MustStringFormatter(format))

	// For "example", set the log level to DEBUG and ERROR.
	for _, level := range []logging.Level{logging.DEBUG, logging.ERROR} {
		logging.SetLevel(level, "example")

		Log.Debug("debug %s", Password("secret"))
		Log.Info("info")
		Log.Notice("notice")
		Log.Warning("warning")
		Log.Error("err")
		Log.Critical("crit")
	}
}
