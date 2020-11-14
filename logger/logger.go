package logger

import (
	"log"
	"os"
)

// Info info logger
var Info = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile).Println

// Warn warning logger
var Warn = log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile).Println

// Error error logger
var Error = log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile).Println
