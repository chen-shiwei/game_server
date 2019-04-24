package logger

import (
	"config"
	"log"
	"os"
)

func init() {
	logFile := config.GetString("logFile") // just for development
	if logF, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666); err != nil {
		panic(err.Error())
	} else {
		log.SetFlags(log.Lshortfile | log.Ltime | log.LstdFlags)
		log.SetOutput(logF)
	}
}
