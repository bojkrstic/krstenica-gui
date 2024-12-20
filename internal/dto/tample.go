package dto

import (
	"time"
)

type Tample struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	City      string    `json:"city"`
	CreatedAt time.Time `json:"created_at"`
}

type TampleCreateReq struct {
	Name string `json:"name"`
	City string `json:"city"`
}

type TampleUpdateReq struct {
	Name   *string `json:"name"`
	City   *string `json:"city"`
	Status *string `json:"status"`
}
