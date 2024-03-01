package models

import (
	"time"
)

type Patient struct {
	ID        uint      `gorm:"id;primaryKey"`
	Name      string    `gorm:"name;size:30;not null"`
	Age       int       `gorm:"age;not null"`
	Gender    int       `gorm:"gender;not null"`
	CreatedAt time.Time `gorm:"created_at;not null"`
	UpdatedAt time.Time `gorm:"updated_at;not null"`
}

func (Patient) TableName() string {
	return "patient"
}
