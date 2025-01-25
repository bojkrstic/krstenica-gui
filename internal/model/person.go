package model

import "database/sql"

type PersonStatus string

const (
	PersonStatusActive   PersonStatus = "active"
	PersonStatusDeleted  PersonStatus = "deleted"
	PersonStatusInactive PersonStatus = "inactive"
)

type Person struct {
	ID         int64        `gorm:"column:id"`
	FirstName  string       `gorm:"column:first_name"`
	LastName   string       `gorm:"column:last_name"`
	BriefName  string       `gorm:"column:brief_name"`
	Occupation string       `gorm:"column:occupation"`
	Religion   string       `gorm:"column:religion"`
	Address    string       `gorm:"column:address"`
	Country    string       `gorm:"column:country"`
	Role       string       `gorm:"column:role"`
	Status     string       `gorm:"column:status"`
	City       string       `gorm:"column:city"`
	BirthDate  sql.NullTime `gorm:"column:birth_date"`
	CreatedAt  sql.NullTime `gorm:"column:created_at"`
}

func (Person) TableName() string {
	return "persons"
}
