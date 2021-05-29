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
	// Output special level used for script output
	Output *log.Logger
)

// Init logger level
func Init(level string, with ...string) {
	var debugHandle, infoHandle, warnHandle, errorHandle, outputHandle io.Writer
	debugHandle = os.Stdout
	infoHandle = os.Stdout
	warnHandle = os.Stderr
	errorHandle = os.Stderr
	outputHandle = ioutil.Discard
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

	if contains(with, "out") {
		outputHandle = os.Stdout
	}

	commonFlags := log.LstdFlags | log.Lmicroseconds
	if level == "debug" {
		commonFlags = log.LstdFlags | log.Lmicroseconds | log.Lshortfile
	}

	Debug = log.New(debugHandle, Gray("DBG "), commonFlags)
	Info = log.New(infoHandle, Green("INF "), commonFlags)
	Warning = log.New(warnHandle, Orange("WRN "), commonFlags)
	Error = log.New(errorHandle, Red("ERR "), commonFlags)
	Output = log.New(outputHandle, Purple("OUT "), commonFlags)
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
