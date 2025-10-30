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

type veterinaryVisitRepository struct {
	collection *mongo.Collection
}

func NewVeterinaryVisitRepository(db *mongo.Database) interfaces.VeterinaryVisitRepository {
	return &veterinaryVisitRepository{
		collection: db.Collection("veterinary_visits"),
	}
}

func (r *veterinaryVisitRepository) Create(ctx context.Context, visit *entities.VeterinaryVisit) error {
	visit.CreatedAt = time.Now()
	visit.UpdatedAt = time.Now()
	_, err := r.collection.InsertOne(ctx, visit)
	return err
}

func (r *veterinaryVisitRepository) GetByID(ctx context.Context, id string) (*entities.VeterinaryVisit, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid visit ID")
	}

	var visit entities.VeterinaryVisit
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&visit)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("visit not found")
		}
		return nil, err
	}
	return &visit, nil
}

func (r *veterinaryVisitRepository) Update(ctx context.Context, id string, visit *entities.VeterinaryVisit) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid visit ID")
	}

	visit.UpdatedAt = time.Now()
	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": visit})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("visit not found")
	}
	return nil
}

func (r *veterinaryVisitRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid visit ID")
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("visit not found")
	}
	return nil
}

func (r *veterinaryVisitRepository) GetByAnimalID(ctx context.Context, animalID string, limit, offset int) ([]*entities.VeterinaryVisit, int64, error) {
	filter := bson.M{"animal_id": animalID}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	findOptions := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(offset)).
		SetSort(bson.M{"visit_date": -1})

	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var visits []*entities.VeterinaryVisit
	if err = cursor.All(ctx, &visits); err != nil {
		return nil, 0, err
	}
	return visits, total, nil
}

func (r *veterinaryVisitRepository) GetUpcomingVisits(ctx context.Context, days int) ([]*entities.VeterinaryVisit, error) {
	futureDate := time.Now().AddDate(0, 0, days)

	filter := bson.M{
		"next_visit_date": bson.M{
			"$lte": futureDate,
			"$gte": time.Now(),
		},
	}

	cursor, err := r.collection.Find(ctx, filter, options.Find().SetSort(bson.M{"next_visit_date": 1}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var visits []*entities.VeterinaryVisit
	if err = cursor.All(ctx, &visits); err != nil {
		return nil, err
	}
	return visits, nil
}

func (r *veterinaryVisitRepository) GetByDateRange(ctx context.Context, startDate, endDate string) ([]*entities.VeterinaryVisit, error) {
	start, _ := time.Parse(time.RFC3339, startDate)
	end, _ := time.Parse(time.RFC3339, endDate)

	filter := bson.M{"visit_date": bson.M{"$gte": start, "$lte": end}}

	cursor, err := r.collection.Find(ctx, filter, options.Find().SetSort(bson.M{"visit_date": -1}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var visits []*entities.VeterinaryVisit
	if err = cursor.All(ctx, &visits); err != nil {
		return nil, err
	}
	return visits, nil
}

type vaccinationRepository struct {
	collection *mongo.Collection
}

func NewVaccinationRepository(db *mongo.Database) interfaces.VaccinationRepository {
	return &vaccinationRepository{
		collection: db.Collection("vaccinations"),
	}
}

func (r *vaccinationRepository) Create(ctx context.Context, vaccination *entities.Vaccination) error {
	vaccination.CreatedAt = time.Now()
	_, err := r.collection.InsertOne(ctx, vaccination)
	return err
}

func (r *vaccinationRepository) GetByID(ctx context.Context, id string) (*entities.Vaccination, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid vaccination ID")
	}

	var vaccination entities.Vaccination
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&vaccination)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("vaccination not found")
		}
		return nil, err
	}
	return &vaccination, nil
}

func (r *vaccinationRepository) Update(ctx context.Context, id string, vaccination *entities.Vaccination) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid vaccination ID")
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": vaccination})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("vaccination not found")
	}
	return nil
}

func (r *vaccinationRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid vaccination ID")
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("vaccination not found")
	}
	return nil
}

func (r *vaccinationRepository) GetByAnimalID(ctx context.Context, animalID string) ([]*entities.Vaccination, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"animal_id": animalID}, options.Find().SetSort(bson.M{"date_administered": -1}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var vaccinations []*entities.Vaccination
	if err = cursor.All(ctx, &vaccinations); err != nil {
		return nil, err
	}
	return vaccinations, nil
}

func (r *vaccinationRepository) GetDueSoon(ctx context.Context, days int) ([]*entities.Vaccination, error) {
	futureDate := time.Now().AddDate(0, 0, days)

	filter := bson.M{
		"next_due_date": bson.M{
			"$lte": futureDate,
			"$gte": time.Now(),
		},
	}

	cursor, err := r.collection.Find(ctx, filter, options.Find().SetSort(bson.M{"next_due_date": 1}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var vaccinations []*entities.Vaccination
	if err = cursor.All(ctx, &vaccinations); err != nil {
		return nil, err
	}
	return vaccinations, nil
}
