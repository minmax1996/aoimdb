package logger

import (
	"fmt"
	"log"
	"os"
)

//LogLvl if type for constants
type LogLvl int

//Lvl is current level of debugs
var logLvl LogLvl

const (
	//DebugLvl Debug constant
	DebugLvl LogLvl = iota
	//InfoLvl Info constant
	InfoLvl
	//WarnLvl DebWarnug constant
	WarnLvl
	//ErrorLvl Error constant
	ErrorLvl
)

const (
	errTemplate   = "\033[0;31m[ERROR] %s\033[0m"
	ftlTemplate   = "\033[1;31m[FATAL] %s\033[0m"
	infoTemplate  = "\033[1;34m[INFO] %s\033[0m"
	warnTemplate  = "\033[1;33m[WARN] %s\033[0m"
	debugTemplate = "[DEBUG] %s"
)

//SetLevel sets logger lvl
func SetLevel(lvl LogLvl) {
	if lvl <= DebugLvl {
		logLvl = DebugLvl
	} else {
		logLvl = lvl
	}
}

var (
	_logErr *log.Logger
	_logFtl *log.Logger
	_logInf *log.Logger
	_logWrn *log.Logger
	_logDbg *log.Logger
)

func init() {
	_logErr = log.New(os.Stderr, "", log.Ldate|log.Lmicroseconds|log.Lshortfile)
	_logFtl = log.New(os.Stderr, "", log.Ldate|log.Lmicroseconds|log.Lshortfile)
	_logInf = log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds|log.Lshortfile)
	_logWrn = log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds|log.Lshortfile)
	_logDbg = log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds|log.Lshortfile)

	env := os.Getenv("ENV")
	if env == "dev" {
		SetLevel(DebugLvl)
	} else if env == "staging" {
		SetLevel(InfoLvl)
	} else if env == "prod" {
		SetLevel(WarnLvl)
	}
	//this log used for checking if debugs enabled or not
	Debug("DEBUG mode")
}

// Error error logger
func Error(v interface{}) {
	_ = _logErr.Output(2, fmt.Sprintf(errTemplate, v))
}

// ErrorFormat error logger
func ErrorFormat(format string, a ...interface{}) {
	_ = _logErr.Output(2, fmt.Sprintf(errTemplate, fmt.Sprintf(format, a...)))
}

// Fatal logger, calls os.Exit(1)
func Fatal(v interface{}) {
	_ = _logFtl.Output(2, fmt.Sprintf(ftlTemplate, v))
	os.Exit(1)
}

// FatalFormat logger, calls os.Exit(1)
func FatalFormat(format string, a ...interface{}) {
	_ = _logFtl.Output(2, fmt.Sprintf(ftlTemplate, fmt.Sprintf(format, a...)))
	os.Exit(1)
}

//Info info logger
func Info(v interface{}) {
	if logLvl > InfoLvl {
		return
	}
	_ = _logInf.Output(2, fmt.Sprintf(infoTemplate, v))
}

//InfoFormat info logger
func InfoFormat(format string, a ...interface{}) {
	if logLvl > InfoLvl {
		return
	}
	_ = _logInf.Output(2, fmt.Sprintf(infoTemplate, fmt.Sprintf(format, a...)))
}

//Warn warn logger
func Warn(v interface{}) {
	if logLvl > ErrorLvl {
		return
	}
	_ = _logWrn.Output(2, fmt.Sprintf(warnTemplate, v))
}

//WarnFormat warn logger
func WarnFormat(format string, a ...interface{}) {
	if logLvl > ErrorLvl {
		return
	}
	_ = _logWrn.Output(2, fmt.Sprintf(warnTemplate, fmt.Sprintf(format, a...)))
}

//Debug logger
func Debug(v interface{}) {
	if logLvl > DebugLvl {
		return
	}
	_ = _logDbg.Output(2, fmt.Sprintf(debugTemplate, v))
}

//DebugFormat logger
func DebugFormat(format string, a ...interface{}) {
	if logLvl > DebugLvl {
		return
	}
	_ = _logDbg.Output(2, fmt.Sprintf(debugTemplate, fmt.Sprintf(format, a...)))
}
