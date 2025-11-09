package monitoring

import (
	"context"
	"runtime"
	"time"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories"
	"github.com/sainaif/animalsys/backend/internal/infrastructure/database/mongodb"
	"go.mongodb.org/mongo-driver/bson"
)

var startTime = time.Now()

type MonitoringUseCase struct {
	db               *mongodb.Database
	userRepo         repositories.UserRepository
	animalRepo       repositories.AnimalRepository
	adoptionRepo     repositories.AdoptionRepository
	donationRepo     repositories.DonationRepository
	documentRepo     repositories.DocumentRepository
	auditLogRepo     repositories.AuditLogRepository
}

func NewMonitoringUseCase(
	db *mongodb.Database,
	userRepo repositories.UserRepository,
	animalRepo repositories.AnimalRepository,
	adoptionRepo repositories.AdoptionRepository,
	donationRepo repositories.DonationRepository,
	documentRepo repositories.DocumentRepository,
	auditLogRepo repositories.AuditLogRepository,
) *MonitoringUseCase {
	return &MonitoringUseCase{
		db:           db,
		userRepo:     userRepo,
		animalRepo:   animalRepo,
		adoptionRepo: adoptionRepo,
		donationRepo: donationRepo,
		documentRepo: documentRepo,
		auditLogRepo: auditLogRepo,
	}
}

// GetSystemHealth returns overall system health status
func (uc *MonitoringUseCase) GetSystemHealth(ctx context.Context) (*entities.SystemHealth, error) {
	health := &entities.SystemHealth{
		Timestamp:    time.Now(),
		Dependencies: make(map[string]entities.ServiceInfo),
		Uptime:       int64(time.Since(startTime).Seconds()),
	}

	// Check database health
	dbHealth := uc.checkDatabaseHealth(ctx)
	health.Database = dbHealth

	// Check storage health
	storageHealth := uc.checkStorageHealth()
	health.Storage = storageHealth

	// Check memory health
	memoryHealth := uc.checkMemoryHealth()
	health.Memory = memoryHealth

	// Check API health (placeholder - would need actual metrics collection)
	health.API = entities.APIHealth{
		Status:              entities.HealthStatusHealthy,
		AverageResponseTime: 50,
		RequestsPerMinute:   100,
		ErrorRate:           0.1,
	}

	// Determine overall status
	health.Status = uc.determineOverallHealth(dbHealth, storageHealth, memoryHealth)

	return health, nil
}

// GetUsageStatistics returns system usage statistics
func (uc *MonitoringUseCase) GetUsageStatistics(ctx context.Context) (*entities.UsageStatistics, error) {
	stats := &entities.UsageStatistics{
		Timestamp:    time.Now(),
		EntityCounts: make(map[string]int64),
	}

	// Get user statistics
	userFilter := repositories.UserFilter{Limit: 1, Offset: 0}
	_, totalUsers, _ := uc.userRepo.List(ctx, userFilter)
	stats.TotalUsers = totalUsers
	stats.EntityCounts["users"] = totalUsers

	// Get active users (users who logged in within last 30 days)
	// This is a simplified version - would need proper tracking
	stats.ActiveUsers = totalUsers / 2 // Placeholder

	// Get animal statistics
	animalFilter := repositories.AnimalFilter{Limit: 1, Offset: 0}
	_, totalAnimals, _ := uc.animalRepo.List(ctx, animalFilter)
	stats.TotalAnimals = totalAnimals
	stats.EntityCounts["animals"] = totalAnimals

	// Get adoption statistics
	adoptionFilter := repositories.AdoptionFilter{Limit: 1, Offset: 0}
	_, totalAdoptions, _ := uc.adoptionRepo.List(ctx, adoptionFilter)
	stats.TotalAdoptions = totalAdoptions
	stats.EntityCounts["adoptions"] = totalAdoptions

	// Get donation statistics
	donationFilter := repositories.DonationFilter{Limit: 1, Offset: 0}
	_, totalDonations, _ := uc.donationRepo.List(ctx, &donationFilter)
	stats.TotalDonations = totalDonations
	stats.EntityCounts["donations"] = totalDonations

	// Get document statistics
	documentFilter := repositories.DocumentFilter{Limit: 1, Offset: 0}
	_, totalDocuments, _ := uc.documentRepo.List(ctx, &documentFilter)
	stats.TotalDocuments = totalDocuments
	stats.EntityCounts["documents"] = totalDocuments

	// Get audit log statistics
	auditLogFilter := repositories.AuditLogFilter{Limit: 1, Offset: 0}
	_, totalAuditLogs, _ := uc.auditLogRepo.List(ctx, auditLogFilter)
	stats.TotalAuditLogs = totalAuditLogs
	stats.EntityCounts["audit_logs"] = totalAuditLogs

	// Storage statistics (placeholder)
	stats.StorageUsedBytes = 1024 * 1024 * 100 // 100MB placeholder

	// API call statistics (placeholder - would need actual tracking)
	stats.APICallsToday = 1000
	stats.APICallsThisWeek = 7000
	stats.APICallsThisMonth = 30000

	return stats, nil
}

// GetPerformanceMetrics returns system performance metrics
func (uc *MonitoringUseCase) GetPerformanceMetrics(ctx context.Context) (*entities.PerformanceMetrics, error) {
	metrics := &entities.PerformanceMetrics{
		Timestamp: time.Now(),
	}

	// API metrics (placeholder - would need actual metrics collection)
	metrics.APIMetrics = entities.APIPerformanceMetrics{
		TotalRequests:       10000,
		SuccessfulRequests:  9900,
		FailedRequests:      100,
		AverageResponseTime: 45,
		P50ResponseTime:     30,
		P95ResponseTime:     100,
		P99ResponseTime:     200,
		RequestsByEndpoint:  make(map[string]int64),
		ErrorsByType:        make(map[string]int64),
	}

	// Database metrics
	dbStats := uc.getDatabasePerformanceMetrics(ctx)
	metrics.DatabaseMetrics = dbStats

	// Cache metrics (placeholder)
	metrics.CacheMetrics = entities.CachePerformanceMetrics{
		Hits:       8000,
		Misses:     2000,
		HitRate:    80.0,
		TotalKeys:  500,
		MemoryUsed: 1024 * 1024 * 10, // 10MB
	}

	// Queue metrics (placeholder)
	metrics.QueueMetrics = entities.QueuePerformanceMetrics{
		TotalJobs:     1000,
		PendingJobs:   10,
		CompletedJobs: 980,
		FailedJobs:    10,
	}

	return metrics, nil
}

// GetDatabaseStatistics returns detailed database statistics
func (uc *MonitoringUseCase) GetDatabaseStatistics(ctx context.Context) (*entities.DatabaseStatistics, error) {
	stats := &entities.DatabaseStatistics{
		Timestamp:    time.Now(),
		DatabaseName: uc.db.DB.Name(),
		Collections:  []entities.CollectionStats{},
		Indexes:      []entities.IndexInfo{},
	}

	// Get list of collections
	collections, err := uc.db.DB.ListCollectionNames(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	var totalSize int64
	var totalDocs int64
	var totalIndexSize int64

	// Get stats for each collection
	for _, collName := range collections {
		coll := uc.db.Collection(collName)

		// Get document count
		count, _ := coll.EstimatedDocumentCount(ctx)

		// Get indexes
		indexes, _ := coll.Indexes().List(ctx)
		var indexCount int
		var indexList []entities.IndexInfo

		if indexes != nil {
			defer indexes.Close(ctx)
			for indexes.Next(ctx) {
				var idx bson.M
				if err := indexes.Decode(&idx); err == nil {
					indexCount++

					// Extract index keys
					keys := []string{}
					if keyDoc, ok := idx["key"].(bson.M); ok {
						for k := range keyDoc {
							keys = append(keys, k)
						}
					}

					indexInfo := entities.IndexInfo{
						Collection: collName,
						Name:       idx["name"].(string),
						Keys:       keys,
						Unique:     false,
						Size:       0, // Would need to calculate actual size
					}

					if unique, ok := idx["unique"].(bool); ok {
						indexInfo.Unique = unique
					}

					indexList = append(indexList, indexInfo)
				}
			}
		}

		collStats := entities.CollectionStats{
			Name:          collName,
			DocumentCount: count,
			Size:          0, // Would need actual collection size
			IndexCount:    indexCount,
			IndexSize:     0, // Would need actual index size
		}

		if count > 0 {
			collStats.AverageDocSize = collStats.Size / count
		}

		stats.Collections = append(stats.Collections, collStats)
		stats.Indexes = append(stats.Indexes, indexList...)

		totalDocs += count
		totalSize += collStats.Size
		totalIndexSize += collStats.IndexSize
	}

	stats.TotalDocuments = totalDocs
	stats.TotalSize = totalSize
	stats.IndexSize = totalIndexSize

	return stats, nil
}

// GetSystemConfiguration returns system configuration
func (uc *MonitoringUseCase) GetSystemConfiguration(ctx context.Context) (*entities.SystemConfiguration, error) {
	config := &entities.SystemConfiguration{
		Environment:     "production", // Would come from config
		Version:         "1.0.0",      // Would come from config
		GoVersion:       runtime.Version(),
		DatabaseType:    "MongoDB",
		DatabaseVersion: uc.getDatabaseVersion(ctx),
		Features:        make(map[string]bool),
		Settings:        make(map[string]string),
		StartTime:       startTime,
	}

	// Feature flags (would come from settings)
	config.Features["audit_logging"] = true
	config.Features["email_notifications"] = true
	config.Features["sms_notifications"] = false
	config.Features["file_uploads"] = true

	// Settings (sanitized - no secrets)
	config.Settings["max_upload_size"] = "10MB"
	config.Settings["session_timeout"] = "15m"
	config.Settings["audit_log_retention"] = "90d"

	return config, nil
}

// Helper methods

func (uc *MonitoringUseCase) checkDatabaseHealth(ctx context.Context) entities.DatabaseHealth {
	health := entities.DatabaseHealth{
		Status:    entities.HealthStatusHealthy,
		Connected: false,
	}

	start := time.Now()
	err := uc.db.Client.Ping(ctx, nil)
	responseTime := time.Since(start).Milliseconds()

	health.ResponseTime = responseTime
	health.Connected = (err == nil)

	if err != nil {
		health.Status = entities.HealthStatusUnhealthy
		return health
	}

	// Get collection count
	collections, _ := uc.db.DB.ListCollectionNames(ctx, bson.D{})
	health.CollectionCount = len(collections)

	// Response time thresholds
	if responseTime > 1000 {
		health.Status = entities.HealthStatusUnhealthy
	} else if responseTime > 500 {
		health.Status = entities.HealthStatusDegraded
	}

	return health
}

func (uc *MonitoringUseCase) checkStorageHealth() entities.StorageHealth {
	// Placeholder - would need actual filesystem stats
	total := int64(1024 * 1024 * 1024 * 100) // 100GB
	used := int64(1024 * 1024 * 1024 * 20)   // 20GB
	available := total - used

	health := entities.StorageHealth{
		Status:         entities.HealthStatusHealthy,
		TotalSpace:     total,
		UsedSpace:      used,
		AvailableSpace: available,
		UsagePercent:   float64(used) / float64(total) * 100,
	}

	if health.UsagePercent > 90 {
		health.Status = entities.HealthStatusUnhealthy
	} else if health.UsagePercent > 80 {
		health.Status = entities.HealthStatusDegraded
	}

	return health
}

func (uc *MonitoringUseCase) checkMemoryHealth() entities.MemoryHealth {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	health := entities.MemoryHealth{
		Status:     entities.HealthStatusHealthy,
		Allocated:  m.Alloc,
		TotalAlloc: m.TotalAlloc,
		Sys:        m.Sys,
		NumGC:      m.NumGC,
	}

	// Calculate usage percentage (simplified)
	if m.Sys > 0 {
		health.UsagePercent = float64(m.Alloc) / float64(m.Sys) * 100
	}

	if health.UsagePercent > 90 {
		health.Status = entities.HealthStatusUnhealthy
	} else if health.UsagePercent > 75 {
		health.Status = entities.HealthStatusDegraded
	}

	return health
}

func (uc *MonitoringUseCase) determineOverallHealth(
	db entities.DatabaseHealth,
	storage entities.StorageHealth,
	memory entities.MemoryHealth,
) entities.HealthStatus {
	// If any component is unhealthy, system is unhealthy
	if db.Status == entities.HealthStatusUnhealthy ||
		storage.Status == entities.HealthStatusUnhealthy ||
		memory.Status == entities.HealthStatusUnhealthy {
		return entities.HealthStatusUnhealthy
	}

	// If any component is degraded, system is degraded
	if db.Status == entities.HealthStatusDegraded ||
		storage.Status == entities.HealthStatusDegraded ||
		memory.Status == entities.HealthStatusDegraded {
		return entities.HealthStatusDegraded
	}

	return entities.HealthStatusHealthy
}

func (uc *MonitoringUseCase) getDatabasePerformanceMetrics(ctx context.Context) entities.DatabasePerformanceMetrics {
	// Placeholder - would need actual query metrics collection
	return entities.DatabasePerformanceMetrics{
		TotalQueries:        10000,
		AverageQueryTime:    20,
		SlowQueries:         10,
		ActiveConnections:   5,
		QueriesByCollection: make(map[string]int64),
	}
}

func (uc *MonitoringUseCase) getDatabaseVersion(ctx context.Context) string {
	var result bson.M
	err := uc.db.DB.RunCommand(ctx, bson.D{{Key: "buildInfo", Value: 1}}).Decode(&result)
	if err != nil {
		return "unknown"
	}

	if version, ok := result["version"].(string); ok {
		return version
	}

	return "unknown"
}
