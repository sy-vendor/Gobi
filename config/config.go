package config

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Config 应用配置结构
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Database DatabaseConfig `mapstructure:"database"`
	Security SecurityConfig `mapstructure:"security"`
	Logging  LoggingConfig  `mapstructure:"logging"`
	Cache    CacheConfig    `mapstructure:"cache"`
	Webhook  WebhookConfig  `mapstructure:"webhook"`
	Monitor  MonitorConfig  `mapstructure:"monitor"`
	API      APIConfig      `mapstructure:"api"`
	AI       AIConfig       `mapstructure:"ai"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port            string        `mapstructure:"port"`
	Host            string        `mapstructure:"host"`
	ReadTimeout     time.Duration `mapstructure:"read_timeout"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout"`
	MaxHeaderBytes  int           `mapstructure:"max_header_bytes"`
	GracefulTimeout time.Duration `mapstructure:"graceful_timeout"`
	EnableHTTPS     bool          `mapstructure:"enable_https"`
	CertFile        string        `mapstructure:"cert_file"`
	KeyFile         string        `mapstructure:"key_file"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret                 string `mapstructure:"secret"`
	ExpirationHours        int    `mapstructure:"expiration_hours"`
	RefreshExpirationHours int    `mapstructure:"refresh_expiration_hours"`
	Issuer                 string `mapstructure:"issuer"`
	Audience               string `mapstructure:"audience"`
	Algorithm              string `mapstructure:"algorithm"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Type            string              `mapstructure:"type"`
	DSN             string              `mapstructure:"dsn"`
	MaxOpenConns    int                 `mapstructure:"max_open_conns"`
	MaxIdleConns    int                 `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration       `mapstructure:"conn_max_lifetime"`
	ConnectionPool  DatabasePoolConfig  `mapstructure:"connection_pool"`
	SSL             DatabaseSSLConfig   `mapstructure:"ssl"`
	Retry           DatabaseRetryConfig `mapstructure:"retry"`
	Migration       MigrationConfig     `mapstructure:"migration"`
}

// DatabasePoolConfig 数据库连接池配置
type DatabasePoolConfig struct {
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `mapstructure:"conn_max_idle_time"`
}

// DatabaseSSLConfig 数据库SSL配置
type DatabaseSSLConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	Mode     string `mapstructure:"mode"`
	CertFile string `mapstructure:"cert_file"`
	KeyFile  string `mapstructure:"key_file"`
	CAFile   string `mapstructure:"ca_file"`
}

// DatabaseRetryConfig 数据库重试配置
type DatabaseRetryConfig struct {
	MaxRetries  int           `mapstructure:"max_retries"`
	RetryDelay  time.Duration `mapstructure:"retry_delay"`
	BackoffType string        `mapstructure:"backoff_type"` // linear, exponential
}

// MigrationConfig 数据库迁移配置
type MigrationConfig struct {
	AutoMigrate bool   `mapstructure:"auto_migrate"`
	Path        string `mapstructure:"path"`
	TableName   string `mapstructure:"table_name"`
}

// SecurityConfig 安全配置
type SecurityConfig struct {
	BcryptCost     int                  `mapstructure:"bcrypt_cost"`
	RateLimit      string               `mapstructure:"rate_limit"`
	CORSOrigins    []string             `mapstructure:"cors_origins"`
	AllowedHeaders []string             `mapstructure:"allowed_headers"`
	AllowedMethods []string             `mapstructure:"allowed_methods"`
	TrustedProxies []string             `mapstructure:"trusted_proxies"`
	APIKeyLength   int                  `mapstructure:"api_key_length"`
	PasswordPolicy PasswordPolicyConfig `mapstructure:"password_policy"`
}

// PasswordPolicyConfig 密码策略配置
type PasswordPolicyConfig struct {
	MinLength        int  `mapstructure:"min_length"`
	RequireUppercase bool `mapstructure:"require_uppercase"`
	RequireLowercase bool `mapstructure:"require_lowercase"`
	RequireNumbers   bool `mapstructure:"require_numbers"`
	RequireSymbols   bool `mapstructure:"require_symbols"`
}

// LoggingConfig 日志配置
type LoggingConfig struct {
	Level      string `mapstructure:"level"`
	Format     string `mapstructure:"format"`
	Output     string `mapstructure:"output"`
	FilePath   string `mapstructure:"file_path"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
	Compress   bool   `mapstructure:"compress"`
	Console    bool   `mapstructure:"console"`
}

// CacheConfig 缓存配置
type CacheConfig struct {
	Enabled  bool                `mapstructure:"enabled"`
	TTL      time.Duration       `mapstructure:"ttl"`
	MaxSize  int                 `mapstructure:"max_size"`
	Strategy CacheStrategyConfig `mapstructure:"strategy"`
}

// CacheStrategyConfig 缓存策略配置
type CacheStrategyConfig struct {
	SimpleQueryTTL      time.Duration `mapstructure:"simple_query_ttl"`
	ComplexQueryTTL     time.Duration `mapstructure:"complex_query_ttl"`
	MaxCacheSize        int           `mapstructure:"max_cache_size"`
	HotCacheEnabled     bool          `mapstructure:"hot_cache_enabled"`
	HotCacheRatio       float64       `mapstructure:"hot_cache_ratio"`
	PromotionThreshold  int           `mapstructure:"promotion_threshold"`
	BusinessHoursStart  int           `mapstructure:"business_hours_start"`
	BusinessHoursEnd    int           `mapstructure:"business_hours_end"`
	AdaptiveTTL         bool          `mapstructure:"adaptive_ttl"`
	CacheWarmup         bool          `mapstructure:"cache_warmup"`
	MaintenanceInterval time.Duration `mapstructure:"maintenance_interval"`
	EvictionPolicy      string        `mapstructure:"eviction_policy"`
	CompressionEnabled  bool          `mapstructure:"compression_enabled"`
	MetricsEnabled      bool          `mapstructure:"metrics_enabled"`
}

// WebhookConfig Webhook配置
type WebhookConfig struct {
	Timeout    time.Duration     `mapstructure:"timeout"`
	MaxRetries int               `mapstructure:"max_retries"`
	RetryDelay time.Duration     `mapstructure:"retry_delay"`
	MaxPayload int               `mapstructure:"max_payload"`
	VerifySSL  bool              `mapstructure:"verify_ssl"`
	Headers    map[string]string `mapstructure:"headers"`
}

// MonitorConfig 监控配置
type MonitorConfig struct {
	Enabled     bool           `mapstructure:"enabled"`
	MetricsPort string         `mapstructure:"metrics_port"`
	HealthCheck bool           `mapstructure:"health_check"`
	Profiling   bool           `mapstructure:"profiling"`
	Alerting    AlertingConfig `mapstructure:"alerting"`
}

// AlertingConfig 告警配置
type AlertingConfig struct {
	Enabled    bool               `mapstructure:"enabled"`
	Channels   []string           `mapstructure:"channels"`
	Thresholds map[string]float64 `mapstructure:"thresholds"`
	Cooldown   time.Duration      `mapstructure:"cooldown"`
}

// APIConfig API配置
type APIConfig struct {
	Version         string `mapstructure:"version"`
	Prefix          string `mapstructure:"prefix"`
	DefaultLimit    int    `mapstructure:"default_limit"`
	MaxLimit        int    `mapstructure:"max_limit"`
	EnableSwagger   bool   `mapstructure:"enable_swagger"`
	EnableMetrics   bool   `mapstructure:"enable_metrics"`
	EnableProfiling bool   `mapstructure:"enable_profiling"`
}

type AIConfig struct {
	DeepSeekAPIKey string `mapstructure:"deepseek_api_key"`
}

// ConfigManager 配置管理器
type ConfigManager struct {
	config     *Config
	viper      *viper.Viper
	watcher    *fsnotify.Watcher
	mutex      sync.RWMutex
	callbacks  []func(*Config)
	configPath string
	env        string
}

var (
	AppConfig     *Config
	configManager *ConfigManager
	once          sync.Once
)

// NewConfigManager 创建配置管理器
func NewConfigManager() *ConfigManager {
	return &ConfigManager{
		viper:     viper.New(),
		callbacks: make([]func(*Config), 0),
	}
}

// LoadConfig 加载配置
func LoadConfig() error {
	var err error
	once.Do(func() {
		configManager = NewConfigManager()
		err = configManager.load()
	})
	return err
}

// GetConfig 获取配置
func GetConfig() *Config {
	configManager.mutex.RLock()
	defer configManager.mutex.RUnlock()
	return configManager.config
}

// OnConfigChange 注册配置变更回调
func OnConfigChange(callback func(*Config)) {
	configManager.mutex.Lock()
	defer configManager.mutex.Unlock()
	configManager.callbacks = append(configManager.callbacks, callback)
}

// load 加载配置
func (cm *ConfigManager) load() error {
	// 设置环境
	cm.env = getEnvironment()

	// 设置配置文件路径
	cm.configPath = getConfigPath()

	// 配置Viper
	cm.setupViper()

	// 读取配置文件
	if err := cm.viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	// 根据环境获取对应的配置节
	cm.viper.SetConfigName("config")

	// 解析配置
	config := &Config{}
	if err := cm.viper.UnmarshalKey(cm.env, config); err != nil {
		return fmt.Errorf("failed to unmarshal config for env %s: %w", cm.env, err)
	}

	// 设置默认值
	cm.setDefaults(config)

	// 处理敏感信息（在验证之前）
	cm.processSensitiveData(config)

	// 验证配置
	if err := cm.validateConfig(config); err != nil {
		return fmt.Errorf("config validation failed: %w", err)
	}

	// 更新配置
	cm.mutex.Lock()
	cm.config = config
	AppConfig = config
	cm.mutex.Unlock()

	// 自动设置AI Key到环境变量
	if config.AI.DeepSeekAPIKey != "" {
		os.Setenv("DEEPSEEK_API_KEY", config.AI.DeepSeekAPIKey)
	}

	// 启动文件监听
	go cm.watchConfigFile()

	fmt.Printf("Loaded config for env: %s, port: %s, db type: %s\n",
		cm.env, config.Server.Port, config.Database.Type)

	return nil
}

// setupViper 设置Viper
func (cm *ConfigManager) setupViper() {
	cm.viper.SetConfigName("config")
	cm.viper.SetConfigType("yaml")
	cm.viper.AddConfigPath(cm.configPath)
	cm.viper.AddConfigPath(".")
	cm.viper.AddConfigPath("./config")

	// 环境变量支持
	cm.viper.SetEnvPrefix("GOBI")
	cm.viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	cm.viper.AutomaticEnv()

	// 绑定环境变量
	cm.bindEnvironmentVariables()
}

// bindEnvironmentVariables 绑定环境变量
func (cm *ConfigManager) bindEnvironmentVariables() {
	// 服务器配置
	cm.viper.BindEnv("server.port", "GOBI_SERVER_PORT")
	cm.viper.BindEnv("server.host", "GOBI_SERVER_HOST")

	// JWT配置
	cm.viper.BindEnv("jwt.secret", "GOBI_JWT_SECRET")
	cm.viper.BindEnv("jwt.expiration_hours", "GOBI_JWT_EXPIRATION_HOURS")

	// 数据库配置
	cm.viper.BindEnv("database.type", "GOBI_DATABASE_TYPE")
	cm.viper.BindEnv("database.dsn", "GOBI_DATABASE_DSN")

	// 安全配置
	cm.viper.BindEnv("security.bcrypt_cost", "GOBI_SECURITY_BCRYPT_COST")
	cm.viper.BindEnv("security.rate_limit", "GOBI_SECURITY_RATE_LIMIT")
}

// setDefaults 设置默认值
func (cm *ConfigManager) setDefaults(config *Config) {
	// 服务器默认值
	if config.Server.Port == "" {
		config.Server.Port = "8080"
	}
	if config.Server.Host == "" {
		config.Server.Host = "0.0.0.0"
	}
	if config.Server.ReadTimeout == 0 {
		config.Server.ReadTimeout = 30 * time.Second
	}
	if config.Server.WriteTimeout == 0 {
		config.Server.WriteTimeout = 30 * time.Second
	}
	if config.Server.MaxHeaderBytes == 0 {
		config.Server.MaxHeaderBytes = 1 << 20 // 1MB
	}
	if config.Server.GracefulTimeout == 0 {
		config.Server.GracefulTimeout = 30 * time.Second
	}

	// JWT默认值
	if config.JWT.ExpirationHours == 0 {
		config.JWT.ExpirationHours = 24
	}
	if config.JWT.RefreshExpirationHours == 0 {
		config.JWT.RefreshExpirationHours = 168 // 7天
	}
	if config.JWT.Algorithm == "" {
		config.JWT.Algorithm = "HS256"
	}

	// 数据库默认值
	if config.Database.Type == "" {
		config.Database.Type = "sqlite"
	}
	if config.Database.DSN == "" {
		config.Database.DSN = "gobi.db"
	}
	if config.Database.MaxOpenConns == 0 {
		config.Database.MaxOpenConns = 25
	}
	if config.Database.MaxIdleConns == 0 {
		config.Database.MaxIdleConns = 5
	}
	if config.Database.ConnMaxLifetime == 0 {
		config.Database.ConnMaxLifetime = 300 * time.Second
	}

	// 安全默认值
	if config.Security.BcryptCost == 0 {
		config.Security.BcryptCost = 12
	}
	if config.Security.RateLimit == "" {
		config.Security.RateLimit = "100-M"
	}
	if config.Security.APIKeyLength == 0 {
		config.Security.APIKeyLength = 32
	}

	// 日志默认值
	if config.Logging.Level == "" {
		config.Logging.Level = "info"
	}
	if config.Logging.Format == "" {
		config.Logging.Format = "json"
	}
	if config.Logging.Output == "" {
		config.Logging.Output = "stdout"
	}

	// 缓存默认值
	if config.Cache.TTL == 0 {
		config.Cache.TTL = 300 * time.Second
	}
	if config.Cache.MaxSize == 0 {
		config.Cache.MaxSize = 1000
	}

	// 缓存策略默认值
	if config.Cache.Strategy.SimpleQueryTTL == 0 {
		config.Cache.Strategy.SimpleQueryTTL = 300 * time.Second
	}
	if config.Cache.Strategy.ComplexQueryTTL == 0 {
		config.Cache.Strategy.ComplexQueryTTL = 600 * time.Second
	}
	if config.Cache.Strategy.MaxCacheSize == 0 {
		config.Cache.Strategy.MaxCacheSize = 1000
	}
	if config.Cache.Strategy.HotCacheRatio == 0 {
		config.Cache.Strategy.HotCacheRatio = 0.2
	}
	if config.Cache.Strategy.PromotionThreshold == 0 {
		config.Cache.Strategy.PromotionThreshold = 3
	}
	if config.Cache.Strategy.BusinessHoursStart == 0 {
		config.Cache.Strategy.BusinessHoursStart = 9
	}
	if config.Cache.Strategy.BusinessHoursEnd == 0 {
		config.Cache.Strategy.BusinessHoursEnd = 17
	}
	if config.Cache.Strategy.MaintenanceInterval == 0 {
		config.Cache.Strategy.MaintenanceInterval = 300 * time.Second
	}
	if config.Cache.Strategy.EvictionPolicy == "" {
		config.Cache.Strategy.EvictionPolicy = "lru"
	}

	// Webhook默认值
	if config.Webhook.Timeout == 0 {
		config.Webhook.Timeout = 30 * time.Second
	}
	if config.Webhook.MaxRetries == 0 {
		config.Webhook.MaxRetries = 3
	}
	if config.Webhook.RetryDelay == 0 {
		config.Webhook.RetryDelay = 5 * time.Second
	}
	if config.Webhook.MaxPayload == 0 {
		config.Webhook.MaxPayload = 1024 * 1024 // 1MB
	}
}

// validateConfig 验证配置
func (cm *ConfigManager) validateConfig(config *Config) error {
	var errors []string

	// 验证服务器配置
	if config.Server.Port == "" {
		errors = append(errors, "server.port is required")
	}
	if port, err := strconv.Atoi(config.Server.Port); err != nil || port <= 0 || port > 65535 {
		errors = append(errors, "server.port must be a valid port number (1-65535)")
	}

	// 验证JWT配置
	if config.JWT.Secret == "" {
		errors = append(errors, "jwt.secret is required")
	}
	if len(config.JWT.Secret) < 32 {
		errors = append(errors, "jwt.secret must be at least 32 characters long")
	}
	if config.JWT.ExpirationHours <= 0 {
		errors = append(errors, "jwt.expiration_hours must be positive")
	}

	// 验证数据库配置
	if config.Database.Type == "" {
		errors = append(errors, "database.type is required")
	}
	if !isValidDatabaseType(config.Database.Type) {
		errors = append(errors, "database.type must be one of: sqlite, mysql, postgres")
	}
	if config.Database.DSN == "" {
		errors = append(errors, "database.dsn is required")
	}

	// 验证安全配置
	if config.Security.BcryptCost < 4 || config.Security.BcryptCost > 31 {
		errors = append(errors, "security.bcrypt_cost must be between 4 and 31")
	}
	if config.Security.APIKeyLength < 16 || config.Security.APIKeyLength > 64 {
		errors = append(errors, "security.api_key_length must be between 16 and 64")
	}

	// 验证缓存配置
	if config.Cache.Strategy.HotCacheRatio < 0 || config.Cache.Strategy.HotCacheRatio > 1 {
		errors = append(errors, "cache.strategy.hot_cache_ratio must be between 0 and 1")
	}
	if config.Cache.Strategy.PromotionThreshold < 1 {
		errors = append(errors, "cache.strategy.promotion_threshold must be positive")
	}

	// 验证Webhook配置
	if config.Webhook.MaxRetries < 0 {
		errors = append(errors, "webhook.max_retries must be non-negative")
	}
	if config.Webhook.MaxPayload < 1024 {
		errors = append(errors, "webhook.max_payload must be at least 1KB")
	}

	if len(errors) > 0 {
		return fmt.Errorf("configuration validation failed: %s", strings.Join(errors, "; "))
	}

	return nil
}

// processSensitiveData 处理敏感数据
func (cm *ConfigManager) processSensitiveData(config *Config) {
	// 生成JWT密钥（如果未设置）
	if config.JWT.Secret == "" || strings.Contains(config.JWT.Secret, "default") {
		config.JWT.Secret = generateSecureSecret(32)
	}

	// 处理数据库DSN中的敏感信息
	if strings.Contains(config.Database.DSN, "password") {
		// 在生产环境中，应该从环境变量或密钥管理服务获取
		fmt.Println("Warning: Database DSN contains password, consider using environment variables")
	}
}

// watchConfigFile 监听配置文件变化
func (cm *ConfigManager) watchConfigFile() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Printf("Failed to create config watcher: %v\n", err)
		return
	}
	defer watcher.Close()

	cm.watcher = watcher

	// 监听配置文件目录
	if err := watcher.Add(cm.configPath); err != nil {
		fmt.Printf("Failed to watch config directory: %v\n", err)
		return
	}

	fmt.Println("Config file watcher started")

	for {
		select {
		case event := <-watcher.Events:
			if event.Op&fsnotify.Write == fsnotify.Write {
				if strings.HasSuffix(event.Name, "config.yaml") {
					fmt.Println("Config file changed, reloading...")
					if err := cm.reload(); err != nil {
						fmt.Printf("Failed to reload config: %v\n", err)
					}
				}
			}
		case err := <-watcher.Errors:
			fmt.Printf("Config watcher error: %v\n", err)
		}
	}
}

// reload 重新加载配置
func (cm *ConfigManager) reload() error {
	// 重新读取配置文件
	if err := cm.viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	// 解析配置
	config := &Config{}
	if err := cm.viper.Unmarshal(config); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// 设置默认值
	cm.setDefaults(config)

	// 验证配置
	if err := cm.validateConfig(config); err != nil {
		return fmt.Errorf("config validation failed: %w", err)
	}

	// 处理敏感数据
	cm.processSensitiveData(config)

	// 更新配置
	cm.mutex.Lock()
	cm.config = config
	AppConfig = config
	cm.mutex.Unlock()

	// 通知回调函数
	cm.notifyCallbacks(config)

	fmt.Printf("Config reloaded successfully for env: %s\n", cm.env)

	return nil
}

// notifyCallbacks 通知配置变更回调
func (cm *ConfigManager) notifyCallbacks(config *Config) {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	for _, callback := range cm.callbacks {
		go callback(config)
	}
}

// 辅助函数
func getEnvironment() string {
	env := os.Getenv("GOBI_ENV")
	if env == "" {
		env = "default"
	}
	return env
}

func getConfigPath() string {
	configPath := os.Getenv("GOBI_CONFIG_PATH")
	if configPath == "" {
		configPath = "./config"
	}
	return configPath
}

func isValidDatabaseType(dbType string) bool {
	validTypes := []string{"sqlite", "mysql", "postgres"}
	for _, validType := range validTypes {
		if dbType == validType {
			return true
		}
	}
	return false
}

func generateSecureSecret(length int) string {
	bytes := make([]byte, length)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// 配置工具函数
func (c *Config) IsDevelopment() bool {
	return getEnvironment() == "dev"
}

func (c *Config) IsProduction() bool {
	return getEnvironment() == "prod"
}

func (c *Config) IsTest() bool {
	return getEnvironment() == "test"
}

func (c *Config) GetDatabaseDSN() string {
	return c.Database.DSN
}

func (c *Config) GetJWTSecret() string {
	return c.JWT.Secret
}

func (c *Config) GetServerPort() string {
	return c.Server.Port
}

func (c *Config) GetCacheTTL() time.Duration {
	return c.Cache.TTL
}

func (c *Config) GetRateLimit() string {
	return c.Security.RateLimit
}

// 配置验证函数
func (c *Config) ValidatePassword(password string) error {
	policy := c.Security.PasswordPolicy

	if len(password) < policy.MinLength {
		return fmt.Errorf("password must be at least %d characters long", policy.MinLength)
	}

	if policy.RequireUppercase && !containsUppercase(password) {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}

	if policy.RequireLowercase && !containsLowercase(password) {
		return fmt.Errorf("password must contain at least one lowercase letter")
	}

	if policy.RequireNumbers && !containsNumbers(password) {
		return fmt.Errorf("password must contain at least one number")
	}

	if policy.RequireSymbols && !containsSymbols(password) {
		return fmt.Errorf("password must contain at least one symbol")
	}

	return nil
}

func containsUppercase(s string) bool {
	for _, r := range s {
		if r >= 'A' && r <= 'Z' {
			return true
		}
	}
	return false
}

func containsLowercase(s string) bool {
	for _, r := range s {
		if r >= 'a' && r <= 'z' {
			return true
		}
	}
	return false
}

func containsNumbers(s string) bool {
	for _, r := range s {
		if r >= '0' && r <= '9' {
			return true
		}
	}
	return false
}

func containsSymbols(s string) bool {
	symbols := "!@#$%^&*()_+-=[]{}|;:,.<>?"
	for _, r := range s {
		if strings.ContainsRune(symbols, r) {
			return true
		}
	}
	return false
}
