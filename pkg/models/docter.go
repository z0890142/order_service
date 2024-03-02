package models

import "time"

type Doctor struct {
	ID        uint      `gorm:"id;primaryKey" json:"id"`
	Username  string    `gorm:"username;size:30;not null" json:"username"`
	Password  string    `gorm:"password;not null" json:"password"`
	CreatedAt time.Time `gorm:"created_at;not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"updated_at;not null" json:"updated_at"`
}

func (Doctor) TableName() string {
	return "doctor"
}
