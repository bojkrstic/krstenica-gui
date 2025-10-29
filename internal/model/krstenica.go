package model

import (
	"database/sql"
)

type KrstenicaStatus string

const (
	KrstenicaStatusActive   KrstenicaStatus = "active"
	KrstenicaStatusDeleted  KrstenicaStatus = "deleted"
	KrstenicaStatusInactive KrstenicaStatus = "inactive"
)

type Krstenica struct {
	ID                  int64         `gorm:"column:id"`
	Book                string        `gorm:"column:book"`
	Page                int64         `gorm:"column:page"`
	CurrentNumber       int64         `gorm:"column:current_number"`
	EparhijaId          sql.NullInt64 `gorm:"column:eparhija_id"`
	EparhijaName        string        `gorm:"column:eparhija_name"`
	TampleId            sql.NullInt64 `gorm:"column:tample_id"`
	TampleName          string        `gorm:"column:tample_name"`
	TampleCity          string        `gorm:"column:tample_city"`
	ParentId            sql.NullInt64 `gorm:"column:parent_id"`
	ParentFirstName     string        `gorm:"column:parent_first_name"`
	ParentLastName      string        `gorm:"column:parent_last_name"`
	ParentOccupation    string        `gorm:"column:parent_occupation"`
	ParentCity          string        `gorm:"column:parent_city"`
	ParentReligion      string        `gorm:"column:parent_religion"`
	GodfatherId         sql.NullInt64 `gorm:"column:godfather_id"`
	GodfatherFirstName  string        `gorm:"column:godfather_first_name"`
	GodfatherLastName   string        `gorm:"column:godfather_last_name"`
	GodfatherOccupation string        `gorm:"column:godfather_occupation"`
	GodfatherCity       string        `gorm:"column:godfather_city"`
	GodfatherReligion   string        `gorm:"column:godfather_religion"`
	ParohId             sql.NullInt64 `gorm:"column:paroh_id"`
	ParohFirstName      string        `gorm:"column:paroh_first_name"`
	ParohLastName       string        `gorm:"column:paroh_last_name"`
	PriestId            sql.NullInt64 `gorm:"column:priest_id"`
	PriestFirstName     string        `gorm:"column:priest_first_name"`
	PriestLastName      string        `gorm:"column:priest_last_name"`
	FirstName           string        `gorm:"column:first_name"`
	LastName            string        `gorm:"column:last_name"`
	Gender              string        `gorm:"column:gender"`
	City                string        `gorm:"column:city"`
	Country             string        `gorm:"column:country"`
	BirthDate           sql.NullTime  `gorm:"column:birth_date"`
	// BirthDate              JSONDate     `gorm:"column:birth_date" json:"birth_date"`
	BirthOrder             string       `gorm:"column:birth_order"`
	PlaceOfBirthday        string       `gorm:"column:place_of_birthday"`
	MunicipalityOfBirthday string       `gorm:"column:municipality_of_birthday"`
	Baptism                sql.NullTime `gorm:"column:baptism"`
	IsChurchMarried        string       `gorm:"column:is_church_married"`
	IsTwin                 string       `gorm:"column:is_twin"`
	HasPhysicalDisability  string       `gorm:"column:has_physical_disability"`
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

type KrstenicaPost struct {
	ID            int64  `gorm:"column:id"`
	Book          string `gorm:"column:book"`
	Page          int64  `gorm:"column:page"`
	CurrentNumber int64  `gorm:"column:current_number"`
	EparhijaId    int64  `gorm:"column:eparhija_id"`
	//EparhijaName           string       `gorm:"column:eparhija_name"`
	TampleId int64 `gorm:"column:tample_id"`
	//TampleName             string       `gorm:"column:tample_name"`
	//TampleCity             string       `gorm:"column:tample_city"`
	ParentId int64 `gorm:"column:parent_id"`
	//ParentFirstName        string       `gorm:"column:parent_first_name"`
	//ParentLastName         string       `gorm:"column:parent_last_name"`
	//ParentOccupation       string       `gorm:"column:parent_occupation"`
	//ParentCity             string       `gorm:"column:parent_city"`
	//ParentReligion         string       `gorm:"column:parent_religion"`
	GodfatherId int64 `gorm:"column:godfather_id"`
	//GodfatherFirstName     string       `gorm:"column:godfather_first_name"`
	//GodfatherLastName      string       `gorm:"column:godfather_last_name"`
	//GodfatherOccupation    string       `gorm:"column:godfather_occupation"`
	//GodfatherCity          string       `gorm:"column:godfather_city"`
	//GodfatherReligion      string       `gorm:"column:godfather_religion"`
	ParohId *int64 `gorm:"column:paroh_id"`
	//ParohFirstName         string       `gorm:"column:paroh_first_name"`
	//ParohLastName          string       `gorm:"column:paroh_last_name"`
	PriestId int64 `gorm:"column:priest_id"`
	//PriestFirstName        string       `gorm:"column:priest_first_name"`
	//PriestLastName         string       `gorm:"column:priest_last_name"`
	FirstName string       `gorm:"column:first_name"`
	LastName  string       `gorm:"column:last_name"`
	Gender    string       `gorm:"column:gender"`
	City      string       `gorm:"column:city"`
	Country   string       `gorm:"column:country"`
	BirthDate sql.NullTime `gorm:"column:birth_date"`
	// BirthDate              JSONDate     `gorm:"column:birth_date" json:"birth_date"`
	BirthOrder             string       `gorm:"column:birth_order"`
	PlaceOfBirthday        string       `gorm:"column:place_of_birthday"`
	MunicipalityOfBirthday string       `gorm:"column:municipality_of_birthday"`
	Baptism                sql.NullTime `gorm:"column:baptism"`
	IsChurchMarried        string       `gorm:"column:is_church_married"`
	IsTwin                 string       `gorm:"column:is_twin"`
	HasPhysicalDisability  string       `gorm:"column:has_physical_disability"`
	Anagrafa               string       `gorm:"column:anagrafa"`
	NumberOfCertificate    int64        `gorm:"column:number_of_certificate"`
	TownOfCertificate      string       `gorm:"column:town_of_certificate"`
	Certificate            sql.NullTime `gorm:"column:certificate"`
	Comment                string       `gorm:"column:comment"`
	Status                 string       `gorm:"column:status"`
	CreatedAt              sql.NullTime `gorm:"column:created_at"`
}

func (KrstenicaPost) TableName() string {
	return "krstenice"
}
