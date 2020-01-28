package logger

import (
	"io"
	"io/ioutil"
	"log"
	"os"
)

var (
	// Debug level
	Debug *log.Logger
	// Info level
	Info *log.Logger
	// Warning level
	Warning *log.Logger
	// Error level
	Error *log.Logger
)

// Init logger level
func Init(level string) {
	var debugHandle, infoHandle, warnHandle, errorHandle io.Writer
	debugHandle = os.Stdout
	infoHandle = os.Stdout
	warnHandle = os.Stderr
	errorHandle = os.Stderr
	switch level {
	case "info":
		debugHandle = ioutil.Discard
	case "warn":
		debugHandle = ioutil.Discard
		infoHandle = ioutil.Discard
	case "error":
		debugHandle = ioutil.Discard
		infoHandle = ioutil.Discard
		warnHandle = ioutil.Discard
	}

	commonFlags := log.LstdFlags | log.Lmicroseconds
	if level == "debug" {
		commonFlags = log.LstdFlags | log.Lmicroseconds | log.Lshortfile
	}

	Debug = log.New(debugHandle, "DBG ", commonFlags)
	Info = log.New(infoHandle, "INF ", commonFlags)
	Warning = log.New(warnHandle, "WRN ", commonFlags)
	Error = log.New(errorHandle, "ERR ", commonFlags)
}
