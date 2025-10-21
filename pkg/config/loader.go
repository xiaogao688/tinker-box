package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config represents the application configuration
type Config struct {
	App           AppConfig           `mapstructure:"app"`
	Server        ServerConfig        `mapstructure:"server"`
	Database      DatabaseConfig      `mapstructure:"database"`
	Cache         CacheConfig         `mapstructure:"cache"`
	Elasticsearch ElasticsearchConfig `mapstructure:"elasticsearch"`
	Kafka         KafkaConfig         `mapstructure:"kafka"`
	Logging       LoggingConfig       `mapstructure:"logging"`
}

// AppConfig represents application configuration
type AppConfig struct {
	Name    string `mapstructure:"name"`
	Version string `mapstructure:"version"`
	Env     string `mapstructure:"env"`
}

// ServerConfig represents server configuration
type ServerConfig struct {
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

// DatabaseConfig represents database configuration
type DatabaseConfig struct {
	MySQL      MySQLConfig      `mapstructure:"mysql"`
	PostgreSQL PostgreSQLConfig `mapstructure:"postgres"`
}

// MySQLConfig represents MySQL configuration
type MySQLConfig struct {
	Enabled         bool          `mapstructure:"enabled"`
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	Username        string        `mapstructure:"username"`
	Password        string        `mapstructure:"password"`
	Database        string        `mapstructure:"database"`
	Charset         string        `mapstructure:"charset"`
	ParseTime       bool          `mapstructure:"parse_time"`
	Loc             string        `mapstructure:"loc"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

// PostgreSQLConfig represents PostgreSQL configuration
type PostgreSQLConfig struct {
	Enabled         bool          `mapstructure:"enabled"`
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	Username        string        `mapstructure:"username"`
	Password        string        `mapstructure:"password"`
	Database        string        `mapstructure:"database"`
	SSLMode         string        `mapstructure:"sslmode"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

// CacheConfig represents cache configuration
type CacheConfig struct {
	Redis RedisConfig `mapstructure:"redis"`
}

// RedisConfig represents Redis configuration
type RedisConfig struct {
	Enabled      bool          `mapstructure:"enabled"`
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	Password     string        `mapstructure:"password"`
	Database     int           `mapstructure:"database"`
	PoolSize     int           `mapstructure:"pool_size"`
	MinIdleConns int           `mapstructure:"min_idle_conns"`
	MaxConnAge   time.Duration `mapstructure:"max_conn_age"`
}

// ElasticsearchConfig represents Elasticsearch configuration
type ElasticsearchConfig struct {
	Enabled             bool          `mapstructure:"enabled"`
	Host                string        `mapstructure:"host"`
	Port                int           `mapstructure:"port"`
	Username            string        `mapstructure:"username"`
	Password            string        `mapstructure:"password"`
	MaxIdleConnsPerHost int           `mapstructure:"max_idle_conns_per_host"`
	Timeout             time.Duration `mapstructure:"timeout"`
}

// KafkaConfig represents Kafka configuration
type KafkaConfig struct {
	Enabled           bool          `mapstructure:"enabled"`
	Host              string        `mapstructure:"host"`
	Port              int           `mapstructure:"port"`
	Username          string        `mapstructure:"username"`
	Password          string        `mapstructure:"password"`
	SecurityProtocol  string        `mapstructure:"security_protocol"`
	SASLMechanism     string        `mapstructure:"sasl_mechanism"`
	SessionTimeout    time.Duration `mapstructure:"session_timeout"`
	HeartbeatInterval time.Duration `mapstructure:"heartbeat_interval"`
}

// LoggingConfig represents logging configuration
type LoggingConfig struct {
	Level    string `mapstructure:"level"`
	Format   string `mapstructure:"format"`
	Output   string `mapstructure:"output"`
	FilePath string `mapstructure:"file_path"`
}

// Load loads configuration using Viper
func Load() (*Config, error) {
	// Set default values
	setDefaults()

	// Configure Viper
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/yourapp")

	// Enable reading from environment variables
	viper.AutomaticEnv()
	viper.SetEnvPrefix("APP")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
		// Config file not found, use defaults and environment variables
	}

	// Unmarshal into struct
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

// setDefaults sets default configuration values
func setDefaults() {
	// App defaults
	viper.SetDefault("app.name", "yourapp")
	viper.SetDefault("app.version", "1.0.0")
	viper.SetDefault("app.env", "development")

	// Server defaults
	viper.SetDefault("server.host", "localhost")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.read_timeout", "30s")
	viper.SetDefault("server.write_timeout", "30s")

	// Database defaults
	viper.SetDefault("database.mysql.enabled", false)
	viper.SetDefault("database.mysql.host", "localhost")
	viper.SetDefault("database.mysql.port", 3306)
	viper.SetDefault("database.mysql.username", "root")
	viper.SetDefault("database.mysql.password", "")
	viper.SetDefault("database.mysql.database", "yourapp")
	viper.SetDefault("database.mysql.charset", "utf8mb4")
	viper.SetDefault("database.mysql.parse_time", true)
	viper.SetDefault("database.mysql.loc", "Local")
	viper.SetDefault("database.mysql.max_idle_conns", 10)
	viper.SetDefault("database.mysql.max_open_conns", 100)
	viper.SetDefault("database.mysql.conn_max_lifetime", "3600s")

	viper.SetDefault("database.postgres.enabled", false)
	viper.SetDefault("database.postgres.host", "localhost")
	viper.SetDefault("database.postgres.port", 5432)
	viper.SetDefault("database.postgres.username", "postgres")
	viper.SetDefault("database.postgres.password", "")
	viper.SetDefault("database.postgres.database", "yourapp")
	viper.SetDefault("database.postgres.sslmode", "disable")
	viper.SetDefault("database.postgres.max_idle_conns", 10)
	viper.SetDefault("database.postgres.max_open_conns", 100)
	viper.SetDefault("database.postgres.conn_max_lifetime", "3600s")

	// Cache defaults
	viper.SetDefault("cache.redis.enabled", false)
	viper.SetDefault("cache.redis.host", "localhost")
	viper.SetDefault("cache.redis.port", 6379)
	viper.SetDefault("cache.redis.password", "")
	viper.SetDefault("cache.redis.database", 0)
	viper.SetDefault("cache.redis.pool_size", 10)
	viper.SetDefault("cache.redis.min_idle_conns", 5)
	viper.SetDefault("cache.redis.max_conn_age", "3600s")

	// Elasticsearch defaults
	viper.SetDefault("elasticsearch.enabled", false)
	viper.SetDefault("elasticsearch.host", "localhost")
	viper.SetDefault("elasticsearch.port", 9200)
	viper.SetDefault("elasticsearch.username", "")
	viper.SetDefault("elasticsearch.password", "")
	viper.SetDefault("elasticsearch.max_idle_conns_per_host", 10)
	viper.SetDefault("elasticsearch.timeout", "30s")

	// Kafka defaults
	viper.SetDefault("kafka.enabled", false)
	viper.SetDefault("kafka.host", "localhost")
	viper.SetDefault("kafka.port", 9092)
	viper.SetDefault("kafka.username", "")
	viper.SetDefault("kafka.password", "")
	viper.SetDefault("kafka.security_protocol", "PLAINTEXT")
	viper.SetDefault("kafka.sasl_mechanism", "PLAIN")
	viper.SetDefault("kafka.session_timeout", "30s")
	viper.SetDefault("kafka.heartbeat_interval", "3s")

	// Logging defaults
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.format", "json")
	viper.SetDefault("logging.output", "stdout")
	viper.SetDefault("logging.file_path", "logs/app.log")
}

// GetString returns a string value from config
func GetString(key string) string {
	return viper.GetString(key)
}

// GetInt returns an int value from config
func GetInt(key string) int {
	return viper.GetInt(key)
}

// GetBool returns a bool value from config
func GetBool(key string) bool {
	return viper.GetBool(key)
}

// GetDuration returns a duration value from config
func GetDuration(key string) time.Duration {
	return viper.GetDuration(key)
}
