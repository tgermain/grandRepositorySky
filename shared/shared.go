package shared

import (
	"github.com/op/go-logging"
	"os"
	"time"
)

//Const parts -----------------------------------------------------------
type MessageType int

const (
	LOOKUP MessageType = iota
	UPDATESUCCESSOR
	UPDATEPREDECESSOR
)

var messageTypes = []string{
	"lookup",
}

func (mt MessageType) String() string {
	return messageTypes[mt]
}

//Global var -----------------------------------------------------------
var LocalId, LocalIp, LocalPort string

//Objects parts ---------------------------------------------------------
type DistantNode struct {
	Id, Ip, Port string
}

type SendingQueueMsg struct {
	DaType      MessageType
	Destination *DistantNode
	Args        map[string]string
}

//example of map init map[string]int{
//     "rsc": 3711,
//     "r":   2138,
//     "gri": 1908,
//     "adg": 912,
// }

//Log part -------------------------------------------------------------

// Example format string. Everything except the message has a custom color
// which is dependent on the log level. Many fields have a custom output
// formatting too, eg. the time returns the hour down to the milli second.
var format = "%{color}%{time:15:04:05.000000} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}"

func SetupLogger() *logging.Logger {
	// Setup one stderr and one file backend and combine them both into one
	// logging backend. By default stderr is used with the standard log flag.

	//stdErr backend
	logBackend := logging.NewLogBackend(os.Stderr, "", 0)

	//file creation and opening
	logFileBaseName := "sampleLog-" + time.Now().Format(time.RFC3339) + ".log"
	logFileName := "./" + logFileBaseName
	logFile, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
	if err != nil {
		panic("Could not open log file: " + err.Error())
	}
	//file backend
	logFileBackend := logging.NewLogBackend(logFile, "", 0)

	logging.SetBackend(logBackend, logFileBackend)
	logging.SetFormatter(logging.MustStringFormatter(format))
	logging.SetLevel(logging.ERROR, "main")

	return logging.MustGetLogger("main")
}
