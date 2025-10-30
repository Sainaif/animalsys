package infrastructure

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Database holds the MongoDB client and database instance
type Database struct {
	Client *mongo.Client
	DB     *mongo.Database
}

// NewDatabase creates a new database connection
func NewDatabase(config *DatabaseConfig) (*Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Build connection URI if not provided
	uri := config.URI
	if uri == "" {
		if config.Username != "" && config.Password != "" {
			uri = fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?authSource=admin",
				config.Username,
				config.Password,
				config.Host,
				config.Port,
				config.Database,
			)
		} else {
			uri = fmt.Sprintf("mongodb://%s:%s/%s",
				config.Host,
				config.Port,
				config.Database,
			)
		}
	}

	// Set client options
	clientOptions := options.Client().
		ApplyURI(uri).
		SetMaxPoolSize(100).
		SetMinPoolSize(10).
		SetMaxConnIdleTime(30 * time.Second).
		SetServerSelectionTimeout(10 * time.Second)

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping the database to verify connection
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	// Get database instance
	database := client.Database(config.Database)

	return &Database{
		Client: client,
		DB:     database,
	}, nil
}

// Close closes the database connection
func (d *Database) Close(ctx context.Context) error {
	if d.Client == nil {
		return nil
	}

	return d.Client.Disconnect(ctx)
}

// HealthCheck checks if the database connection is alive
func (d *Database) HealthCheck(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return d.Client.Ping(ctx, readpref.Primary())
}

// Collection returns a collection from the database
func (d *Database) Collection(name string) *mongo.Collection {
	return d.DB.Collection(name)
}

// Collections returns all collection names
func (d *Database) Collections(ctx context.Context) ([]string, error) {
	return d.DB.ListCollectionNames(ctx, map[string]interface{}{})
}

// EnsureIndexes creates indexes for all collections
func (d *Database) EnsureIndexes(ctx context.Context) error {
	// Users collection indexes
	usersIndexes := []mongo.IndexModel{
		{
			Keys:    map[string]interface{}{"email": 1},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    map[string]interface{}{"username": 1},
			Options: options.Index().SetUnique(true),
		},
	}
	if _, err := d.Collection("users").Indexes().CreateMany(ctx, usersIndexes); err != nil {
		return fmt.Errorf("failed to create users indexes: %w", err)
	}

	// Animals collection indexes
	animalsIndexes := []mongo.IndexModel{
		{
			Keys: map[string]interface{}{"status": 1},
		},
		{
			Keys: map[string]interface{}{"species": 1},
		},
		{
			Keys: map[string]interface{}{"created_at": -1},
		},
	}
	if _, err := d.Collection("animals").Indexes().CreateMany(ctx, animalsIndexes); err != nil {
		return fmt.Errorf("failed to create animals indexes: %w", err)
	}

	// Adoptions collection indexes
	adoptionsIndexes := []mongo.IndexModel{
		{
			Keys: map[string]interface{}{"status": 1},
		},
		{
			Keys: map[string]interface{}{"animal_id": 1},
		},
		{
			Keys: map[string]interface{}{"applicant_id": 1},
		},
	}
	if _, err := d.Collection("adoptions").Indexes().CreateMany(ctx, adoptionsIndexes); err != nil {
		return fmt.Errorf("failed to create adoptions indexes: %w", err)
	}

	// Finances collection indexes
	financesIndexes := []mongo.IndexModel{
		{
			Keys: map[string]interface{}{"type": 1},
		},
		{
			Keys: map[string]interface{}{"category": 1},
		},
		{
			Keys: map[string]interface{}{"date": -1},
		},
	}
	if _, err := d.Collection("finances").Indexes().CreateMany(ctx, financesIndexes); err != nil {
		return fmt.Errorf("failed to create finances indexes: %w", err)
	}

	// Documents collection indexes
	documentsIndexes := []mongo.IndexModel{
		{
			Keys: map[string]interface{}{"type": 1},
		},
		{
			Keys: map[string]interface{}{"related_entity_type": 1, "related_entity_id": 1},
		},
	}
	if _, err := d.Collection("documents").Indexes().CreateMany(ctx, documentsIndexes); err != nil {
		return fmt.Errorf("failed to create documents indexes: %w", err)
	}

	// Inventory collection indexes
	inventoryIndexes := []mongo.IndexModel{
		{
			Keys: map[string]interface{}{"type": 1},
		},
		{
			Keys: map[string]interface{}{"stock_level": 1},
		},
	}
	if _, err := d.Collection("inventory").Indexes().CreateMany(ctx, inventoryIndexes); err != nil {
		return fmt.Errorf("failed to create inventory indexes: %w", err)
	}

	return nil
}
