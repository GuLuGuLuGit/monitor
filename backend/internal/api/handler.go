package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/monitor/openclaw-monitor/internal/collector"
	"github.com/monitor/openclaw-monitor/internal/database"
)

type Handler struct {
	collector *collector.Collector
}

func NewHandler(c *collector.Collector) *Handler {
	return &Handler{collector: c}
}

// GetLatestStatus 获取最新状态
func (h *Handler) GetLatestStatus(c *gin.Context) {
	snapshot, err := database.GetLatestStatusSnapshot()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No status data found"})
		return
	}
	c.JSON(http.StatusOK, snapshot)
}

// GetStatusHistory 获取状态历史
func (h *Handler) GetStatusHistory(c *gin.Context) {
	hours := 24
	if h := c.Query("hours"); h != "" {
		if parsed, err := strconv.Atoi(h); err == nil {
			hours = parsed
		}
	}

	snapshots, err := database.GetStatusHistory(hours)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, snapshots)
}

// GetLatestHealth 获取最新健康检查
func (h *Handler) GetLatestHealth(c *gin.Context) {
	health, err := database.GetLatestHealthCheck()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No health data found"})
		return
	}
	c.JSON(http.StatusOK, health)
}

// GetActiveSessions 获取活动会话
func (h *Handler) GetActiveSessions(c *gin.Context) {
	sessions, err := database.GetActiveSessions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, sessions)
}

// GetAgents 获取所有代理
func (h *Handler) GetAgents(c *gin.Context) {
	agents, err := database.GetAllAgents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, agents)
}

// RefreshData 手动刷新数据
func (h *Handler) RefreshData(c *gin.Context) {
	go h.collector.CollectAll()
	c.JSON(http.StatusOK, gin.H{"message": "Data collection started"})
}

// GetSystemInfo 获取系统信息
func (h *Handler) GetSystemInfo(c *gin.Context) {
	version, err := h.collector.GetOpenClawVersion()
	if err != nil {
		version = "unknown"
	}

	info := gin.H{
		"openclaw_version": version,
		"server_time":      time.Now(),
		"database_status":  "connected",
	}
	c.JSON(http.StatusOK, info)
}

// GetDashboardData 获取仪表板数据
func (h *Handler) GetDashboardData(c *gin.Context) {
	status, _ := database.GetLatestStatusSnapshot()
	health, _ := database.GetLatestHealthCheck()
	agents, _ := database.GetAllAgents()
	sessions, _ := database.GetActiveSessions()

	dashboard := gin.H{
		"status":   status,
		"health":   health,
		"agents":   agents,
		"sessions": sessions,
		"updated":  time.Now(),
	}
	c.JSON(http.StatusOK, dashboard)
}
