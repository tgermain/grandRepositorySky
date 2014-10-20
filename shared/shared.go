package shared

import (
	"github.com/op/go-logging"
	"github.com/tgermain/grandRepositorySky/dataSet"
	"os"
)

//Global var -----------------------------------------------------------
var LocalId, LocalIp, LocalPort string
var Datas = dataSet.MakeDataSet()

//Objects parts ---------------------------------------------------------
type DistantNode struct {
	Id, Ip, Port string
}

//Log part -------------------------------------------------------------

// Example format string. Everything except the message has a custom color
// which is dependent on the log level. Many fields have a custom output
// formatting too, eg. the time returns the hour down to the milli second.
var format = "%{color}%{time:15:04:05.000000} ▶ %{level:.4s} %{id:03x} %{message} %{color:reset} "

func SetupLogger() *logging.Logger {
	// Setup one stderr and one file backend and combine them both into one
	// logging backend. By default stderr is used with the standard log flag.

	//stdErr backend
	logBackend := logging.NewLogBackend(os.Stderr, "", 0)

	//file creation and opening
	logFileBaseName := "mainLog.log"
	logFileName := "./" + logFileBaseName
	logFile, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
	if err != nil {
		panic("Could not open log file: " + err.Error())
	}
	//file backend
	logFileBackend := logging.NewLogBackend(logFile, "", 0)

	logging.SetBackend(logBackend, logFileBackend)
	logging.SetFormatter(logging.MustStringFormatter(format))

	logging.SetLevel(logging.NOTICE, "main")

	Logger = logging.MustGetLogger("main")
	return Logger
}

var Logger *logging.Logger
