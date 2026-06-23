package sqlite

import (
	"context"
	"log"
	"time"

	"github.com/juevigrace/diva-server/storage"
)

func (s *SQLiteStorage) Health(ctx context.Context) storage.HealthResult {
	start := time.Now()
	result := storage.HealthResult{
		Timestamp: start,
		Metadata:  make(map[string]string),
	}

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	pingStart := time.Now()
	err := s.db.PingContext(ctx)
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

	result.Status = "healthy"
	result.Message = "Database is responsive"
	result.Database.IsConnected = true
	result.ResponseTime = time.Since(start)

	stats := s.db.Stats()
	result.Connection = storage.ConnectionHealth{
		TotalConns:    int32(stats.OpenConnections),
		IdleConns:     int32(stats.Idle),
		AcquiredConns: int32(stats.InUse),
		MaxConns:      int32(stats.MaxOpenConnections),
	}

	result.Metadata["config_path"] = s.config.Path
	result.Metadata["driver"] = "sqlite3"

	return result
}
