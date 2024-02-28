package models

import (
	"time"
)

type Patient struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"size:30;not null"`
	Age       int       `gorm:"not null"`
	Gender    int       `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

func (Patient) TableName() string {
	return "patient"
}
