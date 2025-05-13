package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Document struct {
	ID          primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	Filename    string              `bson:"filename" json:"filename"`
	UploadedBy  primitive.ObjectID  `bson:"uploaded_by" json:"uploaded_by"`
	UploadDate  string              `bson:"upload_date" json:"upload_date"`
	Type        string              `bson:"type" json:"type"` // medical, contract, other
	RelatedID   *primitive.ObjectID `bson:"related_id,omitempty" json:"related_id,omitempty"`
	ContentType string              `bson:"content_type" json:"content_type"`
	Size        int64               `bson:"size" json:"size"`
}
