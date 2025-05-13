package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Adoption struct {
	ID                 primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	AnimalID           primitive.ObjectID     `bson:"animal_id" json:"animal_id"`
	UserID             primitive.ObjectID     `bson:"user_id" json:"user_id"`
	Status             string                 `bson:"status" json:"status"` // pending | approved | rejected
	ApplicationData    map[string]interface{} `bson:"application_data" json:"application_data"`
	ContractDocumentID *primitive.ObjectID    `bson:"contract_document_id,omitempty" json:"contract_document_id,omitempty"`
}
