package bootstrap

import (
	"context"
	"fmt"

	"yourapp/internal/global"
	"yourapp/pkg/cache/redisx"
	"yourapp/pkg/logger"
	"yourapp/pkg/messaging/kafka"
	"yourapp/pkg/storage/elasticsearch"
	"yourapp/pkg/storage/mysql"
	"yourapp/pkg/storage/postgres"
)

// Start initializes and starts all application services
func Start(ctx context.Context) error {
	logger.Info("Starting application bootstrap...")

	// Initialize database connections
	if err := initDatabases(ctx); err != nil {
		return fmt.Errorf("failed to initialize databases: %w", err)
	}

	// Initialize cache
	if err := initCache(ctx); err != nil {
		return fmt.Errorf("failed to initialize cache: %w", err)
	}

	// Initialize Elasticsearch
	if err := initElasticsearch(ctx); err != nil {
		return fmt.Errorf("failed to initialize Elasticsearch: %w", err)
	}

	// Initialize Kafka
	if err := initKafka(ctx); err != nil {
		return fmt.Errorf("failed to initialize Kafka: %w", err)
	}

	logger.Info("Application bootstrap completed successfully")
	return nil
}

// initDatabases initializes database connections
func initDatabases(ctx context.Context) error {
	cfg := global.GetConfig()

	// Initialize MySQL
	if cfg.Database.MySQL.Enabled {
		if err := mysql.Init(ctx, cfg.Database.MySQL); err != nil {
			return fmt.Errorf("failed to initialize MySQL: %w", err)
		}
		logger.Info("MySQL connection initialized")
	}

	// Initialize PostgreSQL
	if cfg.Database.PostgreSQL.Enabled {
		if err := postgres.Init(ctx, cfg.Database.PostgreSQL); err != nil {
			return fmt.Errorf("failed to initialize PostgreSQL: %w", err)
		}
		logger.Info("PostgreSQL connection initialized")
	}

	return nil
}

// initCache initializes cache connections
func initCache(ctx context.Context) error {
	cfg := global.GetConfig()

	if cfg.Cache.Redis.Enabled {
		if err := redisx.Init(ctx, cfg.Cache.Redis); err != nil {
			return fmt.Errorf("failed to initialize Redis: %w", err)
		}
		logger.Info("Redis connection initialized")
	}

	return nil
}

// initElasticsearch initializes Elasticsearch connection
func initElasticsearch(ctx context.Context) error {
	cfg := global.GetConfig()

	if cfg.Elasticsearch.Enabled {
		if err := elasticsearch.Init(ctx, cfg.Elasticsearch); err != nil {
			return fmt.Errorf("failed to initialize Elasticsearch: %w", err)
		}
		logger.Info("Elasticsearch connection initialized")
	}

	return nil
}

// initKafka initializes Kafka connection
func initKafka(ctx context.Context) error {
	cfg := global.GetConfig()

	if cfg.Kafka.Enabled {
		if err := kafka.Init(ctx, cfg.Kafka); err != nil {
			return fmt.Errorf("failed to initialize Kafka: %w", err)
		}
		logger.Info("Kafka connection initialized")
	}

	return nil
}
