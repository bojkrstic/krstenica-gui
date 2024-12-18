package model

import "database/sql"

type Tample struct {
	ID        int64        `gorm:"column:id"`
	Name      string       `gorm:"column:name"`
	Status    string       `gorm:"column:status"`
	City      string       `gorm:"column:city"`
	CreatedAt sql.NullTime `gorm:"column:created_at"`
}

func (Tample) TableName() string {
	return "tamples"
}
