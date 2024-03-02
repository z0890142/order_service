package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	PatientID  int                `bson:"patient_id" json:"patient_id"`
	DoctorName string             `bson:"doctor_name" json:"doctor_name"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`

	Content string `bson:"content" json:"content"`
	Status  string `bson:"status" json:"status"`
}
