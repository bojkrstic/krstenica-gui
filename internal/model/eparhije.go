package model

import "database/sql"

type EparhijeStatus string

const (
	EparhijeStatusActive   EparhijeStatus = "active"
	EparhijeStatusDeleted  EparhijeStatus = "deleted"
	EparhijeStatusInactive EparhijeStatus = "inactive"
)

type Eparhija struct {
	ID        int64          `gorm:"column:id"`
	Name      string         `gorm:"column:name"`
	Status    EparhijeStatus `gorm:"column:status"`
	City      string         `gorm:"column:city"`
	CreatedAt sql.NullTime   `gorm:"column:created_at"`
}

func (Eparhija) TableName() string {
	return "eparhije"
}
