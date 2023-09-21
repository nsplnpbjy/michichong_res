package controller

import "github.com/gin-gonic/gin"

func GetEngine() *gin.Engine {
	r := gin.Default()
	r.Use(CORS()) // 使用CORS中间件
	return r
}
