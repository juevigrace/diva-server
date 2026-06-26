package storage

import (
	"time"
)

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
