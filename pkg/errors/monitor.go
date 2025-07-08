package errors

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// ErrorStats 错误统计
type ErrorStats struct {
	TotalErrors    int64                   `json:"total_errors"`
	ErrorCounts    map[ErrorCode]int64     `json:"error_counts"`
	SeverityCounts map[ErrorSeverity]int64 `json:"severity_counts"`
	CategoryCounts map[ErrorCategory]int64 `json:"category_counts"`
	RecentErrors   []*ErrorRecord          `json:"recent_errors"`
	LastErrorTime  time.Time               `json:"last_error_time"`
	ErrorRate      float64                 `json:"error_rate"` // 每分钟错误率
	RetryCount     int64                   `json:"retry_count"`
	SuccessCount   int64                   `json:"success_count"`
	mu             sync.RWMutex
}

// ErrorRecord 错误记录
type ErrorRecord struct {
	Code       ErrorCode              `json:"code"`
	Message    string                 `json:"message"`
	Severity   ErrorSeverity          `json:"severity"`
	Category   ErrorCategory          `json:"category"`
	RequestID  string                 `json:"request_id"`
	Timestamp  time.Time              `json:"timestamp"`
	Details    map[string]interface{} `json:"details,omitempty"`
	Stack      []string               `json:"stack,omitempty"`
	RetryCount int                    `json:"retry_count,omitempty"`
}

// AlertConfig 告警配置
type AlertConfig struct {
	Enabled            bool           `json:"enabled"`
	ErrorRateThreshold float64        `json:"error_rate_threshold"` // 每分钟错误率阈值
	SeverityThreshold  ErrorSeverity  `json:"severity_threshold"`   // 严重程度阈值
	AlertChannels      []AlertChannel `json:"alert_channels"`
	CooldownPeriod     time.Duration  `json:"cooldown_period"` // 告警冷却期
	lastAlertTime      time.Time
	mu                 sync.RWMutex
}

// AlertChannel 告警通道接口
type AlertChannel interface {
	SendAlert(alert *Alert) error
}

// Alert 告警信息
type Alert struct {
	Type      string                 `json:"type"`
	Message   string                 `json:"message"`
	Severity  ErrorSeverity          `json:"severity"`
	Timestamp time.Time              `json:"timestamp"`
	Details   map[string]interface{} `json:"details"`
}

// ErrorMonitor 错误监控器
type ErrorMonitor struct {
	stats    *ErrorStats
	config   *AlertConfig
	channels []AlertChannel
	ctx      context.Context
	cancel   context.CancelFunc
	mu       sync.RWMutex
}

var (
	globalMonitor *ErrorMonitor
	monitorOnce   sync.Once
)

// GetGlobalMonitor 获取全局错误监控器
func GetGlobalMonitor() *ErrorMonitor {
	monitorOnce.Do(func() {
		globalMonitor = NewErrorMonitor()
	})
	return globalMonitor
}

// NewErrorMonitor 创建新的错误监控器
func NewErrorMonitor() *ErrorMonitor {
	ctx, cancel := context.WithCancel(context.Background())

	monitor := &ErrorMonitor{
		stats: &ErrorStats{
			ErrorCounts:    make(map[ErrorCode]int64),
			SeverityCounts: make(map[ErrorSeverity]int64),
			CategoryCounts: make(map[ErrorCategory]int64),
			RecentErrors:   make([]*ErrorRecord, 0, 100),
		},
		config: &AlertConfig{
			Enabled:            true,
			ErrorRateThreshold: 10.0, // 每分钟10个错误
			SeverityThreshold:  SeverityHigh,
			CooldownPeriod:     5 * time.Minute,
		},
		channels: make([]AlertChannel, 0),
		ctx:      ctx,
		cancel:   cancel,
	}

	// 启动监控协程
	go monitor.startMonitoring()

	return monitor
}

// RecordError 记录错误
func (m *ErrorMonitor) RecordError(err *CustomError) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 更新统计信息
	m.stats.TotalErrors++
	m.stats.ErrorCounts[err.Code]++
	m.stats.SeverityCounts[err.Severity]++
	m.stats.CategoryCounts[err.Category]++
	m.stats.LastErrorTime = time.Now()

	// 创建错误记录
	record := &ErrorRecord{
		Code:      err.Code,
		Message:   err.Message,
		Severity:  err.Severity,
		Category:  err.Category,
		RequestID: err.RequestID,
		Timestamp: err.Timestamp,
		Details:   err.Details,
		Stack:     err.Stack,
	}

	// 添加到最近错误列表
	m.stats.RecentErrors = append(m.stats.RecentErrors, record)

	// 保持最近错误列表大小
	if len(m.stats.RecentErrors) > 100 {
		m.stats.RecentErrors = m.stats.RecentErrors[1:]
	}

	// 检查是否需要告警
	m.checkAndSendAlert(err)
}

// RecordSuccess 记录成功操作
func (m *ErrorMonitor) RecordSuccess() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.stats.SuccessCount++
}

// RecordRetry 记录重试操作
func (m *ErrorMonitor) RecordRetry() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.stats.RetryCount++
}

// GetStats 获取统计信息
func (m *ErrorMonitor) GetStats() *ErrorStats {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// 计算错误率
	now := time.Now()
	if !m.stats.LastErrorTime.IsZero() {
		timeDiff := now.Sub(m.stats.LastErrorTime).Minutes()
		if timeDiff > 0 {
			m.stats.ErrorRate = float64(m.stats.TotalErrors) / timeDiff
		}
	}

	// 创建副本
	stats := *m.stats
	stats.ErrorCounts = make(map[ErrorCode]int64)
	stats.SeverityCounts = make(map[ErrorSeverity]int64)
	stats.CategoryCounts = make(map[ErrorCategory]int64)

	for k, v := range m.stats.ErrorCounts {
		stats.ErrorCounts[k] = v
	}
	for k, v := range m.stats.SeverityCounts {
		stats.SeverityCounts[k] = v
	}
	for k, v := range m.stats.CategoryCounts {
		stats.CategoryCounts[k] = v
	}

	return &stats
}

// AddAlertChannel 添加告警通道
func (m *ErrorMonitor) AddAlertChannel(channel AlertChannel) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.channels = append(m.channels, channel)
}

// SetAlertConfig 设置告警配置
func (m *ErrorMonitor) SetAlertConfig(config *AlertConfig) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.config = config
}

// checkAndSendAlert 检查并发送告警
func (m *ErrorMonitor) checkAndSendAlert(err *CustomError) {
	if !m.config.Enabled {
		return
	}

	// 检查冷却期
	m.config.mu.RLock()
	if time.Since(m.config.lastAlertTime) < m.config.CooldownPeriod {
		m.config.mu.RUnlock()
		return
	}
	m.config.mu.RUnlock()

	// 检查严重程度
	if err.Severity < m.config.SeverityThreshold {
		return
	}

	// 检查错误率
	if m.stats.ErrorRate > m.config.ErrorRateThreshold {
		m.sendAlert(&Alert{
			Type:      "HIGH_ERROR_RATE",
			Message:   fmt.Sprintf("High error rate detected: %.2f errors/minute", m.stats.ErrorRate),
			Severity:  err.Severity,
			Timestamp: time.Now(),
			Details: map[string]interface{}{
				"error_rate":   m.stats.ErrorRate,
				"threshold":    m.config.ErrorRateThreshold,
				"total_errors": m.stats.TotalErrors,
			},
		})
		return
	}

	// 检查严重错误
	if err.Severity >= SeverityCritical {
		m.sendAlert(&Alert{
			Type:      "CRITICAL_ERROR",
			Message:   fmt.Sprintf("Critical error occurred: %s", err.Message),
			Severity:  err.Severity,
			Timestamp: time.Now(),
			Details: map[string]interface{}{
				"error_code":    string(err.Code),
				"error_message": err.Message,
				"category":      string(err.Category),
				"request_id":    err.RequestID,
			},
		})
		return
	}
}

// sendAlert 发送告警
func (m *ErrorMonitor) sendAlert(alert *Alert) {
	m.config.mu.Lock()
	m.config.lastAlertTime = time.Now()
	m.config.mu.Unlock()

	for _, channel := range m.channels {
		go func(ch AlertChannel) {
			if err := ch.SendAlert(alert); err != nil {
				// 记录告警发送失败
				fmt.Printf("Failed to send alert: %v\n", err)
			}
		}(channel)
	}
}

// startMonitoring 启动监控
func (m *ErrorMonitor) startMonitoring() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.updateErrorRate()
		case <-m.ctx.Done():
			return
		}
	}
}

// updateErrorRate 更新错误率
func (m *ErrorMonitor) updateErrorRate() {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	if !m.stats.LastErrorTime.IsZero() {
		timeDiff := now.Sub(m.stats.LastErrorTime).Minutes()
		if timeDiff > 0 {
			m.stats.ErrorRate = float64(m.stats.TotalErrors) / timeDiff
		}
	}
}

// Stop 停止监控
func (m *ErrorMonitor) Stop() {
	m.cancel()
}

// Reset 重置统计信息
func (m *ErrorMonitor) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.stats = &ErrorStats{
		ErrorCounts:    make(map[ErrorCode]int64),
		SeverityCounts: make(map[ErrorSeverity]int64),
		CategoryCounts: make(map[ErrorCategory]int64),
		RecentErrors:   make([]*ErrorRecord, 0, 100),
	}
}

// 便捷函数
func RecordError(err *CustomError) {
	if globalMonitor != nil {
		globalMonitor.RecordError(err)
	}
}

func RecordSuccess() {
	if globalMonitor != nil {
		globalMonitor.RecordSuccess()
	}
}

func RecordRetry() {
	if globalMonitor != nil {
		globalMonitor.RecordRetry()
	}
}

func GetErrorStats() *ErrorStats {
	if globalMonitor != nil {
		return globalMonitor.GetStats()
	}
	return &ErrorStats{}
}
