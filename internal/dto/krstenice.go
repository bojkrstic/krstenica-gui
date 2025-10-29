package dto

import (
	"time"
)

type Krstenica struct {
	ID                     int64     `json:"id"`
	Book                   string    `json:"book"`
	Page                   int64     `json:"page"`
	CurrentNumber          int64     `json:"current_number"`
	EparhijaId             *int64    `json:"eparhija_id"`
	EparhijaName           string    `json:"eparhija_name"`
	TampleId               *int64    `json:"tample_id"`
	TampleName             string    `json:"tample_name"`
	TampleCity             string    `json:"tample_city"`
	ParentId               *int64    `json:"parent_id"`
	ParentFirstName        string    `json:"parent_first_name"`
	ParentLastName         string    `json:"parent_last_name"`
	ParentOccupation       string    `json:"parent_occupation"`
	ParentCity             string    `json:"parent_city"`
	ParentReligion         string    `json:"parent_religion"`
	GodfatherId            *int64    `json:"godfather_id"`
	GodfatherFirstName     string    `json:"godfather_first_name"`
	GodfatherLastName      string    `json:"godfather_last_name"`
	GodfatherOccupation    string    `json:"godfather_occupation"`
	GodfatherCity          string    `json:"godfather_city"`
	GodfatherReligion      string    `json:"godfather_religion"`
	ParohId                *int64    `json:"paroh_id"`
	ParohFirstName         string    `json:"paroh_first_name"`
	ParohLastName          string    `json:"paroh_last_name"`
	PriestId               *int64    `json:"priest_id"`
	PriestFirstName        string    `json:"priest_first_name"`
	PriestLastName         string    `json:"priest_last_name"`
	FirstName              string    `json:"first_name"`
	LastName               string    `json:"last_name"`
	Gender                 string    `json:"gender"`
	City                   string    `json:"city"`
	Country                string    `json:"country"`
	BirthDate              time.Time `json:"birth_date"`
	BirthOrder             string    `json:"birth_order"`
	PlaceOfBirthday        string    `json:"place_of_birthday"`
	MunicipalityOfBirthday string    `json:"municipality_of_birthday"`
	Baptism                time.Time `json:"baptism"`
	IsChurchMarried        string    `json:"is_church_married"`
	IsTwin                 string    `json:"is_twin"`
	HasPhysicalDisability  bool      `json:"has_physical_disability"`
	Anagrafa               string    `json:"anagrafa"`
	NumberOfCertificate    int64     `json:"number_of_certificate"`
	TownOfCertificate      string    `json:"town_of_certificate"`
	Certificate            time.Time `json:"certificate"`
	Comment                string    `json:"comment"`
	Status                 string    `json:"status"`
	CreatedAt              time.Time `json:"created_at"`
}

type KrstenicaCreateReq struct {
	Book                   string    `json:"book" form:"book"`
	Page                   int64     `json:"page" form:"page"`
	CurrentNumber          int64     `json:"current_number" form:"current_number"`
	EparhijaId             int64     `json:"eparhija_id" form:"eparhija_id"`
	TampleId               int64     `json:"tample_id" form:"tample_id"`
	ParentId               int64     `json:"parent_id" form:"parent_id"`
	GodfatherId            int64     `json:"godfather_id" form:"godfather_id"`
	ParohId                *int64    `json:"paroh_id" form:"paroh_id"`
	PriestId               int64     `json:"priest_id" form:"priest_id"`
	FirstName              string    `json:"first_name" form:"first_name"`
	LastName               string    `json:"last_name" form:"last_name"`
	Gender                 string    `json:"gender" form:"gender"`
	City                   string    `json:"city" form:"city"`
	Country                string    `json:"country" form:"country"`
	BirthDate              time.Time `json:"birth_date" form:"birth_date" time_format:"2006-01-02T15:04:05Z07:00"`
	BirthOrder             string    `json:"birth_order" form:"birth_order"`
	PlaceOfBirthday        string    `json:"place_of_birthday" form:"place_of_birthday"`
	MunicipalityOfBirthday string    `json:"municipality_of_birthday" form:"municipality_of_birthday"`
	Baptism                time.Time `json:"baptism" form:"baptism" time_format:"2006-01-02T15:04:05Z07:00"`
	IsChurchMarried        string    `json:"is_church_married" form:"is_church_married"`
	IsTwin                 string    `json:"is_twin" form:"is_twin"`
	HasPhysicalDisability  bool      `json:"has_physical_disability" form:"has_physical_disability"`
	Anagrafa               string    `json:"anagrafa" form:"anagrafa"`
	NumberOfCertificate    int64     `json:"number_of_certificate" form:"number_of_certificate"`
	TownOfCertificate      string    `json:"town_of_certificate" form:"town_of_certificate"`
	Certificate            time.Time `json:"certificate" form:"certificate" time_format:"2006-01-02T15:04:05Z07:00"`
	Comment                string    `json:"comment" form:"comment"`
}

type KrstenicaUpdateReq struct {
	Book                   *string    `json:"book" form:"book"`
	Page                   *int64     `json:"page" form:"page"`
	CurrentNumber          *int64     `json:"current_number" form:"current_number"`
	EparhijaId             *int64     `json:"eparhija_id" form:"eparhija_id"`
	TampleId               *int64     `json:"tample_id" form:"tample_id"`
	ParentId               *int64     `json:"parent_id" form:"parent_id"`
	GodfatherId            *int64     `json:"godfather_id" form:"godfather_id"`
	ParohId                *int64     `json:"paroh_id" form:"paroh_id"`
	PriestId               *int64     `json:"priest_id" form:"priest_id"`
	FirstName              *string    `json:"first_name" form:"first_name"`
	LastName               *string    `json:"last_name" form:"last_name"`
	Gender                 *string    `json:"gender" form:"gender"`
	City                   *string    `json:"city" form:"city"`
	Country                *string    `json:"country" form:"country"`
	BirthDate              *time.Time `json:"birth_date" form:"birth_date" time_format:"2006-01-02T15:04:05Z07:00"`
	BirthOrder             *string    `json:"birth_order" form:"birth_order"`
	PlaceOfBirthday        *string    `json:"place_of_birthday" form:"place_of_birthday"`
	MunicipalityOfBirthday *string    `json:"municipality_of_birthday" form:"municipality_of_birthday"`
	Baptism                *time.Time `json:"baptism" form:"baptism" time_format:"2006-01-02T15:04:05Z07:00"`
	IsChurchMarried        *string    `json:"is_church_married" form:"is_church_married"`
	IsTwin                 *string    `json:"is_twin" form:"is_twin"`
	HasPhysicalDisability  *bool      `json:"has_physical_disability" form:"has_physical_disability"`
	Anagrafa               *string    `json:"anagrafa" form:"anagrafa"`
	NumberOfCertificate    *int64     `json:"number_of_certificate" form:"number_of_certificate"`
	TownOfCertificate      *string    `json:"town_of_certificate" form:"town_of_certificate"`
	Certificate            *time.Time `json:"certificate" form:"certificate" time_format:"2006-01-02T15:04:05Z07:00"`
	Comment                *string    `json:"comment" form:"comment"`
	Status                 *string    `json:"status" form:"status"`
}
