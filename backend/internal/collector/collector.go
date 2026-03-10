package collector

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/monitor/openclaw-monitor/internal/database"
	"github.com/monitor/openclaw-monitor/internal/models"
)

type Config struct {
	OpenClawPath string
	NodePath     string
	Timeout      int
}

type Collector struct {
	config Config
}

func NewCollector(config Config) *Collector {
	return &Collector{config: config}
}

// executeCommand 执行 OpenClaw 命令
func (c *Collector) executeCommand(args ...string) (string, error) {
	// 构建完整命令
	cmdArgs := []string{c.config.OpenClawPath + "/openclaw.mjs"}
	cmdArgs = append(cmdArgs, args...)

	cmd := exec.Command(c.config.NodePath, cmdArgs...)
	cmd.Dir = c.config.OpenClawPath

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("command failed: %w, output: %s", err, string(output))
	}

	return string(output), nil
}

// CollectStatus 收集状态信息
func (c *Collector) CollectStatus() error {
	log.Println("Collecting status...")

	output, err := c.executeCommand("status", "--json")
	if err != nil {
		log.Printf("Failed to collect status: %v", err)
		return err
	}

	// 解析 JSON
	var statusData map[string]interface{}
	if err := json.Unmarshal([]byte(output), &statusData); err != nil {
		log.Printf("Failed to parse status JSON: %v", err)
		return err
	}

	// 提取频道信息
	channelsHealthy := 0
	channelsTotal := 0
	if channels, ok := statusData["channels"].([]interface{}); ok {
		channelsTotal = len(channels)
		for _, ch := range channels {
			if chMap, ok := ch.(map[string]interface{}); ok {
				if status, ok := chMap["status"].(string); ok && status == "healthy" {
					channelsHealthy++
				}
			}
		}
	}

	snapshot := &models.StatusSnapshot{
		Timestamp:       time.Now(),
		StatusJSON:      output,
		ChannelsHealthy: channelsHealthy,
		ChannelsTotal:   channelsTotal,
		CreatedAt:       time.Now(),
	}

	if err := database.SaveStatusSnapshot(snapshot); err != nil {
		log.Printf("Failed to save status snapshot: %v", err)
		return err
	}

	log.Printf("Status collected: %d/%d channels healthy", channelsHealthy, channelsTotal)
	return nil
}

// CollectHealth 收集健康检查信息
func (c *Collector) CollectHealth() error {
	log.Println("Collecting health...")

	output, err := c.executeCommand("health", "--json")
	if err != nil {
		log.Printf("Failed to collect health: %v", err)
		// 健康检查失败时记录错误状态
		health := &models.HealthCheck{
			Timestamp:     time.Now(),
			OverallStatus: "error",
			DetailsJSON:   fmt.Sprintf(`{"error": "%s"}`, err.Error()),
			CreatedAt:     time.Now(),
		}
		return database.SaveHealthCheck(health)
	}

	// 解析 JSON
	var healthData map[string]interface{}
	if err := json.Unmarshal([]byte(output), &healthData); err != nil {
		log.Printf("Failed to parse health JSON: %v", err)
		return err
	}

	// 确定整体状态
	overallStatus := "healthy"
	if status, ok := healthData["status"].(string); ok {
		overallStatus = status
	}

	health := &models.HealthCheck{
		Timestamp:     time.Now(),
		OverallStatus: overallStatus,
		DetailsJSON:   output,
		CreatedAt:     time.Now(),
	}

	if err := database.SaveHealthCheck(health); err != nil {
		log.Printf("Failed to save health check: %v", err)
		return err
	}

	log.Printf("Health collected: %s", overallStatus)
	return nil
}

// CollectAgents 收集代理信息
func (c *Collector) CollectAgents() error {
	log.Println("Collecting agents...")

	output, err := c.executeCommand("agents", "list", "--json")
	if err != nil {
		log.Printf("Failed to collect agents: %v", err)
		return err
	}

	// 解析 JSON
	var agentsData []map[string]interface{}
	if err := json.Unmarshal([]byte(output), &agentsData); err != nil {
		log.Printf("Failed to parse agents JSON: %v", err)
		return err
	}

	for _, agentData := range agentsData {
		agentID, _ := agentData["id"].(string)
		agentName, _ := agentData["name"].(string)
		status, _ := agentData["status"].(string)

		bindingsJSON, _ := json.Marshal(agentData["bindings"])

		agent := &models.Agent{
			AgentID:      agentID,
			AgentName:    agentName,
			Status:       status,
			BindingsJSON: string(bindingsJSON),
			UpdatedAt:    time.Now(),
		}

		if err := database.SaveOrUpdateAgent(agent); err != nil {
			log.Printf("Failed to save agent %s: %v", agentID, err)
		}
	}

	log.Printf("Agents collected: %d", len(agentsData))
	return nil
}

// CollectAll 收集所有数据
func (c *Collector) CollectAll() {
	log.Println("Starting data collection cycle...")

	// 收集状态
	if err := c.CollectStatus(); err != nil {
		log.Printf("Status collection failed: %v", err)
	}

	// 收集健康检查
	if err := c.CollectHealth(); err != nil {
		log.Printf("Health collection failed: %v", err)
	}

	// 收集代理信息
	if err := c.CollectAgents(); err != nil {
		log.Printf("Agents collection failed: %v", err)
	}

	// 保存指标
	metric := &models.Metric{
		Timestamp:   time.Now(),
		MetricName:  "collection_completed",
		MetricValue: 1,
		MetricType:  "counter",
		CreatedAt:   time.Now(),
	}
	database.SaveMetric(metric)

	log.Println("Data collection cycle completed")
}

// GetOpenClawVersion 获取 OpenClaw 版本
func (c *Collector) GetOpenClawVersion() (string, error) {
	output, err := c.executeCommand("--version")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(output), nil
}
