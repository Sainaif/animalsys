package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Schedule struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	EmployeeID primitive.ObjectID `bson:"employee_id" json:"employee_id"`
	ShiftDate  string             `bson:"shift_date" json:"shift_date"` // ISO date string
	ShiftTime  string             `bson:"shift_time" json:"shift_time"`
	Tasks      []string           `bson:"tasks" json:"tasks"`
	Status     string             `bson:"status" json:"status"` // assigned | swap_requested | absent_requested | swapped | absent_approved
}
