package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server      ServerConfig      `mapstructure:"server"`
	Database    DatabaseConfig    `mapstructure:"database"`
	Redis       RedisConfig       `mapstructure:"redis"`
	Security    SecurityConfig    `mapstructure:"security"`
	Monitoring  MonitoringConfig  `mapstructure:"monitoring"`
	AI          AIConfig          `mapstructure:"ai"`
}

type ServerConfig struct {
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	Mode         string        `mapstructure:"mode"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

type SecurityConfig struct {
	JWTPrivateKey        string        `mapstructure:"jwt_private_key"`
	JWTPublicKey         string        `mapstructure:"jwt_public_key"`
	JWTSecret            string        `mapstructure:"jwt_secret"`
	AccessTokenExpiry    time.Duration `mapstructure:"access_token_expiry"`
	RefreshTokenExpiry   time.Duration `mapstructure:"refresh_token_expiry"`
	CORSOrigins          []string      `mapstructure:"cors_origins"`
	RateLimit            int           `mapstructure:"rate_limit"`
	RateLimitTTL         int           `mapstructure:"rate_limit_ttl"`
	PasswordMinLength    int           `mapstructure:"password_min_length"`
	RequireEmailVerification bool      `mapstructure:"require_email_verification"`
}

type MonitoringConfig struct {
	MetricsEnabled bool `mapstructure:"metrics_enabled"`
	MetricsPort    int  `mapstructure:"metrics_port"`
}

type AIConfig struct {
	Provider    string  `mapstructure:"provider"`
	APIKey      string  `mapstructure:"api_key"`
	Model       string  `mapstructure:"model"`
	Temperature float64 `mapstructure:"temperature"`
	MaxTokens   int     `mapstructure:"max_tokens"`
}

func LoadConfig(path string) (*Config, error) {
	// Для .env файлов используем специальную обработку
	if strings.HasSuffix(path, ".env") {
		return loadEnvConfig(path)
	}
	
	viper.SetConfigFile(path)
	viper.AutomaticEnv()

	// Set defaults
	setDefaults()
	
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

func loadEnvConfig(path string) (*Config, error) {
	// Читаем .env файл вручную
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read .env file: %w", err)
	}
	
	lines := strings.Split(string(content), "\n")
	envMap := make(map[string]string)
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			envMap[key] = value
		}
	}
	
	// Создаем конфиг из переменных окружения
	config := &Config{
		Server: ServerConfig{
			Host:         getEnv(envMap, "SERVER_HOST", "localhost"),
			Port:         getEnvInt(envMap, "SERVER_PORT", 8080),
			Mode:         getEnv(envMap, "SERVER_MODE", "debug"),
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
			IdleTimeout:  120 * time.Second,
		},
		Database: DatabaseConfig{
			Host:     getEnv(envMap, "DATABASE_HOST", "localhost"),
			Port:     getEnvInt(envMap, "DATABASE_PORT", 5432),
			User:     getEnv(envMap, "DATABASE_USER", "postgres"),
			Password: getEnv(envMap, "DATABASE_PASSWORD", "postgres"),
			DBName:   getEnv(envMap, "DATABASE_DBNAME", "easyhire"),
			SSLMode:  getEnv(envMap, "DATABASE_SSLMODE", "disable"),
		},
		Redis: RedisConfig{
			Host:     getEnv(envMap, "REDIS_HOST", "localhost"),
			Port:     getEnvInt(envMap, "REDIS_PORT", 6379),
			Password: getEnv(envMap, "REDIS_PASSWORD", ""),
			DB:       getEnvInt(envMap, "REDIS_DB", 0),
			PoolSize: getEnvInt(envMap, "REDIS_POOL_SIZE", 10),
		},
		Security: SecurityConfig{
			JWTSecret:          getEnv(envMap, "SECURITY_JWT_SECRET", "test-secret"),
			AccessTokenExpiry:  15 * time.Minute,
			RefreshTokenExpiry: 7 * 24 * time.Hour,
			CORSOrigins:        []string{"*"},
		},
		Monitoring: MonitoringConfig{
			MetricsEnabled: true,
			MetricsPort:    9090,
		},
		AI: AIConfig{
			Provider:    getEnv(envMap, "AI_PROVIDER", "gemini"),
			APIKey:      getEnv(envMap, "AI_API_KEY", ""),
			Model:       getEnv(envMap, "AI_MODEL", "gemini-pro"),
			Temperature: 0.7,
			MaxTokens:   1000,
		},
	}
	
	return config, nil
}

func getEnv(envMap map[string]string, key, defaultValue string) string {
	if val, ok := envMap[key]; ok && val != "" {
		return val
	}
	return defaultValue
}

func getEnvInt(envMap map[string]string, key string, defaultValue int) int {
	val := getEnv(envMap, key, "")
	if val == "" {
		return defaultValue
	}
	
	var result int
	_, err := fmt.Sscanf(val, "%d", &result)
	if err != nil {
		return defaultValue
	}
	return result
}

func setDefaults() {
	viper.SetDefault("server.host", "localhost")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("server.read_timeout", 30*time.Second)
	viper.SetDefault("server.write_timeout", 30*time.Second)
	viper.SetDefault("server.idle_timeout", 120*time.Second)
	
	viper.SetDefault("database.sslmode", "disable")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.db", 0)
	viper.SetDefault("redis.pool_size", 10)
	
	viper.SetDefault("security.jwt_secret", "change-this-in-production")
	viper.SetDefault("security.access_token_expiry", 15*time.Minute)
	viper.SetDefault("security.refresh_token_expiry", 7*24*time.Hour)
	viper.SetDefault("security.cors_origins", []string{"*"})
	viper.SetDefault("security.rate_limit", 100)
	viper.SetDefault("security.rate_limit_ttl", 60)
	viper.SetDefault("security.password_min_length", 8)
	viper.SetDefault("security.require_email_verification", false)
	
	viper.SetDefault("monitoring.metrics_enabled", true)
	viper.SetDefault("monitoring.metrics_port", 9090)
	
	viper.SetDefault("ai.temperature", 0.7)
	viper.SetDefault("ai.max_tokens", 1000)
}
