package dto

import (
	"time"
)

type Priest struct {
	ID        int64     `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	City      string    `json:"city"`
	Title     string    `json:"title"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type PriestCreateReq struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	City      string `json:"city"`
	Title     string `json:"title"`
}

type PriestUpdateReq struct {
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	City      *string `json:"city"`
	Title     *string `json:"title"`
	Status    *string `json:"status"`
}
