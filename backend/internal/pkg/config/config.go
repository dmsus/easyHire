package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	AI       AIConfig
	CodeExec CodeExecutionConfig
	Email    EmailConfig
	Security SecurityConfig
	Monitoring MonitoringConfig
}

type ServerConfig struct {
	Host         string
	Port         int
	Mode         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type DatabaseConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	Name            string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

type RedisConfig struct {
	Host       string
	Port       int
	Password   string
	DB         int
	PoolSize   int
}

type JWTConfig struct {
	Secret           string
	AccessTokenTTL   time.Duration
	RefreshTokenTTL  time.Duration
	Issuer           string
}

type AIConfig struct {
	ServiceURL      string
	Timeout         time.Duration
	GeminiAPIKey    string
	OpenRouterAPIKey string
}

type CodeExecutionConfig struct {
	ExecutorURL     string
	Timeout         time.Duration
	MemoryLimitMB   int
	CPULimit        float64
}

type EmailConfig struct {
	SMTPHost     string
	SMTPPort     int
	SMTPUser     string
	SMTPPassword string
	SMTPFrom     string
}

type SecurityConfig struct {
	CORSAllowedOrigins  []string
	CORSAllowedMethods  []string
	CORSAllowedHeaders  []string
	RateLimitEnabled    bool
	RateLimitRequests   int
}

type MonitoringConfig struct {
	MetricsEnabled     bool
	MetricsPort        int
	HealthCheckEnabled bool
	TracingEnabled     bool
}

func LoadConfig(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	// Set defaults
	setDefaults()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	return parseConfig()
}

func setDefaults() {
	viper.SetDefault("SERVER_HOST", "0.0.0.0")
	viper.SetDefault("SERVER_PORT", 8080)
	viper.SetDefault("SERVER_MODE", "debug")
	viper.SetDefault("SERVER_READ_TIMEOUT", 30)
	viper.SetDefault("SERVER_WRITE_TIMEOUT", 30)
	viper.SetDefault("SERVER_IDLE_TIMEOUT", 120)

	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", 5432)
	viper.SetDefault("DB_USER", "easyhire")
	viper.SetDefault("DB_NAME", "easyhire")
	viper.SetDefault("DB_SSL_MODE", "disable")
	viper.SetDefault("DB_MAX_OPEN_CONNS", 25)
	viper.SetDefault("DB_MAX_IDLE_CONNS", 5)
	viper.SetDefault("DB_CONN_MAX_LIFETIME", 300)

	viper.SetDefault("REDIS_HOST", "localhost")
	viper.SetDefault("REDIS_PORT", 6379)
	viper.SetDefault("REDIS_DB", 0)
	viper.SetDefault("REDIS_POOL_SIZE", 10)

	viper.SetDefault("JWT_ACCESS_TOKEN_TTL", "15m")
	viper.SetDefault("JWT_REFRESH_TOKEN_TTL", "24h")
	viper.SetDefault("JWT_ISSUER", "easyhire-api")

	viper.SetDefault("AI_SERVICE_URL", "http://localhost:8000")
	viper.SetDefault("AI_SERVICE_TIMEOUT", 30)
	viper.SetDefault("CODE_EXECUTOR_URL", "http://localhost:8081")
	viper.SetDefault("CODE_EXECUTOR_TIMEOUT", 30)
	viper.SetDefault("CODE_EXECUTION_MEMORY_LIMIT_MB", 256)
	viper.SetDefault("CODE_EXECUTION_CPU_LIMIT", 0.5)

	viper.SetDefault("SMTP_PORT", 587)
	viper.SetDefault("SMTP_FROM", "support@easyhire.com")

	viper.SetDefault("METRICS_ENABLED", true)
	viper.SetDefault("METRICS_PORT", 9090)
	viper.SetDefault("HEALTH_CHECK_ENABLED", true)
	viper.SetDefault("TRACING_ENABLED", false)

	viper.SetDefault("CORS_ALLOWED_ORIGINS", "http://localhost:3000")
	viper.SetDefault("CORS_ALLOWED_METHODS", "GET,POST,PUT,DELETE,PATCH,OPTIONS")
	viper.SetDefault("CORS_ALLOWED_HEADERS", "Origin,Content-Type,Accept,Authorization")
	viper.SetDefault("RATE_LIMIT_ENABLED", true)
	viper.SetDefault("RATE_LIMIT_REQUESTS_PER_MINUTE", 60)
}

func parseConfig() (*Config, error) {
	var config Config

	// Parse durations
	readTimeout, err := time.ParseDuration(fmt.Sprintf("%ds", viper.GetInt("SERVER_READ_TIMEOUT")))
	if err != nil {
		return nil, err
	}
	writeTimeout, err := time.ParseDuration(fmt.Sprintf("%ds", viper.GetInt("SERVER_WRITE_TIMEOUT")))
	if err != nil {
		return nil, err
	}
	idleTimeout, err := time.ParseDuration(fmt.Sprintf("%ds", viper.GetInt("SERVER_IDLE_TIMEOUT")))
	if err != nil {
		return nil, err
	}

	accessTokenTTL, err := time.ParseDuration(viper.GetString("JWT_ACCESS_TOKEN_TTL"))
	if err != nil {
		return nil, err
	}
	refreshTokenTTL, err := time.ParseDuration(viper.GetString("JWT_REFRESH_TOKEN_TTL"))
	if err != nil {
		return nil, err
	}

	aiTimeout, err := time.ParseDuration(fmt.Sprintf("%ds", viper.GetInt("AI_SERVICE_TIMEOUT")))
	if err != nil {
		return nil, err
	}
	codeExecTimeout, err := time.ParseDuration(fmt.Sprintf("%ds", viper.GetInt("CODE_EXECUTOR_TIMEOUT")))
	if err != nil {
		return nil, err
	}
	connMaxLifetime, err := time.ParseDuration(fmt.Sprintf("%ds", viper.GetInt("DB_CONN_MAX_LIFETIME")))
	if err != nil {
		return nil, err
	}

	config.Server = ServerConfig{
		Host:         viper.GetString("SERVER_HOST"),
		Port:         viper.GetInt("SERVER_PORT"),
		Mode:         viper.GetString("SERVER_MODE"),
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}

	config.Database = DatabaseConfig{
		Host:            viper.GetString("DB_HOST"),
		Port:            viper.GetInt("DB_PORT"),
		User:            viper.GetString("DB_USER"),
		Password:        viper.GetString("DB_PASSWORD"),
		Name:            viper.GetString("DB_NAME"),
		SSLMode:         viper.GetString("DB_SSL_MODE"),
		MaxOpenConns:    viper.GetInt("DB_MAX_OPEN_CONNS"),
		MaxIdleConns:    viper.GetInt("DB_MAX_IDLE_CONNS"),
		ConnMaxLifetime: connMaxLifetime,
	}

	config.Redis = RedisConfig{
		Host:     viper.GetString("REDIS_HOST"),
		Port:     viper.GetInt("REDIS_PORT"),
		Password: viper.GetString("REDIS_PASSWORD"),
		DB:       viper.GetInt("REDIS_DB"),
		PoolSize: viper.GetInt("REDIS_POOL_SIZE"),
	}

	config.JWT = JWTConfig{
		Secret:          viper.GetString("JWT_SECRET"),
		AccessTokenTTL:  accessTokenTTL,
		RefreshTokenTTL: refreshTokenTTL,
		Issuer:          viper.GetString("JWT_ISSUER"),
	}

	config.AI = AIConfig{
		ServiceURL:      viper.GetString("AI_SERVICE_URL"),
		Timeout:         aiTimeout,
		GeminiAPIKey:    viper.GetString("GEMINI_API_KEY"),
		OpenRouterAPIKey: viper.GetString("OPENROUTER_API_KEY"),
	}

	config.CodeExec = CodeExecutionConfig{
		ExecutorURL:   viper.GetString("CODE_EXECUTOR_URL"),
		Timeout:       codeExecTimeout,
		MemoryLimitMB: viper.GetInt("CODE_EXECUTION_MEMORY_LIMIT_MB"),
		CPULimit:      viper.GetFloat64("CODE_EXECUTION_CPU_LIMIT"),
	}

	config.Email = EmailConfig{
		SMTPHost:     viper.GetString("SMTP_HOST"),
		SMTPPort:     viper.GetInt("SMTP_PORT"),
		SMTPUser:     viper.GetString("SMTP_USER"),
		SMTPPassword: viper.GetString("SMTP_PASSWORD"),
		SMTPFrom:     viper.GetString("SMTP_FROM"),
	}

	// Parse CORS settings
	corsOrigins := viper.GetString("CORS_ALLOWED_ORIGINS")
	corsMethods := viper.GetString("CORS_ALLOWED_METHODS")
	corsHeaders := viper.GetString("CORS_ALLOWED_HEADERS")

	config.Security = SecurityConfig{
		CORSAllowedOrigins: parseCSV(corsOrigins),
		CORSAllowedMethods: parseCSV(corsMethods),
		CORSAllowedHeaders: parseCSV(corsHeaders),
		RateLimitEnabled:   viper.GetBool("RATE_LIMIT_ENABLED"),
		RateLimitRequests:  viper.GetInt("RATE_LIMIT_REQUESTS_PER_MINUTE"),
	}

	config.Monitoring = MonitoringConfig{
		MetricsEnabled:     viper.GetBool("METRICS_ENABLED"),
		MetricsPort:        viper.GetInt("METRICS_PORT"),
		HealthCheckEnabled: viper.GetBool("HEALTH_CHECK_ENABLED"),
		TracingEnabled:     viper.GetBool("TRACING_ENABLED"),
	}

	return &config, nil
}

func parseCSV(input string) []string {
	if input == "" {
		return []string{}
	}
	var result []string
	// Simple CSV parsing (can be enhanced)
	for i, s := range splitCSV(input) {
		if s != "" {
			result = append(result, s)
		}
	}
	return result
}

func splitCSV(input string) []string {
	// Simple split by comma, trim spaces
	var result []string
	start := 0
	inQuote := false
	for i, char := range input {
		switch char {
		case '"':
			inQuote = !inQuote
		case ',':
			if !inQuote {
				result = append(result, trimCSVField(input[start:i]))
				start = i + 1
			}
		}
	}
	if start <= len(input) {
		result = append(result, trimCSVField(input[start:]))
	}
	return result
}

func trimCSVField(field string) string {
	field = trimSpaces(field)
	if len(field) > 1 && field[0] == '"' && field[len(field)-1] == '"' {
		field = field[1 : len(field)-1]
	}
	return trimSpaces(field)
}

func trimSpaces(s string) string {
	start := 0
	end := len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t' || s[start] == '\n') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n') {
		end--
	}
	return s[start:end]
}
