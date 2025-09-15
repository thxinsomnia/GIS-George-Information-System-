package models

type Medal struct {
	Id        int64   `json:"id" gorm:"primaryKey"`
	Name      string `json:"medal_name" gorm:"type:varchar(100);not null"`
	Description string `json:"description" gorm:"type:text;not null"`
	Point     int64   `json:"points" gorm:"type:int;not null"`
	Url       string `json:"image_url" gorm:"type:text;not null"`
	CreatedAt string `json:"created_at" gorm:"type:timestamp;not null"`
}