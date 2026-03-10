package models

import (
	"time"
)

// StatusSnapshot 状态快照
type StatusSnapshot struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	Timestamp       time.Time `gorm:"index" json:"timestamp"`
	StatusJSON      string    `gorm:"type:text" json:"status_json"`
	ChannelsHealthy int       `json:"channels_healthy"`
	ChannelsTotal   int       `json:"channels_total"`
	CreatedAt       time.Time `json:"created_at"`
}

// HealthCheck 健康检查
type HealthCheck struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Timestamp     time.Time `gorm:"index" json:"timestamp"`
	OverallStatus string    `json:"overall_status"` // healthy/warning/error
	DetailsJSON   string    `gorm:"type:text" json:"details_json"`
	CreatedAt     time.Time `json:"created_at"`
}

// Session 会话记录
type Session struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	SessionID    string    `gorm:"index" json:"session_id"`
	Recipient    string    `json:"recipient"`
	Channel      string    `json:"channel"`
	Status       string    `json:"status"`
	StartedAt    time.Time `json:"started_at"`
	LastActivity time.Time `json:"last_activity"`
	CreatedAt    time.Time `json:"created_at"`
}

// Agent 代理信息
type Agent struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	AgentID      string    `gorm:"uniqueIndex" json:"agent_id"`
	AgentName    string    `json:"agent_name"`
	Status       string    `json:"status"`
	BindingsJSON string    `gorm:"type:text" json:"bindings_json"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Metric 指标数据
type Metric struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Timestamp   time.Time `gorm:"index" json:"timestamp"`
	MetricName  string    `gorm:"index" json:"metric_name"`
	MetricValue float64   `json:"metric_value"`
	MetricType  string    `json:"metric_type"`
	CreatedAt   time.Time `json:"created_at"`
}

// TableName 指定表名
func (StatusSnapshot) TableName() string {
	return "status_snapshots"
}

func (HealthCheck) TableName() string {
	return "health_checks"
}

func (Session) TableName() string {
	return "sessions"
}

func (Agent) TableName() string {
	return "agents"
}

func (Metric) TableName() string {
	return "metrics"
}
