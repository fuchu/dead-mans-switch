package router

import (
	"dms/controller"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Router(alertmanager *controller.AlertManager) *gin.Engine {
	r := gin.Default()
	r.Use(gin.Recovery())
	r.GET("/health", controller.Health{}.GetHealth)
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	webhook := r.Group("/webhook")
	{
		webhook.POST("/alertmanager", func(c *gin.Context) {
			// 在路由处理函数中调用方法，并传递参数
			controller.ProcessAlert(c, alertmanager)
		})
	}
	return r
}
