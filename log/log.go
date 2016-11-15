package log

import (
	"io"
	"os"

	"github.com/starkandwayne/goutils/ansi"
)

var logtarget io.Writer
var logfilter FilterFlag

//FilterFlag represents a logging channel to turn on
type FilterFlag uint16

const (
	//DEBUG is for development and programming information
	DEBUG FilterFlag = 1 << iota
	//INFO is for verbose information about non-error tasks being performed
	INFO
	//WARN is for things that could suggest that something is wrong
	WARN
	//ERROR is for when things go wrong with the operation of the program.
	//Doesn't necessarily need to be critical to the program continuing to run.
	ERROR
	//NONE represents no filters. Bitwise or'ing this with anything nullifies this
	NONE FilterFlag = 0
)

func init() {
	SetTarget(os.Stderr)
	//Default to having all the logs on. Technically going with "undefined behavior", though
	logfilter = 0xFFFF
}

//SetTarget takes an io.Writer that will be written to by future calls to log functions.
func SetTarget(target io.Writer) {
	logtarget = target
}

//SetFilter takes a bitstring flags representing which filters to turn on
//Example: log.SetFilter(log.INFO, log.ERROR)
func SetFilter(filter ...FilterFlag) {
	logfilter = 0
	for _, f := range filter {
		logfilter |= f
	}
}

//SetFilterLevel sets the minimum filter, turning on that level and all higher levels
//Takes the smallest flag in the given filter, ignoring greater ones
//ERROR > WARN > INFO > DEBUG
func SetFilterLevel(filter FilterFlag) {
	for i := FilterFlag(0x0001); i < FilterFlag(0x8000); i <<= 1 {
		if filter&i != 0 {
			//Take a full filter and erase the parts below the lowest bit of the input filter
			filter = 0xFFFF ^ (i - 1)
			break
		}
	}
	SetFilter(filter)
}

func writeHelper(filter FilterFlag, message string, args ...interface{}) {
	if shouldDisplay(filter) {
		_, err := ansi.Fprintf(logtarget, message+"\n", args...)
		if err != nil {
			panic("Unable to write to log target")
		}
	}
}

//Errorf prints an ERROR message to the log target, taking formatting arguments
func Errorf(message string, args ...interface{}) {
	writeHelper(ERROR, "ERROR: @R{"+message+"}", args...)
}

//Error prints an ERROR message to the log target
func Error(message string) {
	Errorf(message)
}

//Warnf prints a WARN message to the log target, taking formatting arguments
func Warnf(message string, args ...interface{}) {
	writeHelper(WARN, "WARN:  @Y{"+message+"}", args...)
}

//Warn prints a WARN message to the log target
func Warn(message string) {
	Warnf(message)
}

//Infof prints an INFO message to the log target, taking formatting arguments
func Infof(message string, args ...interface{}) {
	writeHelper(INFO, "INFO:  "+message, args...)
}

//Info prints an INFO message to the log target
func Info(message string) {
	Infof(message)
}

//Debugf prints a DEBUG message to the log target, taking formatting arguments
func Debugf(message string, args ...interface{}) {
	writeHelper(DEBUG, "DEBUG: @M{"+message+"}", args...)
}

//Debug prints a DEBUG message to the log target
func Debug(message string) {
	Debugf(message)
}

func shouldDisplay(tocheck FilterFlag) bool {
	return tocheck&logfilter > 0
}
