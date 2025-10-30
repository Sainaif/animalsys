package repository

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/sainaif/animalsys/internal/core/entities"
	"github.com/sainaif/animalsys/internal/core/interfaces"
)

type financeRepository struct {
	collection *mongo.Collection
}

func NewFinanceRepository(db *mongo.Database) interfaces.FinanceRepository {
	return &financeRepository{
		collection: db.Collection("finances"),
	}
}

func (r *financeRepository) Create(ctx context.Context, finance *entities.Finance) error {
	finance.CreatedAt = time.Now()
	finance.UpdatedAt = time.Now()
	_, err := r.collection.InsertOne(ctx, finance)
	return err
}

func (r *financeRepository) GetByID(ctx context.Context, id string) (*entities.Finance, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid finance ID")
	}

	var finance entities.Finance
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&finance)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("finance record not found")
		}
		return nil, err
	}
	return &finance, nil
}

func (r *financeRepository) Update(ctx context.Context, id string, finance *entities.Finance) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid finance ID")
	}

	finance.UpdatedAt = time.Now()
	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": finance})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("finance record not found")
	}
	return nil
}

func (r *financeRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid finance ID")
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("finance record not found")
	}
	return nil
}

func (r *financeRepository) List(ctx context.Context, filter *entities.FinanceFilter) ([]*entities.Finance, int64, error) {
	mongoFilter := bson.M{}

	if filter.Type != "" {
		mongoFilter["type"] = filter.Type
	}
	if filter.Category != "" {
		mongoFilter["category"] = filter.Category
	}
	if filter.FiscalYear > 0 {
		mongoFilter["fiscal_year"] = filter.FiscalYear
	}
	if filter.Quarter > 0 {
		mongoFilter["quarter"] = filter.Quarter
	}
	if !filter.StartDate.IsZero() || !filter.EndDate.IsZero() {
		dateFilter := bson.M{}
		if !filter.StartDate.IsZero() {
			dateFilter["$gte"] = filter.StartDate
		}
		if !filter.EndDate.IsZero() {
			dateFilter["$lte"] = filter.EndDate
		}
		mongoFilter["date"] = dateFilter
	}
	if filter.MinAmount > 0 || filter.MaxAmount > 0 {
		amountFilter := bson.M{}
		if filter.MinAmount > 0 {
			amountFilter["$gte"] = filter.MinAmount
		}
		if filter.MaxAmount > 0 {
			amountFilter["$lte"] = filter.MaxAmount
		}
		mongoFilter["amount"] = amountFilter
	}
	if filter.Search != "" {
		mongoFilter["$or"] = []bson.M{
			{"description": bson.M{"$regex": filter.Search, "$options": "i"}},
			{"reference_number": bson.M{"$regex": filter.Search, "$options": "i"}},
		}
	}

	total, err := r.collection.CountDocuments(ctx, mongoFilter)
	if err != nil {
		return nil, 0, err
	}

	findOptions := options.Find()
	if filter.Limit > 0 {
		findOptions.SetLimit(int64(filter.Limit))
	}
	if filter.Offset > 0 {
		findOptions.SetSkip(int64(filter.Offset))
	}
	sortField := "date"
	if filter.SortBy != "" {
		sortField = filter.SortBy
	}
	sortOrder := -1
	if filter.SortOrder == "asc" {
		sortOrder = 1
	}
	findOptions.SetSort(bson.M{sortField: sortOrder})

	cursor, err := r.collection.Find(ctx, mongoFilter, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var finances []*entities.Finance
	if err = cursor.All(ctx, &finances); err != nil {
		return nil, 0, err
	}
	return finances, total, nil
}

func (r *financeRepository) GetTotalByType(ctx context.Context, financeType entities.FinanceType, startDate, endDate string) (float64, error) {
	start, _ := time.Parse(time.RFC3339, startDate)
	end, _ := time.Parse(time.RFC3339, endDate)

	pipeline := mongo.Pipeline{
		{{"$match", bson.M{
			"type": financeType,
			"date": bson.M{"$gte": start, "$lte": end},
		}}},
		{{"$group", bson.M{
			"_id":   nil,
			"total": bson.M{"$sum": "$amount"},
		}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	var result []struct {
		Total float64 `bson:"total"`
	}
	if err = cursor.All(ctx, &result); err != nil {
		return 0, err
	}

	if len(result) == 0 {
		return 0, nil
	}
	return result[0].Total, nil
}

func (r *financeRepository) GetByCategory(ctx context.Context, category string, startDate, endDate string) ([]*entities.Finance, error) {
	start, _ := time.Parse(time.RFC3339, startDate)
	end, _ := time.Parse(time.RFC3339, endDate)

	filter := bson.M{
		"category": category,
		"date":     bson.M{"$gte": start, "$lte": end},
	}

	cursor, err := r.collection.Find(ctx, filter, options.Find().SetSort(bson.M{"date": -1}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var finances []*entities.Finance
	if err = cursor.All(ctx, &finances); err != nil {
		return nil, err
	}
	return finances, nil
}

func (r *financeRepository) GetSummary(ctx context.Context, startDate, endDate string) (map[string]float64, error) {
	start, _ := time.Parse(time.RFC3339, startDate)
	end, _ := time.Parse(time.RFC3339, endDate)

	pipeline := mongo.Pipeline{
		{{"$match", bson.M{
			"date": bson.M{"$gte": start, "$lte": end},
		}}},
		{{"$group", bson.M{
			"_id":   "$type",
			"total": bson.M{"$sum": "$amount"},
		}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []struct {
		Type  string  `bson:"_id"`
		Total float64 `bson:"total"`
	}
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	summary := make(map[string]float64)
	for _, r := range results {
		summary[r.Type] = r.Total
	}

	return summary, nil
}

func (r *financeRepository) GetByFiscalYear(ctx context.Context, year int) ([]*entities.Finance, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"fiscal_year": year}, options.Find().SetSort(bson.M{"date": 1}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var finances []*entities.Finance
	if err = cursor.All(ctx, &finances); err != nil {
		return nil, err
	}
	return finances, nil
}
