package postgres

import (
	"context"
	"log"
	"time"

	"github.com/juevigrace/diva-server/storage"
)

func (s *PGStorage) Health(ctx context.Context) storage.HealthResult {
	start := time.Now()
	result := storage.HealthResult{
		Timestamp: start,
		Metadata:  make(map[string]string),
	}

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	// Ping the database
	pingStart := time.Now()
	err := s.pool.Ping(ctx)
	pingDuration := time.Since(pingStart)

	result.Database = storage.DatabaseHealth{
		LastPing:     pingStart,
		PingDuration: pingDuration,
	}

	if err != nil {
		result.Status = "unhealthy"
		result.Message = "Database connection failed"
		result.Database.IsConnected = false
		result.Database.ErrorMessage = err.Error()
		result.ResponseTime = time.Since(start)

		log.Printf("Database health check failed: %v", err)
		return result
	}

	// Database is healthy
	result.Status = "healthy"
	result.Message = "Database is responsive"
	result.Database.IsConnected = true
	result.ResponseTime = time.Since(start)

	// Get connection pool statistics
	poolStats := s.pool.Stat()
	result.Connection = storage.ConnectionHealth{
		TotalConns:    poolStats.TotalConns(),
		IdleConns:     poolStats.IdleConns(),
		AcquiredConns: poolStats.AcquiredConns(),
		MaxConns:      poolStats.MaxConns(),
	}

	// Add metadata
	result.Metadata["config_name"] = s.config.Name
	result.Metadata["driver"] = "pgx"

	return result
}
