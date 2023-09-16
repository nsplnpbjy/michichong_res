package main

import (
	"github.com/nsplnp/michichong/controller"
	"github.com/nsplnp/michichong/dboption"
	"github.com/nsplnp/michichong/logutil"
)

// 定义一个CORS中间件，用来设置响应头

func main() {

	//使用log
	logutil.LogInit()
	log := logutil.GetLog()
	defer func() {
		if err := logutil.CloseLog(); err != nil {
			log.Info().Msg("日志关闭失败")
		}
	}()

	//使用数据库
	dboption.DbInit()
	r := controller.ControllerInit()
	r.Run(":8092") // 监听并在 0.0.0.0:8092 上启动服务
}
