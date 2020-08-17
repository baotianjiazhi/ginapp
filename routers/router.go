package routers

import (
	"webapp/controller"
	"webapp/logger"

	"github.com/gin-gonic/gin"
)

func SetUp() *gin.Engine {
	r := gin.New()
	r.LoadHTMLGlob("templates/**/*")
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	// 注册相关路由信息

	r.GET("/", controller.Index)
	return r
}
