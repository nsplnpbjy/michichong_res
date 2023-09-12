package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/nsplnp/michichong/dboption"
	"github.com/nsplnp/michichong/logutil"
)

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
				"msg": "插入失败",
			})
		case true:
			c.JSON(http.StatusOK, gin.H{
				"msg": "插入成功",
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
	r.Run(":8092") // 监听并在 0.0.0.0:8080 上启动服务
}
