package main

import (
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"log"
	"log/syslog"
	"os"
	"time"
)

var syslogWriter *syslog.Writer
var syslogEnable = false
var syslogAddr = "127.0.0.1:514"
var syslogTag = "FakeSSDPServer"
var syslogLevel = INFO
var logLevel = DEBUG
var logTime = true
var syslogInitialized = false

const (
	FATAL = 0
	ERROR = 1
	WARN  = 2
	INFO  = 3
	DEBUG = 4
)

func initSyslog() {
	logInfo("Syslog server is %s\n", syslogAddr)
	var err error
	syslogWriter, err = syslog.Dial("udp", syslogAddr, syslog.LOG_NOTICE|syslog.LOG_USER, syslogTag)
	if err != nil {
		log.Fatal(err)
		return
	}
}

func doLog(text string, level int) {
	if level <= logLevel {
		if logTime {
			text = time.Now().Format("2006-01-02 15:04:05 -0700 ") + text
		}
		if level <= ERROR {
			fmt.Fprintln(os.Stderr, text)
		} else {
			fmt.Fprintln(os.Stdout, text)
		}
	}
	if syslogEnable {
		if !syslogInitialized {
			initSyslog()
			syslogInitialized = true
		}
		if level <= syslogLevel {
			gbk := simplifiedchinese.GBK.NewEncoder()
			gbkText, err := gbk.String(text)
			if err != nil {
				fmt.Fprintln(os.Stderr, "[syslog] Error when convert UTF8 to GBK")
			}
			switch level {
			case FATAL:
				syslogWriter.Emerg(gbkText)
			case ERROR:
				syslogWriter.Err(gbkText)
			case WARN:
				syslogWriter.Warning(gbkText)
			case INFO:
				syslogWriter.Info(gbkText)
			case DEBUG:
				syslogWriter.Debug(gbkText)
			}
		}
	}
}

func logDebug(format string, a ...interface{}) {
	text := fmt.Sprintf("[DEBG] "+format, a...)
	doLog(text, DEBUG)
}

func logInfo(format string, a ...interface{}) {
	text := fmt.Sprintf("[INFO] "+format, a...)
	doLog(text, INFO)
}

func logWarn(format string, a ...interface{}) {
	text := fmt.Sprintf("[WARN] "+format, a...)
	doLog(text, WARN)
}

func logError(format string, a ...interface{}) {
	text := fmt.Sprintf("[ERRO] "+format, a...)
	doLog(text, ERROR)
}

func logFatal(format string, a ...interface{}) {
	text := fmt.Sprintf("[FATA] "+format, a...)
	doLog(text, FATAL)
}
