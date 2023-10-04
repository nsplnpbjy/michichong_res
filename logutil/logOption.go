package logutil // 包名不变

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
)

var (
	logFile   *os.File // 只保留一个全局变量logFile
	returnLog zerolog.Logger
)

func LogInit() error { // 返回error类型的值，而不是使用全局变量
	// 使用log
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	fileName := time.Now().Format("2006-01-02-info.log") // 使用time.Now().Format生成规范的文件名
	if err := os.MkdirAll("./log/", 0755); err != nil {  // 在打开文件之前，创建./log/这个文件夹，如果已经存在，就不会有影响
		return fmt.Errorf("创建文件夹错误: %w", err) // 检查并返回错误信息
	}
	logFile, err := os.OpenFile("./log/"+fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644) // 在文件名前面加上./log/这个相对路径
	if err != nil {
		return fmt.Errorf("打开日志错误: %w", err) 
	}
	returnLog = zerolog.New(logFile).With().Timestamp().Logger() // 使用zerolog库创建日志对象
	returnLog.Info().Msg("logInited")                            // 使用zerolog库记录日志信息
	return nil
}

func CloseLog() error { // 返回error类型的值，而不是布尔值
	if err := logFile.Close(); err != nil {
		return fmt.Errorf("关闭日志错误: %w", err) // 返回错误信息，而不是忽略或覆盖
	}
	return nil
}

func GetLog() *zerolog.Logger {
	return &returnLog
}
