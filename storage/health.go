package storage

import (
	"context"
	"log"
	"time"
)

type HealthStatus struct {
	IsHealthy    bool
	LastCheck    time.Time
	FailureCount int
	LastError    error
}

type HealthResult struct {
	Status       string            `json:"status"`
	Message      string            `json:"message"`
	Timestamp    time.Time         `json:"timestamp"`
	ResponseTime time.Duration     `json:"response_time"`
	Database     DatabaseHealth    `json:"database"`
	Connection   ConnectionHealth  `json:"connection"`
	Metadata     map[string]string `json:"metadata,omitempty"`
}

type DatabaseHealth struct {
	IsConnected  bool          `json:"is_connected"`
	LastPing     time.Time     `json:"last_ping"`
	PingDuration time.Duration `json:"ping_duration"`
	ErrorMessage string        `json:"error_message,omitempty"`
}

type ConnectionHealth struct {
	TotalConns    int32 `json:"total_conns"`
	IdleConns     int32 `json:"idle_conns"`
	AcquiredConns int32 `json:"acquired_conns"`
	MaxConns      int32 `json:"max_conns"`
}

func (s *StorageS) Health(ctx context.Context) HealthResult {
	start := time.Now()
	result := HealthResult{
		Timestamp: start,
		Metadata:  make(map[string]string),
	}

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	// Ping the database
	pingStart := time.Now()
	err := s.pool.Ping(ctx)
	pingDuration := time.Since(pingStart)

	result.Database = DatabaseHealth{
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
	result.Connection = ConnectionHealth{
		TotalConns:    poolStats.TotalConns(),
		IdleConns:     poolStats.IdleConns(),
		AcquiredConns: poolStats.AcquiredConns(),
		MaxConns:      poolStats.MaxConns(),
	}

	// Add metadata
	result.Metadata["config_name"] = s.config.Name
	result.Metadata["driver"] = s.config.Driver

	return result
}
