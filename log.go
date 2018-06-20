package beacon

import (
	"log"
)

// Log is the logging abstraction for the beacon client.
type Log interface {
	// Debug logs a debug message.
	Debug(source NRN, msg string, data ...interface{})
	// Warn logs a warning message.
	Warn(source NRN, msg string, data ...interface{})
	// Error logs an error message.
	Error(source NRN, msg string, err error, data ...interface{})
}

type ConsoleLog struct{}

func (c ConsoleLog) Debug(source NRN, msg string, data ...interface{}) {
	log.Printf("DBG %s - %q  %#v", source, msg, data)
}
func (c ConsoleLog) Warn(source NRN, msg string, data ...interface{}) {
	log.Printf("WRN %s - %q  %#v", source, msg, data)

}
func (c ConsoleLog) Error(source NRN, msg string, err error, data ...interface{}) {
	log.Printf("ERR %s - %q: %s %#v", source, msg, err, data)
}

type EmptyLog struct{}

func (c EmptyLog) Debug(source NRN, msg string, data ...interface{}) {
}
func (c EmptyLog) Warn(source NRN, msg string, data ...interface{}) {

}
func (c EmptyLog) Error(source NRN, msg string, err error, data ...interface{}) {
}
