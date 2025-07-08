package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// ConfigValidator 配置验证器
type ConfigValidator struct {
	errors []string
}

// NewConfigValidator 创建配置验证器
func NewConfigValidator() *ConfigValidator {
	return &ConfigValidator{
		errors: make([]string, 0),
	}
}

// Validate 验证配置
func (cv *ConfigValidator) Validate(config *Config) error {
	cv.errors = cv.errors[:0]

	// 验证服务器配置
	cv.validateServer(config.Server)

	// 验证JWT配置
	cv.validateJWT(config.JWT)

	// 验证数据库配置
	cv.validateDatabase(config.Database)

	// 验证安全配置
	cv.validateSecurity(config.Security)

	// 验证日志配置
	cv.validateLogging(config.Logging)

	// 验证缓存配置
	cv.validateCache(config.Cache)

	// 验证Webhook配置
	cv.validateWebhook(config.Webhook)

	// 验证监控配置
	cv.validateMonitor(config.Monitor)

	// 验证API配置
	cv.validateAPI(config.API)

	if len(cv.errors) > 0 {
		return fmt.Errorf("configuration validation failed: %s", strings.Join(cv.errors, "; "))
	}

	return nil
}

// validateServer 验证服务器配置
func (cv *ConfigValidator) validateServer(config ServerConfig) {
	if config.Port == "" {
		cv.errors = append(cv.errors, "server.port is required")
	}

	if config.ReadTimeout <= 0 {
		cv.errors = append(cv.errors, "server.read_timeout must be positive")
	}

	if config.WriteTimeout <= 0 {
		cv.errors = append(cv.errors, "server.write_timeout must be positive")
	}

	if config.MaxHeaderBytes <= 0 {
		cv.errors = append(cv.errors, "server.max_header_bytes must be positive")
	}

	if config.EnableHTTPS {
		if config.CertFile == "" {
			cv.errors = append(cv.errors, "server.cert_file is required when HTTPS is enabled")
		}
		if config.KeyFile == "" {
			cv.errors = append(cv.errors, "server.key_file is required when HTTPS is enabled")
		}
	}
}

// validateJWT 验证JWT配置
func (cv *ConfigValidator) validateJWT(config JWTConfig) {
	if config.Secret == "" {
		cv.errors = append(cv.errors, "jwt.secret is required")
	}

	if len(config.Secret) < 32 {
		cv.errors = append(cv.errors, "jwt.secret must be at least 32 characters long")
	}

	if config.ExpirationHours <= 0 {
		cv.errors = append(cv.errors, "jwt.expiration_hours must be positive")
	}

	if config.RefreshExpirationHours <= 0 {
		cv.errors = append(cv.errors, "jwt.refresh_expiration_hours must be positive")
	}

	if config.RefreshExpirationHours <= config.ExpirationHours {
		cv.errors = append(cv.errors, "jwt.refresh_expiration_hours must be greater than expiration_hours")
	}

	validAlgorithms := []string{"HS256", "HS384", "HS512", "RS256", "RS384", "RS512"}
	valid := false
	for _, alg := range validAlgorithms {
		if config.Algorithm == alg {
			valid = true
			break
		}
	}
	if !valid {
		cv.errors = append(cv.errors, fmt.Sprintf("jwt.algorithm must be one of: %s", strings.Join(validAlgorithms, ", ")))
	}
}

// validateDatabase 验证数据库配置
func (cv *ConfigValidator) validateDatabase(config DatabaseConfig) {
	if config.Type == "" {
		cv.errors = append(cv.errors, "database.type is required")
	}

	validTypes := []string{"sqlite", "mysql", "postgres"}
	valid := false
	for _, dbType := range validTypes {
		if config.Type == dbType {
			valid = true
			break
		}
	}
	if !valid {
		cv.errors = append(cv.errors, fmt.Sprintf("database.type must be one of: %s", strings.Join(validTypes, ", ")))
	}

	if config.DSN == "" {
		cv.errors = append(cv.errors, "database.dsn is required")
	}

	if config.MaxOpenConns <= 0 {
		cv.errors = append(cv.errors, "database.max_open_conns must be positive")
	}

	if config.MaxIdleConns <= 0 {
		cv.errors = append(cv.errors, "database.max_idle_conns must be positive")
	}

	if config.MaxIdleConns > config.MaxOpenConns {
		cv.errors = append(cv.errors, "database.max_idle_conns cannot be greater than max_open_conns")
	}

	if config.ConnMaxLifetime <= 0 {
		cv.errors = append(cv.errors, "database.conn_max_lifetime must be positive")
	}

	// 验证连接池配置
	if config.ConnectionPool.MaxOpenConns <= 0 {
		cv.errors = append(cv.errors, "database.connection_pool.max_open_conns must be positive")
	}

	if config.ConnectionPool.MaxIdleConns <= 0 {
		cv.errors = append(cv.errors, "database.connection_pool.max_idle_conns must be positive")
	}

	if config.ConnectionPool.ConnMaxLifetime <= 0 {
		cv.errors = append(cv.errors, "database.connection_pool.conn_max_lifetime must be positive")
	}

	// 验证重试配置
	if config.Retry.MaxRetries < 0 {
		cv.errors = append(cv.errors, "database.retry.max_retries must be non-negative")
	}

	if config.Retry.RetryDelay < 0 {
		cv.errors = append(cv.errors, "database.retry.retry_delay must be non-negative")
	}

	validBackoffTypes := []string{"linear", "exponential"}
	valid = false
	for _, backoffType := range validBackoffTypes {
		if config.Retry.BackoffType == backoffType {
			valid = true
			break
		}
	}
	if !valid {
		cv.errors = append(cv.errors, fmt.Sprintf("database.retry.backoff_type must be one of: %s", strings.Join(validBackoffTypes, ", ")))
	}
}

// validateSecurity 验证安全配置
func (cv *ConfigValidator) validateSecurity(config SecurityConfig) {
	if config.BcryptCost < 4 || config.BcryptCost > 31 {
		cv.errors = append(cv.errors, "security.bcrypt_cost must be between 4 and 31")
	}

	if config.APIKeyLength < 16 || config.APIKeyLength > 64 {
		cv.errors = append(cv.errors, "security.api_key_length must be between 16 and 64")
	}

	// 验证密码策略
	if config.PasswordPolicy.MinLength < 4 {
		cv.errors = append(cv.errors, "security.password_policy.min_length must be at least 4")
	}

	if config.PasswordPolicy.MinLength > 128 {
		cv.errors = append(cv.errors, "security.password_policy.min_length cannot exceed 128")
	}
}

// validateLogging 验证日志配置
func (cv *ConfigValidator) validateLogging(config LoggingConfig) {
	validLevels := []string{"debug", "info", "warn", "error", "fatal", "panic"}
	valid := false
	for _, level := range validLevels {
		if config.Level == level {
			valid = true
			break
		}
	}
	if !valid {
		cv.errors = append(cv.errors, fmt.Sprintf("logging.level must be one of: %s", strings.Join(validLevels, ", ")))
	}

	validFormats := []string{"json", "text"}
	valid = false
	for _, format := range validFormats {
		if config.Format == format {
			valid = true
			break
		}
	}
	if !valid {
		cv.errors = append(cv.errors, fmt.Sprintf("logging.format must be one of: %s", strings.Join(validFormats, ", ")))
	}

	validOutputs := []string{"stdout", "stderr", "file"}
	valid = false
	for _, output := range validOutputs {
		if config.Output == output {
			valid = true
			break
		}
	}
	if !valid {
		cv.errors = append(cv.errors, fmt.Sprintf("logging.output must be one of: %s", strings.Join(validOutputs, ", ")))
	}

	if config.Output == "file" && config.FilePath == "" {
		cv.errors = append(cv.errors, "logging.file_path is required when output is file")
	}

	if config.MaxSize <= 0 {
		cv.errors = append(cv.errors, "logging.max_size must be positive")
	}

	if config.MaxBackups < 0 {
		cv.errors = append(cv.errors, "logging.max_backups must be non-negative")
	}

	if config.MaxAge < 0 {
		cv.errors = append(cv.errors, "logging.max_age must be non-negative")
	}
}

// validateCache 验证缓存配置
func (cv *ConfigValidator) validateCache(config CacheConfig) {
	if config.TTL <= 0 {
		cv.errors = append(cv.errors, "cache.ttl must be positive")
	}

	if config.MaxSize <= 0 {
		cv.errors = append(cv.errors, "cache.max_size must be positive")
	}

	// 验证缓存策略
	if config.Strategy.SimpleQueryTTL <= 0 {
		cv.errors = append(cv.errors, "cache.strategy.simple_query_ttl must be positive")
	}

	if config.Strategy.ComplexQueryTTL <= 0 {
		cv.errors = append(cv.errors, "cache.strategy.complex_query_ttl must be positive")
	}

	if config.Strategy.MaxCacheSize <= 0 {
		cv.errors = append(cv.errors, "cache.strategy.max_cache_size must be positive")
	}

	if config.Strategy.HotCacheRatio < 0 || config.Strategy.HotCacheRatio > 1 {
		cv.errors = append(cv.errors, "cache.strategy.hot_cache_ratio must be between 0 and 1")
	}

	if config.Strategy.PromotionThreshold < 1 {
		cv.errors = append(cv.errors, "cache.strategy.promotion_threshold must be positive")
	}

	if config.Strategy.BusinessHoursStart < 0 || config.Strategy.BusinessHoursStart > 23 {
		cv.errors = append(cv.errors, "cache.strategy.business_hours_start must be between 0 and 23")
	}

	if config.Strategy.BusinessHoursEnd < 0 || config.Strategy.BusinessHoursEnd > 23 {
		cv.errors = append(cv.errors, "cache.strategy.business_hours_end must be between 0 and 23")
	}

	if config.Strategy.BusinessHoursEnd <= config.Strategy.BusinessHoursStart {
		cv.errors = append(cv.errors, "cache.strategy.business_hours_end must be greater than business_hours_start")
	}

	if config.Strategy.MaintenanceInterval <= 0 {
		cv.errors = append(cv.errors, "cache.strategy.maintenance_interval must be positive")
	}

	validEvictionPolicies := []string{"lru", "lfu", "fifo"}
	valid := false
	for _, policy := range validEvictionPolicies {
		if config.Strategy.EvictionPolicy == policy {
			valid = true
			break
		}
	}
	if !valid {
		cv.errors = append(cv.errors, fmt.Sprintf("cache.strategy.eviction_policy must be one of: %s", strings.Join(validEvictionPolicies, ", ")))
	}
}

// validateWebhook 验证Webhook配置
func (cv *ConfigValidator) validateWebhook(config WebhookConfig) {
	if config.Timeout <= 0 {
		cv.errors = append(cv.errors, "webhook.timeout must be positive")
	}

	if config.MaxRetries < 0 {
		cv.errors = append(cv.errors, "webhook.max_retries must be non-negative")
	}

	if config.RetryDelay < 0 {
		cv.errors = append(cv.errors, "webhook.retry_delay must be non-negative")
	}

	if config.MaxPayload <= 0 {
		cv.errors = append(cv.errors, "webhook.max_payload must be positive")
	}
}

// validateMonitor 验证监控配置
func (cv *ConfigValidator) validateMonitor(config MonitorConfig) {
	if config.Alerting.Enabled {
		if len(config.Alerting.Channels) == 0 {
			cv.errors = append(cv.errors, "monitor.alerting.channels cannot be empty when alerting is enabled")
		}

		if config.Alerting.Cooldown <= 0 {
			cv.errors = append(cv.errors, "monitor.alerting.cooldown must be positive")
		}

		// 验证阈值
		for metric, threshold := range config.Alerting.Thresholds {
			if threshold < 0 {
				cv.errors = append(cv.errors, fmt.Sprintf("monitor.alerting.thresholds.%s must be non-negative", metric))
			}

			if metric == "error_rate" && threshold > 1 {
				cv.errors = append(cv.errors, "monitor.alerting.thresholds.error_rate must be between 0 and 1")
			}

			if metric == "memory_usage" && threshold > 1 {
				cv.errors = append(cv.errors, "monitor.alerting.thresholds.memory_usage must be between 0 and 1")
			}
		}
	}
}

// validateAPI 验证API配置
func (cv *ConfigValidator) validateAPI(config APIConfig) {
	if config.DefaultLimit <= 0 {
		cv.errors = append(cv.errors, "api.default_limit must be positive")
	}

	if config.MaxLimit <= 0 {
		cv.errors = append(cv.errors, "api.max_limit must be positive")
	}

	if config.DefaultLimit > config.MaxLimit {
		cv.errors = append(cv.errors, "api.default_limit cannot be greater than max_limit")
	}
}

// ConfigExporter 配置导出器
type ConfigExporter struct{}

// NewConfigExporter 创建配置导出器
func NewConfigExporter() *ConfigExporter {
	return &ConfigExporter{}
}

// ExportToYAML 导出配置到YAML文件
func (ce *ConfigExporter) ExportToYAML(config *Config, filepath string) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config to YAML: %w", err)
	}

	if err := os.WriteFile(filepath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// ExportToJSON 导出配置到JSON文件
func (ce *ConfigExporter) ExportToJSON(config *Config, filepath string) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config to JSON: %w", err)
	}

	if err := os.WriteFile(filepath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// ExportToEnv 导出配置到环境变量文件
func (ce *ConfigExporter) ExportToEnv(config *Config, filepath string) error {
	var lines []string

	// 服务器配置
	lines = append(lines, fmt.Sprintf("GOBI_SERVER_PORT=%s", config.Server.Port))
	lines = append(lines, fmt.Sprintf("GOBI_SERVER_HOST=%s", config.Server.Host))

	// JWT配置
	lines = append(lines, fmt.Sprintf("GOBI_JWT_SECRET=%s", config.JWT.Secret))
	lines = append(lines, fmt.Sprintf("GOBI_JWT_EXPIRATION_HOURS=%d", config.JWT.ExpirationHours))

	// 数据库配置
	lines = append(lines, fmt.Sprintf("GOBI_DATABASE_TYPE=%s", config.Database.Type))
	lines = append(lines, fmt.Sprintf("GOBI_DATABASE_DSN=%s", config.Database.DSN))

	// 安全配置
	lines = append(lines, fmt.Sprintf("GOBI_SECURITY_BCRYPT_COST=%d", config.Security.BcryptCost))
	lines = append(lines, fmt.Sprintf("GOBI_SECURITY_RATE_LIMIT=%s", config.Security.RateLimit))

	// 日志配置
	lines = append(lines, fmt.Sprintf("GOBI_LOGGING_LEVEL=%s", config.Logging.Level))
	lines = append(lines, fmt.Sprintf("GOBI_LOGGING_FORMAT=%s", config.Logging.Format))

	// 缓存配置
	lines = append(lines, fmt.Sprintf("GOBI_CACHE_ENABLED=%t", config.Cache.Enabled))
	lines = append(lines, fmt.Sprintf("GOBI_CACHE_TTL=%s", config.Cache.TTL.String()))

	data := strings.Join(lines, "\n") + "\n"

	if err := os.WriteFile(filepath, []byte(data), 0644); err != nil {
		return fmt.Errorf("failed to write env file: %w", err)
	}

	return nil
}

// ConfigImporter 配置导入器
type ConfigImporter struct{}

// NewConfigImporter 创建配置导入器
func NewConfigImporter() *ConfigImporter {
	return &ConfigImporter{}
}

// ImportFromYAML 从YAML文件导入配置
func (ci *ConfigImporter) ImportFromYAML(filepath string) (*Config, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

// ImportFromJSON 从JSON文件导入配置
func (ci *ConfigImporter) ImportFromJSON(filepath string) (*Config, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

// ConfigDiff 配置差异比较器
type ConfigDiff struct {
	Added   map[string]interface{} `json:"added"`
	Removed map[string]interface{} `json:"removed"`
	Changed map[string]interface{} `json:"changed"`
}

// CompareConfigs 比较两个配置的差异
func CompareConfigs(config1, config2 *Config) *ConfigDiff {
	diff := &ConfigDiff{
		Added:   make(map[string]interface{}),
		Removed: make(map[string]interface{}),
		Changed: make(map[string]interface{}),
	}

	// 这里可以实现详细的配置比较逻辑
	// 为了简化，这里只提供基本框架

	return diff
}

// ConfigTemplate 配置模板生成器
type ConfigTemplate struct{}

// NewConfigTemplate 创建配置模板生成器
func NewConfigTemplate() *ConfigTemplate {
	return &ConfigTemplate{}
}

// GenerateTemplate 生成配置模板
func (ct *ConfigTemplate) GenerateTemplate(env string) *Config {
	config := &Config{}

	// 设置默认值
	config.Server.Port = "8080"
	config.Server.Host = "0.0.0.0"
	config.Server.ReadTimeout = 30 * time.Second
	config.Server.WriteTimeout = 30 * time.Second
	config.Server.MaxHeaderBytes = 1 << 20
	config.Server.GracefulTimeout = 30 * time.Second

	config.JWT.ExpirationHours = 24
	config.JWT.RefreshExpirationHours = 168
	config.JWT.Algorithm = "HS256"

	config.Database.MaxOpenConns = 25
	config.Database.MaxIdleConns = 5
	config.Database.ConnMaxLifetime = 300 * time.Second

	config.Security.BcryptCost = 12
	config.Security.APIKeyLength = 32

	config.Logging.Level = "info"
	config.Logging.Format = "json"
	config.Logging.Output = "stdout"

	config.Cache.Enabled = true
	config.Cache.TTL = 300 * time.Second
	config.Cache.MaxSize = 1000

	config.Webhook.Timeout = 30 * time.Second
	config.Webhook.MaxRetries = 3
	config.Webhook.RetryDelay = 5 * time.Second
	config.Webhook.MaxPayload = 1024 * 1024

	config.Monitor.Enabled = true
	config.Monitor.MetricsPort = "9090"
	config.Monitor.HealthCheck = true

	config.API.Version = "v1"
	config.API.Prefix = "/api"
	config.API.DefaultLimit = 20
	config.API.MaxLimit = 100

	// 根据环境调整配置
	switch env {
	case "dev":
		config.Logging.Level = "debug"
		config.Logging.Format = "text"
		config.Security.BcryptCost = 10
		config.Monitor.Profiling = true
		config.API.EnableProfiling = true
	case "prod":
		config.Logging.Level = "warn"
		config.Logging.Output = "file"
		config.Logging.FilePath = "/var/log/gobi/app.log"
		config.Cache.TTL = 600 * time.Second
		config.Cache.MaxSize = 5000
		config.Webhook.MaxRetries = 5
		config.Webhook.RetryDelay = 10 * time.Second
		config.Monitor.Alerting.Enabled = true
		config.API.EnableSwagger = false
	case "test":
		config.Logging.Level = "error"
		config.Cache.Enabled = false
		config.Webhook.Timeout = 5 * time.Second
		config.Webhook.MaxRetries = 1
		config.Webhook.RetryDelay = 1 * time.Second
		config.Monitor.Enabled = false
		config.API.EnableSwagger = false
		config.API.EnableMetrics = false
	}

	return config
}

// SaveTemplate 保存配置模板到文件
func (ct *ConfigTemplate) SaveTemplate(config *Config, filepath string) error {
	exporter := NewConfigExporter()
	return exporter.ExportToYAML(config, filepath)
}
