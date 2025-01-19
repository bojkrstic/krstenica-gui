package model

import "database/sql"

type KrstenicaStatus string

const (
	KrstenicaStatusActive   KrstenicaStatus = "active"
	KrstenicaStatusDeleted  KrstenicaStatus = "deleted"
	KrstenicaStatusInactive KrstenicaStatus = "inactive"
)

type Krstenica struct {
	ID                     int64        `gorm:"column:id"`
	Book                   int64        `gorm:"column:book"`
	Page                   int64        `gorm:"column:page"`
	CurrentNumber          int64        `gorm:"column:current_number"`
	EparhijaId             int64        `gorm:"column:eparhija_id"`
	TampleId               int64        `gorm:"column:tample_id"`
	ParentId               int64        `gorm:"column:parent_id"`
	GodfatherId            int64        `gorm:"column:godfather_id"`
	ParohId                int64        `gorm:"column:paroh_id"`
	PriestId               int64        `gorm:"column:priest_id"`
	FirstName              string       `gorm:"column:first_name"`
	LastName               string       `gorm:"column:last_name"`
	Gender                 string       `gorm:"column:gender"`
	City                   string       `gorm:"column:city"`
	Country                string       `gorm:"column:country"`
	BirthDate              sql.NullTime `gorm:"column:birth_date"`
	BirthOrder             int64        `gorm:"column:birth_order"`
	PlaceOfBirthday        string       `gorm:"column:place_of_birthday"`
	MunicipalityOfBirthday string       `gorm:"column:municipality_of_birthday"`
	Baptism                sql.NullTime `gorm:"column:baptism"`
	IsChurchMarried        bool         `gorm:"column:is_church_married"`
	IsTwin                 bool         `gorm:"column:is_twin"`
	HasPhysicalDisability  bool         `gorm:"column:has_physical_disability"`
	Anagrafa               string       `gorm:"column:anagrafa"`
	NumberOfCertificate    int64        `gorm:"column:number_of_certificate"`
	TownOfCertificate      string       `gorm:"column:town_of_certificate"`
	Certificate            sql.NullTime `gorm:"column:certificate"`
	Comment                string       `gorm:"column:comment"`
	Status                 string       `gorm:"column:status"`
	CreatedAt              sql.NullTime `gorm:"column:created_at"`
}

func (Krstenica) TableName() string {
	return "krstenice"
}
