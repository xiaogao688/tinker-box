package options

import (
	"time"
)

// ServerOptions represents server configuration options
type ServerOptions struct {
	Host         string
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// DatabaseOptions represents database configuration options
type DatabaseOptions struct {
	Host            string
	Port            int
	Username        string
	Password        string
	Database        string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
}

// CacheOptions represents cache configuration options
type CacheOptions struct {
	Host         string
	Port         int
	Password     string
	Database     int
	PoolSize     int
	MinIdleConns int
	MaxConnAge   time.Duration
}

// LoggingOptions represents logging configuration options
type LoggingOptions struct {
	Level    string
	Format   string
	Output   string
	FilePath string
}

// DefaultServerOptions returns default server options
func DefaultServerOptions() *ServerOptions {
	return &ServerOptions{
		Host:         "localhost",
		Port:         8080,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
}

// DefaultDatabaseOptions returns default database options
func DefaultDatabaseOptions() *DatabaseOptions {
	return &DatabaseOptions{
		Host:            "localhost",
		Port:            3306,
		Username:        "root",
		Password:        "",
		Database:        "yourapp",
		MaxIdleConns:    10,
		MaxOpenConns:    100,
		ConnMaxLifetime: 3600 * time.Second,
	}
}

// DefaultCacheOptions returns default cache options
func DefaultCacheOptions() *CacheOptions {
	return &CacheOptions{
		Host:         "localhost",
		Port:         6379,
		Password:     "",
		Database:     0,
		PoolSize:     10,
		MinIdleConns: 5,
		MaxConnAge:   3600 * time.Second,
	}
}

// DefaultLoggingOptions returns default logging options
func DefaultLoggingOptions() *LoggingOptions {
	return &LoggingOptions{
		Level:    "info",
		Format:   "text",
		Output:   "stdout",
		FilePath: "logs/app.log",
	}
}
