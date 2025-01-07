package model

import "database/sql"

type PriesteStatus string

const (
	PriestStatusActive   PriesteStatus = "active"
	PriestStatusDeleted  PriesteStatus = "deleted"
	PriestStatusInactive PriesteStatus = "inactive"
)

type Priest struct {
	ID        int64         `gorm:"column:id"`
	FirstName string        `gorm:"column:first_name"`
	LastName  string        `gorm:"column:last_name"`
	City      string        `gorm:"column:city"`
	Title     string        `gorm:"column:title"`
	Status    PriesteStatus `gorm:"column:status"`
	CreatedAt sql.NullTime  `gorm:"column:created_at"`
}

func (Priest) TableName() string {
	return "priests"
}
