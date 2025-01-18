package dto

import (
	"time"
)

type Person struct {
	ID         int64     `json:"id"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	BriefName  string    `json:"brief_name"`
	Occupation string    `json:"occupation"`
	Religion   string    `json:"religion"`
	Address    string    `json:"address"`
	Country    string    `json:"country"`
	Role       string    `json:"role"`
	Status     string    `json:"status"`
	City       string    `json:"city"`
	CreatedAt  time.Time `json:"created_at"`
}

type PersonCreateReq struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	BriefName  string `json:"brief_name"`
	Occupation string `json:"occupation"`
	Religion   string `json:"religion"`
	Address    string `json:"address"`
	Country    string `json:"country"`
	Role       string `json:"role"`
	City       string `json:"city"`
}

type PersonUpdateReq struct {
	FirstName  *string `json:"first_name"`
	LastName   *string `json:"last_name"`
	BriefName  *string `json:"brief_name"`
	Occupation *string `json:"occupation"`
	Religion   *string `json:"religion"`
	Address    *string `json:"address"`
	Country    *string `json:"country"`
	Role       *string `json:"role"`
	City       *string `json:"city"`
	Status     *string `json:"status"`
}
