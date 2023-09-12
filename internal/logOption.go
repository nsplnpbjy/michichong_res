package internal

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/rs/zerolog"
)

var (
	logFile *os.File
	logErr  error
)

func LogInit() {
	//使用log
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logFile, logErr = os.OpenFile(time.Now().GoString()+"info.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
		fmt.Println("打开日志错误")
	}
	log.SetOutput(logFile)
	log.Print("logInited")
}

func CloseLog() bool {
	if err := logFile.Close(); err != nil {
		return false
	}
	return true
}
