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

	Debug = log.New(debugHandle, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
	Info = log.New(infoHandle, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Warning = log.New(warnHandle, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(errorHandle, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}
