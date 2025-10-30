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

type donorRepository struct {
	collection *mongo.Collection
}

func NewDonorRepository(db *mongo.Database) interfaces.DonorRepository {
	return &donorRepository{
		collection: db.Collection("donors"),
	}
}

func (r *donorRepository) Create(ctx context.Context, donor *entities.Donor) error {
	donor.CreatedAt = time.Now()
	donor.UpdatedAt = time.Now()
	_, err := r.collection.InsertOne(ctx, donor)
	return err
}

func (r *donorRepository) GetByID(ctx context.Context, id string) (*entities.Donor, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid donor ID")
	}

	var donor entities.Donor
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&donor)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("donor not found")
		}
		return nil, err
	}
	return &donor, nil
}

func (r *donorRepository) GetByEmail(ctx context.Context, email string) (*entities.Donor, error) {
	var donor entities.Donor
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&donor)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("donor not found")
		}
		return nil, err
	}
	return &donor, nil
}

func (r *donorRepository) Update(ctx context.Context, id string, donor *entities.Donor) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid donor ID")
	}

	donor.UpdatedAt = time.Now()
	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": donor})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("donor not found")
	}
	return nil
}

func (r *donorRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid donor ID")
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("donor not found")
	}
	return nil
}

func (r *donorRepository) List(ctx context.Context, donorType entities.DonorType, limit, offset int) ([]*entities.Donor, int64, error) {
	filter := bson.M{}
	if donorType != "" {
		filter["type"] = donorType
	}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	findOptions := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(offset)).
		SetSort(bson.M{"total_donated": -1})

	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var donors []*entities.Donor
	if err = cursor.All(ctx, &donors); err != nil {
		return nil, 0, err
	}
	return donors, total, nil
}

func (r *donorRepository) GetMajorDonors(ctx context.Context) ([]*entities.Donor, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"is_major_donor": true}, options.Find().SetSort(bson.M{"total_donated": -1}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var donors []*entities.Donor
	if err = cursor.All(ctx, &donors); err != nil {
		return nil, err
	}
	return donors, nil
}

func (r *donorRepository) UpdateTotals(ctx context.Context, id string, amount float64) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid donor ID")
	}

	now := time.Now()
	update := bson.M{
		"$inc": bson.M{"total_donated": amount, "donation_count": 1},
		"$set": bson.M{"last_donation": now, "updated_at": now},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("donor not found")
	}
	return nil
}

type donationRepository struct {
	collection *mongo.Collection
}

func NewDonationRepository(db *mongo.Database) interfaces.DonationRepository {
	return &donationRepository{
		collection: db.Collection("donations"),
	}
}

func (r *donationRepository) Create(ctx context.Context, donation *entities.Donation) error {
	donation.CreatedAt = time.Now()
	_, err := r.collection.InsertOne(ctx, donation)
	return err
}

func (r *donationRepository) GetByID(ctx context.Context, id string) (*entities.Donation, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid donation ID")
	}

	var donation entities.Donation
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&donation)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("donation not found")
		}
		return nil, err
	}
	return &donation, nil
}

func (r *donationRepository) Update(ctx context.Context, id string, donation *entities.Donation) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid donation ID")
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": donation})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("donation not found")
	}
	return nil
}

func (r *donationRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid donation ID")
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("donation not found")
	}
	return nil
}

func (r *donationRepository) GetByDonorID(ctx context.Context, donorID string, limit, offset int) ([]*entities.Donation, int64, error) {
	filter := bson.M{"donor_id": donorID}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	findOptions := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(offset)).
		SetSort(bson.M{"date": -1})

	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var donations []*entities.Donation
	if err = cursor.All(ctx, &donations); err != nil {
		return nil, 0, err
	}
	return donations, total, nil
}

func (r *donationRepository) GetByDateRange(ctx context.Context, startDate, endDate string) ([]*entities.Donation, error) {
	start, _ := time.Parse(time.RFC3339, startDate)
	end, _ := time.Parse(time.RFC3339, endDate)

	filter := bson.M{"date": bson.M{"$gte": start, "$lte": end}}

	cursor, err := r.collection.Find(ctx, filter, options.Find().SetSort(bson.M{"date": -1}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var donations []*entities.Donation
	if err = cursor.All(ctx, &donations); err != nil {
		return nil, err
	}
	return donations, nil
}

func (r *donationRepository) GetTotalByDonor(ctx context.Context, donorID string) (float64, error) {
	pipeline := mongo.Pipeline{
		{{"$match", bson.M{"donor_id": donorID}}},
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

func (r *donationRepository) GetRecurringDonations(ctx context.Context) ([]*entities.Donation, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"type": entities.DonationTypeRecurring})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var donations []*entities.Donation
	if err = cursor.All(ctx, &donations); err != nil {
		return nil, err
	}
	return donations, nil
}
