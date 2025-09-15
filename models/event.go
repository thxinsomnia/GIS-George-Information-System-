package models

type Event struct {
	EventId     uint   `json:"event_id" gorm:"primaryKey"`
	Name        string `json:"event_name" gorm:"type:varchar(100);not null"`
	Time       string `json:"event_time" gorm:"type:timestamp;not null"`
	CreatedAt   string `json:"created_at" gorm:"type:timestamp;not null"`
	Type       string `json:"event_type" gorm:"type:varchar(50);not null"`
}