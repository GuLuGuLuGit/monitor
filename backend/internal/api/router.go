package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(handler *Handler) *gin.Engine {
	router := gin.Default()

	// CORS 配置
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	router.Use(cors.New(config))

	// API 路由组
	api := router.Group("/api")
	{
		// 状态相关
		status := api.Group("/status")
		{
			status.GET("/latest", handler.GetLatestStatus)
			status.GET("/history", handler.GetStatusHistory)
			status.POST("/refresh", handler.RefreshData)
		}

		// 健康检查
		health := api.Group("/health")
		{
			health.GET("/latest", handler.GetLatestHealth)
		}

		// 会话管理
		sessions := api.Group("/sessions")
		{
			sessions.GET("/active", handler.GetActiveSessions)
		}

		// 代理管理
		agents := api.Group("/agents")
		{
			agents.GET("", handler.GetAgents)
		}

		// 系统信息
		system := api.Group("/system")
		{
			system.GET("/info", handler.GetSystemInfo)
		}

		// 仪表板
		api.GET("/dashboard", handler.GetDashboardData)
	}

	// 静态文件服务
	router.Static("/static", "./frontend")
	router.StaticFile("/", "./frontend/index.html")

	return router
}
