package repositories

import (
	"context"
	"time"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories"
	"github.com/sainaif/animalsys/backend/internal/infrastructure/database/mongodb"
	"github.com/sainaif/animalsys/backend/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type documentRepository struct {
	db *mongodb.Database
}

// NewDocumentRepository creates a new document repository
func NewDocumentRepository(db *mongodb.Database) repositories.DocumentRepository {
	return &documentRepository{db: db}
}

func (r *documentRepository) collection() *mongo.Collection {
	return r.db.DB.Collection("documents")
}

// Create creates a new document
func (r *documentRepository) Create(ctx context.Context, document *entities.Document) error {
	if document.ID.IsZero() {
		document.ID = primitive.NewObjectID()
	}

	_, err := r.collection().InsertOne(ctx, document)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return errors.ErrConflict
		}
		return err
	}

	return nil
}

// FindByID finds a document by ID
func (r *documentRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*entities.Document, error) {
	var document entities.Document
	err := r.collection().FindOne(ctx, bson.M{"_id": id}).Decode(&document)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.ErrNotFound
		}
		return nil, err
	}

	return &document, nil
}

// Update updates a document
func (r *documentRepository) Update(ctx context.Context, document *entities.Document) error {
	document.UpdatedAt = time.Now()

	result, err := r.collection().UpdateOne(
		ctx,
		bson.M{"_id": document.ID},
		bson.M{"$set": document},
	)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.ErrNotFound
	}

	return nil
}

// Delete deletes a document
func (r *documentRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	result, err := r.collection().DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.ErrNotFound
	}

	return nil
}

// List lists documents with filtering and pagination
func (r *documentRepository) List(ctx context.Context, filter *repositories.DocumentFilter) ([]*entities.Document, int64, error) {
	query := bson.M{}

	if !filter.IncludeArchived {
		query["is_archived"] = bson.M{"$ne": true}
	}

	// Access control
	if filter.UserID != nil {
		query["$or"] = []bson.M{
			{"is_public": true},
			{"uploaded_by": filter.UserID},
			{"accessible_by": filter.UserID},
		}
	}

	// Apply filters
	if filter.Type != "" {
		query["type"] = filter.Type
	}

	if filter.RelatedEntity != "" {
		query["related_entity"] = filter.RelatedEntity
	}

	if filter.RelatedEntityID != nil {
		query["related_entity_id"] = filter.RelatedEntityID
	}

	if filter.UploadedBy != nil {
		query["uploaded_by"] = filter.UploadedBy
	}

	if filter.IsPublic != nil {
		query["is_public"] = *filter.IsPublic
	}

	if filter.IsExpired != nil {
		now := time.Now()
		if *filter.IsExpired {
			query["expires_at"] = bson.M{"$lte": now}
		} else {
			query["$or"] = []bson.M{
				{"expires_at": nil},
				{"expires_at": bson.M{"$gt": now}},
			}
		}
	}

	if filter.Search != "" {
		query["$or"] = []bson.M{
			{"title": bson.M{"$regex": filter.Search, "$options": "i"}},
			{"file_name": bson.M{"$regex": filter.Search, "$options": "i"}},
			{"description": bson.M{"$regex": filter.Search, "$options": "i"}},
		}
	}

	if len(filter.Tags) > 0 {
		query["tags"] = bson.M{"$all": filter.Tags}
	}

	// Count total
	total, err := r.collection().CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	// Set up options
	opts := options.Find()

	// Sorting
	sortField := "created_at"
	if filter.SortBy != "" {
		sortField = filter.SortBy
	}
	sortOrder := -1 // descending by default
	if filter.SortOrder == "asc" {
		sortOrder = 1
	}
	opts.SetSort(bson.D{{Key: sortField, Value: sortOrder}})

	// Pagination
	if filter.Limit > 0 {
		opts.SetLimit(filter.Limit)
	}
	if filter.Offset > 0 {
		opts.SetSkip(filter.Offset)
	}

	// Execute query
	cursor, err := r.collection().Find(ctx, query, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var documents []*entities.Document
	if err = cursor.All(ctx, &documents); err != nil {
		return nil, 0, err
	}

	return documents, total, nil
}

// GetByRelatedEntity gets documents by related entity
func (r *documentRepository) GetByRelatedEntity(ctx context.Context, entityType string, entityID primitive.ObjectID) ([]*entities.Document, error) {
	query := bson.M{
		"related_entity":    entityType,
		"related_entity_id": entityID,
	}

	cursor, err := r.collection().Find(ctx, query, options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var documents []*entities.Document
	if err = cursor.All(ctx, &documents); err != nil {
		return nil, err
	}

	return documents, nil
}

// GetByType gets documents by type
func (r *documentRepository) GetByType(ctx context.Context, docType entities.DocumentType) ([]*entities.Document, error) {
	query := bson.M{"type": docType}

	cursor, err := r.collection().Find(ctx, query, options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var documents []*entities.Document
	if err = cursor.All(ctx, &documents); err != nil {
		return nil, err
	}

	return documents, nil
}

// GetPublicDocuments gets all public documents
func (r *documentRepository) GetPublicDocuments(ctx context.Context) ([]*entities.Document, error) {
	query := bson.M{"is_public": true}

	cursor, err := r.collection().Find(ctx, query, options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var documents []*entities.Document
	if err = cursor.All(ctx, &documents); err != nil {
		return nil, err
	}

	return documents, nil
}

// GetByUploader gets documents uploaded by a specific user
func (r *documentRepository) GetByUploader(ctx context.Context, userID primitive.ObjectID) ([]*entities.Document, error) {
	query := bson.M{"uploaded_by": userID}

	cursor, err := r.collection().Find(ctx, query, options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var documents []*entities.Document
	if err = cursor.All(ctx, &documents); err != nil {
		return nil, err
	}

	return documents, nil
}

// GetUserAccessibleDocuments gets documents accessible to a specific user
func (r *documentRepository) GetUserAccessibleDocuments(ctx context.Context, userID primitive.ObjectID) ([]*entities.Document, error) {
	query := bson.M{
		"$or": []bson.M{
			{"is_public": true},
			{"uploaded_by": userID},
			{"accessible_by": userID},
		},
	}

	cursor, err := r.collection().Find(ctx, query, options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var documents []*entities.Document
	if err = cursor.All(ctx, &documents); err != nil {
		return nil, err
	}

	return documents, nil
}

// GetExpiredDocuments gets all expired documents
func (r *documentRepository) GetExpiredDocuments(ctx context.Context) ([]*entities.Document, error) {
	now := time.Now()
	query := bson.M{
		"expires_at": bson.M{"$lte": now},
	}

	cursor, err := r.collection().Find(ctx, query, options.Find().SetSort(bson.D{{Key: "expires_at", Value: 1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var documents []*entities.Document
	if err = cursor.All(ctx, &documents); err != nil {
		return nil, err
	}

	return documents, nil
}

// GetExpiringSoonDocuments gets documents expiring soon (within 30 days)
func (r *documentRepository) GetExpiringSoonDocuments(ctx context.Context) ([]*entities.Document, error) {
	now := time.Now()
	thirtyDaysFromNow := now.AddDate(0, 0, 30)

	query := bson.M{
		"expires_at": bson.M{
			"$gt":  now,
			"$lte": thirtyDaysFromNow,
		},
	}

	cursor, err := r.collection().Find(ctx, query, options.Find().SetSort(bson.D{{Key: "expires_at", Value: 1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var documents []*entities.Document
	if err = cursor.All(ctx, &documents); err != nil {
		return nil, err
	}

	return documents, nil
}

// GetDocumentVersions gets all versions of a document
func (r *documentRepository) GetDocumentVersions(ctx context.Context, documentID primitive.ObjectID) ([]*entities.Document, error) {
	// Get the main document first
	mainDoc, err := r.FindByID(ctx, documentID)
	if err != nil {
		return nil, err
	}

	// Build query to find all versions
	var versions []*entities.Document
	versions = append(versions, mainDoc)

	// Traverse backward through previous versions
	currentDoc := mainDoc
	for currentDoc.PreviousVersion != nil {
		prevDoc, err := r.FindByID(ctx, *currentDoc.PreviousVersion)
		if err != nil {
			if err == errors.ErrNotFound {
				break // Stop if previous version not found
			}
			return nil, err
		}
		versions = append(versions, prevDoc)
		currentDoc = prevDoc
	}

	return versions, nil
}

// IncrementDownloadCount increments the download count for a document
func (r *documentRepository) IncrementDownloadCount(ctx context.Context, id primitive.ObjectID) error {
	now := time.Now()
	_, err := r.collection().UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{
			"$inc": bson.M{"download_count": 1},
			"$set": bson.M{"last_downloaded_at": now},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// GetDocumentStatistics gets document statistics
func (r *documentRepository) GetDocumentStatistics(ctx context.Context) (*repositories.DocumentStatistics, error) {
	stats := &repositories.DocumentStatistics{
		ByType: make(map[string]int64),
	}

	// Total documents
	total, err := r.collection().CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	stats.TotalDocuments = total

	// Public documents
	publicCount, err := r.collection().CountDocuments(ctx, bson.M{"is_public": true})
	if err != nil {
		return nil, err
	}
	stats.PublicDocuments = publicCount

	// Expired documents
	now := time.Now()
	expiredCount, err := r.collection().CountDocuments(ctx, bson.M{
		"expires_at": bson.M{"$lte": now},
	})
	if err != nil {
		return nil, err
	}
	stats.ExpiredDocuments = expiredCount

	// By type
	pipeline := mongo.Pipeline{
		{{Key: "$group", Value: bson.M{
			"_id":   "$type",
			"count": bson.M{"$sum": 1},
		}}},
	}

	cursor, err := r.collection().Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	type typeCount struct {
		Type  string `bson:"_id"`
		Count int64  `bson:"count"`
	}

	var typeCounts []typeCount
	if err = cursor.All(ctx, &typeCounts); err != nil {
		return nil, err
	}

	for _, tc := range typeCounts {
		stats.ByType[tc.Type] = tc.Count
	}

	// Private documents
	privateCount, err := r.collection().CountDocuments(ctx, bson.M{"is_public": false})
	if err != nil {
		return nil, err
	}
	stats.PrivateDocuments = privateCount

	// Total storage used (sum of file sizes) and downloads
	storagePipeline := mongo.Pipeline{
		{{Key: "$group", Value: bson.M{
			"_id":             nil,
			"total_size":      bson.M{"$sum": "$file_size"},
			"total_downloads": bson.M{"$sum": "$download_count"},
		}}},
	}

	var storageResult []struct {
		TotalSize      int64 `bson:"total_size"`
		TotalDownloads int64 `bson:"total_downloads"`
	}

	cursor, err = r.collection().Aggregate(ctx, storagePipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &storageResult); err != nil {
		return nil, err
	}

	if len(storageResult) > 0 {
		stats.TotalSize = storageResult[0].TotalSize
		stats.TotalDownloads = storageResult[0].TotalDownloads
	}

	return stats, nil
}

// EnsureIndexes creates the necessary indexes for the documents collection
func (r *documentRepository) EnsureIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "type", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "related_entity", Value: 1}, {Key: "related_entity_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "uploaded_by", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "is_public", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "expires_at", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "created_at", Value: -1}},
		},
		{
			Keys: bson.D{{Key: "tags", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "previous_version", Value: 1}},
		},
	}

	_, err := r.collection().Indexes().CreateMany(ctx, indexes)
	return err
}
