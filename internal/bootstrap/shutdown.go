package bootstrap

import (
	"context"

	"yourapp/pkg/cache/redisx"
	"yourapp/pkg/logger"
	"yourapp/pkg/messaging/kafka"
	"yourapp/pkg/storage/elasticsearch"
	"yourapp/pkg/storage/mysql"
	"yourapp/pkg/storage/postgres"

	"go.uber.org/zap"
)

// Shutdown gracefully shuts down all application services
func Shutdown(ctx context.Context) error {
	logger.Info("Starting graceful shutdown...")

	// Close database connections
	if err := mysql.Close(); err != nil {
		logger.Error("Error closing MySQL connection", zap.Error(err))
	}

	if err := postgres.Close(); err != nil {
		logger.Error("Error closing PostgreSQL connection", zap.Error(err))
	}

	// Close cache connections
	if err := redisx.Close(); err != nil {
		logger.Error("Error closing Redis connection", zap.Error(err))
	}

	// Close Elasticsearch connections
	if err := elasticsearch.Close(); err != nil {
		logger.Error("Error closing Elasticsearch connection", zap.Error(err))
	}

	// Close Kafka connections
	if err := kafka.Close(); err != nil {
		logger.Error("Error closing Kafka connection", zap.Error(err))
	}

	logger.Info("Graceful shutdown completed")
	return nil
}
