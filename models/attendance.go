package models

import "github.com/google/uuid"

type Attendance struct {
	Id          int64   `json:"id" gorm:"primaryKey"`
	SoldierId   uuid.UUID   `json:"soldier_id" gorm:"type:bigint;not null"`
	EventId     int64   `json:"event_id" gorm:"type:bigint;not null"`
	CreatedAt   string `json:"created_at" gorm:"type:timestamp;not null"`
	ValorEarned int64   `json:"valor_earned" gorm:"type:bigint;not null"`

	User  User  `gorm:"foreignKey:SoldierId;references:Id"`
    Event Event `gorm:"foreignKey:EventId;references:EventId"`
}