package model

import "database/sql"

type TampleStatus string

const (
	TampleStatusActive   TampleStatus = "active"
	TampleStatusDeleted  TampleStatus = "deleted"
	TampleStatusInactive TampleStatus = "inactive"
)

type Tample struct {
	ID        int64        `gorm:"column:id"`
	Name      string       `gorm:"column:name"`
	Status    TampleStatus `gorm:"column:status"`
	City      string       `gorm:"column:city"`
	CreatedAt sql.NullTime `gorm:"column:created_at"`
}

func (Tample) TableName() string {
	return "tamples"
}
