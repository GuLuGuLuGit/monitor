package scheduler

import (
	"fmt"
	"log"

	"github.com/monitor/openclaw-monitor/internal/collector"
	"github.com/monitor/openclaw-monitor/internal/database"
	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	cron      *cron.Cron
	collector *collector.Collector
	interval  int
}

func NewScheduler(c *collector.Collector, interval int) *Scheduler {
	return &Scheduler{
		cron:      cron.New(),
		collector: c,
		interval:  interval,
	}
}

// Start 启动调度器
func (s *Scheduler) Start() {
	log.Printf("Starting scheduler with %d second interval", s.interval)

	// 立即执行一次
	go s.collector.CollectAll()

	// 添加定时任务 - 每 N 秒执行一次
	spec := fmt.Sprintf("@every %ds", s.interval)
	_, err := s.cron.AddFunc(spec, func() {
		s.collector.CollectAll()
	})
	if err != nil {
		log.Fatalf("Failed to add cron job: %v", err)
	}

	// 添加数据清理任务 - 每天凌晨 2 点执行
	_, err = s.cron.AddFunc("0 2 * * *", func() {
		log.Println("Running daily cleanup...")
		if err := database.CleanupOldData(30); err != nil {
			log.Printf("Cleanup failed: %v", err)
		}
	})
	if err != nil {
		log.Fatalf("Failed to add cleanup job: %v", err)
	}

	s.cron.Start()
	log.Println("Scheduler started")
}

// Stop 停止调度器
func (s *Scheduler) Stop() {
	log.Println("Stopping scheduler...")
	s.cron.Stop()
	log.Println("Scheduler stopped")
}
