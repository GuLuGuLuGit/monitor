package database

import (
	"fmt"
	"log"
	"time"

	"github.com/monitor/openclaw-monitor/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB 初始化数据库
func InitDB(dbPath string) error {
	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}

	// 自动迁移
	err = DB.AutoMigrate(
		&models.StatusSnapshot{},
		&models.HealthCheck{},
		&models.Session{},
		&models.Agent{},
		&models.Metric{},
	)
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("Database initialized successfully")
	return nil
}

// CleanupOldData 清理旧数据
func CleanupOldData(retentionDays int) error {
	cutoffTime := time.Now().AddDate(0, 0, -retentionDays)

	// 清理旧的状态快照
	if err := DB.Where("created_at < ?", cutoffTime).Delete(&models.StatusSnapshot{}).Error; err != nil {
		return err
	}

	// 清理旧的健康检查
	if err := DB.Where("created_at < ?", cutoffTime).Delete(&models.HealthCheck{}).Error; err != nil {
		return err
	}

	// 清理旧的指标
	if err := DB.Where("created_at < ?", cutoffTime).Delete(&models.Metric{}).Error; err != nil {
		return err
	}

	log.Printf("Cleaned up data older than %d days", retentionDays)
	return nil
}

// GetLatestStatusSnapshot 获取最新状态快照
func GetLatestStatusSnapshot() (*models.StatusSnapshot, error) {
	var snapshot models.StatusSnapshot
	err := DB.Order("timestamp desc").First(&snapshot).Error
	if err != nil {
		return nil, err
	}
	return &snapshot, nil
}

// GetStatusHistory 获取状态历史
func GetStatusHistory(hours int) ([]models.StatusSnapshot, error) {
	var snapshots []models.StatusSnapshot
	cutoffTime := time.Now().Add(-time.Duration(hours) * time.Hour)
	err := DB.Where("timestamp > ?", cutoffTime).Order("timestamp desc").Find(&snapshots).Error
	return snapshots, err
}

// GetLatestHealthCheck 获取最新健康检查
func GetLatestHealthCheck() (*models.HealthCheck, error) {
	var health models.HealthCheck
	err := DB.Order("timestamp desc").First(&health).Error
	if err != nil {
		return nil, err
	}
	return &health, nil
}

// GetActiveSessions 获取活动会话
func GetActiveSessions() ([]models.Session, error) {
	var sessions []models.Session
	err := DB.Where("status = ?", "active").Order("last_activity desc").Find(&sessions).Error
	return sessions, err
}

// GetAllAgents 获取所有代理
func GetAllAgents() ([]models.Agent, error) {
	var agents []models.Agent
	err := DB.Order("updated_at desc").Find(&agents).Error
	return agents, err
}

// SaveStatusSnapshot 保存状态快照
func SaveStatusSnapshot(snapshot *models.StatusSnapshot) error {
	return DB.Create(snapshot).Error
}

// SaveHealthCheck 保存健康检查
func SaveHealthCheck(health *models.HealthCheck) error {
	return DB.Create(health).Error
}

// SaveOrUpdateAgent 保存或更新代理
func SaveOrUpdateAgent(agent *models.Agent) error {
	var existing models.Agent
	err := DB.Where("agent_id = ?", agent.AgentID).First(&existing).Error
	if err == gorm.ErrRecordNotFound {
		return DB.Create(agent).Error
	}
	return DB.Model(&existing).Updates(agent).Error
}

// SaveMetric 保存指标
func SaveMetric(metric *models.Metric) error {
	return DB.Create(metric).Error
}
