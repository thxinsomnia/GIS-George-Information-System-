package models

import "github.com/google/uuid"


type User struct {
	Id        uuid.UUID `json:"id" gorm:"primaryKey"`
	Name      string     `json:"name" gorm:"type:varchar(100);not null"`
	Email     string     `json:"email" gorm:"type:varchar(100);not null"`
	Password  string     `json:"password" gorm:"type:varchar(100)"`
	Activity  int64      `json:"activitypoint" gorm:"type:int;not null"`
	CreatedAt string     `json:"created_at" gorm:"timestamp;not null"`
	Status    bool       `json:"status" gorm:"type:boolean;not null"`
	Role      string     `json:"role" gorm:"type:varchar(50);not null"`
	TotalValor int64      `json:"total_valor" gorm:"type:bigint;not null"`
	RankID      int64      `json:"-"`
    Rank        Rank      `json:"ranks" gorm:"foreignKey:RankID"`
}