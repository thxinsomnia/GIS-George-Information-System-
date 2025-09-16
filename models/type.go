package models

type Type struct {
	TypeId   uint   `json:"type_id" gorm:"primaryKey"`
	TypeName string `json:"type_name" gorm:"type:varchar(100);not null"`
	Point    string `json:"points" gorm:"type:int;not null"`
	CreatedAt string `json:"created_at" gorm:"type:timestamp;not null"`
}

func (Type) TableName() string {
	return "event_type"}
