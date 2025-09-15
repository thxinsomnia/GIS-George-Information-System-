package models

type Attendance struct {
	Id          int64   `json:"id" gorm:"primaryKey"`
	SoldierId   int64   `json:"soldier_id" gorm:"type:bigint;not null"`
	EventId     int64   `json:"event_id" gorm:"type:bigint;not null"`
	CreatedAt   string `json:"created_at" gorm:"type:timestamp;not null"`
}