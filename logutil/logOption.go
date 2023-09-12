package logutil // 更改包名为logutil，避免使用internal

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
)

var (
	logFile *os.File // 只保留一个全局变量logFile
)

func LogInit() error { // 返回error类型的值，而不是使用全局变量
	// 使用log
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	fileName := time.Now().Format("2006-01-02-info.log") // 使用time.Now().Format生成规范的文件名
	logFile, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("打开日志错误: %w", err) // 返回错误信息，而不是打印或忽略
	}
	log := zerolog.New(logFile).With().Timestamp().Logger() // 使用zerolog库创建日志对象
	log.Info().Msg("logInited")                             // 使用zerolog库记录日志信息
	return nil
}

func CloseLog() error { // 返回error类型的值，而不是布尔值
	if err := logFile.Close(); err != nil {
		return fmt.Errorf("关闭日志错误: %w", err) // 返回错误信息，而不是忽略或覆盖
	}
	return nil
}
