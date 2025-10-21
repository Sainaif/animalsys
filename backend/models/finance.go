package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Finance struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Date        string             `bson:"date" json:"date"`
	Type        string             `bson:"type" json:"type"` // income | expense
	Amount      float64            `bson:"amount" json:"amount"`
	Description string             `bson:"description" json:"description"`
	Category    string             `bson:"category" json:"category"`
	CreatedBy   primitive.ObjectID `bson:"created_by" json:"created_by"`
}
