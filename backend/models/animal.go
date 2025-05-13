package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Animal struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name          string             `bson:"name" json:"name"`
	Species       string             `bson:"species" json:"species"`
	Breed         string             `bson:"breed" json:"breed"`
	Age           int                `bson:"age" json:"age"`
	HealthHistory []string           `bson:"health_history" json:"health_history"`
	Status        string             `bson:"status" json:"status"` // available | adopted | deceased
}
