package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	PatientID string             `bson:"patient_id"`
	DoctorID  string             `bson:"doctor_id"`
	CreatedAt time.Time          `bson:"created_at"`
	Content   string             `bson:"content"`
	Status    string             `bson:"status"`
}
