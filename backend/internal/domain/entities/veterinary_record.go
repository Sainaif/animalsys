package entities

import (
	"time"
)

// VeterinaryRecord represents a combined view of a veterinary visit or vaccination.
type VeterinaryRecord struct {
	RecordType  string           `json:"record_type" bson:"record_type"`
	Visit       *VeterinaryVisit `json:"visit,omitempty" bson:"visit,omitempty"`
	Vaccination *Vaccination     `json:"vaccination,omitempty" bson:"vaccination,omitempty"`
	Date        time.Time        `json:"date" bson:"date"`
}
