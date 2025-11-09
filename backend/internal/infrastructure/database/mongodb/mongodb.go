package mongodb

import (
	"context"
	"fmt"

	"github.com/sainaif/animalsys/backend/internal/infrastructure/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Database wraps MongoDB database connection
type Database struct {
	Client *mongo.Client
	DB     *mongo.Database
}

// Connect establishes connection to MongoDB
func Connect(ctx context.Context, cfg config.DatabaseConfig) (*Database, error) {
	// Set client options
	clientOptions := options.Client().
		ApplyURI(cfg.URI).
		SetServerSelectionTimeout(cfg.Timeout)

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping the database to verify connection
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	db := client.Database(cfg.Database)

	return &Database{
		Client: client,
		DB:     db,
	}, nil
}

// Disconnect closes the database connection
func (d *Database) Disconnect(ctx context.Context) error {
	if d.Client != nil {
		return d.Client.Disconnect(ctx)
	}
	return nil
}

// Ping checks if database connection is alive
func (d *Database) Ping(ctx context.Context) error {
	return d.Client.Ping(ctx, readpref.Primary())
}

// Collection returns a MongoDB collection
func (d *Database) Collection(name string) *mongo.Collection {
	return d.DB.Collection(name)
}

// Collections holds all collection names
var Collections = struct {
	Users                 string
	Animals               string
	VeterinaryVisits      string
	Vaccinations          string
	Contacts              string
	Adoptions             string
	AdoptionApplications  string
	Donors                string
	Donations             string
	Campaigns             string
	Events                string
	EventAttendances      string
	Volunteers            string
	VolunteerHours        string
	VolunteerAssignments  string
	Partners              string
	Transfers             string
	Communications        string
	Templates             string
	Tasks                 string
	Documents             string
	Notifications         string
	Reports               string
	ReportExecutions      string
	Settings              string
	AuditLogs             string
	MedicalConditions     string
	Medications           string
	TreatmentPlans        string
}{
	Users:                "users",
	Animals:              "animals",
	VeterinaryVisits:     "veterinary_visits",
	Vaccinations:         "vaccinations",
	Contacts:             "contacts",
	Adoptions:            "adoptions",
	AdoptionApplications: "adoption_applications",
	Donors:               "donors",
	Donations:            "donations",
	Campaigns:            "campaigns",
	Events:               "events",
	EventAttendances:     "event_attendances",
	Volunteers:           "volunteers",
	VolunteerHours:       "volunteer_hours",
	VolunteerAssignments: "volunteer_assignments",
	Partners:             "partners",
	Transfers:            "transfers",
	Communications:       "communications",
	Templates:            "communication_templates",
	Tasks:                "tasks",
	Documents:            "documents",
	Notifications:        "notifications",
	Reports:              "reports",
	ReportExecutions:     "report_executions",
	Settings:             "foundation_settings",
	AuditLogs:            "audit_logs",
	MedicalConditions:    "medical_conditions",
	Medications:          "medications",
	TreatmentPlans:       "treatment_plans",
}
