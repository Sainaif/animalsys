package entities

import (
	"time"
)

// SystemHealth represents the overall system health status
type SystemHealth struct {
	Status       HealthStatus           `json:"status"`
	Timestamp    time.Time              `json:"timestamp"`
	Database     DatabaseHealth         `json:"database"`
	API          APIHealth              `json:"api"`
	Storage      StorageHealth          `json:"storage"`
	Memory       MemoryHealth           `json:"memory"`
	Dependencies map[string]ServiceInfo `json:"dependencies"`
	Uptime       int64                  `json:"uptime_seconds"`
}

// HealthStatus represents the health status
type HealthStatus string

const (
	HealthStatusHealthy   HealthStatus = "healthy"
	HealthStatusDegraded  HealthStatus = "degraded"
	HealthStatusUnhealthy HealthStatus = "unhealthy"
)

// DatabaseHealth represents database health metrics
type DatabaseHealth struct {
	Status         HealthStatus `json:"status"`
	Connected      bool         `json:"connected"`
	ResponseTime   int64        `json:"response_time_ms"`
	ActiveConns    int32        `json:"active_connections"`
	CollectionCount int          `json:"collection_count"`
	DatabaseSize   int64        `json:"database_size_bytes"`
}

// APIHealth represents API health metrics
type APIHealth struct {
	Status              HealthStatus `json:"status"`
	AverageResponseTime int64        `json:"avg_response_time_ms"`
	RequestsPerMinute   float64      `json:"requests_per_minute"`
	ErrorRate           float64      `json:"error_rate_percent"`
}

// StorageHealth represents storage health metrics
type StorageHealth struct {
	Status         HealthStatus `json:"status"`
	TotalSpace     int64        `json:"total_space_bytes"`
	UsedSpace      int64        `json:"used_space_bytes"`
	AvailableSpace int64        `json:"available_space_bytes"`
	UsagePercent   float64      `json:"usage_percent"`
}

// MemoryHealth represents memory usage metrics
type MemoryHealth struct {
	Status       HealthStatus `json:"status"`
	Allocated    uint64       `json:"allocated_bytes"`
	TotalAlloc   uint64       `json:"total_alloc_bytes"`
	Sys          uint64       `json:"sys_bytes"`
	NumGC        uint32       `json:"num_gc"`
	UsagePercent float64      `json:"usage_percent"`
}

// ServiceInfo represents a dependency service status
type ServiceInfo struct {
	Name      string       `json:"name"`
	Status    HealthStatus `json:"status"`
	Available bool         `json:"available"`
	Message   string       `json:"message,omitempty"`
}

// UsageStatistics represents system usage statistics
type UsageStatistics struct {
	Timestamp          time.Time             `json:"timestamp"`
	TotalUsers         int64                 `json:"total_users"`
	ActiveUsers        int64                 `json:"active_users"`
	TotalAnimals       int64                 `json:"total_animals"`
	TotalAdoptions     int64                 `json:"total_adoptions"`
	TotalDonations     int64                 `json:"total_donations"`
	TotalDocuments     int64                 `json:"total_documents"`
	TotalAuditLogs     int64                 `json:"total_audit_logs"`
	StorageUsedBytes   int64                 `json:"storage_used_bytes"`
	APICallsToday      int64                 `json:"api_calls_today"`
	APICallsThisWeek   int64                 `json:"api_calls_this_week"`
	APICallsThisMonth  int64                 `json:"api_calls_this_month"`
	EntityCounts       map[string]int64      `json:"entity_counts"`
}

// PerformanceMetrics represents system performance metrics
type PerformanceMetrics struct {
	Timestamp           time.Time                  `json:"timestamp"`
	APIMetrics          APIPerformanceMetrics      `json:"api_metrics"`
	DatabaseMetrics     DatabasePerformanceMetrics `json:"database_metrics"`
	CacheMetrics        CachePerformanceMetrics    `json:"cache_metrics"`
	QueueMetrics        QueuePerformanceMetrics    `json:"queue_metrics"`
}

// APIPerformanceMetrics represents API performance metrics
type APIPerformanceMetrics struct {
	TotalRequests       int64              `json:"total_requests"`
	SuccessfulRequests  int64              `json:"successful_requests"`
	FailedRequests      int64              `json:"failed_requests"`
	AverageResponseTime int64              `json:"avg_response_time_ms"`
	P50ResponseTime     int64              `json:"p50_response_time_ms"`
	P95ResponseTime     int64              `json:"p95_response_time_ms"`
	P99ResponseTime     int64              `json:"p99_response_time_ms"`
	RequestsByEndpoint  map[string]int64   `json:"requests_by_endpoint"`
	ErrorsByType        map[string]int64   `json:"errors_by_type"`
}

// DatabasePerformanceMetrics represents database performance metrics
type DatabasePerformanceMetrics struct {
	TotalQueries        int64            `json:"total_queries"`
	AverageQueryTime    int64            `json:"avg_query_time_ms"`
	SlowQueries         int64            `json:"slow_queries"`
	ActiveConnections   int32            `json:"active_connections"`
	QueriesByCollection map[string]int64 `json:"queries_by_collection"`
}

// CachePerformanceMetrics represents cache performance metrics
type CachePerformanceMetrics struct {
	Hits        int64   `json:"hits"`
	Misses      int64   `json:"misses"`
	HitRate     float64 `json:"hit_rate_percent"`
	TotalKeys   int64   `json:"total_keys"`
	MemoryUsed  int64   `json:"memory_used_bytes"`
}

// QueuePerformanceMetrics represents queue/background job metrics
type QueuePerformanceMetrics struct {
	TotalJobs     int64 `json:"total_jobs"`
	PendingJobs   int64 `json:"pending_jobs"`
	CompletedJobs int64 `json:"completed_jobs"`
	FailedJobs    int64 `json:"failed_jobs"`
}

// DatabaseStatistics represents detailed database statistics
type DatabaseStatistics struct {
	Timestamp       time.Time                 `json:"timestamp"`
	DatabaseName    string                    `json:"database_name"`
	Collections     []CollectionStats         `json:"collections"`
	TotalSize       int64                     `json:"total_size_bytes"`
	TotalDocuments  int64                     `json:"total_documents"`
	IndexSize       int64                     `json:"index_size_bytes"`
	Indexes         []IndexInfo               `json:"indexes"`
}

// CollectionStats represents statistics for a collection
type CollectionStats struct {
	Name           string  `json:"name"`
	DocumentCount  int64   `json:"document_count"`
	Size           int64   `json:"size_bytes"`
	AverageDocSize int64   `json:"avg_doc_size_bytes"`
	IndexCount     int     `json:"index_count"`
	IndexSize      int64   `json:"index_size_bytes"`
}

// IndexInfo represents information about a database index
type IndexInfo struct {
	Collection string   `json:"collection"`
	Name       string   `json:"name"`
	Keys       []string `json:"keys"`
	Unique     bool     `json:"unique"`
	Size       int64    `json:"size_bytes"`
}

// SystemConfiguration represents the current system configuration
type SystemConfiguration struct {
	Environment      string            `json:"environment"`
	Version          string            `json:"version"`
	GoVersion        string            `json:"go_version"`
	DatabaseType     string            `json:"database_type"`
	DatabaseVersion  string            `json:"database_version"`
	Features         map[string]bool   `json:"features"`
	Settings         map[string]string `json:"settings"`
	StartTime        time.Time         `json:"start_time"`
}
