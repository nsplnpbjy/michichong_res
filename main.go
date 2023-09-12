package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/nsplnp/michichong/dboption"
	"github.com/nsplnp/michichong/logutil"
)

// 定义一个CORS中间件，用来设置响应头
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method               // 获取请求方法
		origin := c.Request.Header.Get("Origin") // 获取请求来源
		if origin != "" {
			// 设置响应头，允许跨域请求
			c.Header("Access-Control-Allow-Origin", "*")                                                                                                                          // 允许指定的域名访问，也可以设置为*，表示允许任意域名访问
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")                                                                                   // 允许使用的请求方法
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")                                                             // 允许使用的请求头
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type") // 允许浏览器获取的响应头
			c.Header("Access-Control-Allow-Credentials", "true")                                                                                                                  // 允许发送cookie
			c.Set("content-type", "application/json")                                                                                                                             // 设置返回格式是json
		}

		// 如果是预检请求，直接返回204，表示成功
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}

		c.Next() // 处理请求
	}
}

func main() {

	//使用log
	logutil.LogInit()
	defer func() {
		if err := logutil.CloseLog(); err != nil {
			log.Println("日志关闭失败")
		}
	}()

	//使用数据库
	dboption.DbInit()

	r := gin.Default()
	r.Use(CORS()) // 使用CORS中间件

	//测试使用
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"msg": "pong",
		})
	})

	//查看所有预约
	r.POST("/getRes", func(c *gin.Context) {
		DoGetRes(c)
	})

	//查看某个团的预约
	r.POST("/getSpeRes", func(c *gin.Context) {
		DoGetSpeRes(c)
	})

	//新增预约
	r.POST("/insertRes", func(c *gin.Context) {
		switch DoInsertRes(c) {
		case false:
			c.JSON(http.StatusConflict, gin.H{
				"msg":  "插入失败",
				"code": 400,
			})
		case true:
			c.JSON(http.StatusOK, gin.H{
				"msg":  "插入成功",
				"code": 200,
			})

		}
	})

	//删除预约
	r.POST("/deleteRes", func(c *gin.Context) {
		switch DoDeleteRes(c) {
		case false:
			c.JSON(http.StatusConflict, gin.H{
				"msg": "删除失败",
			})
		case true:
			c.JSON(http.StatusOK, gin.H{
				"msg": "删除成功",
			})

		}

	})

	//参观完毕
	r.POST("/done", func(c *gin.Context) {
		if Done(c) {
			c.JSON(http.StatusOK, gin.H{
				"msg": "成功",
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"msg": "失败",
			})
		}
	})
	r.Run(":8092") // 监听并在 0.0.0.0:8092 上启动服务
}
