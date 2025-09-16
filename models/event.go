package models

import "time"

type Event struct {
	EventId     uint   `json:"event_id" gorm:"primaryKey"`
	EventName   string `json:"event_name" gorm:"type:varchar(100);not null"`
	EventTime       time.Time `json:"event_time" gorm:"type:timestamp;not null"`
	CreatedAt   string `json:"created_at" gorm:"type:timestamp;not null"`
	TypeId   uint `json:"type_id" gorm:"type:varchar(50);not null"`

	Type        Type      `gorm:"foreignKey:TypeId;references:TypeId"`
}