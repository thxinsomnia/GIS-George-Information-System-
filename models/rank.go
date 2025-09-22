package models

type Rank struct {
	Id         int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name       string `json:"name" gorm:"type:varchar(100)"`
	TotalValor int64  `json:"total_valor" gorm:"type:bigint"`
}